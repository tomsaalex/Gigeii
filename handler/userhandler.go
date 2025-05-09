package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"example.com/repository"
	"example.com/service"
	"example.com/templates/base"
	custalerts "example.com/templates/components/alerts"
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
		custalerts.MakeAlertDanger("Invalid username or password. Please try again.").Render(r.Context(), w)
		return
	}

	errorsList := make([]string, 0)
	if userDTO.Email == "" {
		errorsList = append(errorsList, "Email address is missing.")
	}
	if userDTO.Password == "" {
		errorsList = append(errorsList, "Password is missing.")
	}

	if len(errorsList) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		custalerts.MakeMultiLineAlertDanger(errorsList).Render(r.Context(), w)
		return
	}

	loggedUser, err := h.userService.Login(r.Context(), userDTO.Email, userDTO.Password)

	if err != nil {
		var authErr *service.AuthError
		var entityNotFoundErr *repository.EntityNotFoundError

		if errors.As(err, &authErr) {
			w.WriteHeader(http.StatusBadRequest)
			custalerts.MakeAlertDanger("Credentials invalid.").Render(r.Context(), w)
		} else if errors.As(err, &entityNotFoundErr) {
			w.WriteHeader(http.StatusBadRequest)
			custalerts.MakeAlertDanger("No user has that email address.").Render(r.Context(), w)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			custalerts.MakeAlertDanger("Unknown error occurred.").Render(r.Context(), w)
		}
		return
	}

	token, err := h.jwtUtil.GenerateJWT(loggedUser.Email)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		custalerts.MakeAlertDanger("Unknown error occurred.").Render(r.Context(), w)
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
		custalerts.MakeAlertDanger("Invalid username or password. Please try again.").Render(r.Context(), w)
		return
	}

	errorsList := make([]string, 0)
	if userDTO.Email == "" {
		errorsList = append(errorsList, "Email is missing.")
	}
	if userDTO.Username == "" {
		errorsList = append(errorsList, "Username is missing.")
	}
	if userDTO.Password == "" {
		errorsList = append(errorsList, "Password is missing.")
	}
	if userDTO.ConfirmPassword == "" {
		errorsList = append(errorsList, "Confirm Password is missing.")
	} else if userDTO.Password != userDTO.ConfirmPassword {
		errorsList = append(errorsList, "Passwords do not match.")
	}

	if len(errorsList) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		custalerts.MakeMultiLineAlertDanger(errorsList).Render(r.Context(), w)
		return
	}

	user := h.userDTOMapper.RegistrationDTOToUser(userDTO)

	registeredUser, err := h.userService.Register(r.Context(), user, userDTO.Password)
	if err != nil {
		var authErr *service.AuthError
		var duplicateEntityErr *repository.DuplicateEntityError

		if errors.As(err, &authErr) {
			w.WriteHeader(http.StatusInternalServerError)
			custalerts.MakeAlertDanger("Couldn't register user.").Render(r.Context(), w)
		} else if errors.As(err, &duplicateEntityErr) {
			w.WriteHeader(http.StatusBadRequest)
			custalerts.MakeAlertDanger("Email is alreay in use by a different user.").Render(r.Context(), w)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			custalerts.MakeAlertDanger("An unknown error occurred.").Render(r.Context(), w)
		}
		return
	}

	token, err := h.jwtUtil.GenerateJWT(registeredUser.Email)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		custalerts.MakeAlertDanger("Couldn't register user.").Render(r.Context(), w)
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
