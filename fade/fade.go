package main

import (
	"flag"
	"image"
	"image/draw"
	"log"
	"time"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/colornames"
	"github.com/stefan-muehlebach/ledgrid"
)

type ColorType int

const (
	Red ColorType = iota
	Green
	Blue
	NumColors
)

const (
	width    = 10
	height   = 10
	defHost  = "raspi-2"
	defPort  = 5333
	defGroup = colornames.Greens
)

func main() {
	var host string
	var port uint
	var colorGroup colornames.ColorGroup = defGroup

	var gc *gg.Context
	var ledGrid *ledgrid.LedGrid
	var pixelClient *ledgrid.PixelClient
	var uniColor *image.Uniform
	var prevColor, nextColor, currColor color.Color
	var radius float64

	flag.StringVar(&host, "host", defHost, "Controller hostname")
	flag.UintVar(&port, "port", defPort, "Controller port")
	flag.Var(&colorGroup, "colors", "Color group")
	flag.Parse()

	gc = gg.NewContext(10, 10)
	gc.SetStrokeWidth(1.5)
	pixelClient = ledgrid.NewPixelClient(host, port)
    pixelClient.SetGamma(0, 1.0)
	ledGrid = ledgrid.NewLedGrid(image.Rect(0, 0, width, height))
	prevColor = color.Black
	uniColor = image.NewUniform(prevColor)

	for _, colorName := range colornames.Groups[colorGroup] {
		log.Printf("[%s]", colorName)
		nextColor = colornames.Map[colorName]
		for t := 0.0; t <= 1.0; t += 0.05 {
			currColor = prevColor.Interpolate(nextColor, f(t))
			if t <= 0.5 {
				radius = (f(2*t) * 5)
			} else {
				radius = ((1 - f(2*t-1)) * 5)
			}
			gc.SetFillColor(color.Black)
			gc.Clear()
			gc.SetStrokeColor(currColor)
			gc.SetFillColor(currColor)
			gc.DrawCircle(5, 5, radius)
			gc.Fill()
			draw.Draw(ledGrid, ledGrid.Bounds(), gc.Image(), image.Point{}, draw.Src)
			pixelClient.Draw(ledGrid)
			time.Sleep(60 * time.Millisecond)
		}
		prevColor = nextColor
	}

	draw.Draw(ledGrid, ledGrid.Bounds(), uniColor, image.Point{}, draw.Src)
	pixelClient.Draw(ledGrid)
	time.Sleep(5 * time.Millisecond)

	pixelClient.Close()
}

func f(t float64) float64 {
	return 3*t*t - 2*t*t*t
}
