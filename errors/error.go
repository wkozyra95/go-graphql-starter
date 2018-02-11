// Package errors error module.
package errors

import (
	"fmt"
)

var (
	// ErrNotFound error not found.
	ErrNotFound = fmt.Errorf("notfound")
	// ErrUnauthorized error unauthorized request.
	ErrUnauthorized = fmt.Errorf("unauthorized")
	// ErrNotLoggedIn error not logged in.
	ErrNotLoggedIn = fmt.Errorf("notloggedin")
	// ErrMalformed error malformed request.
	ErrMalformed = fmt.Errorf("malformed")
	// ErrExpired token expired error.
	ErrExpired = fmt.Errorf("expired")
	// ErrFormError form error.
	ErrFormError = fmt.Errorf("formerror")
	// ErrInternalServerError Internal Server Error.
	ErrInternalServerError = fmt.Errorf("internal")
)
