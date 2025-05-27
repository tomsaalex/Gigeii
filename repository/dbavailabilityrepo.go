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
	availability model.Availability,
) (*model.Availability, error) {
	tx, err := r.connPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	addAvailabilityParams := r.mapper.AvailabilityToAddAvailabilityParams(availability)

	dbAvailability, err := qtx.CreateAvailability(ctx, addAvailabilityParams)

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
		dbAvailabilityTime, err := qtx.AddAvailabilityHour(ctx, hour)
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

	return modelAvailability, tx.Commit(ctx)
}

func (r *DBAvailabilityRepository) ShiftPrecedenceAbove(ctx context.Context, precedenceThreshold int32) error {
	return r.queries.ShiftPrecedenceAbove(ctx, precedenceThreshold)
}

func (r *DBAvailabilityRepository) GetConflictingAvailabilities(
	ctx context.Context,
	availability model.Availability,
) ([]model.Availability, error) {
	searchParams := r.mapper.AvailabilityToFindAvailabilityConflictsParams(availability)
	conflictingAvailabilities, err := r.queries.FindAvailabilityConflicts(ctx, *searchParams)

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

func (r *DBAvailabilityRepository) GetByID(ctx context.Context, availabilityID uuid.UUID) (*model.Availability, error) {
	availabilityWithHourRows, err := r.queries.GetAvailabilityByID(ctx, uuidToPgtype(availabilityID))

	if err != nil {
		return nil, &EntityNotFoundError{
			Message: fmt.Sprintf("No Availability found with ID: %s", availabilityID.String()),
		}
	}

	modelAvailability := r.mapper.DBAvailabilityWithHourToAvailability(availabilityWithHourRows)
	return modelAvailability, nil
}

func (r *DBAvailabilityRepository) Update(
	ctx context.Context,
	availability model.Availability,
) (*model.Availability, error) {
	tx, err := r.connPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	updateAvailabilityParams := r.mapper.AvailabilityToUpdateAvailabilityParams(availability)
	dbAvailability, err := qtx.UpdateAvailability(ctx, updateAvailabilityParams)
	if err != nil {
		return nil, &RepositoryError{
			Message: fmt.Sprintf("Availability update failed unexpectedly: %s", err.Error()),
		}
	}

	updatedAvailability := r.mapper.DBAvailabilityToAvailability(dbAvailability)

	err = qtx.DeleteAvailabilityHoursForAvailability(ctx, uuidToPgtype(availability.ID))

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
		dbAvailabilityTime, err := qtx.AddAvailabilityHour(ctx, hour)
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

	return updatedAvailability, tx.Commit(ctx)
}

func (r *DBAvailabilityRepository) Delete(ctx context.Context, availabilityID uuid.UUID) (*model.Availability, error) {
	dbAvailability, err := r.queries.DeleteAvailability(ctx, uuidToPgtype(availabilityID))

	if err != nil {
		return nil, &RepositoryError{
			Message: "Availability deletion failed unexpectedly",
		}
	}

	modelAvailability := r.mapper.DBAvailabilityToAvailability(dbAvailability)
	return modelAvailability, nil
}
