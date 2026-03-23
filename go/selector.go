package thehunted

type Selector interface {
	SelectStart(g View) *SelectedStart
	SelectLoadout(g View) *SelectedLoadout
}
