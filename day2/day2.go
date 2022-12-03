package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Item int

const (
	Rock     Item = 1
	Paper    Item = 2
	Scissors Item = 3
)

func parse(d string) Item {
	if d == "A" || d == "X" {
		return Rock
	} else if d == "B" || d == "Y" {
		return Paper
	} else {
		return Scissors
	}
}

func beats(my, opponent Item) bool {
	return (my == Rock && opponent == Scissors) || (my-opponent) == 1
}

func roundScore(act []Item) int {
	my := act[1]
	opponent := act[0]
	if beats(my, opponent) {
		return int(my + 6)
	} else if my == opponent {
		return int(my + 3)
	} else {
		return int(my)
	}
}

func calcTotalScore(values [][]Item) int {
	totalScore := 0

	for _, p := range values {
		totalScore = totalScore + roundScore(p)
	}

	return totalScore
}

func substitute(values [][]Item) {
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
	rawContent, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(rawContent)
	lines := strings.Split(strings.TrimSpace(content), "\n")

	values := make([][]Item, 0, len(lines))

	for _, l := range lines {
		line := strings.Split(l, " ")
		v := []Item{parse(line[0]), parse(line[1])}
		values = append(values, v)
	}

	fmt.Println("Part 1:", calcTotalScore(values))
	substitute(values)
	fmt.Println("Part 2:", calcTotalScore(values))
}
