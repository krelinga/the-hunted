package thehunted

import "errors"

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