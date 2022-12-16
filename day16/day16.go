package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Valve struct {
	flowRate int
	leadsTo  []string
	idx      int
}

type ValvePair struct {
	a string
	b string
}
type Valves map[string]*Valve
type Distances map[ValvePair]int
type Visited uint64

type Next struct {
	idx         string
	timeLeft    int
	accumulated int
	rate        int
	visited     Visited
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func parseValves(data string) Valves {
	valveRegex := regexp.MustCompile("Valve ([A-Z]+) has flow rate=(\\d+); tunnels? leads? to valves? (.+)")

	allMatches := valveRegex.FindAllStringSubmatch(data, -1)
	valves := make(Valves, len(allMatches))
	for i, m := range allMatches {
		flowRate, _ := strconv.Atoi(m[2])
		valves[m[1]] = &Valve{
			flowRate: flowRate,
			leadsTo:  strings.Split(m[3], ", "),
			idx:      i,
		}
	}

	return valves
}

func prune(valves Valves) Valves {
	for k, v := range valves {
		if v.flowRate == 0 {
			delete(valves, k)
		}
	}

	return valves
}

func mkPair(a, b string) ValvePair {
	return ValvePair{a: a, b: b}
}

func getOrInf(a, b string, distances Distances) int {
	d, ok := distances[mkPair(a, b)]
	if !ok {
		return math.MaxInt / 4
	} else {
		return d
	}
}

func calculateDistances(valves Valves) Distances {
	distances := make(Distances)

	for k, v := range valves {
		distances[mkPair(k, k)] = 0
		for _, e := range v.leadsTo {
			distances[mkPair(k, e)] = 1
		}
	}

	for k := range valves {
		for i := range valves {
			for j := range valves {
				ij := getOrInf(i, j, distances)
				ik := getOrInf(i, k, distances)
				kj := getOrInf(k, j, distances)

				if ij > ik+kj {
					distances[mkPair(i, j)] = ik + kj
				}
			}
		}
	}

	return distances
}

func isVisited(v *Valve, visited Visited) bool {
	return visited&(1<<v.idx) != 0
}

func visit(v *Valve, visited Visited) Visited {
	return visited | (1 << v.idx)
}

func singleWorkerDfs(valves Valves, distances Distances, timeLeft int, visited Visited) int {
	best := 0

	stack := make([]Next, 1)
	stack[0] = Next{idx: "AA", timeLeft: timeLeft, visited: visited}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		terminal := true
		for idx, other := range valves {
			if !isVisited(other, current.visited) && distances[mkPair(current.idx, idx)]+1 < current.timeLeft {
				timeToOpen := distances[mkPair(current.idx, idx)] + 1

				stack = append(stack, Next{
					idx:         idx,
					timeLeft:    current.timeLeft - timeToOpen,
					accumulated: current.rate*timeToOpen + current.accumulated,
					rate:        current.rate + other.flowRate,
					visited:     visit(other, current.visited),
				})

				terminal = false
			}
		}

		if terminal {
			best = max(best, current.accumulated+current.rate*current.timeLeft)
		}
	}

	return best
}

func part1(valves Valves, distances Distances) int {
	return singleWorkerDfs(valves, distances, 30, 0)
}

func part2(valves Valves, distances Distances) int {
	best := 0

	stackMine := make([]Next, 1)
	stackMine[0] = Next{idx: "AA", timeLeft: 26, visited: 0}

	for len(stackMine) > 0 {
		current := stackMine[len(stackMine)-1]
		stackMine = stackMine[:len(stackMine)-1]

		for idx, other := range valves {
			if !isVisited(other, current.visited) && distances[mkPair(current.idx, idx)]+1 < current.timeLeft {
				timeToOpen := distances[mkPair(current.idx, idx)] + 1
				thisStats := Next{
					idx:         idx,
					timeLeft:    current.timeLeft - timeToOpen,
					accumulated: current.rate*timeToOpen + current.accumulated,
					rate:        current.rate + other.flowRate,
					visited:     visit(other, current.visited),
				}

				stackMine = append(stackMine, thisStats)

				elephantBest := singleWorkerDfs(valves, distances, 26, thisStats.visited)
				best = max(best, thisStats.accumulated+thisStats.rate*thisStats.timeLeft+elephantBest)
			}
		}
	}

	return best
}

func main() {
	rawContent, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(rawContent)
	valves := parseValves(content)
	distances := calculateDistances(valves)
	valves = prune(valves)

	fmt.Println("Part 1: ", part1(valves, distances))
	fmt.Println("Part 2: ", part2(valves, distances))
}
