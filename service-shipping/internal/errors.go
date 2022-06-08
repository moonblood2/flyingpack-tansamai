package internal

import "fmt"

// Err contains message and internalError
type Err struct {
	Message       string
	Details       interface{}
	InternalError error
}

type (
	//ErrNotFound use in repository eg. user that specific id, and other request resource that not found eg. path.
	ErrNotFound Err

	//ErrDatabase use in repository.
	ErrDatabase Err

	//ErrExists
	ErrExists Err

	//ErrInvalidInput use in service when validate input.
	ErrInvalidInput Err

	//ErrInternal use when not match any case.
	ErrInternal Err
)

// Error implement error interface for ErrNotFound
func (e ErrNotFound) Error() string {
	return fmt.Sprintf("Error: Message: %v, InternalError: %v\n", e.Message, e.InternalError)
}

// Error implement error interface for ErrDatabase
func (e ErrDatabase) Error() string {
	return fmt.Sprintf("Error: Message: %v, InternalError: %v\n", e.Message, e.InternalError)
}

// Error implement error interface for ErrExists
func (e ErrExists) Error() string {
	return fmt.Sprintf("Error: Message: %v, InternalError: %v\n", e.Message, e.InternalError)
}

// Error implement error interface for ErrInvalidInput
func (e ErrInvalidInput) Error() string {
	return fmt.Sprintf("Error: Message: %v, InternalError: %v\n", e.Message, e.InternalError)
}

// Error implement error interface for ErrInternal
func (e ErrInternal) Error() string {
	return fmt.Sprintf("Error: Message: %v, InternalError: %v\n", e.Message, e.InternalError)
}
