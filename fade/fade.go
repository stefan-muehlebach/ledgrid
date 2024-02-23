package main

import (
	"flag"
	"image"
	"image/draw"
	"log"
	"time"

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
	width   = 10
	height  = 10
	defHost = "raspi-2"
	defPort = 5333
)

func main() {
	var host string
	var port uint

	var ledGrid *ledgrid.LedGrid
	var ctrl *ledgrid.PixelCtrl
	// var img *image.RGBA
	var uniColor *image.Uniform
	var prevColor, nextColor, currColor color.Color

	flag.StringVar(&host, "host", defHost, "Controller hostname")
	flag.UintVar(&port, "port", defPort, "Controller port")
	flag.Parse()

	ctrl = ledgrid.NewPixelCtrl(host, port)
	ledGrid = ledgrid.NewLedGrid(image.Rect(0, 0, width, height))
	// img = image.NewRGBA(image.Rect(0, 0, width, height))
    prevColor = color.Black
	uniColor = image.NewUniform(prevColor)

	for _, colorName := range colornames.Names {
	// for _, colorName := range colornames.Groups[colornames.Greens] {
		log.Printf("[%s]", colorName)
		nextColor = colornames.Map[colorName]
		for t := 0.0; t <= 1.0; t += 0.05 {
			currColor = prevColor.Interpolate(nextColor, t)
			uniColor.C = currColor
			draw.Draw(ledGrid, ledGrid.Bounds(), uniColor, image.Point{}, draw.Src)
			ctrl.Send(ledGrid.Pix)
			time.Sleep(30 * time.Millisecond)
		}
        prevColor = nextColor
		// for row := 0; row < img.Bounds().Dy(); row++ {
		// 	for col := 0; col < img.Bounds().Dx(); col++ {
		// 		x := float64(col) / float64(img.Bounds().Dx()-1)
		// 		img.Set(col, row, c.Dark(x))
		// 	}
		// }
		// draw.Draw(ledGrid, ledGrid.Bounds(), img, image.Point{}, draw.Src)
		// ctrl.Send(ledGrid.Pix)
		// time.Sleep(1000 * time.Millisecond)
	}

	draw.Draw(ledGrid, ledGrid.Bounds(), uniColor, image.Point{}, draw.Src)
	ctrl.Send(ledGrid.Pix)
	time.Sleep(5 * time.Millisecond)

	ctrl.Close()
}
