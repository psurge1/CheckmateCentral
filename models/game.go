package models

import (
	"errors"

	"github.com/psurge1/distributed-chess-server/streaming"
)

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

func NewGame(white *GPlayer, black *GPlayer) (*Game, error) {
	if white.Game != nil {
		return nil, errors.New("white player already belongs to a game")
	}
	if black.Game != nil {
		return nil, errors.New("black player already belongs to a game")
	}
	return &Game{GWhite: white, GBlack: black, Status: StatusPaired, Result: ResultUnknown}, nil
}

func (game *Game) Run() {
	streaming.EstablishGameConnection(game.GWhite.Moves, game.GBlack.Moves)
}
