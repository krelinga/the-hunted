package thehunted

import "errors"

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
		return ErrInvalidCrewQuality
	}
}
