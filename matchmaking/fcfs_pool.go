package matchmaking

import (
	"errors"
	"sync"

	"github.com/psurge1/distributed-chess-server/datastructures"
	"github.com/psurge1/distributed-chess-server/models"
)

type FirstComeFirstServeMPool struct {
	mu      sync.Mutex
	Players datastructures.Queue[*models.GPlayer]
}

func NewFirstComeFirstServeMPool() *FirstComeFirstServeMPool {
	mpool := FirstComeFirstServeMPool{}
	return &mpool
}

func (mpool *FirstComeFirstServeMPool) Join(gplayer *models.GPlayer) bool {
	mpool.mu.Lock()
	defer mpool.mu.Unlock()

	mpool.Players.Enqueue(gplayer)
	return true
}

func (mpool *FirstComeFirstServeMPool) Leave(gplayer *models.GPlayer) bool {
	mpool.mu.Lock()
	defer mpool.mu.Unlock()

	if _, err := mpool.Players.FindAndDelete(gplayer, func(one *models.GPlayer, two *models.GPlayer) bool {
		if one == nil && two == nil {
			return true
		}
		if one == nil || two == nil {
			return false
		}
		if one.Player == nil && two.Player == nil {
			return true
		}
		if one.Player == nil || two.Player == nil {
			return false
		}
		return one.Player.Username == two.Player.Username
	}); err != nil {
		return false
	}

	return true
}

func (mpool *FirstComeFirstServeMPool) ExtractMatch() (*models.Game, error) {
	mpool.mu.Lock()
	defer mpool.mu.Unlock()

	if mpool.Players.Len() < 2 {
		return nil, errors.New("not enough players in pool to match")
	}

	playerOne, err := mpool.Players.Dequeue()
	if err != nil {
		return nil, errors.New("unable to dequeue player from pool")
	}
	playerTwo, err := mpool.Players.Dequeue()
	if err != nil {
		mpool.Players.Enqueue(playerOne)
		return nil, errors.New("unable to dequeue player from pool")
	}

	return &models.Game{GWhite: playerOne, GBlack: playerTwo, Status: models.StatusPaired, Result: models.ResultUnknown}, nil
}
