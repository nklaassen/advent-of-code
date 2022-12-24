package main

import (
	"fmt"
	"math"
	"math/bits"
	"os"
	"strings"
)

const (
	maxTime = 1024
)

func main() {
	f, _ := os.ReadFile("input.txt")
	lines := strings.Split(strings.TrimSpace(string(f)), "\n")

	state := newState()
	state.valley[0] = make([][]BlizzardSet, len(lines))
	for i, line := range lines {
		state.valley[0][i] = make([]BlizzardSet, len(line))
		for j := range line {
			if line[j] == '.' {
				continue
			}
			state.valley[0][i][j] = BlizzardSet(parseBlizzard(line[j]))
		}
	}
	for i := 0; i < maxTime; i++ {
		state.computeNext(i)
	}
	start := SpaceTime{0, 0, 1}
	goal := Pair{26, 120}
	t1 := state.solve(start, goal)
	fmt.Println("part 1:", t1)
	start, goal = SpaceTime{t1, goal.i, goal.j}, Pair{start.i, start.j}
	t2 := state.solve(start, goal)
	fmt.Println("back to start:", t2)
	start, goal = SpaceTime{t2, goal.i, goal.j}, Pair{start.i, start.j}
	t3 := state.solve(start, goal)
	fmt.Println("part 2:", t3)
}

type State struct {
	// indexed by t, i, j
	valley [][][]BlizzardSet
}

func newState() *State {
	return &State{
		valley: make([][][]BlizzardSet, maxTime+1),
	}
}

type SpaceTime struct {
	t, i, j int
}

func (state *State) solve(start SpaceTime, goal Pair) int {
	// bfs in spacetime
	seen := map[SpaceTime]bool{
		start: true,
	}
	q := []SpaceTime{start}
	for len(q) > 0 {
		curr := q[0]
		q = q[1:]
		for _, offset := range []Pair{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {0, 0}} {
			next := SpaceTime{curr.t + 1, curr.i + offset.i, curr.j + offset.j}
			if next.i < 0 || next.i >= len(state.valley[0]) {
				continue
			}
			if next.i == goal.i && next.j == goal.j {
				return next.t
			}
			if state.valley[next.t][next.i][next.j].empty() && !seen[next] {
				seen[next] = true
				q = append(q, next)
			}
		}
	}
	return math.MaxInt
}

func (state *State) computeNext(t int) {
	state.valley[t+1] = make([][]BlizzardSet, len(state.valley[t]))
	for i := range state.valley[t] {
		state.valley[t+1][i] = make([]BlizzardSet, len(state.valley[t][i]))
	}
	for i := range state.valley[t] {
		for j := range state.valley[t][i] {
			for _, b := range state.valley[t][i][j].blizzards() {
				ii, jj := state.moveBlizzard(b, i, j)
				state.valley[t+1][ii][jj].set(b)
			}
		}
	}
}

func (state *State) printAt(t int) {
	for i := 0; i < len(state.valley[t]); i++ {
		for j := 0; j < len(state.valley[t][i]); j++ {
			fmt.Print(state.valley[t][i][j])
		}
		fmt.Println()
	}
}

func (state *State) moveBlizzard(b Blizzard, i, j int) (int, int) {
	var ii, jj int
	switch b {
	case 1:
		// >
		ii, jj = i, j+1
	case 2:
		// v
		ii, jj = i+1, j
	case 4:
		// <
		ii, jj = i, j-1
	case 8:
		// ^
		ii, jj = i-1, j
	case 16:
		// #
		return i, j
	default:
		panic(b)
	}
	numRows := len(state.valley[0])
	numCols := len(state.valley[0][0])
	switch {
	case ii == 0:
		ii = numRows - 2
	case ii == numRows-1:
		ii = 1
	case jj == 0:
		jj = numCols - 2
	case jj == numCols-1:
		jj = 1
	}
	return ii, jj
}

type Blizzard uint8

func (b Blizzard) String() string {
	switch b {
	case 1:
		return ">"
	case 2:
		return "v"
	case 4:
		return "<"
	case 8:
		return "^"
	case 16:
		return "#"
	}
	return " "
}

func parseBlizzard(b byte) Blizzard {
	switch b {
	case '>':
		return 1
	case 'v':
		return 2
	case '<':
		return 4
	case '^':
		return 8
	case '#':
		return 16
	}
	panic(b)
}

type BlizzardSet uint8

func (b BlizzardSet) String() string {
	if pop := bits.OnesCount8(uint8(b)); pop <= 1 {
		return Blizzard(b).String()
	} else {
		return fmt.Sprint(pop)
	}
}

func (b BlizzardSet) blizzards() (blizzards []Blizzard) {
	for bit := uint8(1); bit <= 16; bit <<= 1 {
		if uint8(b)&bit != 0 {
			blizzards = append(blizzards, Blizzard(bit))
		}
	}
	return
}

func (b *BlizzardSet) set(blizzard Blizzard) {
	*b |= BlizzardSet(blizzard)
}

func (b BlizzardSet) empty() bool {
	return uint8(b) == 0
}

type Pair struct {
	i, j int
}
