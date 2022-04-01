package domain

import (
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Custom ErrorType Constants
const (
	ErrInvalidArgument = iota // ErrInvalidArgument is returned when an argument is invalid.
	ErrNotFound               // ErrNotFound is returned when a resource is not found.
	ErrInternalError          // ErrInternalError is returned when an internal error occurs.
	ErrBadRequest             // ErrBadRequest is returned when a bad request is made.
)

// Custom Error Type
type Error struct {
	// Error Type, must be one of the ErrTypes constants
	errType int
	// Orginal Error
	err error

	// User Friendly Message for UI
	uiMsg string
}

// Error interface implementation
// Error will call the original error of custom error
func (e *Error) Error() string {
	return e.err.Error()
}

// uiMsgFromType returns the User Friendly message from error type
func uiMsgFromType(errType int) string {
	switch errType {
	case ErrInvalidArgument:
		return "Invalid argument"
	case ErrNotFound:
		return "Not found"
	case ErrInternalError:
		return "Internal server error, Try again later"
	case ErrBadRequest:
		return "Bad request"
	}
	return "Internal server error, Try again later"
}

// Create New Custom Error with Error Type and Error Message
func NewError(err error, errType int, hint ...string) error {
	// To avoid current NewError Function from stack
	depthOne := 1
	errors.WithStackDepth(err, depthOne)

	return &Error{
		errType: errType,
		err:     fmt.Errorf("%s: %s", hint, err),
		uiMsg:   uiMsgFromType(errType),
	}
}

// Checks if the error is a domain error with type
func ErrIs(err error, errorType int) bool {
	if err == nil {
		return false
	}
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.errType == errorType
}

// Get original error to Wrap
func ErrToWrap(err error) error {
	if err == nil {
		return nil
	}

	e, ok := err.(*Error)
	if !ok {
		log.Warn("Error is not a domain error")
		return err
	}
	return e.err
}

// Log Errors
func ErrLog(err error) {
	if err == nil {
		return
	}
	e, ok := err.(*Error)
	if !ok {
		// Not a domain error
		log.Errorf("%+v", err)
		return
	}
	log.Error(e.err)
}

// ErrFailedGinReq handles error from gin.Context
func ErrFailedGinReq(c *gin.Context, err error) {
	// If the err is nil, log warn
	if err == nil {
		log.Warn("FailedGinReq: Error is nil")
		return
	}

	c.JSON(400, gin.H{
		"status": "error",
		"error": gin.H{
			"message": err.Error(),
		},
	})
}
