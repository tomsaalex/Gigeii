package handler

type AvailabilityDTO struct {
	AvailabilityID          string   `json:"availability_id,omitempty"`
	StartDate               string   `json:"start_date"`
	EndDate                 string   `json:"end_date"`
	Days                    []int32  `json:"days"`
	Hours                   []string `json:"hours"`
	Price                   string   `json:"price"`
	MaxParticipants         int32    `json:"max_participants"`
	PrecedentAvailabilityID string   `json:"prec_availability_id,omitempty"`
	ConflictResolutionMode  bool     `json:"resolve_conflict"`
	Duration               int32    `json:"duration"`
    Precedance             int32    `json:"precedance"`
	Notes                  string   `json:"notes, omitempty"` 
}
