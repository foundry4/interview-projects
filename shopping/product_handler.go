package shopping

import (
	"net/http"
	"encoding/json"

	"github.com/gorilla/mux"
)

// ProductHandler handles all REST calls related to inventory management
type ProductHandler struct {
	Product *ProductService
}

// NewProductHandler creates a new ProductHandler
// needed for catering to all the REST calls towrads the inventory
func NewProductHandler(ps *ProductService) *ProductHandler {
	return &ProductHandler{Product: ps}
}

// Create handler retrives all the Products
func (ph *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {

	p := new(Product)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondWithError(w, http.StatusBadRequest, err)
	}
	defer r.Body.Close()

	p, err := ph.Product.Save(p)
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

// Delete handler deletes the Product based on the id passed
func (ph *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	err := ph.Product.Delete(id)
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

// QueryByID handler retrives the Product based on the id passed
func (ph *ProductHandler) QueryByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	p, err := ph.Product.Fetch(id)
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
func (ph *ProductHandler) Query(w http.ResponseWriter, r *http.Request) {

	p, err := ph.Product.FetchAll()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err)
	}

	Respond(w, http.StatusOK, p)
}
