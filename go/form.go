package thehunted

import (
	"errors"
	"fmt"
)

type Form interface {
	formIsAClosedType()
}

type baseForm struct {}

func (_ *baseForm) formIsAClosedType() {}

type TextFormField string

type SelectFormField[T any] struct {
	Options []T
	Selected int
}

type BoolFormField bool

var ErrInvalidFormField = errors.New("invalid form field")

func (s SelectFormField[T]) Validate() error {
	if s.Selected < 0 || s.Selected >= len(s.Options) {
		return fmt.Errorf("%w: selected index %d is out of range for options %v", ErrInvalidFormField, s.Selected, s.Options)
	}
	return nil
}