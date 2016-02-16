package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/http2" // FIXME 20160214: Remove when Go 1.6 is released

	"github.com/larsmoa/renderdb/repository"
	"github.com/larsmoa/renderdb/repository/sql"
	"github.com/larsmoa/renderdb/routes"

	"github.com/go-martini/martini"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type applicationArgs struct {
	serverAddress      string
	dbConnectionString string
	dbDriver           string
	useHTTP2           bool
	tlsCertFile        string
	tlsKeyFile         string
}

type application struct {
	args   applicationArgs
	db     *sqlx.DB
	repo   repository.Repository
	router *martini.ClassicMartini
}

func (a *application) parseArguments() error {
	flag.StringVar(&a.args.serverAddress, "serverAddress", ":8080", "")
	flag.StringVar(&a.args.dbDriver, "driver", "",
		"Example: 'sqlite3'")
	flag.StringVar(&a.args.dbConnectionString, "datasource", "",
		"Example: 'file:test.db'")
	flag.StringVar(&a.args.tlsCertFile, "cert", "",
		"TLS certificate to use to secure the HTTP link.")
	flag.StringVar(&a.args.tlsKeyFile, "key", "",
		"TLS private key to use to secure the HTTP link.")
	flag.BoolVar(&a.args.useHTTP2, "http2", false,
		"Enable HTTP2 support. Requires TLS certification and private key.")
	flag.Parse()
	return nil
}

func (a *application) initializeDatabase() error {
	if a.args.dbConnectionString != "" {
		db, err := sqlx.Open(a.args.dbDriver, a.args.dbConnectionString)
		if err != nil {
			return err
		}
		a.db = db
	} else {
		// In-memory test database
		db, err := sqlx.Open("sqlite3", ":memory:")
		if err != nil {
			return err
		}
		a.db = db
	}
	return sql.Initialize(a.db)
}

func (a *application) initializeRepository() error {
	var err error
	a.repo, err = repository.NewRepository(a.db)
	return err
}

func (a *application) initializeRoutes() error {
	a.router = martini.Classic()

	staticController := new(routes.StaticController)
	staticController.Init(a.router)

	geomController := new(routes.GeometryController)
	geomController.Init(a.repo, a.router)

	return nil
}

func (a *application) initializeServer() error {
	srv := &http.Server{
		Addr:    a.args.serverAddress,
		Handler: a.router,
	}

	// Use HTTP 2?
	protocolVersion := "1.1"
	if a.args.useHTTP2 {
		protocolVersion = "2"
		http2.ConfigureServer(srv,
			&http2.Server{
				MaxHandlers:          10,
				MaxConcurrentStreams: 50,
			})
	}

	// TLS certificate/key
	if a.args.tlsCertFile != "" && a.args.tlsKeyFile != "" {
		fmt.Printf("Serving at %s using HTTPS/%s...", a.args.serverAddress, protocolVersion)
		return srv.ListenAndServeTLS(a.args.tlsCertFile, a.args.tlsKeyFile)
	} else if a.args.tlsCertFile != "" || a.args.tlsKeyFile != "" {
		return errors.New("Must provide both TLS certificate and private key.")
	} else if a.args.useHTTP2 {
		return errors.New("Must provide TLS certificate and private key when using HTTP/2.")
	}
	fmt.Printf("Serving at %s using HTTP/%s...", a.args.serverAddress, protocolVersion)
	return srv.ListenAndServe()
}

func (a *application) run() (int, error) {
	err := a.parseArguments()
	if err != nil {
		return 1, err
	}

	err = a.initializeDatabase()
	if err != nil {
		return 2, err
	}

	err = a.initializeRepository()
	if err != nil {
		return 3, err
	}

	err = a.initializeRoutes()
	if err != nil {
		return 4, err
	}

	err = a.initializeServer()
	if err != nil {
		return 5, err
	}

	return 0, nil
}

func main() {
	app := application{}
	code, err := app.run()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	os.Exit(code)
}
