package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x int
	y int
}

type FloorType int
type ActionType int
type Direction int

const (
	Empty FloorType = iota
	Wall
)

const (
	ActRight ActionType = iota
	ActLeft
)

const (
	DirRight Direction = iota
	DirDown
	DirLeft
	DirUp
)

type MapData struct {
	width  int
	height int
	data   map[Point]FloorType

	startingPoint Point

	rowStarts    []int
	rowEnds      []int
	columnStarts []int
	columnEnds   []int

	sides []Side
}

type Side struct {
	id   int
	x    int
	y    int
	size int
}

func (s *Side) contains(p Point) bool {
	return s.x <= p.x && p.x < s.x+s.size && s.y <= p.y && p.y < s.y+s.size
}

type Action interface {
	sealed()
}

type Move struct {
	steps int
}

func (m *Move) sealed() {}

type Turn struct {
	action ActionType
}

func (m *Turn) sealed() {}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func parseMap(lines []string) MapData {
	mapData := MapData{
		data:          map[Point]FloorType{},
		height:        len(lines),
		startingPoint: Point{math.MaxInt, 0},
	}

	for _, l := range lines {
		mapData.width = max(mapData.width, len(l))
	}

	mapData.rowStarts = make([]int, mapData.height)
	mapData.rowEnds = make([]int, mapData.height)
	for i := range mapData.rowStarts {
		mapData.rowStarts[i] = math.MaxInt
		mapData.rowEnds[i] = math.MinInt
	}
	mapData.columnStarts = make([]int, mapData.width)
	mapData.columnEnds = make([]int, mapData.width)
	for i := range mapData.columnStarts {
		mapData.columnStarts[i] = math.MaxInt
		mapData.columnEnds[i] = math.MinInt
	}

	for y, line := range lines {
		for x, c := range line {
			if c == '.' {
				if y == 0 {
					mapData.startingPoint.x = min(mapData.startingPoint.x, x)
				}
				mapData.data[Point{x, y}] = Empty
			} else if c == '#' {
				mapData.data[Point{x, y}] = Wall
			}

			if c != ' ' {
				mapData.rowStarts[y] = min(mapData.rowStarts[y], x)
			}

			if c != ' ' && mapData.columnStarts[x] == math.MaxInt {
				mapData.columnStarts[x] = y
			} else if y == 0 && c != ' ' {
				mapData.columnStarts[x] = 0
			}
			if c == ' ' && mapData.columnStarts[x] != math.MaxInt && mapData.columnEnds[x] == math.MinInt {
				mapData.columnEnds[x] = y - 1
			}

			mapData.rowEnds[y] = x
		}

		if len(line) < mapData.width {
			for x := len(line); x < mapData.width; x++ {
				if mapData.columnStarts[x] != math.MaxInt && mapData.columnEnds[x] == math.MinInt {
					mapData.columnEnds[x] = y - 1
				}
			}
		}
	}

	for i, v := range mapData.columnEnds {
		if v == math.MinInt {
			mapData.columnEnds[i] = mapData.height - 1
		}
	}

	if mapData.width == 16 {
		// Example data
		mapData.sides = nil
	} else {
		// My data
		mapData.sides = []Side{
			{0, 50, 0, 50},
			{1, 0, 150, 50},
			{2, 0, 100, 50},
			{3, 50, 50, 50},
			{4, 50, 100, 50},
			{5, 100, 0, 50},
		}
	}

	return mapData
}

func (m *MapData) locate(pt Point) *Side {
	for i := range m.sides {
		if m.sides[i].contains(pt) {
			return &m.sides[i]
		}
	}

	return nil
}

func parseActions(line string) []Action {
	var result []Action

	for i := 0; i < len(line); i++ {
		if line[i] == 'L' {
			result = append(result, &Turn{ActLeft})
		} else if line[i] == 'R' {
			result = append(result, &Turn{ActRight})
		} else {
			j := i
			for ; j < len(line) && line[j] != 'L' && line[j] != 'R'; j++ {
			}
			v, _ := strconv.Atoi(line[i:j])
			result = append(result, &Move{v})
			i = j - 1
		}
	}

	return result
}

func applyTurn(dir Direction, act ActionType) Direction {
	if act == ActRight {
		return (dir + 1) % 4
	} else {
		if dir == 0 {
			return 3
		} else {
			return dir - 1
		}
	}
}

func (d Direction) toOffset() Point {
	if d == DirRight {
		return Point{1, 0}
	} else if d == DirDown {
		return Point{0, 1}
	} else if d == DirLeft {
		return Point{-1, 0}
	} else {
		return Point{0, -1}
	}
}

func (p Point) add(b Point) Point {
	return Point{p.x + b.x, p.y + b.y}
}

func applyMove(
	dir Direction,
	mapData *MapData,
	pos Point,
	steps int,
	step func(dir Direction, mapData *MapData, prevPos Point) (Point, Direction),
) (Point, Direction) {
	for ; steps > 0; steps-- {
		nextPos, newDir := step(dir, mapData, pos)

		floor, isFloor := mapData.data[nextPos]
		if !isFloor {
			fmt.Errorf("the alg is wrong, this should not happen :)")
		}

		if floor == Wall {
			break
		}
		pos = nextPos
		dir = newDir
	}

	return pos, dir
}

func stepPart1(dir Direction, mapData *MapData, prevPos Point) (Point, Direction) {
	nextPos := prevPos.add(dir.toOffset())
	if dir == DirLeft || dir == DirRight {
		if nextPos.x < mapData.rowStarts[nextPos.y] {
			nextPos.x = mapData.rowEnds[nextPos.y]
		} else if nextPos.x > mapData.rowEnds[nextPos.y] {
			nextPos.x = mapData.rowStarts[nextPos.y]
		}
	} else {
		if nextPos.y < mapData.columnStarts[nextPos.x] {
			nextPos.y = mapData.columnEnds[nextPos.x]
		} else if nextPos.y > mapData.columnEnds[nextPos.x] {
			nextPos.y = mapData.columnStarts[nextPos.x]
		}
	}
	return nextPos, dir
}

func (p Point) yAsXOffset(currSide, nextSide *Side) Point {
	return Point{p.y - currSide.y + nextSide.x, 0}
}

func (p Point) yAsYOffset(currSide, nextSide *Side) Point {
	return Point{0, p.y - currSide.y + nextSide.y}
}

func (p Point) xAsXOffset(currSide, nextSide *Side) Point {
	return Point{p.x - currSide.x + nextSide.x, 0}
}

func (p Point) xAsYOffset(currSide, nextSide *Side) Point {
	return Point{0, p.x - currSide.x + nextSide.y}
}

func (p Point) yFarAsXOffset(currSide, nextSide *Side) Point {
	return Point{currSide.y + currSide.size - 1 - p.y + nextSide.x, 0}
}

func (p Point) yFarAsYOffset(currSide, nextSide *Side) Point {
	return Point{0, currSide.y + currSide.size - 1 - p.y + nextSide.y}
}

func (p Point) xFarAsXOffset(currSide, nextSide *Side) Point {
	return Point{currSide.x + currSide.size - 1 - p.x + nextSide.x, 0}
}

func (p Point) xFarAsYOffset(currSide, nextSide *Side) Point {
	return Point{0, currSide.x + currSide.size - 1 - p.x + nextSide.y}
}

func (p Point) leftOf(side *Side) Point {
	return Point{side.x, p.y}
}

func (p Point) topOf(side *Side) Point {
	return Point{p.x, side.y}
}

func (p Point) rightOf(side *Side) Point {
	return Point{side.x + side.size - 1, p.y}
}

func (p Point) bottomOf(side *Side) Point {
	return Point{p.x, side.y + side.size - 1}
}

func stepPart2(dir Direction, mapData *MapData, currPos Point) (Point, Direction) {
	nextPos := currPos.add(dir.toOffset())

	currSide := mapData.locate(currPos)
	nextSide := mapData.locate(nextPos)

	toLeft := dir == DirLeft
	toRight := dir == DirRight
	toTop := dir == DirUp
	toBottom := dir == DirDown

	if nextSide == nil || nextSide.id != currSide.id {
		if currSide.id == 0 {
			if toTop {
				nextSide = &mapData.sides[1]
				dir = DirRight
				nextPos = nextPos.xAsYOffset(currSide, nextSide).leftOf(nextSide)
			}
			if toLeft {
				nextSide = &mapData.sides[2]
				dir = DirRight
				nextPos = nextPos.yFarAsYOffset(currSide, nextSide).leftOf(nextSide)
			}
		}

		if currSide.id == 1 {
			if toLeft {
				nextSide = &mapData.sides[0]
				dir = DirDown
				nextPos = nextPos.yAsXOffset(currSide, nextSide).topOf(nextSide)
			}
			if toBottom {
				nextSide = &mapData.sides[5]
				dir = DirDown
				nextPos = nextPos.xFarAsXOffset(currSide, nextSide).topOf(nextSide)
			}
			if toRight {
				nextSide = &mapData.sides[4]
				dir = DirUp
				nextPos = nextPos.yAsXOffset(currSide, nextSide).bottomOf(nextSide)
			}
		}

		if currSide.id == 2 {
			if toLeft {
				nextSide = &mapData.sides[0]
				dir = DirRight
				nextPos = nextPos.yFarAsYOffset(currSide, nextSide).leftOf(nextSide)
			}
			if toTop {
				nextSide = &mapData.sides[3]
				dir = DirRight
				nextPos = nextPos.xAsYOffset(currSide, nextSide).leftOf(nextSide)
			}
		}

		if currSide.id == 3 {
			if toLeft {
				nextSide = &mapData.sides[2]
				dir = DirDown
				nextPos = nextPos.yAsXOffset(currSide, nextSide).topOf(nextSide)
			}
			if toRight {
				nextSide = &mapData.sides[5]
				dir = DirDown
				nextPos = nextPos.yAsXOffset(currSide, nextSide).bottomOf(nextSide)
			}
		}

		if currSide.id == 4 {
			if toRight {
				nextSide = &mapData.sides[5]
				dir = DirLeft
				nextPos = nextPos.yFarAsYOffset(currSide, nextSide).rightOf(nextSide)
			}
			if toBottom {
				nextSide = &mapData.sides[1]
				dir = DirLeft
				nextPos = nextPos.xAsYOffset(currSide, nextSide).rightOf(nextSide)
			}
		}

		if currSide.id == 5 {
			if toTop {
				nextSide = &mapData.sides[1]
				dir = DirUp
				nextPos = nextPos.xFarAsXOffset(currSide, nextSide).bottomOf(nextSide)
			}
			if toRight {
				nextSide = &mapData.sides[4]
				dir = DirLeft
				nextPos = nextPos.yFarAsYOffset(currSide, nextSide).rightOf(nextSide)
			}
			if toBottom {
				nextSide = &mapData.sides[3]
				dir = DirLeft
				nextPos = nextPos.xAsYOffset(currSide, nextSide).rightOf(nextSide)
			}
		}
	}

	return nextPos, dir
}

func part1(mapData *MapData, actions []Action) int {
	direction := DirRight
	currentPos := mapData.startingPoint

	for _, a := range actions {
		if turn, isTurn := a.(*Turn); isTurn {
			direction = applyTurn(direction, turn.action)
		} else if move, isMove := a.(*Move); isMove {
			currentPos, direction = applyMove(direction, mapData, currentPos, move.steps, stepPart1)
		} else {
			fmt.Errorf("unknown action")
		}
	}

	return (currentPos.y+1)*1000 + (currentPos.x+1)*4 + int(direction)
}

func part2(mapData *MapData, actions []Action) int {
	direction := DirRight
	currentPos := mapData.startingPoint

	for _, a := range actions {
		if turn, isTurn := a.(*Turn); isTurn {
			direction = applyTurn(direction, turn.action)
		} else if move, isMove := a.(*Move); isMove {
			currentPos, direction = applyMove(direction, mapData, currentPos, move.steps, stepPart2)
		} else {
			fmt.Errorf("unknown action")
		}
	}

	return (currentPos.y+1)*1000 + (currentPos.x+1)*4 + int(direction)
}

func main() {
	rawContent, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(rawContent)
	lines := strings.Split(content, "\n")
	mapLines := lines[0 : len(lines)-2]
	actionLine := lines[len(lines)-1]

	mapData := parseMap(mapLines)
	actions := parseActions(actionLine)

	fmt.Println("Part 1: ", part1(&mapData, actions))
	fmt.Println("Part 2: ", part2(&mapData, actions))
}
