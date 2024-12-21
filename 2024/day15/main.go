package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"
)

const PART = 2

var display []byte
var robot Vec2

var ROWS int
var COLS int

var DIRS = map[rune]Vec2{
    '>': {y: 0, x: 1},
    'v': {y: 1, x: 0},
	'<': {y: 0, x: -1},
	'^': {y: -1, x: 0},
}

const (
	blocked = iota
	moved
)

type Box struct {
	left  Vec2
	right Vec2
}

func (b Box) Move(dir Vec2) {
	old_left := display[b.left.Idx()]
	old_right := display[b.right.Idx()]
	left_nbor := b.left.Add(dir).Idx()
	right_nbor := b.right.Add(dir).Idx()
	display[b.left.Idx()] = '.'
	display[b.right.Idx()] = '.'
	display[left_nbor] = old_left
	display[right_nbor] = old_right
}

type Vec2 struct {
	y int
	x int
}

func (v Vec2) Add(w Vec2) Vec2 {
	return Vec2{y: v.y + w.y, x: v.x + w.x}
}

func (v Vec2) Idx() int {
	return v.y*COLS + v.x
}

func (v Vec2) Left() Vec2 {
    return Vec2{y:v.y, x:(v.x-1)}
}

func (v Vec2) Right() Vec2 {
    return Vec2{y:v.y, x:(v.x+1)}
}


func main() {
	cmds := readGrid()
    if PART == 2 {
        resizeGrid()
    }
	for ip := 0; ip < len(cmds); ip++ {
		instr := rune(cmds[ip])
		fmt.Printf("Move %c\n", instr)
		checkNbor(&robot, DIRS[instr])
		printGrid()
        time.Sleep(10*time.Millisecond)
        fmt.Printf("\033[0;0H")
        clear()
	}
	sum := gps()
	fmt.Printf("Sum of GPS: %d\n", sum)
}

func gps() int {
	total := 0
	for j := range ROWS {
		for i := range COLS {
			if display[j*COLS+i] == 'O' || display[j*COLS+i] == '[' {
				total += (100 * j) + i
			}
		}
	}
	return total
}


func checkNbor(cur *Vec2, dir Vec2) int {
    n := cur.Add(dir)
	nbor := display[n.Idx()]
	switch nbor {
	case '#':
		return blocked
	case '.':
		old := display[cur.Idx()]
		display[cur.Idx()] = '.'
		cur.y += dir.y
		cur.x += dir.x
		display[cur.Idx()] = old
		return moved
	case 'O':
		res := checkNbor(&n, dir)
		if res == blocked {
			return blocked
		} else {
			return checkNbor(cur, dir)
		}
	case '[':
        box := Box{left: n, right: n.Right()}
        res := checkBox(box, dir)
        if res == blocked {
			return blocked
		} else {
			return checkNbor(cur, dir)
		}
    case ']':
        box := Box{left: n.Left(), right: n}
        res := checkBox(box, dir)
        if res == blocked {
			return blocked
		} else {
			return checkNbor(cur, dir)
		}

	default:
		panic(fmt.Sprintf("checkNbor: invalid nbor %c", nbor))
	}
}

type Queue []Box

func (q Queue) Add(b Box) Queue {
    q = append(q, b)
    return q
} 

func (q Queue) PopRight() (Queue, Box) {
    b := q[len(q)-1]
    q = q[:len(q)-1]
    return q, b
}

func (q Queue) PopLeft() (Queue, Box) {
    b := q[0]
    q = q[1:]
    return q, b
}


func checkBox(b Box, d Vec2) int {
    queue := Queue{b}
    stack := Queue{b}

    visited := make(map[Box]bool)
    for len(queue) > 0 {

        var cur Box
        queue, cur = queue.PopLeft()

        if visited[cur] {
            continue
        }
        visited[cur] = true

        ln := cur.left.Add(d)
        rn := cur.right.Add(d)

        if display[ln.Idx()] == '#' || display[rn.Idx()] == '#' {
            return blocked
        }
        
        if ln != cur.left && ln != cur.right {
            if display[ln.Idx()] == '[' {
                box := Box{left: ln, right: ln.Right()}
                queue = queue.Add(box)
                stack = stack.Add(box)
            }
            if display[ln.Idx()] == ']' {
                box := Box{left: ln.Left(), right: ln}
                queue = queue.Add(box)
                stack = stack.Add(box)
            }
        }


        if rn != cur.left && rn != cur.right {
            if display[rn.Idx()] == '[' {
                box := Box{left: rn, right: rn.Right()}
                queue = queue.Add(box)
                stack = stack.Add(box)
            }
            if display[rn.Idx()] == ']' {
                box := Box{left: rn.Left(), right: rn}
                queue = queue.Add(box)
                stack = stack.Add(box)
            }
        }
    }

    visited = make(map[Box]bool)
    for len(stack) > 0 {
        var cur Box
        stack, cur = stack.PopRight()
        if visited[cur] {
            continue
        }
        visited[cur] = true
        cur.Move(d)
    }

    return moved
}

func resizeGrid() []byte {
	new_display := make([]byte, ROWS*COLS*2)
	for j := range ROWS {
		for i := range COLS {
			idx := j*COLS + i
			old := display[idx]
			new_idx := j*2*COLS + 2*i
			switch old {
			case '@':
				new_display[new_idx] = '@'
				new_display[new_idx+1] = '.'
			case '#':
				new_display[new_idx] = '#'
				new_display[new_idx+1] = '#'
			case '.':
				new_display[new_idx] = '.'
				new_display[new_idx+1] = '.'
			case 'O':
				new_display[new_idx] = '['
				new_display[new_idx+1] = ']'
			}
		}
	}
	COLS *= 2
	display = new_display
	robot.x *= 2
	return new_display
}

func readGrid() []byte {
	file := readFile("input.txt")
	split := bytes.Split(file, []byte{'\n', '\n'})
	assert(len(split) == 2, "readGrid: too many parts in split")
	grid, cmds := split[0], split[1]

	rows := bytes.Split(grid, []byte{'\n'})
	ROWS = len(rows)
	COLS = len(rows[0])
	display = make([]byte, ROWS*COLS)

	for j, row := range rows {
		for i, cell := range row {
			display[j*COLS+i] = cell
			if cell == '@' {
				robot.y = j
				robot.x = i
			}
		}
	}

	cmds = bytes.ReplaceAll(cmds, []byte{'\n'}, []byte{})

	return cmds
}

func printGrid() {
	for j := range ROWS {
		s := j * COLS
		e := j*COLS + COLS
		fmt.Printf("%s\n", display[s:e])
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
