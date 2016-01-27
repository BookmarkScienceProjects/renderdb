package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/larsmoa/renderdb/geometry"
	"github.com/larsmoa/renderdb/routes"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type applicationArgs struct {
	dbConnectionString string
}

type application struct {
	args   applicationArgs
	db     *sqlx.DB
	repo   geometry.Repository
	router *mux.Router
}

func (a *application) parseArguments() error {
	flag.StringVar(&a.args.dbConnectionString, "database", "", "-d")
	flag.Parse()
	return nil
}

func (a *application) initializeDatabase() error {
	if a.args.dbConnectionString != "" {
		db, err := sqlx.Open("postgres", a.args.dbConnectionString)
		a.db = db
		return err
	} else {
		// In-memory test database
		db, err := sqlx.Open("sqlite3", ":memory:")

		db.MustExec(`
            CREATE TABLE geometry_objects(
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                geometry_text STRING NOT NULL,
                metadata STRING NOT NULL
            )`)
		a.db = db
		return err
	}
}

func (a *application) initializeRepository() error {
	var err error
	a.repo, err = geometry.NewRepository(a.db)
	return err
}

func (a *application) initializeRoutes() error {
	a.router = mux.NewRouter()

	geomController := new(routes.GeometryController)
	geomController.Init(a.repo, a.router)

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

	err = a.initializeRepository()
	if err != nil {
		return 3, err
	}

	err = a.initializeRoutes()
	if err != nil {
		return 4, err
	}

	http.Handle("/", a.router)
	fmt.Printf("Serving...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return 5, err
	}

	return 0, nil
}

func main() {
	app := application{}
	app.run()
}
