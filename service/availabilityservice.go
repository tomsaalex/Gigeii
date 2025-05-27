package service

import (
	"context"
	"errors"
	"fmt"

	"example.com/model"
	"example.com/repository"
	"github.com/google/uuid"
)

type AvailabilityService struct {
	availabilityRepo repository.AvailabilityRepository
}

func NewAvailabilityService(availabilityRepo repository.AvailabilityRepository) *AvailabilityService {
	return &AvailabilityService{
		availabilityRepo: availabilityRepo,
	}
}

func (s *AvailabilityService) AddAvailability(
	ctx context.Context,
	availability model.Availability,
	precAvailabilityID uuid.UUID,
	conflictResolutionMode bool,
) (*model.Availability, []model.Availability, error) {
	// TODO: Business logic validation.
	conflictingAvailabilities, err := s.availabilityRepo.GetConflictingAvailabilities(ctx, availability)
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
			foundAvailability, err := s.availabilityRepo.GetByID(ctx, precAvailabilityID)

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

		s.availabilityRepo.ShiftPrecedenceAbove(ctx, availability.Precedance-1)
	}

	addedAvailability, err := s.availabilityRepo.Add(ctx, availability)
	return addedAvailability, conflictingAvailabilities, err
}

func (s *AvailabilityService) UpdateAvailability(
	ctx context.Context,
	availability model.Availability,
) (*model.Availability, error) {

	return s.availabilityRepo.Update(ctx, availability)
}

func (s *AvailabilityService) DeleteAvailability(
	ctx context.Context,
	availabilityID uuid.UUID,
) (*model.Availability, error) {
	return s.availabilityRepo.Delete(ctx, availabilityID)
}
