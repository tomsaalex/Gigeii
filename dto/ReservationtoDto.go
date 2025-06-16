package dto

import (
	"example.com/model"
	"github.com/google/uuid"
)

func ReservationToDTO(res *model.Reservation) *ReservationDTO {
	return &ReservationDTO{
		ReservationReference:         res.ReservationReference,
		ExternalReservationReference: res.ExternalReservationReference,
		ResellerID:                   res.ResellerID.String(),
		AvailabilityID:               res.AvailabilityID.String(),
		DateTime:                     res.DateTime,
		Quantity:                     res.Quantity,
		Status:                       res.Status,
	}
}

func (dto ReservationDTO) ToModel() (model.Reservation, error) {
	availabilityID, err := uuid.Parse(dto.AvailabilityID)
	if err != nil {
		return model.Reservation{}, err
	}
	resellerID, err := uuid.Parse(dto.ResellerID)
	if err != nil {
		return model.Reservation{}, err
	}

	return model.Reservation{
		ReservationReference:         dto.ReservationReference,
		ExternalReservationReference: dto.ExternalReservationReference,
		ResellerID:                   resellerID,
		AvailabilityID:               availabilityID,
		DateTime:                     dto.DateTime,
		Quantity:                     dto.Quantity,
		Status:                       dto.Status,
	}, nil
}
