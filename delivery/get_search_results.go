package delivery

import (
	"github.com/MasoudHeydari/eps-api/db"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) GetSearchResults(c echo.Context) error {
	dto := new(GetAllSearchResults)
	if err := c.Bind(dto); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if dto.SQID == 0 {
		return c.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}
	serps, err := db.GetAllResult(c.Request().Context(), s.db, dto.SQID, dto.Page)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, serps)
}
