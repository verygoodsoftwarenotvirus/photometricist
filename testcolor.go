package main

import (
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
)

func main() {
    purple := colorful.Hex("#9E4DFF")
    fmt.Println("")
    fmt.Println("purple: ")
    fmt.Println(purple.DistanceLuv(colorful.Hex("#7036b5")))
    fmt.Println(purple.DistanceLuv(colorful.Hex("#7137b6")))

    black := colorful.Hex("#191818")
    fmt.Println("")
    fmt.Println("black: ")
    fmt.Println(black.DistanceLuv(colorful.Hex("#141313")))
    fmt.Println(black.DistanceLuv(colorful.Hex("#141414")))

    white := colorful.Hex("#FFFDF7")
    fmt.Println("")
    fmt.Println("white: ")
    fmt.Println(white.DistanceLuv(colorful.Hex("#f4f2ed")))
    fmt.Println(white.DistanceLuv(colorful.Hex("#f5f3ee")))

    tan := colorful.Hex("#D18D12")
    fmt.Println("")
    fmt.Println("tan: ")
    fmt.Println(tan.DistanceLuv(colorful.Hex("#b77c0f")))
    fmt.Println(tan.DistanceLuv(colorful.Hex("#b87d10")))

    brown := colorful.Hex("#795000")
    fmt.Println("")
    fmt.Println("brown: ")
    fmt.Println(brown.DistanceLuv(colorful.Hex("#4d3300")))
    fmt.Println(brown.DistanceLuv(colorful.Hex("#4e3400")))

    gray := colorful.Hex("#7D7C7A")
    fmt.Println("")
    fmt.Println("gray: ")
    fmt.Println(gray.DistanceLuv(colorful.Hex("#575655")))
    fmt.Println(gray.DistanceLuv(colorful.Hex("#585756")))

    yellow := colorful.Hex("#FFF000")
    fmt.Println("")
    fmt.Println("yellow ")
    fmt.Println(yellow.DistanceLuv(colorful.Hex("#d1c400")))
    fmt.Println(yellow.DistanceLuv(colorful.Hex("#d2c500")))

    orange := colorful.Hex("#C27B13")
    fmt.Println("")
    fmt.Println("orange ")
    fmt.Println(orange.DistanceLuv(colorful.Hex("#66410a")))
    fmt.Println(orange.DistanceLuv(colorful.Hex("#67420b")))

    gold := colorful.Hex("#C1B000")
    fmt.Println("")
    fmt.Println("gold: ")
    fmt.Println(gold.DistanceLuv(colorful.Hex("#a49500")))
    fmt.Println(gold.DistanceLuv(colorful.Hex("#a59600")))

    green := colorful.Hex("#1CBD2A")
    fmt.Println("")
    fmt.Println("green: ")
    fmt.Println(green.DistanceLuv(colorful.Hex("#07330b")))
    fmt.Println(green.DistanceLuv(colorful.Hex("#08340c")))

    pink := colorful.Hex("#FFBECC")
    fmt.Println("")
    fmt.Println("pink: ")
    fmt.Println(pink.DistanceLuv(colorful.Hex("#d19ba7")))
    fmt.Println(pink.DistanceLuv(colorful.Hex("#d29ca8")))

    blue := colorful.Hex("#3F4AFF")
    fmt.Println("")
    fmt.Println("blue: ")
    fmt.Println(blue.DistanceLuv(colorful.Hex("#222789")))
    fmt.Println(blue.DistanceLuv(colorful.Hex("#23288a")))

    red := colorful.Hex("#FF260C")
    fmt.Println("")
    fmt.Println("red: ")
    fmt.Println(red.DistanceLuv(colorful.Hex("#5e0e04")))
    fmt.Println(red.DistanceLuv(colorful.Hex("#5f0f05")))
}
