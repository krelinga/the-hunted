package thehunted

import (
	"errors"
	"fmt"
	"math/rand"
)

type DiceD6 int

func (d DiceD6) String() string {
	if err := d.Validate(); err != nil {
		return fmt.Sprintf("Invalid D6 (%d)", d)
	}
	return fmt.Sprintf("%d", d)
}

var ErrInvalidDice = errors.New("invalid Dice value")

func (d DiceD6) Validate() error {
	if d < 1 || d > 6 {
		return fmt.Errorf("%w: die value %d is out of range (1-6)", ErrInvalidDice, d)
	}
	return nil
}

func (d DiceD6) Must() {
	if err := d.Validate(); err != nil {
		panic(err)
	}
}

func (d DiceD6) AsInt() int {
	d.Must()
	return int(d)
}

type Result2D6 struct {
	Dice1 DiceD6
	Dice2 DiceD6
}

func (r Result2D6) String() string {
	return fmt.Sprintf("%s + %s = %d", r.Dice1, r.Dice2, r.AsInt())
}

func (r Result2D6) Validate() error {
	if err := r.Dice1.Validate(); err != nil {
		return fmt.Errorf("invalid first die: %w", err)
	}
	if err := r.Dice2.Validate(); err != nil {
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
	return r.Dice1.AsInt() + r.Dice2.AsInt()
}

type Roller interface {
	RollD6() DiceD6
	Roll2D6() Result2D6
}

type RandomRoller struct{}

func (r RandomRoller) RollD6() DiceD6 {
	return DiceD6(rand.Intn(6) + 1)
}

func (r RandomRoller) Roll2D6() Result2D6 {
	return Result2D6{
		Dice1: r.RollD6(),
		Dice2: r.RollD6(),
	}
}
