package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kilianlievens/advent-of-code-2023/advent"
)

func main() {
	exampleOneInput := advent.Read("./input/day_six/example_one.txt")
	res := race(exampleOneInput, false)
	superRes := race(exampleOneInput, true)
	fmt.Printf("Example input: standard %d, super %d\n", res, superRes) // 288, 71503

	puzzleInput := advent.Read("./input/day_six/puzzle_one.txt")
	res = race(puzzleInput, false)
	superRes = race(puzzleInput, true)
	fmt.Printf("Puzzle input: standard %d, super %d\n", res, superRes) // 800280, 45128024
}

func parseNumbers(s string, superRace bool) []int {
	var numbers []int
	var superNumberString string

	segments := strings.Split(s, " ")
	for _, s := range segments {
		number, err := strconv.Atoi(s)
		if err == nil {
			if superRace {
				superNumberString += s
			}

			if !superRace {
				numbers = append(numbers, number)
			}
		}
	}

	if superRace {
		superNumber, _ := strconv.Atoi(superNumberString)
		numbers = append(numbers, superNumber)
	}

	return numbers
}

func checkWin(waitTime int, totalTime int, winDistance int) bool {
	distance := waitTime * (totalTime - waitTime)

	return distance > winDistance
}

func getWinningCombos(time int, distance int) int {
	winners := 0

	for i := 1; i <= time; i++ {
		won := checkWin(i, time, distance)
		if won {
			winners++
		}
	}

	return winners
}

func race(input []string, superRace bool) int {
	multiplication := 1

	times := parseNumbers(input[0], superRace)
	distances := parseNumbers(input[1], superRace)

	for i := 0; i < len(times); i++ {
		multiplication *= getWinningCombos(times[i], distances[i])
	}

	return multiplication
}
