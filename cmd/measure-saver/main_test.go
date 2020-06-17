package main

import (
	"errors"
	"io/ioutil"
	"reflect"
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

func Test_readKeysFile(t *testing.T) {
	testfile, err := ioutil.ReadFile("testdata/keys.txt")
	if err != nil {
		t.Fatalf("Cannot read test file")
	}
	got, err := readKeysFiles(testfile)
	expected := map[string]bool{
		"foo": true,
		"bar": true,
	}
	if err != nil {
		t.Fatalf("readKeysFile() error = %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("readKeysFile() = %v, want %v", got, expected)
	}
}
