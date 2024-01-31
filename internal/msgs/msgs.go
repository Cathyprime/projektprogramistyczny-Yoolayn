package msgs

import (
	"errors"
	"net/http"

	"github.com/charmbracelet/log"
)

type respError struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Content string `json:"reason"`
}

// error
var (
	ErrBadOptions        = errors.New("bad options provided for db operation")
	ErrDecode            = errors.New("failed decoding response")
	ErrEncryption        = errors.New("failed to encrypt the password")
	ErrFailedToGetParams = errors.New("failed to read params")
	ErrForbidden         = errors.New("action is forbidden")
	ErrInternal          = errors.New("an internal error has occurred")
	ErrNotFound          = errors.New("resource not found")
	ErrObjectIDConv      = errors.New("failed creating objectid from string")
	ErrTaken             = errors.New("name is taken")
	ErrTypeConn          = errors.New("Connection error")
	ErrUpdateFailed      = errors.New("failed to update the user")
	ErrUserCreation      = errors.New("failed creating user")
	ErrWrongEmailFormat  = errors.New("email not formated properly")
	ErrWrongFormat       = errors.New("wrong body format")
	ErrNotAuthorized     = errors.New("credentials not authorized")
	ErrDeleteFailed      = errors.New("failed to delete the user")
)

// debug
var (
	DebugStruct      = errors.New("the value of struct:")
	DebugSkippedLoop = errors.New("Loop skipped")
	DebugJSON        = errors.New("the value of json:")
)

var msgmap = map[error]int{
	ErrBadOptions:        http.StatusBadRequest,
	ErrDecode:            http.StatusInternalServerError,
	ErrFailedToGetParams: http.StatusInternalServerError,
	ErrInternal:          http.StatusInternalServerError,
	ErrObjectIDConv:      http.StatusBadRequest,
	ErrTaken:             http.StatusBadRequest,
	ErrUpdateFailed:      http.StatusBadRequest,
	ErrWrongEmailFormat:  http.StatusBadRequest,
	ErrWrongFormat:       http.StatusBadRequest,
	ErrEncryption:        http.StatusBadRequest,
	ErrForbidden:         http.StatusForbidden,
	ErrNotFound:          http.StatusNotFound,
	ErrNotAuthorized:     http.StatusUnauthorized,
	ErrDeleteFailed:      http.StatusBadRequest,
}

func ReportError(err error, content string, info ...any) (int, respError) {
	log.Error(err, info...)
	return msgmap[err], respError{
		Code:    msgmap[err],
		Error:   err.Error(),
		Content: content,
	}
}
