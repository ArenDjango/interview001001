package common

import "fmt"

type ProcessingError struct {
	Status uint8
	Err    error
}

func (e *ProcessingError) Error() string {
	return fmt.Sprintf("status: %d, error: %v", e.Status, e.Err)
}

func (e *ProcessingError) Unwrap() error {
	return e.Err
}

func NewError(status uint8, err error) *ProcessingError {
	return &ProcessingError{
		Status: status,
		Err:    err,
	}
}
