package main

import (
	// "bufio"
	// "strings"
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

// func readColorConfig(filename string) []colorful.Color {
// 	colors := []colorful.Color{}
// 	file, err := os.Open(filename)
// 	closeIfError(err)
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		colorValues := strings.Split(strings.Trim(strings.Split(scanner.Text(), "]")[0], "[ "), ",")
// 		colorName := strings.Trim(strings.Split(scanner.Text(), "]")[1], " ")

// 		fmt.Println(" ")
// 		fmt.Println("red: ", colorValues[0])
// 		fmt.Println("green: ", colorValues[1])
// 		fmt.Println("blue: ", colorValues[2])
// 		fmt.Println("name: ", colorName)
// 	}

// 	err = scanner.Err()
// 	closeIfError(err)
// 	return colors
// }

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

func openImage(filename string) image.Image {
	imgfile, err := os.Open(filename)
	closeIfError(err)

	defer imgfile.Close()

	img, _, err := image.Decode(imgfile)
	closeIfError(err)

	return img
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

func closeIfError(err error) {
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	saveAs := "sample.png"
	testingListOfImages := true
	k := 5

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
			img := openImage(saveAs)

			croppedImg := cropImage(img, 50, "cropped.png")
			resizedImg := resizeImage(croppedImg, 200, "resized.png")

			clusters := createClusters(k, resizedImg)

			fmt.Printf("\n%v: \n", location)
			analyzeCluster(clusters[0])
			analyzeCluster(clusters[1])
			analyzeCluster(clusters[2])
			analyzeCluster(clusters[3])
			analyzeCluster(clusters[4])

			deleteFileByLocation(saveAs)
		}
	} else {
		url := "http://i.imgur.com/WpsnGdF.jpg"
		downloadImageFromUrl(url, saveAs)
		img := openImage(saveAs)

		croppedImg := cropImage(img, 50, "cropped.png")
		resizedImg := resizeImage(croppedImg, 200, "resized.png")

		clusters := createClusters(k, resizedImg)

		fmt.Printf("\n%v: \n", url)
		analyzeCluster(clusters[0])
		analyzeCluster(clusters[1])
		analyzeCluster(clusters[2])
		analyzeCluster(clusters[3])
		analyzeCluster(clusters[4])

		deleteFileByLocation(saveAs)
	}
}
