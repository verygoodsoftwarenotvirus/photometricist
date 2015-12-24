package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
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
	"strings"
	// "sync"
	"time"
)

// These have to be imported to avoid unknown image format errors.
import _ "image/jpeg"
import _ "image/gif"

type ColorDefinitions struct {
	Definitions []ColorDefinition `json:"colors"`
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

func closeIfError(statement string, err error) {
	if err != nil {
		log.Fatal(statement, ":\n", err)
	}
}

func skipIfError(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

func euclidianDistance(pOne int, pTwo int, qOne int, qTwo int) float64 {
	// https://en.wikipedia.org/wiki/Euclidean_distance#Two_dimensions
	return math.Sqrt(math.Pow(float64(qOne-pOne), 2) + math.Pow(float64(qTwo-pTwo), 2))
}

func retrieveColorBoundaries(filename string) ColorDefinitions {
	var colors ColorDefinitions
	colorDefinitionFile, err := os.Open(filename)
	jsonParser := json.NewDecoder(colorDefinitionFile)
	err = jsonParser.Decode(&colors)
	closeIfError("Error decoding new_colors.json", err)
	return colors
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

func deleteFileByLocation(location string) {
	err := os.Remove(location)
	closeIfError("Error deleting file", err)
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
		debugPointOutline := image.Rect(point["X"], point["Y"], point["X"]+5, point["Y"]+5)
		debugPoint := image.Rect(point["X"]+1, point["Y"]+1, point["X"]+4, point["Y"]+4)
		draw.Draw(debugImg, debugPointOutline, &image.Uniform{color.White}, image.Point{X: point["X"], Y: point["Y"]}, draw.Src)
		draw.Draw(debugImg, debugPoint, &image.Uniform{color.Black}, image.Point{X: point["X"], Y: point["Y"]}, draw.Src)
	}
	if len(filename) > 0 {
		out, err := os.Create(filename)
		closeIfError("Error creating file for debug image", err)

		err = png.Encode(out, debugImg)
		closeIfError("Error encoding debug image", err)
	}
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
	finalColor := colorful.Color{(redTotal / pixelTotal) / 255.0, (greenTotal / pixelTotal) / 255.0, (blueTotal / pixelTotal) / 255.0}
	namedResults := []string{}

	for _, color := range definedColors {
		h, s, v := finalColor.Hsv()

		if color.Min.Hue <= h && h <= color.Max.Hue && color.Min.Value <= v && v <= color.Max.Value {
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

func buildRow(results map[string][]string) []string {
	returnSlice := []string{""}
	for hex, names := range results {
		returnSlice = append(returnSlice, hex)
		returnSlice = append(returnSlice, strings.Join(names, ","))
	}
	return returnSlice
}

func main() {
	start := time.Now()
	rand.Seed(time.Now().Unix())

	// filenames
	saveAs := "sample.png"
	croppedFilename := "cropped.png"
	resizedFilename := "resized.png"

	k := 5
	numberOfImages := 0
	numberOfColors := 0

	colorDefinitions := retrieveColorBoundaries("colordefs.json")

	csvfile, err := os.Create("slt_output.csv")
	closeIfError("Error creating output CSV file", err)

	sourceFile, err := os.Open("skusandimages.csv")
	closeIfError("Error opening input file", err)

	writer := csv.NewWriter(csvfile)
	reader := csv.NewReader(sourceFile)
	lines, err := reader.ReadAll()
	closeIfError("Error reading input CSV", err)
	err = writer.Write([]string{
		"SKU",
		"imageUrl",
		"Gen. Color 0",
		"Matches for 0",
		"Gen. Color 1",
		"Matches for 1",
		"Gen. Color 2",
		"Matches for 2",
		"Gen. Color 3",
		"Matches for 3",
		"Gen. Color 4",
		"Matches for 4",
	})
	closeIfError("Error writing headers", err)

	for lineNumber, line := range lines {
		if lineNumber == 0 {
			/* skip headers */
		} else if lineNumber < 15 {

			// Row indices.
			imageIndex := 1
			skuIndex := 0

			sku := line[skuIndex]
			imageUrl := line[imageIndex]

			downloadImageFromUrl(imageUrl, saveAs)
			img, shouldSkip := openImage(saveAs)

			if !shouldSkip {
				croppedImg := cropImage(img, 50, croppedFilename)
				resizedImg := resizeImage(croppedImg, 200, resizedFilename)
				clusters := createClusters(k, resizedImg)
				results := make(map[string][]string, k)

				analyzeCluster(clusters[0], colorDefinitions.Definitions, results)
				analyzeCluster(clusters[1], colorDefinitions.Definitions, results)
				analyzeCluster(clusters[2], colorDefinitions.Definitions, results)
				analyzeCluster(clusters[3], colorDefinitions.Definitions, results)
				analyzeCluster(clusters[4], colorDefinitions.Definitions, results)

				// fmt.Println(results)
				newRow := []string{sku, imageUrl}
				newRow = append(newRow, buildRow(results)...)
				err = writer.Write(newRow)
				// 	sku,
				// 	imageUrl,
				// 	color0,
				// 	strings.Join(results[0], ","),
				// 	color1,
				// 	strings.Join(results[1], ","),
				// 	color2,
				// 	strings.Join(results[2], ","),
				// 	color3,
				// 	strings.Join(results[3], ","),
				// 	color4,
				// 	strings.Join(results[4], ","),
				// })
				// closeIfError("Error writing results as CSV row", err)

				deleteFileByLocation(saveAs)
				numberOfImages += 1
				numberOfColors += k
			} // else {
			// err = writer.Write([]string{""})
			// closeIfError("Error writing blank line", err)
			//	}
		}
	}
	writer.Flush()

	deleteFileByLocation(croppedFilename)
	deleteFileByLocation(resizedFilename)

	elapsed := time.Since(start)
	log.Printf("Processing %v colors from %v images took %s", numberOfColors, numberOfImages, elapsed)
}
