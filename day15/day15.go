package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
)

type Point struct {
	x int
	y int
}

type Sensor struct {
	sensor      Point
	beacon      Point
	sensorRange int
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func (a *Point) distanceTo(b Point) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func (a *Sensor) covers(b Point) bool {
	return a.sensor.distanceTo(b) <= a.sensorRange
}

func (a *Sensor) maxXAt(y int) int {
	return a.sensor.x + a.sensorRange - abs(a.sensor.y-y)
}

func parseAll(data string) []Sensor {
	sensorRegex := regexp.MustCompile("Sensor at x=(-?\\d+), y=(-?\\d+): closest beacon is at x=(-?\\d+), y=(-?\\d+)")

	matches := sensorRegex.FindAllStringSubmatch(data, -1)
	sensors := make([]Sensor, len(matches))

	for i, single := range matches {
		sx, _ := strconv.Atoi(single[1])
		sy, _ := strconv.Atoi(single[2])
		bx, _ := strconv.Atoi(single[3])
		by, _ := strconv.Atoi(single[4])

		sensor := Point{x: sx, y: sy}
		beacon := Point{x: bx, y: by}
		sensors[i] = Sensor{
			sensor:      sensor,
			beacon:      beacon,
			sensorRange: sensor.distanceTo(beacon),
		}
	}

	return sensors
}

func isCoveredByAny(p Point, sensors []Sensor) bool {
	for _, s := range sensors {
		if s.covers(p) {
			return true
		}
	}

	return false
}

func xMinMax(sensors []Sensor) (int, int) {
	min, max := math.MaxInt, math.MinInt
	for _, s := range sensors {
		if s.sensor.x-s.sensorRange < min {
			min = s.sensor.x - s.sensorRange
		}

		if s.sensor.x+s.sensorRange > max {
			max = s.sensor.x + s.sensorRange
		}
	}

	return min, max
}

func countBeaconsOn(sensors []Sensor, line int) int {
	uniqueBeacons := make(map[Point]struct{})
	for _, s := range sensors {
		uniqueBeacons[s.beacon] = struct{}{}
	}

	total := 0
	for b := range uniqueBeacons {
		if b.y == line {
			total++
		}
	}
	return total
}

func sortSensorsByXMin(sensors []Sensor) {
	sort.Slice(sensors, func(i, j int) bool {
		left := sensors[i].sensor.x - sensors[i].sensorRange
		right := sensors[j].sensor.x - sensors[j].sensorRange
		return left < right
	})
}

func part1(sensors []Sensor) int {
	const row = 2000000

	min, max := xMinMax(sensors)
	result := 0

	for x := min; x <= max; x++ {
		if isCoveredByAny(Point{x: x, y: row}, sensors) {
			result++
		}
	}

	return result - countBeaconsOn(sensors, row)
}

func part2(sensors []Sensor) int64 {
	const maxRange = 4000000

	for y := 0; y <= maxRange; y++ {
		x := 0
		moved := true
		for moved {
			moved = false
			for _, s := range sensors {
				if s.covers(Point{x: x, y: y}) {
					if s.maxXAt(y)+1 > x {
						x = s.maxXAt(y) + 1
						moved = true
					}
				}
			}
		}

		if x < maxRange {
			return int64(x)*int64(maxRange) + int64(y)
		}
	}

	return -1
}

func main() {
	rawContent, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(rawContent)

	sensors := parseAll(content)
	sortSensorsByXMin(sensors)

	fmt.Println("Part 1: ", part1(sensors))
	fmt.Println("Part 2: ", part2(sensors))
}
