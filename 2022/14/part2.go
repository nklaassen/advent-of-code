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

	state := newState()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		state.pushPath(parsePath(scanner.Text()))
	}

	sandAtRest := 0
	for state.dropSand(false) {
		sandAtRest++
	}
	fmt.Println(state)
	fmt.Println(sandAtRest)
}

func parsePath(str string) path {
	var p path
	parts := strings.Split(str, " ")
	for i := 0; i < len(parts); i += 2 {
		coords := strings.Split(parts[i], ",")
		x, err := strconv.Atoi(coords[0])
		if err != nil {
			log.Fatal(err)
		}
		y, err := strconv.Atoi(coords[1])
		if err != nil {
			log.Fatal(err)
		}
		p = append(p, pair{x, y})
	}
	return p
}

type State struct {
	paths                  Paths
	blocked                map[pair]byte
	minX, maxX, minY, maxY int
}

func newState() *State {
	return &State{
		blocked: make(map[pair]byte),
		minX:    999,
	}
}

func (s *State) pushPath(p path) {
	s.paths = append(s.paths, p)
	for i := 1; i < len(p); i++ {
		prev := p[i-1]
		coords := p[i]
		if coords.x != prev.x {
			for x := min(coords.x, prev.x); x <= max(coords.x, prev.x); x++ {
				s.blocked[pair{x, coords.y}] = '#'
			}
		} else {
			for y := min(coords.y, prev.y); y <= max(coords.y, prev.y); y++ {
				s.blocked[pair{coords.x, y}] = '#'
			}
		}
		s.minX = min(s.minX, coords.x)
		s.maxX = max(s.maxX, coords.x)
		s.maxY = max(s.maxY, coords.y)
	}
}

func (s *State) dropSand(trace bool) bool {
	if _, ok := s.blocked[pair{500, 0}]; ok {
		return false
	}
	x := 500
	y := s.minY
	for y <= s.maxY {
		if trace {
			s.blocked[pair{x, y}] = '~'
		}
		if _, b := s.blocked[pair{x, y + 1}]; !b {
			y++
			continue
		}
		if _, b := s.blocked[pair{x - 1, y + 1}]; !b {
			x--
			y++
			continue
		}
		if _, b := s.blocked[pair{x + 1, y + 1}]; !b {
			x++
			y++
			continue
		}
		// all options blocked, stay put
		s.blocked[pair{x, y}] = 'o'
		return true
	}
	s.blocked[pair{x, y}] = 'o'
	return true
}

func (s *State) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d -> %d\n    ", s.minX, s.maxX)
	for x := s.minX; x <= s.maxX; x++ {
		fmt.Fprintf(&b, "%d", x%10)
	}
	fmt.Fprintln(&b)
	for y := s.minY; y <= s.maxY+1; y++ {
		fmt.Fprintf(&b, "%3d ", y)
		for x := s.minX; x <= s.maxX; x++ {
			if c, ok := s.blocked[pair{x, y}]; ok {
				fmt.Fprintf(&b, "%c", c)
			} else {
				fmt.Fprint(&b, ".")
			}
		}
		fmt.Fprintln(&b)
	}
	fmt.Fprint(&b, "    ")
	for x := s.minX; x <= s.maxX; x++ {
		fmt.Fprintf(&b, "%d", x%10)
	}
	fmt.Fprintf(&b, "\n%d -> %d\n", s.minX, s.maxX)
	return b.String()
}

type Paths []path

type path []pair

type pair struct {
	x, y int
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
