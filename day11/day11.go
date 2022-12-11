package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Monkey struct {
	id          int
	items       []uint64
	op          string
	mulByOld    bool
	operand     uint64
	divisibleBy uint64
	ifTrue      int
	ifFalse     int

	inspected int
}

func parseMonkey(data []string) Monkey {
	id, _ := strconv.Atoi(data[1])
	items := strings.Split(data[2], ", ")
	parsedItems := make([]uint64, len(items))
	for i, n := range items {
		item, _ := strconv.Atoi(n)
		parsedItems[i] = uint64(item)
	}

	operand, operandFail := strconv.Atoi(data[4])
	divisibleBy, _ := strconv.Atoi(data[5])
	ifTrue, _ := strconv.Atoi(data[6])
	ifFalse, _ := strconv.Atoi(data[7])

	return Monkey{
		id:          id,
		items:       parsedItems,
		op:          data[3],
		mulByOld:    operandFail != nil,
		operand:     uint64(operand),
		divisibleBy: uint64(divisibleBy),
		ifTrue:      ifTrue,
		ifFalse:     ifFalse,

		inspected: 0,
	}
}

func applyOpToItem(monkey *Monkey, item uint64) uint64 {
	if monkey.op == "+" {
		if monkey.mulByOld {
			return item + item
		} else {
			return item + monkey.operand
		}
	} else {
		if monkey.mulByOld {
			return item * item
		} else {
			return item * monkey.operand
		}
	}
}

func findMaxInspects(monkeys []Monkey) (int, int) {
	idx, value := 0, 0

	for i, m := range monkeys {
		if value < m.inspected {
			idx = i
			value = m.inspected
		}
	}

	return idx, value
}

func calcCoprime(monkeys []Monkey) uint64 {
	var total uint64 = 1
	for _, m := range monkeys {
		total *= m.divisibleBy
	}
	return total
}

func manualMonkey(monkeys []Monkey, monkey *Monkey, divide bool, coprime uint64) {
	itemsToProcess := monkey.items
	monkey.items = make([]uint64, 0)

	for _, item := range itemsToProcess {
		item = applyOpToItem(monkey, item)
		item = item % coprime
		if divide {
			item = item / 3
		}

		if item%monkey.divisibleBy == 0 {
			monkeys[monkey.ifTrue].items = append(monkeys[monkey.ifTrue].items, item)
		} else {
			monkeys[monkey.ifFalse].items = append(monkeys[monkey.ifFalse].items, item)
		}
	}

	monkey.inspected += len(itemsToProcess)
}

func solve(monkeys []Monkey, rounds int, divide bool) uint64 {
	coprime := calcCoprime(monkeys)

	for r := 0; r < rounds; r++ {
		for i := range monkeys {
			manualMonkey(monkeys, &monkeys[i], divide, coprime)
		}
	}

	top1Idx, top1 := findMaxInspects(monkeys)
	monkeys[top1Idx].inspected = 0
	_, top2 := findMaxInspects(monkeys)

	return uint64(top1) * uint64(top2)
}

func part1(monkeys []Monkey) uint64 {
	return solve(monkeys, 20, true)
}

func part2(monkeys []Monkey) uint64 {
	return solve(monkeys, 10000, false)
}

func main() {
	rawContent, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(rawContent)

	monkeyRegex := regexp.MustCompile("(?m)Monkey (\\d+):\\n\\s+Starting items: ([\\d ,]+)$\\n\\s+Operation: new = old (.) (.+)\\n\\s+ Test: divisible by (\\d+)\\n\\s+If true: throw to monkey (\\d+)\\n\\s+If false: throw to monkey (\\d+)")
	monkeyMatches := monkeyRegex.FindAllStringSubmatch(content, -1)

	monkeys := make([]Monkey, len(monkeyMatches))
	for i, m := range monkeyMatches {
		monkeys[i] = parseMonkey(m)
	}

	fmt.Println("Part 1: ", part1(monkeys))

	for i, m := range monkeyMatches {
		monkeys[i] = parseMonkey(m)
	}
	fmt.Println("Part 2: ", part2(monkeys))
}
