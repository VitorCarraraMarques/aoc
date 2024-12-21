package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"maps"
	"os"
	"os/exec"
)

const VISUALIZATION = true

var display []byte
var reindeer = Reindeer{dir: 2}
var ROWS int
var COLS int

var DIRS = []Vec2{
	{y: 0, x: 1},
	{y: 1, x: 0},
	{y: 0, x: -1},
	{y: -1, x: 0},
}

type Reindeer struct {
	pos   Vec2
	goal  Vec2
	dir   int
	score int
}

func (r *Reindeer) TurnLeft() {
	r.dir = (((r.dir - 1) % 4) + 4) % 4
	r.score += 1000
}

func (r *Reindeer) TurnRight() {
	r.dir = (((r.dir + 1) % 4) + 4) % 4
	r.score += 1000
}

func (r *Reindeer) Move() {
	r.pos = r.pos.Add(DIRS[r.dir])
	r.score = r.score + 1
	if display[r.pos.Idx()] != '#' {
		display[r.pos.Idx()] = '@'
	}
}

func (r *Reindeer) Nbors() []Reindeer {
	var nbors []Reindeer

	front := *r
	front.Move()
	if display[front.pos.Idx()] != '#' {
		nbors = append(nbors, front)
	}

	back := *r
	back.TurnRight()
	back.TurnRight()
	back.Move()
	if display[back.pos.Idx()] != '#' {
		nbors = append(nbors, back)
	}

	left := *r
	left.TurnLeft()
	left.Move()
	if display[left.pos.Idx()] != '#' {
		nbors = append(nbors, left)
	}

	right := *r
	right.TurnRight()
	right.Move()
	if display[right.pos.Idx()] != '#' {
		nbors = append(nbors, right)
	}

	return nbors
}

type Path map[Vec2]bool

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

func (v Vec2) InBound() bool {
	return v.x >= 0 && v.x < COLS && v.y >= 0 && v.y < ROWS
}

type Item struct {
	value Reindeer
	cost  int
	index int
	path  Path
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].cost < pq[j].cost
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) update(item *Item, cost int) {
	item.cost = cost
	heap.Fix(pq, item.index)
}

func main() {
	readGrid()
	ends := search(reindeer.pos, reindeer.goal)
	clear()
	for i, item := range ends {
		fmt.Printf("--- Path #%d ----\n", i)
		fmt.Printf("Score: %d\n", item.value.score)
		printPath(item.path)
	}

	best_score := 130536
	for _, item := range ends {
		if item.value.score < best_score {
			best_score = item.value.score
		}
	}

	unique := make(Path)
	for _, item := range ends {
		if item.value.score == best_score {
			for p := range item.path {
				unique[p] = true
			}
		}
	}
	fmt.Printf("Number of Seats: %d\n", len(unique))
}

func printPath(p Path) {
	for j := range ROWS {
		for i := range COLS {
			pos := Vec2{j, i}
			if p[pos] {
				fmt.Printf("@")
			} else if display[pos.Idx()] == '#' {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func search(start Vec2, goal Vec2) []Item {
	best_score := 130536
	var ends []Item

	type P struct {
		V Vec2
		D int
	}
	costs := make(map[P]int)
	d := reindeer.dir
	costs[P{start, d}] = 0

	it := Item{
		value: reindeer,
		index: 0,
		cost:  costs[P{start, d}],
		path:  make(Path),
	}
	pq := PriorityQueue{&it}
	heap.Init(&pq)

	for i := 0; pq.Len() > 0; i++ {
		item := heap.Pop(&pq).(*Item)

		// clear()
		// fmt.Printf("#%d Iteration | %d Items Enqueued | Current Best Score: %d | Current Score: %d\n", i, pq.Len(), best_score, item.value.score)
		// if VISUALIZATION {
		// 	printGrid()
		// 	time.Sleep(100 * time.Microsecond)
		// }

		if item.value.score > best_score {
			continue
		}

		item.path[item.value.pos] = true
		if item.value.pos == goal {
			ends = append(ends, *item)
			best_score = item.value.score
		}

		nbors := item.value.Nbors()
		for _, n := range nbors {
			if item.path[n.pos] {
				continue
			}

			if c, ok := costs[P{n.pos, n.dir}]; ok {
				if n.score > c {
					continue
				}
			}
			costs[P{n.pos, n.dir}] = n.score

			if display[n.pos.Idx()] == '#' {
				continue
			}

			p_nbor := make(Path)
			maps.Copy(p_nbor, item.path)
			nbor := Item{
				value: n,
				cost:  n.score,
				path:  p_nbor,
			}
			heap.Push(&pq, &nbor)
		}
	}
	return ends
}

func readGrid() {
	file := readFile("input.txt")
	rows := bytes.Split(file, []byte{'\n'})
	ROWS = len(rows)
	COLS = len(rows[0])
	display = make([]byte, ROWS*COLS)

	for j, row := range rows {
		for i, cell := range row {
			display[j*COLS+i] = cell
			if cell == 'S' {
				reindeer.pos.y = j
				reindeer.pos.x = i
			}
			if cell == 'E' {
				reindeer.goal.y = j
				reindeer.goal.x = i
			}
		}
	}
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
