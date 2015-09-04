package main

import (
	"encoding/csv"
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
	"time"
)

// These have to be imported to avoid unknown image format errors.
import _ "image/jpeg"
import _ "image/gif"

type ColorDefinition struct {
	name        string
	hex         string
	minDistance float64
	maxDistance float64
}

func retrieveColorDefinitions(filename string) []ColorDefinition {
	// This is gross and needs to be automated.
	Black := ColorDefinition{
		name:        "Black",
		hex:         "#191818",
		minDistance: 0.023842433019540137,
		maxDistance: 0.02104211515520766,
	}
	Brown := ColorDefinition{
		name:        "Brown",
		hex:         "#795000",
		minDistance: 0.19809101648616032,
		maxDistance: 0.19210336287428736,
	}
	Blue := ColorDefinition{
		name:        "Blue",
		hex:         "#3F4AFF",
		minDistance: 0.4641466052266115,
		maxDistance: 0.4633987252122273,
	}
	Gold := ColorDefinition{
		name:        "Gold",
		hex:         "#C1B000",
		minDistance: 0.13133286731508037,
		maxDistance: 0.1264689296432074,
	}
	Gray := ColorDefinition{
		name:        "Gray",
		hex:         "#7D7C7A",
		minDistance: 0.15423603423913348,
		maxDistance: 0.15004028649601187,
	}
	Green := ColorDefinition{
		name:        "Green",
		hex:         "#1CBD2A",
		minDistance: 0.7492784109976135,
		maxDistance: 0.7456682300440184,
	}
	Orange := ColorDefinition{
		name:        "Orange",
		hex:         "#C27B13",
		minDistance: 0.3696758700641386,
		maxDistance: 0.3666353779425447,
	}
	Pink := ColorDefinition{
		name:        "Pink",
		hex:         "#FFBECC",
		minDistance: 0.14258279796428422,
		maxDistance: 0.13907449495014637,
	}
	Purple := ColorDefinition{
		name:        "Purple",
		hex:         "#9E4DFF",
		minDistance: 0.28147884511085836,
		maxDistance: 0.2798940121040469,
	}
	Red := ColorDefinition{
		name:        "Red",
		hex:         "#FF260C",
		minDistance: 0.6663975631872985,
		maxDistance: 0.6639093437480922,
	}
	Tan := ColorDefinition{
		name:        "Tan",
		hex:         "#D18D12",
		minDistance: 0.0981207988444481,
		maxDistance: 0.0948660417438293,
	}
	White := ColorDefinition{
		name:        "White",
		hex:         "#FFFDF7",
		minDistance: 0.03822537621306147,
		maxDistance: 0.03477705237396113,
	}
	Yellow := ColorDefinition{
		name:        "Yellow",
		hex:         "#FFF000",
		minDistance: 0.20183243554226693,
		maxDistance: 0.1971941953430347,
	}
	return []ColorDefinition{Black, Brown, Blue, Gold, Gray, Green, Orange, Pink, Purple, Red, Tan, White, Yellow}
}

func downloadImageFromUrl(url string, saveAs string) {
	response, err := http.Get(url)
	closeIfError(err)
	defer response.Body.Close()

	file, err := os.Create(saveAs)
	closeIfError(err)

	// dump the response body to the file
	_, err = io.Copy(file, response.Body)
	closeIfError(err)

	file.Close()
}

func deleteFileByLocation(location string) {
	err := os.Remove(location)
	closeIfError(err)
}

func openImage(filename string) (result image.Image, skip bool) {
	imgfile, err := os.Open(filename)
	shouldSkip := closeIfError(err)
	if shouldSkip {
		return nil, true
	}

	defer imgfile.Close()

	img, _, err := image.Decode(imgfile)
	shouldSkip = closeIfError(err)
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
		closeIfError(err)
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
		closeIfError(err)

		err = png.Encode(out, resizedImage)
		closeIfError(err)
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
		closeIfError(err)

		err = png.Encode(out, debugImg)
		closeIfError(err)
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

	// createDebugImage("debug.png", bounds, clusterPoints)

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

func analyzeCluster(cluster []color.Color, definedColors []ColorDefinition) (generatedColor string, matches []string) {
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
	results := []string{}

	for _, color := range definedColors {
		comparison, _ := colorful.Hex(color.hex)
		if color.minDistance <= finalColor.DistanceLab(comparison) && finalColor.DistanceLab(comparison) <= color.maxDistance {
			results = append(results, color.name)
		} else {
			if (color.maxDistance - finalColor.DistanceLab(comparison)) < 0.003 {
				fmt.Printf("\nmin: %v\ndist: %v\nmax: %v\n", color.minDistance, finalColor.DistanceLab(comparison), color.maxDistance)
			}
		}
	}

	return finalColor.Hex(), results
}

func euclidianDistance(pOne int, pTwo int, qOne int, qTwo int) float64 {
	// https://en.wikipedia.org/wiki/Euclidean_distance#Two_dimensions
	return math.Sqrt(math.Pow(float64(qOne-pOne), 2) + math.Pow(float64(qTwo-pTwo), 2))
}

func closeIfError(err error) bool {
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err)
		return true
	}
	return false
}

func main() {
	start := time.Now()
	rand.Seed(time.Now().Unix())
	saveAs := "sample.png"
	testingCSV := true
	testingListOfImages := false
	k := 5
	numberOfImages := 0

	colorDefinitions := retrieveColorDefinitions("")

	if testingListOfImages {
		imageLocations := []string{
			"http://i.imgur.com/WpsnGdF.jpg",
			"http://i.imgur.com/0fOqb1G.jpg",
			"http://i.imgur.com/fKVX14Z.jpg",
			"http://i.imgur.com/9MaHnRC.jpg",
			"http://i.imgur.com/weEh6sY.jpg",
			"http://i.imgur.com/jU8hj6v.png",
		}
		for _, location := range imageLocations {
			downloadImageFromUrl(location, saveAs)
			img, _ := openImage(saveAs)

			croppedImg := cropImage(img, 50, "cropped.png")
			resizedImg := resizeImage(croppedImg, 200, "resized.png")

			clusters := createClusters(k, resizedImg)

			fmt.Printf("\n%v: \n", location)
			analyzeCluster(clusters[0], colorDefinitions)
			analyzeCluster(clusters[1], colorDefinitions)
			analyzeCluster(clusters[2], colorDefinitions)
			analyzeCluster(clusters[3], colorDefinitions)
			analyzeCluster(clusters[4], colorDefinitions)

			deleteFileByLocation(saveAs)
		}
	} else if testingCSV {
		csvfile, err := os.Create("slt_output.csv")
		closeIfError(err)

		sourceFile, err := os.Open("skusandimages.csv")
		closeIfError(err)

		writer := csv.NewWriter(csvfile)
		reader := csv.NewReader(sourceFile)
		lines, err := reader.ReadAll()
		closeIfError(err)
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
		closeIfError(err)

		for lineNumber, line := range lines {
			if lineNumber > 10 {
				sku := line[0]
				imageUrl := line[1]

				downloadImageFromUrl(imageUrl, saveAs)
				img, shouldSkip := openImage(saveAs)

				if !shouldSkip {
					croppedImg := cropImage(img, 50, "cropped.png")
					resizedImg := resizeImage(croppedImg, 200, "resized.png")
					clusters := createClusters(k, resizedImg)

					color0, matches0 := analyzeCluster(clusters[0], colorDefinitions)
					color1, matches1 := analyzeCluster(clusters[1], colorDefinitions)
					color2, matches2 := analyzeCluster(clusters[2], colorDefinitions)
					color3, matches3 := analyzeCluster(clusters[3], colorDefinitions)
					color4, matches4 := analyzeCluster(clusters[4], colorDefinitions)

					err := writer.Write([]string{
						sku,
						imageUrl,
						color0,
						strings.Join(matches0, ","),
						color1,
						strings.Join(matches1, ","),
						color2,
						strings.Join(matches2, ","),
						color3,
						strings.Join(matches3, ","),
						color4,
						strings.Join(matches4, ","),
					})
					closeIfError(err)

					deleteFileByLocation(saveAs)
					numberOfImages += 1
				} else {
					err := writer.Write([]string{
						sku,
						imageUrl,
						"image",
						"skipped",
						"because",
						"of",
						"an",
						"unknown",
						"format",
						"error",
						"sorry!",
						":(",
					})
					closeIfError(err)
				}
			}
		}
		writer.Flush()
	} else {
		url := "http://i.imgur.com/WpsnGdF.jpg"
		downloadImageFromUrl(url, saveAs)
		img, _ := openImage(saveAs)

		croppedImg := cropImage(img, 50, "cropped.png")
		resizedImg := resizeImage(croppedImg, 200, "resized.png")

		clusters := createClusters(k, resizedImg)

		fmt.Printf("\n%v: \n", url)
		analyzeCluster(clusters[0], colorDefinitions)
		analyzeCluster(clusters[1], colorDefinitions)
		analyzeCluster(clusters[2], colorDefinitions)
		analyzeCluster(clusters[3], colorDefinitions)
		analyzeCluster(clusters[4], colorDefinitions)

		deleteFileByLocation(saveAs)
	}
	elapsed := time.Since(start)
	log.Printf("Processing %v images took %s", numberOfImages, elapsed)
}
