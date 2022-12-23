package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var input = flag.String("input", "input.txt", "input file")

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	flag.Parse()

	input, err := parseInput(*input)
	if err != nil {
		return err
	}

	state := newState(input)

	state.run()
	state.output.printAll()

	fmt.Println(state.pos, state.dir)
	finalRow := int(real(state.pos)) + 1
	finalCol := int(-imag(state.pos)) + 1
	var finalDir int
	switch state.dir {
	case -1i:
		finalDir = 0
	case 1:
		finalDir = 1
	case 1i:
		finalDir = 2
	case -1:
		finalDir = 3
	}
	fmt.Println(finalRow, finalCol, finalDir)
	fmt.Println(1000*finalRow + 4*finalCol + finalDir)
	return nil
}

func parseInput(filename string) (*Input, error) {
	f, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	parts := bytes.Split(f, []byte("\n\n"))
	m := parts[0]
	path := strings.TrimSpace(string(parts[1]))

	var input Input

	input.rows = bytes.Split(m, []byte("\n"))

	var dir string
	var tiles int
	r := strings.NewReader("R" + path)
	scan := func() error {
		_, err := fmt.Fscanf(r, "%1s%d", &dir, &tiles)
		return err
	}
	for err = scan(); err == nil; err = scan() {
		input.path = append(input.path, Step{
			dist:    tiles,
			turn:    parseTurn(dir),
			rawTurn: dir,
		})
	}
	if !errors.Is(err, io.EOF) {
		return nil, err
	}
	input.path[0].rawTurn = ""
	input.path[0].turn = 1

	return &input, nil
}

type Input struct {
	rows [][]byte
	path []Step
}

func (i *Input) Clone() *Input {
	rowsCopy := make([][]byte, len(i.rows))
	for i, row := range i.rows {
		rowsCopy[i] = make([]byte, len(row))
		copy(rowsCopy[i], row)
	}
	return &Input{
		rows: rowsCopy,
		path: i.path,
	}
}

type Step struct {
	dist    int
	turn    complex64
	rawTurn string
}

func (s *Step) String() string {
	return fmt.Sprintf("%s %d", s.rawTurn, s.dist)
}

func parseTurn(d string) complex64 {
	switch d {
	case "R":
		return 1i
	case "L":
		return -1i
	default:
		panic(fmt.Sprintf("unknown direction %q", d))
	}
}

type State struct {
	pos, dir      complex64
	input, output *Input
}

func newState(input *Input) *State {
	s := &State{
		dir:    -1i,
		pos:    0,
		input:  input,
		output: input.Clone(),
	}
	r := input.rows[0]
	for i, b := range r {
		if b == '.' {
			s.pos = pos(0, i)
			break
		}
	}
	return s
}

func (s *State) run() {
	for _, step := range s.input.path {
		s.dir *= step.turn
		for i := 0; i < step.dist; i++ {
			newPos := s.pos + s.dir
			newPos, newDir := s.warp(newPos, s.dir)
			if s.input.rows[row(newPos)][col(newPos)] == '#' {
				break
			}
			s.output.rows[row(newPos)][col(newPos)] = dirChar(newDir)
			s.pos, s.dir = newPos, newDir
		}
		fmt.Println(&step)
		s.output.printNear(row(s.pos))
	}
}

func dirChar(dir complex64) byte {
	switch dir {
	case 1:
		return 'v'
	case -1:
		return '^'
	case 1i:
		return '<'
	case -1i:
		return '>'
	default:
		panic("unexpected dir")
	}
}

func (s *State) warp(newPos complex64, dir complex64) (complex64, complex64) {
	r := row(newPos)
	c := col(newPos)
	if r >= 0 && r < len(s.input.rows) &&
		c >= 0 && c < len(s.input.rows[r]) &&
		s.input.rows[r][c] != ' ' {
		//fmt.Println("no warping necessary")
		return newPos, dir
	}
	fmt.Println("warping required")
	ogDir := dir
	defer func() {
		fmt.Printf("warped from %v to %v\n", newPos, pos(r, c))
		fmt.Printf("dir warped from %v to %v\n", ogDir, dir)
	}()
	switch {
	// top
	case r == -1 && between(c, 50, 100):
		r = 150 + c - 50
		c = 0
		dir *= 1i
	case r == -1 && between(c, 100, 150):
		c = c - 100
		r = 199
		// same
	// top lr
	case c == 49 && between(r, 0, 50):
		r = 149 - r
		c = 0
		dir *= -1
	case c == 150 && between(r, 0, 50):
		r = 149 - r
		c = 99
		dir *= -1
	// bottom of top
	case r == 50 && between(c, 100, 150) && dir == 1:
		r = 50 + c - 100
		c = 99
		dir *= 1i
	// upper mid lr
	case c == 49 && between(r, 50, 100) && dir == 1i:
		c = r - 50
		r = 100
		dir *= -1i
	case c == 100 && between(r, 50, 100) && dir == -1i:
		c = 100 + r - 50
		r = 49
		dir *= -1i
	// top of bottom left
	case r == 99 && between(c, 0, 50) && dir == -1:
		r = 50 + c
		c = 50
		dir *= 1i
	// lower mid lr
	case c == -1 && between(r, 100, 150):
		r = 49 - (r - 100)
		c = 50
		dir *= -1
	case c == 100 && between(r, 100, 150):
		r = 49 - (r - 100)
		c = 149
		dir *= -1
	// bottom of mid
	case r == 150 && between(c, 50, 100):
		r = 150 + c - 50
		c = 49
		dir *= 1i
	// bottom lr
	case c == -1 && between(r, 150, 200):
		c = 50 + r - 150
		r = 0
		dir *= -1i
	case c == 50 && between(r, 150, 200):
		c = 50 + r - 150
		r = 149
		dir *= -1i
	// bottom
	case r == 200 && between(c, 0, 50):
		c = 100 + c
		r = 0
		// same
	default:
		panic(fmt.Sprintf("unexpected r %d c %d", r, c))
	}
	return pos(r, c), dir
}

func between(x, low, high int) bool {
	return x >= low && x < high
}

func (input *Input) printNear(j int) {
	min := j - 5
	if min < 0 {
		min = 0
	}
	for i := min; i < min+10 && i < len(input.rows); i++ {
		fmt.Println(string(input.rows[i]))
	}
	fmt.Println()
}

func (input *Input) printAll() {
	for i := 0; i < len(input.rows); i++ {
		fmt.Println(string(input.rows[i]))
	}
	fmt.Println()
}

func row(c complex64) int {
	return int(real(c))
}

func col(c complex64) int {
	return -int(imag(c))
}

func pos(r, c int) complex64 {
	return complex(float32(r), float32(-c))
}
