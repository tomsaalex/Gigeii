package service

import (
	"context"

	"example.com/dto"
	"github.com/google/uuid"
)

type ResellerService interface {
	Register(ctx context.Context, req dto.RegisterResellerRequest) (*dto.ResellerResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.ResellerResponse, error)
	ListAll(ctx context.Context) ([]dto.ResellerResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.ResellerResponse, error)
	GetByUsername(ctx context.Context, username string) (*dto.ResellerResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}