// measure-upload is a REST API which receives measurements from the
// M-Lab Measure Chrome Extension in JSON format and stores them into a
// PostgreSQL database.
package main

import (
	"flag"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/m-lab/go/flagx"
	"github.com/m-lab/go/rtx"
	"github.com/m-lab/measure-upload-service/internal"
	"github.com/m-lab/measure-upload-service/internal/model"
)

const (
	DefaultDBAddr = ":5432"
	DefaultDBUser = "postgres"
	DefaultDBPass = "secret"
	DefaultDBName = "measure-app"
)

var (
	flagDBAddr = flag.String("db.addr", DefaultDBAddr,
		"Address of the PostgreSQL database to use.")
	flagDBUser = flag.String("db.user", DefaultDBUser,
		"Username to use to connect to the database.")
	flagDBPass = flag.String("db.pass", DefaultDBPass,
		"Password to use to connect to the database.")
	flagDBName = flag.String("db.name", DefaultDBName,
		"Name of the database to use.")
	flagCreate = flag.Bool("db.create", false,
		"Create database tables on startup")
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
	if *flagCreate {
		rtx.Must(createSchema(db), "Cannot create database schema")
	}

	// Initialize the handler.
	uploadHandler := internal.UploadHandler{
		DBConn: db,
	}

	// Initialize the Echo server.
	e := echo.New()
	e.Validator = &internal.CustomValidator{
		Validator: validator.New(),
	}

	// Endpoints' routing.
	e.POST("/tests", uploadHandler.PostUpload)

	// Start the Echo server.
	e.Logger.Fatal(e.Start(":1323"))
}

// createSchema creates database schema for the Measurement model.
func createSchema(db *pg.DB) error {

	err := db.CreateTable((*model.Measurement)(nil), &orm.CreateTableOptions{})
	if err != nil {
		return err
	}

	return nil
}
