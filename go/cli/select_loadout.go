package main

import (
	"fmt"
	"maps"
	"slices"

	"charm.land/huh/v2"
	thehunted "github.com/krelinga/the-hunted/go"
)

func handleSelectLoadout(form *thehunted.SelectLoadoutForm) error {
	loadoutOptions := []huh.Option[int]{}
	for i, loadout := range form.Overall.Options {
		loadoutOptions = append(loadoutOptions, huh.NewOption(loadout.String(), i))
	}

	tubeLocs := []thehunted.TorpLoc{}
	for loc := range form.Layout {
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
	groups := []*huh.Group{
		huh.NewGroup(
			huh.NewSelect[int]().
				Options(loadoutOptions...).
				Value(&form.Overall.Selected),
		).Title("Overall Loadout"),
	}
	for i, loc := range tubeLocs {
		groups = append(groups, huh.NewGroup(
			huh.NewSelect[thehunted.TorpType]().
				OptionsFunc(func() []huh.Option[thehunted.TorpType] {
					overall := maps.Clone(form.Overall.Options[form.Overall.Selected])
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
	if aftReloads, ok := form.Layout[thehunted.NewTorpLocReload(thehunted.FacingAft)]; ok {
		capacity := aftReloads.Capacity
		optFunc := func() []huh.Option[thehunted.TorpType] {
			overall := maps.Clone(form.Overall.Options[form.Overall.Selected])
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
		return fmt.Errorf("error running form: %w", err)
	}

	fwdReloads := form.Layout[thehunted.NewTorpLocReload(thehunted.FacingFwd)].Items
	maps.Copy(fwdReloads, form.Overall.Options[form.Overall.Selected])
	for i, loc := range tubeLocs {
		form.Layout[loc].Items[selectedTorps[i]]++
		fwdReloads[selectedTorps[i]]--
	}
	for _, torpType := range aftReloadSlots {
		form.Layout[thehunted.NewTorpLocReload(thehunted.FacingAft)].Items[torpType]++
		fwdReloads[torpType]--
	}

	return nil
}
