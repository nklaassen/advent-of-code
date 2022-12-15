package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	interestingY = 2000000
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	var s state
	s.beaconsOnLine = make(map[int]bool)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var sensor, beacon pair
		fmt.Sscanf(scanner.Text(), "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d",
			&sensor.x, &sensor.y, &beacon.x, &beacon.y)
		s.pushSensorBeacon(sensor, beacon)
	}

	fmt.Println(s.numCannotContainBeacon())
}

type state struct {
	sensors       []pair
	beacons       []pair
	spans         []pair
	beaconsOnLine map[int]bool
}

func (s *state) numCannotContainBeacon() int {
	a := s.spans[0]
	var merged []pair
	for i := 1; i < len(s.spans); i++ {
		b := s.spans[i]
		if a.y >= b.x {
			a = mergeSpans(a, b)
		} else {
			merged = append(merged, a)
			a = b
		}
	}
	merged = append(merged, a)

	count := 0
	for _, span := range merged {
		count += span.y - span.x + 1
	}
	count -= len(s.beaconsOnLine)
	return count
}

func mergeSpans(a, b pair) pair {
	return pair{a.x, max(a.y, b.y)}
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (s *state) pushSensorBeacon(sensor, beacon pair) {
	s.sensors = append(s.sensors, sensor)
	s.beacons = append(s.beacons, beacon)

	if beacon.y == interestingY {
		s.beaconsOnLine[beacon.x] = true
	}

	dBeacon := distance(sensor, beacon)
	dLine := abs(sensor.y - interestingY)

	spare := dBeacon - dLine

	// don't care about sensors too far to see the line
	if spare < 0 {
		return
	}

	x1, x2 := sensor.x, sensor.x
	x1 -= spare
	x2 += spare
	s.pushSpan(pair{x1, x2})
}

func (s *state) pushSpan(span pair) {
	for i := range s.spans {
		if span.x < s.spans[i].x {
			s.spans = append(s.spans[:i+1], s.spans[i:]...)
			s.spans[i] = span
			return
		}
	}
	s.spans = append(s.spans, span)
}

func distance(a, b pair) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type pair struct {
	x, y int
}
