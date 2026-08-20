// Package streaming is used for websocket utils
package streaming

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/notnil/chess"
)

func EstablishGameConnection(whitePlayer *websocket.Conn, blackPlayer *websocket.Conn) chess.Outcome {
	defer whitePlayer.Close()
	defer blackPlayer.Close()

	err := whitePlayer.WriteMessage(1, []byte("White"))
	if err != nil {
		log.Println(err)
		return chess.NoOutcome
	}
	err = blackPlayer.WriteMessage(1, []byte("Black"))
	if err != nil {
		log.Println(err)
		return chess.NoOutcome
	}
	game := chess.NewGame(chess.UseNotation(chess.UCINotation{}))
	for game.Outcome() == chess.NoOutcome {
		messageType, p, err := whitePlayer.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		// perform game move validation
		err = game.MoveStr(string(p))
		for err != nil {
			messageType, p, err = whitePlayer.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
			err = game.MoveStr(string(p))
		}
		if err := blackPlayer.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			break
		}

		log.Println(game.Position().Board().Draw())

		messageType, p, err = blackPlayer.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		// perform game move validation
		err = game.MoveStr(string(p))
		for err != nil {
			messageType, p, err = blackPlayer.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
			err = game.MoveStr(string(p))
		}
		if err := whitePlayer.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			break
		}

		log.Println(game.Position().Board().Draw())
	}
	whitePlayer.WriteMessage(1, []byte(game.Outcome().String()))
	blackPlayer.WriteMessage(1, []byte(game.Outcome().String()))
	return game.Outcome()
}

func EstablishChatConnection(sender *websocket.Conn, receiver *websocket.Conn) {
	defer sender.Close()
	defer receiver.Close()
	for {
		messageType, p, err := sender.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		if err := receiver.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			break
		}
		log.Println(messageType, string(p))
	}
}
