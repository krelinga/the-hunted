package thehunted

import (
	"errors"
	"fmt"
)

type Torpedo int

const (
	TorpedoG7a Torpedo = iota + 1
	TorpedoG7e
	TorpedoG7es
	TorpedoG7eFalke
)

var ErrInvalidTorpedo = errors.New("invalid torpedo")

func (t Torpedo) Validate() error {
	switch t {
	case TorpedoG7a, TorpedoG7e, TorpedoG7es, TorpedoG7eFalke:
		return nil
	default:
		return fmt.Errorf("%w: %d", ErrInvalidTorpedo, t)
	}
}

func (t Torpedo) String() string {
	switch t {
	case TorpedoG7a:
		return "G7a Steam"
	case TorpedoG7e:
		return "G7e Electric"
	case TorpedoG7es:
		return "G7es Zaunkönig"
	case TorpedoG7eFalke:
		return "G7e Falke"
	default:
		return fmt.Sprintf("Unknown torpedo (%d)", t)
	}
}