package repository

import (
	"context"
	"fmt"

	"example.com/db"
	"example.com/model"
)

type DBUserRepository struct {
	queries *db.Queries
	mapper  UserMapperDB
}

func (r *DBUserRepository) Add(ctx context.Context, user model.User) (*model.User, error) {
	addUserParams := r.mapper.UserToAddUserParams(user)
	dbUser, err := r.queries.AddUser(ctx, addUserParams)
	if err != nil {
		return nil, &DuplicateEntityError{
			Message: fmt.Sprintf("Add: Email \"%s\" is already taken by a different user.", user.Email),
		}
	}

	modelUser := r.mapper.DBUserToUser(dbUser)
	return modelUser, nil
}
func (r *DBUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	foundUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, &EntityNotFoundError{Message: fmt.Sprintf("GetUserByEmail: No user matches email \"%s\"", email)}
	}

	convertedUser := r.mapper.DBUserToUser(foundUser)
	return convertedUser, nil
}
