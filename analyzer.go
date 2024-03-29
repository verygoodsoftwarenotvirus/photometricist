package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"github.com/disintegration/imaging"
	"github.com/lucasb-eyer/go-colorful"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// These have to be imported to avoid unknown image format errors.
import _ "image/jpeg"
import _ "image/gif"

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

// helper functions
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

// file handling functions

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

func downloadImageFromUrl(url string, saveAs string) {
	file, err := os.Create(saveAs)
	closeIfError("Error creating file to save downloaded image to", err)

	response, err := http.Get(url)
	closeIfError("Error downloading image", err)
	defer response.Body.Close()

	// dump the response body to the file
	_, err = io.Copy(file, response.Body)
	closeIfError("Error copying download response to opened file.", err)

	file.Close()
}

func openImage(filename string) (result image.Image, skip bool) {
	imgfile, err := os.Open(filename)
	shouldSkip := skipIfError(err)
	if shouldSkip {
		return nil, true
	}

	defer imgfile.Close()

	img, _, err := image.Decode(imgfile)
	shouldSkip = skipIfError(err)
	if shouldSkip {
		return nil, true
	}

	return img, false
}

// image handling functions

func cropImage(img image.Image, percentage float64, filename string) image.Image {
	bounds := img.Bounds()
	width := int(float64(bounds.Max.X) * (percentage * .01))
	height := int(float64(bounds.Max.Y) * (percentage * .01))

	if len(filename) > 0 {
		out, err := os.Create(filename)
		err = png.Encode(out, imaging.CropCenter(img, width, height))
		closeIfError("Error encoding image as PNG", err)
	}

	return imaging.CropCenter(img, width, height)
}

func resizeImage(img image.Image, maxWidth int, filename string) image.Image {
	bounds := img.Bounds()
	modifier := float64(maxWidth) / float64(bounds.Max.X)
	height := int(float64(bounds.Max.Y) * modifier)

	resizedImage := imaging.Resize(img, maxWidth, height, imaging.Lanczos)
	if len(filename) > 0 {
		out, err := os.Create(filename)
		closeIfError("Error creating file for resized image", err)

		err = png.Encode(out, resizedImage)
		closeIfError("Error encoding resized image", err)
	}
	return resizedImage
}

func createDebugImage(filename string, bounds image.Rectangle, clusterPoints []map[string]int) {
	debugImgOutline := image.Rect(0, 0, bounds.Max.X, bounds.Max.Y)
	debugImg := image.NewRGBA(debugImgOutline)
	draw.Draw(debugImg, debugImg.Bounds(), &image.Uniform{color.Transparent}, image.ZP, draw.Src)

	for _, point := range clusterPoints {
		// magic numbers ahoy!
		debugPointOutline := image.Rect(point["X"], point["Y"], point["X"]+5, point["Y"]+5)
		debugPoint := image.Rect(point["X"]+1, point["Y"]+1, point["X"]+4, point["Y"]+4)
		basePoint := image.Point{X: point["X"], Y: point["Y"]}
		draw.Draw(debugImg, debugPointOutline, &image.Uniform{color.White}, basePoint, draw.Src)
		draw.Draw(debugImg, debugPoint, &image.Uniform{color.Black}, basePoint, draw.Src)
	}
	if len(filename) > 0 {
		out, err := os.Create(filename)
		closeIfError("Error creating file for debug image", err)

		err = png.Encode(out, debugImg)
		closeIfError("Error encoding debug image", err)
	}
}

// analysis functions

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

func analyzeImages(line []string, config Configuration, currentImageNumber int, writer *csv.Writer) {
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
	}
}

func analyzeColors() {
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

	for lineNumber, line := range lines {
		if lineNumber == 0 {
			/* skip headers */
		} else if lineNumber < 2 {
			analyzeImages(line, configuration, numberOfImages, writer)

			numberOfImages += 1
		}
	}

	writer.Flush()

	elapsed := time.Since(start)
	log.Printf("Processing %v colors from %v images took %s", numberOfImages*configuration.K, numberOfImages, elapsed)
}
