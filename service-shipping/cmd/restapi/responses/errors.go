package responses

import (
	"encoding/json"
	"fmt"
	"github.com/jna-distribution/service-shipping/internal"
	"log"
	"net/http"
)

type (
	// ErrorResponse struct.
	ErrorResponse struct {
		Status  int         `json:"-"`
		Message string      `json:"message"`
		Details interface{} `json:"details,omitempty"`
	}

	// ErrBadRequest response 400
	ErrBadRequest internal.Err

	// ErrUnauthorized response 401
	ErrUnauthorized internal.Err

	// ErrForbidden response 403
	ErrForbidden internal.Err

	// ErrMethodNotAllowed response 405
	ErrMethodNotAllowed internal.Err

	//ErrInternalServer response 500
	ErrInternalServer internal.Err
)

// Error method for implement error interface(golang builtin interface).
func (e ErrorResponse) Error() string {
	return fmt.Sprintf("[%v]: %v", e.Status, e.Message)
}

// Error implement error interface for ErrBadRequest
func (e ErrBadRequest) Error() string {
	return fmt.Sprintf("Error: Message: %v, InternalError: %v\n", e.Message, e.InternalError)
}

// Error implement error interface for ErrUnauthorized
func (e ErrUnauthorized) Error() string {
	return fmt.Sprintf("Error: Message: %v, InternalError: %v\n", e.Message, e.InternalError)
}

// Error implement error interface for ErrForbidden
func (e ErrForbidden) Error() string {
	return fmt.Sprintf("Error: Message: %v, InternalError: %v\n", e.Message, e.InternalError)
}

// Error implement error interface for ErrInternalServer
func (e ErrInternalServer) Error() string {
	return fmt.Sprintf("Error: Message: %v, InternalError: %v\n", e.Message, e.InternalError)
}

// Error implement error interface for ErrMethodNotAllowed
func (e ErrMethodNotAllowed) Error() string {
	return fmt.Sprintf("Error: Message: %v, InternalError: %v\n", e.Message, e.InternalError)
}

// BuildErrorResponse function to build ErrorResponse error type with proper status code.
func BuildErrorResponse(err error) ErrorResponse {
	var errRes ErrorResponse
	var internalErr error

	switch err.(type) {
	case internal.ErrNotFound:
		e := err.(internal.ErrNotFound)
		//Default message.
		if e.Message == "" {
			e.Message = "The requested resource was not found."
		}
		//Save InternalError
		internalErr = e.InternalError
		errRes = ErrorResponse{
			Status:  http.StatusNotFound,
			Message: e.Message,
			Details: e.Details,
		}

	case internal.ErrDatabase:
		e := err.(internal.ErrDatabase)
		//Default message.
		if e.Message == "" {
			e.Message = "Internal server error."
		}
		//Save InternalError
		internalErr = e.InternalError
		errRes = ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: e.Message,
		}

	case internal.ErrExists:
		e := err.(internal.ErrExists)
		if e.Message == "" {
			e.Message = "Resource already exists."
		}
		errRes = ErrorResponse{
			Status:  http.StatusConflict,
			Message: e.Message,
		}

	case internal.ErrInvalidInput:
		//Invalid input show internal error that return from ozzo-validation to user.
		e := err.(internal.ErrInvalidInput)
		//Set default message.
		if e.Message == "" {
			e.Message = "There is some problem with the data you submitted."
		}
		//Save InternalError
		internalErr = e.InternalError
		errRes = ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: e.Message,
			Details: e.Details,
		}

	case ErrBadRequest:
		e := err.(ErrBadRequest)
		switch e.InternalError.(type) {
		case *json.SyntaxError:
			errRes = ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "The requested input is invalid JSON syntax.",
			}
		default:
			errRes = ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Bad request.",
			}
		}

	case ErrUnauthorized:
		e := err.(ErrUnauthorized)
		//Set default message.
		if e.Message == "" {
			e.Message = "This action not allowed to perform."
		}
		//Save InternalError
		internalErr = e.InternalError
		errRes = ErrorResponse{
			Status:  http.StatusUnauthorized,
			Message: e.Message,
			Details: e.Details,
		}

	case ErrForbidden:
		e := err.(ErrForbidden)
		//Set default message.
		if e.Message == "" {
			e.Message = "Access to this resource is forbidden."
		}
		//Save InternalError
		internalErr = e.InternalError
		errRes = ErrorResponse{
			Status:  http.StatusForbidden,
			Message: e.Message,
			Details: e.Details,
		}

	case ErrMethodNotAllowed:
		e := err.(ErrMethodNotAllowed)
		//Default message.
		if e.Message == "" {
			e.Message = "The requested method not allowed."
		}
		//Save InternalError.
		internalErr = e.InternalError
		errRes = ErrorResponse{
			Status:  http.StatusMethodNotAllowed,
			Message: e.Message,
		}

	case ErrInternalServer:
		e := err.(ErrInternalServer)
		//Default message.
		if e.Message == "" {
			e.Message = "Internal server error."
		}
		//Save InternalError
		internalErr = e.InternalError
		errRes = ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: e.Message,
		}

	default:
		//Default, not match any type of error, let's check log.
		errRes = ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error.",
		}
	}

	log.Printf("ErrorResponse: %v\n", errRes)
	log.Printf("Error:			type: %T, value: %+v\n", err, err)
	log.Printf("InternalError:	type: %T, value: %+v\n", internalErr, internalErr)
	return errRes
}
