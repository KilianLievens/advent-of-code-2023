package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/kilianlievens/advent-of-code-2023/advent"
)

type history struct {
	Sequence   []int
	DiffLayers [][]int
}

func (h *history) PopulateDiff() {
	allZeroes := true
	var curDiffLayer []int
	lastDiffLayer := h.DiffLayers[len(h.DiffLayers)-1]

	for i := 0; i < len(lastDiffLayer)-1; i++ {
		diff := lastDiffLayer[i+1] - lastDiffLayer[i]
		curDiffLayer = append(curDiffLayer, diff)

		if diff != 0 {
			allZeroes = false
		}
	}

	h.DiffLayers = append(h.DiffLayers, curDiffLayer)

	if !allZeroes {
		h.PopulateDiff()
	}
}

func (h *history) Extrapolate() {
	// -2 because we don't need to do the zero layer
	for i := len(h.DiffLayers) - 2; i >= 0; i-- {
		curDifflayer := h.DiffLayers[i]
		lastDifflayer := h.DiffLayers[i+1]

		next := curDifflayer[len(curDifflayer)-1] + lastDifflayer[len(lastDifflayer)-1]
		h.DiffLayers[i] = append(curDifflayer, next)
	}
}

func main() {
	exampleOneInput := advent.Read("./input/day_nine/example_one.txt")
	prediction, rPrediction := predict(exampleOneInput)
	fmt.Printf("Example one input: prediction %d, reverse prediction %d\n", prediction, rPrediction) // 114, 2

	puzzleInput := advent.Read("./input/day_nine/puzzle_one.txt")
	prediction, rPrediction = predict(puzzleInput)
	fmt.Printf("Puzzle input: prediction %d, reverse prediction %d\n", prediction, rPrediction) // 1637452937, 908
}

func predict(input []string) (int, int) {
	var sequences []*history

	for _, line := range input {
		numStrings := strings.Split(line, " ")

		var sequence []int
		for _, ns := range numStrings {
			n, _ := strconv.Atoi(ns)
			sequence = append(sequence, n)
		}

		sequences = append(sequences, &history{Sequence: sequence, DiffLayers: [][]int{sequence}})
	}

	var sum int
	var reversedSum int
	for _, s := range sequences {
		s.PopulateDiff()
		s.Extrapolate()

		sum += s.DiffLayers[0][len(s.DiffLayers[0])-1]

		reversed := s.Sequence
		slices.Reverse(reversed)

		s.DiffLayers = [][]int{reversed}
		s.PopulateDiff()
		s.Extrapolate()

		reversedSum += s.DiffLayers[0][len(s.DiffLayers[0])-1]
	}

	return sum, reversedSum
}
