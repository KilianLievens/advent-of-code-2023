package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kilianlievens/advent-of-code-2023/advent"
)

type state struct {
	Input string
	Nums  string
}

var cache map[state]int = make(map[state]int)

func sliceIsEqual(a, b []int, fullEqual bool) bool {
	if fullEqual && (len(a) != len(b)) {
		return false
	}

	if len(b) < len(a) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func getSpringLengths(input string) []int {
	var lengths []int

	springs := strings.Split(input, ".")
	for _, s := range springs {
		if s != "" {
			lengths = append(lengths, len(s))
		}
	}

	return lengths
}

func getIsValidResult(input string, nums []int) bool {
	springLengths := getSpringLengths(input)

	return sliceIsEqual(springLengths, nums, true)
}

func getPossiblyValid(input string, nums []int) bool {
	segments := strings.SplitN(input, "?", 2)

	if len(segments) == 1 {
		return getIsValidResult(input, nums)
	}

	springLengths := getSpringLengths(segments[0])

	if len(springLengths) == 0 {
		return true
	}

	if len(springLengths) > len(nums) {
		return false
	}

	if len(springLengths) > 1 && springLengths[0] < nums[0] {
		return false
	}

	missingNumbers := nums[len(springLengths):]
	sumMissingNumbers := -1 // account one less . in the middle
	for _, n := range missingNumbers {
		sumMissingNumbers += n + 1
	}

	if sumMissingNumbers > (len(segments[1]) + 1) {
		return false
	}

	last := springLengths[len(springLengths)-1]
	prunedSpringLengths := springLengths[:len(springLengths)-1]

	if last > nums[len(springLengths)-1] {
		return false
	}

	return sliceIsEqual(prunedSpringLengths, nums, false)
}

func popInput(input string, nums []int) (string, []int) {
	trimmed := strings.Trim(input, ".")

	segments := strings.SplitN(trimmed, "?", 2)
	springs := strings.Split(segments[0], ".")

	if len(springs) == 1 {
		return trimmed, nums
	}

	if len(springs[0]) == nums[0] {
		return trimmed[len(springs[0]):], nums[1:]
	}

	return trimmed, nums

}

func getAllReplacementPossibilites(input string, nums []int) int {
	s := state{
		Input: input,
		Nums:  fmt.Sprintf("%v\n", nums),
	}

	c, ok := cache[s]
	if ok {
		return c
	}

	curatedInput, curatedNums := popInput(input, nums)

	firstPossibility := strings.Replace(curatedInput, "?", "#", 1)
	secondPossibility := strings.Replace(curatedInput, "?", ".", 1)

	if firstPossibility == secondPossibility {
		if getIsValidResult(firstPossibility, curatedNums) {
			return 1
		}

		return 0
	}

	res := 0

	if getPossiblyValid(firstPossibility, curatedNums) {
		res += getAllReplacementPossibilites(firstPossibility, curatedNums)
	}

	if getPossiblyValid(secondPossibility, curatedNums) {
		res += getAllReplacementPossibilites(secondPossibility, curatedNums)
	}

	cache[s] = res

	return res
}

func getPossibilites(springs string, nums []int, foldCoefficient int) int {
	var unfoldedNums []int
	for i := 0; i < foldCoefficient; i++ {
		unfoldedNums = append(unfoldedNums, nums...)
	}

	var unfoldedSprings string = springs
	for i := 1; i < foldCoefficient; i++ {
		unfoldedSprings += "?" + springs
	}

	return getAllReplacementPossibilites(unfoldedSprings, unfoldedNums)
}

func main() {
	exampleOneInput := advent.Read("./input/day_twelve/example_one.txt")
	arrangements, folded := getArrangements(exampleOneInput)
	fmt.Printf("Example input: arrangements %d folded arrangements %d\n", arrangements, folded) // 21, 525152

	puzzleInput := advent.Read("./input/day_twelve/puzzle_one.txt")
	arrangements, folded = getArrangements(puzzleInput)
	fmt.Printf("Puzzle input: arrangements %d folded arrangements %d\n", arrangements, folded) // 7922, 18093821750095
}

func getArrangements(input []string) (int, int) {
	var sum int
	var foldedSum int

	for _, line := range input {
		segments := strings.Split(line, " ")

		springs := segments[0]
		rawNums := strings.Split(segments[1], ",")

		var nums []int
		for _, rn := range rawNums {
			num, _ := strconv.Atoi(rn)
			nums = append(nums, num)
		}

		// Part 1
		noFoldPossibilities := getPossibilites(springs, nums, 1)
		sum += noFoldPossibilities

		// Part 2
		fiveFoldPossibilities := getPossibilites(springs, nums, 5)
		foldedSum += fiveFoldPossibilities
	}

	return sum, foldedSum
}
