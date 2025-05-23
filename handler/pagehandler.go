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
