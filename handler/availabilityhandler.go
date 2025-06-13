package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

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

func (h *AvailabilityHandler) getAllAvailabilities(w http.ResponseWriter, r *http.Request) {
	availabilities, err := h.availabilityService.GetAllAvailabilities(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		custalerts.MakeAlertDanger("Failed to fetch availabilities.").Render(r.Context(), w)
		return
	}

	response := make([]AvailabilityDTO, 0, len(availabilities))
	for _, avail := range availabilities {
		days := make([]int32, 0)
		for i := 0; i < 7; i++ {
			if avail.Days&(1<<i) != 0 {
				days = append(days, int32(6-i)) // Convert bitmask to day indices (0=Sun, 6=Sat)

			}
		}
		hours := make([]string, 0, len(avail.Hours))
		for _, h := range avail.Hours {
			hours = append(hours, fmt.Sprintf("%02d:%02d", h.Hour, h.Minute))
		}
		response = append(response, AvailabilityDTO{
			AvailabilityID:          avail.ID.String(),
			StartDate:               avail.StartDate.Format("2006-01-02"),
			EndDate:                 avail.EndDate.Format("2006-01-02"),
			Days:                    days,
			Hours:                   hours,
			Price:                   fmt.Sprintf("%d.%02d", avail.Price/100, avail.Price%100),
			MaxParticipants:         avail.MaxParticipants,
			Notes:                   avail.Notes, // Include if notes are added to AvailabilityDTO
			ConflictResolutionMode:  false,       // Default value
			Duration: 			  int32(avail.Duration.Hours()),
			Precedance: avail.Precedance,
			PrecedentAvailabilityID: "",          // Default empty
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *AvailabilityHandler) getAvailabilityByID(w http.ResponseWriter, r *http.Request) {
	availabilityIDString := chi.URLParam(r, "id")
	availabilityID, err := uuid.Parse(availabilityIDString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		custalerts.MakeAlertDanger("Couldn't process Availability ID").Render(r.Context(), w)
		return
	}

	avail, err := h.availabilityService.GetAvailabilityByID(r.Context(), availabilityID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		custalerts.MakeAlertDanger("Availability not found").Render(r.Context(), w)
		return
	}

	// Map model to DTO
	dto := h.availabilityDTOMapper.AvailabilityToDTO(avail)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto)
}
func (h *AvailabilityHandler) getAvailabilitiesInRange(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("fromDateTime")
	toStr := r.URL.Query().Get("toDateTime")

	if fromStr == "" || toStr == "" {
		writeErrorResponse(w, http.StatusBadRequest, "VALIDATION_FAILURE", "fromDateTime and toDateTime query parameters are required")
		return
	}

	// Validăm că sunt în format 
	if _, err := time.Parse(time.RFC3339, fromStr); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "VALIDATION_FAILURE", "Invalid fromDateTime format: "+err.Error())
		return
	}

	if _, err := time.Parse(time.RFC3339, toStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeErrorResponse(w, http.StatusBadRequest, "VALIDATION_FAILURE", "Invalid toDateTime format: "+err.Error())
		return
	}

	// Doar trimitem stringurile mai departe
	availabilities, err := h.availabilityService.GetAvailabilitiesInRange(
		r.Context(),
		fromStr,
		toStr,
	)
	if err != nil {
		//fmt.Println("Error fetching availabilities:", err)
		writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_SYSTEM_FAILURE", "Internal error while fetching availabilities")
		return
	}

	items := make([]AvailabilityItem, 0, len(availabilities))
	for _, a := range availabilities {
		items = append(items, AvailabilityItem{
			DateTime:  a.DateTime.Format(time.RFC3339),
			Vacancies: a.Vacancies,
			Price:     a.Price,
		})
	}

	response := AvailabilityAPIResponse{}
	response.Data.Availabilities = items

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
