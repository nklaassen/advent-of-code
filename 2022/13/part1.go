package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	linePairs := bytes.Split(f, []byte{'\n', '\n'})

	sum := 0
	for i, linePair := range linePairs {
		index := i + 1

		lines := bytes.Split(linePair, []byte{'\n'})

		first := parseNode(string(lines[0]))
		second := parseNode(string(lines[1]))

		fmt.Println(first)
		fmt.Println(second)

		if order := compare(first, second); order < 0 {
			sum += index
			fmt.Println("in order")
		} else if order > 0 {
			fmt.Println("out of order")
		} else if order == 0 {
			panic("pairs compared equal")
		}
		fmt.Println()
	}

	fmt.Println(sum)
}

type Node interface {
	String() string
}

type IntNode struct {
	x      int
	length int
}

func (i IntNode) String() string {
	return fmt.Sprintf("%d", i.x)
}

type ListNode struct {
	nodes  []Node
	length int
}

func (l ListNode) String() string {
	return fmt.Sprintf("%v", l.nodes)
}

func parseNode(str string) Node {
	switch str[0] {
	case '[':
		return parseList(str)
	default:
		return parseInt(str)
	}
	return nil
}

func parseList(str string) ListNode {
	var list []Node
	// start at 1 because first char is known [
	i := 1
	for {
		c := str[i]
		switch c {
		case '[':
			listNode := parseList(str[i:])
			list = append(list, listNode)
			i += listNode.length
		case ']':
			return ListNode{
				nodes:  list,
				length: i + 1,
			}
		case ',':
			i++
		default:
			intNode := parseInt(str[i:])
			list = append(list, intNode)
			i += intNode.length
		}
	}
}

func parseInt(str string) IntNode {
	var x int
	_, err := fmt.Sscanf(str, "%d", &x)
	if err != nil {
		log.Fatal(err)
	}
	return IntNode{
		x:      x,
		length: (x / 10) + 1,
	}
}

func compare(first, second Node) int {
	firstInt, firstIsInt := first.(IntNode)
	secondInt, secondIsInt := second.(IntNode)
	if firstIsInt && secondIsInt {
		if firstInt.x < secondInt.x {
			return -1
		} else if firstInt.x > secondInt.x {
			return 1
		}
		return 0
	}
	if firstIsInt {
		return compare(ListNode{
			nodes: []Node{first},
		}, second)
	}
	if secondIsInt {
		return compare(first, ListNode{
			nodes: []Node{second},
		})
	}
	firstList := first.(ListNode)
	secondList := second.(ListNode)

	// make sure first is shortest, else swap
	if len(firstList.nodes) > len(secondList.nodes) {
		return -compare(second, first)
	}

	for i, fn := range firstList.nodes {
		sn := secondList.nodes[i]
		if order := compare(fn, sn); order != 0 {
			return order
		}
	}

	// everything matches, it's a tie
	if len(firstList.nodes) == len(secondList.nodes) {
		return 0
	}

	// second is longer, it should sort later, it's in order
	return -1
}
