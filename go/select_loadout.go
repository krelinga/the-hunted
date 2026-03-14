package thehunted

import (
	"fmt"
	"maps"
	"slices"
)

type SelectLoadoutForm struct {
	baseForm
	Loadout SelectFormField[Loadout]
}

func (f *SelectLoadoutForm) Validate() error {
	if err := f.Loadout.Validate(); err != nil {
		return fmt.Errorf("%w: invalid loadout selection", err)
	}
	return nil
}

type LoadoutSelectedEvent struct {
	baseEvent
	Loadout Loadout
}

func (e LoadoutSelectedEvent) String() string {
	return fmt.Sprintf("Loadout selected: %s", e.Loadout)
}

func (g *Game) formForSelectLoadout() Form {
	defLoadout := g.UBoat.UBoatType.DefaultLoadout(g.startPatrolDate)
	var extraE, extraA []Loadout

	for i := 1; i <= 4; i++ {
		if defLoadout[TorpedoG7a]-i >= 0 {
			l := maps.Clone(defLoadout)
			l[TorpedoG7a] -= i
			l[TorpedoG7e] += i
			extraE = append(extraE, l)
		}
		if defLoadout[TorpedoG7e]-i >= 0 {
			l := maps.Clone(defLoadout)
			l[TorpedoG7e] -= i
			l[TorpedoG7a] += i
			extraA = append(extraA, l)
		}
	}
	slices.Reverse(extraE)
	var options []Loadout
	options = append(options, extraE...)
	options = append(options, defLoadout)
	options = append(options, extraA...)

	return &SelectLoadoutForm{
		Loadout: SelectFormField[Loadout]{
			Options: options,
		},
	}
}

func (g *Game) advanceFromSelectLoadout(form Form) ([]Event, error) {
	selectLoadoutForm, ok := form.(*SelectLoadoutForm)
	if !ok {
		return nil, fmt.Errorf("%w: expected *SelectLoadoutForm, got %T", ErrUnexpectedForm, form)
	}
	if err := selectLoadoutForm.Validate(); err != nil {
		return nil, err
	}
	events := []Event{}
	loadout := selectLoadoutForm.Loadout.Options[selectLoadoutForm.Loadout.Selected]
	g.nextLoadout = loadout
	events = append(events, LoadoutSelectedEvent{Loadout: loadout})
	g.gameState = GameStateInPort
	events = append(events, GameStateSetEvent{GameState: g.gameState})
	return events, nil
}
