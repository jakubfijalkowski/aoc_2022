package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type TaggedNumber struct {
	value int64
	index int64
}

func wrap(v int64, l int) int64 {
	l -= 1
	if v < 0 {
		mul := math.Ceil(float64(-v) / float64(l))
		v = v + int64(mul)*int64(l)
	}

	if v > int64(l) {
		v = v % int64(l)
	}

	return v
}

func printNumbers(numbers []TaggedNumber) {
	for _, n := range numbers {
		fmt.Printf("%v, ", n.value)
	}
	fmt.Println()
}

func findToMove(numbers []TaggedNumber, idx int64) int64 {
	for i, v := range numbers {
		if v.index == idx {
			return int64(i)
		}
	}
	panic("in the pentagram")
}

func mix(numbers []TaggedNumber) []TaggedNumber {
	for idx := int64(0); idx < int64(len(numbers)); idx++ {
		i := findToMove(numbers, idx)
		v := numbers[i]
		j := wrap(i+v.value, len(numbers))

		if j < i {
			a := numbers[:j]
			b := numbers[j:i]
			c := numbers[i+1:]

			newNumbers := make([]TaggedNumber, 0, len(numbers))
			newNumbers = append(newNumbers, a...)
			newNumbers = append(newNumbers, v)
			newNumbers = append(newNumbers, b...)
			newNumbers = append(newNumbers, c...)

			numbers = newNumbers
		} else if j > i {
			a := numbers[:i]
			b := numbers[i+1 : j+1]
			c := numbers[j+1:]

			newNumbers := make([]TaggedNumber, 0, len(numbers))
			newNumbers = append(newNumbers, a...)
			newNumbers = append(newNumbers, b...)
			newNumbers = append(newNumbers, v)
			newNumbers = append(newNumbers, c...)

			numbers = newNumbers
		}
	}
	return numbers
}

func findResult(numbers []TaggedNumber) int64 {
	zeroBase := 0
	for i, v := range numbers {
		if v.value == 0 {
			zeroBase = i
			break
		}
	}

	th1000 := numbers[(zeroBase+1000)%len(numbers)].value
	th2000 := numbers[(zeroBase+2000)%len(numbers)].value
	th3000 := numbers[(zeroBase+3000)%len(numbers)].value

	return th1000 + th2000 + th3000
}

func part1(numbers []TaggedNumber) int64 {
	numbers = mix(numbers)
	return findResult(numbers)
}

func part2(numbers []TaggedNumber) int64 {
	for i := 0; i < 10; i++ {
		numbers = mix(numbers)
	}
	return findResult(numbers)
}

func main() {
	rawContent, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(rawContent)
	lines := strings.Split(content, "\n")

	numbersPart1 := make([]TaggedNumber, len(lines))
	for i, l := range lines {
		v, _ := strconv.Atoi(l)
		numbersPart1[i] = TaggedNumber{int64(v), int64(i)}
	}

	numbersPart2 := make([]TaggedNumber, len(lines))
	for i, v := range numbersPart1 {
		numbersPart2[i] = TaggedNumber{v.value * 811589153, v.index}
	}

	fmt.Println("Part 1: ", part1(numbersPart1))
	fmt.Println("Part 2: ", part2(numbersPart2))
}
