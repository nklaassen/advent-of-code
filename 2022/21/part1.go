package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var input = flag.String("input", "input.txt", "input file")

func main() {
	f, err := os.Open(*input)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	monkeys := make(map[string]Monkey)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		name := parts[0]
		remaining := strings.TrimSpace(parts[1])
		parts = strings.Split(remaining, " ")
		if len(parts) == 1 {
			num, err := strconv.Atoi(parts[0])
			if err != nil {
				log.Fatal(err)
			}
			monkeys[name] = Monkey{
				name: name,
				num:  num,
			}
			continue
		}
		monkeys[name] = Monkey{
			name:   name,
			op:     parts[1],
			inputs: [2]string{parts[0], parts[2]},
		}
	}

	state := State{
		monkeys: monkeys,
	}
	fmt.Println(state.resolve("root"))
}

type State struct {
	monkeys map[string]Monkey
}

func (s *State) resolve(name string) int {
	m := s.monkeys[name]
	if m.op == "" {
		return m.num
	}
	a := s.resolve(m.inputs[0])
	b := s.resolve(m.inputs[1])
	switch m.op {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		return a / b
	default:
		panic("unknown op")
	}
}

type Monkey struct {
	name   string
	num    int
	op     string
	inputs [2]string
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
