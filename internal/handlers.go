package internal

import (
	"fmt"
	"net/http"

	"github.com/go-pg/pg"
	"github.com/labstack/echo"
	"github.com/m-lab/measure-upload-service/internal/model"
)

type UploadHandler struct {
	DBConn *pg.DB
}

func (h *UploadHandler) PostUpload(c echo.Context) error {
	// Verify the provided appID is among the allowed ones.
	appID := c.QueryParam("appid")
	if appID == "" {
		return echo.NewHTTPError(http.StatusBadRequest,
			"parameter appid not provided.")
	}

	var m model.Measurement
	err := c.Bind(&m)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = h.DBConn.Insert(&m)
	if err != nil {
		fmt.Printf(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, m)
}
