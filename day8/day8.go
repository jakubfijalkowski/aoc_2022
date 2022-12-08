package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func parse(lines []string) [][]int {
	result := make([][]int, len(lines))
	for i, l := range lines {
		single := make([]int, len(l))

		for j, c := range l {
			single[j] = int(c - '0')
		}

		result[i] = single
	}
	return result
}

func isVisible(plantMap [][]int, x, y, oX, oY int) bool {
	initial := plantMap[y][x]
	x += oX
	y += oY
	for x >= 0 && y >= 0 && x < len(plantMap[0]) && y < len(plantMap) {
		current := plantMap[y][x]
		if current >= initial {
			return false
		}

		x += oX
		y += oY
	}
	return true
}

func howManyCanItSee(plantMap [][]int, x, y, oX, oY int) int {
	see := 0
	initial := plantMap[y][x]
	x += oX
	y += oY
	for x >= 0 && y >= 0 && x < len(plantMap[0]) && y < len(plantMap) {
		current := plantMap[y][x]
		see++
		if current >= initial {
			return see
		}

		x += oX
		y += oY
	}
	return see
}

func scenicScore(plantMap [][]int, x, y int) int {
	return howManyCanItSee(plantMap, x, y, -1, 0) * howManyCanItSee(plantMap, x, y, 1, 0) * howManyCanItSee(plantMap, x, y, 0, -1) * howManyCanItSee(plantMap, x, y, 0, 1)
}

func part1(plantMap [][]int) int {
	total := len(plantMap)*2 + len(plantMap[0])*2 - 4

	for y := 1; y < len(plantMap)-1; y++ {
		for x := 1; x < len(plantMap[0])-1; x++ {
			if isVisible(plantMap, x, y, -1, 0) || isVisible(plantMap, x, y, 1, 0) || isVisible(plantMap, x, y, 0, -1) || isVisible(plantMap, x, y, 0, 1) {
				total++
			}
		}
	}

	return total
}

func part2(plantMap [][]int) int {
	total := 0

	for y := 1; y < len(plantMap)-1; y++ {
		for x := 1; x < len(plantMap[0])-1; x++ {
			score := scenicScore(plantMap, x, y)
			if score > total {
				total = score
			}
		}
	}

	return total
}

func main() {
	rawContent, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(rawContent)
	lines := strings.Split(content, "\n")

	plantMap := parse(lines)

	fmt.Println("Part 1: ", part1(plantMap))
	fmt.Println("Part 2: ", part2(plantMap))

}
