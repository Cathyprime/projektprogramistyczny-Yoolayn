package msgs

import (
	"errors"
)

// error
var (
	ErrBadOptions        = errors.New("bad options provided for db operation")
	ErrDecode            = errors.New("failed decoding response")
	ErrFailedToGetParams = errors.New("failed to read params")
	ErrForbidden         = errors.New("action is forbidden")
	ErrInternal          = errors.New("an internal error has occurred")
	ErrNotFound          = errors.New("resource not found")
	ErrObjectIDConv      = errors.New("failed creating objectid from string")
	ErrTypeConn          = errors.New("Connection error: ")
	ErrUpdateFailed      = errors.New("failed to update the user")
	ErrUserCreation      = errors.New("failed creating user")
)

// debug
const (
	DebugStruct = "the value of struct: "
)
