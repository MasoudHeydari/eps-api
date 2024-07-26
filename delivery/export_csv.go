package delivery

import (
	"net/http"

	"github.com/MasoudHeydari/eps-api/db"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) ExportCSV(c echo.Context) error {
	dto := new(ExportCSV)
	if err := c.Bind(dto); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if dto.SQID == 0 {
		return c.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}
	csvAbsPatch, fileName, err := db.ExportCSV(c.Request().Context(), s.db, dto.SQID, s.fileNameMaxLen)
	if err != nil {
		logrus.Info("ExportCSV: ", err)
		return c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	return c.Attachment(csvAbsPatch, fileName)
}
