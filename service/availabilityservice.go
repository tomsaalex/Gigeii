package service

import (
	"context"
	"errors"
	"fmt"

	"example.com/db"
	"example.com/model"
	"example.com/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AvailabilityService struct {
	availabilityRepo repository.AvailabilityRepository
	connPool         *pgxpool.Pool
	queries          *db.Queries
}

func NewAvailabilityService(
	availabilityRepo repository.AvailabilityRepository,
	connPool *pgxpool.Pool,
	queries *db.Queries,
) *AvailabilityService {
	return &AvailabilityService{
		connPool:         connPool,
		queries:          queries,
		availabilityRepo: availabilityRepo,
	}
}

func (s *AvailabilityService) AddAvailability(
	ctx context.Context,
	availability model.Availability,
	precAvailabilityID uuid.UUID,
	conflictResolutionMode bool,
) (*model.Availability, []model.Availability, error) {
	tx, err := s.connPool.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	conflictingAvailabilities, err := s.availabilityRepo.GetConflictingAvailabilities(ctx, qtx, availability)
	if err != nil {
		var re *repository.RepositoryError
		if errors.As(err, &re) {
			return nil, nil, re
		}
	}

	if len(conflictingAvailabilities) != 0 {
		if !conflictResolutionMode {
			return nil, conflictingAvailabilities, &UnhandledConflictError{
				Message: "Unhandled availability conflicts found",
			}
		}

		if precAvailabilityID != uuid.Nil {
			foundAvailability, err := s.availabilityRepo.GetByID(ctx, qtx, precAvailabilityID)

			if err != nil {
				return nil, conflictingAvailabilities, &repository.EntityNotFoundError{
					Message: fmt.Sprintf(
						"Couldn't find the Availability with ID for conflict resolution: %s",
						precAvailabilityID,
					),
				}
			}

			availability.Precedance = foundAvailability.Precedance + 1
		} else {
			availability.Precedance = 0
		}

		s.availabilityRepo.ShiftPrecedenceAbove(ctx, qtx, availability.Precedance-1)
	}

	addedAvailability, err := s.availabilityRepo.Add(ctx, qtx, availability)
	if err != nil {
		return nil, nil, err
	}

	return addedAvailability, conflictingAvailabilities, tx.Commit(ctx)
}

func (s *AvailabilityService) UpdateAvailability(
	ctx context.Context,
	availability model.Availability,
	precAvailabilityID uuid.UUID,
	conflictResolutionMode bool,
) (*model.Availability, []model.Availability, error) {
	tx, err := s.connPool.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	conflictingAvailabilities, err := s.availabilityRepo.GetConflictingAvailabilities(ctx, qtx, availability)
	if err != nil {
		var re *repository.RepositoryError
		if errors.As(err, &re) {
			return nil, nil, re
		}
	}

	if len(conflictingAvailabilities) != 0 {
		if !conflictResolutionMode {
			return nil, conflictingAvailabilities, &UnhandledConflictError{
				Message: "Unhandled availability conflicts found",
			}
		}

		if precAvailabilityID != uuid.Nil {
			foundAvailability, err := s.availabilityRepo.GetByID(ctx, qtx, precAvailabilityID)

			if err != nil {
				return nil, conflictingAvailabilities, &repository.EntityNotFoundError{
					Message: fmt.Sprintf(
						"Couldn't find the Availability with ID for conflict resolution: %s",
						precAvailabilityID,
					),
				}
			}

			availability.Precedance = foundAvailability.Precedance + 1
		} else {
			availability.Precedance = 0
		}

		s.availabilityRepo.ShiftPrecedenceAbove(ctx, qtx, availability.Precedance-1)
	}

	updatedAvailability, err := s.availabilityRepo.Update(ctx, qtx, availability)

	if err != nil {
		return nil, nil, err
	}

	return updatedAvailability, conflictingAvailabilities, tx.Commit(ctx)
}

func (s *AvailabilityService) DeleteAvailability(
	ctx context.Context,
	availabilityID uuid.UUID,
) (*model.Availability, error) {
	return s.availabilityRepo.Delete(ctx, nil, availabilityID)
}
