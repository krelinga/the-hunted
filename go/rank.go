package thehunted

import (
	"errors"
	"fmt"
)

type Rank int

const (
	RankOltzS = iota + 1
	RankKptLt
	RankKKpt
	RankFKpt
	RankKptzS
)

var ErrInvalidRank = errors.New("invalid rank")

func (r Rank) Validate() error {
	switch r {
	case RankOltzS, RankKptLt, RankKKpt, RankFKpt, RankKptzS:
		return nil
	default:
		return ErrInvalidRank
	}
}

func (r Rank) String() string {
	switch r {
	case RankOltzS:
		return "Oberleutnant zur See (1)"
	case RankKptLt:
		return "Kapitänleutnant (2)"
	case RankKKpt:
		return "Korvettenkapitän (3)"
	case RankFKpt:
		return "Fregattenkapitän (4)"
	case RankKptzS:
		return "Kapitän zur See (5)"
	default:
		return fmt.Sprintf("Invalid  Rank (%d)", r)
	}
}