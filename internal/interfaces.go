package internal

type Database interface {
	Close() error
	Insert(...interface{}) error
}
