package main

import (
	"context"
	"errors"
	"io/ioutil"
	"reflect"
	"testing"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/m-lab/measure-saver/internal"
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

func (db *mockDB) Exec(interface{}, ...interface{}) (orm.Result, error) {
	return nil, nil
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

func Test_main(t *testing.T) {
	oldPgConnect := pgConnect
	pgConnect = func(opt *pg.Options) internal.Database {
		return &mockDB{}
	}

	_, cancel := context.WithCancel(context.Background())

	// Test without TLS.
	go main()
	time.Sleep(1 * time.Second)
	cancel()

	// Test with TLS.
	*flagTLSCert = "testdata/cert.pem"
	*flagTLSKey = "testdata/key.pem"
	*flagListenAddr = "127.0.0.1:0"
	go main()
	time.Sleep(1 * time.Second)
	cancel()

	pgConnect = oldPgConnect
}

func Test_initEchoServer(t *testing.T) {
	e := initEchoServer()
	if e == nil {
		t.Fatalf("initEchoServer() returned nil.")
	}

	flagKeysFile.Set("testdata/keys.txt")
	e = initEchoServer()
	if e == nil {
		t.Fatalf("initEchoServer() returned nil.")
	}
}
