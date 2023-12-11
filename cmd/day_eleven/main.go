package main

import (
	"fmt"

	"github.com/kilianlievens/advent-of-code-2023/advent"
)

const (
	galaxy = iota
	space
)

var symbolMap map[rune]int = map[rune]int{
	'.': space,
	'#': galaxy,
}

func getExpandingY(universe [][]int) []int {
	var knownYExpansion []int
	for y, row := range universe {
		spaceOnly := true
		for _, s := range row {
			if s == galaxy {
				spaceOnly = false

				break
			}
		}

		if spaceOnly {
			knownYExpansion = append(knownYExpansion, y)
		}
	}

	return knownYExpansion
}

func predictGalaxies(
	galaxies []advent.Coord,
	expandingY, expandingX []int,
	expandCoefficient int,
) []advent.Coord {
	var expandedGalaxies []advent.Coord
	for _, g := range galaxies {
		toAddY := 0
		for _, y := range expandingY {
			if y < g.Y {
				toAddY++
			}
		}

		toAddX := 0
		for _, x := range expandingX {
			if x < g.X {
				toAddX++
			}
		}

		expandedGalaxies = append(expandedGalaxies, advent.Coord{
			Y: g.Y + toAddY*(expandCoefficient-1),
			X: g.X + toAddX*(expandCoefficient-1),
		})
	}

	return expandedGalaxies
}

func calcManhattanDistanceSum(galaxies []advent.Coord) int {
	sum := 0
	for i, og := range galaxies {
		for _, ng := range galaxies[i+1:] {
			sum += advent.CalcManhattanDistance(og, ng)
		}
	}

	return sum
}

func main() {
	exampleOneInput := advent.Read("./input/day_eleven/example_one.txt")
	distances := getGalaxyDistances(exampleOneInput)
	fmt.Printf("Example input: distances %v\n", distances) // [374 1030 8410 82000210]

	puzzleInput := advent.Read("./input/day_eleven/puzzle_one.txt")
	distances = getGalaxyDistances(puzzleInput)
	fmt.Printf("Puzzle input: distances %v\n", distances) // [9681886 16010894 87212234 791134099634]
}

func getGalaxyDistances(input []string) []int {
	var originalUniverse [][]int
	var originalGalaxies []advent.Coord

	// Parse
	for y, line := range input {
		var universeRow []int

		for x, r := range line {
			s := symbolMap[r]
			universeRow = append(universeRow, s)

			if s == galaxy {
				originalGalaxies = append(originalGalaxies, advent.Coord{Y: y, X: x})
			}
		}

		originalUniverse = append(originalUniverse, universeRow)
	}

	// Expansion predictions
	expandingY := getExpandingY(originalUniverse)
	tUniverse := advent.Transpose2D(originalUniverse)
	expandingX := getExpandingY(tUniverse)

	// Get distances
	return []int{
		calcManhattanDistanceSum(predictGalaxies(originalGalaxies, expandingY, expandingX, 2)),
		calcManhattanDistanceSum(predictGalaxies(originalGalaxies, expandingY, expandingX, 10)),
		calcManhattanDistanceSum(predictGalaxies(originalGalaxies, expandingY, expandingX, 100)),
		calcManhattanDistanceSum(predictGalaxies(originalGalaxies, expandingY, expandingX, 1_000_000)),
	}
}
