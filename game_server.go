package main

import (
	"log"
	"net/http"

	"github.com/psurge1/distributed-chess-server/internal/httpapi"
)

func main() {
	root := http.NewServeMux()

	userHandlers := &httpapi.UsersHandlers{}
	sessionHandlers := &httpapi.SessionsHandlers{}
	tournamentHandlers := &httpapi.TournamentsHandlers{}
	gameHandlers := httpapi.NewGamesHandlers()

	users := httpapi.NewUsersMux(userHandlers)
	sessions := httpapi.NewSessionsMux(sessionHandlers)
	tournaments := httpapi.NewTournamentsMux(tournamentHandlers)
	games := httpapi.NewGamesMux(gameHandlers)

	root.Handle("/users", users)
	root.Handle("/users/", users)
	root.Handle("/sessions", sessions)
	root.Handle("/sessions/", sessions)
	root.Handle("/tournaments", tournaments)
	root.Handle("/tournaments/", tournaments)
	root.Handle("/games", games)
	root.Handle("/games/", games)

	root.HandleFunc("/{$}", httpapi.Home)
	root.HandleFunc("/joingame", gameHandlers.JoinLegacyGame)

	log.Fatal(http.ListenAndServe(":8080", root))
}
