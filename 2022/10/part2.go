package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	screenWidth  = 40
	screenHeight = 6
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	x := 1
	cycle := 1

	var x_t []int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		cmd := parts[0]
		var dCycle int
		var dx int
		switch cmd {
		case "noop":
			dCycle = 1
			dx = 0
		case "addx":
			dCycle = 2
			dx, err = strconv.Atoi(parts[1])
			if err != nil {
				log.Fatal(err)
			}
		}
		for i := 0; i < dCycle; i++ {
			rayX := (cycle - 1) % 40
			if abs(rayX-x) < 2 {
				print("#")
			} else {
				print(".")
			}
			if cycle%40 == 0 {
				println()
			}
			cycle++
			x_t = append(x_t, x)
		}
		x += dx
	}
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x, y int) int {
	return -max(-x, -y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
