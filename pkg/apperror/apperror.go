// Package apperror defines typed errors for the application.
//
// Why typed errors instead of plain fmt.Errorf?
// A plain error string gives the handler no information about what went wrong.
// Everything becomes a 500 or the handler has to parse strings — both are bad.
//
// With typed errors, the service says exactly what happened:
//   return apperror.NotFound("job not found")
//
// The handler checks the type and maps it to an HTTP status:
//   if apperror.IsNotFound(err) → 404
//
// No string parsing. No guessing. Clean separation between business logic and HTTP concerns.
package apperror

import "fmt"

type errorCode string

const (
	codeNotFound        errorCode = "NOT_FOUND"
	codeBadRequest      errorCode = "BAD_REQUEST"
	codeInternal        errorCode = "INTERNAL"
	codeConflict        errorCode = "CONFLICT"
	codeValidationError errorCode = "VALIDATION_ERROR"
)

// AppError is the error type all services in this application return.
type AppError struct {
	Code    errorCode
	Message string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Constructors — services use these to return errors.

func NotFound(msg string) *AppError {
	return &AppError{Code: codeNotFound, Message: msg}
}

func BadRequest(msg string) *AppError {
	return &AppError{Code: codeBadRequest, Message: msg}
}

func Internal(msg string) *AppError {
	return &AppError{Code: codeInternal, Message: msg}
}

func Conflict(msg string) *AppError {
	return &AppError{Code: codeConflict, Message: msg}
}

// Type checks — handlers use these to decide which HTTP status to return.

func IsNotFound(err error) bool {
	e, ok := err.(*AppError)
	return ok && e.Code == codeNotFound
}

func IsBadRequest(err error) bool {
	e, ok := err.(*AppError)
	return ok && e.Code == codeBadRequest
}

func IsInternal(err error) bool {
	e, ok := err.(*AppError)
	return ok && e.Code == codeInternal
}

func IsConflict(err error) bool {
	e, ok := err.(*AppError)
	return ok && e.Code == codeConflict
}
