package repository

import (
	"context"
	"time"

	"example.com/db"
	"example.com/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AvailabilityMapperDB struct {
}

func uuidToPgtype(id uuid.UUID) pgtype.UUID {
	if id == uuid.Nil {
		return pgtype.UUID{
			Valid: false,
		}
	}

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




func timeOfDayToPgTimestampTz(td model.TimeOfDay) pgtype.Timestamptz {
	timestamp := time.Date(1970, 1, 1, int(td.Hour), int(td.Minute), 0, 0, time.Local)
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
		timestamp := timeOfDayToPgTimestampTz(td)
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
	return db.CreateAvailabilityParams{
		StartDate:       timeToPgtype(availability.StartDate),
		EndDate:         timeToPgtype(availability.EndDate),
		Days:            availability.Days,
		Price:           availability.Price,
		MaxParticipants: availability.MaxParticipants,
		Precedance:      availability.Precedance,
		CreatedBy:       uuidToPgtype(availability.CreatedBy),
		// TODO: Actually add the duration here
		Duration: pgtype.Interval{Microseconds: int64(availability.Duration.Minutes()) * 60 * 1000000, Valid: true},
		Notes:    pgtype.Text{String: availability.Notes, Valid: availability.Notes != ""},
	}
}

func (m *AvailabilityMapperDB) AvailabilityToFindAvailabilityConflictsParams(
	availability model.Availability,
) *db.FindAvailabilityConflictsParams {
	hours := make([]pgtype.Timestamptz, 0)
	for _, hour := range availability.Hours {
		hours = append(hours, timeOfDayToPgTimestampTz(hour))
	}

	return &db.FindAvailabilityConflictsParams{
		StartDate:      timeToPgtype(availability.StartDate),
		EndDate:        timeToPgtype(availability.EndDate),
		Days:           availability.Days,
		Hours:          hours,
		AvailabilityID: uuidToPgtype(availability.ID),
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

func (m *AvailabilityMapperDB) DBAvailabilityWithHourToAvailability(
	availabilityWithHourRows []db.GetAvailabilityByIDRow,
) *model.Availability {
	modelHours := make([]model.TimeOfDay, 0)

	for _, availWithHour := range availabilityWithHourRows {
		extractedTime := availWithHour.Hour.Time
		timeOfDay := model.TimeOfDay{
			Hour:   int32(extractedTime.Hour()),
			Minute: int32(extractedTime.Minute()),
		}
		modelHours = append(modelHours, timeOfDay)
	}
	
	

	return &model.Availability{
		ID:              availabilityWithHourRows[0].AvailabilityID.Bytes,
		StartDate:       availabilityWithHourRows[0].StartDate.Time,
		EndDate:         pgTypeToTime(availabilityWithHourRows[0].EndDate),
		Days:            availabilityWithHourRows[0].Days,
		Price:           availabilityWithHourRows[0].Price,
		MaxParticipants: availabilityWithHourRows[0].MaxParticipants,
		Precedance:      availabilityWithHourRows[0].Precedance,
		CreatedBy:       availabilityWithHourRows[0].CreatedBy.Bytes,
		Hours:           modelHours,
		Duration:        time.Duration(availabilityWithHourRows[0].Duration.Microseconds) * time.Microsecond,
		Notes:           availabilityWithHourRows[0].Notes.String,
	}
}

func (m *AvailabilityMapperDB) AvailabilityConflictsToAvailabilities(
	ctx context.Context,
	dbAvailabilities []db.FindAvailabilityConflictsRow,
) []model.Availability {
	availabilityMap := make(map[uuid.UUID]*model.Availability)

	for _, row := range dbAvailabilities {
		availabilityID := row.ID.Bytes

		if _, exists := availabilityMap[availabilityID]; !exists {
			availabilityMap[availabilityID] = &model.Availability{
				ID:              row.ID.Bytes,
				StartDate:       pgTypeToTime(row.StartDate),
				EndDate:         pgTypeToTime(row.EndDate),
				Days:            row.Days,
				Price:           row.Price,
				MaxParticipants: row.MaxParticipants,
				Precedance:      row.Precedance,
				CreatedBy:       row.CreatedBy.Bytes,
				Hours:           []model.TimeOfDay{},
			}
		}

		hour := row.Hour.Time
		timeOfDay := model.TimeOfDay{
			Hour:   int32(hour.Hour()),
			Minute: int32(hour.Minute()),
		}
		availabilityMap[availabilityID].Hours = append(availabilityMap[availabilityID].Hours, timeOfDay)
	}

	availabilities := make([]model.Availability, 0, len(availabilityMap))
	for _, availability := range availabilityMap {
		availabilities = append(availabilities, *availability)
	}

	return availabilities
}

func (m *AvailabilityMapperDB) AvailabilityToUpdateAvailabilityParams(
	availability model.Availability,
) db.UpdateAvailabilityParams {
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
		Duration: pgtype.Interval{Microseconds: int64(availability.Duration.Minutes()) * 60 * 1000000, Valid: true},
		Notes:    pgtype.Text{String: availability.Notes, Valid: availability.Notes != ""},
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
		Duration:        time.Duration(dbAvailability.Duration.Microseconds) * time.Microsecond,
		Notes:           dbAvailability.Notes.String,
	}

	return &availability
}

func (m *AvailabilityMapperDB) DBAvailabilitiesToAvailabilities(
	dbRows []db.GetAllAvailabilitiesRow,
) []model.Availability {
	availabilityMap := make(map[uuid.UUID]*model.Availability)
	for _, row := range dbRows {
		id := row.ID
		uuidID, err := uuid.FromBytes(id.Bytes[:])
		if err != nil {
			// handle error
		}
		avail, exists := availabilityMap[uuidID]
		if !exists {
			avail = &model.Availability{
				ID:              uuidID,
				StartDate:       row.StartDate.Time,
				EndDate:         row.EndDate.Time,
				Days:            row.Days,
				Price:           row.Price,
				MaxParticipants: row.MaxParticipants,
				Precedance:      row.Precedance,
				CreatedBy:       uuid.Must(uuid.FromBytes(row.CreatedBy.Bytes[:])),
				Duration:        time.Duration(row.Duration.Microseconds * 1000),
				Notes:           row.Notes.String,
				Hours:           []model.TimeOfDay{},
			}
			availabilityMap[uuidID] = avail
		}
		// Always append the hour, if not NULL
		if row.Hour.Valid {
			hour := row.Hour.Time.Hour()
			minute := row.Hour.Time.Minute()
			avail.Hours = append(avail.Hours, model.TimeOfDay{Hour: int32(hour), Minute: int32(minute)})
		}
	}
	// Flatten map to slice
	result := make([]model.Availability, 0, len(availabilityMap))
	for _, avail := range availabilityMap {
		result = append(result, *avail)
	}
	return result
}

