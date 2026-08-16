package main

import "github.com/gorilla/websocket"

type Player struct {
	Username string
	Elo      int
}

type GPlayer struct {
	Player *Player
	Moves  *websocket.Conn
	Game   *Game
}
