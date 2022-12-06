package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Cmd struct {
	count int
	from  int
	to    int
}

func countStacks(l string) int {
	return len(l) - strings.Count(l, " ")
}

func parseStacks(lines []string) [][]string {
	specLine := lines[len(lines)-1]
	lines = lines[:len(lines)-1]

	stacks := countStacks(specLine)

	result := make([][]string, stacks)
	for i := range result {
		result[i] = make([]string, 0, 1)
	}

	for i := len(lines) - 1; i >= 0; i-- {
		l := lines[i]
		for j := 0; j < stacks; j++ {
			if len(l) >= j*4+3 {
				spec := strings.TrimSpace(l[(j * 4):(j*4 + 3)])
				if spec != "" {
					result[j] = append(result[j], spec[1:2])
				}
			}
		}

	}

	return result
}

func parseCmd(l string) Cmd {
	parts := strings.Split(l, " ")
	count, _ := strconv.Atoi(parts[1])
	from, _ := strconv.Atoi(parts[3])
	to, _ := strconv.Atoi(parts[5])
	return Cmd{
		count: count,
		from:  from - 1,
		to:    to - 1,
	}
}

func parseAllCmds(lines []string) []Cmd {
	result := make([]Cmd, 0, len(lines))
	for _, l := range lines {
		result = append(result, parseCmd(l))
	}
	return result
}

func applyV1(cmd Cmd, stacks [][]string) {
	toMove := stacks[cmd.from][len(stacks[cmd.from])-cmd.count:]
	stacks[cmd.from] = stacks[cmd.from][:len(stacks[cmd.from])-cmd.count]
	for i, _ := range toMove {
		stacks[cmd.to] = append(stacks[cmd.to], toMove[len(toMove)-i-1])
	}
}

func applyV2(cmd Cmd, stacks [][]string) {
	toMove := stacks[cmd.from][len(stacks[cmd.from])-cmd.count:]
	stacks[cmd.from] = stacks[cmd.from][:len(stacks[cmd.from])-cmd.count]
	stacks[cmd.to] = append(stacks[cmd.to], toMove...)
}

func part1(stacks [][]string, cmds []Cmd) string {
	for _, c := range cmds {
		applyV1(c, stacks)
	}

	var ret string
	for _, s := range stacks {
		if len(s) > 0 {
			ret += s[len(s)-1]
		}
	}

	return ret
}

func part2(stacks [][]string, cmds []Cmd) string {
	for _, c := range cmds {
		applyV2(c, stacks)
	}

	var ret string
	for _, s := range stacks {
		if len(s) > 0 {
			ret += s[len(s)-1]
		}
	}

	return ret
}

func main() {
	rawContent, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(rawContent)
	lines := strings.Split(content, "\n\n")

	stacks := parseStacks(strings.Split(lines[0], "\n"))
	cmds := parseAllCmds(strings.Split(lines[1], "\n"))

	fmt.Println("Part 1: ", part1(stacks, cmds))

	stacks = parseStacks(strings.Split(lines[0], "\n"))
	fmt.Println("Part 2: ", part2(stacks, cmds))
}
