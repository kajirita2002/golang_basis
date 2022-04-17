package code

import (
	"errors"
	"fmt"
)

type privateError struct {
	code Code
	msg  string
}

func (e privateError) Error() string {
	return fmt.Sprintf("Code: %s, Msg: %s", e.code, e.msg)
}

// GetCode takes error and returns Code accordingly.
func GetCode(err error) Code {
	if err == nil {
		return OK
	}
	var e privateError
	if errors.As(err, &e) {
		return e.code
	}
	return Unknown
}
