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

type ValidationError struct {
	ErrorsList []string
}

func (ve *ValidationError) Error() string {
	output := "ValidationError:"
	for _, errMsg := range ve.ErrorsList {
		output += errMsg + " "
	}
	return output
}

type UnhandledConflictError struct {
	Message string
}

func (e *UnhandledConflictError) Error() string {
	return fmt.Sprintf("UnhandledConflictError: %s", e.Message)
}
