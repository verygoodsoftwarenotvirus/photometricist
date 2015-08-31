package main

import (
	"fmt"
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
	"strconv"
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

type ColorRange struct {
	minRed   float64
	maxRed   float64
	minBlue  float64
	maxBlue  float64
	minGreen float64
	maxGreen float64
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

	if error != nil {
		log.Fatal(error)
		return nil, error
	}

	defer imgfile.Close()

	img, _, error = image.Decode(imgfile)

	if error != nil {
		log.Fatal(error)
		return nil, error
	}

	return img, nil
}

func cropImage(img image.Image, percentage float64) (croppedImg image.Image, error error) {
	// https://github.com/LiterallyElvis/color-analyzer/blob/master/analysis_objects.py#L63
	percentage = math.Max(percentage, 99.0)

	return nil, nil
}

func createDebugImage(filename string, bounds image.Rectangle, clusterPoints []map[string]int) {
	out, err := os.Create(filename)
	closeIfError(err)
	debugImgOutline := image.Rect(0, 0, bounds.Max.X, bounds.Max.Y)
	debugImg := image.NewGray(debugImgOutline)
	draw.Draw(debugImg, debugImg.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)

	for _, point := range clusterPoints {
		debugPoint := image.Rect(point["X"], point["Y"], point["X"]+5, point["Y"]+5)
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

	createDebugImage("debug.png", bounds, clusterPoints)

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

func analyzeClusters(clusters map[int][]color.Color, definedColors []DefinedColor) {
	for cluster := range clusters {
		redTotal := float64(0.0)
		greenTotal := float64(0.0)
		blueTotal := float64(0.0)
		pixelTotal := float64(0.0)
		for pixel := range clusters[cluster] {
			r, g, b, _ := clusters[cluster][pixel].RGBA()

			redTotal += float64(r >> 8)
			greenTotal += float64(g >> 8)
			blueTotal += float64(b >> 8)
			pixelTotal += 1
		}

		// result := calculateFinalColorValues(redTotal, greenTotal, blueTotal, pixelTotal)

		// for _, color := range definedColors {
		// 	compare := createComparisonFromDefinedColor(color)
		// 	if colorMatchesRanges(result, compare) {
		// 		fmt.Println("\n!!!!!!!!!!!!!!!!!\n")
		// 		fmt.Println("Matched!: ", color.name)
		// 		fmt.Println(uint16(compare.minRed), " <=> ", uint16(result.red), " <=> ", uint16(compare.maxRed))
		// 		fmt.Println(uint16(compare.minGreen), " <=> ", uint16(result.green), " <=> ", uint16(compare.maxGreen))
		// 		fmt.Println(uint16(compare.minBlue), " <=> ", uint16(result.blue), " <=> ", uint16(compare.maxBlue))
		// 		fmt.Println("\n!!!!!!!!!!!!!!!!!\n")
		// 	} else {
		// 		fmt.Println("\nUnmatched!", color.name)
		// 		fmt.Println(uint16(compare.minRed), " <=> ", uint16(result.red), " <=> ", uint16(compare.maxRed))
		// 		fmt.Println(uint16(compare.minGreen), " <=> ", uint16(result.green), " <=> ", uint16(compare.maxGreen))
		// 		fmt.Println(uint16(compare.minBlue), " <=> ", uint16(result.blue), " <=> ", uint16(compare.maxBlue))
		// 	}
		// }
	}
}

func colorMatchesRanges(color GeneratedColor, compare ColorRange) (result bool) {
	if compare.minRed > color.red || color.red > compare.maxRed {
		return false
	} else if compare.minGreen > color.green || color.green > compare.maxGreen {
		return false
	} else if compare.minBlue > color.blue || color.blue > compare.maxBlue {
		return false
	} else {
		return true
	}
}

func createComparisonFromDefinedColor(color DefinedColor) (result ColorRange) {
	rawHex := color.hex[1:len(color.hex)]
	rgb, err := strconv.ParseUint(string(rawHex), 16, 32)
	closeIfError(err)
	red, green, blue := float64(rgb>>16), float64((rgb>>8)&0xFF), float64(rgb&0xFF)
	return ColorRange{
		minRed:   math.Max(red-(red*(color.variance*.01)), 0),
		maxRed:   math.Min(red+(red*(color.variance*.01)), 255),
		minGreen: math.Max(green-(green*(color.variance*.01)), 0),
		maxGreen: math.Min(green+(green*(color.variance*.01)), 255),
		minBlue:  math.Max(blue-(blue*(color.variance*.01)), 0),
		maxBlue:  math.Min(blue+(blue*(color.variance*.01)), 255),
	}
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
	url := "http://www.surlatable.com/images/customers/c1079/PRO-1748698/PRO-1748698_hopup/main_variation_Default_view_2_715x715."
	saveAs := "sample.png"
	// testingListOfImages := false
	k := 5
	// saveAs := "sample_images/red.png"
	// colorConfigFile := "colors.json"

	err := downloadImageFromUrl(url, saveAs)
	closeIfError(err)

	img, err := openImage(saveAs)
	closeIfError(err)

	definedColors := readColorConfig("")
	clusters := createClusters(k, img)
	analyzeClusters(clusters, definedColors)
}
