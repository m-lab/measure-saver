package internal

import (
	"fmt"
	"net/http"

	"github.com/go-pg/pg"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/m-lab/measure-upload-service/internal/model"
)

type UploadHandler struct {
	DBConn *pg.DB
}

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate object.
	if err = c.Validate(&m); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.DBConn.Insert(&m)
	if err != nil {
		fmt.Printf(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, m)
}
