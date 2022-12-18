package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	deltas = [...]Coord{
		Coord{1, 0, 0},
		Coord{-1, 0, 0},
		Coord{0, 1, 0},
		Coord{0, -1, 0},
		Coord{0, 0, 1},
		Coord{0, 0, -1},
	}
)

type Coord struct {
	x, y, z int
}

func (c *Coord) add(o *Coord) Coord {
	return Coord{
		c.x + o.x,
		c.y + o.y,
		c.z + o.z,
	}
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

func minParts(a, b *Coord) Coord {
	return Coord{
		min(a.x, b.x),
		min(a.y, b.y),
		min(a.z, b.z),
	}
}

func maxParts(a, b *Coord) Coord {
	return Coord{
		max(a.x, b.x),
		max(a.y, b.y),
		max(a.z, b.z),
	}
}

type State struct {
	cubes      map[Coord]bool
	sa         int
	mins, maxs Coord
}

func newState() *State {
	return &State{
		cubes: make(map[Coord]bool),
	}
}

func (s *State) pushCube(coord Coord) {
	s.cubes[coord] = true
	s.mins = minParts(&s.mins, &coord)
	s.maxs = maxParts(&s.maxs, &coord)
	for _, delta := range deltas {
		if s.cubes[coord.add(&delta)] {
			s.sa -= 1
		} else {
			s.sa += 1
		}
	}
}

func (s *State) cullVoids() {
	for x := s.mins.x + 1; x < s.maxs.x; x++ {
		for y := s.mins.y + 1; y < s.maxs.y; y++ {
			for z := s.mins.z + 1; z < s.maxs.z; z++ {
				coord := Coord{x, y, z}
				if !s.cubes[coord] {
					s.cullVoid(coord)
				}
			}
		}
	}
}

func (s *State) cullVoid(start Coord) {
	void := make(map[Coord]bool)
	voidSa := 0
	exterior := false
	stack := []Coord{start}
	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if void[curr] {
			// already counted
			continue
		}
		void[curr] = true
		if curr.x == s.mins.x || curr.y == s.mins.y || curr.z == s.mins.z ||
			curr.x == s.maxs.x || curr.y == s.maxs.y || curr.z == s.maxs.z {
			exterior = true
		}
		for _, delta := range deltas {
			coord := curr.add(&delta)
			occupied := false
			if coord.x < s.mins.x || coord.y < s.mins.y || coord.z < s.mins.z ||
				coord.x > s.maxs.x || coord.y > s.maxs.y || coord.z > s.maxs.z {
				occupied = true
			}
			if void[coord] {
				voidSa -= 1
				occupied = true
			} else {
				voidSa += 1
			}
			if s.cubes[coord] {
				occupied = true
			}
			if !occupied {
				stack = append(stack, coord)
			}
		}
	}

	if !exterior {
		s.sa -= voidSa
	}

	for coord := range void {
		s.cubes[coord] = true
	}

	location := "interior"
	if exterior {
		location = "exterior"
	}
	fmt.Printf("found %s void of size %d with sa %d at %v\n", location, len(void), voidSa, start)
}

func main() {
	s := newState()

	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var curr Coord
	scan := func() error {
		_, err := fmt.Fscanf(f, "%d,%d,%d\n", &curr.x, &curr.y, &curr.z)
		return err
	}
	for err := scan(); err == nil; err = scan() {
		s.pushCube(curr)
	}
	if err != nil && !errors.Is(err, io.EOF) {
		log.Fatal(err)
	}

	fmt.Println(s.sa)
	s.cullVoids()
	fmt.Println(s.sa)
}
