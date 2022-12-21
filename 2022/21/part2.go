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
	flag.Parse()

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
				num:  float64(num),
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

	x := state.resolveX(monkeys["root"].inputs[1], 0)
	key := monkeys["root"].inputs[0]
	fn := func(x float64) float64 { return state.resolveX(key, x) }
	fmt.Println(monkeys["root"])
	fmt.Println(x)
	lower := 0
	higher := 9999999999999
	for lower < higher {
		mid := (lower + higher) / 2
		res := fn(float64(mid))
		fmt.Printf("%d %f %f\n", mid, res, x-res)
		if res == x {
			break
		}
		if res < x {
			higher = mid
		} else {
			lower = mid
		}
	}
}

type State struct {
	monkeys map[string]Monkey
}

func (s *State) resolveX(name string, x float64) float64 {
	if name == "humn" {
		return x
	}

	m := s.monkeys[name]
	if m.op == "" {
		return m.num
	}
	a := s.resolveX(m.inputs[0], x)
	b := s.resolveX(m.inputs[1], x)

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
	num    float64
	op     string
	inputs [2]string
}

func max(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}

func min(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}
