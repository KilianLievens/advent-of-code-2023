package main

import (
	"strconv"
	"strings"

	"github.com/kilianlievens/advent-of-code-2023/advent"
)

var textNumberReplacer = strings.NewReplacer(
	"one", "1",
	"two", "2",
	"three", "3",
	"four", "4",
	"five", "5",
	"six", "6",
	"seven", "7",
	"eight", "8",
	"nine", "9",
)

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
			line = textNumberReplacer.Replace(line)
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
