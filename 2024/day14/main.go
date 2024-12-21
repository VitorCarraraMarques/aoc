package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

const PART = 2

const ROWS = 103
const COLS = 101

// const ROWS = 7
// const COLS = 11

// ansi escape code for reseting print position in terminal
var RESET = fmt.Sprintf("\033[%dA\033[%dD", ROWS+3, COLS)

type Vec2 struct {
	x int
	y int
}

type Robot struct {
	pos Vec2
	vel Vec2
}

func (r *Robot) Update(dt int) {
	//(a % b) + b) % b => mod even negative numbers
	r.pos.x = ((r.pos.x+r.vel.x*dt)%COLS + COLS) % COLS
	r.pos.y = ((r.pos.y+r.vel.y*dt)%ROWS + ROWS) % ROWS
}

type Grid struct {
	Robots map[Vec2]int
}

func (g *Grid) Step(rs []*Robot, dt int) {
	for _, r := range rs {
		g.Robots[r.pos]--
		r.Update(dt)
		g.Robots[r.pos]++
	}
}

func (g *Grid) HTML() string {
	var s string
	s += "<div>"
	for j := range ROWS {
		s += "<div style=\"display:flex;flex-direction:flex-row;gap:0;\">"
		for i := range COLS {
			s += "<div style=\"height:10px;width:10px;border:1px solid black;"
			if g.Robots[Vec2{x: i, y: j}] > 0 {
				s += "background-color: black;\""
			} else {
				s += "background-color: white;\""
			}
			s += "></div>"
		}
		s += "</div>"
	}
	s += "</div>"
	return s
}

func (g *Grid) Print() {
	for j := range ROWS {
		for i := range COLS {
			if g.Robots[Vec2{x: i, y: j}] > 0 {
				fmt.Printf("##")
			} else {
				fmt.Printf("  ")
			}
		}
		fmt.Printf("\n")
	}
}

func NewGrid(rs []*Robot) Grid {
	g := Grid{
		Robots: make(map[Vec2]int),
	}
	for _, r := range rs {
		g.Robots[r.pos]++
	}
	return g
}

func main() {
    // see solution on terminal when t = 8149
	solve()

    // see solution on browser on localhost:8000/gen/8149
    //serve()
}

func serve() {
    http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {w.WriteHeader(404)})
    http.HandleFunc("/gen/{t}", handle)
    http.ListenAndServe("localhost:8000", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	var rs = readRobots()
	var g = NewGrid(rs)
	t, _ := strconv.Atoi(r.PathValue("t"))
	g.Step(rs, t)
	fmt.Printf("t = %d\n", t)
	fmt.Fprintf(w, g.HTML())
}

func solve() {
	rbs := readRobots()
	if PART == 1 {
		part1(rbs)
	} else if PART == 2 {
		g := NewGrid(rbs)
		g.Print()
		for i := 0; ; {
			var dt int
			_, err := fmt.Scanf("%d", &dt)
			if err != nil {
				fmt.Printf("\n")
				dt = 1
			}
			i += dt
			fmt.Printf("t = %d\n", i)
			fmt.Printf("%s", RESET)
			g.Step(rbs, dt)
			g.Print()
		}
	}
}

func part1(rbs []*Robot) {
	row_half := (ROWS / 2)
	col_half := (COLS / 2)
	q1, q2, q3, q4 := 0, 0, 0, 0
	for _, r := range rbs {
		r.Update(100)

		if r.pos.x < col_half && r.pos.y < row_half {
			q1++
		} else if r.pos.x > col_half && r.pos.y < row_half {
			q2++
		} else if r.pos.x < col_half && r.pos.y > row_half {
			q3++
		} else if r.pos.x > col_half && r.pos.y > row_half {
			q4++
		}
	}
	fmt.Printf("Final Safety Factor: %d\n", q1*q2*q3*q4)
}

func readRobots() []*Robot {
	file := readFile("input.txt")
	lines := bytes.Split(file, []byte{'\n'})
	lines = lines[:len(lines)-1]
	rbs := make([]*Robot, len(lines))
	for i, ln := range lines {
		f := bytes.Fields(ln)
		rb := Robot{}

		p := bytes.Split(f[0][2:], []byte{','})
		rb.pos.x, _ = strconv.Atoi(string(p[0]))
		rb.pos.y, _ = strconv.Atoi(string(p[1]))

		v := bytes.Split(f[1][2:], []byte{','})
		rb.vel.x, _ = strconv.Atoi(string(v[0]))
		rb.vel.y, _ = strconv.Atoi(string(v[1]))

		rbs[i] = &rb
	}
	return rbs

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
