package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func cIsBetween(a, b, c int) bool {
	return c >= a && c <= b
}

func overlap(fl, fu, sl, su int) bool {
	return cIsBetween(fl, fu, sl) ||
		cIsBetween(fl, fu, su) ||
		cIsBetween(sl, su, fl) ||
		cIsBetween(sl, su, fu)
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	count := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		assignments := strings.Split(line, ",")
		bounds := strings.Split(assignments[0], "-")
		fl, _ := strconv.Atoi(bounds[0])
		fu, _ := strconv.Atoi(bounds[1])
		bounds = strings.Split(assignments[1], "-")
		sl, _ := strconv.Atoi(bounds[0])
		su, _ := strconv.Atoi(bounds[1])

		if overlap(fl, fu, sl, su) {
			count++
		}
	}

	fmt.Println(count)
}
