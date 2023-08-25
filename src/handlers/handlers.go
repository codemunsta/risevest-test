package handlers

import (
	"encoding/json"
	"net/http"
)

func NewApi(writer http.ResponseWriter, request *http.Request) {
	response := map[string]interface{}{
		"User":    request.RemoteAddr,
		"Message": "Pong",
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}
