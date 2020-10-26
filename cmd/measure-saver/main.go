// measure-saver is a REST API which receives measurements from the
// M-Lab Measure Chrome Extension in JSON format and stores them into a
// PostgreSQL database.
package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/m-lab/go/flagx"
	"github.com/m-lab/go/httpx"
	"github.com/m-lab/go/rtx"
	"github.com/m-lab/measure-saver/internal"
	"github.com/m-lab/measure-saver/internal/measurements"
	"github.com/m-lab/measure-saver/internal/model"
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
	DefaultDBName = "measure-saver"
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
	flagDBURL = flag.String("db.url", "", "Database connection URL")

	flagTLSCert = flag.String("tls.cert", "",
		"Path to the TLS certificate file to use.")
	flagTLSKey = flag.String("tls.key", "",
		"Path to the TLS key file to use.")

	flagKeysFile = flagx.FileBytesArray{}

	pgConnect = func(opt *pg.Options) internal.Database {
		return pg.Connect(opt)
	}

	ctx, cancel = context.WithCancel(context.Background())
)

func init() {
	flag.Var(&flagKeysFile, "keys.file",
		"Text file containing the allowed API keys.")
}

func initEchoServer() *echo.Echo {
	// Initialize the Echo server.
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Validator = &measurements.Validator{
		Validator: validator.New(),
	}

	keysFiles := flagKeysFile.Get()
	if len(keysFiles) != 0 {
		// Read API keys file and set up key auth middleware.
		allowedKeys, err := readKeysFiles(keysFiles...)
		rtx.Must(err, "Cannot read keys file")

		e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
			KeyLookup: "query:key",
			Validator: func(key string, c echo.Context) (bool, error) {
				_, ok := allowedKeys[key]
				return ok, nil
			},
		}))
	} else {
		fmt.Println("Warning: no keys file specified. The endpoints will be " +
			"unprotected.")
	}

	return e
}

func main() {
	flag.Parse()
	rtx.Must(flagx.ArgsFromEnv(flag.CommandLine), "Could not parse env args")

	// Initialize database connection.
	// If a complete URL has been provided, use it instead of the individual
	// user/pass/host flags.
	var options *pg.Options
	var err error
	if *flagDBURL != "" {
		options, err = pg.ParseURL(*flagDBURL)
		rtx.Must(err, "Invalid connection URL provided")
	} else {
		options = &pg.Options{
			Addr:     *flagDBAddr,
			User:     *flagDBUser,
			Password: *flagDBPass,
			Database: *flagDBName,
		}
	}
	db := pgConnect(options)
	defer db.Close()

	// Check connection has been successful.
	// go-pg provides a connection pool, so connections aren't actually made
	// until some query is executed, thus we run a simple SELECT 1 to verify
	// the connection works.
	_, err = db.Exec("SELECT 1")
	rtx.Must(err, "Connection to the database failed")

	// Create schema if needed.
	rtx.Must(createSchema(db), "Cannot create database schema")

	// Initialize the handler.
	measurementsHandler := measurements.Handler{
		DB: db,
	}

	// Initialize Echo server.
	e := initEchoServer()

	// Endpoints' routing.
	e.POST("/v0/measurements", measurementsHandler.Post)

	// Start the Echo server.
	e.Server.Addr = *flagListenAddr
	if *flagTLSCert != "" && *flagTLSKey != "" {
		rtx.Must(httpx.ListenAndServeTLSAsync(e.Server, *flagTLSCert, *flagTLSKey),
			"Cannot initialize TLS server")
	} else {
		rtx.Must(httpx.ListenAndServeAsync(e.Server),
			"Cannot initialize server")
	}

	<-ctx.Done()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
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

// readKeysFiles reads the keys files and returns a map of keys, one per
// line. Lines beginning with '#' are treated as comments and ignored.
func readKeysFiles(files ...[]byte) (map[string]bool, error) {
	keys := map[string]bool{}
	for _, file := range files {
		scanner := bufio.NewScanner(bytes.NewReader(file))
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "#") &&
				strings.TrimSpace(line) != "" {
				keys[line] = true
			}
		}
	}
	return keys, nil
}
