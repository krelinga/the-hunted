package thehunted

type Driver interface {
	StartGame(g GameView) *StartGameResult
	SelectLoadout(g GameView) *SelectLoadoutResult
	SelectDefensivePosture(g GameView) *SelectDefensivePostureResult
}