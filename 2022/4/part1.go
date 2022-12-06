package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func firstContainsSecond(fl, fu, sl, su int) bool {
	return fl <= sl && fu >= su
}

func fullyContains(fl, fu, sl, su int) bool {
	return firstContainsSecond(fl, fu, sl, su) || firstContainsSecond(sl, su, fl, fu)
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

		if fullyContains(fl, fu, sl, su) {
			count++
		}
	}

	fmt.Println(count)
}
