package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/dto"
	"example.com/service"
	"github.com/go-chi/chi/v5"
)

//request structure
type ReserveRequest struct {
	Data struct {
		DateTime                  string `json:"dateTime"`
		Quantity                  int    `json:"quantity"`
		ExternalReservationReference string `json:"externalReservationReference"`
	} `json:"data"`
}

type CancelRequest struct {
	Data struct {
		ReservationReference         string `json:"reservationReference"`
		ExternalReservationReference string `json:"externalReservationReference"`
	} `json:"data"`
}

type ViewResponse struct {
	Data struct {
		ReservationReference         string `json:"reservationReference"`
		ExternalReservationReference string `json:"externalReservationReference"`
		DateTime                     string `json:"dateTime"`
		Quantity                     int32  `json:"quantity"`
		Status                       string `json:"status"`
	} `json:"data"`
}

type ReserveResponse struct {
	Data struct {
		ReservationReference string `json:"reservationReference"`
	} `json:"data"`
}

type ErrorResponse struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}


type ReservationHandler struct {
	ReservationService service.ReservationService
	AvailabilityService service.AvailabilityService
}


func NewReservationHandler(reservationService service.ReservationService, availabilityService service.AvailabilityService) *ReservationHandler {
	return &ReservationHandler{
		ReservationService: reservationService,
		AvailabilityService: availabilityService,
	}
}





func (h *ReservationHandler) Reserve(w http.ResponseWriter, r *http.Request) {
	var req ReserveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "VALIDATION_FAILURE", "Invalid JSON")
		return
	}
	//fmt.Println("Reserve request:", req)



	reseller, ok := r.Context().Value(ResellerContextKey).(*dto.ResellerResponse)
	if !ok {
		h.writeError(w, "AUTHORIZATION_FAILURE", "Unauthorized")
		return
	}

	dateTime, err := time.Parse(time.RFC3339, req.Data.DateTime)
	if err != nil {
		h.writeError(w, "VALIDATION_FAILURE", "Invalid dateTime format")
		return
	}

	id, err := h.AvailabilityService.GetAvailabilityIDForReservation(r.Context(), dateTime)
	if err != nil {
		h.writeError(w, "VALIDATION_FAILURE", "Invalid availability")
		fmt.Println(err)
		return
	}
		//fmt.Println("Availability ID for reservation:", id)
	reservation := dto.ReservationDTO{
		ResellerID:                   reseller.Id,
		AvailabilityID:               id.String(),
		DateTime:                     dateTime,
		Quantity:                     int32(req.Data.Quantity),
		ExternalReservationReference: req.Data.ExternalReservationReference,
	}

	result, err := h.ReservationService.ReserveOrUpdate(r.Context(), reservation)
	if err != nil {
		if err.Error() == "NO_AVAILABILITY" {
			h.writeError(w, "NO_AVAILABILITY", err.Error())
			return
		}
		h.writeError(w, "INTERNAL_SYSTEM_FAILURE", err.Error())
		return
	}

	resp := ReserveResponse{}
	resp.Data.ReservationReference = result.ReservationReference
	h.writeJSON(w, resp)
}



func (h *ReservationHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	var req CancelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "VALIDATION_FAILURE", "Invalid JSON")
		return
	}

	_, ok := r.Context().Value(ResellerContextKey).(*dto.ResellerResponse)
	if !ok {
		h.writeError(w, "AUTHORIZATION_FAILURE", "Unauthorized")
		return
	}

	err := h.ReservationService.CancelReservation(r.Context(), req.Data.ReservationReference, req.Data.ExternalReservationReference)
	if err != nil {
		switch err.Error() {
		case "INVALID_RESERVATION":
			h.writeError(w, "INVALID_RESERVATION", err.Error())
		case "RESERVATION_ALREADY_CANCELED":
			h.writeError(w, "RESERVATION_ALREADY_CANCELED", err.Error())
		default:
			h.writeError(w, "INTERNAL_SYSTEM_FAILURE", err.Error())
		}
		return
	}

	// Success 200 OK cu gol data:
	h.writeJSON(w, map[string]interface{}{"data": map[string]interface{}{}})
}





func (h *ReservationHandler) ViewReservation(w http.ResponseWriter, r *http.Request) {
	reservationReference := chi.URLParam(r, "reservationReference")

	_, ok := r.Context().Value(ResellerContextKey).(*dto.ResellerResponse)
	if !ok {
		h.writeError(w, "AUTHORIZATION_FAILURE", "Unauthorized")
		return
	}

	result, err := h.ReservationService.GetByReservationReference(r.Context(), reservationReference)
	if err != nil {
		h.writeError(w, "INVALID_RESERVATION", err.Error())
		return
	}

	resp := ViewResponse{}
	resp.Data.ReservationReference = result.ReservationReference
	resp.Data.ExternalReservationReference = result.ExternalReservationReference
	resp.Data.DateTime = result.DateTime.Format(time.RFC3339)
	resp.Data.Quantity = result.Quantity
	resp.Data.Status = result.Status

	h.writeJSON(w, resp)
}





func (h *ReservationHandler) writeError(w http.ResponseWriter, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ErrorResponse{
		ErrorCode:    code,
		ErrorMessage: message,
	})
}

func (h *ReservationHandler) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
