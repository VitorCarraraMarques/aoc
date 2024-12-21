package main

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

const (
    empty byte = iota
    obstacle
    marked
    guard
    barrel
)

const (
    up byte = iota
    right
    down
    left
)

var dirs = [][]int{
    {-1, 0}, // UP_____^
    {0, 1},  // RIGHT__>
    {1, 0},  // DOWN___v
    {0, -1}, // LEFT___<
}

type Position struct {
    Row int
    Col int
}

func (p *Position) Sub(q Position) (int, int) {
    a := p.Row - q.Row
    if a != 0 {
        a = a / abs(a)
    }
    b := p.Col - q.Col
    if b != 0 {
        b = b / abs(b)
    }
    return a, b
}


type State struct {
    Pos Position
    Dir byte
}

func (s *State) Turn() {
    s.Dir = s.NextDir()
}

func (s *State) NextDir() byte {
    return (s.Dir + 1) % 4
}

func (s *State) Inc() State {
    dir := dirs[s.Dir]
    row := s.Pos.Row + dir[0]
    col := s.Pos.Col + dir[1]
    return State{Pos: Position{row, col}, Dir: s.Dir} 
}

func (s *State) String() string {
    return fmt.Sprintf("{Y: %d, X: %d} {dY: %d, dX: %d}", s.Pos.Row, s.Pos.Col, dirs[s.Dir][0], dirs[s.Dir][1])
}

type Guard struct {
    State
}

type Path []State

func (g *Guard) IsObstacleInFront(m [][]byte) bool {
    x := g.Pos.Row + dirs[g.Dir][0]
    y := g.Pos.Col + dirs[g.Dir][1]
    return m[x][y] == obstacle
}

func (g *Guard) Step(m [][]byte) int {
    isnew := 0
    g.Pos.Row += dirs[g.Dir][0]
    g.Pos.Col += dirs[g.Dir][1]
    cell := m[g.Pos.Row][g.Pos.Col] 
    if cell != marked { 
        m[g.Pos.Row][g.Pos.Col] = marked
        isnew = 1
    }
    return isnew
}

func (g *Guard) IsOoBInFront(m [][]byte) bool {
    x, y := g.Pos.Row+dirs[g.Dir][0], g.Pos.Col+dirs[g.Dir][1]
    return IsOoB(x, y, m)
}


func (g *Guard) Walk(m [][]byte) (Path, int) {
    var steps Path
    visited := make(map[State]bool)
    sum := 0

    for i := 0;;i++ {
        printMap(m, *g)
        time.Sleep(time.Millisecond * 200)
        visited[g.State] = true

        if g.IsOoBInFront(m) {
            break
        }

        if g.IsObstacleInFront(m) {
            g.Turn()
        } else {
            sum += g.Step(m)
            steps = append(steps, g.State)
        }

        if visited[g.State] {
            return Path{}, -1
        }
    }
    return steps, sum
}

func main() {
    m, g := readMap()
    init := g.State
    steps, unique := g.Walk(m)
    usteps := make(map[Position]bool)
    for _, s := range steps {
        usteps[s.Pos] = true
    }
    barrels := make(map[Position]bool)
    for step := range usteps {
        var loop [][]byte
        for _, row := range m {
            cp := make([]byte, len(row))
            copy(cp, row)
            loop = append(loop, cp)
        }
        g.State = init
        row := step.Row
        col := step.Col
        loop[row][col] = obstacle
        _, isLoop := g.Walk(loop)
        if isLoop == -1 {
            pos := Position{row, col}
            barrels[pos] = true
        }
    }
    fmt.Printf("Part One: %d\n", unique)
    fmt.Printf("Part Two: %d\n", len(barrels))
}


func IsOoB(row, col int, m [][]byte) bool {
    w, h := len(m[0]), len(m)
    return (col < 0 || col >= w || row < 0 || row >= h)
}

func printMap(m [][]byte, g Guard) {
    ROWS := len(m)
    COLS := len(m[0])
    for _, row := range m {
        for _, cell := range row {
            printCell(cell, g)
        }
        fmt.Printf("\n")
    }
    fmt.Printf("\033[%dA\033[%dD", ROWS, COLS)
}

func printCell(cell byte, g Guard) {
    switch cell {
    case empty:
        fmt.Printf("%s", ".")
    case obstacle:
        fmt.Printf("%s", "#")
    case marked:
        fmt.Printf("%s", "X")
    case guard:
        printGuard(g)
    }
}

func printGuard(g Guard) {
    switch g.Dir {
    case up:
        fmt.Printf("%s", "^")
    case right:
        fmt.Printf("%s", ">")
    case down:
        fmt.Printf("%s", "v")
    case left:
        fmt.Printf("%s", "<")
    }
}

func readMap() ([][]byte, Guard) {
    file := readFile("test.txt")
    lines := bytes.Split(file, []byte{'\n'})
    lines = lines[:len(lines)-1]

    mapp := make([][]byte, len(lines))
    gd := Guard{}

    for i, line := range lines {
        row := make([]byte, len(line))
        mapp[i] = row
        for j, letter := range line {
            if letter == '.' {
                mapp[i][j] = empty
            } else if letter == '#' {
                mapp[i][j] = obstacle
            } else if letter == '^' {
                mapp[i][j] = guard
                gd.Pos = Position{Row: i, Col: j}
                gd.Dir = up
            } else if letter == '>' {
                mapp[i][j] = guard
                gd.Pos = Position{Row: i, Col: j}
                gd.Dir = right
            } else if letter == 'v' {
                mapp[i][j] = guard
                gd.Pos = Position{Row: i, Col: j}
                gd.Dir = down
            } else if letter == '<' {
                mapp[i][j] = guard
                gd.Pos = Position{Row: i, Col: j}
                gd.Dir = left
            }
        }
    }
    return mapp, gd
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

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}
