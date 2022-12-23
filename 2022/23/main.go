package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var input = flag.String("input", "input.txt", "input file name")

func main() {
	flag.Parse()

	f, _ := os.ReadFile(*input)
	lines := strings.Split(string(f), "\n")

	s := newState()
	for i, line := range lines {
		for j, c := range line {
			if c == '#' {
				s.pushElf(Pair{j, i})
			}
		}
	}
	for i := 0; i < 10; i++ {
		s.round()
	}
	fmt.Println(s)
	fmt.Println("Part 1:", s.countEmpty())
	round := 11
	for s.round() {
		fmt.Println("movement in round", round)
		round++
	}
	fmt.Println("Part 2:", round)
}

type State struct {
	elves                  map[Pair]bool
	directions             []Move
	minX, maxX, minY, maxY int
}

func newState() *State {
	return &State{
		elves:      make(map[Pair]bool),
		directions: []Move{N, S, W, E},
	}
}

func (s *State) round() bool {
	// map of proposed destinations to pairs of elves proposing to move there
	proposals := make(map[Pair][]Pair)
	for elf := range s.elves {
		hasNeighbour := false
		var proposal Pair
	directionLoop:
		for _, d := range s.directions {
			for _, c := range d.check {
				if s.elves[add(elf, c)] {
					hasNeighbour = true
					continue directionLoop
				}
			}
			if proposal != (Pair{}) {
				continue
			}
			proposal = add(elf, d.dest)
		}
		if hasNeighbour && proposal != (Pair{}) {
			proposals[proposal] = append(proposals[proposal], elf)
		}
	}

	moved := false
	for dest, elves := range proposals {
		if len(elves) != 1 {
			continue
		}
		moved = true
		elf := elves[0]
		s.elves[dest] = true
		delete(s.elves, elf)
		s.minX = min(s.minX, dest.x)
		s.maxX = max(s.maxX, dest.x)
		s.minY = min(s.minY, dest.y)
		s.maxY = max(s.maxY, dest.y)
	}
	s.rotateDirections()
	return moved
}

func (s *State) countEmpty() int {
	width := s.maxX - s.minX + 1
	height := s.maxY - s.minY + 1
	totalArea := width * height
	return totalArea - len(s.elves)
}

func (s *State) pushElf(elf Pair) {
	s.elves[elf] = true
	s.minX = min(s.minX, elf.x)
	s.maxX = max(s.maxX, elf.x)
	s.minY = min(s.minY, elf.y)
	s.maxY = max(s.maxY, elf.y)
}

func (s *State) rotateDirections() {
	s.directions = append(s.directions[1:], s.directions[0])
}

func (s *State) String() string {
	var sb strings.Builder
	for y := s.minY; y <= s.maxY; y++ {
		for x := s.minX; x <= s.maxX; x++ {
			if x == 0 && y == 0 {
				sb.Write([]byte("X"))
				continue
			}
			if s.elves[Pair{x, y}] {
				sb.Write([]byte("#"))
			} else {
				sb.Write([]byte(" "))
			}
		}
		sb.Write([]byte("\n"))
	}
	return sb.String()
}

type Move struct {
	check []Pair
	dest  Pair
}

func newMove(cardinal Pair) Move {
	return Move{
		check: []Pair{
			add(cardinal, Pair{-cardinal.y, -cardinal.x}),
			cardinal,
			add(cardinal, Pair{cardinal.y, cardinal.x}),
		},
		dest: cardinal,
	}
}

var (
	N = newMove(Pair{0, -1})
	S = newMove(Pair{0, 1})
	W = newMove(Pair{-1, 0})
	E = newMove(Pair{1, 0})
)

type Pair struct {
	x, y int
}

func add(a, b Pair) Pair {
	return Pair{
		a.x + b.x,
		a.y + b.y,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
