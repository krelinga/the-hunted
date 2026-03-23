package thehunted

import (
	"errors"

	"github.com/krelinga/the-hunted/go/views"
)

type View interface {
	GetKmdtName() string
	GetKmdtRank() Rank
	GetCrewQuality() CrewQuality
	GetUBoat() UBoatView
	GetPatrols() views.Slice[PatrolView]
	GetStartPatrolDate() PatrolDate
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

func (v viewImpl) GetPatrols() views.Slice[PatrolView] {
	return views.WrapViewerSlice(v.data.Patrols)
}

func (v viewImpl) GetStartPatrolDate() PatrolDate {
	return v.data.StartPatrolDate
}

type Game struct {
	Selector    Selector
	EventWriter EventWriter
	Roller      Roller

	data      Data
	nextState gameState
}

func (g *Game) GetView() View {
	return g.data.View()
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
	newState, err := allHandlers[g.nextState](g.GetView(), g.Selector, roller, ew)
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
