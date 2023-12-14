package main

import (
	"errors"
	"fmt"

	"github.com/kilianlievens/advent-of-code-2023/advent"
)

const (
	North = iota
	South
	East
	West
)

const (
	NS = iota
	EW
	NE
	NW
	SE
	SW

	Ground
	Start
)

var pipeMap = map[rune]int{
	'|': NS,
	'-': EW,
	'L': NE,
	'J': NW,
	'F': SE,
	'7': SW,

	'.': Ground,
	'S': Start,
}

func move(dir int, coord advent.Coord, maxX, maxY int) (advent.Coord, bool) {
	switch dir {
	case North:
		return advent.Coord{X: coord.X, Y: coord.Y - 1}, coord.Y > 0
	case South:
		return advent.Coord{X: coord.X, Y: coord.Y + 1}, coord.Y < maxY
	case West:
		return advent.Coord{X: coord.X - 1, Y: coord.Y}, coord.X > 0
	case East:
		return advent.Coord{X: coord.X + 1, Y: coord.Y}, coord.X < maxX
	default:
		return coord, false
	}
}

func getDirection(direction int, pipe int) (int, error) {
	if direction == North {
		switch pipe {
		case NS:
			return North, nil
		case SE:
			return East, nil
		case SW:
			return West, nil
		default:
			return 0, errors.New("invalid movement")
		}
	}

	if direction == South {
		switch pipe {
		case NS:
			return South, nil
		case NE:
			return East, nil
		case NW:
			return West, nil
		default:
			return 0, errors.New("invalid movement")
		}
	}

	if direction == East {
		switch pipe {
		case EW:
			return East, nil
		case NW:
			return North, nil
		case SW:
			return South, nil
		default:
			return 0, errors.New("invalid movement")
		}
	}

	if direction == West {
		switch pipe {
		case EW:
			return West, nil
		case NE:
			return North, nil
		case SE:
			return South, nil
		default:
			return 0, errors.New("invalid movement")
		}
	}

	return 0, errors.New("no valid direction found")
}

// Depends on N-E-S-W order
func parseStart(allowedDirs []int) int {
	if allowedDirs[0] == North {
		switch allowedDirs[1] {
		case East:
			return NE
		case West:
			return NW
		case South:
			return NS
		}
	}
	if allowedDirs[0] == East {
		switch allowedDirs[1] {
		case South:
			return SE
		case West:
			return EW
		}
	}

	return SW
}

func main() {
	exampleOneInput := advent.Read("./input/day_ten/example_one.txt")
	steps, enclosed := pipeItUp(exampleOneInput)
	fmt.Printf("Example one input: steps %d, enclosed %d\n", steps, enclosed) // 8, 1

	exampleTwoInput := advent.Read("./input/day_ten/example_two.txt")
	steps, enclosed = pipeItUp(exampleTwoInput)
	fmt.Printf("Example two input: steps %d, enclosed %d\n", steps, enclosed) // 80, 10

	puzzleInput := advent.Read("./input/day_ten/puzzle_one.txt")
	steps, enclosed = pipeItUp(puzzleInput)
	fmt.Printf("Puzzle input: steps %d, enclosed %d\n", steps, enclosed) // 6923, 529
}

func pipeItUp(input []string) (int, int) {
	var matrix [][]int
	var start advent.Coord

	// Parse matrix
	for y, line := range input {
		var row []int

		for x, r := range line {
			pipe := pipeMap[r]
			row = append(row, pipe)

			if pipe == Start {
				start = advent.Coord{X: x, Y: y}
			}
		}
		matrix = append(matrix, row)
	}

	maxX := len(matrix[0]) - 1
	maxY := len(matrix) - 1

	// Determine allowed start directions
	var allowedStartDirs []int
	for _, d := range []int{North, East, South, West} {
		position, ok := move(d, start, maxX, maxY)
		if !ok {
			continue
		}

		if _, err := getDirection(d, matrix[position.Y][position.X]); err == nil {
			allowedStartDirs = append(allowedStartDirs, d)
		}
	}

	// Path through the loop
	loop := map[advent.Coord]bool{}

	currentPosition, _ := move(allowedStartDirs[0], start, maxX, maxY)
	currentDirection := allowedStartDirs[0]

	steps := 1
	loop[currentPosition] = true

	for currentPosition.X != start.X || currentPosition.Y != start.Y {
		currentDirection, _ = getDirection(currentDirection, matrix[currentPosition.Y][currentPosition.X])
		currentPosition, _ = move(currentDirection, currentPosition, maxX, maxY)

		loop[currentPosition] = true

		steps++
	}

	// Determine enclosed tiles
	// Start must be a proper tile for the enclosed calculation
	matrix[start.Y][start.X] = parseStart(allowedStartDirs)

	enclosed := 0
	for y, row := range matrix {
		var inside = false

		// A vertical to horizontal tile and its corresponding horizontal to vertical tile are
		// considered a pair for this naming
		var firstHalfFound = false
		var firstHalf int

		for x := range row {
			// Is not part of the loop?
			if _, ok := loop[advent.Coord{X: x, Y: y}]; !ok {
				if inside {
					enclosed++
				}

				continue
			}

			// Get the loop segment
			tile := matrix[y][x]

			// Is horizontal tile? Horizontal tiles have no impact on inside calculation
			if tile == EW {
				continue
			}

			// Is Vertical tile? Vertical tiles determine the inside
			if tile == NS {
				inside = !inside
				continue
			}

			// Tile must then be vertical to horizontal. Have we not seen its pair yet?
			if !firstHalfFound {
				firstHalfFound = true
				firstHalf = tile

				continue
			}

			// If the pair together form a vertical
			if (firstHalf == NE && tile == SW) || (firstHalf == SE && tile == NW) {
				inside = !inside
			}

			// If the pair together form a U-turn
			firstHalfFound = false
		}

	}

	return steps / 2, enclosed
}
