package handler

import (
	"math/rand"
	"net/http"

	"example.com/templates/base"
	"example.com/templates/components/editavailcomp"
	"example.com/templates/pages"
	"example.com/templates/pages/addavailpage"
	"github.com/go-chi/chi/v5"
)

type AvailabilityHandler struct {
}

func NewAvailabilityHandler() *AvailabilityHandler {
	return &AvailabilityHandler{}
}

func (h *AvailabilityHandler) Routes(r chi.Router) {
	r.Get("/add-availability", h.addAvailabilityPage)
	r.Get("/test", h.testModal)
	r.Get("/modal", h.serveModal)
	r.Post("/availabilities", h.mockAdd)
	r.Put("/availabilities", h.mockAdd)

	r.Get("/times", h.addTime)
}

func (h *AvailabilityHandler) addAvailabilityPage(w http.ResponseWriter, r *http.Request) {
	base.PageSkeleton(addavailpage.MakeNewAvailabilityAddPage()).Render(r.Context(), w)
}

func (h *AvailabilityHandler) testModal(w http.ResponseWriter, r *http.Request) {
	base.PageSkeleton(pages.Test()).Render(r.Context(), w)
}

func (h *AvailabilityHandler) serveModal(w http.ResponseWriter, r *http.Request) {
	editavailcomp.MakeNewAvailabilityEditComponent(true).Render(r.Context(), w)
}

func (h *AvailabilityHandler) mockAdd(w http.ResponseWriter, r *http.Request) {
	if rand.Intn(2) == 1 {
		w.Header().Set("HX-Response-Status", "success")
		w.WriteHeader(http.StatusOK)
	} else {
		w.Header().Set("HX-Response-Status", "error")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class=\"alert alert-danger\">Generic error.</div>"))
	}
}

func (h *AvailabilityHandler) addTime(w http.ResponseWriter, r *http.Request) {
	time := r.URL.Query().Get("time-temp")
	if time == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	editavailcomp.TimeComponent(time).Render(r.Context(), w)
}
