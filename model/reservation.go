package model

import (
	"time"

	"github.com/google/uuid"
)




type Reservation struct {
    ID                        uuid.UUID
    AvailabilityID            uuid.UUID
    ReservationReference      string
    ExternalReservationReference string
    ResellerID                uuid.UUID
    DateTime                  time.Time
    Quantity                  int32
    Status                    string    //CONFIRMED, CANCELLED
    CreatedAt                 time.Time
    UpdatedAt                 time.Time
}
