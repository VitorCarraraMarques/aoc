package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var ALLOWED = map[string]bool{
	",": true,
	"(": true,
	")": true,
    "0": true,
	"1": true,
	"2": true,
	"3": true,
	"4": true,
	"5": true,
	"6": true,
	"7": true,
	"8": true,
	"9": true,
}

func main() {
	p := os.Args[1]
	if p == "1" {
		one()
	} else {
		two()
	}
}

func one() {
	input := readFile("input.txt")
	idx := 0
    sum := 0
	for idx < len(input) {
		w := input[idx]
		if w == 'm' {
			sub := string(input[idx:idx+4])
			if sub == "mul(" {
                var cur string
                idx += 4
                i := 0
                for i <= 7 {
                    a := string(input[idx+i])
                    if !ALLOWED[a] {
                        fmt.Printf("   Illegal Char: %s\n", a)
                        break
                    }
                    if a == ")" {
                        fmt.Printf("   Parsed Expr : %s\n", cur)
                        sum += eval(cur)
                        break
                    }
                    cur += a
                    i++
                }
                idx += i
			}
        }
        idx++
	}
    fmt.Printf("FINAL ANSWER: %d\n", sum)
}

func two() {
	input := readFile("input.txt")
	idx := 0
    sum := 0
    ignore := false
	for idx < len(input) {
		w := input[idx]
		if w == 'm' {
            if ignore {
                idx++
                continue
            }
			sub := string(input[idx:idx+4])
			if sub == "mul(" {
                var cur string
                idx += 4
                i := 0
                for i <= 7 {
                    a := string(input[idx+i])
                    if !ALLOWED[a] {
                        fmt.Printf("   Illegal Char: %s\n", a)
                        break
                    }
                    if a == ")" {
                        fmt.Printf("   Parsed Expr : %s\n", cur)
                        sum += eval(cur)
                        break
                    }
                    cur += a
                    i++
                }
                idx += i
			}
        }
        if w == 'd' {
			do := string(input[idx:idx+4])
			if do == "do()" {
                ignore = false
                idx += 4
                continue
            }
            dont := string(input[idx:idx+7])
            if dont == "don't()" {
                ignore = true
                idx += 7
                continue
            }
        }
        idx++
	}
    fmt.Printf("FINAL ANSWER: %d\n", sum)
}

func eval(expr string) int {
    ns := strings.Split(expr, ",")
    a, _ := strconv.Atoi(ns[0])
    b, _ := strconv.Atoi(ns[1])
    return a*b

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
