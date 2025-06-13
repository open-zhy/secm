package errors

import (
	e "errors"
	"fmt"
)

func Wrapf(err error, str string, vars ...any) error {
	return e.Join(fmt.Errorf(str, vars...), err)
}

func New(str string, vars ...any) error {
	return fmt.Errorf(str, vars...)
}
