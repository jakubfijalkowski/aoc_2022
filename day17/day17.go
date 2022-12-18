package main

import (
	"fmt"
	"log"
	"os"
)

type Direction int

const (
	Left Direction = iota
	Right
)

const ChamberSize = 7

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func max(a, b int64) int64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a, b int64) int64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func parseMovements(data []byte) []Direction {
	result := make([]Direction, len(data))

	for i, c := range data {
		if c == '<' {
			result[i] = Left
		} else {
			result[i] = Right
		}
	}

	return result
}

type Brick struct {
	brickMask [][]bool
	sizeX     int64
	sizeY     int64
	x         int64
	y         int64
}

func produceBrick(t int64) Brick {
	if t%5 == 0 {
		return Brick{
			brickMask: [][]bool{{true, true, true, true}},
			sizeX:     4,
			sizeY:     1,
		}
	} else if t%5 == 1 {
		return Brick{
			brickMask: [][]bool{{false, true, false}, {true, true, true}, {false, true, false}},
			sizeX:     3,
			sizeY:     3,
		}
	} else if t%5 == 2 {
		return Brick{
			brickMask: [][]bool{{false, false, true}, {false, false, true}, {true, true, true}},
			sizeX:     3,
			sizeY:     3,
		}
	} else if t%5 == 3 {
		return Brick{
			brickMask: [][]bool{{true}, {true}, {true}, {true}},
			sizeX:     1,
			sizeY:     4,
		}
	} else if t%5 == 4 {
		return Brick{
			brickMask: [][]bool{{true, true}, {true, true}},
			sizeX:     2,
			sizeY:     2,
		}
	}
	panic("cannot happen")
}

func (b *Brick) toPositionInside(x, y int64) (int64, int64) {
	return x - b.x, y - b.y
}

func (b *Brick) overlaps(other *Brick) bool {
	if !(min(b.x+b.sizeX, other.x+other.sizeX) > max(b.x, other.x) && min(b.y+b.sizeY, other.y+other.sizeY) > max(b.y, other.y)) {
		return false
	}

	for y := int64(0); y < b.sizeY; y++ {
		for x := int64(0); x < b.sizeX; x++ {
			otherX, otherY := other.toPositionInside(b.x+x, b.y+y)
			if otherY >= 0 && otherX >= 0 && otherY < other.sizeY && otherX < other.sizeX {
				if b.brickMask[y][x] && other.brickMask[otherY][otherX] {
					return true
				}
			}
		}
	}

	return false
}

func (b *Brick) move(d Direction) {
	if d == Left && b.x > 0 {
		b.x--
	} else if d == Right && b.x+b.sizeX < ChamberSize {
		b.x++
	}
}

func (b *Brick) fall() {
	b.y++
}

func (b *Brick) moveTo(x, y int64) {
	b.x = x
	b.y = y
}

type Chamber struct {
	bricks      []Brick
	topBrickY   int64
	movementIdx int64
}

func createChamber() Chamber {
	return Chamber{
		bricks:    []Brick{},
		topBrickY: 0,
	}
}

func (c *Chamber) isStopped(b *Brick) bool {
	if b.y+b.sizeY > 0 {
		return true
	}

	for i := len(c.bricks) - 1; i >= 0; i-- {
		if b.overlaps(&c.bricks[i]) {
			return true
		}
	}

	return false
}

func (chamber *Chamber) doBrick(movements []Direction, i int64) {
	brick := produceBrick(i)
	brick.moveTo(2, chamber.topBrickY-3-brick.sizeY)

	touchdown := false
	for !touchdown {
		orgX, orgY := brick.x, brick.y
		brick.move(movements[int(chamber.movementIdx%int64(len(movements)))])
		if chamber.isStopped(&brick) {
			brick.moveTo(orgX, orgY)
		}
		chamber.movementIdx++

		orgX, orgY = brick.x, brick.y
		brick.fall()
		if chamber.isStopped(&brick) {
			brick.moveTo(orgX, orgY)
			touchdown = true
		}
	}

	chamber.bricks = append(chamber.bricks, brick)
	chamber.topBrickY = min(chamber.topBrickY, brick.y)
}

func (chamber *Chamber) freeUpSpace() {
	if len(chamber.bricks) > 1000 {
		chamber.bricks = chamber.bricks[len(chamber.bricks)-1000:]
	}
}

func part1(movements []Direction) int64 {
	const blocks = 2022
	chamber := createChamber()

	for i := int64(0); i < blocks; i++ {
		chamber.doBrick(movements, i)
		chamber.freeUpSpace()
	}

	return -chamber.topBrickY
}

func findCycle(movements []Direction, pts []int64, targetBlocks int64) (int64, int64) {
	for offset := int64(0); offset < int64(len(pts)-len(movements)*2); offset++ {
		for jump := int64(1); offset+4*jump < int64(len(pts)); jump++ {
			offsetPts := pts[offset]
			baseCycle := pts[offset+jump] - offsetPts
			found := (targetBlocks-offset)%jump == 0

			for i := int64(1); found && i*jump+offset < int64(len(pts)); i++ {
				currPts := pts[i*jump+offset]
				if currPts != baseCycle*i+offsetPts {
					found = false
					break
				}
			}

			if found {
				return offset, jump
			}
		}
	}
	panic("cycle not found")
}

func part2(movements []Direction) int64 {
	const trialBlocks = 100000
	const targetBlocks = 1000000000000
	chamber := createChamber()
	pts := make([]int64, trialBlocks)

	for i := int64(0); i < trialBlocks; i++ {
		pts[i] = -chamber.topBrickY
		chamber.doBrick(movements, i)
		chamber.freeUpSpace()
	}

	offset, jump := findCycle(movements, pts, targetBlocks)

	multiply := (targetBlocks - offset) / jump
	base := pts[offset+jump] - pts[offset]
	total := pts[offset] + multiply*base

	return total
}

func main() {
	rawContent, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}

	movements := parseMovements(rawContent)

	fmt.Println("Part 1: ", part1(movements))
	fmt.Println("Part 2: ", part2(movements))
}
