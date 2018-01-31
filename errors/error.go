// Package errors error module.
package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	// ErrNotFound error not found.
	ErrNotFound = "notfound"
	// ErrUnauthorized error unauthorized request.
	ErrUnauthorized = "unauthorized"
	// ErrNotLoggedIn error not logged in.
	ErrNotLoggedIn = "notloggedin"
	// ErrMalformed error malformed request.
	ErrMalformed = "malformed"
	// ErrExpired token expired error.
	ErrExpired = "expired"
	// ErrFormError form error.
	ErrFormError = "formerror"
	// ErrInternalServerError Internal Server Error.
	ErrInternalServerError = "internal"
)

// SerializableError serializable error.
type SerializableError interface {
	json.Marshaler
	Error() string
}

// JSONError json error.
type JSONError map[string]SerializableError

// MarshalJSON marshal json.
func (je JSONError) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]SerializableError(je))
}

// Error error.
func (je JSONError) Error() string {
	return fmt.Sprintf("%+v", map[string]SerializableError(je))
}

// TextError text error.
type TextError string

// MarshalJSON marshal json.
func (te TextError) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(te))
}

// Error error.
func (te TextError) Error() string {
	return string(te)
}

// MessageError message error.
type MessageError struct {
	Msg  string
	Code int
	JSON JSONError
}

// New message error constructor.
func New(message string, code int) MessageError {
	return MessageError{
		Msg:  message,
		Code: code,
		JSON: JSONError(map[string]SerializableError{}),
	}
}

// NewSimple message error constructor with single string reason.
func NewSimple(message string, code int, err string) MessageError {
	return MessageError{
		Msg:  message,
		Code: code,
		JSON: JSONError(map[string]SerializableError{
			"reason": TextError(err),
		}),
	}
}

// Empty empty message error constructor.
func Empty() MessageError {
	return MessageError{
		Code: http.StatusBadRequest,
		JSON: JSONError(map[string]SerializableError{}),
	}
}

// MarshalJSON marshal json.
func (me MessageError) MarshalJSON() ([]byte, error) {
	return me.JSON.MarshalJSON()
}

// Error error.
func (me MessageError) Error() string {
	return fmt.Sprintf("%s %v", me.Msg, me.JSON)
}
