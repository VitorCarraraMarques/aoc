package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
    day := os.Args[1]
    if day == "1" {
        one()
    } else {
        two()
    }
}

func one() {
    reports := readReports()
    safe_count := 0
    for _, levels := range reports {
        if isSafe(levels) {
            safe_count++
            fmt.Printf("Const Change: %v\n", levels)
        }
    }
    fmt.Printf("FINAL ANSWER: %d\n", safe_count)
}

func two(){
    reports := readReports()
    safe_count := 0
    var unsafe [][]int
    for _, levels := range reports {
        if isSafe(levels) {
            safe_count++
            fmt.Printf("SAFE: %v\n", levels)
        } else {
            unsafe = append(unsafe, levels)
        }
    }
    fmt.Printf("\n---------------------------------\n")
    for _, levels := range unsafe {
        fmt.Printf("\nLEVELS BEFORE: %v\n", levels)
        for i := 0; i < len(levels); i++ {
            sub := make([]int, len(levels))
            copy(sub, levels)
            sub = append(sub[:i], sub[i+1:]...)
            fmt.Printf("\tSUB: %v\n", sub)
            if isSafe(sub){
                safe_count++
                break
            }
        } 
        fmt.Printf("LEVELS AFTER: %v\n\n", levels)
    }
    fmt.Printf("FINAL ANSWER: %d\n", safe_count)
}

func isSafe(lvls []int) bool {
    return checkChange(lvls) && checkDiff(lvls)
}

func checkChange(lvls []int) bool {
    change := "default"
    for i := 0; i < len(lvls) - 1; i++ {
        cur := lvls[i+1]
        prev := lvls[i]
        var cur_change string
        if cur > prev {
            cur_change = "inc"
        } else if cur < prev {
            cur_change = "dec"
        } else {
            cur_change = "eq"
        }
        if change != "default" && change != cur_change {
            return false    
        }
        change = cur_change
    }
    return true
}

func checkDiff(lvls []int) bool {
    for i := 0; i < len(lvls) - 1; i++ {
        cur := lvls[i+1]
        prev := lvls[i]
        diff := math.Abs(float64(cur - prev)) 
        if diff < 1 || diff > 3 {
            return false
        }
    }
    return true
}

func readReports() [][]int {
	input := string(readFile("input.txt"))
	lines := strings.Split(input, "\n")
    reports := make([][]int, len(lines) - 1)

    for n, line := range lines {
        if len(line) < 2 {
            continue 
        }
        words := strings.Fields(line)
        levels := make([]int, len(words))
        for i, w := range words {
            n, _ := strconv.Atoi(w)
            levels[i] = n
        }
        reports[n] = levels
    }
    return reports
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



