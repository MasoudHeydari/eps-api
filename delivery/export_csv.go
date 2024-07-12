package delivery

import (
	"github.com/MasoudHeydari/eps-api/db"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Server) ExportCSV(c echo.Context) error {
	dto := new(ExportCSV)
	if err := c.Bind(dto); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if dto.SQID == 0 {
		return c.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}
	csvAbsPatch, fileName, err := db.ExportCSV(c.Request().Context(), s.db, dto.SQID)
	if err != nil {
		logrus.Info("ExportCSV: ", err)
		return c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	return c.Attachment(csvAbsPatch, fileName)
}
