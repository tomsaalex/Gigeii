package model

import (
	"time"

	"github.com/google/uuid"
)

type TimeOfDay struct {
	Hour   int32
	Minute int32
}

type Availability struct {
	ID              uuid.UUID
	StartDate       time.Time
	EndDate         time.Time
	Days            int32
	Hours           []TimeOfDay
	Price           int32
	MaxParticipants int32
	Precedance      int32
	CreatedBy       uuid.UUID
	Duration        time.Duration
	Notes           string //added for additional information
}
