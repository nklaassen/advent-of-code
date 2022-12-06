package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	numStacks        int = 9
	maxHeight            = 8
	colWidth             = 4
	crateOffset          = 1
	firstCommandLine     = 10
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	stacks := make([][]byte, numStacks)

	scanner := bufio.NewScanner(f)
	for lineNum := 0; scanner.Scan(); lineNum++ {
		line := scanner.Text()

		if lineNum < maxHeight {
			for stackNum := 0; stackNum < numStacks; stackNum++ {
				if crate := line[colWidth*stackNum+crateOffset]; crate != ' ' {
					stacks[stackNum] = append(stacks[stackNum], crate)
				}
			}
		}
		if lineNum == maxHeight {
			// stacks were originally read top to bottom
			for _, stack := range stacks {
				reverse(stack)
			}
		}
		if lineNum >= firstCommandLine {
			var num int
			var from int
			var to int
			_, err := fmt.Sscanf(line, "move %d from %d to %d", &num, &from, &to)
			if err != nil {
				log.Fatal(err)
			}
			move(num, &stacks[from-1], &stacks[to-1])
		}
	}
	for _, stack := range stacks {
		fmt.Printf("%c", pop(&stack))
	}
	fmt.Println()
}

func move(num int, from, to *[]byte) {
	for i := 0; i < num; i++ {
		*to = append(*to, pop(from))
	}
}

func pop(s *[]byte) byte {
	i := len(*s) - 1
	top := (*s)[i]
	*s = (*s)[:i]
	return top
}

func reverse(s []byte) {
	for i := 0; i < len(s)/2; i++ {
		j := len(s) - i - 1
		s[i], s[j] = s[j], s[i]
	}
}
