package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Direction int

const (
	Left Direction = iota
	Up
	Right
	Down
)

type Move struct {
	dir    Direction
	length int
}

type Pos struct {
	x int
	y int
}

func abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
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

func parse(l string) Move {
	split := strings.Split(l, " ")
	length, _ := strconv.Atoi(split[1])
	if split[0] == "R" {
		return Move{dir: Right, length: length}
	} else if split[0] == "U" {
		return Move{dir: Up, length: length}
	} else if split[0] == "D" {
		return Move{dir: Down, length: length}
	} else {
		return Move{dir: Left, length: length}
	}
}

func toOffset(m Move) (int, int) {
	if m.dir == Left {
		return -m.length, 0
	} else if m.dir == Up {
		return 0, m.length
	} else if m.dir == Right {
		return m.length, 0
	} else {
		return 0, -m.length
	}
}

func areAdjacent(h, t Pos) bool {
	return abs(h.x-t.x) <= 1 && abs(h.y-t.y) <= 1
}

func isCovering(h, t Pos) bool {
	return h.x == t.x && h.y == t.y
}

func getMove(h, t Pos) (int, int) {
	ox := clamp(h.x-t.x, -1, 1)
	oy := clamp(h.y-t.y, -1, 1)

	return ox, oy
}

func add(p Pos, x, y int) Pos {
	return Pos{
		x: p.x + x,
		y: p.y + y,
	}
}

func part1(moves []Move) int {
	visited := make(map[Pos]struct{})

	head := Pos{x: 0, y: 0}
	tail := Pos{x: 0, y: 0}
	visited[tail] = struct{}{}

	for _, m := range moves {
		ox, oy := toOffset(m)
		head = add(head, ox, oy)

		for !areAdjacent(head, tail) {
			ox, oy = getMove(head, tail)
			tail = add(tail, ox, oy)
			visited[tail] = struct{}{}
		}
	}

	return len(visited)
}

func part2(moves []Move) int {
	visited := make(map[Pos]struct{})

	snake := make([]Pos, 10)
	for i, _ := range snake {
		snake[i] = Pos{x: 0, y: 0}
	}
	visited[snake[len(snake)-1]] = struct{}{}

	for _, m := range moves {
		dx, dy := toOffset(m)
		newHead := add(snake[0], dx, dy)
		for !isCovering(snake[0], newHead) {
			hx, hy := getMove(newHead, snake[0])
			snake[0] = add(snake[0], hx, hy)

			for i := 1; i < len(snake); i++ {
				if !areAdjacent(snake[i-1], snake[i]) {
					ox, oy := getMove(snake[i-1], snake[i])
					snake[i] = add(snake[i], ox, oy)
				}
			}

			visited[snake[len(snake)-1]] = struct{}{}
		}
	}

	return len(visited)
}

func main() {
	rawContent, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(rawContent)
	lines := strings.Split(content, "\n")

	moves := make([]Move, len(lines))
	for i, l := range lines {
		moves[i] = parse(l)
	}

	fmt.Println("Part 1: ", part1(moves))
	fmt.Println("Part 2: ", part2(moves))
}
