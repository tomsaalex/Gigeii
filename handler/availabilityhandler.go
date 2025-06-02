package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"example.com/repository"
	"example.com/service"
	custalerts "example.com/templates/components/alerts"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AvailabilityHandler struct {
	userService           *service.UserService
	availabilityService   *service.AvailabilityService
	availabilityDTOMapper *AvailabilityDTOMapper
}

func NewAvailabilityHandler(
	userService *service.UserService,
	availabilityService *service.AvailabilityService,
) *AvailabilityHandler {
	return &AvailabilityHandler{
		userService:           userService,
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
			custalerts.MakeMultiLineAlertDanger(valErr.ErrorsList).Render(r.Context(), w)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			custalerts.MakeAlertDanger("An unknown error occurred").Render(r.Context(), w)
		}
		return
	}

	precAvailabilityID, err := uuid.Parse(availabilityDTO.PrecedentAvailabilityID)

	if err != nil {
		precAvailabilityID = uuid.Nil
	}

	conflictResolutionMode := availabilityDTO.ConflictResolutionMode

	ownerEmail, emailAttached := r.Context().Value(UserEmailKey).(string)
	if !emailAttached {
		w.WriteHeader(http.StatusUnauthorized)
		custalerts.MakeAlertDanger("You must be authenticated to use this endpoint.").Render(r.Context(), w)
		return
	}

	availabilityOwner, err := h.userService.GetUserByEmail(r.Context(), ownerEmail)
	if !emailAttached {
		w.WriteHeader(http.StatusUnauthorized)
		custalerts.MakeAlertDanger("Cannot add an Availability for the user you are logged in as.").
			Render(r.Context(), w)
		return
	}
	availability.CreatedBy = availabilityOwner.Id

	_, conflictingAvailabilities, err := h.availabilityService.AddAvailability(
		r.Context(),
		*availability,
		precAvailabilityID,
		conflictResolutionMode,
	)

	if err != nil {
		var re *repository.RepositoryError
		var uce *service.UnhandledConflictError
		var enf *repository.EntityNotFoundError

		if errors.As(err, &re) {
			w.WriteHeader(http.StatusInternalServerError)
			custalerts.MakeAlertDanger("Server error. Failed to add Availability.").Render(r.Context(), w)
			return
		}
		if errors.As(err, &uce) {
			w.WriteHeader(http.StatusConflict)
			conflictsDisplay := []string{"The Availability conflicts with following Availabilities:"}
			for _, conflictAvailability := range conflictingAvailabilities {
				conflictsDisplay = append(conflictsDisplay, conflictAvailability.ID.String())
			}
			custalerts.MakeMultiLineAlertDanger(conflictsDisplay).Render(r.Context(), w)
			return
		}
		if errors.As(err, &enf) {
			w.WriteHeader(http.StatusNotFound)
			custalerts.MakeAlertDanger("The Availability offered for conflict resolution wasn't found").
				Render(r.Context(), w)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
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

	precAvailabilityID, err := uuid.Parse(availabilityDTO.PrecedentAvailabilityID)

	if err != nil {
		precAvailabilityID = uuid.Nil
	}

	conflictResolutionMode := availabilityDTO.ConflictResolutionMode

	ownerEmail, emailAttached := r.Context().Value(UserEmailKey).(string)
	if !emailAttached {
		w.WriteHeader(http.StatusUnauthorized)
		custalerts.MakeAlertDanger("You must be authenticated to use this endpoint.").Render(r.Context(), w)
		return
	}

	availabilityOwner, err := h.userService.GetUserByEmail(r.Context(), ownerEmail)
	if !emailAttached {
		w.WriteHeader(http.StatusUnauthorized)
		custalerts.MakeAlertDanger("Cannot add an Availability for the user you are logged in as.").
			Render(r.Context(), w)
		return
	}
	availability.CreatedBy = availabilityOwner.Id

	_, conflictingAvailabilities, err := h.availabilityService.UpdateAvailability(
		r.Context(),
		*availability,
		precAvailabilityID,
		conflictResolutionMode,
	)

	if err != nil {
		var re *repository.RepositoryError
		var uce *service.UnhandledConflictError
		var enf *repository.EntityNotFoundError

		if errors.As(err, &re) {
			w.WriteHeader(http.StatusInternalServerError)
			// A common cause for this is if you're trying to update an availability that doesn't exist. Maybe that should get its own error.
			custalerts.MakeAlertDanger("Server error. Failed to update Availability.").
				Render(r.Context(), w)
			return
		}
		if errors.As(err, &uce) {
			w.WriteHeader(http.StatusConflict)
			conflictsDisplay := []string{"The Availability conflicts with following Availabilities:"}
			for _, conflictAvailability := range conflictingAvailabilities {
				conflictsDisplay = append(conflictsDisplay, conflictAvailability.ID.String())
			}
			custalerts.MakeMultiLineAlertDanger(conflictsDisplay).Render(r.Context(), w)
			return
		}
		if errors.As(err, &enf) {
			w.WriteHeader(http.StatusNotFound)
			custalerts.MakeAlertDanger("The Availability offered for conflict resolution wasn't found").
				Render(r.Context(), w)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		custalerts.MakeAlertDanger(err.Error()).Render(r.Context(), w)
		return
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
