package main

import (
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
	"time"
)

// These have to be imported to avoid unknown image format errors.
import _ "image/jpeg"
import _ "image/gif"

type DefinedColor struct {
	name     string
	hex      string
	variance float64
}

type GeneratedColor struct {
	red   float64
	green float64
	blue  float64
}

func readColorConfig(configLocation string) (definedColors []DefinedColor) {
	/*
	   This should read the color values from the config file, but I don't know
	   how to do that yet, so we're going to just create a bunch of structs manually
	   and plop them in a slice and feel bad about it until we make it better. Cool?
	*/

	black := DefinedColor{name: "Black", hex: "#191818", variance: 100.0}
	brown := DefinedColor{name: "Brown", hex: "#795000", variance: 36.0}
	blue := DefinedColor{name: "Blue", hex: "#3F4AFF", variance: 46.0}
	gold := DefinedColor{name: "Gold", hex: "#C1B000", variance: 15.0}
	gray := DefinedColor{name: "Gray", hex: "#7D7C7A", variance: 30.0}
	green := DefinedColor{name: "Green", hex: "#1CBD2A", variance: 73.0}
	orange := DefinedColor{name: "Orange", hex: "#C27B13", variance: 47.0}
	pink := DefinedColor{name: "Pink", hex: "#FFBECC", variance: 18.0}
	purple := DefinedColor{name: "Purple", hex: "#9E4DFF", variance: 29.0}
	red := DefinedColor{name: "Red", hex: "#FF260C", variance: 63.0}
	tan := DefinedColor{name: "Tan", hex: "#D18D12", variance: 12.0}
	white := DefinedColor{name: "White", hex: "#FFFDF7", variance: 4.0}
	yellow := DefinedColor{name: "Yellow", hex: "#FFF000", variance: 18.0}

	return []DefinedColor{black, brown, blue, gold, gray, green, orange, pink, purple, red, tan, white, yellow}
}

func downloadImageFromUrl(url string, saveAs string) (error error) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer response.Body.Close()

	// open a file for writing
	file, err := os.Create(saveAs)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// dump the response body to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
		return err
	}

	file.Close()
	return nil
}

func deleteImageByLocation(location string) (error error) {
	error = os.Remove(location)
	if error != nil {
		return error
	}
	return nil
}

func openImage(filename string) (img image.Image, error error) {
	imgfile, error := os.Open(filename)
	closeIfError(error)

	defer imgfile.Close()

	img, _, error = image.Decode(imgfile)
	closeIfError(error)

	return img, error
}

func cropImage(img image.Image, percentage float64) (croppedImg image.Image, error error) {
	// https://github.com/LiterallyElvis/color-analyzer/blob/master/analysis_objects.py#L63
	bounds := img.Bounds()
	width := int(float64(bounds.Max.X) * (percentage * .01))
	height := int(float64(bounds.Max.Y) * (percentage * .01))

	out, err := os.Create("testcrop.png")
	err = png.Encode(out, imaging.CropCenter(img, width, height))
	closeIfError(err)

	return imaging.CropCenter(img, width, height), nil
}

func createDebugImage(filename string, bounds image.Rectangle, clusterPoints []map[string]int) {
	out, err := os.Create(filename)
	closeIfError(err)
	debugImgOutline := image.Rect(0, 0, bounds.Max.X, bounds.Max.Y)
	// debugImg := image.NewRGBA(debugImgOutline)
	debugImg := image.NewGray(debugImgOutline)
	draw.Draw(debugImg, debugImg.Bounds(), &image.Uniform{color.Transparent}, image.ZP, draw.Src)

	for _, point := range clusterPoints {
		debugPointOutline := image.Rect(point["X"], point["Y"], point["X"]+5, point["Y"]+5)
		debugPoint := image.Rect(point["X"]+1, point["Y"]+1, point["X"]+4, point["Y"]+4)
		draw.Draw(debugImg, debugPointOutline, &image.Uniform{color.White}, image.Point{X: point["X"], Y: point["Y"]}, draw.Src)
		draw.Draw(debugImg, debugPoint, &image.Uniform{color.Black}, image.Point{X: point["X"], Y: point["Y"]}, draw.Src)
	}

	err = png.Encode(out, debugImg)
	closeIfError(err)
}

func createClusters(numberOfClusters int, img image.Image) (completeClusters map[int][]color.Color) {
	// everything below this line seems ucked up
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
	// everything above this line seems fucked up

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

func analyzeCluster(cluster []color.Color, definedColors []DefinedColor) {
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
	smallestDistance := math.MaxFloat64
	closestColor := ""
	for _, definedColor := range definedColors {
		tempColor, _ := colorful.Hex(definedColor.hex)
		if finalColor.DistanceCIE94(tempColor) < smallestDistance {
			closestColor = definedColor.name
		}
	}
	fmt.Println(closestColor)
}

func calculateFinalColorValues(red float64, green float64, blue float64, total float64) (color GeneratedColor) {
	return GeneratedColor{red: red / total, green: green / total, blue: blue / total}
}

func euclidianDistance(pOne int, pTwo int, qOne int, qTwo int) float64 {
	// https://en.wikipedia.org/wiki/Euclidean_distance#Two_dimensions
	return math.Sqrt(math.Pow(float64(qOne-pOne), 2) + math.Pow(float64(qTwo-pTwo), 2))
}

func closeIfError(error error) {
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	deleteImageByLocation("sample.png")
	deleteImageByLocation("debug.png")
	url := "http://i.imgur.com/WpsnGdF.jpg"
	saveAs := "sample.png"
	// testingListOfImages := false
	k := 5
	// saveAs := "sample_images/green.png"
	// colorConfigFile := "colors.json"

	err := downloadImageFromUrl(url, saveAs)
	closeIfError(err)

	img, err := openImage(saveAs)
	closeIfError(err)

	croppedImg, err := cropImage(img, 50)
	resizedImg := imaging.Resize(croppedImg, 200, 200, imaging.Lanczos)

	definedColors := readColorConfig("")
	clusters := createClusters(k, resizedImg)

	go analyzeCluster(clusters[0], definedColors)
	go analyzeCluster(clusters[1], definedColors)
	go analyzeCluster(clusters[2], definedColors)
	go analyzeCluster(clusters[3], definedColors)
	go analyzeCluster(clusters[4], definedColors)
}
