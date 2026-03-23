package thehunted

import (
	"errors"
	"fmt"
)

type View interface {
	GetKmdtName() string
	GetKmdtRank() Rank
	GetCrewQuality() CrewQuality
	GetUBoat() UBoatView
	GetPatrols() PatrolsView
}

type Data struct {
	KmdtName    string
	KmdtRank    Rank
	CrewQuality CrewQuality
	UBoat       *UBoatData
	Patrols     PatrolsData
	// TODO: rename to NextPatrolDate.
	StartPatrolDate PatrolDate
}

func (g *Data) GetKmdtName() string {
	return g.KmdtName
}

func (g *Data) GetKmdtRank() Rank {
	return g.KmdtRank
}

func (g *Data) GetCrewQuality() CrewQuality {
	return g.CrewQuality
}

func (g *Data) GetUBoat() UBoatView {
	return g.UBoat
}

func (g *Data) GetPatrols() PatrolsView {
	return g.Patrols
}

func (g *Data) GetStartPatrolDate() PatrolDate {
	return g.StartPatrolDate
}

type Game struct {
	Selector    Selector
	EventWriter EventWriter
	Roller      Roller

	data      Data
	gameState GameState
	nextState gameState
}

func (g *Game) GetView() View {
	return &g.data
}

func (g *Game) Form() Form {
	switch g.gameState {
	case GameStateNotStarted:
		return g.formForNotStarted()
	case GameStateSelectLoadout:
		return g.formForSelectLoadout()
	default:
		panic(fmt.Sprintf("unexpected game state %v", g.gameState))
	}
}

func (g *Game) Advance(form Form) error {
	switch g.gameState {
	case GameStateNotStarted:
		return g.advanceFromNotStarted(form)
	case GameStateSelectLoadout:
		return g.advanceFromSelectLoadout(form)
	default:
		panic(fmt.Sprintf("unexpected game state %v", g.gameState))
	}
}

func (g *Game) Next() error {
	if g.Done() {
		panic("game is already done")
	}
	roller := g.Roller
	if roller == nil {
		roller = RandomRoller{}
	}
	ew := applyEventToGame{G: g}
	newState, err := allHandlers[g.nextState](g.GetView(), roller, ew)
	if err != nil {
		return err
	}
	g.nextState = newState
	return nil
}

func (g *Game) Done() bool {
	return g.gameState == GameStateFinished
}

func (g *Game) setGameState(newState GameState) {
	g.gameState = newState
	g.EventWriter.WriteEvent(GameStateSetEvent{
		GameState: newState,
	})
}

type GameState int

const (
	GameStateNotStarted GameState = iota
	GameStateSelectLoadout
	GameStateFinished
)

func (gs GameState) String() string {
	switch gs {
	case GameStateNotStarted:
		return "Not Started"
	case GameStateSelectLoadout:
		return "Select Loadout"
	case GameStateFinished:
		return "Finished"
	default:
		return fmt.Sprintf("Unknown GameState (%d)", gs)
	}
}

var ErrUnexpectedForm = errors.New("unexpected form")

type GameStateSetEvent struct {
	baseEvent
	GameState GameState
}

func (e GameStateSetEvent) apply(gd *Data) {
	// No GameData fields are affected by a game state change, so this is a no-op.
	// TODO: remove this event type eventually.
}

func (e GameStateSetEvent) String() string {
	return fmt.Sprintf("Game state set: %s", e.GameState)
}

type gameState int

const (
	gameStateStart gameState = iota
	gameStateSelectLoadout
	gameStateStartPatrol
	gameStateDone
)

type handler func(g View, r Roller, ew EventWriter) (gameState, error)

var allHandlers = map[gameState]handler{
	gameStateStart:         handleStart,
	gameStateSelectLoadout: handleSelectLoadout,
	gameStateStartPatrol:   handleStartPatrol,
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
