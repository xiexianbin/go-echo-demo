package ecode

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	once      sync.Once
	_messages atomic.Value         // stored map[string]map[int]string
	_codes    = map[int]struct{}{} // store registered code
)

// Register register code message map
func Register(cm ...map[Code]string) {
	for _, m := range cm {
		_messages.Store(m)
	}
}

// New new a Code by int value
func New(e int) Code {
	if e <= 0 {
		panic("code must great than 0")
	}

	if _, ok := _codes[e]; ok {
		panic(fmt.Sprintf("code: %d already exist!", e))
	}
	return Code(e)
}

// ECode ecode interface
type ECode interface {
	// Code return error code
	Code() int

	// Name return error code string, auto generated by "stringer -type=Code"
	String() string

	// Error ecode string for error.Error interface
	Error() string

	// HTTPCode return http status_code for Code
	HTTPCode() int

	// Message return code message
	Message() string
}

// stringer ref https://pkg.go.dev/golang.org/x/tools/cmd/stringer
//
//go:generate stringer -type=Code
type Code int

// Code return error code
func (e Code) Code() int {
	return int(e)
}

// Error ecode string for error.Error interface
// equal to e.String()
func (e Code) Error() string {
	return e.String()
}

// HTTPCode return http status_code for Code
func (c Code) HTTPCode() int {
	if c == CodeOK || strings.Contains(c.String(), "OK") { // 200
		return http.StatusOK
	} else if strings.Contains(c.String(), "Created") { // 201
		return http.StatusCreated
	} else if strings.Contains(c.String(), "Accepted") { // 202
		return http.StatusAccepted
	} else if strings.Contains(c.String(), "BadRequest") { // 400
		return http.StatusBadRequest
	} else if strings.Contains(c.String(), "Unauthorized") { // 401
		return http.StatusUnauthorized
	} else if strings.Contains(c.String(), "Forbidden") { // 403
		return http.StatusForbidden
	} else if strings.Contains(c.String(), "NotFound") { // 404
		return http.StatusNotFound
	} else if strings.Contains(c.String(), "Conflict") { // 409
		return http.StatusConflict
	}

	return http.StatusBadRequest
}

// Message return code message
func (e Code) Message() string {
	if m, ok := _messages.Load().(map[int]string); ok {
		if msg, ok := m[e.Code()]; ok {
			return msg
		}
	}
	return e.Error()
}

func init() {
	once.Do(func() {
		Register(
			codeCommonMessages,
		)
	})
}