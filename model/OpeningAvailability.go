package model

import "time"

type OpeningAvailability struct {
	DateTime  time.Time
	Vacancies int32
	Price     int32
	Precedance int32
}
