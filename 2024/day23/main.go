package main

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
)


func main() {
    part2()
}

func part2() {
    adj := readAdj()
    sets := make(map[string]bool)
    for k := range adj {
        check(adj, k, k, sets)
    }
    maxlen := 0
    var big string
    for k := range sets {
        if len(k) > maxlen {
            maxlen = len(k)
            big = k
        }
    }
    out := ""
    for i := 0; i < len(big); i += 2 {
        out += big[i:i+2]
        out += ","
    }
    fmt.Printf("%s\n", out)
}



func check(adj map[string][]string, src, set string, sets map[string]bool){
    nbors := adj[src]
neighbors:
    for _, nbor := range nbors {
        if strings.Contains(set, nbor) {
            continue
        }
        if !slices.Contains(adj[nbor], src) {
            continue
        }
        for i := 0; i < len(set); i += 2 {
            item := set[i:i+2]
            if !slices.Contains(adj[nbor], item) {
                continue neighbors
            }
        }
        set += nbor
        set = strsort(set, 2)
        sets[set] = true
        check(adj, nbor, set, sets)
    }
}

func part1() {
    adj := readAdj()
    sets := make(map[string]bool)
    for src, nbors := range adj {
        if len(nbors) >= 2 {
            for i := 0; i < len(nbors); i++ {
                for j := i+1; j < len(nbors); j++ {
                    a := nbors[i]
                    b := nbors[j]
                    if a_nbors, ok := adj[a] ;ok && slices.Contains(a_nbors, b) {
                        set := src+ a + b
                        set = strsort(set, 2)
                        sets[set] = true
                    }
                }
            }
        }
    }

    count := 0
counter:
    for set := range sets {
        for i := 0; i < len(set); i += 2 {
            if set[i] == 't' {
                count++
                continue counter
            }
        }
    }
    fmt.Printf("Count: %d\n", count)
}

type SortString []string
func (a SortString ) Len() int           { return len(a) }
func (a SortString ) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortString ) Less(i, j int) bool { return a[i] < a[j] }
func strsort(s string, size int) string {
    chop := []string{}
    for i := 0; i < len(s); i += size {
        chop = append(chop, s[i:i+size])
    }
    sorted := SortString(chop)
    sort.Sort(&sorted)
    out := ""
    for _, l := range sorted {
        out += l 
    } 
    return out
}

func readAdj() map[string][]string {
    file := readFile("input.txt")
    lines := bytes.Split(file, []byte{'\n'})
    lines = lines[:len(lines)-1]
    adj := make(map[string][]string)
    for _, line := range lines {
        sp := bytes.Split(line, []byte{'-'})    
        left := string(sp[0])
        right := string(sp[1])
        adj[left] = append(adj[left], right)
        adj[right] = append(adj[right], left)
    }
    return adj
}


func readFile(name string) []byte{
    file, err := os.ReadFile(name)
    if err != nil {
        panic(fmt.Sprintf("readFile %s: %v", name, err))
    }
    return file
}
