package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var input = flag.String("input", "input.txt", "input file")

func main() {
	flag.Parse()

	f, _ := os.Open(*input)
	scanner := bufio.NewScanner(f)
	sum := 0
	for scanner.Scan() {
		sum += parseSnafu(scanner.Text())
	}
	fmt.Println(renderSnafu(sum))
}

var values = map[byte]int{
	'=': -2,
	'-': -1,
	'0': 0,
	'1': 1,
	'2': 2,
}

func parseSnafu(str string) int {
	var n int
	m := 1
	for i := len(str) - 1; i >= 0; i-- {
		n += values[str[i]] * m
		m *= 5
	}
	return n
}

type output struct {
	char  string
	carry int
}

var outputs = map[int]output{
	0: {"0", 0},
	1: {"1", 0},
	2: {"2", 0},
	3: {"=", 1},
	4: {"-", 1},
}

func renderSnafu(n int) string {
	var str string
	carry := 0
	for n+carry > 0 {
		r := (n + carry) % 5
		o := outputs[r]
		str = o.char + str
		carry = o.carry
		n /= 5
	}
	return str
}
