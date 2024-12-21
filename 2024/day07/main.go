package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
)

const RESULT = 145397611075341


type Equation struct {
	Res  int
	Nums []int
}

var OPS = []string{"+", "*", "||"}

func main() {
	eqs := readEquations()
    
    total := 0
    sum := make(chan int, len(eqs))
    var wg sync.WaitGroup
	for _, eq := range eqs {
        wg.Add(1)
        go handleEquation(sum, &wg, eq)
	}

    wg.Wait()
    close(sum)

    for n := range sum {
        total += n    
    }

    assert(total == RESULT, "different than the expected :(\n")
    fmt.Printf("Part Two: %d\n", total)
}

func handleEquation(ch chan<- int, wg *sync.WaitGroup, eq Equation) {
    exprs := parse(eq.Nums)
    for _, expr := range exprs {
        res := eval(expr)
        if res == eq.Res {
            ch <- res
            break
        }
    }
    wg.Done()
}

func parse(ns []int) []string {
    var exprs []string

    N := float64(len(ns) - 1)
    nops := int(math.Pow(3, N))
    
    for ops := 0; ops < nops; ops++ {
        var s string
        for i := 0; i < len(ns)-1; i++ {
            div := int(math.Pow(3, float64(i)))
            trunc := ops / div
            op := trunc % 3
            s += fmt.Sprintf("%d ", ns[i])
            s += fmt.Sprintf("%s ", OPS[op])
        }
        s += fmt.Sprintf("%d", ns[len(ns)-1])
        exprs = append(exprs, s)
    }
    return exprs
}


func parseBinary(ns []int) []string {
    var exprs []string
    var N uint16 = 2 << (len(ns) - 1)
    var ops uint16
    for ; ops < N; ops++ {
        var s string
        var i uint16
        for ; i < uint16(len(ns)-1); i++ {
            op := (ops >> i) & 1
            s += fmt.Sprintf("%d ", ns[i])
            s += fmt.Sprintf("%s ", OPS[op])
        }
        s += fmt.Sprintf("%d", ns[len(ns)-1])
        exprs = append(exprs, s)
    }
    return exprs
}

func eval(expr string) int {
    fields := strings.Fields(expr) 
    for len(fields) >= 3 {
        var cur string
        sl := fields[0]
        sr := fields[2]

        op := fields[1]

        if op == "||" {
            cur = sl + sr
        } else {
            left, _ := strconv.Atoi(sl)
            right, _ := strconv.Atoi(sr)
            if op == "+" {
                cur = strconv.Itoa(left + right)
            } else if op == "*" {
                cur = strconv.Itoa(left * right)
            }
        }
        fields = fields[2:]
        if len(fields) > 0 {
            fields[0] = cur
        }
    }
    n, _ := strconv.Atoi(fields[0])
    return n
}

func readEquations() []Equation {
	file := readFile("input.txt")
	lines := bytes.Split(file, []byte{'\n'})
	lines = lines[:len(lines)-1]
	eqs := make([]Equation, len(lines))
	for i, line := range lines {
		parts := bytes.Split(line, []byte{':'})
		b_res, b_expr := parts[0], parts[1]
		res, _ := strconv.Atoi(string(b_res))
		f := bytes.Fields(b_expr)
		expr := make([]int, len(f))
		for j, b := range f {
			n, _ := strconv.Atoi(string(b))
			expr[j] = n
		}
		eqs[i] =  Equation{
			Res:  res,
			Nums: expr,
		}
	}
	return eqs
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
