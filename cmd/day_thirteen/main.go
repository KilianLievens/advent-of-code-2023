package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kilianlievens/advent-of-code-2023/advent"
)

func readWithWhitespace(fileName string) []string {
	body, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	return strings.Split(string(body), "\n")
}

func getMirrorY(matrix [][]string, requireSmudged bool) (int, error) {
	for y := 0; y < len(matrix)-1; y++ {
		isMirrorLine := true
		isSmudged := false

		maxDist := advent.MinInt(y, len(matrix)-y-2)
		for yDist := 0; yDist <= maxDist; yDist++ {
			for x := 0; x < len(matrix[0]); x++ {
				if matrix[y-yDist][x] != matrix[y+1+yDist][x] {
					if requireSmudged && !isSmudged {
						isSmudged = true

						continue
					}

					isMirrorLine = false

					break
				}
			}

			if !isMirrorLine {
				break
			}
		}

		if isMirrorLine {
			if requireSmudged && !isSmudged {
				continue
			}

			return y, nil
		}
	}

	return 0, errors.New("no mirror line found")
}

func getPatternSum(patterns [][][]string, requireSmudged bool) int {
	var sum int

	for i, p := range patterns {
		y, err := getMirrorY(p, requireSmudged)
		if err == nil {
			sum += 100 * (y + 1)

			continue
		}

		tP := advent.Transpose2D[string](p)
		x, err := getMirrorY(tP, requireSmudged)
		if err == nil {
			sum += x + 1

			continue
		}

		log.Fatalf("Could not find mirror line for pattern %d", i)
	}

	return sum
}

func main() {
	exampleOneInput := readWithWhitespace("./input/day_thirteen/example_one.txt")
	mirrors, smudgedMirrors := detectMirrors(exampleOneInput)
	fmt.Printf("Example input: mirrors %d smudged mirrors %d\n", mirrors, smudgedMirrors) // 405, 400

	puzzleInput := readWithWhitespace("./input/day_thirteen/puzzle_one.txt")
	mirrors, smudgedMirrors = detectMirrors(puzzleInput)
	fmt.Printf("Puzzle input: mirrors %d smudged mirrors %d\n", mirrors, smudgedMirrors) // 34821, 36919
}

func detectMirrors(input []string) (int, int) {
	var patterns [][][]string
	var currentPattern [][]string

	// Parse
	for _, line := range input {
		if line != "" {
			currentPattern = append(currentPattern, strings.Split(line, ""))
		}
		if line == "" {
			patterns = append(patterns, currentPattern)
			currentPattern = [][]string{}
		}
	}

	return getPatternSum(patterns, false), getPatternSum(patterns, true)
}
