package domain

import "fmt"

type ErrorCode string

const (
	ErrCodeEntityNotFound ErrorCode = "EntityNotFound"
	ErrCodeInvalidInput   ErrorCode = "InvalidInput"
)

type DomainError struct {
	Err  error
	Code ErrorCode
}

func (e DomainError) Error() string {

	return fmt.Sprintf("%s - %s", e.Err.Error(), e.Code)
}

func NewDomainError(code ErrorCode, err error) DomainError {
	return DomainError{Code: code, Err: err}
}

func EntityNotFoundError(entityType string, entityID string) DomainError {
	return NewDomainError(ErrCodeEntityNotFound, fmt.Errorf("entity %s not found", entityType))
}

func InvalidInputError(input string, msg string) DomainError {
	return NewDomainError(ErrCodeInvalidInput, fmt.Errorf("error for input - %s: %s", input, msg))
}
