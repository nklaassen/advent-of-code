package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	maxSize          = 100000
	totalSpace       = 70000000
	totalNeededSpace = 30000000
)

func main() {
	f, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	file := string(f)

	root := newDir("/", nil)
	state := &State{
		root: root,
		curr: root,
	}

	commands := strings.Split(file, "$")
	for _, command := range commands {
		processCommand(state, strings.TrimSpace(command))
	}

	currentSpace := totalSpace - root.size
	fmt.Println("current free space:", currentSpace)
	needSpace := totalNeededSpace - currentSpace
	fmt.Println("need extra space:", needSpace)

	minSizeToDelete := int64(totalNeededSpace)
	dfs(root, func(d *Dir, depth int) {
		if d.size < needSpace {
			// not big enough
			// we could optimize the dfs here but its plenty fast enough
			return
		}
		if d.size < minSizeToDelete {
			minSizeToDelete = d.size
		}
	}, 0)
	fmt.Println("minimum size directory to delete:", minSizeToDelete)
}

func dfs(curr *Dir, f func(*Dir, int), depth int) {
	f(curr, depth)
	for _, child := range curr.children {
		dfs(child, f, depth+1)
	}
}

func processCommand(state *State, commandAndOutput string) {
	lines := strings.Split(commandAndOutput, "\n")
	cmdAndArgs := strings.Split(lines[0], " ")
	cmd := cmdAndArgs[0]
	args := cmdAndArgs[1:]
	output := lines[1:]

	switch cmd {
	case "ls":
		processLs(state, output)
	case "cd":
		processCd(state, args[0])
	}
}

func processLs(state *State, outputLines []string) {
	for _, line := range outputLines {
		parts := strings.Split(line, " ")
		switch parts[0] {
		case "dir":
			dirName := parts[1]
			if _, ok := state.curr.children[dirName]; ok {
				// already know about this dir
				continue
			}
			state.curr.children[dirName] = newDir(dirName, state.curr)
		default:
			sizeStr, fileName := parts[0], parts[1]
			size, err := strconv.Atoi(sizeStr)
			if err != nil {
				log.Fatal(err)
			}
			state.curr.files[fileName] = &File{
				name: fileName,
				size: int64(size),
			}
			curr := state.curr
			for curr != nil {
				curr.size += int64(size)
				curr = curr.parent
			}
		}
	}
}

func processCd(state *State, dir string) {
	switch dir {
	case "/":
		state.curr = state.root
	case "..":
		state.curr = state.curr.parent
	default:
		state.curr = state.curr.children[dir]
	}
}

type State struct {
	root *Dir
	curr *Dir
}

func newDir(name string, parent *Dir) *Dir {
	return &Dir{
		name:     name,
		parent:   parent,
		children: make(map[string]*Dir),
		files:    make(map[string]*File),
	}
}

type Dir struct {
	name     string
	size     int64
	children map[string]*Dir
	files    map[string]*File
	parent   *Dir
}

type File struct {
	name string
	size int64
}
