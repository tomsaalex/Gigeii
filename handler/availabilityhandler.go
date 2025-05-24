package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"example.com/service"
	custalerts "example.com/templates/components/alerts"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AvailabilityHandler struct {
	availabilityService   *service.AvailabilityService
	availabilityDTOMapper *AvailabilityDTOMapper
}

func NewAvailabilityHandler(availabilityService *service.AvailabilityService) *AvailabilityHandler {
	return &AvailabilityHandler{
		availabilityService:   availabilityService,
		availabilityDTOMapper: &AvailabilityDTOMapper{},
	}
}

func (h *AvailabilityHandler) addAvailability(w http.ResponseWriter, r *http.Request) {
	var availabilityDTO AvailabilityDTO
	err := json.NewDecoder(r.Body).Decode(&availabilityDTO)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		custalerts.MakeAlertDanger("Couldn't process Availability data").Render(r.Context(), w)
		return
	}

	availability, err := h.availabilityDTOMapper.AvailabilityDTOToAvailability(availabilityDTO)

	if err != nil {
		var valErr *service.ValidationError
		if errors.As(err, &valErr) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println(valErr)
			custalerts.MakeMultiLineAlertDanger(valErr.ErrorsList).Render(r.Context(), w)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			custalerts.MakeAlertDanger("An unknown error occurred").Render(r.Context(), w)
		}
		return
	}

	_, err = h.availabilityService.AddAvailability(r.Context(), *availability)

	if err != nil {
		custalerts.MakeAlertDanger(err.Error()).Render(r.Context(), w)
		return
	}

	// TODO: Not sure what to return here, really. I guess a redirect to something else.
	w.Header().Set("HX-Redirect", "/")
}

func (h *AvailabilityHandler) updateAvailability(w http.ResponseWriter, r *http.Request) {
	availabilityIDString := chi.URLParam(r, "id")

	_, err := uuid.Parse(availabilityIDString)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		custalerts.MakeAlertDanger("Couldn't process Availability ID").Render(r.Context(), w)
		return
	}

	var availabilityDTO AvailabilityDTO
	err = json.NewDecoder(r.Body).Decode(&availabilityDTO)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		custalerts.MakeAlertDanger("Couldn't process Availability data").Render(r.Context(), w)
		return
	}

	if availabilityIDString != availabilityDTO.AvailabilityID {
		w.WriteHeader(http.StatusBadRequest)
		custalerts.MakeAlertDanger("The ID in the URI doesn't match the ID in the request body").Render(r.Context(), w)
		return
	}

	availability, err := h.availabilityDTOMapper.AvailabilityDTOToAvailability(availabilityDTO)

	if err != nil {
		var valErr *service.ValidationError
		if errors.As(err, &valErr) {
			w.WriteHeader(http.StatusBadRequest)
			custalerts.MakeMultiLineAlertDanger(valErr.ErrorsList).Render(r.Context(), w)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			custalerts.MakeAlertDanger("An unknown error occurred").Render(r.Context(), w)
		}
		return
	}

	_, err = h.availabilityService.UpdateAvailability(r.Context(), *availability)
	if err != nil {
		// TODO: Handle errors here

		fmt.Printf("Rollback transaction: %d\n", availability.Hours[len(availability.Hours)-1].Hour)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *AvailabilityHandler) deleteAvailability(w http.ResponseWriter, r *http.Request) {
	availabilityIDString := chi.URLParam(r, "id")

	availabilityID, err := uuid.Parse(availabilityIDString)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		custalerts.MakeAlertDanger("Couldn't process Availability ID").Render(r.Context(), w)
		return
	}

	_, err = h.availabilityService.DeleteAvailability(r.Context(), availabilityID)

	if err != nil {
		// TODO: Handle errors here
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
