package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Instr struct {
	lef string
	rig string
	op  string
	out string
}

type Output struct {
	name string
	val  int
}

type SortOutput []Output
func (a SortOutput) Len() int           { return len(a) }
func (a SortOutput) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortOutput) Less(i, j int) bool { return a[i].name > a[j].name }

func main() {
    fmt.Printf("x + y + c = S, c\n")
    x, y, c := 1, 0, 0
    s, co := fulladder(x,y,c)
    fmt.Printf("%d + %d + %d = %d, %d\n", x, y, c, s, co)

    x, y, c = 0, 1, 0
    s, co = fulladder(x,y,c)
    fmt.Printf("%d + %d + %d = %d, %d\n", x, y, c, s, co)

    x, y, c = 0, 0, 1
    s, co = fulladder(x,y,c)
    fmt.Printf("%d + %d + %d = %d, %d\n", x, y, c, s, co)

    x, y, c = 1, 1, 0
    s, co = fulladder(x,y,c)
    fmt.Printf("%d + %d + %d = %d, %d\n", x, y, c, s, co)

    x, y, c = 1, 0, 1
    s, co = fulladder(x,y,c)
    fmt.Printf("%d + %d + %d = %d, %d\n", x, y, c, s, co)

    x, y, c = 0, 1, 1
    s, co = fulladder(x,y,c)
    fmt.Printf("%d + %d + %d = %d, %d\n", x, y, c, s, co)

    x, y, c = 1, 1, 1
    s, co = fulladder(x,y,c)
    fmt.Printf("%d + %d + %d = %d, %d\n", x, y, c, s, co)

    return
	wires, instrs := readInput()
    setWires(wires, 'y', "1")

    for x := range 44 {
        xs := "1"
        for range x {
            xs += "0"
        }
        setWires(wires, 'x', xs)
        setWires(wires, 'y', xs)

        x_input := wireToBinary(wires, 'x')
        y_input := wireToBinary(wires, 'y')
        fmt.Printf("----\n")
        fmt.Printf("InputX: Binary = %s | Decimal = %d\n", x_input, toDecimal(x_input))
        fmt.Printf("InputY: Binary = %s | Decimal = %d\n", y_input, toDecimal(y_input))
        out := eval(wires, instrs)
        fmt.Printf("Output: Binary = %s | Decimal = %d | Expected = %d\n", out, toDecimal(out), toDecimal(x_input) + toDecimal(y_input))
    }
}


func fulladder(x, y, cin int) (s, cout int) {
    s1 := x ^ y
    c1 := x & y
    c2 := s1 & cin
    s = s1 ^ cin
    cout = c1 | c2
    return
}


func eval(wires map[string]int, instrs []Instr) string {
	done := false
	for !done {
		done = iterinstr(wires, instrs)
	}
    sout := wireToBinary(wires, 'z')
	return sout
}

func iterinstr(wires map[string]int, instrs []Instr) bool {
	done := true
	for _, instr := range instrs {
		l, ok1 := wires[instr.lef]
		r, ok2 := wires[instr.rig]
		if !ok1 || !ok2 {
			done = false
			continue
		}
		o := instr.out
		switch instr.op {
		case "AND":
			res := l & r
			wires[o] = res
		case "OR":
			res := l | r
			wires[o] = res
		case "XOR":
			res := l ^ r
			wires[o] = res
		}
	}
	return done
}

func vizPrintInstrs(instrs []Instr) {
    fmt.Printf("digraph {\n")
    for _, instr := range instrs {
        sop := fmt.Sprintf("%s_%s_%s", instr.lef, instr.op, instr.rig)
        fmt.Printf(" %s -> %s\n", instr.lef, sop)
        fmt.Printf(" %s -> %s\n", instr.rig, sop)
        fmt.Printf(" %s -> %s", sop, instr.out)
        if instr.out[0] == 'z' {
            fmt.Printf(" [color=red]")
        } 
        fmt.Printf("\n")
    }
    fmt.Printf("}\n")

}

func setWires(wires map[string]int, letter byte, val string) {
    count := countWires(wires, letter)
    n := len(val)
    if count < n {
        panic(fmt.Sprintf("setWires: %d exceeds the maximum amount of bits", count))
    }
    if count > n {
        for range count - n {
            val = "0" + val
        }
    }
    n = len(val) - 1
    for i := n; i >= 0; i-- {
        d := n - i
        k := fmt.Sprintf("%c%.2d", letter, d)
        num, _ := strconv.Atoi(string(val[i]))
        wires[k] = num
    }
}

func countWires(wires map[string]int, letter byte) int {
    count := 0
    for k := range wires {
        if k[0] == letter {
            count++
        }
    }
    return count
}

func wireToBinary(wires map[string]int, letter byte) string {
	var out []Output
	for name, val := range wires {
		if name[0] == letter {
			o := Output{name, val}
			out = append(out, o)
		}
	}

	sort.Sort(SortOutput(out))
	sout := ""
	for _, o := range out {
		sout += fmt.Sprintf("%d", o.val)
	}
    return sout
}

func printWires(wires map[string]int){
    xc := countWires(wires, 'x')
    for i := range xc {
        k := fmt.Sprintf("x%.2d", i)
        v := wires[k]
        fmt.Printf("%s: %d\n", k, v)
    }
    yc := countWires(wires, 'y')
    for i := range yc {
        k := fmt.Sprintf("y%.2d", i)
        v := wires[k]
        fmt.Printf("%s: %d\n", k, v)
    }
}


func toDecimal(bin string) int {
    dec := 0
    n := len(bin)-1
    for i := n; i >= 0; i-- {
        exp := n - i
        if bin[i] == '1' {
            dec += pow(2, exp)
        }
    }
    return dec
}

func readInput() (map[string]int, []Instr) {
	file := readFile("input.txt")
	sp := bytes.Split(file, []byte{'\n', '\n'})
	up, bot := sp[0], sp[1]

	uplns := bytes.Split(up, []byte{'\n'})
	gates := make(map[string]int)
	for _, ln := range uplns {
		sp := bytes.Split(ln, []byte{':', ' '})
		k := strings.TrimSpace(string(sp[0]))
		v, _ := strconv.Atoi(string(sp[1][0]))
		gates[k] = v
	}

	botlns := bytes.Split(bot, []byte{'\n'})
	botlns = botlns[:len(botlns)-1]
	instrs := make([]Instr, len(botlns))
	for i, ln := range botlns {
		instr := Instr{}

		sp := bytes.Split(ln, []byte{'-', '>'})
		instr.out = strings.TrimSpace(string(sp[1]))

		op := sp[0]
		ps := bytes.Split(op, []byte{' '})
		instr.lef = strings.TrimSpace(string(ps[0]))
		instr.op = strings.TrimSpace(string(ps[1]))
		instr.rig = strings.TrimSpace(string(ps[2]))

		instrs[i] = instr
	}

	return gates, instrs
}

func readFile(name string) []byte {
	file, err := os.ReadFile(name)
	if err != nil {
		panic(fmt.Sprintf("readFile %s: %v", name, err))
	}
	return file
}

func pow(a, b int) int {
    out := 1
    for b > 0 {
        out *= a
        b--
    }
    return out
}
