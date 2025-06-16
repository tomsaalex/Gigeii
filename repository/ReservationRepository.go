package repository

import (
	"context"

	"example.com/db"
	"example.com/model"
	"github.com/jackc/pgx/v5/pgxpool"
)


type ReservationRepository interface {
	ReserveOrUpdate(ctx context.Context, reservation model.Reservation)(*model.Reservation, error)
	GetByReservationReference(ctx context.Context, reservationReference string) (*model.Reservation, error)
	GetByExternalReservationReference(ctx context.Context, reseller_id string,externalreservationReference string) (*model.Reservation, error)
	CancelReservation(ctx context.Context, reservationReference string)(error)
}




func NewDbReservationRepository(queries *db.Queries, connPool *pgxpool.Pool) ReservationRepository {
	return &DBReservationRepository{
		queries:  queries,
		mapper:   &ReservationMapper{},
		connPool: connPool,
	}
}

