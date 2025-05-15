package repository

import (
	"context"

	"example.com/db"
	"example.com/model"
	"github.com/google/uuid"
)

type AvailabilityRepository interface {
	Add(ctx context.Context, availability model.Availability) (*model.Availability, error)
	GetByProduct(ctx context.Context, productId uuid.UUID) ([]model.Availability, error)
}

func NewDBAvailabilityRepository(queries *db.Queries) AvailabilityRepository {
	return &DBAvailabilityRepository{
		queries: queries,
		mapper:  AvailabilityMapperDB{},
	}
}
