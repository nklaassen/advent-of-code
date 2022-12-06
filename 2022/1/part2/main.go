package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	curr := 0
	var elfCounts []int

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := scanner.Text()
		if len(l) == 0 {
			// empty line, next elf
			elfCounts = append(elfCounts, curr)
			curr = 0
			continue
		}
		x, err := strconv.Atoi(l)
		if err != nil {
			log.Fatal(err)
		}
		curr += x
	}
	elfCounts = append(elfCounts, curr)

	sort.Ints(elfCounts)

	l := len(elfCounts)

	topThreeSum := 0
	for i := 0; i < 3; i++ {
		topThreeSum += elfCounts[l-1-i]
	}
	fmt.Println(topThreeSum)
}
