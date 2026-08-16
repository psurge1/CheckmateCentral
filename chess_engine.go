package main

import (
	"bufio"
	"log"
	"os"

	"github.com/notnil/chess"
)

func getUserInput(scanner *bufio.Scanner, message string) string {
	log.Print(message)
	line := ""
	if scanner.Scan() {
		line = scanner.Text()
	}
	return line
}

func playGame() {
	scanner := bufio.NewScanner(os.Stdin)
	game := chess.NewGame(chess.UseNotation(chess.UCINotation{}))
	for game.Outcome() == chess.NoOutcome {
		err := game.MoveStr(getUserInput(scanner, "Move: "))
		for err != nil {
			log.Println("Invalid Move!")
			err = game.MoveStr(getUserInput(scanner, "Move: "))
		}
		log.Println(game.Position().Board().Draw())
	}
	log.Println(game.Position().Board().Draw())
}
