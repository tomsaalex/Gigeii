package repository

import (
	"context"

	"example.com/db"
	"example.com/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AvailabilityRepository interface {
	Add(ctx context.Context, availability model.Availability) (*model.Availability, error)
	Update(ctx context.Context, availability model.Availability) (*model.Availability, error)
	Delete(ctx context.Context, availabilityID uuid.UUID) (*model.Availability, error)
}

func NewDBAvailabilityRepository(connPool *pgxpool.Pool, queries *db.Queries) AvailabilityRepository {
	return &DBAvailabilityRepository{
		queries:  queries,
		connPool: connPool,
		mapper:   AvailabilityMapperDB{},
	}
}
