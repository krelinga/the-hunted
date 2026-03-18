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

type GameOptions struct {
	Roller      Roller
	EventWriter EventWriter
	Driver      Selector
}

type Game interface {
	GameView

	Form() Form
	Advance(form Form) error
	IsFinished() bool
}

type gameImpl struct {
	GameData
	Options   GameOptions
	gameState GameState
}

func (g *gameImpl) writeEvent(event Event) {
	if g.Options.EventWriter != nil {
		g.Options.EventWriter.WriteEvent(event)
	}
}

func (g *gameImpl) setGameState(gameState GameState) {
	g.gameState = gameState
	g.writeEvent(GameStateSetEvent{GameState: g.gameState})
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

func (g *gameImpl) IsFinished() bool {
	return g.gameState == GameStateFinished
}

type GameStateSetEvent struct {
	baseEvent
	GameState GameState
}

func (e GameStateSetEvent) String() string {
	return fmt.Sprintf("Game state set: %s", e.GameState)
}

func NewGame(options GameOptions) Game {
	return &gameImpl{
		Options: options,
	}
}
