package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func priority(c byte) int {
	if c >= 'a' {
		return 1 + int(c-'a')
	}
	return 27 + int(c-'A')
}

func commonPriority(pack string) int {
	for i := 0; i < len(pack)/2; i++ {
		for j := len(pack) / 2; j < len(pack); j++ {
			if pack[i] == pack[j] {
				return priority(pack[i])
			}
		}
	}
	return 999999999999999999
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	sum := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		pack := scanner.Text()
		sum += commonPriority(pack)
	}
	fmt.Println(sum)
}
