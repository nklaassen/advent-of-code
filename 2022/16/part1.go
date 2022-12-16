package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type valveSet string

func (v valveSet) push(valve string) valveSet {
	str := string(v)
	for i := 0; i < len(str); i += 2 {
		chunk := str[i : i+2]
		if valve > chunk {
			continue
		}
		if valve == chunk {
			return v
		}
		prefix := str[:i]
		suffix := str[i:]
		str = prefix + valve + suffix
		return valveSet(str)
	}
	str += valve
	return valveSet(str)
}

/*
func (v *valveSet) pop(valve string) valveSet {
	str := string(*v)
	for i := 0; i < len(str); i += 2 {
		chunk := str[i : i+2]
		if chunk != valve {
			continue
		}
		prefix := str[:i]
		suffix := str[i+2:]
		str = prefix + suffix
		*v = valveSet(str)
		return
	}
}
*/

func (v valveSet) contains(valve string) bool {
	str := string(v)
	for i := 0; i < len(str); i += 2 {
		chunk := str[i : i+2]
		if chunk == valve {
			return true
		}
	}
	return false
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	state := &State{
		valves:  make(map[string]int),
		tunnels: make(map[string][]string),
		M:       make(map[key]memo),
	}

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
		tunnels := strings.Split(strings.TrimSpace(parts[1]), ",")
		state.valves[valve] = rate
		for _, tunnel := range tunnels {
			state.tunnels[valve] = append(state.tunnels[valve], strings.TrimSpace(tunnel))
		}
	}

	var alreadyOpened valveSet
	fmt.Println(state.maxPressure("AA", 30, alreadyOpened))

	k := key{
		position: "AA",
		t:        30,
	}
	for k != (key{}) {
		m := state.M[k]
		fmt.Println(k, m)
		k = m.nextKey
	}
	bp()
}

func bp() {}

type memo struct {
	canRelease int
	opened     bool
	nextKey    key
}

type key struct {
	position      string
	t             int
	alreadyOpened valveSet
}

func (state *State) maxPressure(valve string, t int, alreadyOpened valveSet) (canRelease int) {
	if t <= 1 {
		return 0
	}
	if t == 2 {
		return state.valves[valve]
	}
	k := key{valve, t, alreadyOpened}
	if memo, ok := state.M[k]; ok {
		return memo.canRelease
	}
	/*
		defer func() {
			fmt.Println(valve, t, canRelease, alreadyOpened)
		}()
	*/

	canRelease = 0
	opened := false
	var nextKey key
	if rate := state.valves[valve]; rate > 0 && !alreadyOpened.contains(valve) {
		opened = true

		inner := alreadyOpened.push(valve)

		for _, tunnel := range state.tunnels[valve] {
			c := state.maxPressure(tunnel, t-2, inner)
			if c > canRelease {
				canRelease = c
				nextKey = key{
					position:      tunnel,
					t:             t - 2,
					alreadyOpened: inner,
				}
			}
		}

		canRelease += rate * (t - 1)
	}

	for _, tunnel := range state.tunnels[valve] {
		c := state.maxPressure(tunnel, t-1, alreadyOpened)
		if c > canRelease {
			opened = false
			canRelease = c
			nextKey = key{
				position:      tunnel,
				t:             t - 1,
				alreadyOpened: alreadyOpened,
			}
		}
	}

	state.M[k] = memo{
		canRelease: canRelease,
		opened:     opened,
		nextKey:    nextKey,
	}
	return canRelease
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

type State struct {
	valves  map[string]int
	tunnels map[string][]string

	M map[key]memo
}
