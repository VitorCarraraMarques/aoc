package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)


func main() {
	input := string(readInput("input.txt"))
	lines := strings.Split(input, "\n")

    left := make([]int, len(lines))
    count := make(map[int]int, len(lines))
	for i, line := range lines {
		words := strings.Fields(line)
		if len(words) == 2 {
			a, _ := strconv.Atoi(words[0])
            left[i] = a

			b, _ := strconv.Atoi(words[1])
            count[b]++
		}
	}
    fmt.Printf("LEFT LIST: %v\n", left)
    fmt.Printf("COUNT MAP: %v\n", count)

    total := 0 
    for _, val := range left {
        scr := count[val] * val
        fmt.Printf("    Val: %d\n", val)
        fmt.Printf("  Count: %d\n", count[val])
        fmt.Printf("    Scr: %d\n\n", scr)
        total += scr
    }

    fmt.Printf("TOTAL: %v\n", total)

		
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

