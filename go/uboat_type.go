package thehunted

import (
	"errors"
	"fmt"
)

type UBoatType int

const (
	UBoatTypeVIIB = iota + 1
	UBoatTypeVIIC
	UBoatTypeVIICFlak
	UBoatTypeVIIC41
	UBoatTypeVIID
	UBoatTypeIXB
	UBoatTypeIXC
	UBoatTypeIXC40
	UBoatTypeIXD2
	UBoatTypeIXD42
	UBoatTypeXB
	UBoatTypeXII
	UBoatTypeXIV
	UBoatTypeXXI
)

var ErrInvalidUBoatType = errors.New("invalid u-boat type")

func (u UBoatType) Validate() error {
	switch u {
	case UBoatTypeVIIB, UBoatTypeVIIC, UBoatTypeVIICFlak, UBoatTypeVIIC41, UBoatTypeVIID,
		UBoatTypeIXB, UBoatTypeIXC, UBoatTypeIXC40, UBoatTypeIXD2, UBoatTypeIXD42,
		UBoatTypeXB, UBoatTypeXII, UBoatTypeXIV, UBoatTypeXXI:
		return nil
	default:
		return ErrInvalidUBoatType
	}
}

func (u UBoatType) String() string {
	switch u {
	case UBoatTypeVIIB:
		return "VIIB"
	case UBoatTypeVIIC:
		return "VIIC"
	case UBoatTypeVIICFlak:
		return "VIIC Flak"
	case UBoatTypeVIIC41:
		return "VIIC/41"
	case UBoatTypeVIID:
		return "VIID"
	case UBoatTypeIXB:
		return "IXB"
	case UBoatTypeIXC:
		return "IXC"
	case UBoatTypeIXC40:
		return "IXC/40"
	case UBoatTypeIXD2:
		return "IXD-2"
	case UBoatTypeIXD42:
		return "IXD/42"
	case UBoatTypeXB:
		return "XB"
	case UBoatTypeXII:
		return "XII"
	case UBoatTypeXIV:
		return "XIV"
	case UBoatTypeXXI:
		return "XXI"
	default:
		return "Invalid U-Boat Type"
	}
}

func (u UBoatType) FirstPatrolDate() (PatrolDate, error) {
	if err := u.Validate(); err != nil {
		return 0, err
	}
	switch u {
	case UBoatTypeVIIB, UBoatTypeVIIC, UBoatTypeVIICFlak, UBoatTypeVIID, UBoatTypeIXB,
		UBoatTypeIXC, UBoatTypeIXC40, UBoatTypeIXD2, UBoatTypeXB, UBoatTypeXII, UBoatTypeXIV:
		return PatrolDateJul43, nil
	case UBoatTypeVIIC41:
		return PatrolDateApr44, nil
	case UBoatTypeIXD42:
		return PatrolDateMar45, nil
	case UBoatTypeXXI:
		return PatrolDateFeb45, nil
	default:
		panic(fmt.Sprintf("unhandled u-boat type %v", u))
	}
}