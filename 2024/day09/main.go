package main

import (
	"fmt"
	"os"
	"strconv"
)

const PART = 2

func main() {
    d := readDiskMap()
    b := blocks(d)
    if PART == 1 {
        b = fragment(b)
    } else if PART == 2 {
        b = move(b)
    }
    s := checksum(b)
    fmt.Printf("Part %d: Checksum = %d\n", PART, s)
}

func move(b []int) []int {
    rstart := len(b) - 1
    for rstart > 1 {
        for b[rstart] == -1 {
            rstart--
        }
        size := fileSize(b, rstart)

        for lstart := 0; lstart < rstart; {
            for b[lstart] != -1 {
                lstart++
            }
            if lstart >= rstart {
                break
            }
            empty := emptySize(b, lstart, rstart)
            if empty >= size {
                for i := 0; i < size; i++ {
                    b[lstart+i], b[rstart-i] = b[rstart-i], b[lstart+i]
                }        
                break
            } 
            lstart += empty
        }
        rstart -= size
    }
    return b
}
func emptySize(b []int, left int, stop int) int {
    l := left+1
    empty := 1
    for l < stop && b[l] == -1 {
        empty++
        l++
    }
    return empty
}

func fileSize(b []int, right int) int {
    r := right - 1
    curf := b[right]
    size := 1
    for r > 0 && b[r] == curf {
        size++
        r--
    }
    return size
}

func fragment(b []int) []int {
    left := 0
    right := len(b) - 1
    for left <= right {
        for b[left] == -1 {
            for b[right] == -1 {
                right--
            }
            b[left] = b[right]
            b[right] = -1
            left++
            right--
        }
        left++
    }
    return b
}

func blocks(d []int) []int{
    var blocks []int
    file := 0
    for i, n := range d {
        if i % 2 == 0 {
            for j := 0; j < n; j++ {
                blocks = append(blocks, file)
            }
            file++
        } else {
            for j := 0; j < n; j++ {
                blocks = append(blocks, -1)
            } 
        }
    }
    return blocks
}

func checksum(b []int) int {
    sum := 0
    for i, n := range b {
        if n == -1 {
            continue
        }
        sum += i*n
    }
    return sum

}

func readDiskMap() []int {
    rs := readFile("input.txt")
    d := make([]int, len(rs))
    for i, r := range rs {
        n, _ := strconv.Atoi(string(r))
        d[i] = n
    }
    return d
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
