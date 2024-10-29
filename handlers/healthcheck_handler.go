package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aelpxy/dbctl/structs"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := structs.Response{Message: "dbctl API is ready to serve you!"}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
