package main

import (
	"fmt"

	"charm.land/huh/v2"
	thehunted "github.com/krelinga/the-hunted/go"
)

func handleSelectLoadout(form *thehunted.SelectLoadoutForm) error {
	loadoutOptions := []huh.Option[int]{}
	for i, loadout := range form.Loadout.Options {
		loadoutOptions = append(loadoutOptions, huh.NewOption(loadout.String(), i))
	}
	huhForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Select Loadout").
				Options(loadoutOptions...).
				Value(&form.Loadout.Selected),
		),
	)
	if err := huhForm.Run(); err != nil {
		return fmt.Errorf("error running form: %w", err)
	}
	return nil
}