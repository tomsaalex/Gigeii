package model

import "github.com/google/uuid"

type Reseller struct {
	Id          uuid.UUID
	Name        string
	Username    string
	PasswordHash []byte
	Email       string
}