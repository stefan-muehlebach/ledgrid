package main

import (
	"log"
	"os"
    // "image"
    "image/gif"
    "fmt"
)

func main() {
    fh, err := os.Open("torus.gif")
    if err != nil {
        log.Fatalf("Couldn't open file: %v", err)
    }
    gifData, err := gif.DecodeAll(fh)
    if err != nil {
        log.Fatalf("Couldn't decode file: %v", err)
    }
    fmt.Printf("gifData: %+v\n", gifData)
    for i := range gifData.Image {
        fmt.Printf("%3d: %+v\n", i, gifData.Image[i])
    }
    fh.Close()
}
