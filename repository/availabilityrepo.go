package repository

import (
	"context"
	"time"

	"example.com/db"
	"example.com/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AvailabilityRepository interface {
	Add(ctx context.Context, queries *db.Queries, availability model.Availability) (*model.Availability, error)
	GetConflictingAvailabilities(
		ctx context.Context,
		queries *db.Queries,
		availability model.Availability,
	) ([]model.Availability, error)
	ShiftPrecedenceAbove(ctx context.Context, queries *db.Queries, precedenceThreshold int32) error
	GetByID(ctx context.Context, queries *db.Queries, availabilityID uuid.UUID) (*model.Availability, error)
	Update(ctx context.Context, queries *db.Queries, availability model.Availability) (*model.Availability, error)
	Delete(ctx context.Context, queries *db.Queries, availabilityID uuid.UUID) (*model.Availability, error)
	GetAllAvailabilities(ctx context.Context) ([]model.Availability, error)

	
	GetAvailabilitiesInRange(ctx context.Context, 	from, to string) ([]model.OpeningAvailability, error)
	GetAvailableVacancies(ctx context.Context, availabilityID uuid.UUID) (int32, error)
	GetAvailabilityIdForReservation(ctx context.Context, queries *db.Queries, date, hour time.Time) (uuid.UUID, error)

	
}

func NewDBAvailabilityRepository(connPool *pgxpool.Pool, queries *db.Queries) AvailabilityRepository {
	return &DBAvailabilityRepository{
		queries:  queries,
		connPool: connPool,
		mapper:   AvailabilityMapperDB{},
	}
}
