package handler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"example.com/model"
	"example.com/service"
	"github.com/google/uuid"
)

type AvailabilityDTOMapper struct {
}

func DaysSliceToBitmap(days []int32) (int32, error) {
	daysBitmap := int32(0)
	minDayValue := int32(0)
	maxDayValue := int32(6)
	lastDay := int32(6)

	for _, day := range days {
		if day < minDayValue || day > maxDayValue {
			return 0, fmt.Errorf("The provided days are invalid")
		}

		daysBitmap |= 1 << (int(lastDay - day))
	}

	return daysBitmap, nil
}

func HoursSliceToTimeOfDaySlice(stringHours []string) ([]model.TimeOfDay, error) {
	hours := make([]model.TimeOfDay, 0)

	for _, t := range stringHours {
		hour, err := strconv.Atoi(t[:2])
		if err != nil || hour < 0 || hour > 23 {
			return nil, fmt.Errorf("Some of the hours provided are not valid")
		}

		minute, err := strconv.Atoi(t[3:5])
		if err != nil || minute < 0 || minute > 59 {
			return nil, fmt.Errorf("Some of the hours provided are not valid")
		}
		hours = append(hours, model.TimeOfDay{Hour: int32(hour), Minute: int32(minute)})
	}

	return hours, nil
}

func StringPriceToInt(strPrice string) (int, error) {
	numSubstrings := 2
	parts := strings.SplitN(strPrice, ".", numSubstrings)

	// Handle the whole number part
	whole := parts[0]
	if whole == "" {
		return 0, fmt.Errorf("Price is invalid")
	}

	// Handle the fractional (decimal) part
	fraction := "00"
	if len(parts) == numSubstrings {
		raw := parts[1] + "00" // pad in case it's too short
		fraction = raw[:2]     // truncate to exactly 2 digits
	}

	// Combine into a string like "1234"
	combined := whole + fraction

	// Convert to int
	return strconv.Atoi(combined)
}

func (m *AvailabilityDTOMapper) AvailabilityDTOToAvailability(
	availabilityDTO AvailabilityDTO,
) (*model.Availability, error) {
	errorsList := make([]string, 0)
	if availabilityDTO.StartDate == "" {
		errorsList = append(errorsList, "You need to select a start date.")
	}
	if availabilityDTO.EndDate == "" {
		errorsList = append(errorsList, "You need to select an end date.")
	}
	if availabilityDTO.Price == "" {
		errorsList = append(errorsList, "You need to enter a price.")
	}
	if availabilityDTO.MaxParticipants <= 0 {
		errorsList = append(errorsList, "You need to enter a positive number of max participants.")
	}

	startDate, err := time.Parse("2006-01-02", availabilityDTO.StartDate)
	if err != nil {
		errorsList = append(errorsList, "Start date is invalid")
	}

	endDate, err := time.Parse("2006-01-02", availabilityDTO.EndDate)
	if err != nil {
		errorsList = append(errorsList, "End date is invalid")
	}

	if endDate.Before(startDate) {
    errorsList = append(errorsList, "Start date must be before end date")
}


	daysBitmap, err := DaysSliceToBitmap(availabilityDTO.Days)
	if err != nil {
		errorsList = append(errorsList, err.Error())
	}

	hoursSlice, err := HoursSliceToTimeOfDaySlice(availabilityDTO.Hours)
	if err != nil {
		errorsList = append(errorsList, err.Error())
	}

	price, err := StringPriceToInt(availabilityDTO.Price)
	if err != nil {
		errorsList = append(errorsList, "Price is invalid")
	}

	availabilityID, err := uuid.Parse(availabilityDTO.AvailabilityID)

	if err != nil {
		availabilityID = uuid.Nil
	}

	if len(errorsList) > 0 {
		return nil, &service.ValidationError{
			ErrorsList: errorsList,
		}
	}

	return &model.Availability{
		ID:              availabilityID,
		StartDate:       startDate,
		EndDate:         endDate,
		Days:            daysBitmap,
		Hours:           hoursSlice,
		Price:           int32(price),
		MaxParticipants: availabilityDTO.MaxParticipants,
		Duration:        time.Duration(availabilityDTO.Duration) * time.Minute,
        Precedance:      availabilityDTO.Precedance,
        Notes:           availabilityDTO.Notes,
	}, nil
}

func (m *AvailabilityDTOMapper) AvailabilityToDTO(
    avail *model.Availability,
) AvailabilityDTO {
    // Convert the bitmask days back into a slice
    days := make([]int32, 0)
    for i := 0; i < 7; i++ {
        if avail.Days&(1<<i) != 0 {
            days = append(days, int32(6-i))
        }
    }
    hours := make([]string, 0, len(avail.Hours))
    for _, h := range avail.Hours {
        hours = append(hours, fmt.Sprintf("%02d:%02d", h.Hour, h.Minute))
    }
    price := fmt.Sprintf("%d.%02d", avail.Price/100, avail.Price%100)

    return AvailabilityDTO{
        AvailabilityID:  avail.ID.String(),
        StartDate:       avail.StartDate.Format("2006-01-02"),
        EndDate:         avail.EndDate.Format("2006-01-02"),
        Days:            days,
        Hours:           hours,
        Price:           price,
        MaxParticipants: avail.MaxParticipants,
		Duration:        int32(avail.Duration.Minutes()),
        Precedance:      avail.Precedance,
        Notes:           avail.Notes,
    }
}


func MapOpeningAvailabilityToAPIItem(a model.OpeningAvailability) AvailabilityItem {
	return AvailabilityItem{
		DateTime:  a.DateTime.Format(time.RFC3339),
		Vacancies: a.Vacancies,
		Price:     a.Price,
	}
}
