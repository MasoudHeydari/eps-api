package delivery

import (
	"net/http"

	"github.com/MasoudHeydari/eps-api/db"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) GetAllSearchQueries(c echo.Context) error {
	dto := new(GetAllSearchQueries)
	if err := c.Bind(dto); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	sqs, err := db.GetAllSearchQueries(c.Request().Context(), s.db, dto.Page)
	if err != nil {
		logrus.Info("GetAllSearchQueries: %w", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, echo.Map{"search_queries": sqs})
}
