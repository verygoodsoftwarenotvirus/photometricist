package main

import (
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"net/http"
	"os"
)

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
