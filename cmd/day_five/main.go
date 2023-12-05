package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/kilianlievens/advent-of-code-2023/advent"
)

type conversion struct {
	SourceStart int
	DestStart   int
	RangeLength int
}

type seedPair struct {
	Start       int
	RangeLength int
}

func main() {
	exampleOneInput := advent.Read("./input/day_five/example_one.txt")
	location := findLocation(exampleOneInput, false)
	locationWithRange := findLocation(exampleOneInput, true)
	fmt.Printf("Example input: location %d, with range location %d\n", location, locationWithRange) // 35, 46

	puzzleInput := advent.Read("./input/day_five/puzzle_one.txt")
	location = findLocation(puzzleInput, false)
	locationWithRange = findLocation(puzzleInput, true)
	fmt.Printf("Puzzle input: location %d, with range location %d\n", location, locationWithRange) // 318728750, 37384986
}

func convert(convs []conversion, input int) int {
	for _, conv := range convs {
		if input >= conv.SourceStart && input < conv.SourceStart+conv.RangeLength {
			diff := input - conv.SourceStart
			return conv.DestStart + diff
		}
	}

	return input
}

func findLocation(input []string, seedRange bool) int {
	lowestLocation := math.MaxInt32

	currentType := "seeds"

	var seeds []seedPair
	var seedToSoil []conversion
	var soilToFertilizer []conversion
	var fertilizerToWater []conversion
	var waterToLight []conversion
	var ligthToTemperature []conversion
	var temperatureToHumidity []conversion
	var humidityToLocation []conversion

	// Parse
	for _, line := range input {
		if strings.Contains(line, "seed-to-soil") {
			currentType = "seed-to-soil"
			continue
		}

		if strings.Contains(line, "soil-to-fertilizer") {
			currentType = "soil-to-fertilizer"
			continue
		}

		if strings.Contains(line, "fertilizer-to-water") {
			currentType = "fertilizer-to-water"
			continue
		}

		if strings.Contains(line, "water-to-light") {
			currentType = "water-to-light"
			continue
		}

		if strings.Contains(line, "light-to-temperature") {
			currentType = "light-to-temperature"
			continue
		}

		if strings.Contains(line, "temperature-to-humidity") {
			currentType = "temperature-to-humidity"
			continue
		}

		if strings.Contains(line, "humidity-to-location") {
			currentType = "humidity-to-location"
			continue
		}

		if currentType == "seeds" {
			seedsRaw := strings.Split(line, ": ")
			seedStrings := strings.Split(seedsRaw[1], " ")

			if !seedRange {
				for _, seedString := range seedStrings {
					seed, _ := strconv.Atoi(seedString)
					seeds = append(seeds, seedPair{
						Start:       seed,
						RangeLength: 1,
					})
				}

				continue
			}

			for i := 0; i < len(seedStrings); i += 2 {
				startSeed, _ := strconv.Atoi(seedStrings[i])
				seedRangeLength, _ := strconv.Atoi(seedStrings[i+1])
				seeds = append(seeds, seedPair{
					Start:       startSeed,
					RangeLength: seedRangeLength,
				})
			}

			continue
		}

		rawConversion := strings.Split(line, " ")
		destStart, _ := strconv.Atoi(rawConversion[0])
		sourceStart, _ := strconv.Atoi(rawConversion[1])
		rangeLength, _ := strconv.Atoi(rawConversion[2])
		conversion := conversion{
			SourceStart: sourceStart,
			DestStart:   destStart,
			RangeLength: rangeLength,
		}

		switch currentType {
		case "seed-to-soil":
			seedToSoil = append(seedToSoil, conversion)
		case "soil-to-fertilizer":
			soilToFertilizer = append(soilToFertilizer, conversion)
		case "fertilizer-to-water":
			fertilizerToWater = append(fertilizerToWater, conversion)
		case "water-to-light":
			waterToLight = append(waterToLight, conversion)
		case "light-to-temperature":
			ligthToTemperature = append(ligthToTemperature, conversion)
		case "temperature-to-humidity":
			temperatureToHumidity = append(temperatureToHumidity, conversion)
		case "humidity-to-location":
			humidityToLocation = append(humidityToLocation, conversion)
		}
	}

	seedToLocation := func(pair seedPair, c chan<- int) {
		low := math.MaxInt32

		for i := pair.Start; i < pair.Start+pair.RangeLength; i++ {
			soil := convert(seedToSoil, i)
			fert := convert(soilToFertilizer, soil)
			wate := convert(fertilizerToWater, fert)
			ligh := convert(waterToLight, wate)
			temp := convert(ligthToTemperature, ligh)
			humi := convert(temperatureToHumidity, temp)
			loca := convert(humidityToLocation, humi)

			if loca < low {
				low = loca
			}
		}

		c <- low
	}

	// Execute
	c := make(chan int)

	for _, seed := range seeds {
		go seedToLocation(seed, c)
	}

	for range seeds {
		low := <-c
		if low < lowestLocation {
			lowestLocation = low
		}
	}

	return lowestLocation
}
