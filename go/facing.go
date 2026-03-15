package thehunted

type Facing bool

const (
	FacingFwd Facing = false
	FacingAft Facing = true
)

func (f Facing) String() string {
	if f == FacingFwd {
		return "Forward"
	}
	return "Aft"
}