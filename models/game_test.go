package models

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestNewGameReservesPlayers(t *testing.T) {
	white, black, other := &GPlayer{}, &GPlayer{}, &GPlayer{}
	game, err := NewGame(white, black)
	if err != nil {
		t.Fatal(err)
	}
	if game.GWhite != white || game.GBlack != black || white.Game != game || black.Game != game {
		t.Fatal("game and player references do not agree")
	}
	if game.Status != StatusPaired || game.Result != ResultUnknown {
		t.Fatal("new game must be paired with an unknown result")
	}
	if _, err := NewGame(white, other); err == nil {
		t.Fatal("player already paired as white was accepted")
	}
	if _, err := NewGame(other, black); err == nil {
		t.Fatal("player already paired as black was accepted")
	}
	if other.Game != nil || white.Game != game || black.Game != game {
		t.Fatal("failed pairing changed game ownership")
	}
}

func TestNewGameRejectsInvalidPlayers(t *testing.T) {
	player := &GPlayer{}
	for _, players := range [][2]*GPlayer{{nil, player}, {player, nil}, {nil, nil}, {player, player}} {
		if game, err := NewGame(players[0], players[1]); err == nil || game != nil {
			t.Fatal("invalid players were accepted")
		}
	}
	if player.Game != nil {
		t.Fatal("invalid pairing reserved the player")
	}
}

func TestGameRunRecordsOutcomeAndReleasesPlayers(t *testing.T) {
	drawMoves := make([]string, 0, 16)
	for range 4 {
		drawMoves = append(drawMoves, "g1f3", "g8f6", "f3g1", "f6g8")
	}
	cases := []struct {
		name    string
		moves   []string
		outcome string
		result  GameResult
	}{
		{"white checkmate", []string{"e2e4", "e7e5", "f1c4", "b8c6", "d1h5", "g8f6", "h5f7"}, "1-0", ResultWhiteWon},
		{"black checkmate", []string{"f2f3", "e7e5", "g2g4", "d8h4"}, "0-1", ResultBlackWon},
		{"fivefold repetition", drawMoves, "1/2-1/2", ResultDraw},
		{"disconnect without outcome", nil, "*", ResultUnknown},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			connections := make(chan *websocket.Conn, 2)
			upgrader := websocket.Upgrader{}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				conn, err := upgrader.Upgrade(w, r, nil)
				if err != nil {
					t.Error(err)
					return
				}
				connections <- conn
			}))
			t.Cleanup(server.Close)
			dialer := websocket.Dialer{HandshakeTimeout: 3 * time.Second}
			url := "ws" + server.URL[len("http"):]
			whiteClient, _, err := dialer.Dial(url, nil)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { whiteClient.Close() })
			white := &GPlayer{Moves: <-connections}
			t.Cleanup(func() { white.Moves.Close() })
			blackClient, _, err := dialer.Dial(url, nil)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { blackClient.Close() })
			black := &GPlayer{Moves: <-connections}
			t.Cleanup(func() { black.Moves.Close() })
			game, err := NewGame(white, black)
			if err != nil {
				t.Fatal(err)
			}
			done := make(chan struct{})
			go func() {
				game.Run()
				close(done)
			}()
			readGameMessage(t, whiteClient, "White")
			readGameMessage(t, blackClient, "Black")
			for i, move := range tc.moves {
				current, opponent := whiteClient, blackClient
				if i%2 == 1 {
					current, opponent = blackClient, whiteClient
				}
				if err := current.WriteMessage(websocket.TextMessage, []byte(move)); err != nil {
					t.Fatal(err)
				}
				readGameMessage(t, opponent, move)
			}
			if len(tc.moves) == 0 {
				whiteClient.Close()
				readGameMessage(t, blackClient, tc.outcome)
			} else {
				readGameMessage(t, whiteClient, tc.outcome)
				readGameMessage(t, blackClient, tc.outcome)
			}
			select {
			case <-done:
			case <-time.After(3 * time.Second):
				t.Fatal("game did not finish")
			}
			if game.Status != StatusFinished || game.Result != tc.result {
				t.Fatalf("status/result = %v/%v, want %v/%v", game.Status, game.Result, StatusFinished, tc.result)
			}
			if white.Game != nil || black.Game != nil {
				t.Fatal("finished game did not release players")
			}
			if _, err := NewGame(white, black); err != nil {
				t.Fatalf("players cannot be paired again: %v", err)
			}
		})
	}
}

func readGameMessage(t *testing.T, conn *websocket.Conn, want string) {
	t.Helper()
	if err := conn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		t.Fatal(err)
	}
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	if string(message) != want {
		t.Fatalf("message = %q, want %q", message, want)
	}
}
