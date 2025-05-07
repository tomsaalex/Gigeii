package service

import "example.com/repository"

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
