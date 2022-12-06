package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type game struct {
	opponent, outcome             string
	opponentOffset, outcomeOffset int
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var games []game
	for scanner.Scan() {
		l := scanner.Text()
		parts := strings.Split(l, " ")
		games = append(games, game{
			opponent:       parts[0],
			outcome:        parts[1],
			opponentOffset: int(parts[0][0]) - 'A',
			outcomeOffset:  int(parts[1][0]) - 'X',
		})
	}

	score := 0

	moves := []string{"rock", "paper", "scissors"}
	outcomes := []string{"lose", "draw", "win"}

	for _, g := range games {
		fmt.Printf("opp: %s outcome: %s", moves[g.opponentOffset], outcomes[g.outcomeOffset])

		outcomePoints := 3 * g.outcomeOffset
		fmt.Printf(" outcomePoints: %d", outcomePoints)
		score += outcomePoints

		// if lose, rotate back in moves, draw stays same, win rotates forward
		offset := g.outcomeOffset - 1

		move := (g.opponentOffset + offset) % 3
		if move < 0 {
			move += 3
		}

		fmt.Printf(" move: %s", moves[move])

		movePoints := move + 1
		fmt.Printf(" movePoints: %d", movePoints)
		score += movePoints

		fmt.Printf(" score: %d\n", score)
	}

	fmt.Println(score)
}
