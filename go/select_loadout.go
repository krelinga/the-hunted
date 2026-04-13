package thehunted

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

type SelectedLoadout struct {
	Layout map[TorpLoc]TorpCountsData
}

func (s *SelectedLoadout) Validate(g View) error {
	for torpLoc, torpCounts := range s.Layout {
		if !g.GetUBoat().GetUBoatType().HasTorpLoc(torpLoc) {
			return fmt.Errorf("%w: invalid torpedo location %s", ErrInvalidSelection, torpLoc)
		}
		if torpLoc.IsTube() && torpCounts.Total() > 1 {
			return fmt.Errorf("%w: torpedo location %s is a tube and cannot have more than 1 torpedo", ErrInvalidSelection, torpLoc)
		} else {
			switch torpLoc.Facing() {
			case FacingFwd:
				if torpCounts.Total() > g.GetUBoat().GetUBoatType().FwdReloads() {
					return fmt.Errorf("%w: total count %d for torpedo location %s exceeds forward reload capacity of %d", ErrInvalidSelection, torpCounts.Total(), torpLoc, g.GetUBoat().GetUBoatType().FwdReloads())
				}
			case FacingAft:
				if torpCounts.Total() > g.GetUBoat().GetUBoatType().AftReloads() {
					return fmt.Errorf("%w: total count %d for torpedo location %s exceeds aft reload capacity of %d", ErrInvalidSelection, torpCounts.Total(), torpLoc, g.GetUBoat().GetUBoatType().AftReloads())
				}
			default:
				panic("invalid facing")
			}
		}
		for torpType, count := range torpCounts {
			if count < 0 {
				return fmt.Errorf("%w: negative count %d for torpedo type %s at location %s", ErrInvalidSelection, count, torpType, torpLoc)
			}
		}
	}
	return nil

}

type LoadoutChangedEvent struct {
	baseEvent
	TorpLoc TorpLoc
	TorpCounts TorpCountsView
}

func (e LoadoutChangedEvent) apply(gd *Data) {
	for k, v := range e.TorpCounts.All() {
		gd.UBoat.TorpLayout[e.TorpLoc][k] = v
	}
}

func (e LoadoutChangedEvent) String() string {
	types := slices.Collect(e.TorpCounts.Keys())
	slices.Sort(types)
	deltas := []string{}
	for _, torpType := range types {
		delta, _ := e.TorpCounts.Find(torpType)
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

func handleSelectLoadout(g View, s Selector, r Roller, ew EventWriter) (gameState, error) {
	selected := s.SelectLoadout(g)
	if selected == nil {
		return 0, errNoChange
	}
	if err := selected.Validate(g); err != nil {
		return 0, err
	}

	locs := slices.Collect(maps.Keys(selected.Layout))
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
		ew.WriteEvent(LoadoutChangedEvent{
			TorpLoc: loc,
			TorpCounts: selected.Layout[loc],
		})
	}
	startPatrol(g, r, ew)
	// TODO: implement more.
	return gameStateDone, nil
}

func PermuteLoadouts(defLoadout TorpCountsView) []TorpCountsView {
	var extraE, extraA []TorpCountsView

	for i := 1; i <= 4; i++ {
		if count, _ := defLoadout.Find(TorpTypeG7a); count-i >= 0 {
			l := defLoadout.Clone()
			l[TorpTypeG7a] -= i
			l[TorpTypeG7e] += i
			extraE = append(extraE, l)
		}
		if count, _ := defLoadout.Find(TorpTypeG7e); count-i >= 0 {
			l := defLoadout.Clone()
			l[TorpTypeG7e] -= i
			l[TorpTypeG7a] += i
			extraA = append(extraA, l)
		}
	}
	slices.Reverse(extraE)
	var options []TorpCountsView
	options = append(options, extraE...)
	options = append(options, defLoadout)
	options = append(options, extraA...)
	return options
}