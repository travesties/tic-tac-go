package game

import (
	"fmt"
	"net/http"
)

type Game struct {
	GameOver bool       `json:"gameOver"`
	XIsNext  bool       `json:"xIsNext"`
	Board    [][]string `json:"board"`
}

type PlayerMove struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type InvalidMove struct {
	Status int
	Msg    string
}

func (mr *InvalidMove) Error() string {
	return mr.Msg
}

func NewGame() *Game {
	game := Game{
		GameOver: false,
		XIsNext:  true,
		Board:    [][]string{{"", "", ""}, {"", "", ""}, {"", "", ""}},
	}
	return &game
}

func ApplyMove(move *PlayerMove, game *Game) error {
	if game.GameOver {
		msg := fmt.Sprintf("invalid move (x:%v, y:%v) game is over. DELETE /game to reset.", move.X, move.Y)
		return &InvalidMove{Status: http.StatusBadRequest, Msg: msg}
	}
	if move.X < 0 || move.X > 2 || move.Y < 0 || move.Y > 2 {
		msg := fmt.Sprintf("invalid move (x:%v, y:%v) 0 <= x,y <= 2", move.X, move.Y)
		return &InvalidMove{Status: http.StatusBadRequest, Msg: msg}
	}
	if game.Board[move.X][move.Y] != "" {
		msg := fmt.Sprintf("invalid move (x:%v, y:%v) already occupied", move.X, move.Y)
		return &InvalidMove{Status: http.StatusBadRequest, Msg: msg}
	}

	var playerMark string
	if game.XIsNext {
		playerMark = "X"
	} else {
		playerMark = "O"
	}

	game.Board[move.X][move.Y] = playerMark
	return nil
}

func PlayerWon(player string, move *PlayerMove, game *Game) bool {
	// Check rows and colums for this move
	winOnRow := true
	for rc := 0; rc < 3; rc++ {
		if game.Board[move.X][rc] != player || game.Board[rc][move.Y] != player {
			winOnRow = false
			break
		}
	}

	// Check diagonal if this move is on one
	winOnCol := true
	for x, y := 0, 0; x < 3; x, y = x+1, y+1 {
		if game.Board[x][y] != player {
			winOnCol = false
			break
		}
	}

	return winOnRow || winOnCol
}
