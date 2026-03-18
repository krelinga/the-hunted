package thehunted

import "fmt"

type StartGameForm struct {
	baseForm
	UBoatType SelectFormField[UBoatType]
	UBoatID   TextFormField
	KmdtName  TextFormField
}

type SelectedStart struct {
	UBoatType UBoatType
	UBoatID   string
	KmdtName  string
}

func (f *StartGameForm) Validate() error {
	if err := f.UBoatType.Validate(); err != nil {
		return fmt.Errorf("%w: invalid u-boat type", err)
	}
	return nil
}

func (g *gameImpl) formForNotStarted() Form {
	return &StartGameForm{
		UBoatType: SelectFormField[UBoatType]{
			Options: []UBoatType{
				UBoatTypeVIIB, UBoatTypeVIIC, UBoatTypeVIICFlak, UBoatTypeVIIC41, UBoatTypeVIID,
				UBoatTypeIXB, UBoatTypeIXC, UBoatTypeIXC40, UBoatTypeIXD2, UBoatTypeIXD42,
				UBoatTypeXB, UBoatTypeXII, UBoatTypeXIV, UBoatTypeXXI,
			},
		},
	}
}

type KmdtNamedEvent struct {
	baseEvent
	KmdtName string
}

func (e KmdtNamedEvent) String() string {
	return fmt.Sprintf("Kommandant named: %s", e.KmdtName)
}

type UBoatTypeSelectedEvent struct {
	baseEvent
	UBoatType UBoatType
}

func (e UBoatTypeSelectedEvent) String() string {
	return fmt.Sprintf("U-boat type selected: %s", e.UBoatType)
}

type FirstPatrolDateSetEvent struct {
	baseEvent
	FirstPatrolDate PatrolDate
	UBoatType       UBoatType
}

func (e FirstPatrolDateSetEvent) String() string {
	return fmt.Sprintf("First patrol date set: %s (based on u-boat type %s)", e.FirstPatrolDate, e.UBoatType)
}

type StartingRankSetEvent struct {
	baseEvent
	D6         DiceD6
	Rank       Rank
	PatrolDate PatrolDate
}

func (e StartingRankSetEvent) String() string {
	return fmt.Sprintf("Starting rank set: %s (based on d6 roll %s and patrol date %s)", e.Rank, e.D6, e.PatrolDate)
}

type CrewQualitySetEvent struct {
	baseEvent
	CrewQuality CrewQuality
}

func (e CrewQualitySetEvent) String() string {
	return fmt.Sprintf("Crew quality set: %s", e.CrewQuality)
}

func (g *gameImpl) advanceFromNotStarted(form Form) error {
	startGameForm, ok := form.(*StartGameForm)
	if !ok {
		return fmt.Errorf("%w: expected *StartGameForm, got %T", ErrUnexpectedForm, form)
	}
	if err := startGameForm.Validate(); err != nil {
		return err
	}
	g.KmdtName = string(startGameForm.KmdtName)
	g.writeEvent(KmdtNamedEvent{KmdtName: g.KmdtName})
	uboatType := startGameForm.UBoatType.Options[startGameForm.UBoatType.Selected]
	g.UBoat = NewUBoatData(uboatType, string(startGameForm.UBoatID))
	g.writeEvent(UBoatTypeSelectedEvent{UBoatType: uboatType})
	g.StartPatrolDate = uboatType.FirstPatrolDate()
	g.writeEvent(FirstPatrolDateSetEvent{FirstPatrolDate: g.StartPatrolDate, UBoatType: uboatType})
	rankD6 := g.Options.Roller.RollD6()
	var rankThreshold DiceD6
	if g.StartPatrolDate.Year() <= 1943 {
		rankThreshold = 4
	} else {
		rankThreshold = 5
	}
	if rankD6 <= rankThreshold {
		g.KmdtRank = RankOltzS
	} else {
		g.KmdtRank = RankKptLt
	}
	g.writeEvent(StartingRankSetEvent{D6: rankD6, Rank: g.KmdtRank, PatrolDate: g.StartPatrolDate})
	g.CrewQuality = CrewQualityTrained
	g.writeEvent(CrewQualitySetEvent{CrewQuality: g.CrewQuality})
	g.setGameState(GameStateSelectLoadout)
	return nil
}

func handleStart(g *gameImpl) (gameState, error) {
	return gameStateDone, nil // TODO
}