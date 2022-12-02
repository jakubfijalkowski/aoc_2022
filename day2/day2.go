package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func parse(d string) int {
	if d == "A" || d == "X" {
		return 1
	} else if d == "B" || d == "Y" {
		return 2
	} else {
		return 3
	}
}

func roundScore(act []int) int {
	won := (act[1] == 1 && act[0] == 3) || (act[1]-act[0]) == 1
	draw := act[1] == act[0]
	if won {
		return act[1] + 6
	} else if draw {
		return act[1] + 3
	} else {
		return act[1]
	}
}

func calcTotalScore(values [][]int) int {
	totalScore := 0

	for _, p := range values {
		totalScore = totalScore + roundScore(p)
	}

	return totalScore
}

func substitute(values [][]int) {
	for i := 0; i < len(values); i++ {
		if values[i][1] == 1 {
			values[i][1] = values[i][0] - 1
		} else if values[i][1] == 2 {
			values[i][1] = values[i][0]
		} else {
			values[i][1] = values[i][0] + 1
		}

		if values[i][1] == 0 {
			values[i][1] = 3
		} else if values[i][1] == 4 {
			values[i][1] = 1
		}
	}
}

func main() {
	rawContent, err := ioutil.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(rawContent)
	lines := strings.Split(strings.TrimSpace(content), "\n")

	values := make([][]int, 0, len(lines))

	for _, l := range lines {
		line := strings.Split(l, " ")
		v := []int{parse(line[0]), parse(line[1])}
		values = append(values, v)
	}

	fmt.Println("Part 1:", calcTotalScore(values))
	substitute(values)
	fmt.Println("Part 2:", calcTotalScore(values))
}
