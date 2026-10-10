package httpapi

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/psurge1/distributed-chess-server/models"
)

type GamesHandlers struct {
	upgrader    websocket.Upgrader
	connections []*models.GPlayer
	games       []*models.Game
}

func NewGamesHandlers() *GamesHandlers {
	return &GamesHandlers{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (h *GamesHandlers) List(w http.ResponseWriter, r *http.Request) {
	// TODO: List games filtered by tournament or player.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *GamesHandlers) Get(w http.ResponseWriter, r *http.Request) {
	// TODO: Get the game identified by r.PathValue("gameID").
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *GamesHandlers) ListMoves(w http.ResponseWriter, r *http.Request) {
	// TODO: Get the game's move history.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *GamesHandlers) SubmitMove(w http.ResponseWriter, r *http.Request) {
	// TODO: Submit a move if HTTP move submission is enabled.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *GamesHandlers) Resign(w http.ResponseWriter, r *http.Request) {
	// TODO: Resign the game as the signed-in player.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *GamesHandlers) Stream(w http.ResponseWriter, r *http.Request) {
	// TODO: Upgrade to a WebSocket for an existing game.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// JoinLegacyGame preserves the existing /joingame WebSocket demo.
func (h *GamesHandlers) JoinLegacyGame(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Unable to initiate websocket upgrade from server"))
		return
	}

	gplayer := &models.GPlayer{Moves: conn}
	h.connections = append(h.connections, gplayer)
	connectionCount := len(h.connections)
	if connectionCount > 0 && connectionCount%2 == 0 {
		whitePlayer := h.connections[connectionCount-1]
		blackPlayer := h.connections[connectionCount-2]
		game, err := models.NewGame(whitePlayer, blackPlayer)
		if err != nil {
			log.Println("failed to create game with players! players already in game?")
			return
		}
		h.games = append(h.games, game)
		go game.Run()
	}
}
