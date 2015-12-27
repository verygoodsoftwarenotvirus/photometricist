package main

import (
	"encoding/csv"
	"github.com/lucasb-eyer/go-colorful"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Configuration struct {
	K                 int               `json:"k"`
	InputFilename     string            `json:"input_filename"`
	OutputFilename    string            `json:"output_filename"`
	DownloadDirectory string            `json:"download_dir"`
	ColorDefinitions  []ColorDefinition `json:"colors"`
	TimeString        string
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

func analyzeImages(line []string, config Configuration, currentImageNumber int, writer *csv.Writer, wg sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	// Row indices.
	imageIndex := 1
	skuIndex := 0

	sku := line[skuIndex]
	imageUrl := line[imageIndex]

	// generate filenames and download.
	downloadFilename, resizedFilename, croppedFilename := buildFilenames(config, currentImageNumber)
	downloadImageFromUrl(imageUrl, downloadFilename)
	img, shouldSkip := openImage(downloadFilename)

	if !shouldSkip {
		croppedImg := cropImage(img, 50, croppedFilename)
		resizedImg := resizeImage(croppedImg, 200, resizedFilename)

		clusters := createClusters(config.K, resizedImg)

		results := make(map[string][]string, config.K)
		for x := 0; x < config.K; x++ {
			analyzeCluster(clusters[x], config.ColorDefinitions, results)
		}

		newRow := []string{sku, imageUrl}
		newRow = append(newRow, buildRow(results)...)

		err := writer.Write(newRow)
		closeIfError("Error occurred writing csv row", err)

		deleteFileByLocation(downloadFilename)
		deleteFileByLocation(resizedFilename)
		deleteFileByLocation(croppedFilename)

		wg.Done()
	} else {
		wg.Done()
	}
}

func main() {
	runtime.GOMAXPROCS(2)
	start := time.Now()
	rand.Seed(time.Now().Unix())

	// filenames
	configFilename := "config.json"

	// arbitrary variables
	numberOfImages := 0

	// config things
	configuration := retrieveConfiguration(configFilename)
	configuration.TimeString = strconv.FormatInt(start.UnixNano(), 10)
	inputFilename := configuration.InputFilename
	outputFilename := configuration.OutputFilename

	if configuration.DownloadDirectory == "" {
		var err error
		configuration.DownloadDirectory, err = filepath.Abs(filepath.Dir(os.Args[0]))
		closeIfError("Error getting current directory", err)
	} else {
		ensureFolderExistence(configuration.DownloadDirectory)
	}
	os.Chdir(configuration.DownloadDirectory)

	// CSV stuff
	lines := readInputFile(inputFilename)
	writer := setupOutputFileWriter(outputFilename)

	// concurreny stuff
	var wg sync.WaitGroup

	for lineNumber, line := range lines {
		if lineNumber == 0 {
			/* skip headers */
		} else if lineNumber < 2 {
			go analyzeImages(line, configuration, numberOfImages, writer, wg)

			numberOfImages += 1
		}
	}

	wg.Wait()
	writer.Flush()

	elapsed := time.Since(start)
	log.Printf("Processing %v colors from %v images took %s", numberOfImages*configuration.K, numberOfImages, elapsed)
}
