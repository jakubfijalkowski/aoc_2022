package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

type Point struct {
	x int
	y int
}

type Direction int

const (
	North Direction = iota
	South
	West
	East

	NorthEast
	SouthEast
	SouthWest
	NorthWest
)

type Elf struct {
	id       int
	position Point
	proposed Point
}

type ElfLocations map[Point]struct{}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (p Point) move(dir Direction) Point {
	if dir == North {
		return Point{p.x, p.y - 1}
	} else if dir == NorthEast {
		return Point{p.x + 1, p.y - 1}
	} else if dir == East {
		return Point{p.x + 1, p.y}
	} else if dir == SouthEast {
		return Point{p.x + 1, p.y + 1}
	} else if dir == South {
		return Point{p.x, p.y + 1}
	} else if dir == SouthWest {
		return Point{p.x - 1, p.y + 1}
	} else if dir == West {
		return Point{p.x - 1, p.y}
	} else if dir == NorthWest {
		return Point{p.x - 1, p.y - 1}
	} else {
		panic("in the pentagram")
	}
}

func parseElves(lines []string) []Elf {
	var elves []Elf

	i := 0
	for y, l := range lines {
		for x, c := range l {
			if c == '#' {
				elves = append(elves, Elf{
					id:       i,
					position: Point{x, y},
				})
				i++
			}
		}
	}

	return elves
}

func toLocations(elves []Elf) ElfLocations {
	loc := ElfLocations{}
	for _, e := range elves {
		loc[e.position] = struct{}{}
	}
	return loc
}

func (loc *ElfLocations) isOccupied(p Point) bool {
	_, ok := (*loc)[p]
	return ok
}

type Proposal struct {
	elf         int
	shouldAbort bool
}

func getChecks(dir Direction) []Direction {
	if dir == North {
		return []Direction{North, NorthEast, NorthWest}
	} else if dir == South {
		return []Direction{South, SouthEast, SouthWest}
	} else if dir == West {
		return []Direction{West, NorthWest, SouthWest}
	} else if dir == East {
		return []Direction{East, NorthEast, SouthEast}
	}
	return nil
}

func diffuse(elves []Elf, startingDirection Direction) bool {
	alreadyTaken := toLocations(elves)
	proposed := map[Point]Proposal{}

	for i := range elves {
		e := &elves[i]
		e.proposed = e.position
		doesMove := false

		if !alreadyTaken.isOccupied(e.position.move(North)) && !alreadyTaken.isOccupied(e.position.move(NorthEast)) && !alreadyTaken.isOccupied(e.position.move(East)) &&
			!alreadyTaken.isOccupied(e.position.move(SouthEast)) && !alreadyTaken.isOccupied(e.position.move(South)) && !alreadyTaken.isOccupied(e.position.move(SouthWest)) &&
			!alreadyTaken.isOccupied(e.position.move(West)) && !alreadyTaken.isOccupied(e.position.move(NorthWest)) {
			continue
		}

		for d := 0; d < 4; d++ {
			proposedDir := (startingDirection + Direction(d)) % 4
			checks := getChecks(proposedDir)
			if !alreadyTaken.isOccupied(e.position.move(checks[0])) && !alreadyTaken.isOccupied(e.position.move(checks[1])) && !alreadyTaken.isOccupied(e.position.move(checks[2])) {
				e.proposed = e.position.move(proposedDir)
				doesMove = true
				break
			}
		}

		if doesMove {
			prevProposal, newTaken := proposed[e.proposed]
			if newTaken {
				proposed[e.proposed] = Proposal{prevProposal.elf, true}
				e.proposed = e.position
			} else {
				proposed[e.proposed] = Proposal{i, false}
			}
		}
	}

	for _, p := range proposed {
		if p.shouldAbort {
			elves[p.elf].proposed = elves[p.elf].position
		}
	}

	anyMoved := false

	for i := range elves {
		if elves[i].position != elves[i].proposed {
			anyMoved = true
		}

		elves[i].position = elves[i].proposed
	}

	return anyMoved
}

func getSize(elves []Elf) (Point, Point) {
	minX, minY := math.MaxInt, math.MaxInt
	maxX, maxY := math.MinInt, math.MinInt

	for _, e := range elves {
		minX = min(minX, e.position.x)
		minY = min(minY, e.position.y)

		maxX = max(maxX, e.position.x)
		maxY = max(maxY, e.position.y)
	}

	return Point{minX, minY}, Point{maxX, maxY}
}

func print(elves []Elf) {
	topLeft, rightBottom := getSize(elves)
	locs := toLocations(elves)

	for y := topLeft.y; y <= rightBottom.y; y++ {
		for x := topLeft.x; x <= rightBottom.x; x++ {
			if locs.isOccupied(Point{x, y}) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func part1(elves []Elf) int {
	for i := 0; i < 10; i++ {
		diffuse(elves, Direction(i%4))
	}

	topLeft, rightBottom := getSize(elves)
	width := rightBottom.x - topLeft.x + 1
	height := rightBottom.y - topLeft.y + 1

	return width*height - len(elves)
}

func part2(elves []Elf) int {
	i := 0
	for ; diffuse(elves, Direction(i%4)); i++ {
	}
	return i + 1
}

func main() {
	rawContent, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(rawContent)
	lines := strings.Split(content, "\n")

	elves := parseElves(lines)

	fmt.Println("Part 1: ", part1(elves))

	elves = parseElves(lines)
	fmt.Println("Part 2: ", part2(elves))
}
