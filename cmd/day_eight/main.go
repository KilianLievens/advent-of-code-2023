package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/kilianlievens/advent-of-code-2023/advent"
)

type node struct {
	Value    string
	EndNode  bool
	LeftKey  string
	RightKey string
}

type loop struct {
	Offset uint64
	Length uint64
}

func main() {
	var steps, ghostSteps uint64

	exampleOneInput := advent.Read("./input/day_eight/example_one.txt")
	steps, _ = navigate(exampleOneInput)
	fmt.Printf("Example one input: steps %d\n", steps) // 2

	exampleTwoInput := advent.Read("./input/day_eight/example_two.txt")
	steps, _ = navigate(exampleTwoInput)
	fmt.Printf("Example two input: steps %d\n", steps) // 6

	exampleThreeInput := advent.Read("./input/day_eight/example_three.txt")
	steps, ghostSteps = navigate(exampleThreeInput)
	fmt.Printf("Example three input: steps %d, ghostSteps %d\n", steps, ghostSteps) // 2, 6

	puzzleInput := advent.Read("./input/day_eight/puzzle_one.txt")
	steps, ghostSteps = navigate(puzzleInput)
	fmt.Printf("Puzzle input: steps %d, ghostSteps %d\n", steps, ghostSteps) // 18157, 14299763833181
}

func step(instructions []string, net map[string]node, position string, steps uint64) (string, bool) {
	instruction := instructions[steps%uint64(len(instructions))]

	curNode := net[position]
	if instruction == "L" {
		return curNode.LeftKey, curNode.EndNode
	}
	return curNode.RightKey, curNode.EndNode
}

func gcd(a, b uint64) uint64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}

	return a
}

func lcm(a, b uint64, integers ...uint64) uint64 {
	result := a * b / gcd(a, b)

	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}

	return result
}

func navigate(input []string) (uint64, uint64) {
	// Parse
	instructions := strings.Split(input[0], "")
	nodesRaw := input[1:]
	net := map[string]node{}
	var starts []string

	for _, r := range nodesRaw {
		value := r[0:3]
		net[value] = node{
			Value:    value,
			LeftKey:  r[7:10],
			RightKey: r[12:15],
			EndNode:  value[2:3] == "Z",
		}

		if value[2:3] == "A" {
			starts = append(starts, value)
		}
	}

	// Part one
	position := starts[0]
	steps := uint64(0)

	for i := uint64(0); i < math.MaxInt64; i++ {
		var done bool
		position, done = step(instructions, net, position, i)
		if done {
			steps = i
			break
		}
	}

	// Part two not applicable?
	if len(starts) < 2 {
		return steps, 0
	}

	// Part two
	positions := starts
	var loops []loop

	for _, p := range positions {
		position := p
		foundDone := false
		originalFoundDoneSteps := uint64(0)
		loopSteps := uint64(0)

		for i := uint64(0); i < math.MaxInt64; i++ {
			var done bool
			position, done = step(instructions, net, position, i)

			if done {
				if foundDone == true {
					loopSteps = i - originalFoundDoneSteps
					break
				}

				foundDone = true
				originalFoundDoneSteps = i
			}
		}

		loops = append(loops, loop{
			Offset: originalFoundDoneSteps,
			Length: loopSteps,
		})
	}

	// We don't need the offset. AoC was merciful?
	var loopLenghts []uint64
	for _, l := range loops {
		loopLenghts = append(loopLenghts, l.Length)
	}

	ghostSteps := lcm(loopLenghts[0], loopLenghts[1], loopLenghts[2:]...)

	return steps, ghostSteps
}
