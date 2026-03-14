package thehunted

import (
	"errors"
	"fmt"
)

type CrewQuality int

const (
	CrewQualityGreen CrewQuality = iota + 1
	CrewQualityTrained
	CrewQualityVeteran
	CrewQualityElite
)

var ErrInvalidCrewQuality = errors.New("invalid crew quality")

func (cq CrewQuality) Validate() error {
	switch cq {
	case CrewQualityGreen, CrewQualityTrained, CrewQualityVeteran, CrewQualityElite:
		return nil
	default:
		return fmt.Errorf("%w: %d", ErrInvalidCrewQuality, cq)
	}
}

func (cq CrewQuality) String() string {
	switch cq {
	case CrewQualityGreen:
		return "Green"
	case CrewQualityTrained:
		return "Trained"
	case CrewQualityVeteran:
		return "Veteran"
	case CrewQualityElite:
		return "Elite"
	default:
		return fmt.Sprintf("Unknown crew quality (%d)", cq)
	}
}