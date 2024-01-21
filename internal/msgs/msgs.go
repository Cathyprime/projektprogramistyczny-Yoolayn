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
)

// debug
const (
	DebugStruct = "the value of struct: "
)
