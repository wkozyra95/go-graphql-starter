// Package errors error module.
package errors

import (
	"fmt"
)

var (
	// NotFound error not found.
	NotFound = fmt.Errorf("notfound")
	// ErrUnauthorized error unauthorized request.
	Unauthorized = fmt.Errorf("unauthorized")
	// ErrNotLoggedIn error not logged in.
	NotLoggedIn = fmt.Errorf("notloggedin")
	// ErrMalformed error malformed request.
	Malformed = fmt.Errorf("malformed")
	// ErrExpired token expired error.
	Expired = fmt.Errorf("expired")
	// ErrFormError form error.
	FormError = fmt.Errorf("formerror")
	// ErrInternalServerError Internal Server Error.
	InternalServerError = fmt.Errorf("internal")
)
