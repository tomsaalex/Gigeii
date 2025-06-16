package repository

import (
	"time"

	"example.com/db"
	"example.com/model"
	"github.com/jackc/pgx/v5/pgtype"
)

type ReservationMapper struct{}


// Convert from db.ReserveOrUpdateReservationRow (sqlc) to model.Reservation
func (m *ReservationMapper) DBReserveOrUpdateReservationToModel(row db.ReserveOrUpdateReservationRow) model.Reservation {
	return model.Reservation{
		ID:                           row.ID.Bytes,
		ReservationReference:         row.ReservationReference,
		ExternalReservationReference: row.ExternalReservationReference,
		ResellerID:                   row.ResellerID.Bytes,
		AvailabilityID:               row.AvailabilityID.Bytes,
		DateTime:                     row.DateTime.Time,
		Quantity:                     row.Quantity,
		Status:                       row.Status,
		CreatedAt:                    row.CreatedAt.Time,
		UpdatedAt:                    row.UpdatedAt.Time,
	}
}

func (m *ReservationMapper) DBCancelReservationToModel(row db.CancelReservationRow) model.Reservation {
	return model.Reservation{
		ID:                           row.ID.Bytes,
		ReservationReference:         row.ReservationReference,
		ExternalReservationReference: row.ExternalReservationReference,
		ResellerID:                   row.ResellerID.Bytes,
		AvailabilityID:               row.AvailabilityID.Bytes,
		DateTime:                     row.DateTime.Time,
		Quantity:                     row.Quantity,
		Status:                       row.Status,
		CreatedAt:                    row.CreatedAt.Time,
		UpdatedAt:                    row.UpdatedAt.Time,
	}
}

func (m *ReservationMapper) DBGetReservationByReferenceToModel(row db.GetReservationByReferenceRow) model.Reservation {
	return model.Reservation{
		ID:                           row.ID.Bytes,
		ReservationReference:         row.ReservationReference,
		ExternalReservationReference: row.ExternalReservationReference,
		ResellerID:                   row.ResellerID.Bytes,
		AvailabilityID:               row.AvailabilityID.Bytes,
		DateTime:                     row.DateTime.Time,
		Quantity:                     row.Quantity,
		Status:                       row.Status,
		CreatedAt:                    row.CreatedAt.Time,
		UpdatedAt:                    row.UpdatedAt.Time,
	}
}

func (m *ReservationMapper) DBGetReservationByExternalReferenceToModel(row db.GetReservationByExternalReferenceRow) model.Reservation {
	return model.Reservation{
		ID:                           row.ID.Bytes,
		ReservationReference:         row.ReservationReference,
		ExternalReservationReference: row.ExternalReservationReference,
		ResellerID:                   row.ResellerID.Bytes,
		AvailabilityID:               row.AvailabilityID.Bytes,
		DateTime:                     row.DateTime.Time,
		Quantity:                     row.Quantity,
		Status:                       row.Status,
		CreatedAt:                    row.CreatedAt.Time,
		UpdatedAt:                    row.UpdatedAt.Time,
	}
}


func (m *ReservationMapper) ModelToReserveOrUpdateParams(res model.Reservation) db.ReserveOrUpdateReservationParams {
	return db.ReserveOrUpdateReservationParams{
		ReservationReference:         res.ReservationReference,
		ExternalReservationReference: res.ExternalReservationReference,
		ResellerID:                   uuidToPgtype(res.ResellerID),
		AvailabilityID:               uuidToPgtype(res.AvailabilityID),
		DateTime:                     timestamptzToPgtype(res.DateTime),
		Quantity:                     res.Quantity,
	}
}



func timestamptzToPgtype(t time.Time) (out pgtype.Timestamptz) {
	out.Valid = true
	out.Time = t
	return
}