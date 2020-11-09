package shopping

import (
	"net/http"

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
