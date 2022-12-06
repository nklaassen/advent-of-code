package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	markerLen int = 4
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var b strings.Builder
	_, err = io.Copy(&b, f)
	if err != nil {
		log.Fatal(err)
	}

	input := b.String()

	for i := markerLen; i < len(input); i++ {
		if uniq(input[i-markerLen : i]) {
			fmt.Println(i)
			return
		}
	}
}

func uniq(s string) bool {
	for i := range s {
		c1 := s[i]
		for j := i + 1; j < len(s); j++ {
			c2 := s[j]
			if c1 == c2 {
				return false
			}
		}
	}
	return true
}
