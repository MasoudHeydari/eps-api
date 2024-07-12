package delivery

import (
	"github.com/MasoudHeydari/eps-api/db"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Server) GetAllSearchQueries(c echo.Context) error {
	sqs, err := db.GetAllSearchQueries(c.Request().Context(), s.db)
	if err != nil {
		logrus.Info("GetAllSearchQueries: %w", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, echo.Map{"search_queries": sqs})
}
