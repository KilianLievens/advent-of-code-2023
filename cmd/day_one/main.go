package main

import (
	"strconv"
	"strings"

	"github.com/kilianlievens/advent-of-code-2023/advent"
)

func injectNumbers(input string) string {
	input = strings.ReplaceAll(input, "one", "one1one")
	input = strings.ReplaceAll(input, "two", "two2two")
	input = strings.ReplaceAll(input, "three", "three3three")
	input = strings.ReplaceAll(input, "four", "four4four")
	input = strings.ReplaceAll(input, "five", "five5five")
	input = strings.ReplaceAll(input, "six", "six6six")
	input = strings.ReplaceAll(input, "seven", "seven7seven")
	input = strings.ReplaceAll(input, "eight", "eight8eight")
	input = strings.ReplaceAll(input, "nine", "nine9nine")
	return input
}

func main() {
	exampleOneInput := advent.Read("./input/day_one/example_one.txt")
	println(calibrate(exampleOneInput, false)) // 142
	exampleTwoInput := advent.Read("./input/day_one/example_two.txt")
	println(calibrate(exampleTwoInput, true)) // 281

	puzzleInput := advent.Read("./input/day_one/puzzle_one.txt")
	println(calibrate(puzzleInput, false)) // 53334
	println(calibrate(puzzleInput, true))  // 52834
}

func calibrate(input []string, parseText bool) int {
	var calibrationSum int

	for _, line := range input {
		if parseText {
			line = injectNumbers(line)
		}

		foundFirst := false
		var firstDigit rune
		var lastDigit rune

		for _, char := range line {
			if _, err := strconv.Atoi(string(char)); err != nil {
				continue
			}

			if !foundFirst {
				firstDigit = char
				foundFirst = true
			}

			lastDigit = char
		}

		calibration, _ := strconv.Atoi(string([]rune{firstDigit, lastDigit}))
		calibrationSum += calibration
	}

	return calibrationSum
}
