package main

import (
	"database/sql"
	"flag"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type applicationArgs struct {
	dbConnectionString string
}

type application struct {
	args   applicationArgs
	db     *sql.DB
	router *mux.Router
}

func (a *application) parseArguments() error {
	flag.StringVar(&a.args.dbConnectionString, "database", "postgres://username:password@localhost", "-d")
	flag.Parse()

	return nil
}

func (a *application) initializeDatabase() error {
	db, err := sql.Open("postgres", a.args.dbConnectionString)
	a.db = db
	return err
}

func (a *application) initializeRoutes() error {
	a.router = mux.NewRouter()
	return nil
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

	err = a.initializeRoutes()
	if err != nil {
		return 3, err
	}
	return 0, nil
}

func main() {
	app := application{}
	app.run()
}
