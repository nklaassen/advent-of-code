package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sort"
	"strings"

	"golang.org/x/exp/slices"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

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

	rate := make([]int8, N)
	for valve, r := range valveRates {
		i := valveNumbers[valve]
		rate[i] = int8(r)
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

	distVec := make([][]int8, N)
	for i := range distVec {
		distVec[i] = make([]int8, N)
	}
	for p, d := range dist {
		distVec[p.x][p.y] = int8(d)
		distVec[p.y][p.x] = int8(d)

	}

	// list all valves we care about (rate > 0)
	var valvesWithRate []int8
	for i, r := range rate {
		if r > 0 {
			valvesWithRate = append(valvesWithRate, int8(i))
		}
	}

	state := State{
		dist:    dist,
		distVec: distVec,
		rate:    rate,
		valves:  valvesWithRate,
		M:       make(map[KeyHash]Memo, 1100000),
	}

	m := state.maxRelease(Key{
		t: Pair{26, 26},
	})

	fmt.Println(m.best)
	fmt.Println(len(state.M))
}

func (s *State) maxRelease(key Key) Memo {
	var m Memo
	if key.t.x < 3 && key.t.y < 3 {
		return m
	}

	if m, ok := s.M[key.hash()]; ok {
		return m
	}

	for i, x := range s.valves {
		dx := s.distVec[key.pos.x][x] + 1
		txRemaining := key.t.x - dx
		if txRemaining < 1 {
			continue
		}
		if key.opened.has(int(x)) {
			continue
		}
		for j, y := range s.valves {
			if j == i {
				continue
			}
			dy := s.distVec[key.pos.y][y] + 1
			tyRemaining := key.t.y - dy
			if tyRemaining < 1 {
				continue
			}
			if key.opened.has(int(y)) {
				continue
			}

			nextKey := Key{
				pos:    Pair{x, y},
				t:      Pair{key.t.x - dx, key.t.y - dy},
				opened: key.opened.put2(int(x), int(y)),
			}
			nextKey.normalize()

			next := s.maxRelease(nextKey)

			currBest := next.best
			currBest += int16(s.rate[x]) * int16(txRemaining)
			currBest += int16(s.rate[y]) * int16(tyRemaining)

			if currBest > m.best {
				m.best = currBest
				m.nextKey = nextKey
			}
		}
	}

	if m.best == 0 {
		return m
	}

	s.M[key.hash()] = m
	return m
}

type State struct {
	dist    map[Pair]int
	distVec [][]int8
	rate    []int8
	valves  []int8
	M       map[KeyHash]Memo
}

type Key struct {
	pos    Pair
	t      Pair
	opened Bitset
}

func (k *Key) hash() KeyHash {
	var x uint32
	x |= uint32(k.pos.x) << 0
	x |= uint32(k.pos.y) << 8
	x |= uint32(k.t.x) << 16
	x |= uint32(k.t.y) << 24
	return KeyHash{x, uint32(k.opened), uint32(k.opened >> 32)}
}

type KeyHash [3]uint32

func (k *Key) normalize() {
	if k.pos.y < k.pos.x {
		k.pos.x, k.pos.y = k.pos.y, k.pos.x
		k.t.x, k.t.y = k.t.y, k.t.x
	}
}

type Memo struct {
	best    int16
	nextKey Key
}

type Pair struct {
	x, y int8
}

func orderedPair(x, y int) Pair {
	if x < y {
		return Pair{int8(x), int8(y)}
	}
	return Pair{int8(y), int8(x)}
}

type Bitset uint64

func (b Bitset) put(x int) Bitset {
	return b | (1 << x)
}

func (b Bitset) put2(x, y int) Bitset {
	return b | (1 << x) | (1 << y)
}

func (b *Bitset) has(x int) bool {
	return (*b)&(1<<x) != 0
}

func (b *Bitset) del(x int) Bitset {
	return (*b) & ^(1 << x)
}

func (b Bitset) count() int {
	c := 0
	for ; b != 0; b >>= 1 {
		if b&1 != 0 {
			c++
		}
	}
	return c
}

func (b Bitset) String() string {
	return fmt.Sprintf("%b", b)
}
