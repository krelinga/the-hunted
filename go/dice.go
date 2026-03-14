package thehunted

import (
	"errors"
	"fmt"
	"math/rand"
)

type D6 int

var ErrInvalidDie = errors.New("invalid die")

func (d D6) Validate() error {
	if d < 1 || d > 6 {
		return fmt.Errorf("%w: die value %d is out of range (1-6)", ErrInvalidDie, d)
	}
	return nil
}

type Roller interface {
	RollD6() D6
}

type RandomRoller struct {}

func (r RandomRoller) RollD6() D6 {
	return D6(rand.Intn(6) + 1)
}