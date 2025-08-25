package config

import (
	"fmt"
	"log"
)

var NilError Error

type Error struct {
	Code Code
	Err  error
}

func NewError(code Code, err error) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func NewUnknownError(err error) Error {
	log.Printf("UNKNOWN ERROR: %s\n", err)
	return NewError(CodeUnknown, err)
}

func (e Error) Error() string {
	if e.Code == 0 {
		return ""
	}

	return fmt.Sprintf("%d: %s", e.Code, e.Err)
}

func (e Error) Unwrap() error {
	return e.Err
}

func (e Error) IsZero() bool {
	return e.Code == 0
}
