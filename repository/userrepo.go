package repository

import (
	"context"

	"example.com/db"
	"example.com/model"
)

type UserRepository interface {
	Add(ctx context.Context, user model.User) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

func NewDBUserRepository(queries *db.Queries) UserRepository {
	return &DBUserRepository{
		queries: queries,
		mapper:  UserMapperDB{},
	}
}
