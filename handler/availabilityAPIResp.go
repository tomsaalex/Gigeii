package handler

type AvailabilityAPIResponse struct {
	Data struct {
		Availabilities []AvailabilityItem `json:"availabilities"`
	} `json:"data"`
}

type AvailabilityItem struct {
	DateTime  string `json:"dateTime"`
	Vacancies int32  `json:"vacancies"`
	Price     int32  `json:"price"`
}
