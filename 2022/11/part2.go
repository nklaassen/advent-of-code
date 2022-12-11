package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	totalRounds = 10000
)

var megaDivisor int

func main() {
	f, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	monkeyChunks := strings.Split(string(f), "\n\n")
	monkeys := make([]*Monkey, len(monkeyChunks))
	for i, monkeyChunk := range monkeyChunks {
		monkeys[i] = parseMonkey(monkeyChunk)
	}
	monkeys[6].operation = func(old int) int { return old * old }

	megaDivisor = int(1)
	for _, m := range monkeys {
		megaDivisor *= int(m.divisor)
	}

	for i := 0; i < totalRounds; i++ {
		round(monkeys)
	}

	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].numInspections > monkeys[j].numInspections
	})

	monkeyBusiness := monkeys[0].numInspections * monkeys[1].numInspections
	fmt.Println(monkeyBusiness)
}

func parseMonkey(monkeyChunk string) *Monkey {
	lines := strings.Split(monkeyChunk, "\n")

	header := lines[0]
	startingItems := lines[1]
	operationLine := lines[2]
	testLine := lines[3]
	ifTrue := lines[4]
	ifFalse := lines[5]

	m := &Monkey{}

	// number
	fmt.Sscanf(header, "Monkey %d:", &m.number)

	// items
	itemsStr := strings.Split(startingItems, ":")
	itemStrs := strings.Split(itemsStr[1], ",")
	for _, itemStr := range itemStrs {
		i, err := strconv.Atoi(strings.TrimSpace(itemStr))
		if err != nil {
			log.Fatal(err)
		}
		m.items = append(m.items, int(i))
	}

	// operation
	var operation string
	var operand int
	fmt.Sscanf(strings.TrimSpace(operationLine), "Operation: new = old %s %d", &operation, &operand)
	switch operation {
	case "*":
		m.operation = func(old int) int { return old * operand }
	case "+":
		m.operation = func(old int) int { return old + operand }
	}

	// test
	fmt.Sscanf(strings.TrimSpace(testLine), "Test: divisible by %d", &m.divisor)
	m.test = func(worry int) bool { return worry%m.divisor == 0 }

	// ifTrue ifFalse
	fmt.Sscanf(strings.TrimSpace(ifTrue), "If true: throw to monkey %d", &m.trueTarget)
	fmt.Sscanf(strings.TrimSpace(ifFalse), "If false: throw to monkey %d", &m.falseTarget)

	return m
}

func round(monkeys []*Monkey) {
	for _, monkey := range monkeys {
		monkey.inspectAll(monkeys)
	}
}

type Monkey struct {
	number                  int
	items                   []int
	operation               func(int) int
	divisor                 int
	test                    func(int) bool
	trueTarget, falseTarget int

	numInspections int
}

func (m *Monkey) String() string {
	return fmt.Sprintf("%v", m.items)
}

func (m *Monkey) inspectAll(allMonkeys []*Monkey) {
	for _, item := range m.items {
		m.inspect(allMonkeys, item)
	}
	m.items = nil
}

func (m *Monkey) inspect(allMonkeys []*Monkey, itemWorry int) {
	m.numInspections++
	newWorry := m.operation(itemWorry)
	newWorry %= megaDivisor
	var targetMonkeyIndex int
	if m.test(newWorry) {
		targetMonkeyIndex = m.trueTarget
	} else {
		targetMonkeyIndex = m.falseTarget
	}
	targetItems := &allMonkeys[targetMonkeyIndex].items
	*targetItems = append(*targetItems, newWorry)
}
