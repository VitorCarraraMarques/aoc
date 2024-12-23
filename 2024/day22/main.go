package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func main() {
    n_iter := 2000
    secs := readInitial()
    buyers := make(map[int]map[string]int)
    for _, sec := range secs {
        n := sec 
        m := iter(n, n_iter)
        buyers[sec] = m
    }

    var seqs []string
    for a := -9; a <= 9; a++ {
        for b := -9; b <= 9; b++ {
            for c := -9; c <= 9; c++ {
                for d := -9; d <= 9; d++ {
                    seqs = append(seqs, 
                        strconv.Itoa(a) +
                        strconv.Itoa(b) +
                        strconv.Itoa(c) +
                        strconv.Itoa(d),
                    )
                }
            }
        }
    }

    bananas := 0
    for _, seq := range seqs {
        var curban int
        for cur_buyer := range buyers {
            curban += buyers[cur_buyer][seq]
        }
        bananas = max(bananas, curban)
    }
    fmt.Printf("Bananas: %d\n", bananas)
}

type Line struct {
    price int
    diff int
}

func iter(f, n_iter int) map[string]int {
    var ls []Line
    var prev int
    n := f
    for i := range n_iter {
        dig := n % 10
        diff := 69
        if i != 0 {
            diff = dig - prev
        }
        ls = append(ls, Line{dig, diff})

        prev = dig
        n = nxtsec(n)
    }

    prices := ls
    m := make(map[string]int)
    for i := 0; i < len(prices) - 4; i++ {
        k := ""
        k += strconv.Itoa(prices[i+1].diff)
        k += strconv.Itoa(prices[i+2].diff)
        k += strconv.Itoa(prices[i+3].diff)
        k += strconv.Itoa(prices[i+4].diff)
        if _, ok := m[k]; ok {
            continue
        }
        m[k] = prices[i+4].price
    }
    return m
}

func nxtsec(sec int) int {
    val := sec * 64
    sec = mix(sec, val)
    sec = prune(sec)

    val = sec / 32
    sec = mix(sec, val)
    sec = prune(sec)

    val = sec * 2048
    sec = mix(sec, val)
    sec = prune(sec)
    
    return sec
}

func mix(sec, val int) int {
    return val ^ sec
}

func prune(sec int) int {
    return sec % 16777216
}

func readInitial() []int {
    file := readFile("input.txt")
    lines := bytes.Split(file, []byte{'\n'})
    lines = lines[:len(lines)-1]
    nums := make([]int, len(lines))
    for i, line := range lines {
        n, _ := strconv.Atoi(string(line))
        nums[i] = n
    }
    return nums
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

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
