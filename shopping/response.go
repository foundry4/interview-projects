package shopping

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Respond transforms a model type to JSON and sends it back to the client
func Respond(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

// RespondWithError sends a response with error back to the client.
func RespondWithError(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		Respond(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	Respond(w, http.StatusBadRequest, nil)
}
