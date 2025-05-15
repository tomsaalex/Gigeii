package repository

import (
	"time"

	"example.com/db"
	"example.com/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AvailabilityMapperDB struct {
}

func uuidToPgtype(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
}

func timeToPgtype(date time.Time) pgtype.Date {
	return pgtype.Date{
		Time:  date,
		Valid: true,
	}
}

func pgTypeToTime(date pgtype.Date) time.Time {
	return date.Time
}

func (m *AvailabilityMapperDB) AvailabilityToAddAvailabilityParams(
	availability model.Availability,
) db.CreateAvailabilityParams {
	return db.CreateAvailabilityParams{
		ProductID:        uuidToPgtype(availability.ProductID),
		StartDate:        timeToPgtype(availability.StartDate),
		EndDate:          timeToPgtype(availability.EndDate),
		AvailabilityType: availability.AvailabilityType,
		Days:             availability.Days,
		Hours:            availability.Hours,
		Price:            availability.Price,
		MaxParticipants:  availability.MaxParticipants,
		Precedance:       availability.Precedance,
		CreatedBy:        uuidToPgtype(availability.CreatedBy),
	}
}

func (m *AvailabilityMapperDB) DBAvailabilityToAvailability(
	dbAvailability db.Availability,
) *model.Availability {
	availability := model.Availability{
		Id:               dbAvailability.ID.Bytes,
		ProductID:        dbAvailability.ProductID.Bytes,
		StartDate:        pgTypeToTime(dbAvailability.StartDate),
		EndDate:          pgTypeToTime(dbAvailability.EndDate),
		AvailabilityType: dbAvailability.AvailabilityType,
		Days:             dbAvailability.Days,
		Hours:            dbAvailability.Hours,
		Price:            dbAvailability.Price,
		MaxParticipants:  dbAvailability.MaxParticipants,
		Precedance:       dbAvailability.Precedance,
		CreatedBy:        dbAvailability.CreatedBy.Bytes,
	}

	return &availability
}
