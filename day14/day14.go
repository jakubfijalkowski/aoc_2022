package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type EntityType int
type MapData map[Point]EntityType

const (
	Rock EntityType = iota
	Air
	Sand
)

type Point struct {
	x int
	y int
}

func clamp(x, a, b int) int {
	if x < a {
		return a
	} else if x > b {
		return b
	} else {
		return x
	}
}

func (p *Point) move(x, y int) Point {
	return Point{x: p.x + x, y: p.y + y}
}

func parsePoint(p string) Point {
	parts := strings.Split(p, ",")
	x, _ := strconv.Atoi(parts[0])
	y, _ := strconv.Atoi(parts[1])
	return Point{x: x, y: y}
}

func parsePath(line string) []Point {
	rawPoints := strings.Split(line, " -> ")
	path := make([]Point, len(rawPoints))
	for i := range rawPoints {
		path[i] = parsePoint(rawPoints[i])
	}
	return path
}

func drawPath(path []Point, mapData MapData) {
	startPos := path[0]

	for i := 1; i < len(path); i++ {
		endPos := path[i]

		xOffset := clamp(endPos.x-startPos.x, -1, 1)
		yOffset := clamp(endPos.y-startPos.y, -1, 1)

		for startPos.x != (endPos.x+xOffset) || startPos.y != (endPos.y+yOffset) {
			mapData[startPos] = Rock

			startPos.x += xOffset
			startPos.y += yOffset
		}

		startPos = endPos
	}
}

func findCaveEnd(mapData MapData) int {
	maxY := 0
	for k := range mapData {
		if k.y > maxY {
			maxY = k.y
		}
	}
	return maxY
}

func canMoveTo(mapData MapData, pt Point) bool {
	_, ok := mapData[pt]
	return !ok
}

func countSand(mapData MapData) int {
	total := 0
	for _, v := range mapData {
		if v == Sand {
			total++
		}
	}
	return total
}

func part1(mapData MapData) int {
	caveEnd := findCaveEnd(mapData)

	for true {
		sand := Point{x: 500, y: 0}

		for sand.y < caveEnd {
			if canMoveTo(mapData, sand.move(0, 1)) {
				sand = sand.move(0, 1)
			} else if canMoveTo(mapData, sand.move(-1, 1)) {

				sand = sand.move(-1, 1)
			} else if canMoveTo(mapData, sand.move(1, 1)) {
				sand = sand.move(1, 1)
			} else {
				mapData[sand] = Sand
				break
			}
		}

		if sand.y == caveEnd {
			break
		}
	}

	return countSand(mapData)
}

func part2(mapData MapData) int {
	caveEnd := findCaveEnd(mapData) + 2

	for true {
		sand := Point{x: 500, y: 0}
		if !canMoveTo(mapData, sand) {
			break
		}

		for sand.y < caveEnd {
			if sand.y + 1 == caveEnd {
				mapData[sand] = Sand
				break
			} else if canMoveTo(mapData, sand.move(0, 1)) {
				sand = sand.move(0, 1)
			} else if canMoveTo(mapData, sand.move(-1, 1)) {
				sand = sand.move(-1, 1)
			} else if canMoveTo(mapData, sand.move(1, 1)) {
				sand = sand.move(1, 1)
			} else {
				mapData[sand] = Sand
				break
			}
		}
	}

	return countSand(mapData)
}

func main() {
	rawContent, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(rawContent)
	lines := strings.Split(content, "\n")

	mapData := make(MapData)
	for _, l := range lines {
		drawPath(parsePath(l), mapData)
	}

	fmt.Println("Part 1: ", part1(mapData))
	fmt.Println("Part 2: ", part2(mapData)) // We can re-use the map :)
}
