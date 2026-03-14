package thehunted

import (
	"errors"
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

func (u UBoatType) Must() {
	if err := u.Validate(); err != nil {
		panic(err)
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

func (u UBoatType) FirstPatrolDate() PatrolDate {
	u.Must()
	switch u {
	case UBoatTypeVIIB, UBoatTypeVIIC, UBoatTypeVIICFlak, UBoatTypeVIID, UBoatTypeIXB,
		UBoatTypeIXC, UBoatTypeIXC40, UBoatTypeIXD2, UBoatTypeXB, UBoatTypeXII, UBoatTypeXIV:
		return PatrolDateJul43
	case UBoatTypeVIIC41:
		return PatrolDateApr44
	case UBoatTypeIXD42:
		return PatrolDateMar45
	case UBoatTypeXXI:
		return PatrolDateApr45
	default:
		panic("unreachable code")
	}
}

func (u UBoatType) FwdTubes() int {
	u.Must()
	switch u {
	case UBoatTypeXB, UBoatTypeXIV:
		return 0
	case UBoatTypeXII, UBoatTypeXXI:
		return 6
	default:
		return 4
	}
}

func (u UBoatType) AftTubes() int {
	u.Must()
	switch u {
	case UBoatTypeXIV, UBoatTypeXXI:
		return 0
	case UBoatTypeVIIB, UBoatTypeVIIC, UBoatTypeVIIC41, UBoatTypeVIID, UBoatTypeVIICFlak:
		return 1
	default:
		return 2
	}
}

func (u UBoatType) FwdReloads() int {
	u.Must()
	switch u {
	case UBoatTypeXB, UBoatTypeXIV:
		return 0
	case UBoatTypeXII:
		return 10
	case UBoatTypeIXB, UBoatTypeIXC, UBoatTypeIXC40, UBoatTypeIXD42:
		return 14
	case UBoatTypeIXD2:
		return 16
	case UBoatTypeXXI:
		return 17
	default:
		return 8
	}
}

func (u UBoatType) AftReloads() int {
	u.Must()
	switch u {
	case UBoatTypeXIV, UBoatTypeXXI:
		return 0
	case UBoatTypeVIIB, UBoatTypeVIIC, UBoatTypeVIIC41, UBoatTypeVIID, UBoatTypeVIICFlak:
		return 1
	case UBoatTypeXB:
		return 3
	default:
		return 2
	}
}

func (u UBoatType) DeckGunAmmo() int {
	u.Must()
	switch u {
	case UBoatTypeIXB, UBoatTypeIXC, UBoatTypeIXC40:
		return 5
	case UBoatTypeIXD2:
		return 6
	case UBoatTypeIXD42:
		return 7
	case UBoatTypeVIIB, UBoatTypeVIIC, UBoatTypeVIIC41, UBoatTypeVIID, UBoatTypeXII:
		return 10
	default:
		return 0
	}
}

func (u UBoatType) HasDeckGun() bool {
	return u.DeckGunAmmo() > 0
}

type UBoat struct {
	UBoatType              UBoatType
	ID                     string
	FwdTubes, AftTubes     []*Torpedo
	FwdReloads, AftReloads map[Torpedo]int
	HasDeckGun             bool
	DeckGunAmmo            int
}

func NewUBoat(uBoatType UBoatType, id string) UBoat {
	ub := UBoat{
		UBoatType: uBoatType,
		ID:        id,
		FwdTubes:  make([]*Torpedo, uBoatType.FwdTubes()),
		AftTubes:  make([]*Torpedo, uBoatType.AftTubes()),
		FwdReloads: make(map[Torpedo]int),
		AftReloads: make(map[Torpedo]int),
		HasDeckGun: uBoatType.HasDeckGun(),
		DeckGunAmmo: uBoatType.DeckGunAmmo(),
	}
	return ub
}

type DeckGunRemovedEvent struct {
	baseEvent
}

func (e DeckGunRemovedEvent) String() string {
	return "Deck gun removed"
}

func (u *UBoat) RemoveDeckGun() []Event {
	if !u.HasDeckGun {
		return nil
	}
	u.HasDeckGun = false
	u.DeckGunAmmo = 0
	return []Event{DeckGunRemovedEvent{}}
}
