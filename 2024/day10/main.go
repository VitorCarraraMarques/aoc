package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

type Position struct {
    Row int
    Col int
}

type Path []Position

const PART = 2

var ROWS int
var COLS int

func main() {
    mt := readMapTop()
    res := 0
    for j := 0; j < ROWS; j++ {
        for i := 0; i < COLS; i++ {
            if mt[j][i] != 0 {
                continue
            }
            cur := Position{j, i}
            path := Path{cur}
            trails := trail(mt, path)
            var n int 
            if PART == 1 {
                n = score(trails)
            } else if PART == 2 {
                n = rating(trails)
            }
            // printScore(j, i, mt, trails, n)
            res += n 
        }
    }
    fmt.Printf("Part %d: %d\n", PART, res)
}

func rating(paths []Path) int {
    return len(paths)
}

func score(paths []Path) int {
    ends := make(map[Position]bool)
    for _, p := range paths {
        last := p[len(p)-1]
        ends[last] = true
    }
    sc := 0
    for range ends {
        sc++
    }
    return sc
}

func trail(maptop [][]int, path Path) []Path {
    assert(len(path) > 0, "trail: empty path")

    var paths []Path

    last := path[len(path)-1]
    j := last.Row
    i := last.Col
    cur := maptop[j][i]

    if cur == 9 {
        return []Path{path}
    }

    dirs := [][]int{
        { 0,  1},
        { 1,  0},
        { 0, -1},
        {-1,  0},
    }

    for _, dir := range dirs {
        row := j + dir[0]
        col := i + dir[1]

        if !inbounds(row, col) {
            continue
        }

        nbor_pos := Position{Row: row, Col: col}
        nbor_val := maptop[row][col]
        if nbor_val != cur + 1 {
            continue
        }

        if visited(path, row, col) {
            continue
        }

        next := make(Path, len(path))
        copy(next, path)
        next = append(next, nbor_pos)
        paths = append(paths, trail(maptop, next)...)
    }
    return paths
}

func visited(path Path, row, col int) bool {
    for _, p := range path {
        if p.Row == row && p.Col == col {
            return true
        }
    }
    return false
}

func inbounds(j, i int) bool {
    return j >= 0 && j < ROWS && i >= 0 && i < COLS
}

func printScore(j,i int, mt [][]int, trails []Path, n int) {
    fmt.Printf("====== (ROW: %d, COL: %d) ========\n", j, i)
    for k, p := range trails {
        fmt.Printf("   Path #%d\n", k)
        printPath(mt, p)
    }
    fmt.Printf("  > Score: %d\n", n)
}


func printPath(m [][]int, path Path) {
    for _, p := range path {
        fmt.Printf("    (Row: %d, Col: %d) => %d\n", p.Row, p.Col, m[p.Row][p.Col])
    }
}

func readMapTop() [][]int{
    file := readFile("input.txt")

    lines := bytes.Split(file, []byte{'\n'})
    lines = lines[:len(lines)-1]

    ROWS = len(lines)
    COLS = len(lines[0])

    maptop := make([][]int, len(lines))
    for j, line := range lines {
        row := make([]int, len(line))
        for i, letter := range line {
            n, _ := strconv.Atoi(string(letter))
            row[i] = n
        }
        maptop[j] = row
    }
    return maptop
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
