package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Element interface {
	asArray() *Array
	print()
}

type Value struct {
	value int
}

type Array struct {
	elements []Element
}

func (v *Value) asArray() *Array {
	elements := make([]Element, 1)
	elements[0] = v
	return &Array{elements: elements}
}

func (v *Value) print() {
	fmt.Print(v.value)
}

func (v *Array) asArray() *Array {
	return v
}
func (v *Array) print() {
	fmt.Print("[")
	for i := range v.elements {
		if i > 0 {
			fmt.Print(", ")
		}
		v.elements[i].print()
	}
	fmt.Print("]")
}

func mkArray() Array {
	return Array{elements: make([]Element, 0)}
}

func mkDivider(v int) Array {
	arr1 := mkArray()
	arr2 := mkArray()
	arr2.append(&Value{value: v})
	arr1.append(&arr2)
	return arr1
}

func isDecoder(e Element, expected int) bool {
	arr1, ok1 := e.(*Array)
	if !ok1 || len(arr1.elements) != 1 {
		return false
	}

	arr2, ok2 := arr1.elements[0].(*Array)

	if !ok2 || len(arr2.elements) != 1 {
		return false
	}

	v, ok2 := arr2.elements[0].(*Value)
	return ok2 && v.value == expected
}

func (v *Array) append(el Element) {
	v.elements = append(v.elements, el)
}

type Pair struct {
	a Array
	b Array
}

func parseArray(line string) Array {
	root := mkArray()

	stack := make([]*Array, 1)
	stack[0] = &root

	for i := 1; i < len(line); i++ {
		if line[i] == ']' {
			stack = stack[:len(stack)-1]
		} else if line[i] == '[' {
			newTop := mkArray()
			stack[len(stack)-1].append(&newTop)
			stack = append(stack, &newTop)
		} else if line[i] == ',' {
		} else {
			length := 0
			for line[i+length] != ',' && line[i+length] != ']' {
				length++
			}
			val, _ := strconv.Atoi(line[i : i+length])
			stack[len(stack)-1].append(&Value{value: val})
			i += length - 1
		}
	}

	return root
}

func parsePair(p string) Pair {
	lines := strings.Split(p, "\n")
	return Pair{a: parseArray(lines[0]), b: parseArray(lines[1])}
}

type Result int

const (
	Left  Result = -1
	Equal        = 0
	Right        = 1
)

func compare(a Element, b Element) Result {
	aVal, aOk := a.(*Value)
	bVal, bOk := b.(*Value)
	if aOk && bOk {
		if aVal.value < bVal.value {
			return Left
		} else if aVal.value > bVal.value {
			return Right
		} else {
			return Equal
		}
	}

	aArr := a.asArray()
	bArr := b.asArray()

	arrLen := len(aArr.elements)
	if len(bArr.elements) > arrLen {
		arrLen = len(bArr.elements)
	}

	for i := 0; i < arrLen; i++ {
		if i >= len(aArr.elements) {
			return Left
		} else if i >= len(bArr.elements) {
			return Right
		} else {
			cmp := compare(aArr.elements[i], bArr.elements[i])
			if cmp != Equal {
				return cmp
			}
		}
	}

	return Equal
}

func flatten(pairs []Pair) []Array {
	result := make([]Array, len(pairs)*2)
	for i, p := range pairs {
		result[i*2+0] = p.a
		result[i*2+1] = p.b
	}
	return result
}

func part1(pairs []Pair) int {
	total := 0
	for i, p := range pairs {
		if compare(&p.a, &p.b) == Left {
			total += i + 1
		}
	}
	return total
}

func part2(pairs []Pair) int {
	flat := flatten(pairs)
	flat = append(flat, mkDivider(2), mkDivider(6))

	sort.Slice(flat, func(i, j int) bool {
		return compare(&flat[i], &flat[j]) == Left
	})

	result := 1
	for i, a := range flat {
		if isDecoder(&a, 2) || isDecoder(&a, 6) {
			result *= i + 1
		}
	}

	return result
}

func main() {
	rawContent, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(rawContent)
	lines := strings.Split(content, "\n\n")
	pairs := make([]Pair, len(lines))
	for i, l := range lines {
		pairs[i] = parsePair(l)
	}

	fmt.Println("Part 1: ", part1(pairs))
	fmt.Println("Part 2: ", part2(pairs))
}
