package main

import (
	"log"
	"math"
	"strings"
)

func closeIfError(statement string, err error) {
	if err != nil {
		log.Fatal(statement, ":\n", err)
	}
}

func skipIfError(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}

func euclidianDistance(pOne int, pTwo int, qOne int, qTwo int) float64 {
	// https://en.wikipedia.org/wiki/Euclidean_distance#Two_dimensions
	return math.Sqrt(math.Pow(float64(qOne-pOne), 2) + math.Pow(float64(qTwo-pTwo), 2))
}

func buildRow(results map[string][]string) []string {
	returnSlice := []string{""}
	for hex, names := range results {
		returnSlice = append(returnSlice, hex)
		returnSlice = append(returnSlice, strings.Join(names, ","))
	}
	return returnSlice
}
