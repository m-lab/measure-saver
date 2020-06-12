package internal

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/m-lab/measure-upload-service/internal"
	"github.com/m-lab/measure-upload-service/internal/model"
)

// Handler is the handler for the /v0/measurements endpoint.
type Handler struct {
	DB internal.Database
}

// Validator is a custom validator wrapping go-playground/validator.
type Validator struct {
	Validator *validator.Validate
}

// Validate validates the passed struct according to its "validate"
// annotations.
func (cv *Validator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

// Post is the HTTP handler for POST /tests.
func (h *Handler) Post(c echo.Context) error {
	// TODO: Verify the provided appID is among the allowed ones.
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

	err = h.DB.Insert(&m)
	if err != nil {
		fmt.Printf(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, m)
}
