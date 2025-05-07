package handler

import (
	"net/http"

	"example.com/service"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService *service.UserService
	jwtUtil     *service.JwtUtil
}

func NewUserHandler(userService *service.UserService, jwtUtil *service.JwtUtil) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtUtil:     jwtUtil,
	}
}

func (h *UserHandler) Routes(r chi.Router) {

}

func LoginUser(w http.ResponseWriter, r *http.Request) {

}

func RegisterUser(w http.ResponseWriter, r *http.Request) {

}
