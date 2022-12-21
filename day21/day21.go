package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Context struct {
	nodes   map[string]Node
	visited map[string]struct{}
}

func (ctx *Context) visit(name string) {
	ctx.visited[name] = struct{}{}
}

func (ctx *Context) clear() {
	ctx.visited = map[string]struct{}{}
}

func (ctx *Context) didVisit(n string) bool {
	_, ok := ctx.visited[n]
	return ok
}

type Node interface {
	evaluate(ctx *Context) float64
}

type Number struct {
	value float64
}

func (n *Number) evaluate(ctx *Context) float64 {
	return n.value
}

type Variable struct {
	name string
}

func (v *Variable) evaluate(ctx *Context) float64 {
	ctx.visit(v.name)
	node, ok := ctx.nodes[v.name]
	if !ok {
		fmt.Errorf("cannot find node %v\n", v.name)
	}
	return node.evaluate(ctx)
}

type Operation struct {
	a  Node
	b  Node
	op string
}

func (op *Operation) evaluate(ctx *Context) float64 {
	if op.op == "+" {
		return op.a.evaluate(ctx) + op.b.evaluate(ctx)
	} else if op.op == "-" {
		return op.a.evaluate(ctx) - op.b.evaluate(ctx)
	} else if op.op == "*" {
		return op.a.evaluate(ctx) * op.b.evaluate(ctx)
	} else if op.op == "/" {
		return op.a.evaluate(ctx) / op.b.evaluate(ctx)
	} else {
		fmt.Errorf("operation %v is not supported", op.op)
		return 0
	}
}

func parseScenario(data string) Context {
	nodes := map[string]Node{}
	regex := regexp.MustCompile("(.+): ((\\d+)|(.+) (.) (.+))")
	allMatches := regex.FindAllStringSubmatch(data, -1)

	for _, match := range allMatches {
		if match[3] == "" {
			nodes[match[1]] = &Operation{
				a:  &Variable{match[4]},
				b:  &Variable{match[6]},
				op: match[5],
			}
		} else {
			val, _ := strconv.Atoi(match[3])
			nodes[match[1]] = &Number{float64(val)}
		}
	}
	return Context{nodes, map[string]struct{}{}}
}

func part1(ctx *Context) int64 {
	return int64(ctx.nodes["root"].evaluate(ctx))
}

type Step struct {
	leftSelected  bool
	oppositeValue float64
	operation     string
}

func reverse(steps []Step, target float64) float64 {
	s := steps[0]
	if s.operation == "" {
		return target
	} else if s.operation == "=" {
		target = s.oppositeValue
	} else if s.operation == "+" {
		target = target - s.oppositeValue
	} else if s.operation == "-" {
		if s.leftSelected {
			target = target + s.oppositeValue
		} else {
			target = s.oppositeValue - target
		}
	} else if s.operation == "*" {
		target = target / s.oppositeValue
	} else if s.operation == "/" {
		if s.leftSelected {
			target = target * s.oppositeValue
		} else {
			target = target / s.oppositeValue
		}
	} else {
		fmt.Errorf("unknown operation %v", s.operation)
	}
	return reverse(steps[1:], target)
}

func part2(ctx *Context) int64 {
	var steps []Step
	nextToVisit := ctx.nodes["root"].(*Operation)
	for nextToVisit != nil {
		ctx.clear()

		left := nextToVisit.a.evaluate(ctx)
		leftVisited := ctx.didVisit("humn")
		right := nextToVisit.b.evaluate(ctx)

		if leftVisited {
			steps = append(steps, Step{true, right, nextToVisit.op})
			nextVar := nextToVisit.a.(*Variable)
			nextToVisit, _ = ctx.nodes[nextVar.name].(*Operation)
		} else {
			steps = append(steps, Step{false, left, nextToVisit.op})
			nextVar := nextToVisit.b.(*Variable)
			nextToVisit, _ = ctx.nodes[nextVar.name].(*Operation)
		}
	}

	steps[0].operation = "="
	steps = append(steps, Step{false, 0, ""})

	return int64(reverse(steps, 0))
}

func main() {
	rawContent, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(rawContent)

	context := parseScenario(content)

	fmt.Println("Part 1: ", part1(&context))
	fmt.Println("Part 2: ", part2(&context))
}
