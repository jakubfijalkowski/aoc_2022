package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type ElfRange struct {
	from int
	to   int
}

type RangePair struct {
	a ElfRange
	b ElfRange
}

func parseRange(r string) ElfRange {
	split := strings.Split(r, "-")
	from, _ := strconv.Atoi(split[0])
	to, _ := strconv.Atoi(split[1])
	return ElfRange{
		from: from,
		to:   to,
	}
}

func contains(a, b ElfRange) bool {
	return a.from <= b.from && b.from <= a.to &&
		a.from <= b.to && b.to <= a.to
}

func overlaps(p RangePair) bool {
	return p.a.from <= p.b.to && p.a.to >= p.b.from
}

func parsePair(l string) RangePair {
	split := strings.Split(l, ",")
	return RangePair{a: parseRange(split[0]), b: parseRange(split[1])}
}

func part1(pairs []RangePair) int {
	total := 0
	for _, r := range pairs {
		if contains(r.a, r.b) || contains(r.b, r.a) {
			total = total + 1
		}
	}
	return total
}

func part2(pairs []RangePair) int {
	total := 0
	for _, r := range pairs {
		if overlaps(r) {
			total = total + 1
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
	lines := strings.Split(strings.TrimSpace(content), "\n")
	pairs := make([]RangePair, 0, len(lines))
	for _, l := range lines {
		pairs = append(pairs, parsePair(l))
	}

	fmt.Println("Part 1:", part1(pairs))
	fmt.Println("Part 2:", part2(pairs))
}
