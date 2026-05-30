// Package response defines the standard envelope for all API responses.
//
// Every endpoint returns either a SuccessResponse or an ErrorResponse.
// These structs mirror the schemas defined in api/openapi.yaml exactly.
// oapi-codegen will also generate these from the YAML — but we define them
// here manually for Phase 1 before codegen is wired up.
//
// Why centralise responses?
// Without this, handlers invent their own shapes. One returns {"error": "..."},
// another returns {"message": "..."}. Clients can't rely on a consistent structure.
// With this, every response has the same outer envelope — always.
package response

import (
	"net/http"

	"github.com/Jeremiah-Williams1/inference-orchestrator/pkg/apperror"
	"github.com/gin-gonic/gin"
)

// SuccessResponse is the envelope for all successful responses.
type SuccessResponse struct {
	Message string      `json:"message"`
	Code    string      `json:"code"`
	Data    interface{} `json:"data"`
}

// ErrorResponse is the envelope for all error responses.
type ErrorResponse struct {
	Message string       `json:"message"`
	Code    string       `json:"code"`
	Errors  []FieldError `json:"errors,omitempty"`
}

// FieldError represents a single field-level validation failure.
// Used when a request has multiple invalid fields and you want to
// tell the caller exactly which ones and why.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// OK sends a 200 with a data payload.
func OK(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Message: message,
		Code:    "OK",
		Data:    data,
	})
}

// Accepted sends a 202 — the job was received and queued, not yet processed.
// 202 is the correct status for async operations: "I have it, I'll process it."
// 200 would imply the work is done, which it isn't.
func Accepted(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusAccepted, SuccessResponse{
		Message: message,
		Code:    "ACCEPTED",
		Data:    data,
	})
}

// Err maps an apperror to the correct HTTP status and sends the error response.
// This is the only place where AppError codes are translated to HTTP status codes.
// Handlers never call c.JSON for errors directly — they call this.
func Err(c *gin.Context, err error) {
	switch {
	case apperror.IsNotFound(err):
		c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error(), Code: "NOT_FOUND"})
	case apperror.IsBadRequest(err):
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error(), Code: "BAD_REQUEST"})
	case apperror.IsConflict(err):
		c.JSON(http.StatusConflict, ErrorResponse{Message: err.Error(), Code: "CONFLICT"})
	default:
		// Never leak internal error details to the caller.
		// Log the real error (the middleware does this), return a generic message.
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "internal server error", Code: "INTERNAL"})
	}
}

// ValidationErr sends a 400 for request binding failures.
// Kept separate from Err because gin's binding errors are not AppErrors —
// they come from the binding library and carry field-level detail we want to surface.
func ValidationErr(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Message: "validation failed",
		Code:    "VALIDATION_ERROR",
		Errors:  []FieldError{{Field: "request", Message: err.Error()}},
	})
}
