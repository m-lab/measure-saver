package main

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/go-pg/pg/orm"
	"github.com/labstack/echo"
)

type mockDB struct {
	mustFail bool
}

func (db *mockDB) Insert(o ...interface{}) error {
	return nil
}

func (db *mockDB) CreateTable(model interface{}, opt *orm.CreateTableOptions) error {
	if db.mustFail {
		return errors.New("failed to create table")
	}

	return nil
}

func (db *mockDB) Close() error {
	return nil
}

func Test_createSchema(t *testing.T) {
	err := createSchema(&mockDB{})
	if err != nil {
		t.Errorf("createSchema() error = %v", err)
	}

	err = createSchema(&mockDB{mustFail: true})
	if err == nil {
		t.Errorf("createSchema() expected error, got nil")
	}
}

func Test_apiVersion(t *testing.T) {
	e := echo.New()
	e.Pre(apiVersion)

	t.Run("redirect-to-versioned-path", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/endpoint", nil)
		rec := httptest.NewRecorder()

		req.Header.Set("version", "v0")
		req.URL.Path = "/endpoint"

		e.ServeHTTP(rec, req)

		if req.URL.Path != "/v0/endpoint" {
			t.Errorf("apiVersion() returned wrong path %s", req.URL.Path)
		}
	})

	t.Run("no-redirect-if-version-unspecified", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/endpoint", nil)
		rec := httptest.NewRecorder()

		req.URL.Path = "/endpoint"

		e.ServeHTTP(rec, req)

		if req.URL.Path != "/endpoint" {
			t.Errorf("apiVersion() returned wrong path %s", req.URL.Path)
		}
	})

}
