package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/kilianlievens/advent-of-code-2023/advent"
)

type direction struct {
	XMutation    int
	YMutation    int
	EndOfTheLine func(*advent.Coord, [][]rune) bool
}

type hashDetection struct {
	LastSeen  int
	TimesSeen int
}

func rollRocks(matrix [][]rune, rocks []*advent.Coord, dir direction) bool {
	var changed bool

	for _, rock := range rocks {
		if dir.EndOfTheLine(rock, matrix) {
			continue
		}

		rolledY := rock.Y + dir.YMutation
		rolledX := rock.X + dir.XMutation

		if matrix[rolledY][rolledX] != '.' {
			continue
		}

		changed = true
		matrix[rolledY][rolledX] = 'O'
		matrix[rock.Y][rock.X] = '.'
		rock.Y = rock.Y + dir.YMutation
		rock.X = rock.X + dir.XMutation
	}

	return changed
}

func getLoad(rocks []*advent.Coord, maxY int) int {
	var load int
	for _, rock := range rocks {
		load += maxY - rock.Y
	}

	return load
}

func hashMatrix(matrix [][]rune) string {
	sha := sha256.New()
	sha.Write([]byte(fmt.Sprintf("%v", matrix)))

	return hex.EncodeToString(sha.Sum(nil))
}

func main() {
	exampleOneInput := advent.Read("./input/day_fourteen/example_one.txt")
	load, cycledLoad := doTheRolyPoly(exampleOneInput)
	fmt.Printf("Example input: load %d cycled load %d\n", load, cycledLoad) // 136, 64

	puzzleInput := advent.Read("./input/day_fourteen/puzzle_one.txt")
	load, cycledLoad = doTheRolyPoly(puzzleInput)
	fmt.Printf("Puzzle input: load %d cycled load %d\n", load, cycledLoad) // 113424, 96003
}

func doTheRolyPoly(input []string) (int, int) {
	var matrix [][]rune
	var roundedRocks []*advent.Coord

	// Parse
	for y, line := range input {
		var row []rune
		for x, r := range line {
			row = append(row, r)

			if r == 'O' {
				roundedRocks = append(roundedRocks, &advent.Coord{Y: y, X: x})
			}
		}
		matrix = append(matrix, row)
	}

	// Roll directions
	north := direction{
		XMutation:    0,
		YMutation:    -1,
		EndOfTheLine: func(rock *advent.Coord, _ [][]rune) bool { return rock.Y == 0 },
	}
	south := direction{
		XMutation:    0,
		YMutation:    +1,
		EndOfTheLine: func(rock *advent.Coord, matrix [][]rune) bool { return rock.Y == len(matrix)-1 },
	}
	west := direction{
		XMutation:    -1,
		YMutation:    0,
		EndOfTheLine: func(rock *advent.Coord, _ [][]rune) bool { return rock.X == 0 },
	}
	east := direction{
		XMutation:    +1,
		YMutation:    0,
		EndOfTheLine: func(rock *advent.Coord, matrix [][]rune) bool { return rock.X == len(matrix[0])-1 },
	}

	// Part one
	partOneMatrix := matrix
	partOneRocks := roundedRocks
	notDone := true
	for notDone == true {
		notDone = rollRocks(partOneMatrix, partOneRocks, north)
	}

	singleTiltLoad := getLoad(partOneRocks, len(partOneMatrix))

	// Part two
	partTwoMatrix := matrix
	partTwoRocks := roundedRocks
	dirOrder := [4]direction{north, west, south, east}
	seenHashes := map[string]*hashDetection{}
	iterations := 1000000000

	for i := 0; i < iterations; i++ {
		for _, dir := range dirOrder {
			dirNotDone := true
			for dirNotDone == true {
				dirNotDone = rollRocks(partTwoMatrix, partTwoRocks, dir)
			}
		}

		hash := hashMatrix(partTwoMatrix)
		hd, ok := seenHashes[hash]
		if !ok {
			seenHashes[hash] = &hashDetection{LastSeen: i, TimesSeen: 1}

			continue
		}

		if hd.TimesSeen > 3 {
			loopLength := i - hd.LastSeen
			if (iterations-1)%loopLength == i%loopLength {
				break
			}
		}

		hd.LastSeen = i
		hd.TimesSeen++
	}

	cycleTiltLoad := getLoad(partTwoRocks, len(partTwoMatrix))

	return singleTiltLoad, cycleTiltLoad
}
