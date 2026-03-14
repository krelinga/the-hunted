package thehunted

import "errors"

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

func (p PatrolDate) Year() int {
	if err := p.Validate(); err != nil {
		return 0
	}
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