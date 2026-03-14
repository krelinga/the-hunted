package thehunted

import (
	"errors"
	"fmt"
)

type TorpType int

const (
	TorpTypeG7a TorpType = iota + 1
	TorpTypeG7e
	TorpTypeG7esZaunkonig
	TorpTypeG7esZaunkonigII
	TorpTypeG7eFalke
)

var ErrInvalidTorpType = errors.New("invalid torp type")

func (t TorpType) Validate() error {
	switch t {
	case TorpTypeG7a, TorpTypeG7e, TorpTypeG7esZaunkonig, TorpTypeG7esZaunkonigII, TorpTypeG7eFalke:
		return nil
	default:
		return fmt.Errorf("%w: %d", ErrInvalidTorpType, t)
	}
}

func (t TorpType) String() string {
	switch t {
	case TorpTypeG7a:
		return "G7a Steam"
	case TorpTypeG7e:
		return "G7e Electric"
	case TorpTypeG7esZaunkonig:
		return "G7es Zaunkönig"
	case TorpTypeG7esZaunkonigII:
		return "G7es Zaunkönig II"
	case TorpTypeG7eFalke:
		return "G7e Falke"
	default:
		return fmt.Sprintf("Unknown torpedo (%d)", t)
	}
}