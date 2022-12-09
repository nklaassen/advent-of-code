package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	numKnots = 10
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
	knots   [10]pair
	visited map[pair]struct{}
}

func newState() *state {
	return &state{
		visited: map[pair]struct{}{
			pair{0, 0}: struct{}{},
		},
	}
}

func (s *state) apply(dx, dy, count int) {
	for c := 0; c < count; c++ {
		s.knots[0].x += dx
		s.knots[0].y += dy
		for i := 1; i < numKnots; i++ {
			follow(&s.knots[i-1], &s.knots[i])
		}
		s.visited[s.knots[numKnots-1]] = struct{}{}
	}
}

func follow(leader, follower *pair) {
	if distance(leader, follower) < 2 {
		return
	}
	moveToward(leader, follower)
}

func moveToward(leader, follower *pair) {
	if leader.x != follower.x {
		if follower.x > leader.x {
			follower.x--
		} else {
			follower.x++
		}
	}
	if leader.y != follower.y {
		if follower.y > leader.y {
			follower.y--
		} else {
			follower.y++
		}
	}
}

func distance(leader, follower *pair) int {
	return max(abs(leader.x-follower.x), abs(leader.y-follower.y))
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
