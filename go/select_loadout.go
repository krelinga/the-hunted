package thehunted

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

type SelectLoadoutForm struct {
	baseForm
	Overall SelectFormField[Loadout]
	Layout LayoutFormField[TorpLoc, TorpType]
}

func (f *SelectLoadoutForm) Validate() error {
	if err := f.Overall.Validate(); err != nil {
		return fmt.Errorf("%w: invalid loadout selection", err)
	}
	if err := f.Layout.Validate(); err != nil {
		return fmt.Errorf("%w: invalid layout selection", err)
	}

	totals := map[TorpType]int{}
	for _, locForm := range f.Layout {
		for torpType, count := range locForm.Items {
			totals[torpType] += count
		}
	}
	selectedLoadout := f.Overall.Options[f.Overall.Selected]
	for torpType, total := range totals {
		if total != selectedLoadout[torpType] {
			return fmt.Errorf("%w: total count %d for torpedo type %s does not match selected loadout count %d", ErrInvalidFormField, total, torpType, selectedLoadout[torpType])
		}
	}

	return nil
}

func (g *Game) formForSelectLoadout() Form {
	defLoadout := g.UBoat.UBoatType.DefaultLoadout(g.startPatrolDate)
	var extraE, extraA []Loadout

	for i := 1; i <= 4; i++ {
		if defLoadout[TorpTypeG7a]-i >= 0 {
			l := maps.Clone(defLoadout)
			l[TorpTypeG7a] -= i
			l[TorpTypeG7e] += i
			extraE = append(extraE, l)
		}
		if defLoadout[TorpTypeG7e]-i >= 0 {
			l := maps.Clone(defLoadout)
			l[TorpTypeG7e] -= i
			l[TorpTypeG7a] += i
			extraA = append(extraA, l)
		}
	}
	slices.Reverse(extraE)
	var options []Loadout
	options = append(options, extraE...)
	options = append(options, defLoadout)
	options = append(options, extraA...)

	form := &SelectLoadoutForm{
		Overall: SelectFormField[Loadout]{
			Options: options,
			Selected: len(extraE),
		},
		Layout: LayoutFormField[TorpLoc, TorpType]{},
	}
	if fwdTubes := g.UBoat.UBoatType.FwdTubes(); fwdTubes > 0 {
		for i := 1; i <= fwdTubes; i++ {
			form.Layout[NewTorpLocTube(FacingFwd, i)] = NewLayoutFormLoc[TorpType](1)
		}
	}
	if aftTubes := g.UBoat.UBoatType.AftTubes(); aftTubes > 0 {
		for i := 1; i <= aftTubes; i++ {
			form.Layout[NewTorpLocTube(FacingAft, i)] = NewLayoutFormLoc[TorpType](1)
		}
	}
	if fwdReloads := g.UBoat.UBoatType.FwdReloads(); fwdReloads > 0 {
		form.Layout[NewTorpLocReload(FacingFwd)] = NewLayoutFormLoc[TorpType](fwdReloads)
	}
	if aftReloads := g.UBoat.UBoatType.AftReloads(); aftReloads > 0 {
		form.Layout[NewTorpLocReload(FacingAft)] = NewLayoutFormLoc[TorpType](aftReloads)
	}

	return form
}

type LoadoutChangedEvent struct {
	baseEvent
	TorpLoc TorpLoc
	Loadout Loadout
}

func (e LoadoutChangedEvent) String() string {
	types := slices.Collect(maps.Keys(e.Loadout))
	slices.Sort(types)
	deltas := []string{}
	for _, torpType := range types {
		delta := e.Loadout[torpType]
		if delta == 0 {
			continue
		}
		var verb string
		if delta > 0 {
			verb = "Added"
		} else {
			verb = "Removed"
			delta = -delta
		}
		deltas = append(deltas, fmt.Sprintf("%s %d %s", verb, delta, torpType))
	}
	return fmt.Sprintf("Changed loadout for %s: %s", e.TorpLoc, strings.Join(deltas, ", "))
}


func (g *Game) advanceFromSelectLoadout(form Form) error {
	selectLoadoutForm, ok := form.(*SelectLoadoutForm)
	if !ok {
		return fmt.Errorf("%w: expected *SelectLoadoutForm, got %T", ErrUnexpectedForm, form)
	}
	if err := selectLoadoutForm.Validate(); err != nil {
		return err
	}
	locs := slices.Collect(maps.Keys(selectLoadoutForm.Layout))
	slices.SortFunc(locs, func(a, b TorpLoc) int {
		if a.IsTube() != b.IsTube() {
			if a.IsTube() {
				return -1
			}
			return 1
		}
		if a.Facing() != b.Facing() {
			if a.Facing() == FacingFwd {
				return -1
			}
			return 1
		}
		if a.IsTube() && b.IsTube() {
			aTube, _ := a.Tube()
			bTube, _ := b.Tube()
			return aTube - bTube
		}
		return 0
	})
	for _, loc := range locs {
		loadout := Loadout(selectLoadoutForm.Layout[loc].Items)
		g.UBoat.Torpedos[loc] = maps.Clone(loadout)
		g.writeEvent(LoadoutChangedEvent{
			TorpLoc: loc,
			Loadout: maps.Clone(loadout),
		})
	}
	g.gameState = GameStateInPort
	g.writeEvent(GameStateSetEvent{GameState: g.gameState})
	return nil
}
