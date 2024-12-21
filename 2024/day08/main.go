package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
)

type Position struct {
	Row int
	Col int
}

type Antenna struct {
	Position
	Frequency rune
}


const PART = 2

var ROWS int
var COLS int

func main() {
	mp := readMap()
	anti := make(map[Position]bool)
	for _, as := range mp {
		pos := perms(as)
		for _, p := range pos {
			anti[p] = true
		}
	}
    fmt.Printf("ROWS = %d, COLS = %d\n", ROWS, COLS)
    fmt.Printf("PART %d: There is %d unique antinode positions\n", PART, len(anti))
}

func perms(ants []Antenna) []Position {
	N := len(ants)
	var pos []Position
	for i := 0; i < N; i++ {
		for j := i + 1; j < N; j++ {
			p0 := ants[i].Position
            p1 := ants[j].Position
            if PART == 1 {
                a, b := equidist(p0, p1)
                if inBound(a) { pos = append(pos, a) }
                if inBound(b) { pos = append(pos, b) }
            } else if PART == 2 {
                l := inline(p0, p1)
                pos = append(pos, l...)}
		}
	}
	return pos
}

func equidist(p0, p1 Position) (Position, Position) {
    x0, y0 := p0.Col, p0.Row
    x1, y1 := p1.Col, p1.Row
    dy := (y1 - y0)
    dx := (x1 - x0)

    ax := x0 + 2*dx 
    ay := y0 + 2*dy 
    a :=  Position{ay, ax}
    
    bx := x0 - dx
    by := y0 - dy
    b := Position{by, bx}

    return a, b
}

func inline(p0, p1 Position) []Position {
    var res []Position

    x0, y0 := p0.Col, p0.Row
    x1, y1 := p1.Col, p1.Row

    dy := (y1 - y0)
    dx := (x1 - x0)
    m  := float64(dy) / float64(dx)
    b  := (float64(y0) - m*float64(x0))

    for x := 0.0; x < float64(COLS); x++ {
        y := m*x + b
        i, d := math.Modf(y)
        if fEqual(d, 0.0) {
            p := Position{
                Col: int(x), Row: int(i),
            }
            if inBound(p) {
                res = append(res, p)
            }         
        }
        if fEqual(d, 1.0) {
            p := Position{
                Col: int(x), Row: int(i) + 1,
            }
            if inBound(p) {
                res = append(res, p)
            }         
        }
    }

    if len(res) == 0 {
        fmt.Printf("------\n")
        fmt.Printf("P0       : %+v\n", p0)
        fmt.Printf("           InBound ? %t\n\n", inBound(p0))

        fmt.Printf("P1       : %+v\n", p1)
        fmt.Printf("           InBound ? %t\n\n", inBound(p1))

        fmt.Printf("Equation : y = %.1f * x + %.1f\n\n", m, b)

        y := m*float64(x0)+b
        fmt.Printf("f(p0)    : %.5f = %.5f * %d + %.5f\n", y, m, x0, b)
        i, d := math.Modf(y)
        fmt.Printf("            - y.Int = %.10f, y.Dec = %.10f\n", i, d)
        fmt.Printf("            y.Dec == 0.0 ? %t\n\n", fEqual(d, 0.0))

        y = m*float64(x0)+b
        fmt.Printf("f(p1)    : %.5f = %.5f * %d + %.5f\n", y, m, x1, b)
        i, d = math.Modf(y)
        fmt.Printf("            - y.Int = %.10f, y.Dec = %.10f\n", i, d)
        fmt.Printf("            y.Dec == 0.0 ? %t\n", fEqual(d, 0.0))
    }

    return res
}

func printAnti(anti map[Position]bool) {
    for j := 0; j < ROWS; j++ {
        for i := 0; i < COLS; i++ {
            cur := Position{Row: j, Col: i} 
            if anti[cur] {
                fmt.Printf("#")
            } else {
                fmt.Printf(".")
            }
        } 
        fmt.Printf("\n")
    }
    fmt.Printf("\n")
}

func fEqual(a, b float64) bool {
    diff := math.Abs(a - b)
    eps := 1e-6
    return diff <= eps
}

func inBound(p Position) bool {
    x, y := p.Col, p.Row
	return (x >= 0 && x < COLS) && (y >= 0 && y < ROWS)
}

func readMap() map[rune][]Antenna {
	file := readFile("input.txt")
	lines := bytes.Split(file, []byte{'\n'})
	lines = lines[:len(lines)-1]
	ROWS = len(lines)
	COLS = len(lines[0])

	mp := make(map[rune][]Antenna)
	for row, line := range lines {
		rs := bytes.Runes(line)
		for col, r := range rs {
			if r != '.' {
				if _, ok := mp[r]; ok {
					mp[r] = append(mp[r], Antenna{Position{row, col}, r})
				} else {
					as := []Antenna{
						{Position{row, col}, r},
					}
					mp[r] = as
				}
			}
		}
	}
	return mp
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
