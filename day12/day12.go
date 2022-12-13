package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

type Pos struct {
	x int
	y int
}

func add(p1, p2 Pos) Pos {
	return Pos{
		x: p1.x + p2.x,
		y: p1.y + p2.y,
	}
}

type Map struct {
	elevation []int
	start     Pos
	end       Pos
	width     int
	height    int
}

func index(mapData *Map, x, y int) int {
	return y*mapData.width + x
}

func posIndex(mapData *Map, p Pos) int {
	return p.y*mapData.width + p.x
}

func inRange(mapData *Map, p Pos) bool {
	return p.x >= 0 && p.y >= 0 && p.x < mapData.width && p.y < mapData.height
}

func parse(lines []string) Map {
	mapData := Map{
		height: len(lines),
		width:  len(lines[0]),
	}
	mapData.elevation = make([]int, mapData.width*mapData.height)

	for y, l := range lines {
		runes := []rune(l)
		for x, c := range runes {
			if c == 'S' {
				mapData.elevation[index(&mapData, x, y)] = 0
				mapData.start = Pos{x: x, y: y}
			} else if c == 'E' {
				mapData.elevation[index(&mapData, x, y)] = 'z' - 'a'
				mapData.end = Pos{x: x, y: y}
			} else {
				mapData.elevation[index(&mapData, x, y)] = int(c - 'a')
			}
		}
	}
	return mapData
}

type ToCheck struct {
	steps int
	from  int
	pos   Pos
}

func findShortest(mapData *Map) int {
	queue := make([]ToCheck, 1)
	queue[0] = ToCheck{steps: 0, pos: mapData.start}

	visited := make([]int, mapData.width*mapData.height)
	for i := range visited {
		visited[i] = -1
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		idx := posIndex(mapData, curr.pos)
		if inRange(mapData, curr.pos) && visited[idx] == -1 && mapData.elevation[idx] <= curr.from+1 {
			visited[idx] = curr.steps
			queue = append(queue,
				ToCheck{steps: curr.steps + 1, pos: add(curr.pos, Pos{x: -1, y: 0}), from: mapData.elevation[idx]},
				ToCheck{steps: curr.steps + 1, pos: add(curr.pos, Pos{x: 1, y: 0}), from: mapData.elevation[idx]},
				ToCheck{steps: curr.steps + 1, pos: add(curr.pos, Pos{x: 0, y: -1}), from: mapData.elevation[idx]},
				ToCheck{steps: curr.steps + 1, pos: add(curr.pos, Pos{x: 0, y: 1}), from: mapData.elevation[idx]})
			if mapData.end == curr.pos {
				return visited[idx]
			}
		}
	}
	return -1
}

func part1(mapData *Map) int {
	return findShortest(mapData)
}

func part2(mapData *Map) int {
	min := math.MaxInt
	for y := 0; y < mapData.height; y++ {
		for x := 0; x < mapData.width; x++ {
			mapData.start = Pos{x: x, y: y}
			if mapData.elevation[posIndex(mapData, mapData.start)] == 0 {
				shortest := findShortest(mapData)
				if shortest > -1 && shortest < min {
					min = shortest
				}
			}
		}
	}
	return min
}

func main() {
	rawContent, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(rawContent)
	lines := strings.Split(content, "\n")

	mapData := parse(lines)

	fmt.Println("Part 1: ", part1(&mapData))
	fmt.Println("Part 2: ", part2(&mapData))
}
