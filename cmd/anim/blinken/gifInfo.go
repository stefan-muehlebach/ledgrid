package main

import (
	"log"
	"os"
    "image/gif"
    "fmt"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatalf("usage: %s <file> [...]", os.Args[0])
    }
    for _, fileName := range os.Args[1:] {
        fh, err := os.Open(fileName)
        if err != nil {
            log.Fatalf("Couldn't open file: %v", err)
        }
        gifData, err := gif.DecodeAll(fh)
        if err != nil {
            log.Fatalf("Couldn't decode file: %v", err)
        }
        fh.Close()
        fmt.Printf("GIF informations from '%s':\n", fileName)
        fmt.Printf("  Number of images: %d\n", len(gifData.Image))
        fmt.Printf("  Image size: %dx%d\n", gifData.Config.Width, gifData.Config.Height)
        fmt.Printf("  Loop count: %d\n", gifData.LoopCount)
        fmt.Printf("  Size of color palette: %d\n", len(gifData.Image[0].Palette))
    }
}
