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
	Up Direction = iota
	Right
	Down
	Left
	Wait
)

type Blizzard struct {
	pos Point
	dir Direction
}

type BlizzardMap map[Point]struct{}

type MapData struct {
	blizzards []Blizzard
	mapTaken  BlizzardMap
	width     int
	height    int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func (p Point) add(x, y int) Point {
	return Point{p.x + x, p.y + y}
}

func (p Point) advance(dir Direction) Point {
	if dir == Up {
		return p.add(0, -1)
	} else if dir == Right {
		return p.add(1, 0)
	} else if dir == Down {
		return p.add(0, 1)
	} else if dir == Left {
		return p.add(-1, 0)
	} else {
		panic("in the pentagram")
	}
}

func (p Point) distanceTo(p2 Point) int {
	return abs(p.x-p2.x) + abs(p.y-p2.y)
}

func (b *Blizzard) advance() {
	b.pos = b.pos.advance(b.dir)
}

func (b *Blizzard) wrap(mapData *MapData) {
	if b.pos.x == 0 {
		b.pos.x = mapData.width - 2
	} else if b.pos.x == mapData.width-1 {
		b.pos.x = 1
	} else if b.pos.y == 0 {
		b.pos.y = mapData.height - 2
	} else if b.pos.y == mapData.height-1 {
		b.pos.y = 1
	}
}

func (mapData *MapData) doBlizzardStep() MapData {
	newBlizzards := make([]Blizzard, len(mapData.blizzards))
	newTaken := BlizzardMap{}

	for i, b := range mapData.blizzards {
		b.advance()
		b.wrap(mapData)
		newBlizzards[i] = b
		newTaken[b.pos] = struct{}{}
	}

	return MapData{newBlizzards, newTaken, mapData.width, mapData.height}
}

func (mapData *MapData) isValid(pt Point) bool {
	_, isTaken := mapData.mapTaken[pt]
	return pt == mapData.initialStart() || pt == mapData.finishLine() || (pt.x > 0 && pt.y > 0 && pt.x < mapData.width-1 && pt.y < mapData.height-1 && !isTaken)
}

func (mapData *MapData) initialStart() Point {
	return Point{1, 0}
}

func (mapData *MapData) finishLine() Point {
	return Point{mapData.width - 2, mapData.height - 1}
}

func parseMap(lines []string) MapData {
	var blizzards []Blizzard
	taken := BlizzardMap{}
	for y, line := range lines {
		for x, c := range line {
			if c == '<' {
				taken[Point{x, y}] = struct{}{}
				blizzards = append(blizzards, Blizzard{Point{x, y}, Left})
			} else if c == '^' {
				taken[Point{x, y}] = struct{}{}
				blizzards = append(blizzards, Blizzard{Point{x, y}, Up})
			} else if c == '>' {
				taken[Point{x, y}] = struct{}{}
				blizzards = append(blizzards, Blizzard{Point{x, y}, Right})
			} else if c == 'v' {
				taken[Point{x, y}] = struct{}{}
				blizzards = append(blizzards, Blizzard{Point{x, y}, Down})
			}
		}
	}
	return MapData{blizzards, taken, len(lines[0]), len(lines)}
}

type Step struct {
	mapData *MapData
	pos     Point
	steps   int
}

type Stats struct {
	pos   Point
	steps int
}

type Visited map[Stats]struct{}

func (mapData *MapData) goTo(startFrom, endAt Point) (int, *MapData) {
	steps := []Step{Step{mapData, startFrom, 0}}
	visited := Visited{}

	for len(steps) > 0 {
		step := steps[0]
		steps = steps[1:]

		if step.pos == endAt {
			return step.steps, step.mapData
		} else if !step.mapData.isValid(step.pos) {
			continue
		} else if _, beenThere := visited[Stats{step.pos, step.steps}]; beenThere {
			continue
		}

		visited[Stats{step.pos, step.steps}] = struct{}{}

		newMap := step.mapData.doBlizzardStep()

		steps = append(
			steps,
			Step{&newMap, step.pos.advance(Down), step.steps + 1},
			Step{&newMap, step.pos.advance(Right), step.steps + 1},
			Step{&newMap, step.pos, step.steps + 1},
			Step{&newMap, step.pos.advance(Up), step.steps + 1},
			Step{&newMap, step.pos.advance(Left), step.steps + 1},
		)
	}

	return math.MaxInt, nil
}

func part1(mapData *MapData) int {
	steps, _ := mapData.goTo(mapData.initialStart(), mapData.finishLine())
	return steps
}

func part2(map0 *MapData) int {
	steps1, map1 := map0.goTo(map0.initialStart(), map0.finishLine())
	steps2, map2 := map1.goTo(map1.finishLine(), map1.initialStart())
	steps3, _ := map2.goTo(map2.initialStart(), map2.finishLine())
	return steps1 + steps2 + steps3
}

func main() {
	rawContent, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(rawContent)
	lines := strings.Split(content, "\n")

	mapData := parseMap(lines)

	fmt.Println("Part 1: ", part1(&mapData))
	fmt.Println("Part 2: ", part2(&mapData))
}
