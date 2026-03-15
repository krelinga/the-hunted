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

type LayoutFormField[LocT, ItemT comparable] map[LocT]LayoutFormLoc[ItemT]

func (l LayoutFormField[LocT, ItemT]) Validate() error {
	for loc, locForm := range l {
		if err := locForm.validate(); err != nil {
			return fmt.Errorf("%w: invalid layout for location %v: %v", ErrInvalidFormField, loc, err)
		}
	}
	return nil
}

type LayoutFormLoc[ItemT comparable] struct {
	Capacity int
	Items map[ItemT]int
}

func (l LayoutFormLoc[ItemT]) validate() error {
	total := 0
	for item, count := range l.Items {
		if count < 0 {
			return fmt.Errorf("item %v has negative count %d", item, count)
		}
		total += count
	}
	if total > l.Capacity {
		return fmt.Errorf("total item count %d exceeds capacity %d", total, l.Capacity)
	}
	return nil
}

func NewLayoutFormLoc[ItemT comparable](capacity int) LayoutFormLoc[ItemT] {
	return LayoutFormLoc[ItemT]{
		Capacity: capacity,
		Items: map[ItemT]int{},
	}
}