package service

import "fmt"

type ServiceError struct {
	Message string
}

func (se *ServiceError) Error() string {
	return fmt.Sprintf("ServiceError: %s", se.Message)
}

type AuthError struct {
	Message string
}

func (ae *AuthError) Error() string {
	return fmt.Sprintf("AuthError: %s", ae.Message)
}
