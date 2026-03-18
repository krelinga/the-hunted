package thehunted

type Driver interface {
	StartGame(g *Game) *StartGameResult
	SelectLoadout(g *Game) *SelectLoadoutResult
	SelectDefensivePosture(g *Game) *SelectDefensivePostureResult
}