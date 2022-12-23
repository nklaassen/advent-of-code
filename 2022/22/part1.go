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
	r := input.rows[1]
	for i, b := range r {
		if b == '.' {
			s.pos += complex(0, -float32(i))
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
			newRow, newCol := row(newPos), col(newPos)
			if newRow < 0 || newRow >= len(s.input.rows) ||
				newCol < 0 || newCol >= len(s.input.rows[newRow]) ||
				s.input.rows[row(newPos)][col(newPos)] == ' ' {
				newPos = s.warp()
			}
			if s.input.rows[row(newPos)][col(newPos)] == '#' {
				break
			}
			s.output.rows[row(newPos)][col(newPos)] = 'O'
			s.pos = newPos
		}
		fmt.Println(&step)
		s.output.printNear(row(s.pos))
	}
}

func (s *State) warp() complex64 {
	var newPos complex64
	if real(s.dir) > 0 {
		newPos = complex(0, imag(s.pos))
	}
	if real(s.dir) < 0 {
		newPos = complex(float32(len(s.input.rows)-1), imag(s.pos))
	}
	if imag(s.dir) < 0 {
		newPos = complex(real(s.pos), 0)
	}
	if imag(s.dir) > 0 {
		newPos = complex(real(s.pos), -float32(len(s.input.rows[row(s.pos)])-1))
	}
	for col(newPos) >= len(s.input.rows[row(newPos)]) ||
		s.input.rows[row(newPos)][col(newPos)] == ' ' {
		newPos += s.dir
	}
	return newPos
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

func row(c complex64) int {
	return int(real(c))
}

func col(c complex64) int {
	return -int(imag(c))
}
