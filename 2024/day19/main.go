package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type Piece []byte
type Chopped struct {
	head Piece
	tail Piece
}


var TOWS []Piece
var DESGS []Piece

func main() {
	TOWS, DESGS = towels()

    m := make(map[string]int)
    for _, des := range DESGS {
        m[string(des)] = count(des, 0)
    }

    count_any := 0
    count_total := 0
    for k, v := range m {
        fmt.Printf("%s : %d\n", k, v)
        if v > 0 {
            count_any++
        }
        count_total += v
    }
    fmt.Printf("Any: %d\n", count_any)
    fmt.Printf("Total: %d\n", count_total)
}


var COUNT_CACHE = map[string]int{}
func count(d Piece, c int) int {
	if cp, ok := COUNT_CACHE[string(d)]; ok {
		return cp
	}

	m, ok := chop(d)
    if !ok {
        COUNT_CACHE[string(d)] = 0
        return 0
    }

	for _, mt := range m {
		if len(mt.tail) == 0 {
			c += 1
		}
		c += count(mt.tail, 0)
	}

    COUNT_CACHE[string(d)] = c
	return c
}

var CHOP_CACHE = map[string][]Chopped{}
func chop(d Piece) ([]Chopped, bool) {
	if cp, ok := CHOP_CACHE[string(d)]; ok {
		return cp, len(cp) > 0
	}
	var cuts []Chopped
	for _, tow := range TOWS {
		if len(tow) > len(d) {
			continue
		}
		head := string(d[:len(tow)])
		tail := d[len(tow):]
		if head == string(tow) {
			cuts = append(cuts, Chopped{tow, tail})
		}
	}
	CHOP_CACHE[string(d)] = cuts
	return cuts, len(cuts) > 0
}

func towels() ([]Piece, []Piece) {
	var towels [][]byte
	var designs [][]byte
	file := readFile("input.txt")
	top, bot, found := bytes.Cut(file, []byte{'\n', '\n'})
	if !found {
		panic("read towels: could not split input file")
	}

	towels = bytes.Split(top, []byte{','})
	tows := make([]Piece, len(towels))
	for i := range towels {
		tows[i] = Piece(bytes.TrimSpace(towels[i]))
	}

	designs = bytes.Split(bot, []byte{'\n'})
	designs = designs[:len(designs)-1]
	desgs := make([]Piece, len(designs))
	for i := range designs {
		desgs[i] = Piece(designs[i])
	}

	return tows, desgs
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
