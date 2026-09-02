// Package matchmaking is used to pair players in the general game pool or in each tournament pool
package matchmaking

import (
	"github.com/psurge1/distributed-chess-server/models"
	ts "github.com/psurge1/distributed-chess-server/taskscheduler"
)

type MatchmakingPool interface {
	Join(gplayer *models.GPlayer) bool
	Leave(gplayer *models.GPlayer) bool
	ExtractMatches() ([]*models.Game, error)
}

type Matchmaker struct {
	mpool         MatchmakingPool   // owned by the matchmaker (only one pool per matchmaker and vise versa)
	gameScheduler *ts.TaskScheduler // not owned (many matchmakers can use the same game scheduler)
}

func (mm *Matchmaker) TryMatch() {
	games, err := mm.mpool.ExtractMatches()
	if err != nil {
		return
	}

	// convert game objects to jobs
	numGames := len(games)
	jobSlice := make([]ts.Job, numGames)
	for i, game := range games {
		jobSlice[i] = game
	}
	batch := &ts.Batch{Jobs: jobSlice, Count: numGames}

	// submit games as jobs to game scheduler
	retryCount := 5
	for retry, err := 0, mm.gameScheduler.Submit(batch); err != nil && retry < retryCount; retry++ {
		err = mm.gameScheduler.Submit(batch)
	}
}
