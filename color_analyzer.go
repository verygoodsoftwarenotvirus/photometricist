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

type ColorBoundary struct {
    name string
    minHue float64
    maxHue float64
    minSaturation float64
    maxSaturation float64
    minValue float64
    maxValue float64
}

type ColorBoundaries []ColorBoundary

func retrieveColorBoundaries() ColorBoundaries{
    return ColorBoundaries{
        {
            name: "Black",
            minHue: 0.0,
            maxHue: 0.0,
            minSaturation: 0.0,
            maxSaturation: 1.0,
            minValue: 0.0,
            maxValue: 0.16,
        },
        {
            name: "Brown",
            minHue: 9.0,
            maxHue: 45.0,
            minSaturation: 0.857142,
            maxSaturation: 1.0,
            minValue: 0.3325,
            maxValue: 0.68,
        },
        {
            name: "Blue",
            minHue: 161.0,
            maxHue: 255.0,
            minSaturation: 0.5185,
            maxSaturation: 0.5,
            minValue: 0.27,
            maxValue: 1.0,
        },
        {
            name: "Gold",
            minHue: 45.0,
            maxHue: 55.0,
            minSaturation: 0.8,
            maxSaturation: 0.69,
            minValue: 0.81,
            maxValue: 0.99,
        },
        {
            name: "Gray",
            minHue: 0.0,
            maxHue: 0.0,
            minSaturation: 0.0,
            maxSaturation: 0.3,
            minValue: 0.15,
            maxValue: 0.24,
        },
        {
            name: "Green",
            minHue: 64.0,
            maxHue: 141.0,
            minSaturation: 0.620689,
            maxSaturation: 0.5,
            minValue: 0.29,
            maxValue: 1.0,
        },
        {
            name: "Orange",
            minHue: 18.0,
            maxHue: 38.0,
            minSaturation: 0.823529,
            maxSaturation: 0.5,
            minValue: 0.34,
            maxValue: 1.0,
        },
        {
            name: "Pink",
            minHue: 289.0,
            maxHue: 347.0,
            minSaturation: 0.461538,
            maxSaturation: 0.5,
            minValue: 0.26,
            maxValue: 1.0,
        },
        {
            name: "Purple",
            minHue: 255.0,
            maxHue: 289.0,
            minSaturation: 0.39,
            maxSaturation: 0.5,
            minValue: 0.25,
            maxValue: 1.0,
        },
        {
            name: "Red",
            minHue: 0.0,
            maxHue: 18.0,
            minSaturation: 0.78,
            maxSaturation: 0.5,
            minValue: 0.33,
            maxValue: 1.0,
        },
        {
            name: "Red",
            minHue: 347.0,
            maxHue: 360.0,
            minSaturation: 0.78,
            maxSaturation: 0.5,
            minValue: 0.3,
            maxValue: 1.0,
        },
        {
            name: "Tan",
            minHue: 9.0,
            maxHue: 17.0,
            minSaturation: 0.620689,
            maxSaturation: 0.8,
            minValue: 0.493,
            maxValue: 1.0,
        },
        {
            name: "White",
            minHue: 0.0,
            maxHue: 0.0,
            minSaturation: 0.0,
            maxSaturation: 0.0,
            minValue: 0.9,
            maxValue: 1.0,
        },
        {
            name: "Yellow",
            minHue: 39.0,
            maxHue: 67.0,
            minSaturation: 0.857142,
            maxSaturation: 0.5,
            minValue: 0.35,
            maxValue: 1.0,
        },
    }
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

func analyzeCluster(cluster []color.Color, definedColors ColorBoundaries) (generatedColor string, matches []string) {
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
        h, s, v := finalColor.Hsv()

        if color.minHue <= h && h <= color.maxHue && color.minValue <= v && v <= color.maxValue{
            if color.minSaturation < color.maxSaturation{
                if color.minSaturation <= s && s <= color.maxSaturation{
                    results = append(results, color.name)
                }
            } else {
                if color.minSaturation >= s && s >= color.maxSaturation{
                    results = append(results, color.name)
                }
            }
        }
	}

	return finalColor.Hex(), results
}

func euclidianDistance(pOne int, pTwo int, qOne int, qTwo int) float64 {
	// https://en.wikipedia.org/wiki/Euclidean_distance#Two_dimensions
	return math.Sqrt(math.Pow(float64(qOne-pOne), 2) + math.Pow(float64(qTwo-pTwo), 2))
}

func skipIfError(err error) bool {
    if err != nil {
        fmt.Println(err)
        return true
    }
    return false
}

func closeIfError(err error) {
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
        os.Exit(1)
	}
}

func main() {
	start := time.Now()
	rand.Seed(time.Now().Unix())
	saveAs := "sample.png"
	testingCSV := true
	testingListOfImages := false
	k := 5
	numberOfImages := 0
    numberOfColors := 0

	colorDefinitions := retrieveColorBoundaries()

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
            if lineNumber == 0{
                // skip headers
            } else if lineNumber < 1000 {
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
                    numberOfColors += k
				} else {
                    err := writer.Write([]string{
                        "",
                    })
                    closeIfError(err)
					// err := writer.Write([]string{
					// 	sku,
					// 	imageUrl,
					// 	"image",
					// 	"skipped",
					// })
					// closeIfError(err)
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
	log.Printf("Processing %v colors from %v images took %s", numberOfColors, numberOfImages, elapsed)
}
