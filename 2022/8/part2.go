package main

import (
	"bytes"
	"log"
	"os"
)

const (
	size = 99
)

func main() {
	file, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	var heights [size][size]int

	rows := bytes.Split(file, []byte{'\n'})
	for row := 0; row < size; row++ {
		for col := 0; col < size; col++ {
			heights[row][col] = int(rows[row][col] - '0')
		}
	}

	max := 0
	for row := 1; row < size-1; row++ {
		for col := 1; col < size-1; col++ {
			s := scenicScore(heights, row, col)
			if s > max {
				max = s
			}
		}
	}

	println(max)
}

func scenicScore(heights [size][size]int, row, col int) int {
	h := heights[row][col]
	score := 1
	for r := row - 1; r >= 0; r-- {
		if heights[r][col] >= h || r == 0 {
			score *= row - r
			break
		}
	}
	for r := row + 1; r < size; r++ {
		if heights[r][col] >= h || r == size-1 {
			score *= r - row
			break
		}
	}
	for c := col - 1; c >= 0; c-- {
		if heights[row][c] >= h || c == 0 {
			score *= col - c
			break
		}
	}
	for c := col + 1; c < size; c++ {
		if heights[row][c] >= h || c == size-1 {
			score *= c - col
			break
		}
	}
	return score
}
