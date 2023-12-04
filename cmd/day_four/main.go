package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/kilianlievens/advent-of-code-2023/advent"
)

func main() {
	exampleOneInput := advent.Read("./input/day_four/example_one.txt")
	points, scratchCards := win(exampleOneInput)
	fmt.Printf("Example one input: points %d, scratch cards %d\n", points, scratchCards) // 13, 30

	puzzleInput := advent.Read("./input/day_four/puzzle_one.txt")
	points, scratchCards = win(puzzleInput)
	fmt.Printf("Puzzle input: points %d, scratch cards %d\n", points, scratchCards) // 23235, 5920640
}

func win(input []string) (int, int) {
	points := 0
	scratchCards := 0
	executions := map[int]int{}

	for y, line := range input {
		matches := 0

		allNum := strings.Split(line, ": ")[1]
		numSegments := strings.Split(allNum, " | ")
		rawWiningNum := numSegments[0]
		rawScratchNum := numSegments[1]

		winningNum := map[int]bool{}
		interWinningNum := strings.Split(rawWiningNum, " ")
		for _, num := range interWinningNum {
			numInt, err := strconv.Atoi(num)
			if err != nil {
				continue // empty string
			}

			winningNum[numInt] = true
		}

		interScratchNum := strings.Split(rawScratchNum, " ")
		for _, num := range interScratchNum {
			numInt, err := strconv.Atoi(num)
			if err != nil {
				continue // empty string
			}

			if _, ok := winningNum[numInt]; ok {
				matches++
			}
		}

		points += int(math.Pow(2, float64(matches-1)))

		executionTimes := executions[y] + 1 // +1 for initial execution
		for i := 0; i < executionTimes; i++ {
			scratchCards++

			for i := 1; i <= matches; i++ {
				executions[y+i]++
			}
		}
	}

	return points, scratchCards
}
