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

func commonPriority(pack1, pack2, pack3 string) int {
	s := make(map[byte]struct{})
	for i := range pack1 {
		c1 := pack1[i]
		for j := range pack2 {
			c2 := pack2[j]
			if c1 == c2 {
				s[c1] = struct{}{}
				break
			}
		}
	}
	for k := range pack3 {
		c3 := pack3[k]
		if _, ok := s[c3]; ok {
			return priority(c3)
		}
	}
	return 999999999999
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	sum := 0
	var group []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		pack := scanner.Text()
		group = append(group, pack)
		if len(group) == 3 {
			sum += commonPriority(group[0], group[1], group[2])
			group = nil
		}
	}
	fmt.Println(sum)
}
