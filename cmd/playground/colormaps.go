//go:build ignore

package main

import (
	// "fmt"
	"fmt"
	"image"
	"image/png"

	// "image/color"
	"log"
	"os"
)

var (
	numColors       = 49
	numColumns      = 3
	numColorFields  = 8
	firstColorMap   = image.Point{33, 18}
	colorFieldSize  = image.Point{65, 10}
	colorFieldMp    = colorFieldSize.Div(2)
	colorFieldYDist = image.Point{0, colorFieldSize.Y}
	colorMapSize    = image.Point{65, 80}
	colorMapDist    = image.Point{374, 132}.Sub(firstColorMap)
	colorMapXDist   = image.Point{colorMapDist.X, 0}
	colorMapYDist   = image.Point{0, colorMapDist.Y}
    colorNames = []string{
        "brbg", "prgn", "piyg", "puor", "rdbu", "rdgy",
        "rdylbu", "rdylgn", "spectral", "accent", "dark2", "paired",
        "pastel1", "pastel2", "set1", "set2", "set3", "blues",
        "bugn", " bupu", "gnpu", "greens", "greys", "orrd",
        "oranges", "pubu", "pubugn", "purd", "purples", "rdpu",
        "reds", "ylgn", "ylgnbu", "ylorbr", "ylorrd", "moreland",
        "bentcoolwarm", "jet", "turbo", "parula", "chromajs", "viridis",
        "plasma", "magma", "inferno", "whylrd", "ylrd", "gnpu",
        "sand",
    }
)

func main() {
	fh, err := os.Open("colormaps.png")
	if err != nil {
		log.Fatalf("Couldn't open 'colormaps.png': %v", err)
	}
	defer fh.Close()
	img, err := png.Decode(fh)
	if err != nil {
		log.Fatalf("Couldn't decode file: %v", err)
	}

	for i := range numColors {
		fmt.Printf("  {\n")
		fmt.Printf("    \"ID\": %d,\n", i+90)
		fmt.Printf("    \"Title\": \"%s\",\n", colorNames[i])
		fmt.Printf("    \"Colors\": [\n")
		row := i / numColumns
		col := i % numColumns
		refPt := firstColorMap.Add(colorMapXDist.Mul(col)).Add(colorMapYDist.Mul(row))
		for j := range numColorFields {
			pt := refPt.Add(colorFieldYDist.Mul(j)).Add(colorFieldMp)
			col := img.(*image.RGBA).RGBAAt(pt.X, pt.Y)
			fmt.Printf("        \"%02X%02X%02X\"", col.R, col.G, col.B)
			if j < numColorFields-1 {
				fmt.Printf(",\n")
			} else {
				fmt.Printf("\n")
			}
		}
		fmt.Printf("    ]\n")
		fmt.Printf("  },\n")
	}
}
