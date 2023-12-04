package main

import (
	"fmt"
	"strconv"

	a "github.com/kilianlievens/advent-of-code-2023/advent"
)

type Gear struct {
	neighbours int
	value      int
}

func main() {
	exampleOneInput := a.Read("./input/day_three/example_one.txt")
	partSum, gearSum := fixEngine(exampleOneInput)
	fmt.Printf("Example input: partSum %d, gearSum %d\n", partSum, gearSum) // 4361, 467835

	puzzleInput := a.Read("./input/day_three/puzzle_one.txt")
	partSum, gearSum = fixEngine(puzzleInput)
	fmt.Printf("Puzzle input: partSum %d, gearSum %d\n", partSum, gearSum) // 533784, 78826761
}

func fixEngine(input []string) (int, int) {
	partSum := 0
	gearSum := 0

	var gears map[a.Coord]*Gear = make(map[a.Coord]*Gear)
	var symbols map[a.Coord]bool = make(map[a.Coord]bool)

	// Populate symbol matrix and find all potential gears
	for y, line := range input {
		for x, char := range line {
			_, err := strconv.Atoi(string(char))
			isSymbol := err != nil && char != '.'
			symbols[a.Coord{X: x, Y: y}] = isSymbol

			if char == '*' {
				newGear := Gear{neighbours: 0, value: 1}
				gears[a.Coord{X: x, Y: y}] = &newGear
			}
		}
	}

	// Parse numbers and detect/calculate surrounding symbols and gears
	for y, line := range input {
		previousWasNumber := false
		currentNumberString := ""

		line = line + "." // Extra column (to detect numbers at the end of the line)

		for x, char := range line {
			_, err := strconv.Atoi(string(char))
			// Current char is a number
			if err == nil {
				currentNumberString += string(char)
				previousWasNumber = true

				continue
			}

			// Last char was a number and this char is not a number
			if previousWasNumber {
				surroundingSymbol := false
				currentNumber, _ := strconv.Atoi(currentNumberString)

				for i := max(y-1, 0); i <= min(y+1, len(input)-1); i++ {
					for j := max(x-len(currentNumberString)-1, 0); j <= min(x, len(line)-1); j++ {
						isSymbol, ok := symbols[a.Coord{X: j, Y: i}]
						if ok && isSymbol {
							surroundingSymbol = true
						}

						gear, ok := gears[a.Coord{X: j, Y: i}]
						if ok {
							gear.neighbours++
							gear.value *= currentNumber
						}
					}
				}

				if surroundingSymbol {
					partSum += currentNumber
				}
			}

			currentNumberString = ""
			previousWasNumber = false
		}
	}

	for _, gear := range gears {
		if gear.neighbours == 2 {
			gearSum += gear.value
		}
	}

	return partSum, gearSum
}
