package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
    keys, locks := readKeyLock()
    fit := 0
    for _, lock := range locks {
        for _, key := range keys {
            if testKey(key, lock) {
                fit++
            }
        }
    }
    fmt.Printf("There is %d keys/locks pairs that fit together\n", fit)
}

func testKey(key []int, lock []int) bool {
    if len(lock) != len(key) {
        return false
    }
    for i := range len(lock) {
        if lock[i] + key[i] > 5 {
            return false
        }
    }
    return true
}

func readKeyLock() (keys [][]int, locks [][]int) {
    file := readFile("input.txt")
    items := bytes.Split(file, []byte{'\n', '\n'})
    fmt.Printf("items: \n%v\n", items)
    full := "#####"
    empty := "....."
    for n, item := range items {
        lines := bytes.Split(item, []byte{'\n'})
        if n == len(items)-1{
            lines = lines[:len(lines)-1]
        }
        fmt.Printf("lines: \n%v\n", lines)
        top := string(lines[0])
        bot := string(lines[len(lines)-1])
        lines = lines[1:len(lines)-1]

        item := []int{}
        for i := range len(lines[0]) {
            height := 0
            for j := range len(lines) {
                if lines[j][i] == '#' {
                    height++
                }            
            }
            item = append(item, height)
        }

        if top == full && bot == empty {
            locks = append(locks, item)
        }
        if bot == full && top == empty {
            keys = append(keys, item)
        }
    }
    return
}

func readFile(name string) []byte {
    file, err := os.ReadFile(name)
    if err != nil {
        panic(fmt.Sprintf("readFile %s: %v", name, err))
    }
    return file
}

