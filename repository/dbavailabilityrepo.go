package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"example.com/db"
	"example.com/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBAvailabilityRepository struct {
	queries  *db.Queries
	connPool *pgxpool.Pool
	mapper   AvailabilityMapperDB
}

func (r *DBAvailabilityRepository) Add(
	ctx context.Context,
	qtx *db.Queries,
	availability model.Availability,
) (*model.Availability, error) {
	var queries *db.Queries
	if qtx == nil {
		tx, err := r.connPool.Begin(ctx)
		if err != nil {
			return nil, err
		}
		defer tx.Rollback(ctx)

		queries = r.queries.WithTx(tx)
	} else {
		queries = qtx
	}

	addAvailabilityParams := r.mapper.AvailabilityToAddAvailabilityParams(availability)
	fmt.Println("AddAvailabilityParams:", addAvailabilityParams)

	dbAvailability, err := queries.CreateAvailability(ctx, addAvailabilityParams)

	if err != nil {
		return nil, &RepositoryError{
			Message: fmt.Sprintf("Failed to insert Availability: %s", err.Error()),
		}
	}

	modelAvailability := r.mapper.DBAvailabilityToAvailability(dbAvailability)
	availabilityID, err := uuid.FromBytes(dbAvailability.ID.Bytes[:])
	dbHourParams := r.mapper.AvailabilityToAddAvailabilityHourParams(
		availabilityID,
		availability.Hours,
	)

	for _, hour := range dbHourParams {
		dbAvailabilityTime, err := queries.AddAvailabilityHour(ctx, hour)
		if err != nil {
			return nil, &RepositoryError{
				Message: fmt.Sprintf("Failed to insert Availability's hours: %s", err.Error()),
			}
		}

		modelAvailability.Hours = append(
			modelAvailability.Hours,
			r.mapper.DBAvailabilityHourToTimeOfDay(dbAvailabilityTime),
		)
	}

	return modelAvailability, nil
}

func (r *DBAvailabilityRepository) ShiftPrecedenceAbove(
	ctx context.Context,
	qtx *db.Queries,
	precedenceThreshold int32,
) error {
	var queries *db.Queries
	if qtx == nil {
		queries = r.queries
	} else {
		queries = qtx
	}

	return queries.ShiftPrecedenceAbove(ctx, precedenceThreshold)
}

func (r *DBAvailabilityRepository) GetConflictingAvailabilities(
	ctx context.Context,
	qtx *db.Queries,
	availability model.Availability,
) ([]model.Availability, error) {
	var queries *db.Queries
	if qtx == nil {
		queries = r.queries
	} else {
		queries = qtx
	}
	searchParams := r.mapper.AvailabilityToFindAvailabilityConflictsParams(availability)
	conflictingAvailabilities, err := queries.FindAvailabilityConflicts(ctx, *searchParams)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []model.Availability{}, &EntityNotFoundError{
				Message: "No conflicting availabilities found",
			}
		}
		return nil, &RepositoryError{
			Message: "Check for conflicting availabilities failed",
		}
	}

	modelAvailabilities := r.mapper.AvailabilityConflictsToAvailabilities(ctx, conflictingAvailabilities)
	return modelAvailabilities, nil
}

func (r *DBAvailabilityRepository) GetByID(
	ctx context.Context,
	qtx *db.Queries,
	availabilityID uuid.UUID,
) (*model.Availability, error) {
	var queries *db.Queries
	if qtx == nil {
		queries = r.queries
	} else {
		queries = qtx
	}

	availabilityWithHourRows, err := queries.GetAvailabilityByID(ctx, uuidToPgtype(availabilityID))

	if err != nil || len(availabilityWithHourRows) == 0 {
		return nil, &EntityNotFoundError{
			Message: fmt.Sprintf("No Availability found with ID: %s", availabilityID.String()),
		}
	}

	modelAvailability := r.mapper.DBAvailabilityWithHourToAvailability(availabilityWithHourRows)
	return modelAvailability, nil
}

func (r *DBAvailabilityRepository) Update(
	ctx context.Context,
	qtx *db.Queries,
	availability model.Availability,
) (*model.Availability, error) {
	var queries *db.Queries
	if qtx == nil {
		tx, err := r.connPool.Begin(ctx)
		if err != nil {
			return nil, err
		}
		defer tx.Rollback(ctx)

		queries = r.queries.WithTx(tx)
	} else {
		queries = qtx
	}

	updateAvailabilityParams := r.mapper.AvailabilityToUpdateAvailabilityParams(availability)

	dbAvailability, err := queries.UpdateAvailability(ctx, updateAvailabilityParams)
	if err != nil {
		return nil, &RepositoryError{
			Message: fmt.Sprintf("Availability update failed unexpectedly: %s", err.Error()),
		}
	}

	updatedAvailability := r.mapper.DBAvailabilityToAvailability(dbAvailability)

	err = queries.DeleteAvailabilityHoursForAvailability(ctx, uuidToPgtype(availability.ID))

	if err != nil {
		return nil, &RepositoryError{
			Message: fmt.Sprintf("Availability Hours update failed unexpectedly: %s", err.Error()),
		}
	}

	dbHourParams := r.mapper.AvailabilityToAddAvailabilityHourParams(
		availability.ID,
		availability.Hours,
	)

	for _, hour := range dbHourParams {
		dbAvailabilityTime, err := queries.AddAvailabilityHour(ctx, hour)
		if err != nil {
			return nil, &RepositoryError{
				Message: fmt.Sprintf("Failed to insert Availability's hours: %s", err.Error()),
			}
		}

		updatedAvailability.Hours = append(
			updatedAvailability.Hours,
			r.mapper.DBAvailabilityHourToTimeOfDay(dbAvailabilityTime),
		)
	}

	return updatedAvailability, nil
}

func (r *DBAvailabilityRepository) Delete(
	ctx context.Context,
	qtx *db.Queries,
	availabilityID uuid.UUID,
) (*model.Availability, error) {
	var queries *db.Queries
	if qtx == nil {
		queries = r.queries
	} else {
		queries = qtx
	}

	dbAvailability, err := queries.DeleteAvailability(ctx, uuidToPgtype(availabilityID))

	if err != nil {
		return nil, &RepositoryError{
			Message: "Availability deletion failed unexpectedly",
		}
	}

	modelAvailability := r.mapper.DBAvailabilityToAvailability(dbAvailability)
	return modelAvailability, nil
}

func (r *DBAvailabilityRepository) GetAllAvailabilities(ctx context.Context) ([]model.Availability, error) {
	dbRows, err := r.queries.GetAllAvailabilities(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapper.DBAvailabilitiesToAvailabilities(dbRows), nil
}
func (r *DBAvailabilityRepository) GetAvailabilitiesInRange(
    ctx context.Context,
    from, to string,
) ([]model.OpeningAvailability, error) {
    fromTime, err := time.Parse(time.RFC3339, from)
    if err != nil {
        //fmt.Println("Error parsing from time:", err)
        return nil, fmt.Errorf("invalid from time format: %w", err)
    }

    toTime, err := time.Parse(time.RFC3339, to)
    if err != nil {
        //fmt.Println("Error parsing to time:", err)
        return nil, fmt.Errorf("invalid to time format: %w", err)
    }

    if !toTime.After(fromTime) {
        //fmt.Println("Error: toTime is not after fromTime")
        return nil, fmt.Errorf("toDateTime (%s) must be after fromDateTime (%s)", toTime, fromTime)
    }

    // fmt.Println("Raw fromDateTime:", from)
    // fmt.Println("Raw toDateTime:", to)
    // fmt.Println("Parsed fromTime:", fromTime)
    // fmt.Println("Parsed toTime:", toTime)
    // fmt.Println("fromTime Location:", fromTime.Location())
    // fmt.Println("toTime Location:", toTime.Location())

    // Detect full-day query
    isFullDayQuery := (fromTime.Hour() == 0 && fromTime.Minute() == 0 && fromTime.Second() == 0 &&
        toTime.Hour() == 0 && toTime.Minute() == 0 && toTime.Second() == 0)
    if isFullDayQuery {
        fmt.Println("Full-day query detected, returning all hours in date range")
    }

    // Truncate to date using local date components
    fromDate := time.Date(fromTime.Year(), fromTime.Month(), fromTime.Day(), 0, 0, 0, 0, time.UTC)
    toDate := time.Date(toTime.Year(), toTime.Month(), toTime.Day(), 0, 0, 0, 0, time.UTC)
    //fmt.Println("Truncated fromDate:", fromDate)
    //fmt.Println("Truncated toDate:", toDate)

    var fromDatePg, toDatePg pgtype.Date
    if err := fromDatePg.Scan(fromDate); err != nil {
        fmt.Println("Error scanning from date:", err)
        return nil, fmt.Errorf("failed to scan from date: %w", err)
    }
    if err := toDatePg.Scan(toDate); err != nil {
        fmt.Println("Error scanning to date:", err)
        return nil, fmt.Errorf("failed to scan to date: %w", err)
    }

    // Set time range for query (use local time for hours)
    var fromHourPg, toHourPg pgtype.Time
    if isFullDayQuery {
        fromHourPg = pgtype.Time{Microseconds: 0, Valid: true}
        toHourPg = pgtype.Time{Microseconds: 0, Valid: true}
    } else {
        fromHourPg = pgtype.Time{
            Microseconds: int64(fromTime.Hour())*3600*1000000 + int64(fromTime.Minute())*60*1000000,
            Valid:        true,
        }
        toHourPg = pgtype.Time{
            Microseconds: int64(toTime.Hour())*3600*1000000 + int64(toTime.Minute())*60*1000000,
            Valid:        true,
        }
    }

    //fmt.Println("From Hour (pgtype):", fromHourPg)
    //fmt.Println("To Hour (pgtype):", toHourPg)

    params := db.GetAvailabilitiesInRangeParams{
        
        EndDate:   fromDatePg,   // $2: to date (e.g., 2025-09-16)
		StartDate: toDatePg, // $1: from date (e.g., 2025-09-14)
        Column3:   fromHourPg, // $3: from hour (e.g., 72000000000 for 20:00:00)
        Column4:   toHourPg,   // $4: to hour (e.g., 28800000000 for 08:00:00)
    }
    //fmt.Println("Query Params:", params)

    rows, err := r.queries.GetAvailabilitiesInRange(ctx, params)
    if err != nil {
       // fmt.Println("Error getting availabilities in range:", err)
        return nil, fmt.Errorf("failed to get availabilities in range: %w", err)
    }

	//fmt.Println("Rows returned:", rows)

    resultMap := make(map[string]model.OpeningAvailability)

for _, row := range rows {
	startDateTime := row.StartDate.Time.UTC()
	hourTime := time.Unix(0, int64(row.HourMicroseconds)*1000).UTC()

	dateTime := time.Date(
		startDateTime.Year(), startDateTime.Month(), startDateTime.Day(),
		hourTime.Hour(), hourTime.Minute(), hourTime.Second(), 0, time.UTC,
	)

	// Construim cheia după zi și oră:
	key := dateTime.Format("2006-01-02 15:04")

	existing, exists := resultMap[key]
	if !exists || row.Precedance > existing.Precedance {
		resultMap[key] = model.OpeningAvailability{
			DateTime:  dateTime,
			Vacancies: row.MaxParticipants,
			Price:     row.Price,
			Precedance: row.Precedance, 
		}
	}
}

// convertim map-ul în slice
result := make([]model.OpeningAvailability, 0, len(resultMap))
for _, availability := range resultMap {
	result = append(result, availability)
}

return result, nil
}


/*
func (r *DBAvailabilityRepository) WithSerializableTx(
	ctx context.Context,
	f func(repo AvailabilityRepository) error,
) error {
	tx, err := r.connPool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)
	txRepo := &DBAvailabilityRepository{
		queries:  qtx,
		connPool: r.connPool, // can be omitted
		mapper:   r.mapper,
	}

	if err := f(txRepo); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
*/
