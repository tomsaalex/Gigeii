package handler

import (
	"example.com/service"
	"github.com/go-chi/chi/v5"
)

type RouteDependencies struct {
	UserHandler *UserHandler
	PageHandler *PageHandler
	JwtHelper   *service.JwtUtil
}

func SetupRoutes(dep RouteDependencies) *chi.Mux {
	r := chi.NewRouter()
	r.Use(JWTContextMiddleware(dep.JwtHelper))

	// Guest-only (login/register)
	r.Group(func(r chi.Router) {
		r.With(RedirectIfAuthenticated(dep.JwtHelper)).Get("/login", dep.UserHandler.loginPage)
		r.With(RedirectIfAuthenticated(dep.JwtHelper)).Get("/register", dep.UserHandler.registerPage)
		r.Post("/auth/login", dep.UserHandler.loginUser)
		r.Post("/auth/register", dep.UserHandler.registerUser)
		r.With(RequireAuth(dep.JwtHelper)).Get("/logout", dep.UserHandler.Logout)
		r.With(RequireAuth(dep.JwtHelper)).Post("/logout", dep.UserHandler.Logout)
	})

	// Protected home page
	r.Group(func(r chi.Router) {
		r.With(RequireAuth(dep.JwtHelper)).Get("/", dep.PageHandler.homePage)
		r.With(RequireAuth(dep.JwtHelper)).Get("/calendar", dep.PageHandler.fullCalendarPage)

	})

	return r
}
