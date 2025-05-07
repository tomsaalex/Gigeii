package model

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID
	Fullname string
	Email    string
	Role     string
	PassHash []byte
	PassSalt []byte
}
