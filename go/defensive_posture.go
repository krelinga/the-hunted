package thehunted

import (
	"errors"
	"fmt"
)

type DefensivePosture int

const (
	DefensivePostureSurfaced = iota + 1
	DefensivePostureSchnorcheling
)

func (dp DefensivePosture) String() string {
	switch dp {
	case DefensivePostureSurfaced:
		return "Surfaced"
	case DefensivePostureSchnorcheling:
		return "Schnorcheling"
	default:
		return fmt.Sprintf("Invalid Defensive Posture (%d)", dp)
	}
}

var ErrInvalidDefensivePosture = errors.New("invalid defensive posture")

func (dp DefensivePosture) Validate() error {
	switch dp {
	case DefensivePostureSurfaced, DefensivePostureSchnorcheling:
		return nil
	default:
		return fmt.Errorf("%w: %d", ErrInvalidDefensivePosture, dp)
	}
}

func (dp DefensivePosture) Must() {
	if err := dp.Validate(); err != nil {
		panic(err)
	}
}

type SelectDefensivePostureForm struct {
	baseForm
	DefensivePosture DefensivePosture
}

func (f *SelectDefensivePostureForm) Validate() error {
	return f.DefensivePosture.Validate()
}
