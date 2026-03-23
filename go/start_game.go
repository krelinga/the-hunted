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

func (ss SelectedStart) Validate() error {
	if err := ss.UBoatType.Validate(); err != nil {
		return fmt.Errorf("%w: invalid u-boat type", err)
	}
	if ss.UBoatID == "" {
		return fmt.Errorf("%w: u-boat ID cannot be empty", ErrInvalidFormField)
	}
	if ss.KmdtName == "" {
		return fmt.Errorf("%w: kommandant name cannot be empty", ErrInvalidFormField)
	}
	return nil
}

func (f *StartGameForm) Validate() error {
	if err := f.UBoatType.Validate(); err != nil {
		return fmt.Errorf("%w: invalid u-boat type", err)
	}
	return nil
}

func (g *Game) formForNotStarted() Form {
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

func (e KmdtNamedEvent) apply(gd *Data) {
	gd.KmdtName = e.KmdtName
}

func (e KmdtNamedEvent) String() string {
	return fmt.Sprintf("Kommandant named: %s", e.KmdtName)
}

type NewUBoatEvent struct {
	baseEvent
	UBoatType UBoatType
	UBoatID   string
}

func (e NewUBoatEvent) apply(gd *Data) {
	gd.UBoat = NewUBoatData(e.UBoatType, e.UBoatID)
}

func (e NewUBoatEvent) String() string {
	return fmt.Sprintf("U-boat type selected: %s", e.UBoatType)
}

type FirstPatrolDateSetEvent struct {
	baseEvent
	FirstPatrolDate PatrolDate
	UBoatType       UBoatType
}

func (e FirstPatrolDateSetEvent) apply(gd *Data) {
	gd.StartPatrolDate = e.FirstPatrolDate
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

func (e StartingRankSetEvent) apply(gd *Data) {
	gd.KmdtRank = e.Rank
}

func (e StartingRankSetEvent) String() string {
	return fmt.Sprintf("Starting rank set: %s (based on d6 roll %s and patrol date %s)", e.Rank, e.D6, e.PatrolDate)
}

type CrewQualitySetEvent struct {
	baseEvent
	CrewQuality CrewQuality
}

func (e CrewQualitySetEvent) apply(gd *Data) {
	gd.CrewQuality = e.CrewQuality
}

func (e CrewQualitySetEvent) String() string {
	return fmt.Sprintf("Crew quality set: %s", e.CrewQuality)
}

func (g *Game) advanceFromNotStarted(form Form) error {
	startGameForm, ok := form.(*StartGameForm)
	if !ok {
		return fmt.Errorf("%w: expected *StartGameForm, got %T", ErrUnexpectedForm, form)
	}
	if err := startGameForm.Validate(); err != nil {
		return err
	}
	g.data.KmdtName = string(startGameForm.KmdtName)
	g.EventWriter.WriteEvent(KmdtNamedEvent{KmdtName: g.data.KmdtName})
	uboatType := startGameForm.UBoatType.Options[startGameForm.UBoatType.Selected]
	g.data.UBoat = NewUBoatData(uboatType, string(startGameForm.UBoatID))
	g.EventWriter.WriteEvent(NewUBoatEvent{UBoatType: uboatType, UBoatID: string(startGameForm.UBoatID)})
	g.data.StartPatrolDate = uboatType.FirstPatrolDate()
	g.EventWriter.WriteEvent(FirstPatrolDateSetEvent{FirstPatrolDate: g.data.StartPatrolDate, UBoatType: uboatType})
	rankD6 := g.Roller.RollD6()
	var rankThreshold DiceD6
	if g.data.StartPatrolDate.Year() <= 1943 {
		rankThreshold = 4
	} else {
		rankThreshold = 5
	}
	if rankD6 <= rankThreshold {
		g.data.KmdtRank = RankOltzS
	} else {
		g.data.KmdtRank = RankKptLt
	}
	g.EventWriter.WriteEvent(StartingRankSetEvent{
		D6:         rankD6,
		Rank:       g.data.KmdtRank,
		PatrolDate: g.data.StartPatrolDate,
	})
	g.data.CrewQuality = CrewQualityTrained
	g.EventWriter.WriteEvent(CrewQualitySetEvent{CrewQuality: g.data.CrewQuality})
	g.setGameState(GameStateSelectLoadout)
	return nil
}

func handleStart(g View, r Roller, ew EventWriter) (gameState, error) {
	return gameStateDone, nil // TODO
}
