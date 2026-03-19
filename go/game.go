package thehunted

import (
	"errors"
	"fmt"
)

type GameView interface {
	GetKmdtName() string
	GetKmdtRank() Rank
	GetCrewQuality() CrewQuality
	GetUBoat() UBoatView
	GetPatrols() PatrolsView
}

type GameData struct {
	KmdtName    string
	KmdtRank    Rank
	CrewQuality CrewQuality
	UBoat       *UBoatData
	Patrols     PatrolsData
	// TODO: rename to NextPatrolDate.
	StartPatrolDate PatrolDate
}

func (g *GameData) GetKmdtName() string {
	return g.KmdtName
}

func (g *GameData) GetKmdtRank() Rank {
	return g.KmdtRank
}

func (g *GameData) GetCrewQuality() CrewQuality {
	return g.CrewQuality
}

func (g *GameData) GetUBoat() UBoatView {
	return g.UBoat
}

func (g *GameData) GetPatrols() PatrolsView {
	return g.Patrols
}

func (g *GameData) GetStartPatrolDate() PatrolDate {
	return g.StartPatrolDate
}

type Game interface {
	GameView

	SetSelector(selector Selector) Game
	SetEventWriter(eventWriter EventWriter) Game
	SetRoller(roller Roller) Game

	Form() Form
	Advance(form Form) error
	Next() error
	Done() bool
}

type gameImpl struct {
	GameData

	selector Selector
	eventWriter EventWriter
	roller Roller

	gameState GameState
	nextState gameState
}

func (g *gameImpl) SetSelector(selector Selector) Game {
	g.selector = selector
	return g
}

func (g *gameImpl) SetEventWriter(eventWriter EventWriter) Game {
	g.eventWriter = eventWriter
	return g
}

func (g *gameImpl) SetRoller(roller Roller) Game {
	g.roller = roller
	return g
}

func (g *gameImpl) setGameState(gameState GameState) {
	g.gameState = gameState
	g.eventWriter.WriteEvent(GameStateSetEvent{GameState: g.gameState})
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

func (g *gameImpl) Form() Form {
	switch g.gameState {
	case GameStateNotStarted:
		return g.formForNotStarted()
	case GameStateSelectLoadout:
		return g.formForSelectLoadout()
	default:
		panic(fmt.Sprintf("unexpected game state %v", g.gameState))
	}
}

func (g *gameImpl) Advance(form Form) error {
	switch g.gameState {
	case GameStateNotStarted:
		return g.advanceFromNotStarted(form)
	case GameStateSelectLoadout:
		return g.advanceFromSelectLoadout(form)
	default:
		panic(fmt.Sprintf("unexpected game state %v", g.gameState))
	}
}

func (g *gameImpl) Next() error {
	if g.Done() {
		panic("game is already done")
	}
	newState, err := allHandlers[g.nextState](g)
	if err != nil {
		return err
	}
	g.nextState = newState
	return nil
}

func (g *gameImpl) Done() bool {
	return g.gameState == GameStateFinished
}

type GameStateSetEvent struct {
	baseEvent
	GameState GameState
}

func (e GameStateSetEvent) String() string {
	return fmt.Sprintf("Game state set: %s", e.GameState)
}

func NewGame() Game {
	gi := &gameImpl{}
	return gi.SetRoller(RandomRoller{})
}

type gameState int

const (
	gameStateStart gameState = iota
	gameStateSelectLoadout
	gameStateStartPatrol
	gameStateDone
)

type handler func(g *gameImpl) (gameState, error)

var allHandlers = map[gameState]handler{
	gameStateStart:         handleStart,
	gameStateSelectLoadout: handleSelectLoadout,
	gameStateStartPatrol:   handleStartPatrol,
}
