package thehunted

import (
	"errors"
)

type View interface {
	GetKmdtName() string
	GetKmdtRank() Rank
	GetCrewQuality() CrewQuality
	GetUBoat() UBoatView
	GetPatrols() PatrolsView
	GetStartPatrolDate() PatrolDate
}

type Data struct {
	KmdtName    string
	KmdtRank    Rank
	CrewQuality CrewQuality
	UBoat       *UBoatData
	Patrols     Patrols
	// TODO: rename to NextPatrolDate.
	StartPatrolDate PatrolDate
}

func (v *Data) GetKmdtName() string {
	return v.KmdtName
}

func (v *Data) GetKmdtRank() Rank {
	return v.KmdtRank
}

func (v *Data) GetCrewQuality() CrewQuality {
	return v.CrewQuality
}

func (v *Data) GetUBoat() UBoatView {
	return v.UBoat
}

func (v *Data) GetPatrols() PatrolsView {
	return v.Patrols
}

func (v *Data) GetStartPatrolDate() PatrolDate {
	return v.StartPatrolDate
}

type Game struct {
	Selector    Selector
	EventWriter EventWriter
	Roller      Roller

	data      Data
	nextState gameState
}

var errNoChange = errors.New("no change in game state")

func (g *Game) Next() error {
	if g.Done() {
		panic("game is already done")
	}
	roller := g.Roller
	if roller == nil {
		roller = RandomRoller{}
	}
	ew := applyEventToGame{G: g}
	newState, err := allHandlers[g.nextState](&g.data, g.Selector, roller, ew)
	if err == errNoChange {
		return nil
	} else if err != nil {
		return err
	}
	g.nextState = newState
	return nil
}

func (g *Game) Done() bool {
	return g.nextState == gameStateDone
}

var ErrUnexpectedForm = errors.New("unexpected form")

type gameState int

const (
	gameStateStart gameState = iota
	gameStateSelectLoadout
	gameStateStartPatrol
	gameStateDone
)

type handler func(g View, s Selector, r Roller, ew EventWriter) (gameState, error)

var allHandlers = map[gameState]handler{
	gameStateStart:         handleStart,
	gameStateSelectLoadout: handleSelectLoadout,
}

type applyEventToGame struct {
	G *Game
}

func (e applyEventToGame) WriteEvent(ev Event) {
	ev.apply(&e.G.data)
	if e.G.EventWriter != nil {
		e.G.EventWriter.WriteEvent(ev)
	}
}

var ErrInvalidSelection = errors.New("invalid selection")
