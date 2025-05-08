package repository

import "fmt"

type RepositoryError struct {
	Message string
}

func (e *RepositoryError) Error() string {
	return fmt.Sprintf("RepositoryError: %s", e.Message)
}

type EntityNotFoundError struct {
	Message string
}

func (e *EntityNotFoundError) Error() string {
	return fmt.Sprintf("EntityNotFoundError: %s", e.Message)
}

type EntityDBMappingError struct {
	Message string
}

func (e *EntityDBMappingError) Error() string {
	return fmt.Sprintf("EntityDBMappingError: %s", e.Message)
}

type DuplicateEntityError struct {
	Message string
}

func (e *DuplicateEntityError) Error() string {
	return fmt.Sprintf("DuplicateEntityError: %s", e.Message)
}
