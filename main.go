package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/larsmoa/renderdb/geometry"
	"github.com/larsmoa/renderdb/routes"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type applicationArgs struct {
	serverAddress      string
	dbConnectionString string
	dbDriver           string
}

type application struct {
	args   applicationArgs
	db     *sqlx.DB
	repo   geometry.Repository
	router *mux.Router
}

func (a *application) parseArguments() error {
	flag.StringVar(&a.args.serverAddress, "serverAddress", ":8080", "")
	flag.StringVar(&a.args.dbDriver, "driver", "", "")
	flag.StringVar(&a.args.dbConnectionString, "datasource", "", "")
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

	_, err := a.db.Exec(`
            CREATE TABLE IF NOT EXISTS geometry_objects(
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                bounds_x_min REAL NOT NULL,
                bounds_y_min REAL NOT NULL,
                bounds_z_min REAL NOT NULL,
                bounds_x_max REAL NOT NULL,
                bounds_y_max REAL NOT NULL,
                bounds_z_max REAL NOT NULL,
                geometry_data BLOB NOT NULL,
                metadata STRING NOT NULL)`)
	return err
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
	fmt.Printf("Serving at %s...", a.args.serverAddress)
	err = http.ListenAndServe(a.args.serverAddress, nil)
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
