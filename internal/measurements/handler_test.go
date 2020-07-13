package measurements

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-pg/pg/orm"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/m-lab/measure-saver/internal/model"
)

var (
	testMeasurement = model.Measurement{
		ID:         0,
		BrowserID:  "test",
		DeviceType: "device",
		Notes:      "notes",
		Download:   1000,
		Upload:     1000,
		Latency:    10,
		Results:    model.Results{},
	}
)

type mockDB struct {
	mustFail bool
	objects  []interface{}
}

func (db *mockDB) Insert(o ...interface{}) error {
	if db.mustFail {
		return errors.New("error while inserting entity")
	}

	db.objects = append(db.objects, o)
	return nil
}

func (db *mockDB) CreateTable(model interface{}, opt *orm.CreateTableOptions) error {
	return nil
}

func (db *mockDB) Close() error {
	return nil
}

func (db *mockDB) Exec(interface{}, ...interface{}) (orm.Result, error) {
	return nil, nil
}

func TestHandler_Post(t *testing.T) {
	db := &mockDB{}
	h := &Handler{
		DB: db,
	}
	e := echo.New()
	e.Validator = &Validator{
		Validator: validator.New(),
	}

	jsonBody, err := json.Marshal(testMeasurement)
	if err != nil {
		t.Fatalf("Cannot marshal test Measurement")
	}

	t.Run("ok", func(t *testing.T) {
		// Send test request.
		req := httptest.NewRequest(http.MethodPost, "/v0/measurements",
			strings.NewReader(string(jsonBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err = h.Post(c); err != nil {
			t.Fatalf("TestsHandler.Post() error = %v", err)
		}

		if rec.Code != 200 {
			t.Fatalf("TestsHandler.Post() code = %v, expected 200", rec.Code)
		}

		// Unmarshal result and verify it.
		var result model.Measurement
		json.Unmarshal(rec.Body.Bytes(), &result)

		if !reflect.DeepEqual(testMeasurement, result) {
			t.Errorf("TestsHandler.Post() expected: %v, got %v", testMeasurement, result)
		}
	})

	t.Run("error-body-is-not-json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v0/measurements",
			strings.NewReader("thisisnotjson"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err = h.Post(c); err == nil {
			t.Fatalf("TestsHandler.Post() expected error, got nil")
		}

		if err.(*echo.HTTPError).Code != http.StatusBadRequest {
			t.Errorf("TestsHandler.Post() status code = %v, expected %v",
				err.(*echo.HTTPError).Code, http.StatusBadRequest)
		}
	})

	t.Run("error-missing-required-fields", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v0/measurements",
			strings.NewReader(`{"BrowserID": "test"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err = h.Post(c); err == nil {
			t.Fatalf("TestsHandler.Post() expected error, got nil")
		}

		if err.(*echo.HTTPError).Code != http.StatusBadRequest {
			t.Errorf("TestsHandler.Post() status code = %v, expected %v",
				err.(*echo.HTTPError).Code, http.StatusBadRequest)
		}
	})

	t.Run("error-insert-failure", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v0/measurements",
			strings.NewReader(string(jsonBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db.mustFail = true
		if err = h.Post(c); err == nil {
			t.Fatalf("TestsHandler.Post() expected error, got nil")
		}

		if err.(*echo.HTTPError).Code != http.StatusInternalServerError {
			t.Errorf("TestsHandler.Post() status code = %v, expected %v",
				err.(*echo.HTTPError).Code, 400)
		}
		db.mustFail = false
	})
}
