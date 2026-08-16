package main

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
