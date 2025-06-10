package repository

import (
	"context"

	"example.com/db"
	"example.com/model"
	"github.com/google/uuid"
)

type ResellerRepository interface {
	Add(ctx context.Context, reseller model.Reseller) (*model.Reseller, error)
	GetByEmail(ctx context.Context, email string) (*model.Reseller, error)
	GetByUsername(ctx context.Context, username string) (*model.Reseller, error)
	GetById(ctx context.Context, id uuid.UUID) (*model.Reseller, error)
	GetAll(ctx context.Context) ([]model.Reseller, error)
	DeleteReseller(ctx context.Context, id uuid.UUID) error

}


func NewDbResellerRepository(queries *db.Queries) ResellerRepository {
	return &DBResellerRepository{
		queries: queries,
		mapper:  ResellerMapperDB{},
	}
}