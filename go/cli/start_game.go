package main

import (
	"fmt"

	"charm.land/huh/v2"
	thehunted "github.com/krelinga/the-hunted/go"
)

func handleStartGame(form *thehunted.StartGameForm) error {
	uboatOptions := []huh.Option[int]{}
	for i, option := range form.UBoatType.Options {
		uboatOptions = append(uboatOptions, huh.NewOption(option.String(), i))
	}
	var kmdtName, uBoatID string
	var uboatTypeIndex int
	huhForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Kommandant Name").
				Value(&kmdtName),
			huh.NewInput().
				Title("U-boat ID").
				Value(&uBoatID),
			huh.NewSelect[int]().
				Title("U-boat Type").
				Options(uboatOptions...).
				Value(&uboatTypeIndex),
		),
	)
	if err := huhForm.Run(); err != nil {
		return fmt.Errorf("error running form: %w", err)
	}
	form.KmdtName = thehunted.TextFormField(kmdtName)
	form.UBoatID = thehunted.TextFormField(uBoatID)
	form.UBoatType.Selected = uboatTypeIndex
	return nil
}
