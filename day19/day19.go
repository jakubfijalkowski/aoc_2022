package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"
	"sync/atomic"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func parse(l string) int {
	val, err := strconv.Atoi(l)
	if err != nil {
		fmt.Errorf("cannot parse number from %v: %v", l, err)
	}
	return val
}

type Material struct {
	ore      int
	clay     int
	obsidian int
	geode    int
}

type Robot struct {
	cost       Material
	production Material
}

type Blueprint struct {
	id     int
	robots []Robot
}

type Factory struct {
	materials  Material
	production Material
	inBuild    Material
}

type Cache struct {
	a uint64
	b uint64
}

func mkCache(f Factory, time int) Cache {
	a := uint64(0)
	a = a*10000 + uint64(f.materials.ore)
	a = a*10000 + uint64(f.materials.clay)
	a = a*10000 + uint64(f.materials.obsidian)
	a = a*10000 + uint64(f.materials.geode)

	b := uint64(0)
	b = b*10000 + uint64(f.production.ore)
	b = b*10000 + uint64(f.production.clay)
	b = b*10000 + uint64(f.production.obsidian)
	b = b*10000 + uint64(f.production.geode)
	b = b*1000 + uint64(time)
	return Cache{a, b}
}

func (a Material) add(b Material) Material {
	return Material{a.ore + b.ore, a.clay + b.clay, a.obsidian + b.obsidian, a.geode + b.geode}
}

func (a Material) sub(b Material) Material {
	return Material{a.ore - b.ore, a.clay - b.clay, a.obsidian - b.obsidian, a.geode - b.geode}
}

func (a Material) isValid() bool {
	return a.ore >= 0 && a.clay >= 0 && a.obsidian >= 0 && a.geode >= 0
}

func (f Factory) canProduce(r Robot) bool {
	return f.materials.sub(r.cost).isValid()
}

func makesSenseToProduce(amount int, b Blueprint, what func(r Robot) int) bool {
	for _, r := range b.robots {
		if amount < what(r) {
			return true
		}
	}
	return false
}

func (f Factory) makesSenseToProduce(r Robot, b Blueprint) bool {
	if r.production.ore > 0 {
		return makesSenseToProduce(f.production.ore, b, func(r Robot) int { return r.cost.ore })
	} else if r.production.clay > 0 {
		return makesSenseToProduce(f.production.clay, b, func(r Robot) int { return r.cost.clay })
	} else if r.production.obsidian > 0 {
		return makesSenseToProduce(f.production.obsidian, b, func(r Robot) int { return r.cost.obsidian })
	} else {
		return true
	}
}

func (f Factory) produce(r Robot) Factory {
	if !f.canProduce(r) {
		panic("non-negative material invariant")
	}

	return Factory{f.materials.sub(r.cost), f.production, f.inBuild.add(r.production)}
}

func (f Factory) accumulate() Factory {
	return Factory{f.materials.add(f.production), f.production.add(f.inBuild), Material{0, 0, 0, 0}}
}

func initFactory() Factory {
	return Factory{production: Material{1, 0, 0, 0}}
}

func maximizeGeodes(blueprint Blueprint, time int, factory Factory, lastMax *int, cache map[Cache]struct{}) int {
	cacheEntry := mkCache(factory, time)
	if time == 1 {
		*lastMax = max(*lastMax, factory.materials.geode+factory.production.geode)
		return factory.materials.geode + factory.production.geode
	} else if _, ok := cache[cacheEntry]; ok {
		return 0
	}
	cache[cacheEntry] = struct{}{}

	// Place an upper bound on how many geodes we can produce - basically from now on we produce only geode robots
	maxRobots := factory.production.geode + time
	maxToProduce := time * (factory.production.geode + maxRobots) / 2
	if factory.materials.geode+maxToProduce < *lastMax {
		return 0
	}

	currentMax := maximizeGeodes(blueprint, time-1, factory.accumulate(), lastMax, cache) // Simulate "do not produce any robot"
	*lastMax = max(*lastMax, currentMax)

	for i := len(blueprint.robots) - 1; i >= 0; i-- {
		r := blueprint.robots[i]
		if factory.canProduce(r) && factory.makesSenseToProduce(r, blueprint) {
			currentMax = max(
				currentMax,
				maximizeGeodes(blueprint, time-1, factory.produce(r).accumulate(), lastMax, cache),
			)
			*lastMax = max(*lastMax, currentMax)
		}
	}

	*lastMax = max(*lastMax, currentMax)
	return currentMax
}

func doOne(b Blueprint, time int) int {
	lastMax := 0
	cache := map[Cache]struct{}{}
	return maximizeGeodes(b, time, initFactory(), &lastMax, cache)
}

func part1(blueprints []Blueprint) int {
	var wg sync.WaitGroup
	wg.Add(len(blueprints))

	var total int32 = 0
	for _, blp := range blueprints {
		b := blp
		go func() {
			produced := doOne(b, 24)
			atomic.AddInt32(&total, int32(produced*b.id))
			wg.Done()
		}()
	}
	wg.Wait()
	return int(total)
}

func part2(blueprints []Blueprint) int {
	total := 1
	for i := 0; i < 3; i++ {
		total *= doOne(blueprints[i], 32)
	}
	return total
}

func parseBlueprints(data string) []Blueprint {
	blueprintRegex := regexp.MustCompile(".+ (\\d+):.+ore.+ (\\d+).+clay.+ (\\d+) .+obsidian.+ (\\d+) .+ (\\d+).+geode.+ (\\d+) .+ (\\d+) .+")
	matches := blueprintRegex.FindAllStringSubmatch(data, -1)

	result := make([]Blueprint, len(matches))
	for i, m := range matches {
		result[i] = Blueprint{
			id: parse(m[1]),
			robots: []Robot{
				{
					cost:       Material{parse(m[2]), 0, 0, 0},
					production: Material{1, 0, 0, 0},
				},
				{
					cost:       Material{parse(m[3]), 0, 0, 0},
					production: Material{0, 1, 0, 0},
				},
				{
					cost:       Material{parse(m[4]), parse(m[5]), 0, 0},
					production: Material{0, 0, 1, 0},
				},
				{
					cost:       Material{parse(m[6]), 0, parse(m[7]), 0},
					production: Material{0, 0, 0, 1},
				},
			},
		}
	}
	return result
}

func main() {
	rawContent, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(rawContent)
	blueprints := parseBlueprints(content)

	fmt.Println("Part 1: ", part1(blueprints))
	fmt.Println("Part 2: ", part2(blueprints))
}
