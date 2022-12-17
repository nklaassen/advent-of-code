package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"golang.org/x/exp/slices"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	valveRates := make(map[string]int)
	var valves []string
	tunnels := make(map[string][]string)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		var valve string
		var rate int
		_, err := fmt.Sscanf(line, "Valve %s has flow rate=%d;", &valve, &rate)
		if err != nil {
			log.Fatal(err)
		}
		parts := strings.Split(line, "valves")
		parts = strings.Split(strings.TrimSpace(parts[1]), ",")
		valveRates[valve] = rate
		valves = append(valves, valve)
		for _, tunnel := range parts {
			tunnels[valve] = append(tunnels[valve], strings.TrimSpace(tunnel))
		}
		i++
	}

	N := len(valves)

	sort.Strings(valves)
	valveNumbers := make(map[string]int, N)
	for i, valve := range valves {
		valveNumbers[valve] = i
	}

	rate := make([]int, N)
	for valve, r := range valveRates {
		i := valveNumbers[valve]
		rate[i] = r
	}

	dist := make(map[Pair]int, N*N)

	// add baseline "1 or inf" distances
	for i := 0; i < N; i++ {
		for j := i + 1; j < N; j++ {
			d := 1000
			if slices.Contains(tunnels[valves[i]], valves[j]) {
				d = 1
			}
			dist[orderedPair(i, j)] = d
		}
	}

	// floyd-warshall min distance between each pair of nodes
	for k := 0; k < N; k++ {
		for i := 0; i < N; i++ {
			for j := 0; j < N; j++ {
				d := dist[orderedPair(i, k)] + dist[orderedPair(k, j)]
				if dist[orderedPair(i, j)] > d {
					dist[orderedPair(i, j)] = d
				}
			}
		}
	}

	// remove all nodes where rate is 0, except the start
	for i := 0; i < N; i++ {
		for j := i + 1; j < N; j++ {
			if (i == 0 && rate[j] == 0) || (i > 0 && rate[i] == 0 || rate[j] == 0) {
				delete(dist, orderedPair(i, j))
			}
		}
	}

	// list all valves we care about (rate > 0)
	var valvesWithRate []int
	for i, r := range rate {
		if r > 0 {
			valvesWithRate = append(valvesWithRate, i)
		}
	}

	state := State{
		dist:   dist,
		rate:   rate,
		valves: valvesWithRate,
		M:      make(map[Key]Memo),
	}

	m := state.maxRelease(Key{
		t: Pair{26, 26},
	})

	fmt.Println(m.best)

	key := m.nextKey
	for key != (Key{}) {
		m = state.M[key]
		fmt.Println(key, m.best)
		key = m.nextKey
	}
}

func (s *State) maxRelease(key Key) Memo {
	if m, ok := s.M[key]; ok {
		return m
	}
	var m Memo

	if key.t.x <= 1 && key.t.y <= 1 {
		return m
	}

	for _, x := range s.valves {
		if key.opened.has(x) {
			continue
		}
		dx := s.dist[orderedPair(key.pos.x, x)] + 1
		if key.t.x-dx < 1 {
			continue
		}
		for _, y := range s.valves {
			if y == x {
				continue
			}
			if key.opened.has(y) {
				continue
			}
			dy := s.dist[orderedPair(key.pos.y, y)] + 1
			if key.t.y-dy < 1 {
				continue
			}

			nextKey := Key{
				pos:    Pair{x, y},
				t:      Pair{key.t.x - dx, key.t.y - dy},
				opened: key.opened.put(x).put(y),
			}
			nextKey.normalize()

			next := s.maxRelease(nextKey)

			currBest := next.best
			txRemaining := key.t.x - dx
			tyRemaining := key.t.y - dy
			currBest += s.rate[x] * txRemaining
			currBest += s.rate[y] * tyRemaining

			if currBest > m.best {
				m.best = currBest
				m.nextKey = nextKey
			}
		}
	}

	s.M[key] = m
	return m
}

type State struct {
	dist   map[Pair]int
	rate   []int
	valves []int
	M      map[Key]Memo
}

type Key struct {
	pos    Pair
	t      Pair
	opened Bitset
}

func (k *Key) normalize() {
	if k.pos.y < k.pos.x {
		k.pos.x, k.pos.y = k.pos.y, k.pos.x
		k.t.x, k.t.y = k.t.y, k.t.x
	}
}

type Memo struct {
	best    int
	nextKey Key
}

type Pair struct {
	x, y int
}

func orderedPair(x, y int) Pair {
	if x < y {
		return Pair{x, y}
	}
	return Pair{y, x}
}

type Bitset uint64

func (b Bitset) put(x int) Bitset {
	b |= 1 << x
	return b
}

func (b Bitset) has(x int) bool {
	return b&(1<<x) != 0
}

func (b Bitset) del(x int) Bitset {
	b &= ^(1 << x)
	return b
}

func (b Bitset) count() int {
	c := 0
	for b != 0 {
		if b&1 != 0 {
			c++
		}
		b >>= 1
	}
	return c
}

func (b Bitset) String() string {
	return fmt.Sprintf("%b", b)
}
