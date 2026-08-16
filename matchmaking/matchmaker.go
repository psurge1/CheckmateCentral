// Package matchmaking is used to pair players in the general game pool or in each tournament pool
package matchmaking

import (
	"github.com/psurge1/distributed-chess-server/models"
	ts "github.com/psurge1/distributed-chess-server/taskscheduler"
)

type MatchmakingPool interface {
	Join(gplayer *models.GPlayer) bool
	Leave(gplayer *models.GPlayer) bool
	ExtractMatch() (*models.Game, error)
}

type Matchmaker struct {
	mpool         MatchmakingPool   // owned by the matchmaker (only one pool per matchmaker and vise versa)
	gameScheduler *ts.TaskScheduler // not owned (many matchmakers can use the same game scheduler)
}

func (mm *Matchmaker) tryMatch() {
	//	game, err := mm.mpool.ExtractMatch()
	//	if err != nil {
	//		return
	//	}

	// TODO handle implementation of gameScheduler.Submit
	//	if err := mm.gameScheduler.Submit(game); err != nil {
	//		// TODO: handle scheduling failure
	//	}
}
