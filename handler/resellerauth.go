package handler

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"example.com/dto"
	"example.com/service"
)


const ResellerContextKey contextKey = "reseller"

func BasicAuth(resellerService service.ResellerService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				writeErrorResponse(w, http.StatusUnauthorized, "AUTHORIZATION_FAILURE", "Missing authorization header")
				return
			}

			encoded := strings.TrimPrefix(authHeader, "Basic ")
			decoded, err := base64.StdEncoding.DecodeString(encoded)
			if err != nil {
				writeErrorResponse(w, http.StatusUnauthorized, "AUTHORIZATION_FAILURE", "Malformed authorization header")
				return
			}

			parts := strings.SplitN(string(decoded), ":", 2)
			if len(parts) != 2 {
				writeErrorResponse(w, http.StatusUnauthorized, "AUTHORIZATION_FAILURE", "Invalid authorization format")
				return
			}

			username := parts[0]
			password := parts[1]

			// Refolosim Login-ul tău deja existent
			reseller, err := resellerService.Login(r.Context(), dto.LoginRequest{
				Username: username,
				Password: password,
			})
			if err != nil {
				writeErrorResponse(w, http.StatusUnauthorized, "AUTHORIZATION_FAILURE", "Invalid credentials")
				return
			}

			// Dacă vrei să ai reseller-ul mai târziu în context
			ctx := context.WithValue(r.Context(), ResellerContextKey, reseller)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
