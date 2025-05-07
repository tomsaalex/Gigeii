package repository

import (
	"context"

	"example.com/db"
	"example.com/model"
)

type DBUserRepository struct {
	queries *db.Queries
	mapper  UserMapperDB
}

func (r *DBUserRepository) Add(ctx context.Context, user model.User) (model.User, error) {
	return model.User{}, nil
}
func (r *DBUserRepository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	return model.User{}, nil
}
