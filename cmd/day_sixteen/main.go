package main

import (
	"fmt"

	"github.com/kilianlievens/advent-of-code-2023/advent"
)

const (
	emptySpace = iota
	forwardMirror
	backwardMirror
	verticalSplitter
	horizontalSplitter
)

const (
	N = iota
	E
	S
	W
)

var runeMap = map[rune]int{
	'.':  emptySpace,
	'/':  forwardMirror,
	'\\': backwardMirror,
	'|':  verticalSplitter,
	'-':  horizontalSplitter,
}

type tile struct {
	Energized      bool
	Type           int
	DirectionsDone map[int]bool
}

type direction struct {
	Name           int
	XMutation      int
	YMutation      int
	EndOfTheLine   func(advent.Coord, [][]*tile) bool
	ForwardMirror  *direction
	BackwardMirror *direction
}

type movingBeam struct {
	Dir   direction
	Coord advent.Coord
}

func pathBeam(b movingBeam, matrix [][]*tile) []movingBeam {
	if b.Dir.EndOfTheLine(b.Coord, matrix) {
		return []movingBeam{}
	}

	nextBeam := movingBeam{Coord: advent.Coord{Y: b.Coord.Y + b.Dir.YMutation, X: b.Coord.X + b.Dir.XMutation}}
	nextTile := matrix[nextBeam.Coord.Y][nextBeam.Coord.X]

	_, doneDirection := nextTile.DirectionsDone[b.Dir.Name]
	if doneDirection {
		return []movingBeam{}
	}
	nextTile.Energized = true
	nextTile.DirectionsDone[b.Dir.Name] = true

	var vertical bool
	switch b.Dir.Name {
	case N, S:
		vertical = true
	case E, W:
		vertical = false
	}

	if nextTile.Type == emptySpace ||
		(nextTile.Type == verticalSplitter && vertical) ||
		(nextTile.Type == horizontalSplitter && !vertical) {
		nextBeam.Dir = b.Dir
		return []movingBeam{nextBeam}
	}

	if nextTile.Type == forwardMirror {
		nextBeam.Dir = *(b.Dir.ForwardMirror)
		return []movingBeam{nextBeam}
	}

	if nextTile.Type == backwardMirror {
		nextBeam.Dir = *(b.Dir.BackwardMirror)
		return []movingBeam{nextBeam}
	}

	if nextTile.Type == verticalSplitter || nextTile.Type == horizontalSplitter {
		nextBeamA := nextBeam
		nextBeamA.Dir = *(b.Dir.ForwardMirror)

		nextBeamB := nextBeam
		nextBeamB.Dir = *(b.Dir.BackwardMirror)

		return []movingBeam{nextBeamA, nextBeamB}
	}

	panic("tile type shouldn't be possible")
}

func calcBeamEffect(start movingBeam, matrix [][]*tile) int {
	beams := []movingBeam{start}

	for len(beams) > 0 {
		beams = append(beams[1:], pathBeam(beams[0], matrix)...)
	}

	var energized int
	for _, line := range matrix {
		for _, t := range line {
			if t.Energized {
				energized++
			}
		}
	}

	return energized
}

func createMatrix(input []string) [][]*tile {
	var matrix [][]*tile

	// Parse
	for _, line := range input {
		var row []*tile
		for _, r := range line {
			row = append(row, &tile{Type: runeMap[r], DirectionsDone: make(map[int]bool)})
		}
		matrix = append(matrix, row)
	}

	return matrix
}

func main() {
	exampleOneInput := advent.Read("./input/day_sixteen/example_one.txt")
	energized, maxEnergized := energize(exampleOneInput)
	fmt.Printf("Example input: energized %d maxEnergized %d\n", energized, maxEnergized) // 46, 51

	puzzleInput := advent.Read("./input/day_sixteen/puzzle_one.txt")
	energized, maxEnergized = energize(puzzleInput)
	fmt.Printf("Puzzle input: energized %d maxEnergized %d\n", energized, maxEnergized) // 8323, 8491
}

func energize(input []string) (int, int) {
	matrix := createMatrix(input)

	// Initialize directions
	north := direction{
		Name:         N,
		XMutation:    0,
		YMutation:    -1,
		EndOfTheLine: func(beam advent.Coord, _ [][]*tile) bool { return beam.Y == 0 },
	}
	south := direction{
		Name:         S,
		XMutation:    0,
		YMutation:    +1,
		EndOfTheLine: func(beam advent.Coord, matrix [][]*tile) bool { return beam.Y == len(matrix)-1 },
	}
	west := direction{
		Name:         W,
		XMutation:    -1,
		YMutation:    0,
		EndOfTheLine: func(beam advent.Coord, _ [][]*tile) bool { return beam.X == 0 },
	}
	east := direction{
		Name:         E,
		XMutation:    +1,
		YMutation:    0,
		EndOfTheLine: func(beam advent.Coord, matrix [][]*tile) bool { return beam.X == len(matrix[0])-1 },
	}
	north.ForwardMirror = &east
	north.BackwardMirror = &west
	south.ForwardMirror = &west
	south.BackwardMirror = &east
	east.ForwardMirror = &north
	east.BackwardMirror = &south
	west.ForwardMirror = &south
	west.BackwardMirror = &north

	// Part one
	partOneMatrix := matrix
	// Starts off the matrix!
	energized := calcBeamEffect(movingBeam{Dir: east, Coord: advent.Coord{Y: 0, X: -1}}, partOneMatrix)

	// Part two
	maxEnergized := 0

	for i := 0; i < len(matrix); i++ {
		maxEnergized = advent.MaxInt(
			maxEnergized, calcBeamEffect(movingBeam{Dir: east, Coord: advent.Coord{Y: i, X: -1}}, createMatrix(input)),
		)

		maxEnergized = advent.MaxInt(
			maxEnergized,
			calcBeamEffect(movingBeam{Dir: west, Coord: advent.Coord{Y: i, X: len(matrix[0])}}, createMatrix(input)),
		)
	}

	for i := 0; i < len(matrix[0]); i++ {
		maxEnergized = advent.MaxInt(
			maxEnergized, calcBeamEffect(movingBeam{Dir: south, Coord: advent.Coord{Y: -1, X: i}}, createMatrix(input)),
		)

		maxEnergized = advent.MaxInt(
			maxEnergized,
			calcBeamEffect(movingBeam{Dir: north, Coord: advent.Coord{Y: len(matrix), X: i}}, createMatrix(input)),
		)
	}

	return energized, maxEnergized
}
