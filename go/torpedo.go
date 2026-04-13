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

type TorpLoc struct {
	facing Facing
	tube   int
}

func NewTorpLocTube(f Facing, tube int) TorpLoc {
	if tube < 1 {
		panic(fmt.Sprintf("invalid tube number: %d", tube))
	}
	return TorpLoc{
		facing: f,
		tube:   tube,
	}
}

func NewTorpLocReload(f Facing) TorpLoc {
	return TorpLoc{
		facing: f,
		tube:   0,
	}
}

func (l TorpLoc) String() string {
	if l.tube == 0 {
		return fmt.Sprintf("%s reload", l.facing)
	}
	return fmt.Sprintf("%s tube %d", l.facing, l.tube)
}

func (l TorpLoc) IsTube() bool {
	return l.tube != 0
}

func (l TorpLoc) GetTube() (int, bool) {
	return l.tube, l.IsTube()
}

func (l TorpLoc) IsReload() bool {
	return l.tube == 0
}

func (l TorpLoc) GetFacing() Facing {
	return l.facing
}

// TODO: decide if I like this ordering.
func TorpLocCmp(a, b TorpLoc) int {
	if a.facing != b.facing {
		if a.facing == FacingFwd {
			return -1
		}
		return 1
	}
	return a.tube - b.tube
}
