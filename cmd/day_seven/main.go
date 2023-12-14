package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/kilianlievens/advent-of-code-2023/advent"
)

const (
	highCard = iota
	onePair
	twoPair
	threeOfAKind
	fullHouse
	fourOfAKind
	fiveOfAKind
)

var faceMap = map[rune]int{
	'A': 14,
	'K': 13,
	'Q': 12,
	'J': 11,
	'T': 10,
	'9': 9,
	'8': 8,
	'7': 7,
	'6': 6,
	'5': 5,
	'4': 4,
	'3': 3,
	'2': 2,
}

type hand struct {
	Score     int
	Cards     string
	Bid       int
	JokerHand bool
}

type games []hand

func (h games) Len() int           { return len(h) }
func (h games) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h games) Less(i, j int) bool { return lesserHand(h[i], h[j]) }

func main() {
	exampleOneInput := advent.Read("./input/day_seven/example_one.txt")
	winnings, jokerWinnings := playPoker(exampleOneInput)
	fmt.Printf("Example input: winnings %d, jokerWinnings %d\n", winnings, jokerWinnings) // 6440, 5905

	puzzleInput := advent.Read("./input/day_seven/puzzle_one.txt")
	winnings, jokerWinnings = playPoker(puzzleInput)
	fmt.Printf("Puzzle input: winnings %d, jokerWinnings %d\n", winnings, jokerWinnings) // 251058093, 249781879
}

func score(input string) int {
	cards := map[rune]int{}
	for _, r := range input {
		_, ok := cards[r]
		if !ok {
			cards[r] = 1

			continue
		}

		cards[r]++
	}

	pairs := 0
	triples := 0

	for _, value := range cards {
		if value == 5 {
			return fiveOfAKind
		}

		if value == 4 {
			return fourOfAKind
		}

		if value == 3 {
			triples++
		}

		if value == 2 {
			pairs++
		}
	}

	if pairs == 1 && triples == 1 {
		return fullHouse
	}

	if triples == 1 {
		return threeOfAKind
	}

	if pairs == 2 {
		return twoPair
	}

	if pairs == 1 {
		return onePair
	}

	return highCard
}

func adjustScoreForJokers(input string, score int) int {
	if score == fiveOfAKind {
		return fiveOfAKind
	}

	jokers := 0
	for _, i := range input {
		if i == 'J' {
			jokers++
		}
	}

	if score == fourOfAKind {
		if jokers > 0 {
			return fiveOfAKind
		}

		return fourOfAKind
	}

	if score == fullHouse {
		if jokers > 0 {
			return fiveOfAKind
		}

		return fullHouse
	}

	if score == threeOfAKind {
		if jokers > 0 {
			// Can't be 3-2 because that would be a full house
			return fourOfAKind
		}

		return threeOfAKind
	}

	if score == twoPair {
		if jokers == 1 {
			return fullHouse
		}

		if jokers == 2 {
			return fourOfAKind
		}

		return twoPair
	}

	if score == onePair {
		if jokers > 0 {
			return threeOfAKind
		}

		return onePair
	}

	// highCard
	if jokers > 0 {
		return onePair
	}

	return score
}

func lesserHand(handA, handB hand) bool {
	if handA.Score == handB.Score {
		for i := 0; i < len(handA.Cards); i++ {
			a := rune(handA.Cards[i])
			b := rune(handB.Cards[i])

			if a == b {
				continue
			}

			// Are we playing jokerStyle?
			if handA.JokerHand && handB.JokerHand {
				if a == 'J' {
					return true
				}

				if b == 'J' {
					return false
				}
			}

			aV, _ := faceMap[a]
			bV, _ := faceMap[b]

			return aV < bV
		}

		panic("Should not be possible")
	}

	if handA.Score < handB.Score {
		return true
	}

	return false
}

func playPoker(input []string) (int, int) {
	winnings := 0
	jokerWinnings := 0

	var hands games
	var jokerHands games

	for _, line := range input {
		segments := strings.Split(line, " ")
		rawHand := segments[0]
		rawBid := segments[1]
		bid, _ := strconv.Atoi(rawBid)

		hands = append(
			hands, hand{
				Score:     score(rawHand),
				Bid:       bid,
				Cards:     rawHand,
				JokerHand: false,
			},
		)

		jokerHands = append(
			jokerHands, hand{
				Score:     adjustScoreForJokers(rawHand, score(rawHand)),
				Bid:       bid,
				Cards:     rawHand,
				JokerHand: true,
			},
		)
	}

	sort.Sort(hands)
	sort.Sort(jokerHands)

	for i, hand := range hands {
		winnings += hand.Bid * (i + 1)
	}

	for i, jHand := range jokerHands {
		jokerWinnings += jHand.Bid * (i + 1)
	}

	return winnings, jokerWinnings
}
