package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

var GRID []byte

var ROWS int
var COLS int

var DIRS = []Vec2{
    { 0,  1},
    { 1,  0},
    { 0, -1},
    {-1,  0},
}


type Vec2 struct {
    y, x int
}

func (v Vec2) Add(w Vec2) Vec2 {
    return Vec2{v.y+w.y, v.x+w.x}
}

func (v Vec2) Sub(w Vec2) Vec2 {
    return Vec2{v.y-w.y, v.x-w.x}
}

func (v Vec2) Idx() int {
    return v.y*COLS + v.x
}

func (v Vec2) InBound() bool {
    return v.y >= 0 && v.y < ROWS && v.x >= 0 && v.x < COLS
}

func (v Vec2) Nbors() []Vec2{
    var nbors []Vec2
    for _, dir := range DIRS {
        nbor := v.Add(dir)    
        if nbor.InBound() {
            nbors = append(nbors, nbor)
        }
    }
    return nbors
}

type Nbor struct {
    Vec2
    count int
}

func (v Vec2) Nborhd(n int) []Nbor {
    var nborhd []Nbor
    q := &Queue{Nbor{v, 0}}
    visited := make(map[Vec2]bool)
    for q.Len() > 0 {
        cur := q.PopLeft().(Nbor)
        if cur.count > n {
            continue
        }
        if visited[cur.Vec2] {
            continue
        }
        visited[cur.Vec2] = true
        nborhd = append(nborhd, cur)
        nbors := cur.Nbors()
        for _, n := range nbors {
            if !n.InBound() {
                continue
            }
            nc := cur.count + 1
            q.Add(Nbor{n, nc})
        }
    }
    return nborhd
}

type Queue []any

func (q *Queue) Len() int {
    return len(*q)
}

func (q *Queue) Add(i any) {
    queue := append(*q, i)
    *q = queue
}

func (q *Queue) PopLeft() any {
    item := (*q)[0]
    *q = (*q)[1:] 
    return item
}


func main() {
    start, end := grid()
    base := bfs(start, end)
    mapath := toMap(base.path)
    //printGrid()

    // TEST INPUT
    //part1 := countCheats(2, 50, mapath)
    //part2 := countCheats(20, 50, mapath)
    // fmt.Printf("test1 expected PART1 : 1\n")
    // fmt.Printf("test1 actual   PART1 : %d\n", part1)
    // fmt.Printf("test1 expected PART2 : 285\n")
    // fmt.Printf("test1 actual   PART2 : %d\n", part2)

    // ACTUAL INPUT
    part1 := countCheats(2, 100, mapath)
    part2 := countCheats(20, 100, mapath)
    fmt.Printf("input expected PART 1: 1365\n")
    fmt.Printf("input actual PART 1: %d\n", part1)
    fmt.Printf("input actual PART 2 first attempt: 985803\n") 
    fmt.Printf("input actual PART 2: %d\n", part2)
}

type Item struct {
    Vec2
    count int
    path []Vec2
}

func countCheats(steps, tresh int, mapath map[Vec2]int) int {
    total := 0
    for start, idx := range mapath {
        for _, end := range start.Nborhd(steps) {
            n_idx, ok := mapath[end.Vec2] 
            if !ok {
                continue
            }
            if n_idx > idx {
                diff := n_idx - idx
                save := diff - end.count
                if save >= tresh {
                    total++                        
                }
            }
        }
    }
    return total
}

func toMap(path []Vec2) map[Vec2]int {
    mapath := make(map[Vec2]int)
    for i, p := range path {
        mapath[p] = i
    }
    return mapath
}

func bfs(start, end Vec2) Item {
    first := Item{start, 0, []Vec2{start}} 
    q := &Queue{first}
    visited := make(map[Vec2]bool)
    for q.Len() > 0 {
        item := q.PopLeft().(Item)

        if item.Vec2 == end {
            return item
        }

        if visited[item.Vec2]{
            continue
        }
        visited[item.Vec2] = true

        for _, nbor := range item.Nbors() {
            if GRID[nbor.Idx()] == '#' {
                continue
            }
            if visited[nbor]{
                continue
            }
            nc := item.count + 1
            np := copyPath(item.path)
            np[len(np)-1] = nbor
            q.Add(Item{nbor, nc, np})
        }
    }
    panic("bfs: no path found")
}

func copyPath(path []Vec2) []Vec2 {
    cp := make([]Vec2, len(path)+1)
    copy(cp, path)
    return cp
}


func grid() (start Vec2, end Vec2) {
    file := readFile("input.txt")
    lines := bytes.Split(file, []byte{'\n'})
    lines = lines[:len(lines)-1]

    ROWS = len(lines)
    COLS = len(lines[0])
    GRID = make([]byte, ROWS*COLS)
    for j, row := range lines {
        for i, cell := range row {
            idx := j * COLS + i
            GRID[idx] = cell
            if cell == 'S' {
                start = Vec2{j, i}
            }
            if cell == 'E' {
                end = Vec2{j, i}
            }
        }
    }
    return start, end
}

func printGrid() {
    for j := range ROWS {
        fmt.Printf("%s\n", GRID[j*COLS:j*COLS + COLS])
    }
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
