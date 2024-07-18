package delivery

import (
	"fmt"
	"github.com/MasoudHeydari/eps-api/db"
	"github.com/MasoudHeydari/eps-api/ent"
	"github.com/MasoudHeydari/eps-api/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Server) CreateJob(c echo.Context) error {
	//  curl -X POST -w "%{http_code}\n" http://localhost:9999/api/v1/search -H "Content-Type: application/json" -d '{"loc": "NL", "lang": "En", "q": "Golang"}'
	if !s.isAllowed() {
		logrus.Info("rate limit exceeded")
		return c.JSON(http.StatusTooManyRequests, echo.Map{"details": "rate limit exceeded"})
	}
	s.increaseCounter()
	sq := new(searchQ)
	if err := c.Bind(sq); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	sqID, err := db.InsertNewSearchQuery(c.Request().Context(), s.db, sq.LocCode, sq.Language, sq.Query)
	if err != nil {
		fmt.Println("CreateJob", err)
		switch {
		case ent.IsConstraintError(err):
			return c.JSON(http.StatusConflict, err)
		default:
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	query := model.Query{
		Text:     sq.Query,
		Location: sq.LocCode,
		LangCode: sq.Language,
		Depth:    s.queryDepth,
	}
	logrus.Infof("CreateJob: new query going to be started: %+v", query)
	jobID, err := s.agent.CreateJob(query)
	if err != nil {
		logrus.Debugf("server.CreateJob: %v\n", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	err = db.InsertJobID(c.Request().Context(), s.db, sqID, jobID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	logrus.Infof("server.CreateJob: job %q inserted into DB", jobID)
	return c.JSON(http.StatusOK, echo.Map{"sq_id": sqID})
}
