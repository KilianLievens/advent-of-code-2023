package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kilianlievens/advent-of-code-2023/advent"
)

type lens struct {
	Label       string
	FocalLength int
}

func holidayAsciiStringHelperAlgorithm(s string) int {
	currentValue := 0
	for _, r := range s {
		currentValue += int(r)
		currentValue *= 17
		currentValue %= 256
	}

	return currentValue
}

func addLens(list []lens, newLens lens) []lens {
	for i, l := range list {
		if l.Label == newLens.Label {
			list[i] = newLens
			return list
		}
	}

	return append(list, newLens)
}

func removeLens(list []lens, newLens lens) []lens {
	for i, l := range list {
		if l.Label == newLens.Label {
			if i == len(list)-1 {
				return list[:i]
			}

			return append(list[:i], list[i+1:]...)
		}
	}

	return list
}

func getBoxFocusingPower(boxNumber int, lenses []lens) int {
	focusingPower := 0
	for i, l := range lenses {
		focusingPower += (i + 1) * l.FocalLength * (1 + boxNumber)
	}

	return focusingPower
}

func main() {
	exampleOneInput := advent.Read("./input/day_fifteen/example_one.txt")
	hashSum, focalPower := helpTheCuteReindeer(exampleOneInput)
	fmt.Printf("Example input: hashSum %d focal power %d\n", hashSum, focalPower) // 1320, 145

	puzzleInput := advent.Read("./input/day_fifteen/puzzle_one.txt")
	hashSum, focalPower = helpTheCuteReindeer(puzzleInput)
	fmt.Printf("Puzzle input: hashSum %d focal power %d\n", hashSum, focalPower) // 511343, 294474
}

func helpTheCuteReindeer(input []string) (int, int) {
	segments := strings.Split(input[0], ",")

	// Part one
	sum := 0
	for _, seg := range segments {
		sum += holidayAsciiStringHelperAlgorithm(seg)
	}

	// Part two
	power := 0
	holidayAsciiStringHelperManualArrangementProcedureBoxes := map[int][]lens{}

	for _, seg := range segments {
		l := lens{}
		var operation rune
		var foundOperation bool

		for _, r := range seg {
			if foundOperation {
				l.FocalLength, _ = strconv.Atoi(string(r))
				continue
			}

			if r == '-' || r == '=' {
				foundOperation = true
				operation = r
				continue
			}

			l.Label += string(r)
		}

		boxNum := holidayAsciiStringHelperAlgorithm(l.Label)

		currentLenses, _ := holidayAsciiStringHelperManualArrangementProcedureBoxes[boxNum]

		if operation == '-' {
			holidayAsciiStringHelperManualArrangementProcedureBoxes[boxNum] = removeLens(currentLenses, l)
			continue
		}

		holidayAsciiStringHelperManualArrangementProcedureBoxes[boxNum] = addLens(currentLenses, l)
	}

	for boxNum, lenses := range holidayAsciiStringHelperManualArrangementProcedureBoxes {
		power += getBoxFocusingPower(boxNum, lenses)
	}

	return sum, power
}
