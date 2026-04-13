package thehunted

import (
	"errors"
)

type GameView interface {
	GetKmdtName() string
	GetKmdtRank() Rank
	GetCrewQuality() CrewQuality
	GetUBoat() UBoatView
	GetPatrols() PatrolsView
	GetStartPatrolDate() PatrolDate
}

type Game struct {
	KmdtName    string
	KmdtRank    Rank
	CrewQuality CrewQuality
	UBoat       *UBoat
	Patrols     Patrols
	// TODO: rename to NextPatrolDate.
	StartPatrolDate PatrolDate
}

func (v *Game) GetKmdtName() string {
	return v.KmdtName
}

func (v *Game) GetKmdtRank() Rank {
	return v.KmdtRank
}

func (v *Game) GetCrewQuality() CrewQuality {
	return v.CrewQuality
}

func (v *Game) GetUBoat() UBoatView {
	return v.UBoat
}

func (v *Game) GetPatrols() PatrolsView {
	return v.Patrols
}

func (v *Game) GetStartPatrolDate() PatrolDate {
	return v.StartPatrolDate
}

type Engine struct {
	Selector    Selector
	EventWriter EventWriter
	Roller      Roller

	game      Game
	nextState gameState
}

var errNoChange = errors.New("no change in game state")

func (g *Engine) Next() error {
	if g.Done() {
		panic("game is already done")
	}
	roller := g.Roller
	if roller == nil {
		roller = RandomRoller{}
	}
	ew := applyEventToGame{G: g}
	newState, err := allHandlers[g.nextState](&g.game, g.Selector, roller, ew)
	if err == errNoChange {
		return nil
	} else if err != nil {
		return err
	}
	g.nextState = newState
	return nil
}

func (g *Engine) Done() bool {
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

type handler func(g GameView, s Selector, r Roller, ew EventWriter) (gameState, error)

var allHandlers = map[gameState]handler{
	gameStateStart:         handleStart,
	gameStateSelectLoadout: handleSelectLoadout,
}

type applyEventToGame struct {
	G *Engine
}

func (e applyEventToGame) WriteEvent(ev Event) {
	ev.apply(&e.G.game)
	if e.G.EventWriter != nil {
		e.G.EventWriter.WriteEvent(ev)
	}
}

var ErrInvalidSelection = errors.New("invalid selection")
