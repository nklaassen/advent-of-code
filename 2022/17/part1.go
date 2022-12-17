package main

import (
	"bytes"
	"log"
	"os"
)

const (
	totalShapes    = 2022
	width          = 7
	initialX       = 2
	initialYOffset = 3
)

func Print(args ...any) {
	//fmt.Print(args...)
}

func Println(args ...any) {
	//fmt.Println(args...)
}

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

type World struct {
	rocks         map[Pair]bool
	maxRockHeight int
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

func (w *World) print() {
	for y := w.maxRockHeight; y >= 0; y-- {
		Print("|")
		for x := 0; x < width; x++ {
			if w.rocks[Pair{x, y}] {
				Print("#")
			} else {
				Print(" ")
			}
		}
		Println("|")
	}
	Println("_________")
}

func main() {
	f, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	pushes := bytes.TrimSpace(f)

	world := &World{
		rocks: make(map[Pair]bool),
	}

	pushIndex := 0
	for shapeIndex := 0; shapeIndex < totalShapes; shapeIndex++ {
		shape := shapes[shapeIndex%len(shapes)]
		offset := Pair{initialX, world.maxRockHeight + initialYOffset}
		for {
			pushIndex = pushIndex % len(pushes)
			newOffset := offset
			switch pushes[pushIndex] {
			case '<':
				newOffset.x -= 1
				Print("left")
			case '>':
				newOffset.x += 1
				Print("right")
			}
			pushIndex++
			if world.collide(shape, newOffset) {
				Println("*")
				newOffset = offset
			} else {
				Println()
				offset = newOffset
			}
			newOffset.y -= 1
			Print("down")
			if world.collide(shape, newOffset) {
				Println("*")
				break
			} else {
				Println()
			}
			offset = newOffset
		}
		world.restShape(shape, offset)
		//world.print()
	}
	Println(world.maxRockHeight)
}
