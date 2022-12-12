package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
)

const (
	start byte = 'S'
	end        = 'E'

	width      int = 179
	height     int = 41
	totalNodes     = width * height
)

var offsets = []int{
	-width,
	-1,
	1,
	width,
}

func main() {
	f, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	state := &State{
		unvisited: &Nodes{},
		edges:     make(map[int][]int),
		visited:   make(map[int]bool),
	}

	var start, end *Node

	lines := bytes.Split(f, []byte{'\n'})
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			index := i*width + j
			b := lines[i][j]
			elevation := int(b - 'a')
			node := newNode(index, elevation)
			state.unvisited.Push(node)
			state.nodes = append(state.nodes, node)
			switch b {
			case 'S':
				start = node
				start.elevation = 0
				start.distance = 0
			case 'E':
				end = node
				end.elevation = int('z' - 'a')
			}
		}
	}

	for _, node := range state.nodes {
		for _, o := range offsets {
			// avoid out of range indices
			targetIndex := node.index + o
			if targetIndex < 0 || targetIndex >= totalNodes {
				continue
			}
			// no edge to elevations greater than 1 higher
			if state.nodes[targetIndex].elevation-node.elevation > 1 {
				continue
			}
			state.edges[node.index] = append(state.edges[node.index], targetIndex)
		}
	}

	heap.Init(state.unvisited)
	dijkstra(state, end.index)
	fmt.Println(end.distance)
}

func dijkstra(state *State, end int) {
	for {
		curr := heap.Pop(state.unvisited).(*Node)
		fmt.Println("Visiting node", curr.index)
		fmt.Println("  current distance", curr.distance)
		state.visited[curr.index] = true
		if curr.index == end {
			fmt.Println("  found end")
			return
		}
		for _, e := range state.edges[curr.index] {
			fmt.Println("  edge to node", e)
			if state.visited[e] {
				fmt.Println("    edge already visited")
				continue
			}
			n := state.nodes[e]
			fmt.Println("    edge current distance", n.distance)
			n.distance = min(n.distance, curr.distance+1)
			fmt.Println("    edge new distance", n.distance)
			heap.Fix(state.unvisited, n.heapIndex)
		}
		fmt.Println()
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func newNode(index int, elevation int) *Node {
	return &Node{
		elevation: elevation,
		index:     index,
		distance:  math.MaxInt32,
	}
}

type State struct {
	nodes     []*Node
	unvisited *Nodes
	edges     map[int][]int
	visited   map[int]bool
}

var _ heap.Interface = (*Nodes)(nil)

type Nodes struct {
	nodes []*Node
}

func (n *Nodes) Len() int {
	return len(n.nodes)
}

func (n *Nodes) Less(i, j int) bool {
	return n.nodes[i].distance < n.nodes[j].distance
}

func (n *Nodes) Swap(i, j int) {
	n.nodes[i], n.nodes[j] = n.nodes[j], n.nodes[i]
	n.nodes[i].heapIndex, n.nodes[j].heapIndex = n.nodes[j].heapIndex, n.nodes[i].heapIndex
}

func (n *Nodes) Push(x any) {
	node := x.(*Node)
	node.heapIndex = len(n.nodes)
	n.nodes = append(n.nodes, node)
}

func (n *Nodes) Pop() any {
	x := n.nodes[len(n.nodes)-1]
	n.nodes = n.nodes[:len(n.nodes)-1]
	return x
}

type Node struct {
	distance  int
	elevation int
	index     int
	heapIndex int
}
