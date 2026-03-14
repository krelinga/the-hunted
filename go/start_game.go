package thehunted

import "fmt"

type StartGameForm struct {
	baseForm
	UBoatType SelectFormField[UBoatType]
	UBoatID TextFormField
	KmdtName TextFormField
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
	UBoatType UBoatType
}

func (e FirstPatrolDateSetEvent) String() string {
	return fmt.Sprintf("First patrol date set: %s (based on u-boat type %s)", e.FirstPatrolDate, e.UBoatType)
}

type StartingRankSetEvent struct {
	baseEvent
	D6 D6
	Rank Rank
	PatrolDate PatrolDate
}

func (e StartingRankSetEvent) String() string {
	return fmt.Sprintf("Starting rank set: %s (based on d6 roll %s and patrol date %s)", e.Rank, e.D6, e.PatrolDate)
}

func (g *Game) advanceFromNotStarted(form Form) ([]Event, error) {
	startGameForm, ok := form.(*StartGameForm)
	if !ok {
		return nil, fmt.Errorf("%w: expected StartGameForm, got %T", ErrUnexpectedForm, form)
	}
	if err := startGameForm.Validate(); err != nil {
		return nil, err
	}
	var events []Event
	g.kmdtName = string(startGameForm.KmdtName)
	events = append(events, KmdtNamedEvent{KmdtName: g.kmdtName})
	g.startingUBoatType = startGameForm.UBoatType.Options[startGameForm.UBoatType.Selected]
	events = append(events, UBoatTypeSelectedEvent{UBoatType: g.startingUBoatType})
	if firstPatrolDate, err := g.startingUBoatType.FirstPatrolDate(); err != nil {
		panic(fmt.Sprintf("invalid starting u-boat type: %v", err))
	} else {
		g.startPatrolDate = firstPatrolDate
		events = append(events, FirstPatrolDateSetEvent{FirstPatrolDate: g.startPatrolDate, UBoatType: g.startingUBoatType})
	}
	rankD6 := g.Roller.RollD6()
	var rankThreshold D6
	switch g.startPatrolDate {
	case PatrolDateJul43, PatrolDateAug43, PatrolDateSep43, PatrolDateOct43, PatrolDateNov43, PatrolDateDec43:
		rankThreshold = 4
	default:
		rankThreshold = 5
	}
	if rankD6 <= rankThreshold {
		g.kmdtRank = RankOltzS
	} else {
		g.kmdtRank = RankKptLt
	}
	events = append(events, StartingRankSetEvent{D6: rankD6, Rank: g.kmdtRank, PatrolDate: g.startPatrolDate})
	g.gameState = GameStateInPort
	events = append(events, GameStateSetEvent{GameState: g.gameState})
	return events, nil
}