package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var input = flag.String("input", "input.txt", "input file")

const (
	totalTime     = 24
	botTypes      = 4
	resourceTypes = 4
	oreIndex      = 0
	clayIndex     = 1
	obsidianIndex = 2
	geodeIndex    = 3
)

func main() {
	flag.Parse()

	f, err := os.Open(*input)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var blueprint Blueprint
	scan := func() error {
		_, err := fmt.Fscanf(f, "Blueprint %d: "+
			"Each ore robot costs %d ore. "+
			"Each clay robot costs %d ore. "+
			"Each obsidian robot costs %d ore and %d clay. "+
			"Each geode robot costs %d ore and %d obsidian.\n",
			&blueprint.id,
			&blueprint.bots[0].cost[oreIndex],
			&blueprint.bots[1].cost[oreIndex],
			&blueprint.bots[2].cost[oreIndex], &blueprint.bots[2].cost[clayIndex],
			&blueprint.bots[3].cost[oreIndex], &blueprint.bots[3].cost[obsidianIndex],
		)
		return err
	}
	sum := 0
	for err = scan(); err == nil; err = scan() {
		for i := 0; i < botTypes; i++ {
			blueprint.bots[i].output[i] = 1
		}

		s := newState(&blueprint)
		key := newKey()
		geodes := s.maxGeodes(key)
		quality := blueprint.id * geodes
		sum += quality

		fmt.Printf("geodes: %d quality: %d sum: %d\n", geodes, quality, sum)

		for *key != (Key{}) {
			m := s.M[*key]
			fmt.Println(key, m)
			key = &m.nextKey
		}
	}
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		log.Fatal(err)
	}
}

type State struct {
	M             map[Key]Memo
	blueprint     *Blueprint
	maxProduction Supplies
	maxSeen       int
}

func newState(blueprint *Blueprint) *State {
	s := &State{
		M:         make(map[Key]Memo),
		blueprint: blueprint,
	}
	for i := 0; i < resourceTypes; i++ {
		for j := 0; j < botTypes; j++ {
			s.maxProduction[i] = max(s.maxProduction[i], blueprint.bots[j].cost[i])
		}
	}
	s.maxProduction[optimizingIndex] = 127
	return s
}

var optimizingIndex = geodeIndex

func (s *State) maxGeodes(key *Key) (maxGeodes int) {
	defer func() {
		if maxGeodes > s.maxSeen {
			s.maxSeen = maxGeodes
		}
	}()

	timeRemaining := int(totalTime - key.time)
	finalGeodesWithoutPurchasing := int(key.supplies[optimizingIndex]) + timeRemaining*int(key.botCounts[optimizingIndex])
	if timeRemaining < 2 {
		return finalGeodesWithoutPurchasing
	}

	finalGeodesBuyingGeodeBotEveryTurn := int(finalGeodesWithoutPurchasing) + timeRemaining*(timeRemaining-1)/2
	if finalGeodesBuyingGeodeBotEveryTurn < s.maxSeen {
		// Kill this branch that will never be best
		return 0
	}

	m, ok := s.M[*key]
	if ok {
		return m.maxGeodes
	}

	var newSupplies Supplies
	for i := 0; i < botTypes; i++ {
		newSupplies.add(mul(&s.blueprint.bots[i].output, key.botCounts[i]))
	}

	options := []Key{}
	for nextPurchase := 0; nextPurchase < botTypes; nextPurchase++ {
		if key.botCounts[nextPurchase] >= s.maxProduction[nextPurchase] {
			continue
		}
		opt := *key
		opt.time++
		for opt.time < totalTime {
			opt.supplies.sub(&s.blueprint.bots[nextPurchase].cost)
			if opt.supplies.valid() {
				opt.botCounts[nextPurchase]++
				break
			}
			opt.supplies.add(&s.blueprint.bots[nextPurchase].cost)
			opt.supplies.add(&newSupplies)
			opt.time++
		}
		if opt.time < totalTime {
			options = append(options, opt)
		}
	}

	if len(options) == 0 {
		return int(key.supplies[optimizingIndex] + key.botCounts[optimizingIndex]*(totalTime-key.time))
	}

	for i := range options {
		options[i].supplies.add(&newSupplies)
	}

	var nextKey Key
	for i := range options {
		if m := s.maxGeodes(&options[i]); m >= maxGeodes {
			maxGeodes = m
			nextKey = options[i]
		}
	}

	s.M[*key] = Memo{maxGeodes: maxGeodes, nextKey: nextKey}
	return maxGeodes
}

func max[T int8 | int](x, y T) T {
	if x > y {
		return x
	}
	return y
}

type Memo struct {
	maxGeodes int
	nextKey   Key
}

type Key struct {
	supplies, botCounts Supplies
	time                int8
}

func (k Key) String() string {
	return fmt.Sprintf("{supplies: %v bots: %v t: %d}", k.supplies, k.botCounts, k.time)
}

func newKey() *Key {
	s := &Key{}
	s.botCounts[oreIndex] = 1
	return s
}

type Blueprint struct {
	id   int
	bots [botTypes]Bot
}

type Supplies [resourceTypes]int8

func (s *Supplies) valid() bool {
	for _, c := range *s {
		if c < 0 {
			return false
		}
	}
	return true
}

func (s *Supplies) sub(o *Supplies) {
	for i := range s {
		s[i] -= o[i]
	}
}

func (s *Supplies) add(o *Supplies) {
	for i := range s {
		s[i] += o[i]
	}
}

func mul(s *Supplies, n int8) *Supplies {
	var r Supplies
	for i := range s {
		r[i] = s[i] * n
	}
	return &r
}

type Bot struct {
	cost, output Supplies
}
