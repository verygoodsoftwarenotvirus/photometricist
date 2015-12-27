package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"os"
	"strconv"
)

func retrieveConfiguration(filename string) Configuration {
	var config Configuration
	colorDefinitionFile, err := os.Open(filename)
	jsonParser := json.NewDecoder(colorDefinitionFile)
	err = jsonParser.Decode(&config)
	closeIfError("Error decoding new_colors.json", err)
	return config
}

func deleteFileByLocation(location string) {
	err := os.Remove(location)
	closeIfError("Error deleting file", err)
}

func ensureFolderExistence(folderName string) {
	_, err := os.Stat(folderName)
	if os.IsNotExist(err) {
		os.Mkdir(folderName, 0666)
	}
}

func buildFilenames(config Configuration, iteration int) (string, string, string) {
	var buffer bytes.Buffer
	buffer.WriteString(config.TimeString)
	buffer.WriteString("___")
	buffer.WriteString(strconv.Itoa(iteration))
	base := buffer.String()

	var download bytes.Buffer
	var resized bytes.Buffer
	var cropped bytes.Buffer

	download.WriteString(base)
	resized.WriteString(base)
	cropped.WriteString(base)

	download.WriteString(".png")
	resized.WriteString("_resized.png")
	cropped.WriteString("_cropped.png")
	return download.String(), resized.String(), cropped.String()
}

func readInputFile(inputFilename string) [][]string {
	sourceFile, err := os.Open(inputFilename)
	closeIfError("Error opening input file", err)

	reader := csv.NewReader(sourceFile)
	lines, err := reader.ReadAll()
	closeIfError("Error reading input CSV", err)
	return lines
}

func setupOutputFileWriter(outputFilename string) *csv.Writer {
	csvfile, err := os.Create(outputFilename)
	closeIfError("Error creating output CSV file", err)

	writer := csv.NewWriter(csvfile)
	err = writer.Write([]string{
		"SKU", "imageUrl", "Gen. Color 0", "Matches for 0", "Gen. Color 1",
		"Matches for 1", "Gen. Color 2", "Matches for 2", "Gen. Color 3",
		"Matches for 3", "Gen. Color 4", "Matches for 4",
	})
	closeIfError("Error writing headers", err)
	return writer
}
