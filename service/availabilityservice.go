package service

import (
	"context"

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
) (*model.Availability, error) {
	// TODO: Business logic validation.

	if precAvailabilityID != uuid.Nil {
		foundAvailability, err := s.availabilityRepo.GetByID(ctx, precAvailabilityID)

		if err != nil {
			// TODO: Handle later
		}

		if foundAvailability != nil {

		}
	}

	return s.availabilityRepo.Add(ctx, availability)
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
