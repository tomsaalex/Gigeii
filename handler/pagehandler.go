package handler

import (
    "context"
    "net/http"

    "example.com/templates/base"
    "example.com/templates/pages"
)

type PageHandler struct{}

func NewPageHandler() *PageHandler {
	return &PageHandler{}
}

func (h *PageHandler) homePage(w http.ResponseWriter, r *http.Request) {
	if err := base.PageSkeleton(pages.HomePage()).Render(context.Background(), w); err != nil {
		http.Error(w, "Failed to render homepage", http.StatusInternalServerError)
	}
}

func (h *PageHandler) fullCalendarPage(w http.ResponseWriter, r *http.Request) {
    base.PageSkeleton(pages.FullCalendarPage()).Render(context.Background(), w)
}

/*
func (h *AvailabilityHandler) getAvailabilities(w http.ResponseWriter, r *http.Request) {
    availabilities, err := h.availabilityService.GetAllAvailabilities(r.Context())
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        custalerts.MakeAlertDanger("Failed to fetch availabilities.").Render(r.Context(), w)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(availabilities)
}

*/
