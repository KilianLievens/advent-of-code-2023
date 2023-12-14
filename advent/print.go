package advent

import "fmt"

func PrintRuneMatrix(matrix [][]rune) {
	for _, line := range matrix {
		for _, r := range line {
			fmt.Printf("%s", string(r))
		}
		fmt.Printf("\n")
	}
}

func PrintStringMatrix(matrix [][]string) {
	for _, line := range matrix {
		for _, r := range line {
			fmt.Printf("%s", r)
		}
		fmt.Printf("\n")
	}
}
