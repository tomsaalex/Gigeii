package handler

import (
	"encoding/json"
	"net/http"
)

type ErrorResponseAvailability struct {
	ErrorCode    string `json:"errorCode"`    // AUTHORIZATION_FAILURE, VALIDATION_FAILURE, INTERNAL_SYSTEM_FAILURE
	ErrorMessage string `json:"errorMessage"` // mesajul concret
}
func writeErrorResponse(w http.ResponseWriter, code int, errorCode string, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := ErrorResponseAvailability{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	}

	json.NewEncoder(w).Encode(resp)
}