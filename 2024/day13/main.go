package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

const PART = 2

type Machine struct {
	prize Vec2
	A     Vec2
	B     Vec2
}

type Vec2 struct {
	x int
	y int
}

func main() {
	macs := readMachines()
    total := 0
	for _, mac := range macs {
		m, n, err := eval(mac)
        if err == nil {
            total += tokens(m, n)
        }
	}
	fmt.Printf("min tokens = %d\n", total)
}

func eval(mac Machine) (m, n int, err error) {
	px, py := mac.prize.x, mac.prize.y
	ax, ay := mac.A.x, mac.A.y
	bx, by := mac.B.x, mac.B.y
	m = (py*bx - px*by) / (bx*ay - ax*by)
	n = (px - m*ax) / bx
    
    if (m*ax + n*bx != px) || (m*ay + n*by != py){
        m, n, err = 0, 0, fmt.Errorf("eval %+v: no solution", mac)
    }
	return
}

func tokens(m, n int) int {
	return 3*m + n
}

func readMachines() []Machine {
	file := readFile("input.txt")

	blocks := bytes.Split(file, []byte{'\n', '\n'})
	macs := make([]Machine, len(blocks))
	for i, bl := range blocks {
		lines := bytes.Split(bl, []byte{'\n'})
		if len(lines) > 3 {
			lines = lines[:3]
		}
		lines[0] = lines[0][10:]
		lines[1] = lines[1][10:]
		lines[2] = lines[2][7:]

		a := bytes.Split(lines[0], []byte{','})
		b := bytes.Split(lines[1], []byte{','})
		p := bytes.Split(lines[2], []byte{','})

		ax, _ := strconv.Atoi(string(bytes.TrimSpace(a[0])[2:]))
		ay, _ := strconv.Atoi(string(bytes.TrimSpace(a[1])[2:]))

		bx, _ := strconv.Atoi(string(bytes.TrimSpace(b[0])[2:]))
		by, _ := strconv.Atoi(string(bytes.TrimSpace(b[1])[2:]))

		px, _ := strconv.Atoi(string(bytes.TrimSpace(p[0])[2:]))
		py, _ := strconv.Atoi(string(bytes.TrimSpace(p[1])[2:]))
        if PART == 2 {
            px += 10000000000000
            py += 10000000000000
        }

		mac := Machine{
			A:     Vec2{ax, ay},
			B:     Vec2{bx, by},
			prize: Vec2{px, py},
		}
		macs[i] = mac
	}
	return macs
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
