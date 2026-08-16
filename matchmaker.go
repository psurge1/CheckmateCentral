package main

import (
	"errors"
	"slices"
	"sync"

	"github.com/psurge1/distributed-chess-server/datastructures"
)

type MatchmakingPool interface {
	Join(gplayer *GPlayer) bool
	Leave(gplayer *GPlayer) bool
	ExtractMatch() *Game
}

type FirstComeFirstServeMPool struct {
	mu      sync.Mutex
	Players datastructures.Queue[*GPlayer]
}

func (mpool *FirstComeFirstServeMPool) Join(gplayer *GPlayer) bool {
	mpool.mu.Lock()
	defer mpool.mu.Unlock()

	mpool.Players.Enqueue(gplayer)
	return true
}

func (mpool *FirstComeFirstServeMPool) Leave(gplayer *GPlayer) bool {
	mpool.mu.Lock()
	defer mpool.mu.Unlock()

	targetIdx := -1
	for idx, player := range mpool.Players {
		if player.Player.Username == gplayer.Player.Username {
			targetIdx = idx
			break
		}
	}
	if targetIdx != -1 {
		mpool.Players = slices.Delete(mpool.Players, targetIdx, targetIdx+1)
	}
	return true
}

func (mpool *FirstComeFirstServeMPool) ExtractMatch() (*Game, error) {
	mpool.mu.Lock()
	defer mpool.mu.Unlock()

	playerOne, err := mpool.Players.Dequeue()
	if err != nil {
		return nil, errors.New("unable to dequeue player from pool")
	}
	playerTwo, err := mpool.Players.Dequeue()
	if err != nil {
		return nil, errors.New("unable to dequeue player from pool")
	}

	return &Game{GWhite: playerOne, GBlack: playerTwo, Status: StatusPaired, Result: ResultUnknown}, nil
}
