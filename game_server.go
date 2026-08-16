package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/notnil/chess"
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

func establishGameConnection(whitePlayer *websocket.Conn, blackPlayer *websocket.Conn) {
	defer whitePlayer.Close()
	defer blackPlayer.Close()

	whitePlayer.WriteMessage(1, []byte("White"))
	blackPlayer.WriteMessage(1, []byte("Black"))
	game := chess.NewGame(chess.UseNotation(chess.UCINotation{}))
	for game.Outcome() == chess.NoOutcome {
		messageType, p, err := whitePlayer.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// perform game move validation
		err = game.MoveStr(string(p))
		for err != nil {
			messageType, p, err = whitePlayer.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			err = game.MoveStr(string(p))
		}
		if err := blackPlayer.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

		log.Println(game.Position().Board().Draw())

		messageType, p, err = blackPlayer.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// perform game move validation
		err = game.MoveStr(string(p))
		for err != nil {
			messageType, p, err = blackPlayer.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			err = game.MoveStr(string(p))
		}
		if err := whitePlayer.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

		log.Println(game.Position().Board().Draw())
	}
	whitePlayer.WriteMessage(1, []byte(game.Outcome().String()))
	blackPlayer.WriteMessage(1, []byte(game.Outcome().String()))
}

func establishChatConnection(sender *websocket.Conn, receiver *websocket.Conn) {
	defer sender.Close()
	defer receiver.Close()
	for {
		messageType, p, err := sender.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := receiver.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
		log.Println(messageType, string(p))
	}
}

func joinGameHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Unable to initiate websocket upgrade from server"))
		return
	}

	connections = append(connections, conn)
	connectionCount := len(connections)
	if connectionCount > 0 && connectionCount%2 == 0 {
		// start a go routine with both connections
		// go establishChatConnection(connections[connectionCount-1], connections[connectionCount-2])
		// go establishChatConnection(connections[connectionCount-2], connections[connectionCount-1])
		go establishGameConnection(connections[connectionCount-1], connections[connectionCount-2])
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/joingame", joinGameHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
