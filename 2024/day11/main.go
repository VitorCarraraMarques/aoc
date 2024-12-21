package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
)

const N_BLINKS = 75

func main() {
	stones := readStones()
	fmt.Printf("Initial Arrangement: \n%v\n", stones)
	for i := 1; i <= N_BLINKS; i++ {
	    fmt.Printf("   Blinking #%d\n", i)
        stones = blink(stones)
    }
    total := 0
    for _, n := range stones {
        total += n
    }

    fmt.Printf("the total number of Stones: %d\n", total)

}

func blink(stones map[int]int) map[int]int {
    next := make(map[int]int)
    for stone, count := range stones {
		if stone == 0 {
            next[1] += count
		} else if l := digits(stone); l%2 == 0 {
			half := l / 2
			left := stone / pow(10, half)
			right := stone % pow(10, half)
            next[left] += count
            next[right] += count
		} else {
            m := stone * 2024
			next[m] += count
		}
	}
	return next
}

func readStones() map[int]int {
	file := readFile("input.txt")
	fields := bytes.Fields(file)
	stones := make(map[int]int, len(fields))
	for _, b := range fields {
        n, _ :=  strconv.Atoi(string(b))
		stones[n]++
	}
	return stones
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

func digits(a int) int {
	return int(math.Floor(math.Log10(float64(a))) + 1)
}

func pow(base, exp int) int {
    result := 1
    for {
        if exp & 1 == 1 {
            result *= base
        }
        exp >>= 1
        if exp == 0 {
            break
        }
        base *= base
    }
    return result
}
