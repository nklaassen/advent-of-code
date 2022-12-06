package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	max := -1
	curr := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := scanner.Text()
		if len(l) == 0 {
			// empty line, next elf
			curr = 0
			continue
		}
		x, err := strconv.Atoi(l)
		if err != nil {
			log.Fatal(err)
		}
		curr += x
		if curr > max {
			max = curr
		}
	}

	fmt.Println(max)
}
