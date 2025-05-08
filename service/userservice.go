package service

import (
	"context"

	"example.com/model"
	"example.com/repository"
)

type UserService struct {
	userRepo    repository.UserRepository
	argonHelper Argon2idHash
}

func NewUserService(userRepo repository.UserRepository, argonHelper *Argon2idHash) *UserService {
	return &UserService{
		userRepo:    userRepo,
		argonHelper: *argonHelper,
	}
}

func (s *UserService) Register(ctx context.Context, user model.User, password string) (*model.User, error) {
	_, err := s.userRepo.GetByEmail(ctx, user.Email)

	if err == nil {
		return nil, &repository.DuplicateEntityError{Message: "there's already a user using that email address"}
	}

	hashsalt, err := s.argonHelper.GenerateHash([]byte(password), nil)

	if err != nil {
		return nil, &AuthError{Message: "failed to generate hash for user's password"}
	}

	user.PassHash = hashsalt.Hash
	user.PassSalt = hashsalt.Salt

	addedUser, err := s.userRepo.Add(ctx, user)
	return addedUser, err
}

func (s *UserService) Login(ctx context.Context, email string, password string) (*model.User, error) {
	foundUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	err = s.argonHelper.Compare(foundUser.PassHash, foundUser.PassSalt, []byte(password))

	if err != nil {
		return nil, &AuthError{Message: "auth data is incorrect"}
	}

	return foundUser, nil
}
