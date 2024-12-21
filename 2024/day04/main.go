package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if os.Args[1] == "1" {
		one()
	} else {
		two()
	}
}

func one() {
	DIRS := [][]int{
		{-1, -1},
		{0, -1},
		{1, -1},
		{-1, 0},
		{1, 0},
		{-1, 1},
		{0, 1},
		{1, 1},
	}
	search := readSearch()
	sum := 0
	for i := 0; i < len(*search); i++ {
		line := (*search)[i]
		for j := 0; j < len(line); j++ {
			let := line[j]
			if let == 'X' {
				idx := 0
				for _, dir := range DIRS {
					fmt.Printf("\n%d,%d: %s\n", i, j, string(let))
					if checkNeig(search, i, j, dir, idx) {
						sum += 1
					}
				}
			}
		}
	}
	fmt.Printf("FINAL ANSWER: %d\n", sum)
}
func checkNeig(search *[][]byte, i, j int, dir []int, idx int) bool {
	var NEIG = []byte{'M', 'A', 'S'}
	if idx > 2 {
		return false
	}

	exp := NEIG[idx]
	ni := i + dir[0]
	if ni < 0 || ni >= len((*search)[0]) {
		return false
	}
	nj := j + dir[1]
	if nj < 0 || nj >= len((*search)) {
		return false
	}
	neig := (*search)[ni][nj]
	if exp == neig {
		fmt.Printf("%d,%d: %s\n", ni, nj, string(neig))
		if neig == 'S' {
			return true
		}
		return checkNeig(search, ni, nj, dir, idx+1)
	}
	return false
}

func two() {
	search := readSearch()
	sum := 0
	for i := 0; i < len(*search); i++ {
		line := (*search)[i]
		for j := 0; j < len(line); j++ {
			let := line[j]
			if let == 'M' {
				sum += checkBlocks(search, i, j)
			}
		}
	}
	fmt.Printf("FINAL ANSWER: %d\n", sum)
}

func checkBlocks(s *[][]byte, i, j int) int {
	dirs := [][]int{
		{-1, -1},
		{1, -1},
		{-1, 1},
		{1, 1},
	}
	n := 0
	for _, dir := range dirs {
		ci, cj := i+dir[0], j+dir[1]
		if ci < 1 || ci >= len((*s)[0])-1 {
			continue
		}
		if cj < 1 || cj >= len((*s))-1 {
			continue
		}
		center := (*s)[ci][cj]
		if center == '!' {
			continue
		}
		left_diag := string((*s)[ci-1][cj-1]) + string(center) + string((*s)[ci+1][cj+1])
		right_diag := string((*s)[ci+1][cj-1]) + string(center) + string((*s)[ci-1][cj+1])
		if checkDiags(left_diag, right_diag) {
			(*s)[ci][cj] = '!'
			n += 1
		}
	}
	return n
}

func checkDiags(left_diag, right_diag string) bool {
	return (left_diag == "MAS" || left_diag == "SAM") && (right_diag == "MAS" || right_diag == "SAM")
}

func readSearch() *[][]byte {
	file := readFile("input.txt")
	lines := strings.Split(string(file), "\n")
	var search [][]byte
	for _, line := range lines {
		n := len(line)
		letters := make([]byte, n)
		for i := 0; i < n; i++ {
			letters[i] = line[i]
		}
		if len(letters) > 0 {
			search = append(search, letters)
		}
	}
	return &search
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
