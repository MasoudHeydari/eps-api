package delivery

import (
	"fmt"
	"time"

	"github.com/MasoudHeydari/eps-api/config"
	"github.com/MasoudHeydari/eps-api/db"
	"github.com/MasoudHeydari/eps-api/ent"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type Server struct {
	db         *ent.Client
	agent      *Agent
	limiter    *rate.Limiter
	queryDepth int
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
		db:         client,
		agent:      NewAgent(),
		queryDepth: app.QueryDepth,
		limiter: rate.NewLimiter(
			rate.Every(time.Duration(app.Limiter.IntervalMinutes)*time.Minute),
			app.Limiter.Burst,
		),
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
	return e.Start(":9999")
}
