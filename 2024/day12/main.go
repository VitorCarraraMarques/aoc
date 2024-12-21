package main

import (
	"bytes"
	"fmt"
	"os"
)

const PART = 2

type Position struct {
	Row int
	Col int
}

func (p Position) Sub(q Position) (int, int) {
	return p.Row - q.Row, p.Col - q.Col
}

type Path map[Position]bool

type Region struct {
	Path
	Area      int
	Perimeter int
	Border    Path
}

var DIRS = [4][2]int{
	{0, 1},
	{1, 0},
	{0, -1},
	{-1, 0},
}

var DIAGS = [4][2]int{
	{-1, 1},
	{1, 1},
	{1, -1},
	{-1, -1},
}

var ROWS int
var COLS int

func main() {
	farm := readMap()
	m := make(map[Position]*Region)
	regions := make(map[*Region]bool)
	for j := 0; j < ROWS; j++ {
		for i := 0; i < COLS; i++ {
			if _, ok := m[Position{j, i}]; !ok {
				path := make(Path)
				border := make(Path)
				reg := &Region{
					Path:   path,
					Border: border,
				}
				regions[reg] = true
				m[Position{j, i}] = reg
				flood(farm, j, i, reg, m)
			}
		}
	}

	total := 0
	for reg := range regions {
		p := price(reg)
		total += p
	}

	fmt.Printf("\nTOTAL FENCE PRICE: $%d.00\n", total)

}

func flood(farm [][]rune, j, i int, reg *Region, m map[Position]*Region) {
	if reg.Path[Position{Row: j, Col: i}] {
		return
	}

	reg.Path[Position{j, i}] = true
	m[Position{j, i}] = reg
	reg.Area = len(reg.Path)

	cur := farm[j][i]
	for _, dir := range DIRS {
		nj := j + dir[0]
		ni := i + dir[1]
		if !inbound(nj, ni) {
			reg.Perimeter++
			reg.Border[Position{nj, ni}] = true
			continue
		}

		nbor := farm[nj][ni]
		if nbor != cur {
			reg.Perimeter++
			reg.Border[Position{nj, ni}] = true
		} else {
			flood(farm, nj, ni, reg, m)
		}
	}
}

func price(r *Region) int {
	if PART == 1 {
		return r.Perimeter * r.Area
	} else if PART == 2 {
		return sides(r) * r.Area
	}
	panic("unreachable")
}

func sides(r *Region) int {
	v := make(map[Position]bool)
	sides := 0
	for c := range r.Path {
		if v[c] {
			continue
		}
		v[c] = true
		var outs []Position
		for _, dir := range DIRS {
			nbor := Position{
				Row: c.Row + dir[0],
				Col: c.Col + dir[1],
			}
			if r.Border[nbor] {
				outs = append(outs, nbor)
			}
		}

		n := 0
		switch len(outs) {
		case 0:
			ds := 0
			for _, diag := range DIAGS {
				d := Position{c.Row + diag[0], c.Col + diag[1]}
				if r.Border[d] {
					ds++
				}
			}
			n += ds
		case 1:
			dj, di := c.Sub(outs[0])
			if di == 0 {
				if r.Border[Position{c.Row + dj, c.Col - 1}] {
					n++
				}
				if r.Border[Position{c.Row + dj, c.Col + 1}] {
					n++
				}
			}
			if dj == 0 {
				if r.Border[Position{c.Row - 1, c.Col + di}] {
					n++
				}
				if r.Border[Position{c.Row + 1, c.Col + di}] {
					n++
				}
			}
		case 2:
			if outs[0].Row != outs[1].Row && outs[0].Col != outs[1].Col {
				n++

				dj0, di0 := outs[0].Sub(c)
				dj1, di1 := outs[1].Sub(c)
				if r.Border[Position{c.Row-(dj0 + dj1), c.Col-(di0 + di1)}] {
					n++
				}
			}
		case 3:
			n += 2
		case 4:
			n += 4
        }
		sides += n

	}
	return sides
}

func markedge(r *Region, v map[Position]bool, c Position) {
	v[c] = true
	for _, dir := range DIRS {
		nbor := Position{
			Row: c.Row + dir[0],
			Col: c.Col + dir[1],
		}
		if r.Border[nbor] && !v[nbor] {
			markedge(r, v, nbor)
		}
	}
}

func inbound(j, i int) bool {
	return j >= 0 && j < ROWS && i >= 0 && i < COLS
}

func readMap() [][]rune {
	file := readFile("input.txt")
	lines := bytes.Split(file, []byte{'\n'})
	lines = lines[:len(lines)-1]

	ROWS = len(lines)
	COLS = len(lines[0])

	res := make([][]rune, len(lines))
	for j, l := range lines {
		row := make([]rune, len(l))
		for i, w := range l {
			row[i] = rune(w)
		}
		res[j] = row
	}
	return res
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
