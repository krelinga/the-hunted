package thehunted

import (
	"errors"
	"fmt"
)

type Game struct {
	Roller      Roller
	EventWriter EventWriter

	UBoat   UBoat
	Patrols []Patrol

	kmdtName        string
	kmdtRank        Rank
	crewQuality     CrewQuality
	gameState       GameState
	startPatrolDate PatrolDate
}

func (g *Game) writeEvent(event Event) {
	if g.EventWriter != nil {
		g.EventWriter.WriteEvent(event)
	}
}

func (g *Game) setGameState(gameState GameState) {
	g.gameState = gameState
	g.writeEvent(GameStateSetEvent{GameState: g.gameState})
}

func (g *Game) KmdtName() string {
	return g.kmdtName
}

func (g *Game) KmdtRank() Rank {
	return g.kmdtRank
}

func (g *Game) CrewQuality() CrewQuality {
	return g.crewQuality
}

func (g *Game) GameState() GameState {
	return g.gameState
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

func (g *Game) IsFinished() bool {
	return g.gameState == GameStateFinished
}

type GameStateSetEvent struct {
	baseEvent
	GameState GameState
}

func (e GameStateSetEvent) String() string {
	return fmt.Sprintf("Game state set: %s", e.GameState)
}
