package delivery

import (
	"github.com/MasoudHeydari/eps-api/db"
	"github.com/MasoudHeydari/eps-api/ent"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Server) CancelSearchQuery(c echo.Context) error {
	// curl -X PATCH -w "%{http_code}\n" http://localhost:9999/api/v1/search -H "Content-Type: application/json" -d '{"sq_id" : 1}'
	dto := new(CancelSQ)
	if err := c.Bind(dto); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if dto.SQID == 0 {
		return c.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}
	err := db.CancelSQ(c.Request().Context(), s.db, dto.SQID)
	if err != nil {
		logrus.Info("CancelSearchQuery.CancelSQ: %w", err)
		switch {
		case ent.IsNotFound(err):
			return c.JSON(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		default:
			return c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
	}
	return c.NoContent(http.StatusNoContent)
}
