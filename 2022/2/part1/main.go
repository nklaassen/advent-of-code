package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type game struct {
	opponent, player string
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
			parts[0], parts[1],
		})
	}

	score := 0

	numTie := 0
	numWin := 0
	numLose := 0
	for _, g := range games {
		score += int(g.player[0]) - int('W')
		if int(g.opponent[0])-int('A') == int(g.player[0])-int('X') {
			score += 3
			numTie++
		}
		if g.player == "X" && g.opponent == "C" {
			// rock beats scissors
			score += 6
			numWin++
		} else if g.player == "Y" && g.opponent == "A" {
			// paper beats rock
			score += 6
			numWin++
		} else if g.player == "Z" && g.opponent == "B" {
			// scissors beats paper
			score += 6
			numWin++
		} else {
			numLose++
		}
	}

	fmt.Println(score)
}
