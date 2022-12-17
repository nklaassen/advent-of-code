package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

const (
	totalShapes = 1000000000000
	//totalShapes    = 2022
	width          = 7
	initialX       = 2
	initialYOffset = 3
)

type Pair struct {
	x, y int
}

// offsets of rocks indexed from bottom left corner (even if it does not exist)
type Shape []Pair

func (s Shape) positions(base Pair) []Pair {
	var out []Pair
	for _, offset := range s {
		out = append(out, Pair{base.x + offset.x, base.y + offset.y})
	}
	return out
}

var (
	horizontal Shape = []Pair{
		{0, 0},
		{1, 0},
		{2, 0},
		{3, 0},
	}
	plus Shape = []Pair{
		{1, 0},
		{0, 1},
		{1, 1},
		{2, 1},
		{1, 2},
	}
	backwardL Shape = []Pair{
		{0, 0},
		{1, 0},
		{2, 0},
		{2, 1},
		{2, 2},
	}
	vertical Shape = []Pair{
		{0, 0},
		{0, 1},
		{0, 2},
		{0, 3},
	}
	square Shape = []Pair{
		{0, 0},
		{0, 1},
		{1, 0},
		{1, 1},
	}

	shapes = []Shape{
		horizontal,
		plus,
		backwardL,
		vertical,
		square,
	}
)

type Fingerprint [25]int

type World struct {
	rocks         map[Pair]bool
	maxRockHeight int
	pushes        []byte
	lastSeen      map[Fingerprint]int

	shapeIndex, pushIndex int
}

func (w *World) collide(shape Shape, offset Pair) bool {
	for _, pos := range shape.positions(offset) {
		switch {
		case pos.x < 0:
			return true
		case pos.x >= width:
			return true
		case pos.y < 0:
			return true
		case w.rocks[Pair{pos.x, pos.y}]:
			return true
		}
	}
	return false
}

func (w *World) restShape(shape Shape, offset Pair) {
	for _, pos := range shape.positions(offset) {
		w.rocks[pos] = true
		if pos.y >= w.maxRockHeight {
			w.maxRockHeight = pos.y + 1
		}
	}
}

func (w *World) pushShape() int {
	offset := Pair{initialX, w.maxRockHeight + initialYOffset}
	for {
		newOffset := offset
		switch w.pushes[w.pushIndex] {
		case '<':
			newOffset.x -= 1
		case '>':
			newOffset.x += 1
		}
		w.pushIndex++
		w.pushIndex = w.pushIndex % len(w.pushes)
		if w.collide(shapes[w.shapeIndex], newOffset) {
			newOffset = offset
		} else {
			offset = newOffset
		}
		newOffset.y -= 1
		if w.collide(shapes[w.shapeIndex], newOffset) {
			break
		}
		offset = newOffset
	}
	prevHeight := w.maxRockHeight
	w.restShape(shapes[w.shapeIndex], offset)
	w.shapeIndex++
	w.shapeIndex = w.shapeIndex % len(shapes)
	return w.maxRockHeight - prevHeight
}

func (w *World) print() {
	for y := w.maxRockHeight; y >= 0; y-- {
		fmt.Print("|")
		for x := 0; x < width; x++ {
			if w.rocks[Pair{x, y}] {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println("|")
	}
	fmt.Println("_________")
}

func main() {
	f, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	pushes := bytes.TrimSpace(f)

	world := &World{
		rocks:  make(map[Pair]bool),
		pushes: pushes,
	}

	shapeCount := 0
	var fingerprint Fingerprint
	lastSeen := make(map[Fingerprint]int)
	lastHeight := make(map[Fingerprint]int)
	height := 0

	cycleLen, cycleHeight := 0, 0
	for shapeCount < totalShapes {
		dHeight := world.pushShape()
		height += dHeight
		shapeCount++

		copy(fingerprint[1:], fingerprint[:])
		fingerprint[0] = dHeight

		if lastSeen[fingerprint] > 0 {
			cycleLen, cycleHeight = shapeCount-lastSeen[fingerprint], height-lastHeight[fingerprint]
			fmt.Println("detected cycle", cycleLen, cycleHeight)
			break
		}
		lastSeen[fingerprint] = shapeCount
		lastHeight[fingerprint] = height
	}

	remainingShapes := totalShapes - shapeCount
	remainingFullCycles := remainingShapes / cycleLen
	shapeCount += remainingFullCycles * cycleLen
	height += remainingFullCycles * cycleHeight
	for ; shapeCount < totalShapes; shapeCount++ {
		height += world.pushShape()
	}

	fmt.Println(height)
}
