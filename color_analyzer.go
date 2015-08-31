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
)

import _ "image/jpeg"
import _ "image/png"

type DefinedColor struct{
    hex color.Hex
    variance float32
}

type GeneratedColor struct{
    red uint64
    green uint64
    blue uint64
}

func readColorConfig(configLocation string)(definedColors []DefinedColor){
    /*
        This should read the color values from the config file, but I don't know 
        how to do that yet, so we're going to just create a bunch of structs manually
        and plop them in a slice and feel bad about it until we make it better. Cool?
    */

    black  := DefinedColor{ hex: "#191818", variance: 100.0 }
    brown  := DefinedColor{ hex: "#795000", variance: 36.0  }
    blue   := DefinedColor{ hex: "#3F4AFF", variance: 46.0  }
    gold   := DefinedColor{ hex: "#C1B000", variance: 15.0  }
    gray   := DefinedColor{ hex: "#7D7C7A", variance: 100.0 }
    green  := DefinedColor{ hex: "#1CBD2A", variance: 73.0  }
    orange := DefinedColor{ hex: "#C27B13", variance: 47.0  }
    pink   := DefinedColor{ hex: "#FFBECC", variance: 18.0  }
    purple := DefinedColor{ hex: "#9E4DFF", variance: 29.0  }
    red    := DefinedColor{ hex: "#FF260C", variance: 63.0  }
    tan    := DefinedColor{ hex: "#D18D12", variance: 12.0  }
    white  := DefinedColor{ hex: "#FFFDF7", variance: 4.0   }
    yellow := DefinedColor{ hex: "#FFF000", variance: 18.0  }

    return []DefinedColor{black, brown, blue, gold, gray, green, orange, pink, purple, red, tan, white, yellow}
}

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
    // everything below this line seems ucked up
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

func analyzeClusters(clusters map[int][]color.Color, definedColors []DefinedColor){
    for cluster := range clusters{
        redTotal := uint64(0)
        greenTotal := uint64(0)
        blueTotal := uint64(0)
        pixelTotal := uint64(0)
        for pixel := range clusters[cluster]{
            r, g, b, _ := clusters[cluster][pixel].RGBA()
            
            redTotal += uint64(r >> 8)
            greenTotal += uint64(g >> 8)
            blueTotal += uint64(b >> 8)
            pixelTotal += 1
        }
        
        // generatedColor := calculateFinalHueValues(redTotal, blueTotal, greenTotal, pixelTotal)

        for _, color := range definedColors{
            createGeneratedColorFromDefinedColor(color)
        }
    }
}

func createGeneratedColorFromDefinedColor(color DefinedColor){
    rawHex := color.hex[1:len(color.hex)]
    initialR := rawHex[0:2]
    initialG := rawHex[2:4]
    initialB := rawHex[4:6]

    fmt.Println(initialR, " ", initialG, " ", initialB)
}

func calculateFinalHueValues(red uint64, blue uint64, green uint64, total uint64)(color GeneratedColor){
    r :=   red / total
    g := green / total
    b :=  blue / total

    return GeneratedColor{ red: r, green: g, blue: b }
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
    k := 5
    // _ := "sample_images/red.png"
    // colorConfigFile := "colors.json"

    err := downloadImageFromUrl(url, saveAs)
    closeIfError(err)

    img, err := openImage(saveAs)
    closeIfError(err)

    definedColors := readColorConfig("")
    clusters := createClusters(k, img)
    analyzeClusters(clusters, definedColors)
}