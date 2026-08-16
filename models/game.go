package models

type GameStatus int

const (
	StatusPaired GameStatus = iota
	StatusPlaying
	StatusFinished
	StatusUnknown
)

type GameResult int

const (
	ResultWhiteWon GameResult = iota
	ResultBlackWon
	ResultDraw
	ResultUnknown
)

type Game struct {
	GWhite *GPlayer
	GBlack *GPlayer
	Status GameStatus
	Result GameResult
	// TournamentID TournamentID
}

func NewGame(white *GPlayer, black *GPlayer) *Game {
	return &Game{GWhite: white, GBlack: black, Status: StatusPaired, Result: ResultUnknown}
}

func (game *Game) Run() {
}
