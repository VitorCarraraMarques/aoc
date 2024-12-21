package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
    "sort"
)

func main() {
    cor, incor := readInput()
    fmt.Printf("sum middle correct order: %d\n", cor)
    fmt.Printf("sum middle incorrect order: %d\n", incor)
}



func readInput() (int, int) {
	file := string(readFile("input.txt"))
	parts := strings.Split(file, "\n\n")
	rules := readRules(parts[0])
	upt := readUpdates(parts[1])
    cor, incor := 0, 0
    for _, u := range upt {
        temp := make([]int, len(u))
        copy(temp, u)
        update := Update{
            items: temp,
            rules: rules,
        }
        sort.Sort(update)
        if isEqual(update.items, u) {
            m := middle(u)
            cor += m
        } else {
            m := middle(update.items)
            incor += m 
        }

    }
    return cor, incor
}

func isEqual(a, b []int) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        n1 := a[i]
        n2 := b[i]
        if n1 != n2 {
            return false
        }
    }
    return true
}

func middle(arr []int) int {
	idx := len(arr) / 2
	return arr[idx]
}

func readRules(s string) map[int][]int {
	lines := strings.Split(s, "\n")
	rules := make(map[int][]int)
	for _, line := range lines {
		parts := strings.Split(line, "|")
		left, _ := strconv.Atoi(parts[0])
		right, _ := strconv.Atoi(parts[1])
		rules[right] = append(rules[right], left)
	}
	return rules
}

func readUpdates(s string) [][]int {
	lines := strings.Split(s, "\n")
    lines = lines[:len(lines) - 1]
	updates := make([][]int, len(lines))
	for i, line := range lines {
		w := strings.Split(line, ",")
		upt := make([]int, len(w))
		for j, letter := range w {
			n, _ := strconv.Atoi(letter)
			upt[j] = n
		}
		updates[i] = upt
	}
	return updates
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

type Update struct {
    items []int
    rules map[int][]int
}

func (upt Update) Len() int {
    return len(upt.items)
}

func (upt Update) Swap(i, j int) {
    upt.items[i], upt.items[j] = upt.items[j], upt.items[i]
}

func (upt Update) Less(i, j int) bool {
    //the ith item (=a) is less than the jth item (=b), if items[i] is present in rules[b]
    a := upt.items[i]
    b :=  upt.items[j]
    less_than_b := upt.rules[b] 
    for _, n := range less_than_b {
        if n == a {
            return true
        }
    }
    return false
}


