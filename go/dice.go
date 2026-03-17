package thehunted

import (
	"errors"
	"fmt"
	"math/rand"
)

type ResultD6 int

func (d ResultD6) String() string {
	if err := d.Validate(); err != nil {
		return fmt.Sprintf("Invalid D6 (%d)", d)
	}
	return fmt.Sprintf("%d", d)
}

var ErrInvalidResult = errors.New("invalid Result")

func (d ResultD6) Validate() error {
	if d < 1 || d > 6 {
		return fmt.Errorf("%w: die value %d is out of range (1-6)", ErrInvalidResult, d)
	}
	return nil
}

func (d ResultD6) Must() {
	if err := d.Validate(); err != nil {
		panic(err)
	}
}

func (d ResultD6) AsInt() int {
	d.Must()
	return int(d)
}

type Result2D6 struct {
	Die1 ResultD6
	Die2 ResultD6
}

func (r Result2D6) String() string {
	return fmt.Sprintf("%s + %s", r.Die1, r.Die2)
}

func (r Result2D6) Validate() error {
	if err := r.Die1.Validate(); err != nil {
		return fmt.Errorf("invalid first die: %w", err)
	}
	if err := r.Die2.Validate(); err != nil {
		return fmt.Errorf("invalid second die: %w", err)
	}
	return nil
}

func (r Result2D6) Must() {
	if err := r.Validate(); err != nil {
		panic(err)
	}
}

func (r Result2D6) AsInt() int {
	r.Must()
	return r.Die1.AsInt() + r.Die2.AsInt()
}

type Roller interface {
	RollD6() ResultD6
	Roll2D6() Result2D6
}

type RandomRoller struct{}

func (r RandomRoller) RollD6() ResultD6 {
	return ResultD6(rand.Intn(6) + 1)
}

func (r RandomRoller) Roll2D6() Result2D6 {
	return Result2D6{
		Die1: r.RollD6(),
		Die2: r.RollD6(),
	}
}