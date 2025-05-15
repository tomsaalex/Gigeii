package model

import (
	"time"

	"github.com/google/uuid"
)

type Availability struct {
	Id               uuid.UUID
	ProductID        uuid.UUID
	StartDate        time.Time
	EndDate          time.Time
	AvailabilityType string
	Days             int32
	Hours            int32
	Price            int32
	MaxParticipants  int32
	Precedance       int32
	CreatedBy        uuid.UUID
}
