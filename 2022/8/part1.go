package main

import (
	"bytes"
	"fmt"
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

	var visible [size][size]bool

	for i := 0; i < size; i++ {
		visible[i][0] = true
		visible[0][i] = true
		visible[i][size-1] = true
		visible[size-1][i] = true
	}

	// right to left
	for row := 0; row < size; row++ {
		max := 0
		for col := 0; col < size; col++ {
			if h := heights[row][col]; h > max {
				visible[row][col] = true
				max = h
			}
		}
	}

	// left to right
	for row := 0; row < size; row++ {
		max := 0
		for col := size - 1; col >= 0; col-- {
			if h := heights[row][col]; h > max {
				visible[row][col] = true
				max = h
			}
		}
	}

	// top to bottom
	for col := 0; col < size; col++ {
		max := 0
		for row := 0; row < size; row++ {
			if h := heights[row][col]; h > max {
				visible[row][col] = true
				max = h
			}
		}
	}

	// bottom to top
	for col := 0; col < size; col++ {
		max := 0
		for row := size - 1; row >= 0; row-- {
			if h := heights[row][col]; h > max {
				visible[row][col] = true
				max = h
			}
		}
	}

	numVisible := 0
	for row := 0; row < size; row++ {
		for col := 0; col < size; col++ {
			if visible[row][col] {
				numVisible++
				print("1")
			} else {
				print("0")
			}
		}
		println()
	}

	fmt.Println(numVisible)
}
