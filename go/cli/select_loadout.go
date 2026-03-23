package main

import (
	"fmt"
	"log"
	"slices"

	"charm.land/huh/v2"
	thehunted "github.com/krelinga/the-hunted/go"
	"github.com/krelinga/the-hunted/go/views"
)

func (_ selector) SelectLoadout(g thehunted.View) *thehunted.SelectedLoadout {
	uboatLoadouts := thehunted.PermuteLoadouts(g.GetUBoat().GetUBoatType().DefaultLoadout(g.GetStartPatrolDate()).View())
	loadoutOptions := []huh.Option[int]{}
	for i, loadout := range uboatLoadouts {
		loadoutOptions = append(loadoutOptions, huh.NewOption(loadout.String(), i))
	}

	tubeLocs := []thehunted.TorpLoc{}
	for loc := range g.GetUBoat().GetUBoatType().TorpLocs() {
		if loc.IsTube() {
			tubeLocs = append(tubeLocs, loc)
		}
	}
	slices.SortFunc(tubeLocs, func(a, b thehunted.TorpLoc) int {
		if a.Facing() != b.Facing() {
			if a.Facing() == thehunted.FacingFwd {
				return -1
			}
			return 1
		}
		aTube, aIsTube := a.Tube()
		bTube, bIsTube := b.Tube()
		if !aIsTube || !bIsTube {
			panic("invalid torpedo location: not a tube")
		}
		return int(aTube) - int(bTube)
	})
	selectedTorps := make([]thehunted.TorpType, len(tubeLocs))
	var selectedLoadoutOptionIdx int
	groups := []*huh.Group{
		huh.NewGroup(
			huh.NewSelect[int]().
				Options(loadoutOptions...).
				Value(&selectedLoadoutOptionIdx),
		).Title("Overall Loadout"),
	}
	for i, loc := range tubeLocs {
		groups = append(groups, huh.NewGroup(
			huh.NewSelect[thehunted.TorpType]().
				OptionsFunc(func() []huh.Option[thehunted.TorpType] {
					overall := views.MapClone(uboatLoadouts[selectedLoadoutOptionIdx])
					for j := 0; j < len(selectedTorps) && j < i; j++ {
						overall[selectedTorps[j]]--
					}
					torpTypes := []thehunted.TorpType{}
					for torpType, count := range overall {
						if count > 0 {
							torpTypes = append(torpTypes, torpType)
						}
					}
					slices.Sort(torpTypes)
					options := []huh.Option[thehunted.TorpType]{}
					for _, torpType := range torpTypes {
						options = append(options, huh.NewOption(torpType.String(), torpType))
					}
					return options
				}, nil).
				Value(&selectedTorps[i]),
		).Title(fmt.Sprintf("Select Torpedo for %s", loc)))
	}
	var aftReloadSlots []thehunted.TorpType
	if capacity := g.GetUBoat().GetUBoatType().AftReloads(); capacity > 0 {
		optFunc := func() []huh.Option[thehunted.TorpType] {
			overall := views.MapClone(uboatLoadouts[selectedLoadoutOptionIdx])
			for j := 0; j < len(selectedTorps); j++ {
				overall[selectedTorps[j]]--
			}
			torpTypes := []thehunted.TorpType{}
			for torpType, count := range overall {
				for j := 0; j < count && j < capacity; j++ {
					torpTypes = append(torpTypes, torpType)
				}
			}
			slices.Sort(torpTypes)
			options := []huh.Option[thehunted.TorpType]{}
			for _, torpType := range torpTypes {
				// TODO: using the same torpType more than once causes a bug in the form.
				// We need to figure out a way to allow multiple selection of identical values without breaking the form.
				options = append(options, huh.NewOption(torpType.String(), torpType))
			}
			return options
		}
		if capacity == 1 {
			aftReloadSlots = make([]thehunted.TorpType, 1)
			groups = append(groups, huh.NewGroup(
				huh.NewSelect[thehunted.TorpType]().
					OptionsFunc(optFunc, nil).
					Value(&aftReloadSlots[0]),
			).Title("Select Torpedo for Aft Reload"))
		} else {
			groups = append(groups, huh.NewGroup(
				huh.NewMultiSelect[thehunted.TorpType]().
					OptionsFunc(optFunc, nil).
					Value(&aftReloadSlots).
					Validate(func(got []thehunted.TorpType) error {
						if len(got) != capacity {
							return fmt.Errorf("must select exactly %d torpedoes for aft reload, got %d", capacity, len(got))
						}
						return nil
					}),
			).Title(fmt.Sprintf("Select %d Torpedos for Aft Reload", capacity)))
		}
	}

	huhForm := huh.NewForm(groups...)
	if err := huhForm.Run(); err != nil {
		log.Fatalf("error running form: %v", err)
	}

	layout := map[thehunted.TorpLoc]thehunted.TorpCountsData{}
	fwdReloads := thehunted.TorpCountsData(views.MapClone(uboatLoadouts[selectedLoadoutOptionIdx]))
	for i, loc := range tubeLocs {
		layout[loc] = thehunted.TorpCountsData{selectedTorps[i]: 1}
		fwdReloads[selectedTorps[i]]--
	}
	if len(aftReloadSlots) > 0 {
		layout[thehunted.NewTorpLocReload(thehunted.FacingAft)] = thehunted.TorpCountsData{}
	}
	for _, torpType := range aftReloadSlots {
		layout[thehunted.NewTorpLocReload(thehunted.FacingAft)][torpType]++
		fwdReloads[torpType]--
	}
	layout[thehunted.NewTorpLocReload(thehunted.FacingFwd)] = fwdReloads

	return &thehunted.SelectedLoadout{Layout: layout}
}
