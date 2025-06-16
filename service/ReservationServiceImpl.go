package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"example.com/db"
	"example.com/dto"
	"example.com/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationServiceImpl struct {
	repo         repository.ReservationRepository
	availability repository.AvailabilityRepository
	dbPool       *pgxpool.Pool
}


// CancelReservation implements ReservationService.
func (s *ReservationServiceImpl) CancelReservation(ctx context.Context, reservationReference string, externalReservationReference string) error {
	tx, err := s.dbPool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.New(tx)
	txRepo := repository.NewDbReservationRepository(qtx, s.dbPool)

	//caut rezervarea dupa reference
	existing, err := txRepo.GetByReservationReference(ctx, reservationReference)
	if err != nil {
		return errors.New("INVALID_RESERVATION")
	}

	// Verific daca externalReservationReference corespunde
	if existing.ExternalReservationReference != externalReservationReference {
		return errors.New("INVALID_RESERVATION")
	}

	//  Dacă deja e CANCELlED
	if existing.Status == "CANCELED" {
		return errors.New("RESERVATION_ALREADY_CANCELED")
	}

	// Fac cancelul
	err = txRepo.CancelReservation(ctx, reservationReference)
	if err != nil {
		return fmt.Errorf("failed to cancel reservation: %w", err)
	}

	//Commit
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit error: %w", err)
	}

	return nil
}

// GetByExternalReservationReference implements ReservationService.
func (r *ReservationServiceImpl) GetByExternalReservationReference(ctx context.Context, resellerID string, externalReservationReference string) (*dto.ReservationDTO, error) {
	reservation, err := r.repo.GetByExternalReservationReference(ctx, resellerID, externalReservationReference)
	if err != nil {
		return nil,err
	}
	return dto.ReservationToDTO(reservation), nil

}

// GetByReservationReference implements ReservationService.
func (r *ReservationServiceImpl) GetByReservationReference(ctx context.Context, reservationReference string) (*dto.ReservationDTO, error) {
	reservation,err:= r.repo.GetByReservationReference(ctx, reservationReference)
	if err != nil {
		return nil,err}

	return dto.ReservationToDTO(reservation), nil
}

// ReserveOrUpdate implements ReservationService.
func (s *ReservationServiceImpl) ReserveOrUpdate(ctx context.Context, reservation dto.ReservationDTO) (*dto.ReservationDTO, error) {
	
	// conversie DTO -> model
	modelReservation, err := reservation.ToModel()
	if err != nil {
		return nil, err
	}

	const maxAttempts = 10

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		tx, err := s.dbPool.BeginTx(ctx, pgx.TxOptions{
			IsoLevel: pgx.Serializable,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to begin transaction: %w", err)
		}

		qtx := db.New(tx)
		txRepo := repository.NewDbReservationRepository(qtx, s.dbPool)
		txAvailabilityRepo := repository.NewDBAvailabilityRepository(s.dbPool, qtx)

		// verificam daca exista deja rezervare cu externalReservationReference
		existing, err := txRepo.GetByExternalReservationReference(ctx, modelReservation.ResellerID.String(), modelReservation.ExternalReservationReference)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			tx.Rollback(ctx)
			return nil, fmt.Errorf("db error: %w", err)
		}

		if existing != nil {
			// verificam daca e aceeasi rezervare
			if existing.DateTime.Equal(modelReservation.DateTime) &&
				existing.Quantity == modelReservation.Quantity &&
				existing.AvailabilityID == modelReservation.AvailabilityID {
				if err := tx.Commit(ctx); err != nil {
					return nil, fmt.Errorf("commit error: %w", err)
				}
				return dto.ReservationToDTO(existing), nil
			}

			// daca e diferita, o anulam
			if err := txRepo.CancelReservation(ctx, existing.ReservationReference); err != nil {
				tx.Rollback(ctx)
				return nil, fmt.Errorf("failed to cancel old reservation: %w", err)
			}
		}

		// verificam locurile disponibile
		availableSpots, err := txAvailabilityRepo.GetAvailableVacancies(ctx, modelReservation.AvailabilityID)
		if err != nil {
			tx.Rollback(ctx)
			return nil, fmt.Errorf("failed to check availability: %w", err)
		}
		if availableSpots < modelReservation.Quantity {
			tx.Rollback(ctx)
			return nil, fmt.Errorf("NO_AVAILABILITY: requested %d but only %d left", modelReservation.Quantity, availableSpots)
		}

		// cream o noua rezervare
		modelReservation.ReservationReference = uuid.New().String()
		newReservation, err := txRepo.ReserveOrUpdate(ctx, modelReservation)
		if err != nil {
			tx.Rollback(ctx)
			return nil, fmt.Errorf("failed to create reservation: %w", err)
		}

		if err := tx.Commit(ctx); err != nil {
			if pgxErr, ok := err.(*pgconn.PgError); ok && pgxErr.Code == "40001" {
				time.Sleep(50 * time.Millisecond)
				continue
			}
			return nil, fmt.Errorf("commit error: %w", err)
		}

		return dto.ReservationToDTO(newReservation), nil
	}

	return nil, errors.New("too many retries due to concurrent modifications")
}
