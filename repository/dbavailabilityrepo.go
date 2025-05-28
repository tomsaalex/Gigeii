package repository

import (
	"context"
	"errors"
	"fmt"

	"example.com/db"
	"example.com/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
