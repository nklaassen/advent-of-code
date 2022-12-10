package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var interestingCycles []int = []int{20, 60, 100, 140, 180, 220}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	x := 1
	cycle := 1

	interestingCycleMap := make(map[int]struct{})
	for _, c := range interestingCycles {
		interestingCycleMap[c] = struct{}{}
	}

	sum := 0

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
		//fmt.Println(cmd, dx)
		for i := 0; i < dCycle; i++ {
			if _, ok := interestingCycleMap[cycle]; ok {
				fmt.Println(cycle, x)
				fmt.Println(cycle * x)
				sum += cycle * x
			}
			cycle++
		}
		x += dx
	}
	fmt.Println(sum)
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
