package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)

	s := newState()

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		dir := parts[0]
		count, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Fatal(err)
		}
		switch dir {
		case "L":
			s.apply(-1, 0, count)
		case "R":
			s.apply(1, 0, count)
		case "U":
			s.apply(0, 1, count)
		case "D":
			s.apply(0, -1, count)
		}
	}

	fmt.Println(len(s.visited))
}

type pair struct {
	x, y int
}

type state struct {
	head, tail pair
	visited    map[pair]struct{}
}

func newState() *state {
	return &state{
		visited: map[pair]struct{}{
			pair{0, 0}: struct{}{},
		},
	}
}

func (s *state) apply(dx, dy, count int) {
	for i := 0; i < count; i++ {
		s.head.x += dx
		s.head.y += dy
		if s.distance() < 2 {
			continue
		}
		s.moveToward()
		s.visited[s.tail] = struct{}{}
	}
}

func (s *state) moveToward() {
	if s.head.x != s.tail.x {
		if s.tail.x > s.head.x {
			s.tail.x--
		} else {
			s.tail.x++
		}
	}
	if s.head.y != s.tail.y {
		if s.tail.y > s.head.y {
			s.tail.y--
		} else {
			s.tail.y++
		}
	}
}

func (s *state) distance() int {
	return max(abs(s.head.x-s.tail.x), abs(s.head.y-s.tail.y))
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
