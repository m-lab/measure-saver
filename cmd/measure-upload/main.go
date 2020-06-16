// measure-upload is a REST API which receives measurements from the
// M-Lab Measure Chrome Extension in JSON format and stores them into a
// PostgreSQL database.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/m-lab/go/flagx"
	"github.com/m-lab/go/rtx"
	"github.com/m-lab/measure-upload-service/internal"
	"github.com/m-lab/measure-upload-service/internal/measurements"
	"github.com/m-lab/measure-upload-service/internal/model"
)

const (
	// DefaultListenAddr is the default address to listen on for incoming
	// connections.
	DefaultListenAddr = ":1323"

	// DefaultDBAddr is the default address of the PostgreSQL database.
	DefaultDBAddr = ":5432"

	// DefaultDBUser is the default user for the PostgreSQL database.
	DefaultDBUser = "postgres"

	// DefaultDBPass is the default password for the PostgreSQL database.
	DefaultDBPass = "secret"

	// DefaultDBName is the default PostgreSQL database name.
	DefaultDBName = "measure-app"
)

var (
	flagListenAddr = flag.String("listenaddr", DefaultListenAddr,
		"Address to listen for incoming connection on.")
	flagDBAddr = flag.String("db.addr", DefaultDBAddr,
		"Address of the PostgreSQL database to use.")
	flagDBUser = flag.String("db.user", DefaultDBUser,
		"Username to use to connect to the database.")
	flagDBPass = flag.String("db.pass", DefaultDBPass,
		"Password to use to connect to the database.")
	flagDBName = flag.String("db.name", DefaultDBName,
		"Name of the database to use.")
	flagKeysFile = flag.String("keys.file", "",
		"Text file containing the allowed API keys.")
)

func main() {
	flag.Parse()
	rtx.Must(flagx.ArgsFromEnv(flag.CommandLine), "Could not parse env args")

	// Initialize database connection.
	db := pg.Connect(&pg.Options{
		Addr:     *flagDBAddr,
		User:     *flagDBUser,
		Password: *flagDBPass,
		Database: *flagDBName,
	})
	defer db.Close()

	// Create schema if needed.
	rtx.Must(createSchema(db), "Cannot create database schema")

	// Initialize the handler.
	measurementsHandler := measurements.Handler{
		DB: db,
	}

	// Initialize the Echo server.
	e := echo.New()
	e.Validator = &measurements.Validator{
		Validator: validator.New(),
	}

	if *flagKeysFile != "" {
		// Read API keys file and set up key auth middleware.
		allowedKeys, err := readKeysFile(*flagKeysFile)
		rtx.Must(err, "Cannot read keys file")

		e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
			KeyLookup: "query:key",
			Validator: func(key string, c echo.Context) (bool, error) {
				return contains(allowedKeys, key), nil
			},
		}))
	} else {
		fmt.Println("Warning: no keys file specified. The endpoints will be " +
			"unprotected.")
	}

	// Endpoints' routing.
	e.POST("/v0/measurements", measurementsHandler.Post)

	// Start the Echo server.
	e.Logger.Fatal(e.Start(*flagListenAddr))
}

// createSchema creates database schema for the Measurement model.
func createSchema(db internal.Database) error {
	err := db.CreateTable((*model.Measurement)(nil), &orm.CreateTableOptions{
		IfNotExists: true,
	})
	if err != nil {
		return err
	}

	return nil
}

// readKeysFile reads the keys file and returns a slice of strings, one per
// line.
func readKeysFile(path string) ([]string, error) {
	var keys []string

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		keys = append(keys, scanner.Text())
	}
	return keys, nil
}

func contains(slice []string, el string) bool {
	for _, e := range slice {
		if e == el {
			return true
		}
	}
	return false
}
