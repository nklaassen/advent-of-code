package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	minY = 0
	maxY = 4000000
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	var s state

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var sensor, beacon pair
		fmt.Sscanf(scanner.Text(), "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d",
			&sensor.x, &sensor.y, &beacon.x, &beacon.y)
		s.pushSensorBeacon(sensor, beacon)
	}

	// naively check every y value since we already know how to count possible
	// locations per line
	y := minY
	for ; y <= maxY; y++ {
		fmt.Println("checking", y)
		if num := s.numCannotContainBeacon(y); num != maxY+1 {
			break
		}
	}

	x := s.spans[0].y + 1
	freq := x*maxY + y
	fmt.Println(freq)
}

type state struct {
	sensors []pair
	beacons []pair
	spans   []pair
}

func (s *state) numCannotContainBeacon(y int) int {
	s.reset()

	for i := range s.sensors {
		sensor := s.sensors[i]
		beacon := s.beacons[i]

		dBeacon := distance(sensor, beacon)
		dLine := abs(sensor.y - y)

		spare := dBeacon - dLine

		// don't care about sensors too far to see the line
		if spare < 0 {
			continue
		}

		x1, x2 := sensor.x, sensor.x
		x1 -= spare
		x2 += spare
		if x1 < 0 {
			x1 = 0
		}
		if x2 > maxY {
			x2 = maxY
		}
		if x2 < x1 {
			continue
		}
		s.pushSpan(pair{x1, x2})
	}

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
	s.spans = merged

	count := 0
	for _, span := range merged {
		count += span.y - span.x + 1
	}
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

func (s *state) reset() {
	s.spans = nil
}

func (s *state) pushSensorBeacon(sensor, beacon pair) {
	s.sensors = append(s.sensors, sensor)
	s.beacons = append(s.beacons, beacon)
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
