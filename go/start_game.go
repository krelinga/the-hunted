package thehunted

import "fmt"

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
		return fmt.Errorf("%w: u-boat ID cannot be empty", ErrInvalidSelection)
	}
	if ss.KmdtName == "" {
		return fmt.Errorf("%w: kommandant name cannot be empty", ErrInvalidSelection)
	}
	return nil
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
	gd.UBoat = NewUBoat(e.UBoatType, e.UBoatID)
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

func handleStart(g View, s Selector, r Roller, ew EventWriter) (gameState, error) {
	selected := s.SelectStart(g)
	if selected == nil {
		return 0, errNoChange
	}
	if err := selected.Validate(); err != nil {
		return 0, err
	}

	ew.WriteEvent(KmdtNamedEvent{
		KmdtName: selected.KmdtName,
	})

	ew.WriteEvent(NewUBoatEvent{
		UBoatType: selected.UBoatType,
		UBoatID:   selected.UBoatID,
	})

	patrolDate := selected.UBoatType.FirstPatrolDate()
	ew.WriteEvent(FirstPatrolDateSetEvent{
		FirstPatrolDate: patrolDate,
		UBoatType:       selected.UBoatType,
	})

	rankD6 := r.RollD6()
	var rankThreshold DiceD6
	if patrolDate.Year() <= 1943 {
		rankThreshold = 4
	} else {
		rankThreshold = 5
	}
	var rank Rank
	if rankD6 <= rankThreshold {
		rank = RankOltzS
	} else {
		rank = RankKptLt
	}
	ew.WriteEvent(StartingRankSetEvent{
		D6:         rankD6,
		Rank:       rank,
		PatrolDate: patrolDate,
	})

	ew.WriteEvent(CrewQualitySetEvent{
		CrewQuality: CrewQualityTrained,
	})

	return gameStateSelectLoadout, nil
}
