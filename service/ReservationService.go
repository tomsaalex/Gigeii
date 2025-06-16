package service

import (
	"context"

	"example.com/dto"
	"example.com/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationService interface {
	ReserveOrUpdate(ctx context.Context, reservation dto.ReservationDTO) (*dto.ReservationDTO, error)
	GetByReservationReference(ctx context.Context, reservationReference string) (*dto.ReservationDTO, error)
	GetByExternalReservationReference(ctx context.Context, resellerID string, externalReservationReference string) (*dto.ReservationDTO, error)
	CancelReservation(ctx context.Context, reservationReference string, externalReservationReference string ) error
}

func NewReservationService(repo repository.ReservationRepository, repoAvailability repository.AvailabilityRepository, connPool *pgxpool.Pool) ReservationService {
	return &ReservationServiceImpl{
		repo: repo,
		availability: repoAvailability,
		dbPool: connPool,
	}
}