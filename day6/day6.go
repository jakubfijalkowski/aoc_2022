package main

import (
	"fmt"
	"log"
	"os"
)

func isSubstringUnique(data string, i, count int) bool {
	lastChars := make(map[byte]struct{})

	for j := 0; j < count; j++ {
		lastChars[data[i-j]] = struct{}{}
	}

	return len(lastChars) == count
}

func part1(data string) int {
	for i := 3; i < len(data); i++ {
		if isSubstringUnique(data, i, 4) {
			return i + 1
		}
	}

	return -1
}

func part2(data string) int {
	for i := 13; i < len(data); i++ {
		if isSubstringUnique(data, i, 14) {
			return i + 1
		}
	}

	return -1
}

func main() {

	rawContent, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(rawContent)

	fmt.Println("Part 1: ", part1(content))
	fmt.Println("Part 2: ", part2(content))
}
