package repository

import (
	"context"

	"example.com/db"
	"example.com/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBReservationRepository struct {
	queries  *db.Queries
	mapper   *ReservationMapper
	connPool *pgxpool.Pool
}



// GetByExternalReservationReference implements ReservationRepository.
func (d *DBReservationRepository) GetByExternalReservationReference(ctx context.Context, reseller_id string, externalreservationReference string) (*model.Reservation, error) {
	uuidReseller, err := uuid.Parse(reseller_id)
	if err != nil {
		return nil, err
	}

	params := db.GetReservationByExternalReferenceParams{
		ResellerID: pgtype.UUID{Bytes: uuidReseller, Valid: true},
		ExternalReservationReference: externalreservationReference,
	}

	row, err := d.queries.GetReservationByExternalReference(ctx, params)
	if err != nil {
		return nil, err
	}

	result := d.mapper.DBGetReservationByExternalReferenceToModel(row)
	return &result, nil
}

// GetByReservationReference implements ReservationRepository.
func (d *DBReservationRepository) GetByReservationReference(ctx context.Context, reservationReference string) (*model.Reservation, error) {
	row, err := d.queries.GetReservationByReference(ctx, reservationReference)
	if err != nil {
		return nil, err
	}
	result := d.mapper.DBGetReservationByReferenceToModel(row)
	return &result, nil
}

// ReserveOrUpdate implements ReservationRepository.
func (d *DBReservationRepository) ReserveOrUpdate(ctx context.Context, reservation model.Reservation) (*model.Reservation, error) {
	params := d.mapper.ModelToReserveOrUpdateParams(reservation)

	row, err := d.queries.ReserveOrUpdateReservation(ctx, params)
	if err != nil {
		return nil, err
	}

	result := d.mapper.DBReserveOrUpdateReservationToModel(row)
	return &result, nil
}
// CancelReservation implements ReservationRepository.
func (d *DBReservationRepository) CancelReservation(ctx context.Context, reservationReference string) error {
	_, err := d.queries.CancelReservation(ctx, reservationReference)
	return err	
}