package shopping

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// CartHandler handles all REST calls related to basket/cart management
type CartHandler struct {
	Cart *ProductService
}

// NewCartHandler creates a new ProductHandler
// needed for catering to all the REST calls towards the cart
func NewCartHandler(cs *ProductService) *CartHandler {
	return &CartHandler{Cart: cs}

}

// Create handler retrives all the Products
func (ch *CartHandler) Create(w http.ResponseWriter, r *http.Request) {

	p := new(Product)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondWithError(w, http.StatusBadRequest, err)
	}
	defer r.Body.Close()

	p, err := ch.Cart.Save(p)
	if err != nil {
		switch err {
		case ErrorCorruptDb:
			RespondWithError(w, http.StatusInternalServerError, err)
		default:
			RespondWithError(w, http.StatusBadRequest, err)
		}
	}

	Respond(w, http.StatusCreated, p)
}

// QueryByID handler retrives the Product based on the id passed
func (ch *CartHandler) QueryByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	p, err := ch.Cart.Fetch(id)
	if err != nil {
		switch err {
		case ErrorCorruptDb:
			RespondWithError(w, http.StatusInternalServerError, err)
		default:
			RespondWithError(w, http.StatusNotFound, err)
		}
	}

	Respond(w, http.StatusOK, p)
}

// Query handler retrives all the Products
func (ch *CartHandler) Query(w http.ResponseWriter, r *http.Request) {

	p, err := ch.Cart.FetchAll()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err)
	}

	Respond(w, http.StatusOK, p)
}

// Delete handler deletes the Product based on the id passed
func (ch *CartHandler) Delete(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	err := ch.Cart.Delete(id)
	if err != nil {
		switch err {
		case ErrorCorruptDb:
			RespondWithError(w, http.StatusInternalServerError, err)
		default:
			RespondWithError(w, http.StatusNotFound, err)
		}
	}

	Respond(w, http.StatusNoContent, nil)
}

// Checkout hadlers enables to checkout items from the cart
func (ch *CartHandler) Checkout(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	err := ch.Cart.Delete(id)
	if err != nil {
		switch err {
		case ErrorCorruptDb:
			RespondWithError(w, http.StatusInternalServerError, err)
		default:
			RespondWithError(w, http.StatusNotFound, err)
		}
	}

	msg := fmt.Sprintf("Item with id %s has been successfully checked out. Please check your email for confirmation.", id)
	Respond(w, http.StatusOK, msg)
}
