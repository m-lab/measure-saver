package main

import (
	"errors"
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
	got, err := readKeysFile("testdata/keys.txt")
	expected := []string{"foo", "bar"}
	if err != nil {
		t.Fatalf("readKeysFile() error = %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("readKeysFile() = %v, want %v", got, expected)
	}

	got, err = readKeysFile("testdata/thisdoesnotexist")
	if err == nil {
		t.Fatalf("readKeysFile() expected err, got nil")
	}
}

func Test_contains(t *testing.T) {
	testslice := []string{"foo", "bar"}
	tests := []struct {
		name  string
		slice []string
		el    string
		want  bool
	}{
		{
			name:  "element-exists",
			slice: testslice,
			el:    "foo",
			want:  true,
		},
		{
			name:  "element-does-not-exist",
			slice: testslice,
			el:    "baz",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.slice, tt.el); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
