package handlers

import (
	"net/http"
	"encoding/json"
)

func RespondWithError(w http.ResponseWriter, statusCode int, errMsg string  ) {
	RespondWithJson(w, statusCode, map[string]string{"error": errMsg})
}

func RespondWithJson(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}
