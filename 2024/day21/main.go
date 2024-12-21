package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

var N_ROBOTS = 25

var NUM_PAD = [][]byte{
	{'#', '0', 'A'},
	{'1', '2', '3'},
	{'4', '5', '6'},
	{'7', '8', '9'},
}
var NUM_START = Vec2{0, 2}

var DIR_PAD = [][]byte{
	{'<', 'v', '>'},
	{'#', '^', 'A'},
}
var DIR_START = Vec2{1, 2}

var TEST = [][]byte{ /////// expected length
	{'0', '2', '9', 'A'}, // 68
	{'9', '8', '0', 'A'}, // 60
	{'1', '7', '9', 'A'}, // 68
	{'4', '5', '6', 'A'}, // 64
	{'3', '7', '9', 'A'}, // 64
}

var TEST2 = [][]byte{ ////// expected length
	{'1', '5', '9', 'A'}, // 82
	{'3', '7', '5', 'A'}, // 70
	{'6', '1', '3', 'A'}, // 62
	{'8', '9', '4', 'A'}, // 78
	{'0', '8', '0', 'A'}, // 60
}

var INPUT = [][]byte{
	{'3', '4', '1', 'A'},
	{'0', '8', '3', 'A'},
	{'8', '0', '2', 'A'},
	{'9', '7', '3', 'A'},
	{'7', '8', '0', 'A'},
}

type Dir struct {
	sym byte
	dir Vec2
}

var DIRS = []Dir{
	{'<', Vec2{0, -1}},
	{'v', Vec2{1, 0}},
	{'^', Vec2{-1, 0}},
	{'>', Vec2{0, 1}},
}

type Vec2 struct {
	y, x int
}

func (v Vec2) Add(w Vec2) Vec2 {
	return Vec2{v.y + w.y, v.x + w.x}
}

type PadType int

const (
	NumPad PadType = iota
	DirPad
)

func main() {
	codes := INPUT

	total := 0
	for _, code := range codes {
		seq := findSequence(NUM_PAD, NumPad, code, NUM_START)
		l := countSequences(DIR_PAD, seq, 1)
		n, _ := strconv.Atoi(string(code[:len(code)-1]))
		total += n * l
	}
	fmt.Printf("Expected: 248566068436630\n")
	fmt.Printf("Actual  : %d\n", total)
}

var cache = make(map[string][]int)
func countSequences(pad [][]byte, code []byte, robot int) int {
	key := string(code)
	if val, ok := cache[key]; ok && robot <= len(val) && val[robot-1] != 0 {
		return val[robot-1]
	}
	if _, ok := cache[key]; !ok {
		cache[key] = make([]int, N_ROBOTS)
	}

	seq := findSequence(pad, DirPad, code, DIR_START)
	if robot == N_ROBOTS {
		return len(seq)
	}

	steps := splitSequence(seq)
	count := 0
	for _, step := range steps {
		c := countSequences(pad, step, robot+1)
		count += c
	}
	cache[key][robot-1] = count
	return count
}


func findSequence(pad [][]byte, padtype PadType, code []byte, start Vec2) []byte {
	var out []byte
	pos := start

	for _, point := range code {
		targ := target(pad, point)
		dx, dy := targ.x-pos.x, targ.y-pos.y
		h, v := []byte{}, []byte{}

		for i := 0; i < abs(dx); i++ {
			if dx >= 0 {
				h = append(h, '>')
			} else {
				h = append(h, '<')
			}
		}

		for i := 0; i < abs(dy); i++ {
			if dy >= 0 {
				v = append(v, '^')
			} else {
				v = append(v, 'v')
			}
		}

		switch padtype {
		case NumPad:
			if pos.y == 0 && targ.x == 0 {
				out = append(out, v...)
				out = append(out, h...)
			} else if (pos.x == 0 && targ.y == 0) || dx < 0 {
				out = append(out, h...)
				out = append(out, v...)
			} else {
				out = append(out, v...)
				out = append(out, h...)
			}
		case DirPad:
			if pos.x == 0 && targ.y == 1 {
				out = append(out, h...)
				out = append(out, v...)
			} else if pos.y == 1 && targ.x == 0 {
				out = append(out, v...)
				out = append(out, h...)
			} else if dx < 0 {
				out = append(out, h...)
				out = append(out, v...)
			} else {
				out = append(out, v...)
				out = append(out, h...)
			}
		default:
			panic(fmt.Sprintf("findSequence: invalid padtype %d", padtype))
		}
		pos = targ
		out = append(out, 'A')
	}

	return out
}

func splitSequence(code []byte) [][]byte {
	var result [][]byte
	var current []byte
	for _, point := range code {
		current = append(current, point)
		if point == 'A' {
			result = append(result, current)
			current = []byte{}
		}
	}
	return result
}

func target(pad [][]byte, code byte) Vec2 {
	rows := len(pad)
	cols := len(pad[0])
	var targ Vec2
	for j := range rows {
		for i := range cols {
			if pad[j][i] == code {
				targ = Vec2{j, i}
			}
		}
	}
	return targ
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

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
