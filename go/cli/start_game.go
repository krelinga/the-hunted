package main

import (
	"log"

	"charm.land/huh/v2"
	thehunted "github.com/krelinga/the-hunted/go"
)

func (_ selector) SelectStart(g thehunted.GameView) *thehunted.SelectedStart {
	uboatTypes := thehunted.AllUBoatTypes()
	uboatOptions := []huh.Option[int]{}
	for i, option := range uboatTypes.All() {
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
		log.Fatalf("error running form: %v", err)
	}
	return &thehunted.SelectedStart{
		KmdtName:  kmdtName,
		UBoatID:   uBoatID,
		UBoatType: uboatTypes.Get(uboatTypeIndex),
	}
}
