package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// const ROWS = 7
// const COLS = 7

const ROWS = 71
const COLS = 71

var DIRS = [][]int{
	{0, 1},
	{1, 0},
	{0, -1},
	{-1, 0},
}

type Vec2 struct {
	x, y int
}

func (v Vec2) Add(w Vec2) Vec2 {
	return Vec2{y: v.y + w.y, x: v.x + w.x}
}

func (v Vec2) Idx() int {
	return v.y*COLS + v.x
}

func (v Vec2) InBound() bool {
	return v.x >= 0 && v.x < COLS && v.y >= 0 && v.y < ROWS
}

func (v Vec2) Nbors() []Vec2 {
	var nbors []Vec2
	for _, dir := range DIRS {
		n := Vec2{
			x: v.x + dir[1],
			y: v.y + dir[0],
		}
		if n.InBound() {
			nbors = append(nbors, n)
		}
	}
	return nbors
}

func main() {
	g := make([]byte, ROWS*COLS)

    s := 1025

    fb := fallingbytes("input.txt")
    for i := range s {
        gridbyte(g, fb[i])
    }

    start := Vec2{0, 0}
    goal := Vec2{x: COLS - 1, y: ROWS - 1}
    for i := s; i < len(fb); i++ {
        pos := gridbyte(g, fb[i])
        _, err := bfs(g, start, goal)
        if err != nil {
            fmt.Printf("First Byte to Block Path in Position %v\n", pos)
            break
        }
    }
}

type Item struct {
	Vec2
	path []Vec2
}

type Queue []Item

func (q Queue) Add(b Item) Queue {
	q = append(q, b)
	return q
}

func (q Queue) PopRight() (Queue, Item) {
	b := q[len(q)-1]
	q = q[:len(q)-1]
	return q, b
}

func (q Queue) PopLeft() (Queue, Item) {
	b := q[0]
	q = q[1:]
	return q, b
}

func (q Queue) Len() int {
	return len(q)
}

func bfs(grid []byte, start Vec2, goal Vec2) ([]Vec2, error) {
	p := []Vec2{start}
	item := Item{start, p}
	q := Queue{item}

	vis := make(map[Vec2]bool)

	for i := 0; q.Len() > 0; i++ {
		q, item = q.PopLeft()

		if item.Vec2 == goal {
			return item.path, nil
		}

		if vis[item.Vec2] {
			continue
		}
		vis[item.Vec2] = true

		if grid[item.Idx()] == '#' {
			continue
		}

		for _, n := range item.Nbors() {
            if vis[n] {
                continue
            }
			np := make([]Vec2, len(item.path)+1)
			copy(np, item.path)
            np[len(np)-1] = n
            
			nbor := Item{n, np}
			q = q.Add(nbor)
		}
	}
    return nil, fmt.Errorf("bfs: no path found")
}

func gridbyte(grid []byte, line []byte) Vec2 {
    split := bytes.Split(line, []byte{','})
    x, _ := strconv.Atoi(string(split[0]))
    y, _ := strconv.Atoi(string(split[1]))
    p := Vec2{x: x, y: y}
    grid[p.Idx()] = '#'
    return p
}



func fallingbytes(name string) [][]byte {

	file := readFile(name)
	lines := bytes.Split(file, []byte{'\n'})
	lines = lines[:len(lines)-1]
    return lines
}


func printGrid(g []byte) {
	for j := range ROWS {
		fmt.Printf("%s\n", g[j*COLS:j*COLS+COLS])
	}
    fmt.Printf("\n")
}

func readFile(name string) []byte {
	file, err := os.ReadFile(name)
	if err != nil {
		err = fmt.Errorf("open file %s: %s\n", name, err)
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
	return file
}

func assert(expr bool, msg string) {
	if !expr {
		panic(msg)
	}
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
