package exitcode

import (
	"errors"
	"os"
)

type Code int

const (
	OK          Code = 0
	Fail        Code = 1
	InvalidArgs Code = 2
	IOError     Code = 3
	NotFound    Code = 4
)

func (c Code) Int() int {
	return int(c)
}

// FromError maps an error into an exit code
func FromError(err error) Code {
	if err == nil {
		return OK
	}

	if errors.Is(err, os.ErrNotExist) {
		return NotFound
	}
	if errors.Is(err, os.ErrPermission) {
		return IOError
	}

	return Fail
}
