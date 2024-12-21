package main

import (
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"strings"
    "math"
)

func main() {
	input := string(readInput("input.txt"))
	lines := strings.Split(input, "\n")

	lq := make(PriorityQueue, len(lines) - 1)
	rq := make(PriorityQueue, len(lines) - 1)

	for i, line := range lines {
		words := strings.Fields(line)
		if len(words) == 2 {
			a, _ := strconv.Atoi(words[0])
			b, _ := strconv.Atoi(words[1])
            lq[i] = &Item{
                priority: a,
                index:    i,
            }
            rq[i]  = &Item{
                priority: b,
                index:    i,
            }
		}
	}
		
    heap.Init(&lq)
    heap.Init(&rq)

    i := 0
    total := 0
    for lq.Len() > 0 && rq.Len() > 0{
		left := heap.Pop(&lq).(*Item)
		right := heap.Pop(&rq).(*Item)
        diff := math.Abs(float64(left.priority - right.priority))
        total += int(diff)
        i++
	}
    fmt.Printf("TOTAL: %d\n", total)
}

func readInput(name string) []byte {
	file, err := os.ReadFile(name)
	if err != nil {
		err = fmt.Errorf("open file %s: %s\n", name, err)
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
	return file
}



// adapted from https://pkg.go.dev/container/heap#example-package-PriorityQueue
type Item struct {
	priority int
	index    int
}

var a heap.Interface

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
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
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) update(item *Item, priority int) {
	item.priority = priority
	heap.Fix(pq, item.index)
}
