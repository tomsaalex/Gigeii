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

func timeOfDayToPgType(td model.TimeOfDay) pgtype.Timestamptz {
	timestamp := time.Date(1970, 1, 1, int(td.Hour), int(td.Minute), 0, 0, time.UTC)
	return pgtype.Timestamptz{
		Time:  timestamp,
		Valid: true,
	}
}

func (m *AvailabilityMapperDB) AvailabilityToAddAvailabilityHourParams(
	availabilityID uuid.UUID, timesOfDay []model.TimeOfDay,
) []db.AddAvailabilityHourParams {
	dbAddHourParams := make([]db.AddAvailabilityHourParams, 0)
	for _, td := range timesOfDay {
		timestamp := timeOfDayToPgType(td)
		addHourParam := db.AddAvailabilityHourParams{
			AvailabilityID: uuidToPgtype(availabilityID),
			Hour:           timestamp,
		}

		dbAddHourParams = append(dbAddHourParams, addHourParam)
	}
	return dbAddHourParams
}

func (m *AvailabilityMapperDB) AvailabilityToAddAvailabilityParams(
	availability model.Availability,
) db.CreateAvailabilityParams {
	// TODO: This is hardcoded to john doe until we have middleware implemented
	availability.CreatedBy, _ = uuid.Parse("a86b062a-56ad-4805-9963-2c8b0ad79999")
	return db.CreateAvailabilityParams{
		StartDate:       timeToPgtype(availability.StartDate),
		EndDate:         timeToPgtype(availability.EndDate),
		Days:            availability.Days,
		Price:           availability.Price,
		MaxParticipants: availability.MaxParticipants,
		Precedance:      availability.Precedance,
		CreatedBy:       uuidToPgtype(availability.CreatedBy),
		// TODO: Actually add the duration here
		Duration: pgtype.Interval{Months: 0, Valid: true},
	}
}

func (m *AvailabilityMapperDB) DBAvailabilityHourToTimeOfDay(
	dbHour db.AvailabilityHour,
) model.TimeOfDay {
	extractedTime := dbHour.Hour.Time
	return model.TimeOfDay{
		Hour:   int32(extractedTime.Hour()),
		Minute: int32(extractedTime.Minute()),
	}
}

func (m *AvailabilityMapperDB) AvailabilityToUpdateAvailabilityParams(
	availability model.Availability,
) db.UpdateAvailabilityParams {
	// TODO: This is hardcoded to john doe until we have middleware implemented
	availability.CreatedBy, _ = uuid.Parse("a86b062a-56ad-4805-9963-2c8b0ad79999")
	return db.UpdateAvailabilityParams{
		ID:              uuidToPgtype(availability.ID),
		StartDate:       timeToPgtype(availability.StartDate),
		EndDate:         timeToPgtype(availability.EndDate),
		Days:            availability.Days,
		Price:           availability.Price,
		MaxParticipants: availability.MaxParticipants,
		Precedance:      availability.Precedance,
		CreatedBy:       uuidToPgtype(availability.CreatedBy),
		// TODO: Duration here as well
		Duration: pgtype.Interval{Months: 0, Valid: true},
	}
}

func (m *AvailabilityMapperDB) DBAvailabilityToAvailability(
	dbAvailability db.Availability,
) *model.Availability {
	availability := model.Availability{
		ID:              dbAvailability.ID.Bytes,
		StartDate:       pgTypeToTime(dbAvailability.StartDate),
		EndDate:         pgTypeToTime(dbAvailability.EndDate),
		Days:            dbAvailability.Days,
		Price:           dbAvailability.Price,
		MaxParticipants: dbAvailability.MaxParticipants,
		Precedance:      dbAvailability.Precedance,
		CreatedBy:       dbAvailability.CreatedBy.Bytes,
	}

	return &availability
}
