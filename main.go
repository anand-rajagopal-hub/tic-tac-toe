package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

type TicTacToe struct {
	board [3][3]int
	//state GameState
}

func (t *TicTacToe) printBoard() {
	char := "-"
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if t.board[row][col] == 1 {
				char = "O"
			} else if t.board[row][col] == 2 {
				char = "X"
			} else {
				char = "-"
			}
			fmt.Printf("%2s\t", char)
		}
		fmt.Println("")
	}
	fmt.Println("-----------------------")
}

func (t *TicTacToe) hasW(input int) bool {
	return (t.board[0][0]&t.board[0][1]&t.board[0][2] == input) || // top row
		(t.board[0][0]&t.board[1][1]&t.board[2][2] == input) || // top left to bottom right diagonal
		(t.board[0][0]&t.board[1][0]&t.board[2][0] == input) || // first left column
		(t.board[2][0]&t.board[2][1]&t.board[2][2] == input) || // bottom row
		(t.board[0][2]&t.board[1][2]&t.board[2][2] == input) || // rightmost column
		(t.board[0][2]&t.board[1][1]&t.board[2][0] == input) || // right to left diagonal
		(t.board[0][1]&t.board[1][1]&t.board[2][1] == input) || // middle coloumn
		(t.board[1][0]&t.board[1][1]&t.board[1][2] == input) // middle row
}

func (t *TicTacToe) hasWinner() (bool, int) {
	if t.hasW(1) {
		return true, 1
	}
	if t.hasW(2) {
		return true, 2
	}
	return false, -1
}

func noMovesPossible(t *TicTacToe) bool {
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if t.board[row][col] == 0 {
				return false
			}
		}
	}
	return true
}

func minimax(t *TicTacToe, player, r, c int) (int, int, int) {
	//t.printBoard()
	win, who := t.hasWinner()
	if win && who == 1 {
		return -1, r, c
	} else if win && who == 2 {
		return 1, r, c
	} else if noMovesPossible(t) {
		return 0, r, c
	}
	bestScore := -1*(1<<8) - 1
	if player == 1 {
		bestScore = 1 << 8
	}
	i := 0
	j := 0
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if t.board[row][col] == 0 {
				t.board[row][col] = player
				p := 2
				if player == 2 {
					p = 1
				}
				score, _, _ := minimax(t, p, row, col)
				t.board[row][col] = 0
				if player == 2 {
					if score > bestScore {
						bestScore = score
						i = row
						j = col

					}
				} else {
					if score <= bestScore {
						bestScore = score
						i = row
						j = col
					}
				}
			}
		}
	}
	return bestScore, i, j
}

func (t *TicTacToe) nextMove() int {
	score, r, c := minimax(t, 2, 0, 0)
	fmt.Println("score", score, r, c)

	t.board[r][c] = 2
	return score

}

func NewGame() *TicTacToe {
	return &TicTacToe{}
}

func main() {
	t := NewGame()
	t.printBoard()
	for i := 0; i < 10; i++ {
		fmt.Println("Enter Row")

		reader := bufio.NewReader(os.Stdin)
		b, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatalf("%v", err)
		}
		var row uint8
		err = binary.Read(bytes.NewReader(b), binary.BigEndian, &row)
		if err != nil {
			log.Fatalf("%v", err)
		}

		fmt.Println("Enter Col")

		b, err = reader.ReadBytes('\n')
		if err != nil {
			log.Fatalf("%v", err)
		}
		var col uint8
		err = binary.Read(bytes.NewReader(b), binary.BigEndian, &col)
		if err != nil {
			log.Fatalf("%v", err)
		}
		t.board[row-48-1][col-48-1] = 1
		t.printBoard()
		score := t.nextMove()
		t.printBoard()
		fmt.Println(score)
		w, who := t.hasWinner()
		if w {
			fmt.Printf("Game finished. %d won\n", who)
			break
		}
		if noMovesPossible(t) {
			fmt.Println("Game ended in a tie!")
			break
		}
	}
}
