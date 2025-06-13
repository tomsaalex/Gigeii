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
		// aflam precedenta maxima din conflicte
		maxPrecedence := 0
		for _, c := range conflictingAvailabilities {
			if c.Precedance > int32(maxPrecedence) {
				maxPrecedence = int(c.Precedance)
			}
		}
		availability.Precedance = int32(maxPrecedence + 1)

		// facem shift daca vrei, sau sari peste daca nu vrei shift (aici nu e nevoie neaparat)
	} else {
		availability.Precedance = 0
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

	// Calculăm precedența corect
	maxPrecedence := 0
	for _, c := range conflictingAvailabilities {
		if c.Precedance > int32(maxPrecedence) {
			maxPrecedence = int(c.Precedance)
		}
	}
	availability.Precedance = int32(maxPrecedence + 1)

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
func (s *AvailabilityService) GetAllAvailabilities(ctx context.Context) ([]model.Availability, error) {
	return s.availabilityRepo.GetAllAvailabilities(ctx)
}

func (s *AvailabilityService) GetAvailabilityByID(ctx context.Context, id uuid.UUID) (*model.Availability, error) {
	return s.availabilityRepo.GetByID(ctx, nil, id)
}

func (s *AvailabilityService) GetAvailabilitiesInRange(
	ctx context.Context,
	from, to string,
) ([]model.OpeningAvailability, error) {
	_,err:=s.availabilityRepo.GetAvailabilitiesInRange(ctx, from, to)
	if err != nil {
		fmt.Println("Error getting availabilities service in range:", err)
		return nil,err
	}

	return s.availabilityRepo.GetAvailabilitiesInRange(ctx, from, to)
}


