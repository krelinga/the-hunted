package thehunted

type Selector interface {
	SelectStart(g GameView) *SelectedStart
	SelectLoadout(g GameView) *SelectedLoadout
}
