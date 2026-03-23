package thehunted

import (
	"errors"
	"fmt"

	"github.com/krelinga/the-hunted/go/views"
)

type View interface {
	GetKmdtName() string
	GetKmdtRank() Rank
	GetCrewQuality() CrewQuality
	GetUBoat() UBoatView
	GetPatrols() views2.Slice[PatrolView]
}

type Data struct {
	KmdtName    string
	KmdtRank    Rank
	CrewQuality CrewQuality
	UBoat       *UBoatData
	Patrols     []*PatrolData
	// TODO: rename to NextPatrolDate.
	StartPatrolDate PatrolDate
}

func (d *Data) View() View {
	return viewImpl{data: d}
}

type viewImpl struct {
	data *Data
}

func (v viewImpl) GetKmdtName() string {
	return v.data.KmdtName
}

func (v viewImpl) GetKmdtRank() Rank {
	return v.data.KmdtRank
}

func (v viewImpl) GetCrewQuality() CrewQuality {
	return v.data.CrewQuality
}

func (v viewImpl) GetUBoat() UBoatView {
	return v.data.UBoat.View()
}

func (v viewImpl) GetPatrols() views2.Slice[PatrolView] {
	return views2.WrapViewerSlice(v.data.Patrols)
}

func (v viewImpl) GetStartPatrolDate() PatrolDate {
	return v.data.StartPatrolDate
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
	return g.data.View()
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
	newState, err := allHandlers[g.nextState](g.GetView(), g.Selector, roller, ew)
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

type handler func(g View, s Selector, r Roller, ew EventWriter) (gameState, error)

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
