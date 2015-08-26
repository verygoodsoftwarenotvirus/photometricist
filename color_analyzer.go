package main

import (
    "os"
    "io"
    "fmt"
    "log"
    "math"
    "image"
    "net/http"
    "math/rand"
    "image/color"

    // "reflect"
)

import _ "image/jpeg"
import _ "image/png"

func downloadImageFromUrl(url string, saveAs string)(error error){
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

func deleteImageByLocation(location string) (error error){
    error = os.Remove(location)
    if error != nil{
        return error
    }
    return nil
}

func openImage(filename string) (img image.Image, error error){
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

func cropImage(img image.Image, percentage float64)(croppedImg image.Image, error error){
    // https://github.com/LiterallyElvis/color-analyzer/blob/master/analysis_objects.py#L63
    percentage = math.Max(percentage, 99.0)

    return nil, nil
}

func createClusters(numberOfClusters int, img image.Image) (completeClusters map[int][]color.Color){

    // everything below this line seems fucked up
    clusters := make(map[int][]color.Color, numberOfClusters)
    clusterPoints := make([]map[string]int, numberOfClusters)
    bounds := img.Bounds()

    for i := 0; i < numberOfClusters; i++{
        clusters[i] = []color.Color{}
        clusterPoints[i] = map[string]int{
            "X": rand.Intn(bounds.Max.X - bounds.Min.X) + bounds.Min.X,
            "Y": rand.Intn(bounds.Max.Y - bounds.Min.Y) + bounds.Min.Y,
        }
    }
    // everything above this line seems fucked up

    smallestDistanceIndex := math.MaxInt32
    smallestDistance := math.MaxFloat64

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
      for x := bounds.Min.X; x < bounds.Max.X; x++ {

        smallestDistanceIndex = math.MaxInt32
        smallestDistance = math.MaxFloat64

        for index, point := range clusterPoints{
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

func analyzeClusters(clusters map[int][]color.Color){

    for cluster := range clusters{
        redTotal := uint32(0)
        greenTotal := uint32(0)
        blueTotal := uint32(0)
        pixelTotal := uint32(0)
        for pixel := range clusters[cluster]{
            r, g, b, _ := clusters[cluster][pixel].RGBA()
            
            redTotal += uint32(r >> 8)
            greenTotal += uint32(g >> 8)
            blueTotal += uint32(b >> 8)
            pixelTotal += 1
        }
        fmt.Println(cluster, " index cluster averages (R,G,B): ", redTotal / pixelTotal, greenTotal / pixelTotal, blueTotal / pixelTotal)
    }
}

func euclidianDistance(pOne int, pTwo int, qOne int, qTwo int)(float64){
    // https://en.wikipedia.org/wiki/Euclidean_distance#Two_dimensions
    return math.Sqrt( math.Pow( float64(qOne - pOne), 2) + math.Pow(float64(qTwo - pTwo), 2) )
}

func closeIfError(error error){
    if error != nil {
        fmt.Println(error)
        os.Exit(1)
    }
}

func main(){
    url := "http://i.imgur.com/GVzew0Y.jpg"
    saveAs := "sample.png"
    _ := "sample_images/red.png"
    k := 5

    err := downloadImageFromUrl(url, saveAs)
    closeIfError(err)

    img, err := openImage(saveAs)
    closeIfError(err)

    clusters := createClusters(k, img)
    analyzeClusters(clusters)
}