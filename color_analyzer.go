package main

import (
	"bufio"
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

func readColorConfig(filename string) []colorful.Color {
	colors := []colorful.Color{}
	file, err := os.Open(filename)
	closeIfError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		colorValues := strings.Split(strings.Trim(strings.Split(scanner.Text(), "]")[0], "[ "), ",")
		colorName := strings.Trim(strings.Split(scanner.Text(), "]")[1], " ")

		fmt.Println(" ")
		fmt.Println("red: ", colorValues[0])
		fmt.Println("green: ", colorValues[1])
		fmt.Println("blue: ", colorValues[2])
		fmt.Println("name: ", colorName)
	}

	err = scanner.Err()
	closeIfError(err)
	return colors
}

func downloadImageFromUrl(url string, saveAs string) (error error) {
	response, err := http.Get(url)
	closeIfError(err)
	defer response.Body.Close()

	file, err := os.Create(saveAs)
	closeIfError(err)

	// dump the response body to the file
	_, err = io.Copy(file, response.Body)
	closeIfError(err)

	file.Close()
	return nil
}

func deleteImageByLocation(location string) (error error) {
	error = os.Remove(location)
	closeIfError(error)
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
	debugImg := image.NewRGBA(debugImgOutline)
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

func analyzeCluster(cluster []color.Color) {
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
	fmt.Println(finalColor.Hex())
}

func euclidianDistance(pOne int, pTwo int, qOne int, qTwo int) float64 {
	// https://en.wikipedia.org/wiki/Euclidean_distance#Two_dimensions
	return math.Sqrt(math.Pow(float64(qOne-pOne), 2) + math.Pow(float64(qTwo-pTwo), 2))
}

func closeIfError(error error) {
	if error != nil {
		log.Fatal(error)
		fmt.Println(error)
		os.Exit(1)
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	saveAs := "sample.png"
	deleteImageByLocation(saveAs)
	testingListOfImages := false
	k := 5

	if testingListOfImages {
		url := "http://i.imgur.com/WpsnGdF.jpg"
		err := downloadImageFromUrl(url, saveAs)
		closeIfError(err)

		img, err := openImage(saveAs)
		closeIfError(err)

		croppedImg, err := cropImage(img, 50)
		resizedImg := imaging.Resize(croppedImg, 200, 200, imaging.Lanczos)

		clusters := createClusters(k, resizedImg)

		analyzeCluster(clusters[0])
		analyzeCluster(clusters[1])
		analyzeCluster(clusters[2])
		analyzeCluster(clusters[3])
		analyzeCluster(clusters[4])
	} else {
		imageLocations := []string{"http://i.imgur.com/WpsnGdF.jpg"}
		for _, location := range imageLocations {
			err := downloadImageFromUrl(location, saveAs)
			closeIfError(err)

			img, err := openImage(saveAs)
			closeIfError(err)

			croppedImg, err := cropImage(img, 50)
			resizedImg := imaging.Resize(croppedImg, 200, 200, imaging.Lanczos)

			clusters := createClusters(k, resizedImg)

			analyzeCluster(clusters[0])
			analyzeCluster(clusters[1])
			analyzeCluster(clusters[2])
			analyzeCluster(clusters[3])
			analyzeCluster(clusters[4])
			deleteImageByLocation(saveAs)
		}
	}
}
