package thehunted

type Selector interface {
	SelectStart(g GameView) *SelectedStart
	SelectLoadout(g GameView) *SelectedLoadout
	SelectDefensivePosture(g GameView) *SelectedDefensivePosture
}
