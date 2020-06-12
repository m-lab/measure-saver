// measure-upload is a REST API which receives measurements from the
// M-Lab Measure Chrome Extension in JSON format and stores them into a
// PostgreSQL database.
package main

import (
	"errors"
	"testing"

	"github.com/go-pg/pg/orm"
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
