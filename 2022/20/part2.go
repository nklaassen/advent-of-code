package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	decryptionKey = 811589153
)

var input = flag.String("input", "input.txt", "input file")

func main() {
	flag.Parse()
	f, err := os.Open(*input)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var list []Node
	var n Num
	scan := func() error { _, err := fmt.Fscanf(f, "%d\n", &n); return err }
	for err = scan(); err == nil; err = scan() {
		list = append(list, Node{
			val: n * decryptionKey,
		})
	}
	N := len(list)
	var zero *Node
	for i := 1; i < len(list)-1; i++ {
		list[i].prev = &list[i-1]
		list[i].next = &list[i+1]
		if list[i].val == 0 {
			zero = &list[i]
		}
	}
	list[0].next = &list[1]
	list[0].prev = &list[N-1]
	list[N-1].prev = &list[N-2]
	list[N-1].next = &list[0]

	for i := 0; i < 10; i++ {
		fmt.Println("mix", i)
		mix(list)
	}
	fmt.Println(groveSum(zero))
}

func groveSum(zero *Node) Num {
	sum := Num(0)
	curr := zero.next
	for i := 1; i <= 3000; i++ {
		if i%1000 == 0 {
			fmt.Println(curr.val)
			sum += curr.val
		}
		curr = curr.next
	}
	return sum
}

func printList(head *Node) {
	fmt.Print(head.val, " ")
	for curr := head.next; curr != head; curr = curr.next {
		fmt.Print(curr.val, " ")
	}
	fmt.Println()
	fmt.Print(head.val, " ")
	for curr := head.prev; curr != head; curr = curr.prev {
		fmt.Print(curr.val, " ")
	}
	fmt.Println()
}

func mix(list []Node) {
	N := Num(len(list))
	for i := range list {
		//printList(&list[0])
		val := list[i].val
		//fmt.Println("moving", val)
		if val > 0 {
			steps := val % (N - 1)
			moveForward(&list[i], steps)
		}
		if val < 0 {
			steps := (-val) % (N - 1)
			moveBackward(&list[i], steps)
		}
	}
	//printList(&list[0])
}

func moveForward(node *Node, steps Num) {
	node.prev.next = node.next
	node.next.prev = node.prev
	dest := node.next
	for i := steps; i > 0; i-- {
		dest = dest.next
	}
	node.prev = dest.prev
	node.next = dest
	dest.prev.next = node
	dest.prev = node
}

func moveBackward(node *Node, steps Num) {
	node.prev.next = node.next
	node.next.prev = node.prev
	dest := node.prev
	for i := steps; i > 0; i-- {
		dest = dest.prev
	}
	node.prev = dest
	node.next = dest.next
	dest.next.prev = node
	dest.next = node
}

type Num int

type Node struct {
	prev, next *Node
	val        Num
}
