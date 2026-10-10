package matchmaking

import (
	"testing"

	"github.com/psurge1/distributed-chess-server/models"
)

func TestExtractMatchReservesPlayers(t *testing.T) {
	pool := NewFirstComeFirstServeMPool()
	first, second := &models.GPlayer{}, &models.GPlayer{}
	pool.Join(first)
	pool.Join(second)
	game, err := pool.ExtractMatch()
	if err != nil {
		t.Fatal(err)
	}
	if game.GWhite != first || game.GBlack != second || first.Game != game || second.Game != game {
		t.Fatal("FIFO pairing did not reserve both players")
	}
	if pool.Players.Len() != 0 {
		t.Fatal("paired players remain in the pool")
	}
}

func TestExtractMatchRestoresPlayersOnFailure(t *testing.T) {
	pool := NewFirstComeFirstServeMPool()
	busy, opponent, waiting := &models.GPlayer{}, &models.GPlayer{}, &models.GPlayer{}
	if _, err := models.NewGame(busy, opponent); err != nil {
		t.Fatal(err)
	}
	pool.Join(busy)
	pool.Join(waiting)
	if game, err := pool.ExtractMatch(); err == nil || game != nil {
		t.Fatal("pool paired a player already in a game")
	}
	if pool.Players.Len() != 2 || waiting.Game != nil {
		t.Fatal("failed pairing lost players or reserved the waiting player")
	}
}
