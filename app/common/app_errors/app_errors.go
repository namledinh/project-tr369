package apperrors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ErrDB                  = "ErrDB"
	ErrInvalidRequest      = "ErrInvalidRequest"
	ErrInternal            = "ErrInternal"
	ErrEntityNotExist      = "ErrEntityNotExist"
	ErrUnauthorized        = "ErrUnauthorized"
	ErrDeviceNotOnline     = "ErrDeviceNotOnline"
	ErrEntityAlreadyExists = "ErrEntityAlreadyExists"
)

type AppError struct {
	restCode   int
	resMessage string
	resKey     string

	// unexport fields for tracing and debugging purpose
	rootErr error
	log     string
}

// RootError implement the method to trace to the original error
func (e *AppError) RootError() error {
	if err, ok := e.rootErr.(*AppError); ok {
		return err.RootError()
	}
	return e.rootErr
}

func (e *AppError) Error() string { return e.RootError().Error() }

func (e *AppError) HTTPCode() int { return e.restCode }

func (e *AppError) ErrorKey() string { return e.resKey }

func (e *AppError) ErrorMessage() string { return e.resMessage }

func (e *AppError) IsError(key string) bool { return e.resKey == key }

func (e *AppError) ToJSONString() string {
	var err = map[string]interface{}{
		"errror":  e.Error(),
		"message": e.ErrorMessage(),
		"key":     e.ErrorKey(),
		"log":     e.log,
	}

	jsonErr, _ := json.Marshal(err)
	return string(jsonErr)
}

func (e *AppError) IsErrorKey(key string) bool {
	return e.resKey == key
}

func NewErrorResponse(
	statusCode int,
	root error,
	msg string,
	key string,
	log string,
) *AppError {
	return &AppError{
		restCode:   statusCode,
		resMessage: msg,
		resKey:     key,
		rootErr:    root,
		log:        log,
	}
}

func NewDBError(root error, db string) *AppError {
	return NewErrorResponse(
		http.StatusServiceUnavailable,
		root,
		"The server is currently unavailable",
		ErrDB,
		fmt.Sprintf("database error: %s", db))
}

func NewInvalidRequestError(root error, msg, log string) *AppError {
	return NewErrorResponse(
		http.StatusBadRequest,
		root,
		msg,
		ErrInvalidRequest,
		log,
	)
}

func NewInternalError(root error, log string) *AppError {
	return NewErrorResponse(
		http.StatusInternalServerError,
		root,
		"Internal server error",
		ErrInternal,
		log,
	)
}

func NewErrEntityNotExist(entity string) *AppError {
	errf := fmt.Errorf("%s not found", entity)
	return NewErrorResponse(
		http.StatusNotFound,
		errf,
		errf.Error(),
		ErrEntityNotExist,
		errf.Error(),
	)
}

func NewErrInvalidEventPayload(msg string) *AppError {
	return NewErrorResponse(
		http.StatusBadRequest,
		fmt.Errorf("invalid event payload"),
		fmt.Sprintf("invalid event payload: %s", msg),
		ErrInvalidRequest,
		fmt.Sprintf("invalid event payload: %s", msg),
	)
}

func NewErrUnauthorized() *AppError {
	return NewErrorResponse(
		http.StatusUnauthorized,
		fmt.Errorf("unauthorized"),
		"permission denied for action type",
		ErrUnauthorized,
		"permission denied for action type",
	)
}

func NewErrDeviceNotOnline(macAddress string) *AppError {
	return NewErrorResponse(
		http.StatusBadRequest,
		fmt.Errorf("device not online"),
		fmt.Sprintf("device %s is not online", macAddress),
		ErrDeviceNotOnline,
		fmt.Sprintf("device %s is not online", macAddress),
	)
}

func NewErrEntityAlreadyExists(entity string) *AppError {
	errf := fmt.Errorf("%s already exists", entity)
	return NewErrorResponse(
		http.StatusConflict,
		errf,
		errf.Error(),
		ErrEntityAlreadyExists,
		errf.Error(),
	)
}
