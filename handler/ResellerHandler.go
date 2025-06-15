package handler

import (
	"encoding/json"
	"net/http"

	"example.com/dto"
	"example.com/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ResellerHandler struct {
	service service.ResellerService
	availabilityHandler *AvailabilityHandler

}

func NewResellerHandler(service service.ResellerService, availabilityHandler *AvailabilityHandler) *ResellerHandler {
	return &ResellerHandler{
		service: service,
		availabilityHandler: availabilityHandler,
	}
}


func (h *ResellerHandler) Routes(r chi.Router) {
	// r.Route("/api/resellers", func(r chi.Router) {
		// r.Post("/register", h.Register)
		// r.Post("/login", h.Login)
		// r.Get("/", h.ListAll)
		// r.Get("/{id}", h.GetByID)
		// r.Delete("/{id}", h.Delete)
		r.With(BasicAuth(h.service)).Get("/1/availabilities/", h.availabilityHandler.getAvailabilitiesInRange)


	// })
}

// Register reseller
func (h *ResellerHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterResellerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	reseller, err := h.service.Register(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(reseller)
}

// Login reseller
func (h *ResellerHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	reseller, err := h.service.Login(r.Context(), req)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(reseller)
}

// Get all resellers
func (h *ResellerHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	resellers, err := h.service.ListAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resellers)
}

// Get by ID
func (h *ResellerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	reseller, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(reseller)
}

// Delete reseller
func (h *ResellerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		http.Error(w, "Delete failed: "+err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
