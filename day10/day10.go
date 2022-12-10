package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MicroopType int

const (
	Noop MicroopType = iota
	AddX
)

type Microop struct {
	opType MicroopType
	data   interface{}
}

type AddXData struct {
	operand int
}

func abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

func parse(l string, output []Microop) []Microop {
	data := strings.Split(l, " ")
	if data[0] == "noop" {
		output = append(output, Microop{opType: Noop, data: nil})
	} else if data[0] == "addx" {
		operand, _ := strconv.Atoi(data[1])
		output = append(output, Microop{opType: Noop, data: nil}, Microop{opType: AddX, data: AddXData{operand: operand}})
	} else {
		fmt.Errorf("the instruction %v is unknown", data[0])
	}
	return output
}

func execute(xReg int, op Microop) int {
	if op.opType == Noop {
		return xReg
	} else if op.opType == AddX {
		xData, _ := op.data.(AddXData)
		return xReg + xData.operand
	} else {
		fmt.Errorf("unknown operation")
		return 0
	}
}

func part1(ops []Microop) int {
	total := 0
	xReg := 1

	for i, op := range ops {
		i++

		if i == 20 || i == 60 || i == 100 || i == 140 || i == 180 || i == 220 {
			total += i * xReg
		}

		xReg = execute(xReg, op)
	}
	return total
}

func part2(ops []Microop) {
	xReg := 1

	for i, op := range ops {
		if abs(xReg-(i%40)) <= 1 {
			fmt.Print("#")
		} else {
			fmt.Print(".")
		}

		if i == 39 || ((i-39)%40 == 0) {
			fmt.Println()
		}

		xReg = execute(xReg, op)
	}
}

func main() {
	readFile, err := os.Open("data.txt")

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	ops := make([]Microop, 0)

	for fileScanner.Scan() {
		ops = parse(fileScanner.Text(), ops)
	}

	fmt.Println("Part 1: ", part1(ops))
	part2(ops)
}
