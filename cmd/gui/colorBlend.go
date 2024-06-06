//go:build ignore

package main

import (
	// "github.com/stefan-muehlebach/ledgrid"
	"image"
	"image/color"
	"golang.org/x/image/draw"
	"image/png"
	"log"
	"os"
)

func SaveImage(m image.Image) {
	fh, err := os.Create("colorBlend.png")
	if err != nil {
		log.Fatal(err)
	}
	if err = png.Encode(fh, m); err != nil {
		fh.Close()
		log.Fatal(err)
	}
	fh.Close()
}

func SetColor(img draw.Image, row0, row1 int, c color.Color) {
	for row := row0; row < row1; row++ {
		for col := range img.Bounds().Dx() {
			img.Set(col, row, c)
		}
	}
}

func main() {
	rect := image.Rect(0, 0, 100, 100)
	img1 := image.NewNRGBA(rect)
    mask1 := image.NewUniform(color.Alpha{255})
    opts1 := &draw.Options {
        SrcMask: mask1,
    }
    	img2 := image.NewNRGBA(rect)
    mask2 := image.NewUniform(color.Alpha{160})
    opts2 := &draw.Options {
        SrcMask: mask2,
    }
    	img3 := image.NewNRGBA(rect)

	SetColor(img1, 0, 100, color.NRGBA{0xff, 0x00, 0x00, 0xff})
    SetColor(img2, 0, 100, color.NRGBA{0x00, 0xff, 0x00, 0xff})

	draw.BiLinear.Scale(img3, img3.Bounds(), img1, img1.Bounds(), draw.Over, opts1)
    draw.BiLinear.Scale(img3, img3.Bounds(), img2, img2.Bounds(), draw.Over, opts2)
    SaveImage(img3)
}
