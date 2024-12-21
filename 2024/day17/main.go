package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type Regs struct {
	A int
	B int
	C int
}

type OpCode int

const (
	adv = iota
	bxl
	bst
	jnz
	bxc
	out
	bdv
	cdv
)

type Instr struct {
	OpCode
	Operand int
}

type Program []Instr

func (p Program) Dump() string {
	str := ""
	for _, instr := range p {
		str += strconv.Itoa(int(instr.OpCode)) + strconv.Itoa(instr.Operand)
	}
	return str
}

func (p Program) String() string {
	str := ""
	for _, instr := range p {
		str += strconv.Itoa(int(instr.OpCode)) + "," + strconv.Itoa(instr.Operand) + ","
	}
	return str
}

func main() {
	regs, prog := readInstrs()
	//fmt.Printf("Part One: %s\n", run(regs, prog))

    _ = len(prog.String()) / 2 - 1

    //base := 263882790666240 // first that ends with 3, 0
    //base := 266631569736460 // first that ends with 5, 3, 0
    //base := 266906447643900 // first that ends with 5, 5, 3, 0
    //base := 266932217447952 // first that ends with 7, 5, 5, 3, 0
    //base := 266932351665680 // first that ends with 4, 7, 5, 5, 3, 0
    //base := 266932485883436 // first that ends with 4, 4, 7, 5, 5, 3, 0
    //base := 266932586547860 // first that ends with 1, 4, 4, 7, 5, 5, 3, 0
    //base := 266932601229464 // first that ends with 3, 1, 4, 4, 7, 5, 5, 3, 0
    //base := 266932601262232 // first that ends with 0, 3, 1, 4, 4, 7, 5, 5, 3, 0
    // THE ONE THAT GENERATED THE RIGHT ANSWER >>>>>>>>>> base := 266932601389673 // first that ends with 5, 0, 3, 1, 4, 4, 7, 5, 5, 3, 0
    //base := 266932601401168 // first that ends with 7, 5, 0, 3, 1, 4, 4, 7, 5, 5, 3, 0
    //base := 266932601404508 // first that ends with 3, 7, 5, 0, 3, 1, 4, 4, 7, 5, 5, 3, 0
    //base :=                 // first that ends with 1, 3, 7, 5, 0, 3, 1, 4, 4, 7, 5, 5, 3, 0
    //base :=                 // first that ends with 4, 1, 3, 7, 5, 0, 3, 1, 4, 4, 7, 5, 5, 3, 0
    //base := 266932601404774 // first that ends with 2, 4, 1, 3, 7, 5, 0, 3, 1, 4, 4, 7, 5, 5, 3, 0

    // fmt.Printf("A (base 10) | A (base 8) | A (binary) | Output\n")
    // for range 100000 {
    //     regs.A++
    //     output := run(regs, prog)
    //     fmt.Printf("%d| %s | %b | %s\n", regs.A, strconv.FormatInt(int64(regs.A), 8), regs.A, output)
    // }
    regs.A = 266932601404433
    fmt.Printf("\n---Test Run---\n")
    fmt.Printf("Registers: %+v\n", regs)
    fmt.Printf("Program: %+v\n", prog)
    fmt.Printf("Output: %s\n", run(regs, prog))
    fmt.Printf("--------------\n\n")


}

func alice(bs []byte) int {
	str := ""
	//fmt.Printf("******************ALICE*****************\n")
	for _, b := range bs {
		//fmt.Printf("    alicing byte: %d\n", b)
		str += strconv.Itoa(int(b))
		//fmt.Printf("    current string: %s\n", str)
	}
	r, err := strconv.Atoi(string(str))
	//fmt.Printf("final alice: %d\n", r)
	//fmt.Printf("****************************************\n")

	if err != nil {
		panic(fmt.Sprintf("alice: %s", err.Error()))
	}
	return r
}

func run(regs Regs, prog Program) string {
	output := ""
	ip := 0
	for ip < len(prog) {
		switch prog[ip].OpCode {
		case adv:
			regs.A = regs.A / pow(2, combo(regs, prog[ip].Operand))
		case bxl:
			regs.B = regs.B ^ prog[ip].Operand
		case bst:
			regs.B = combo(regs, prog[ip].Operand) % 8
		case jnz:
			if regs.A == 0 {
				break
			}
			ip = prog[ip].Operand
			continue
		case bxc:
			regs.B = regs.B ^ regs.C
		case out:
			output += strconv.Itoa(combo(regs, prog[ip].Operand)%8) + ","
		case bdv:
			regs.B = regs.A / pow(2, combo(regs, prog[ip].Operand))
		case cdv:
			regs.C = regs.A / pow(2, combo(regs, prog[ip].Operand))
		}
		ip++
	}
	return output
}

func combo(regs Regs, oprn int) int {
	if oprn >= 0 && oprn < 4 {
		return oprn
	}
	if oprn == 4 {
		return regs.A
	}
	if oprn == 5 {
		return regs.B
	}
	if oprn == 6 {
		return regs.C
	}
	panic(fmt.Sprintf("combo: invalid operand %d", oprn))
}

func readInstrs() (Regs, Program) {
	file := readFile("input.txt")
	split := bytes.Split(file, []byte{'\n', '\n'})
	top := split[0]
	bot := split[1]

	regs := bytes.Split(top, []byte{'\n'})
	registers := parseRegisters(regs)
	prog := parseProgram(bot)
	return registers, prog
}

func parseRegisters(regs [][]byte) Regs {
	var registers Regs
	Areg, _ := strconv.Atoi(string(bytes.TrimSpace((bytes.Split(regs[0], []byte{':'})[1]))))
	Breg, _ := strconv.Atoi(string(bytes.TrimSpace((bytes.Split(regs[1], []byte{':'})[1]))))
	Creg, _ := strconv.Atoi(string(bytes.TrimSpace((bytes.Split(regs[2], []byte{':'})[1]))))
	registers.A = Areg
	registers.B = Breg
	registers.C = Creg
	return registers
}

func parseProgram(prog []byte) Program {
	var program Program
	ops := bytes.TrimSpace(bytes.Split(prog, []byte{' '})[1])
	for i := 0; i < len(ops)-2; i += 4 {
		code, _ := strconv.Atoi(string(ops[i]))
		oprn, _ := strconv.Atoi(string(ops[i+2]))
		instr := Instr{
			OpCode(code),
			oprn,
		}
		program = append(program, instr)
	}
	return program
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

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func pow(base, exp int) int {
	result := 1
	for {
		if exp&1 == 1 {
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
