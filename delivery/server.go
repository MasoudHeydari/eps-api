package delivery

import (
	"fmt"
	"sync"

	"github.com/MasoudHeydari/eps-api/config"
	"github.com/MasoudHeydari/eps-api/db"
	"github.com/MasoudHeydari/eps-api/ent"
	"github.com/go-co-op/gocron/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type Server struct {
	db             *ent.Client
	agent          *Agent
	queryDepth     int
	limit          int
	counter        int
	fileNameMaxLen int
	mutex          sync.Mutex
}

func Start(cfgPath string) error {
	app, err := config.LoadAndConvert(cfgPath)
	if err != nil {
		return fmt.Errorf("start: %w", err)
	}
	logrus.Infof("Add Config: %+v", app)
	client, err := db.NewDB(app)
	if err != nil {
		return fmt.Errorf("failed to connect to DB, error: %v", err)
	}
	server := Server{
		db:             client,
		agent:          NewAgent(app.APIKey),
		queryDepth:     app.QueryDepth,
		counter:        0,
		fileNameMaxLen: app.FileNameMaxLen,
		limit:          app.Limiter.Burst,
	}
	err = server.setCron(app.Limiter.Hour, app.Limiter.Minute)
	if err != nil {
		return fmt.Errorf("start: failed to create the cron job: %w", err)
	}
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.POST("/api/v1/search", server.CreateJob)
	e.GET("/api/v1/search/:sq_id", server.GetSearchResults)
	e.PATCH("/api/v1/search", server.CancelSearchQuery)
	e.GET("/api/v1/export/:sq_id", server.ExportCSV)
	e.GET("/api/v1/search", server.GetAllSearchQueries)
	go server.PollJob()
	httpAddress := fmt.Sprintf("%s:%s", app.Http.Listen, app.Http.Port)
	return e.Start(httpAddress)
}

func (s *Server) setCron(hour, minute uint) error {
	// create a scheduler
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return fmt.Errorf("setCron: %w", err)
	}
	// add a job to the scheduler
	j, err := scheduler.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(
					hour,
					minute,
					00,
				),
			),
		),
		gocron.NewTask(
			func() {
				s.resetLimiterCounter()
			},
		),
	)
	if err != nil {
		return fmt.Errorf("setCron: %w", err)
	}
	scheduler.Start()
	logrus.Info("cron job started with ID: ", j.ID())
	return nil
}

func (s *Server) isAllowed() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.counter < s.limit
}

func (s *Server) increaseCounter() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.counter++
}

func (s *Server) resetLimiterCounter() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.counter = 0
	logrus.Info("rate limit reset")
}
