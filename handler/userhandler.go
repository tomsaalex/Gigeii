package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"example.com/repository"
	"example.com/service"
	"example.com/templates/base"
	"example.com/templates/pages"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService   *service.UserService
	jwtUtil       *service.JwtUtil
	userDTOMapper *UserDTOMapper
}

func NewUserHandler(userService *service.UserService, jwtUtil *service.JwtUtil) *UserHandler {
	return &UserHandler{
		userService:   userService,
		jwtUtil:       jwtUtil,
		userDTOMapper: &UserDTOMapper{},
	}
}

func (h *UserHandler) Routes(r chi.Router) {
	r.Post("/auth/register", h.registerUser)
	r.Post("/auth/login", h.loginUser)

	r.Get("/register", h.registerPage)
	r.Get("/login", h.loginPage)
}

func (h *UserHandler) registerPage(w http.ResponseWriter, r *http.Request) {
	base.PageSkeleton(pages.RegisterPage()).Render(r.Context(), w)
}

func (h *UserHandler) loginPage(w http.ResponseWriter, r *http.Request) {
	base.PageSkeleton(pages.LoginPage()).Render(r.Context(), w)
}

func (h *UserHandler) loginUser(w http.ResponseWriter, r *http.Request) {
	var userDTO UserLoginDTO
	err := json.NewDecoder(r.Body).Decode(&userDTO)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class=\"alert alert-danger\">Invalid username or password. Please try again.</div>"))
		return
	}

	errHappened := false
	errorsList := ""
	if userDTO.Email == "" {
		errorsList += "Email address missing\n"
		errHappened = true
	}
	if userDTO.Password == "" {
		errorsList += "Password missing\n"
		errHappened = true
	}

	// TODO: Return here if any error happened
	if errHappened {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class=\"alert alert-danger\">" + errorsList + "</div>"))
		return
	}

	loggedUser, err := h.userService.Login(r.Context(), userDTO.Email, userDTO.Password)

	if err != nil {
		var authErr *service.AuthError
		var entityNotFoundErr *repository.EntityNotFoundError

		if errors.As(err, &authErr) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("<div class=\"alert alert-danger\"> Credentials invalid </div>"))
		} else if errors.As(err, &entityNotFoundErr) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("<div class=\"alert alert-danger\"> No user has that email address </div>"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("<div class=\"alert alert-danger\"> Unknown error occurred </div>"))
		}
		return
	}

	token, err := h.jwtUtil.GenerateJWT(loggedUser.Email)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("<div class=\"alert alert-danger\"> Unknown error occurred </div>"))
		return
	}

	cookie := http.Cookie{
		Name:     "authCookie",
		Value:    token,
		Path:     "/",
		MaxAge:   int(h.jwtUtil.TokenTTL),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	fmt.Println()

	http.SetCookie(w, &cookie)
	w.Header().Set("HX-Redirect", "/")
}

func (h *UserHandler) registerUser(w http.ResponseWriter, r *http.Request) {
	var userDTO UserRegistrationDTO
	err := json.NewDecoder(r.Body).Decode(&userDTO)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class=\"alert alert-danger\"> Invalid username or password. Please try again. </div>"))
		return
	}

	errorsList := ""
	if userDTO.Email == "" {
		errorsList += "Email is missing\n"
	}
	if userDTO.Username == "" {
		errorsList += "Username is missing\n"
	}
	if userDTO.Password == "" {
		errorsList += "Password is missing\n"
	}
	if userDTO.ConfirmPassword == "" {
		errorsList += "Confirm Password is missing\n"
	}
	if userDTO.Password != userDTO.ConfirmPassword {
		errorsList += "Passwords do not match\n"
	}

	if errorsList != "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class=\"alert alert-danger\">" + errorsList + "</div>"))
		return
	}

	user := h.userDTOMapper.RegistrationDTOToUser(userDTO)

	registeredUser, err := h.userService.Register(r.Context(), user, userDTO.Password)
	if err != nil {
		var authErr *service.AuthError
		var duplicateEntityErr *repository.DuplicateEntityError

		if errors.As(err, &authErr) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("<div class=\"alert alert-danger\"> Couldn't register user </div>"))
		} else if errors.As(err, &duplicateEntityErr) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("<div class=\"alert alert-danger\"> Email is alreay in use by a different user </div>"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("<div class=\"alert alert-danger\"> An unknown error occurred </div>"))
		}
		return
	}

	token, err := h.jwtUtil.GenerateJWT(registeredUser.Email)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("<div class=\"alert alert-danger\"> Couldn't register user </div>"))

		return
	}

	cookie := http.Cookie{
		Name:     "authCookie",
		Value:    token,
		Path:     "/",
		MaxAge:   int(h.jwtUtil.TokenTTL),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
	w.Header().Set("HX-Redirect", "/")
}
