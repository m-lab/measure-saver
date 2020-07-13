package internal

import "github.com/go-pg/pg/orm"

type Database interface {
	Close() error
	CreateTable(model interface{}, opt *orm.CreateTableOptions) error
	Insert(...interface{}) error
	Exec(interface{}, ...interface{}) (orm.Result, error)
}

type HTTPServer interface {
}
