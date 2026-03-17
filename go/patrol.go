package thehunted

import (
	"slices"
	"errors"
	"fmt"
)

type PatrolDate int

const (
	PatrolDateJul43 = iota + 1
	PatrolDateAug43
	PatrolDateSep43
	PatrolDateOct43
	PatrolDateNov43
	PatrolDateDec43
	PatrolDateJan44
	PatrolDateFeb44
	PatrolDateMar44
	PatrolDateApr44
	PatrolDateMay44
	PatrolDateJun44
	PatrolDateJul44
	PatrolDateAug44
	PatrolDateSep44
	PatrolDateOct44
	PatrolDateNov44
	PatrolDateDec44
	PatrolDateJan45
	PatrolDateFeb45
	PatrolDateMar45
	PatrolDateApr45
	PatrolDateMay45
)

var ErrInvalidPatrolDate = errors.New("invalid patrol date")

func (p PatrolDate) Validate() error {
	if p < PatrolDateJul43 || p > PatrolDateMay45 {
		return ErrInvalidPatrolDate
	}
	return nil
}

func (p PatrolDate) Must() {
	if err := p.Validate(); err != nil {
		panic(err)
	}
}

func (p PatrolDate) Year() int {
	p.Must()
	return 1943 + int((p-7)/12)
}

func (p PatrolDate) String() string {
	switch p {
	case PatrolDateJul43:
		return "July 1943"
	case PatrolDateAug43:
		return "August 1943"
	case PatrolDateSep43:
		return "September 1943"
	case PatrolDateOct43:
		return "October 1943"
	case PatrolDateNov43:
		return "November 1943"
	case PatrolDateDec43:
		return "December 1943"
	case PatrolDateJan44:
		return "January 1944"
	case PatrolDateFeb44:
		return "February 1944"
	case PatrolDateMar44:
		return "March 1944"
	case PatrolDateApr44:
		return "April 1944"
	case PatrolDateMay44:
		return "May 1944"
	case PatrolDateJun44:
		return "June 1944"
	case PatrolDateJul44:
		return "July 1944"
	case PatrolDateAug44:
		return "August 1944"
	case PatrolDateSep44:
		return "September 1944"
	case PatrolDateOct44:
		return "October 1944"
	case PatrolDateNov44:
		return "November 1944"
	case PatrolDateDec44:
		return "December 1944"
	case PatrolDateJan45:
		return "January 1945"
	case PatrolDateFeb45:
		return "February 1945"
	case PatrolDateMar45:
		return "March 1945"
	case PatrolDateApr45:
		return "April 1945"
	case PatrolDateMay45:
		return "May 1945"
	default:
		return "Invalid Patrol Date"
	}
}

type PatrolState struct {
}

type PatrolSpot int

const (
	PatrolSpotAtlantic = iota + 1
	PatrolSpotIndianOcean
	PatrolSpotBritishIsles
	PatrolSpotNorthAmerica
	PatrolSpotMediterranean
	PatrolSpotBrazilianCoast
	PatrolSpotWestAfricanCoast
	PatrolSpotInvasion
	PatrolSpotArctic
	PatrolSpotCaribbean
)

func (p PatrolSpot) String() string {
	switch p {
	case PatrolSpotAtlantic:
		return "Atlantic"
	case PatrolSpotIndianOcean:
		return "Indian Ocean"
	case PatrolSpotBritishIsles:
		return "British Isles"
	case PatrolSpotNorthAmerica:
		return "North America"
	case PatrolSpotMediterranean:
		return "Mediterranean"
	case PatrolSpotBrazilianCoast:
		return "Brazilian Coast"
	case PatrolSpotWestAfricanCoast:
		return "West African Coast"
	case PatrolSpotInvasion:
		return "Invasion"
	case PatrolSpotArctic:
		return "Arctic"
	case PatrolSpotCaribbean:
		return "Caribbean"
	default:
		return fmt.Sprintf("Invalid Patrol Spot (%d)", p)
	}
}

var ErrInvalidPatrolSpot = errors.New("invalid patrol spot")

func (p PatrolSpot) Validate() error {
	switch p {
	case PatrolSpotAtlantic, PatrolSpotIndianOcean, PatrolSpotBritishIsles, PatrolSpotNorthAmerica, PatrolSpotMediterranean, PatrolSpotBrazilianCoast, PatrolSpotWestAfricanCoast, PatrolSpotInvasion, PatrolSpotArctic, PatrolSpotCaribbean:
		return nil
	default:
		return fmt.Errorf("%w: %d", ErrInvalidPatrolSpot, p)
	}
}

func (p PatrolSpot) Must() {
	if err := p.Validate(); err != nil {
		panic(err)
	}
}

func (p PatrolSpot) IsAnyOf(spots ...PatrolSpot) bool {
	return slices.Contains(spots, p)
}

type PatrolSpotAssignmentEvent struct {
	baseEvent
	PatrolSpot PatrolSpot
	Result2D6  Result2D6
	UBoatType  UBoatType
	PatrolDate PatrolDate
}

func (e PatrolSpotAssignmentEvent) String() string {
	return fmt.Sprintf("Patrol spot assigned: %s (rolled %s for %s on %s)", e.PatrolSpot, e.Result2D6, e.UBoatType, e.PatrolDate)
}

type PatrolAssignment struct {
	PatrolSpot  PatrolSpot
	Wolfpack    bool
	AbwehrAgent bool
}

func (g *Game) determinePatrolAssignment() PatrolAssignment {
	var wolfpack, abwehrAgent bool
	result := g.Roller.Roll2D6()
	var spot PatrolSpot
	switch {
	case g.startPatrolDate <= PatrolDateDec43:
		switch result.AsInt() {
		case 2, 5:
			spot = PatrolSpotIndianOcean
		case 3, 6:
			spot = PatrolSpotAtlantic
			wolfpack = true
		case 4:
			spot = PatrolSpotBritishIsles
		case 7, 8:
			spot = PatrolSpotAtlantic
		case 9:
			spot = PatrolSpotNorthAmerica
		case 10:
			spot = PatrolSpotMediterranean
		case 11:
			spot = PatrolSpotBrazilianCoast
		case 12:
			spot = PatrolSpotWestAfricanCoast
		}
	default:
		panic("patrol assignment for dates after December 1943 not implemented yet")
	}

	switch {
		case g.UBoat.UBoatType.IsTypeIX() && spot.IsAnyOf(PatrolSpotArctic, PatrolSpotMediterranean):
			spot = PatrolSpotWestAfricanCoast
		case g.UBoat.UBoatType.IsTypeVII() && spot.IsAnyOf(PatrolSpotWestAfricanCoast, PatrolSpotBrazilianCoast, PatrolSpotIndianOcean):
			spot = PatrolSpotAtlantic
		case g.UBoat.UBoatType.IsTypeVII() && g.UBoat.UBoatType != UBoatTypeVIID && spot == PatrolSpotCaribbean:
			spot = PatrolSpotAtlantic
		case g.UBoat.UBoatType.IsAnyOf(UBoatTypeIXD2, UBoatTypeIXD42) && spot == PatrolSpotAtlantic:
			spot = PatrolSpotIndianOcean
	}

	g.writeEvent(PatrolSpotAssignmentEvent{
		PatrolSpot: spot,
		Result2D6:  result,
		UBoatType:  g.UBoat.UBoatType,
		PatrolDate: g.startPatrolDate,
	})
	return PatrolAssignment{
		PatrolSpot:  spot,
		Wolfpack:    wolfpack,
		AbwehrAgent: abwehrAgent,
	}
}
