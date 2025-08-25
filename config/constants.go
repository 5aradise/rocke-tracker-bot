package config

import "errors"

type Code int

const ( // Codes
	CodeUnknown              Code = 999
	CodeUserWithTgIDExist    Code = 100
	CodeUserHasSub           Code = 101
	CodeUserWithTgIDNotExist Code = 102
	CodeTODO                 Code = -1
)

var ( // Errors
	ErrDuplicate        = errors.New("duplicate of an existing record")
	ErrUniqueConstraint = errors.New("unique constraint failed")
	ErrNotFound         = errors.New("record not found")
	ErrTODO             = errors.New("TODO")
)
