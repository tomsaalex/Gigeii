package repository

import (
	"context"
	"fmt"

	"example.com/db"
	"example.com/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type DBAvailabilityRepository struct {
	queries *db.Queries
	mapper  AvailabilityMapperDB
}

func (r *DBAvailabilityRepository) Add(
	ctx context.Context,
	availability model.Availability,
) (*model.Availability, error) {
	addAvailabilityParams := r.mapper.AvailabilityToAddAvailabilityParams(availability)
	dbAvailability, err := r.queries.CreateAvailability(ctx, addAvailabilityParams)

	if err != nil {
		return nil, &DuplicateEntityError{
			Message: fmt.Sprintf("Add: Availability for product \"%s\" is already taken.", availability.ProductID),
		}
	}

	modelAvailability := r.mapper.DBAvailabilityToAvailability(dbAvailability)
	return modelAvailability, nil
}

func (r *DBAvailabilityRepository) GetByProduct(
	ctx context.Context,
	productId uuid.UUID,
) ([]model.Availability, error) {
	dbUUID := pgtype.UUID{
		Bytes: productId,
		Valid: true,
	}
	dbAvailabilityList, err := r.queries.GetAvailabilitiesByProduct(ctx, dbUUID)
	if err != nil {
		return nil, &EntityNotFoundError{
			Message: fmt.Sprintf("GetByProduct: No availability matches product \"%s\"", productId),
		}
	}

	var modelAvailabilityList []model.Availability
	for _, dbAvailability := range dbAvailabilityList {
		modelAvailability := r.mapper.DBAvailabilityToAvailability(dbAvailability)
		modelAvailabilityList = append(modelAvailabilityList, *modelAvailability)
	}

	return modelAvailabilityList, nil
}
