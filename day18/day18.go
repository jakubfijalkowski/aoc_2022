package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

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

type Cube struct {
	x int
	y int
	z int
}

func mkCube(c Cube, x, y, z int) Cube {
	return Cube{x: c.x + x, y: c.y + y, z: c.z + z}
}

func parseCube(l string) Cube {
	parts := strings.Split(l, ",")
	x, _ := strconv.Atoi(parts[0])
	y, _ := strconv.Atoi(parts[1])
	z, _ := strconv.Atoi(parts[2])
	return Cube{x: x, y: y, z: z}
}

func parseCubes(lines []string) map[Cube]struct{} {
	cubes := make(map[Cube]struct{}, len(lines))
	for _, l := range lines {
		cubes[parseCube(l)] = struct{}{}
	}
	return cubes
}

func isTakenOffset(c Cube, allCubes map[Cube]struct{}) int {
	_, isTakenByCube := allCubes[c]
	if isTakenByCube {
		return 1
	} else {
		return 0
	}
}

func isTaken(c Cube, cubes map[Cube]struct{}) bool {
	_, isTakenByCube := cubes[c]
	return isTakenByCube
}

func countEmptySides(c Cube, allCubes map[Cube]struct{}) int {
	return 6 -
		isTakenOffset(mkCube(c, 1, 0, 0), allCubes) -
		isTakenOffset(mkCube(c, -1, 0, 0), allCubes) -
		isTakenOffset(mkCube(c, 0, 1, 0), allCubes) -
		isTakenOffset(mkCube(c, 0, -1, 0), allCubes) -
		isTakenOffset(mkCube(c, 0, 0, 1), allCubes) -
		isTakenOffset(mkCube(c, 0, 0, -1), allCubes)
}

func countWaterSides(c Cube, water map[Cube]struct{}) int {
	return isTakenOffset(mkCube(c, 1, 0, 0), water) +
		isTakenOffset(mkCube(c, -1, 0, 0), water) +
		isTakenOffset(mkCube(c, 0, 1, 0), water) +
		isTakenOffset(mkCube(c, 0, -1, 0), water) +
		isTakenOffset(mkCube(c, 0, 0, 1), water) +
		isTakenOffset(mkCube(c, 0, 0, -1), water)
}

func addIfNotTaken(water map[Cube]struct{}, toVisit []Cube, current Cube, minX, maxX, minY, maxY, minZ, maxZ int) []Cube {
	if isTaken(current, water) || current.x < minX || current.x > maxX || current.y < minY || current.y > maxY || current.z < minZ || current.z > maxZ {
		return toVisit
	}
	return append(toVisit, current)
}

func flood(allCubes map[Cube]struct{}) map[Cube]struct{} {
	minX, maxX := math.MaxInt, math.MinInt
	minY, maxY := math.MaxInt, math.MinInt
	minZ, maxZ := math.MaxInt, math.MinInt
	for c := range allCubes {
		minX = min(minX, c.x-1)
		minY = min(minY, c.y-1)
		minZ = min(minZ, c.z-1)
		maxX = max(maxX, c.x+1)
		maxY = max(maxY, c.y+1)
		maxZ = max(maxZ, c.z+1)
	}

	water := make(map[Cube]struct{}, (maxX-minX)*(maxY-minY)*(maxZ-minZ))
	for c := range allCubes {
		water[c] = struct{}{}
	}

	toVisit := []Cube{{maxX, maxY, maxZ}}

	for len(toVisit) > 0 {
		current := toVisit[0]
		toVisit = toVisit[1:]
		if _, ok := water[current]; ok {
			continue
		}
		water[current] = struct{}{}

		toVisit = addIfNotTaken(water, toVisit, mkCube(current, 1, 0, 0), minX, maxX, minY, maxY, minZ, maxZ)
		toVisit = addIfNotTaken(water, toVisit, mkCube(current, -1, 0, 0), minX, maxX, minY, maxY, minZ, maxZ)
		toVisit = addIfNotTaken(water, toVisit, mkCube(current, 0, 1, 0), minX, maxX, minY, maxY, minZ, maxZ)
		toVisit = addIfNotTaken(water, toVisit, mkCube(current, 0, -1, 0), minX, maxX, minY, maxY, minZ, maxZ)
		toVisit = addIfNotTaken(water, toVisit, mkCube(current, 0, 0, 1), minX, maxX, minY, maxY, minZ, maxZ)
		toVisit = addIfNotTaken(water, toVisit, mkCube(current, 0, 0, -1), minX, maxX, minY, maxY, minZ, maxZ)
	}

	for c := range allCubes {
		delete(water, c)
	}

	return water
}

func part1(allCubes map[Cube]struct{}) int {
	total := 0
	for c := range allCubes {
		total += countEmptySides(c, allCubes)
	}
	return total
}

func part2(allCubes map[Cube]struct{}) int {
	total := 0
	water := flood(allCubes)
	for c := range allCubes {
		total += countWaterSides(c, water)
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
	cubes := parseCubes(lines)

	fmt.Println("Part 1: ", part1(cubes))
	fmt.Println("Part 2: ", part2(cubes))
}
