package models

import (
	"errors"

	"github.com/notnil/chess"
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
	if white == nil || black == nil {
		return nil, errors.New("both players are required")
	}
	if white == black {
		return nil, errors.New("a player cannot play against themselves")
	}
	if white.Game != nil {
		return nil, errors.New("white player already belongs to a game")
	}
	if black.Game != nil {
		return nil, errors.New("black player already belongs to a game")
	}
	game := &Game{GWhite: white, GBlack: black, Status: StatusPaired, Result: ResultUnknown}
	white.Game = game
	black.Game = game
	return game, nil
}

func (game *Game) Run() {
	game.Status = StatusPlaying
	outcome := streaming.EstablishGameConnection(game.GWhite.Moves, game.GBlack.Moves)
	switch outcome {
	case chess.WhiteWon:
		game.Result = ResultWhiteWon
	case chess.BlackWon:
		game.Result = ResultBlackWon
	case chess.Draw:
		game.Result = ResultDraw
	default:
		game.Result = ResultUnknown
	}
	game.Status = StatusFinished
	game.GWhite.Game = nil
	game.GBlack.Game = nil
}
