package service

import (
	"context"
	"errors"

	"example.com/dto"
	"example.com/model"
	"example.com/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type resellerServiceImpl struct {
	repo repository.ResellerRepository
}

// Delete implements ResellerService.
func (r *resellerServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.repo.DeleteReseller(ctx, id)
}

// GetByID implements ResellerService.
func (r *resellerServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*dto.ResellerResponse, error) {
	reseller, err := r.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := dto.ToResellerResponse(*reseller)
	return &resp, nil
}

// GetByUsername implements ResellerService.
func (r *resellerServiceImpl) GetByUsername(ctx context.Context, username string) (*dto.ResellerResponse, error) {
	reseller, err := r.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	resp := dto.ToResellerResponse(*reseller)
	return &resp, nil
}

// ListAll implements ResellerService.
func (r *resellerServiceImpl) ListAll(ctx context.Context) ([]dto.ResellerResponse, error) {
	resellers, err := r.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.ResellerResponse, len(resellers))
	for i, item := range resellers {
		result[i] = dto.ToResellerResponse(item)
	}
	return result, nil
}

// Login implements ResellerService.
func (r *resellerServiceImpl) Login(ctx context.Context, req dto.LoginRequest) (*dto.ResellerResponse, error) {
	reseller, err := r.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword(reseller.PasswordHash, []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	resp := dto.ToResellerResponse(*reseller)
	return &resp, nil
}

// Register implements ResellerService.
func (r *resellerServiceImpl) Register(ctx context.Context, req dto.RegisterResellerRequest) (*dto.ResellerResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	reseller := model.Reseller{
		Id:           uuid.New(),
		Name:         req.Name,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashed,
	}

	saved, err := r.repo.Add(ctx, reseller)
	if err != nil {
		return nil, err
	}

	resp := dto.ToResellerResponse(*saved)
	return &resp, nil
}



func NewResellerService(repo repository.ResellerRepository) ResellerService {
	return &resellerServiceImpl{repo: repo}
}
