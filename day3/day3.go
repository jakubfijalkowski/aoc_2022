package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

func splitRucksack(line string) (string, string) {
	middle := len(line) / 2
	a := line[0:middle]
	b := line[middle:]
	return a, b
}

func makeRuneMap(a string) map[rune]struct{} {
	m := map[rune]struct{}{}
	for _, r := range []rune(a) {
		m[r] = struct{}{}
	}
	return m
}

func findCommon(a, b *map[rune]struct{}) rune {
	for k := range *a {
		_, ok := (*b)[k]
		if ok {
			return k
		}
	}
	panic("should not happen")
}

func findCommon3(a, b, c *map[rune]struct{}) rune {
	for k := range *a {
		_, ok1 := (*b)[k]
		_, ok2 := (*c)[k]

		if ok1 && ok2 {
			return k
		}
	}
	panic("should not happen")
}

func score(r rune) int {
	if unicode.IsUpper(r) {
		return int(r-'A') + 27
	} else {
		return int(r-'a') + 1
	}
}

func scoreDifference(a, b string) int {
	aMap := makeRuneMap(a)
	bMap := makeRuneMap(b)
	diff := findCommon(&aMap, &bMap)
	return score(diff)
}

func findBadge(a, b, c string) int {
	aMap := makeRuneMap(a)
	bMap := makeRuneMap(b)
	cMap := makeRuneMap(c)
	common := findCommon3(&aMap, &bMap, &cMap)
	return score(common)
}

func part1(lines []string) int {
	total := 0
	for _, line := range lines {
		a, b := splitRucksack(line)
		total = total + scoreDifference(a, b)
	}

	return total
}

func part2(lines []string) int {
	total := 0
	for i := 0; i < len(lines); i = i + 3 {
		badge := findBadge(lines[i], lines[i+1], lines[i+2])
		total = total + badge
	}
	return total
}

func main() {
	rawContent, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(rawContent)
	lines := strings.Split(strings.TrimSpace(content), "\n")

	fmt.Println("Part 1: ", part1(lines))
	fmt.Println("Part 2: ", part2(lines))
}
