package shopping

import (
	"net/http"
)

func (app *App) routeRequests() {

	ps := NewProductService(app.Inventory)
	ph := NewProductHandler(ps)
	// for inventory use
	app.Router.HandleFunc("/products/{id}", ph.QueryByID).Methods(http.MethodGet)
	app.Router.HandleFunc("/products/{id}", ph.Delete).Methods(http.MethodDelete)
	app.Router.HandleFunc("/products", ph.Query).Methods(http.MethodGet)
	app.Router.HandleFunc("/products", ph.Create).Methods(http.MethodPost)


	cart := NewDB(make(map[string]*Product))
	cs := NewProductService(cart)
	ch := NewCartHandler(cs)
	// for basket/cart use
	app.Router.HandleFunc("/cart", ch.Create).Methods(http.MethodPost)
	app.Router.HandleFunc("/cart/{id}", ch.QueryByID).Methods(http.MethodGet)
	app.Router.HandleFunc("/cart", ch.Query).Methods(http.MethodGet)
	app.Router.HandleFunc("/cart/{id}", ch.Query).Methods(http.MethodDelete)
	
	// for checkout
	app.Router.HandleFunc("/checkout/{id}", ch.Checkout).Methods(http.MethodDelete)

}
