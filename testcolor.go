package main

import (
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
)

type Definition struct {
	minHex string
	hex    string
	maxHex string
	name   string
}

func loadColorDefinitions() []Definition {
	return []Definition{
		{
			name:   "Purple",
			minHex: "#7036b5",
			hex:    "#9E4DFF",
			maxHex: "#7137b6",
		},
		{
			name:   "Black",
			minHex: "#141313",
			hex:    "#191818",
			maxHex: "#141414",
		},
		{
			name:   "White",
			minHex: "#f4f2ed",
			hex:    "#FFFDF7",
			maxHex: "#f5f3ee",
		},
		{
			name:   "Tan",
			minHex: "#b77c0f",
			hex:    "#D18D12",
			maxHex: "#b87d10",
		},
		{
			name:   "Brown",
			minHex: "#4d3300",
			hex:    "#795000",
			maxHex: "#4e3400",
		},
		{
			name:   "Gray",
			minHex: "#575655",
			hex:    "#7D7C7A",
			maxHex: "#585756",
		},
		{
			name:   "Yellow",
			minHex: "#d1c400",
			hex:    "#FFF000",
			maxHex: "#d2c500",
		},
		{
			name:   "Orange",
			minHex: "#66410a",
			hex:    "#C27B13",
			maxHex: "#67420b",
		},
		{
			name:   "Gold",
			minHex: "#a49500",
			hex:    "#C1B000",
			maxHex: "#a59600",
		},
		{
			name:   "Green",
			minHex: "#07330b",
			hex:    "#1CBD2A",
			maxHex: "#08340c",
		},
		{
			name:   "Pink",
			minHex: "#d19ba7",
			hex:    "#FFBECC",
			maxHex: "#d29ca8",
		},
		{
			name:   "Blue",
			minHex: "#222789",
			hex:    "#3F4AFF",
			maxHex: "#23288a",
		},
		{
			name:   "Red",
			minHex: "#5e0e04",
			hex:    "#FF260C",
			maxHex: "#5f0f05",
		},
	}
}

func main() {
	colors := loadColorDefinitions()

	maxDifference := 0.0

	for _, colorDef := range colors {
		floor, _ := colorful.Hex(colorDef.minHex)
		color, _ := colorful.Hex(colorDef.hex)
		light, _ := colorful.Hex(colorDef.maxHex)
		fmt.Println("")
		fmt.Printf("%v: \n", colorDef.name)
		minDistance := color.DistanceLab(floor)
		maxDistance := color.DistanceLab(light)
		midpoint := (minDistance + maxDistance) / 2
		difference := minDistance - midpoint
		if difference > maxDifference {
			maxDifference = difference
		}

		// fmt.Println("Difference: ", difference)
		fmt.Println("minDistance: ", minDistance)
		fmt.Println("maxDistance: ", maxDistance)
		// fmt.Println("Midpoint: ", midpoint)
		// fmt.Println("Point:", maxDistance+(minDistance-midpoint))
	}
}
