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

type UBoat struct {
	UBoatType                UBoatType
	ID                       string
	FwdTubes, AftTubes       []*Torpedo
	FwdReloads, AftReloads   map[Torpedo]int
	FwdCapacity, AftCapacity int
	HasDeckGun               bool
	DeckGunAmmo              int
}

func NewUBoat(uBoatType UBoatType, id string) UBoat {
	ub := UBoat{
		UBoatType: uBoatType,
		ID:        id,
	}
	fwd := func(tubes, capacity int) {
		ub.FwdTubes = make([]*Torpedo, tubes)
		ub.FwdReloads = make(map[Torpedo]int)
		ub.FwdCapacity = capacity
	}
	aft := func(tubes, capacity int) {
		ub.AftTubes = make([]*Torpedo, tubes)
		ub.AftReloads = make(map[Torpedo]int)
		ub.AftCapacity = capacity
	}
	deckGun := func(ammo int) {
		ub.HasDeckGun = true
		ub.DeckGunAmmo = ammo
	}
	switch uBoatType {
	case UBoatTypeVIIB, UBoatTypeVIIC, UBoatTypeVIIC41, UBoatTypeVIID:
		fwd(4, 8)
		aft(1, 1)
		deckGun(10)
	case UBoatTypeVIICFlak:
		fwd(4, 8)
		aft(1, 1)
	case UBoatTypeIXB, UBoatTypeIXC, UBoatTypeIXC40:
		fwd(4, 14)
		aft(2, 2)
		deckGun(5)
	case UBoatTypeIXD42:
		fwd(4, 14)
		aft(2, 2)
		deckGun(7)
	case UBoatTypeIXD2:
		fwd(4, 16)
		aft(2, 2)
		deckGun(6)
	case UBoatTypeXB:
		fwd(0, 0)
		aft(2, 3)
		deckGun(8)
	case UBoatTypeXII:
		fwd(6, 10)
		aft(2, 2)
		deckGun(10)
	case UBoatTypeXIV:
		fwd(0, 0)
		aft(0, 0)
	case UBoatTypeXXI:
		fwd(6, 17)
		aft(0, 0)
	}
	return ub
}
