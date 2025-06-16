package dto

import "time"

type ReservationDTO struct {
	ReservationReference         string    `json:"reservationReference"`
	ExternalReservationReference string    `json:"externalReservationReference"`
	ResellerID                   string    `json:"resellerId"`
	AvailabilityID               string    `json:"availabilityId"`
	DateTime                     time.Time `json:"dateTime"`
	Quantity                     int32     `json:"quantity"`
	Status                       string    `json:"status"`
}
