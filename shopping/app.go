package shopping

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const dataDumpFile = "shopping/data/dump/prod/products.json"

// App encapsulates two main components
// viz. router and db
type App struct {
	Router    *mux.Router
	Inventory *DB
}

// NewApp craetes a new App instance each time the application is restarted
// Note: the DB instance is epehmeral and dies with the app shutdown
func NewApp(db *DB) *App {
	return &App{
		Inventory: db,
	}
}

// Init is responsible for initialisiing a new App everytime
// this also includes (in case of thsi app), loading/seeding data to the db from a file
func (app *App) Init() {

	pm := make(map[string]*Product)

	inv := NewDB(pm)

	pa := NewApp(inv)

	// fill the db from data dump
	db, err := Load(dataDumpFile, pa.Inventory)
	if err != nil {
		log.Fatalf("error loading data to db: %v", err)
	}

	app.Inventory = db
	app.Router = mux.NewRouter()

	app.routeRequests()

}

// Run the server to listen and serve on the addr
func (app *App) Run(addr string) {
	fmt.Printf("App started - listening to port %s", addr)
	log.Fatal(http.ListenAndServe(addr, app.Router))
}
