package main

import (
	"github.com/lucasb-eyer/go-colorful"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"
	// "sync"
)

type Configuration struct {
	K                int               `json:"k"`
	InputFilename    string            `json:"input_filename"`
	OutputFilename   string            `json:"output_filename"`
	ColorDefinitions []ColorDefinition `json:"colors"`
}

type ColorDefinition struct {
	Name string   `json:"name"`
	Min  HSVColor `json:"min"`
	Max  HSVColor `json:"max"`
}

type HSVColor struct {
	Hue        float64 `json:"hue"`
	Saturation float64 `json:"saturation"`
	Value      float64 `json:"value"`
}

func createClusters(numberOfClusters int, img image.Image) map[int][]color.Color {
	clusters := make(map[int][]color.Color, numberOfClusters)
	clusterPoints := make([]map[string]int, numberOfClusters)
	bounds := img.Bounds()

	for i := 0; i < numberOfClusters; i++ {
		clusters[i] = []color.Color{}
		clusterPoints[i] = map[string]int{
			"X": rand.Intn(bounds.Max.X-bounds.Min.X) + bounds.Min.X,
			"Y": rand.Intn(bounds.Max.Y-bounds.Min.Y) + bounds.Min.Y,
		}
	}

	smallestDistanceIndex := math.MaxInt32
	smallestDistance := math.MaxFloat64

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {

			smallestDistanceIndex = math.MaxInt32
			smallestDistance = math.MaxFloat64

			for index, point := range clusterPoints {
				distance := euclidianDistance(x, y, point["X"], point["Y"])
				if distance < smallestDistance {
					smallestDistance = distance
					smallestDistanceIndex = index
				}
			}
			clusters[smallestDistanceIndex] = append(clusters[smallestDistanceIndex], img.At(x, y))
		}
	}
	return clusters
}

func analyzeCluster(cluster []color.Color, definedColors []ColorDefinition, results map[string][]string) {
	redTotal := float64(0.0)
	greenTotal := float64(0.0)
	blueTotal := float64(0.0)
	pixelTotal := float64(0.0)

	for pixel := range cluster {
		r, g, b, _ := cluster[pixel].RGBA()

		redTotal += float64(r >> 8)
		greenTotal += float64(g >> 8)
		blueTotal += float64(b >> 8)
		pixelTotal += 1
	}
	finalColor := colorful.Color{
		(redTotal / pixelTotal) / 255.0,
		(greenTotal / pixelTotal) / 255.0,
		(blueTotal / pixelTotal) / 255.0,
	}
	namedResults := []string{}

	for _, color := range definedColors {
		h, s, v := finalColor.Hsv()

		if color.Min.Hue <= h && h <= color.Max.Hue &&
			color.Min.Value <= v && v <= color.Max.Value {
			if color.Min.Saturation < color.Max.Saturation {
				if color.Min.Saturation <= s && s <= color.Max.Saturation {
					namedResults = append(namedResults, color.Name)
				}
			} else {
				if color.Min.Saturation >= s && s >= color.Max.Saturation {
					namedResults = append(namedResults, color.Name)
				}
			}
		}
	}

	results[finalColor.Hex()] = namedResults
}

func main() {
	start := time.Now()
	timeString := strconv.FormatInt(start.UnixNano(), 10)
	rand.Seed(time.Now().Unix())

	// filenames
	configFilename := "config.json"
	// downloadFilename := "sample.png"
	croppedFilename := "cropped.png"
	resizedFilename := "resized.png"
	downloadedImages := []string{}

	// arbitrary variables
	numberOfImages := 0
	numberOfColors := 0

	// config things
	configuration := retrieveConfiguration(configFilename)
	colorDefinitions := configuration.ColorDefinitions
	inputFilename := configuration.InputFilename
	outputFilename := configuration.OutputFilename
	k := configuration.K

	lines := readInputFile(inputFilename)
	writer := setupOutputFileWriter(outputFilename)

	for lineNumber, line := range lines {
		if lineNumber == 0 {
			/* skip headers */
		} else if lineNumber < 2 {

			// Row indices.
			imageIndex := 1
			skuIndex := 0

			sku := line[skuIndex]
			imageUrl := line[imageIndex]

			downloadFilename := buildFilename(timeString, numberOfImages)
			downloadedImages = append(downloadedImages, downloadFilename)
			downloadImageFromUrl(imageUrl, downloadFilename)
			img, shouldSkip := openImage(downloadFilename)

			if !shouldSkip {
				croppedImg := cropImage(img, 50, croppedFilename)
				resizedImg := resizeImage(croppedImg, 200, resizedFilename)
				clusters := createClusters(k, resizedImg)
				results := make(map[string][]string, k)

				analyzeCluster(clusters[0], colorDefinitions, results)
				analyzeCluster(clusters[1], colorDefinitions, results)
				analyzeCluster(clusters[2], colorDefinitions, results)
				analyzeCluster(clusters[3], colorDefinitions, results)
				analyzeCluster(clusters[4], colorDefinitions, results)

				newRow := []string{sku, imageUrl}
				newRow = append(newRow, buildRow(results)...)

				err := writer.Write(newRow)
				closeIfError("Error occurred writing csv row", err)

				numberOfImages += 1
				numberOfColors += k
			}
		}
	}
	writer.Flush()

	for _, file := range downloadedImages {
		deleteFileByLocation(file)
	}
	deleteFileByLocation(croppedFilename)
	deleteFileByLocation(resizedFilename)

	elapsed := time.Since(start)
	log.Printf("Processing %v colors from %v images took %s", numberOfColors, numberOfImages, elapsed)
}
