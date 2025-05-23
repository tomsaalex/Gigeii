package handler

import (
	"context"
	"net/http"

	"example.com/service"
)

type contextKey string

const UserEmailKey contextKey = "userEmail"

func JWTContextMiddleware(jwtUtil *service.JwtUtil) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("authCookie")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			email, err := jwtUtil.ParseJWT(cookie.Value)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), UserEmailKey, email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAuth(jwtUtil *service.JwtUtil) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("authCookie")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			_, err = jwtUtil.ParseJWT(cookie.Value)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RedirectIfAuthenticated(jwtUtil *service.JwtUtil) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("authCookie")
			if err == nil {
				if _, err := jwtUtil.ParseJWT(cookie.Value); err == nil {
					http.Redirect(w, r, "/", http.StatusFound)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
