package main

import (
	"container/heap"
	"fmt"
	"strconv"

	"github.com/kilianlievens/advent-of-code-2023/advent"
)

type nodeKey struct {
	coord         advent.Coord
	lastDirection string
	straightCount int
}

type node struct {
	coord         advent.Coord
	cost          int
	estimatedCost int
	lastDirection string
	path          []advent.Coord
	straightCount int
}

type nodePrioQueue []node

func (q nodePrioQueue) Len() int {
	return len(q)
}

func (q nodePrioQueue) Less(i, j int) bool {
	return q[i].estimatedCost < q[j].estimatedCost
}

func (q nodePrioQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *nodePrioQueue) Push(x any) {
	*q = append(*q, x.(node))
}

func (q *nodePrioQueue) Pop() any {
	old := *q
	n := len(old)
	x := old[n-1]
	*q = old[:n-1]
	return x
}

type world map[advent.Coord]tile

type tile struct {
	coord advent.Coord
	cost  int
	world world
}

type direction struct {
	xMut int
	yMut int
	name string
}

type movement struct {
	coord         advent.Coord
	dir           direction
	cost          int
	straightCount int
}

// Spicy A star
func getToFactory(start, end advent.Coord, w world, maxStraightCount, minFirstMove int) int {
	var prioQ = &nodePrioQueue{}
	var doneNodes = make(map[nodeKey]bool)

	firstNode := node{
		path:  []advent.Coord{start},
		coord: start,
	}

	heap.Push(prioQ, firstNode)

	for prioQ.Len() > 0 {
		curNode := heap.Pop(prioQ).(node)
		curNodeKey := nodeKey{
			lastDirection: curNode.lastDirection,
			coord:         curNode.coord,
			straightCount: curNode.straightCount,
		}
		curTile, _ := w[curNode.coord]

		if curNode.coord == end {
			return curNode.cost
		}

		if _, done := doneNodes[curNodeKey]; done {
			continue
		}

		doneNodes[curNodeKey] = true

		moves := curTile.Path(curNode.lastDirection, curNode.straightCount, maxStraightCount, minFirstMove)
		for _, move := range moves {
			cost := curNode.cost + move.cost

			// Slices are frustrating >:(
			// https://stackoverflow.com/questions/27568213/could-anyone-explain-this-strange-behaviour-of-appending-to-golang-slices
			path := make([]advent.Coord, len(curNode.path)+1)
			copy(path, curNode.path)
			path[len(curNode.path)] = move.coord

			nextNode := node{
				lastDirection: move.dir.name,
				path:          path,
				coord:         move.coord,
				cost:          cost,
				estimatedCost: cost + advent.CalcManhattanDistance(curNode.coord, move.coord),
				straightCount: move.straightCount,
			}
			nextNodeKey := nodeKey{
				lastDirection: nextNode.lastDirection,
				coord:         nextNode.coord,
				straightCount: nextNode.straightCount,
			}

			if _, done := doneNodes[nextNodeKey]; !done {
				heap.Push(prioQ, nextNode)
			}
		}
	}

	panic("Could not find end.")
}

func (t tile) Path(lastMovement string, straightCount int, maxStraightCount int, minFirstMove int) []movement {
	var moves []movement
	directions := []direction{
		{xMut: 0, yMut: -1, name: "N"},
		{xMut: 0, yMut: +1, name: "S"},
		{xMut: +1, yMut: 0, name: "E"},
		{xMut: -1, yMut: 0, name: "W"},
	}

	// First move counts as one
	adjustedStraightCount := maxStraightCount - minFirstMove + 1

	// Same direction
	for _, dir := range directions {
		coord := advent.Coord{Y: t.coord.Y + dir.yMut, X: t.coord.X + dir.xMut}
		tiley, ok := t.world[coord]
		if ok &&
			straightCount < adjustedStraightCount && lastMovement == dir.name && lastMovement != "" {
			moves = append(moves, movement{coord: coord, dir: dir, cost: tiley.cost, straightCount: straightCount + 1})
		}
	}

	// Different direction
	for _, d := range directions {
		if d.name != lastMovement &&
			// Can't go opposite direction
			!(lastMovement == "N" && d.name == "S") && !(lastMovement == "S" && d.name == "N") &&
			!(lastMovement == "W" && d.name == "E") && !(lastMovement == "E" && d.name == "W") {

			intermediaryCost := 0
			for i := 1; i < minFirstMove; i++ {
				stepCoord := advent.Coord{Y: t.coord.Y + (d.yMut * i), X: t.coord.X + (d.xMut * i)}
				stepTile, _ := t.world[stepCoord]
				intermediaryCost += stepTile.cost
			}

			d.yMut *= minFirstMove
			d.xMut *= minFirstMove
			coord := advent.Coord{Y: t.coord.Y + d.yMut, X: t.coord.X + d.xMut}

			if tiley, ok := t.world[coord]; ok {
				moves = append(moves, movement{coord: coord, dir: d, cost: tiley.cost + intermediaryCost, straightCount: 1})
			}
		}
	}

	return moves
}

func main() {
	exampleOneInput := advent.Read("./input/day_seventeen/example_one.txt")
	bestPath, unknown := pathCrucible(exampleOneInput)
	fmt.Printf("Example one input: bestPath %d unknown %d\n", bestPath, unknown) // 102, 94

	exampleTwoInput := advent.Read("./input/day_seventeen/example_two.txt")
	bestPath, unknown = pathCrucible(exampleTwoInput)
	fmt.Printf("Example two input: bestPath %d unknown %d\n", bestPath, unknown) // 59, 71

	puzzleInput := advent.Read("./input/day_seventeen/puzzle_one.txt")
	bestPath, unknown = pathCrucible(puzzleInput)
	fmt.Printf("Puzzle input: bestPath %d unknown %d\n", bestPath, unknown) // 1128, 1268
}

func pathCrucible(input []string) (int, int) {
	var w = make(world)

	for y, line := range input {
		for x, r := range line {
			cost, _ := strconv.Atoi(string(r))

			coord := advent.Coord{Y: y, X: x}
			newTile := tile{
				coord: coord,
				cost:  cost,
				world: w,
			}

			w[coord] = newTile
		}
	}

	leastDistThreeStep := getToFactory(
		advent.Coord{Y: 0, X: 0},
		advent.Coord{Y: len(input) - 1, X: len(input[0]) - 1},
		w,
		3,
		1,
	)

	leastDistTenStep := getToFactory(
		advent.Coord{Y: 0, X: 0},
		advent.Coord{Y: len(input) - 1, X: len(input[0]) - 1},
		w,
		10,
		4,
	)

	return leastDistThreeStep, leastDistTenStep
}
