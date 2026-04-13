package thehunted

import (
	"errors"
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"
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

type UBoatTypesView interface {
	Length() int
	Get(i int) UBoatType
	All() iter.Seq2[int, UBoatType]
}

type UBoatTypes []UBoatType

func (u UBoatTypes) Length() int {
	return len(u)
}

func (u UBoatTypes) Get(i int) UBoatType {
	return u[i]
}

func (u UBoatTypes) All() iter.Seq2[int, UBoatType] {
	return slices.All(u)
}

var allUBoatTypes = UBoatTypes{
	UBoatTypeVIIB,
	UBoatTypeVIIC,
	UBoatTypeVIICFlak,
	UBoatTypeVIIC41,
	UBoatTypeVIID,
	UBoatTypeIXB,
	UBoatTypeIXC,
	UBoatTypeIXC40,
	UBoatTypeIXD2,
	UBoatTypeIXD42,
	UBoatTypeXB,
	UBoatTypeXII,
	UBoatTypeXIV,
	UBoatTypeXXI,
}

func AllUBoatTypes() UBoatTypesView {
	return allUBoatTypes
}

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

func (u UBoatType) HasTorpLoc(loc TorpLoc) bool {
	u.Must()
	if loc.IsTube() {
		switch loc.Facing() {
		case FacingFwd:
			return loc.tube <= u.FwdTubes()
		case FacingAft:
			return loc.tube <= u.AftTubes()
		default:
			panic("invalid facing")
		}
	} else {
		switch loc.Facing() {
		case FacingFwd:
			return u.FwdReloads() > 0
		case FacingAft:
			return u.AftReloads() > 0
		default:
			panic("invalid facing")
		}
	}
}

type TorpCountsView interface {
	Length() int
	Find(torpType TorpType) (int, bool)
	All() iter.Seq2[TorpType, int]
	Keys() iter.Seq[TorpType]
	Clone() TorpCounts
	Equal(TorpCountsView) bool
	String() string
	Total() int
}

type TorpCounts map[TorpType]int

func (d TorpCounts) Length() int {
	return len(d)
}

func (d TorpCounts) Find(torpType TorpType) (int, bool) {
	count, ok := d[torpType]
	return count, ok
}

func (d TorpCounts) All() iter.Seq2[TorpType, int] {
	return maps.All(d)
}

func (d TorpCounts) Keys() iter.Seq[TorpType] {
	return maps.Keys(d)
}

func (d TorpCounts) Clone() TorpCounts {
	return maps.Clone(d)
}

func (d TorpCounts) Equal(other TorpCountsView) bool {
	if other == nil {
		return false
	}
	if d.Length() != other.Length() {
		return false
	}
	for torpType, count := range d {
		if otherCount, _ := other.Find(torpType); count != otherCount {
			return false
		}
	}
	return true
}

func (d TorpCounts) Total() int {
	total := 0
	for _, count := range d {
		total += count
	}
	return total
}

func (d TorpCounts) String() string {
	kinds := []TorpType{TorpTypeG7e, TorpTypeG7a, TorpTypeG7esZaunkonig, TorpTypeG7esZaunkonigII, TorpTypeG7eFalke}
	var parts []string
	for _, kind := range kinds {
		if count, ok := d[kind]; ok && count > 0 {
			parts = append(parts, fmt.Sprintf("%dx %s", count, kind))
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
}

type TorpLayoutView interface {
	Find(loc TorpLoc) (TorpCountsView, bool)
	All() iter.Seq2[TorpLoc, TorpCountsView]
	Keys() iter.Seq[TorpLoc]
	Total() int
}

type TorpLayout map[TorpLoc]TorpCounts

func (d TorpLayout) Find(loc TorpLoc) (TorpCountsView, bool) {
	counts, ok := d[loc]
	return counts, ok
}

func (d TorpLayout) All() iter.Seq2[TorpLoc, TorpCountsView] {
	return func(yield func(TorpLoc, TorpCountsView) bool) {
		for loc, counts := range d {
			if !yield(loc, counts) {
				return
			}
		}
	}
}

func (d TorpLayout) Keys() iter.Seq[TorpLoc] {
	return func(yield func(TorpLoc) bool) {
		for loc := range d {
			if !yield(loc) {
				return
			}
		}
	}
}

func (d TorpLayout) Total() int {
	total := 0
	for _, counts := range d {
		total += counts.Total()
	}
	return total
}

func (u UBoatType) DefaultLoadout(pd PatrolDate) TorpCountsView {
	u.Must()
	pd.Must()

	var special TorpType
	switch {
	case pd < PatrolDateAug43:
		special = TorpTypeG7eFalke
	case pd < PatrolDateApr45:
		special = TorpTypeG7esZaunkonig
	default:
		special = TorpTypeG7esZaunkonigII
	}

	switch u {
	case UBoatTypeVIIB, UBoatTypeVIIC, UBoatTypeVIIC41, UBoatTypeVIID:
		return TorpCounts{
			TorpTypeG7e: 8,
			TorpTypeG7a: 4,
			special:     2,
		}
	case UBoatTypeVIICFlak:
		return TorpCounts{
			TorpTypeG7a: 3,
			special:     2,
		}
	case UBoatTypeIXB, UBoatTypeIXC, UBoatTypeIXC40, UBoatTypeIXD42:
		return TorpCounts{
			TorpTypeG7e: 10,
			TorpTypeG7a: 10,
			special:     2,
		}
	case UBoatTypeIXD2:
		return TorpCounts{
			TorpTypeG7e: 10,
			TorpTypeG7a: 10,
			special:     4,
		}
	case UBoatTypeXB:
		return TorpCounts{
			TorpTypeG7e: 3,
			special:     2,
		}
	case UBoatTypeXII:
		return TorpCounts{
			TorpTypeG7e: 8,
			TorpTypeG7a: 8,
			special:     4,
		}
	case UBoatTypeXIV:
		return nil
	case UBoatTypeXXI:
		return TorpCounts{
			TorpTypeG7e: 8,
			TorpTypeG7a: 9,
			special:     6,
		}
	default:
		panic("unreachable code")
	}
}

func (u UBoatType) IsMinelayer() bool {
	u.Must()
	switch u {
	case UBoatTypeVIID, UBoatTypeXB:
		return true
	default:
		return false
	}
}

func (u UBoatType) IsTypeVII() bool {
	u.Must()
	switch u {
	case UBoatTypeVIIB, UBoatTypeVIIC, UBoatTypeVIICFlak, UBoatTypeVIIC41, UBoatTypeVIID:
		return true
	default:
		return false
	}
}

func (u UBoatType) IsTypeIX() bool {
	u.Must()
	switch u {
	case UBoatTypeIXB, UBoatTypeIXC, UBoatTypeIXC40, UBoatTypeIXD2, UBoatTypeIXD42:
		return true
	default:
		return false
	}
}

func (u UBoatType) IsAnyOf(types ...UBoatType) bool {
	return slices.Contains(types, u)
}

func (u UBoatType) TorpLocs() iter.Seq[TorpLoc] {
	u.Must()
	return func(yield func(TorpLoc) bool) {
		for tube := 1; tube <= u.FwdTubes(); tube++ {
			if !yield(NewTorpLocTube(FacingFwd, tube)) {
				return
			}
		}
		for tube := 1; tube <= u.AftTubes(); tube++ {
			if !yield(NewTorpLocTube(FacingAft, tube)) {
				return
			}
		}
		if u.FwdReloads() > 0 {
			if !yield(NewTorpLocReload(FacingFwd)) {
				return
			}
		}
		if u.AftReloads() > 0 {
			if !yield(NewTorpLocReload(FacingAft)) {
				return
			}
		}
	}
}

type UBoatView interface {
	GetUBoatType() UBoatType
	GetID() string
	GetTorpLayout() TorpLayoutView
	GetHasDeckGun() bool
	GetDeckGunAmmo() int
}

type UBoat struct {
	UBoatType   UBoatType
	ID          string
	TorpLayout  TorpLayout
	HasDeckGun  bool
	DeckGunAmmo int
}

func (d *UBoat) GetUBoatType() UBoatType {
	return d.UBoatType
}

func (d *UBoat) GetID() string {
	return d.ID
}

func (d *UBoat) GetTorpLayout() TorpLayoutView {
	return d.TorpLayout
}

func (d *UBoat) GetHasDeckGun() bool {
	return d.HasDeckGun
}

func (d *UBoat) GetDeckGunAmmo() int {
	return d.DeckGunAmmo
}

func NewUBoat(uBoatType UBoatType, id string) *UBoat {
	ub := &UBoat{
		UBoatType:   uBoatType,
		ID:          id,
		TorpLayout:  make(TorpLayout),
		HasDeckGun:  uBoatType.HasDeckGun(),
		DeckGunAmmo: uBoatType.DeckGunAmmo(),
	}
	for i := 1; i <= uBoatType.FwdTubes(); i++ {
		ub.TorpLayout[NewTorpLocTube(FacingFwd, i)] = TorpCounts{}
	}
	for i := 1; i <= uBoatType.AftTubes(); i++ {
		ub.TorpLayout[NewTorpLocTube(FacingAft, i)] = TorpCounts{}
	}
	if uBoatType.FwdReloads() > 0 {
		ub.TorpLayout[NewTorpLocReload(FacingFwd)] = TorpCounts{}
	}
	if uBoatType.AftReloads() > 0 {
		ub.TorpLayout[NewTorpLocReload(FacingAft)] = TorpCounts{}
	}
	return ub
}
