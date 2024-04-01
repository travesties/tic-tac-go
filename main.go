package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"travesties/tictactoe/game"
)

var currentGame *game.Game

func handleGetGame(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(currentGame)
}

func handleResetGame(w http.ResponseWriter, _ *http.Request) {
	currentGame = game.NewGame()
	json.NewEncoder(w).Encode(currentGame)
}

func handlePlayerMove(w http.ResponseWriter, r *http.Request) {
	var move game.PlayerMove

	// Verify correct JSON body
	err := decodeJSONBody(w, r, &move)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Print(err.Error())
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		return
	}

	// Attempt to apply the posted move
	err = game.ApplyMove(&move, currentGame)
	if err != nil {
		var im *game.InvalidMove
		if errors.As(err, &im) {
			http.Error(w, im.Msg, im.Status)
		} else {
			log.Print(err.Error())
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		return
	}

	var player string
	if currentGame.XIsNext {
		player = "X"
	} else {
		player = "O"
	}

	win := game.PlayerWon(player, &move, currentGame)
	if win {
		currentGame.GameOver = true
		fmt.Fprintf(w, "Player %v wins! DELETE /game to reset", player)
	} else {
		currentGame.XIsNext = !currentGame.XIsNext
		fmt.Fprintf(w, "Move (x:%v, y:%v) applied", move.X, move.Y)
	}
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("GET /game", handleGetGame)
	router.HandleFunc("POST /game", handlePlayerMove)
	router.HandleFunc("DELETE /game", handleResetGame)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	currentGame = game.NewGame()

	fmt.Println("Server listening on port :8080")
	server.ListenAndServe()
}
