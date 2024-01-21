package msgs

import (
	"errors"
)

// error
var (
	ErrBadOptions   = errors.New("bad options provided for db operation")
	ErrDecode       = errors.New("failed decoding response")
	ErrInternal     = errors.New("an internal error has occurred")
	ErrTypeConn     = errors.New("Connection error: ")
	ErrUserCreation = errors.New("failed creating user")
	ErrObjectIDConv = errors.New("failed creating objectid from string")
	ErrNotFound     = errors.New("resource not found")
	ErrForbidden    = errors.New("action is forbidden")
	ErrUpdateFailed = errors.New("failed to update the user")
)

// debug
const (
	DebugStruct = "the value of struct: "
)
