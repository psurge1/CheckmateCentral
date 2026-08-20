package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/psurge1/distributed-chess-server/models"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Home Page"))
	if err != nil {
		log.Println("Failed to write to response")
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// var mpool = matchmaking.NewFirstComeFirstServeMPool()
var (
	connections = make([]*models.GPlayer, 0)
	games       = make([]*models.Game, 0)
)

func joinGameHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Unable to initiate websocket upgrade from server"))
		return
	}

	gplayer := &models.GPlayer{Moves: conn}
	//mpool.Join(gplayer)
	connections = append(connections, gplayer)
	connectionCount := len(connections)
	if connectionCount > 0 && connectionCount%2 == 0 {
		// start a go routine with both connections
		// go establishChatConnection(connections[connectionCount-1], connections[connectionCount-2])
		// go establishChatConnection(connections[connectionCount-2], connections[connectionCount-1])
		whitePlayer := connections[connectionCount-1]
		blackPlayer := connections[connectionCount-2]
		game, err := models.NewGame(whitePlayer, blackPlayer)
		if err != nil {
			log.Println("failed to create game with players! players already in game?")
			return
		}
		games = append(games, game)
		go game.Run()
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/joingame", joinGameHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
