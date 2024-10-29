package handlers

import (
	"fmt"
	"net/http"
)

func RetrieveDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	response := fmt.Sprintf("db_id: %s", r.PathValue("id"))

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(response))
}
