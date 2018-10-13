package handlers

import (
	"net/http"
	"encoding/json"
)

func RespondWithError(writer http.ResponseWriter, code int, message string  ) {
	RespondWithJson(writer, code, map[string]string{"error": message})
}

func RespondWithJson(writer http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(response)
}
