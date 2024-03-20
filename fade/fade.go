package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"time"

	"github.com/stefan-muehlebach/ledgrid"
)

type ColorType int

const (
	Red ColorType = iota
	Green
	Blue
	NumColors
)

var (
	width                = 10
	height               = 10
	defHost              = "raspi-2"
	defPort         uint = 5333
	defGammaValue        = 3.0
	framesPerSecond      = 50
	frameRefreshMs       = 1000 / framesPerSecond
	frameRefreshSec      = float64(frameRefreshMs) / 1000.0
)

//----------------------------------------------------------------------------

// func sphereFunc(x, y, t float64) float64 {
// 	return math.Sqrt(7.0/9.0 - x*x - y*y)
// }

// func testFunc(x, y, t float64) float64 {
// 	return x/2.0 + 0.5
// }

// func verticalBar(x, y, t float64) float64 {
// 	_, f := math.Modf(x/0.25 + 0.5 + t/2)
// 	if f > 0.8 {
// 		return 1.0
// 	} else {
// 		return 0.0
// 	}
// }

type Counter struct {
	size  image.Point
	bits  []bool
	color ledgrid.LedColor
}

func NewCounter(size image.Point, color ledgrid.LedColor) *Counter {
	c := &Counter{}
	c.size = size
	c.bits = make([]bool, c.size.X*c.size.Y)
	c.color = color
	return c
}

func (c *Counter) Update(t float64) {
	for i, b := range c.bits {
		if !b {
			c.bits[i] = true
			break
		} else {
			c.bits[i] = false
		}
	}
}

func (c *Counter) Draw(grid *ledgrid.LedGrid) {
	for i, b := range c.bits {
		if !b {
			continue
		}
		row := i / c.size.X
		col := i % c.size.X
		grid.SetLedColor(col, row, c.color)
	}
}

//----------------------------------------------------------------------------

func main() {
	var host string
	var port uint
	var gammaValue float64

	var client *ledgrid.PixelClient
	var grid *ledgrid.LedGrid
	var pal *ledgrid.PaletteFader
	var shader *ledgrid.Shader
	var ch string
	var palIdx int = 0
	var palName string

	flag.StringVar(&host, "host", defHost, "Controller hostname")
	flag.UintVar(&port, "port", defPort, "Controller port")
	flag.Float64Var(&gammaValue, "gamma", defGammaValue, "Gamma value")
	flag.Parse()

	client = ledgrid.NewPixelClient(host, port)
	client.SetGamma(gammaValue, gammaValue, gammaValue)
	grid = ledgrid.NewLedGrid(image.Rect(0, 0, width, height))

	pal = ledgrid.NewPaletteFader(ledgrid.PaletteMap["Hipster"])
	shader = ledgrid.NewShader(grid.Bounds().Size(), pal, ledgrid.PlasmaShader)

	// counter := NewCounter(grid.Rect.Size(), ledgrid.Red)

	ticker := time.NewTicker(time.Duration(frameRefreshMs) * time.Millisecond)
	go func() {
		t0 := time.Now()
		for t := range ticker.C {
			t1 := t.Sub(t0).Seconds()
			pal.Update(t1)
			shader.Update(t1)
			// counter.Update(t1)
			grid.Clear()
			shader.Draw(grid)
			// counter.Draw(grid)

			client.Draw(grid)
		}
	}()

mainLoop:
	for {
		fmt.Printf("Current palette: %s\n", palName)
		fmt.Printf("  q: prev; w: next\n")
		fmt.Printf("Current gamma value(s) %.4f\n", gammaValue)
		fmt.Printf("  e: decr; r: incr\n")
		fmt.Printf("  x: quit\n")

		_, err := fmt.Scanln(&ch)
		if err != nil {
			log.Fatalf("error reading command:", err)
		}

		switch ch {
		case "q", "w":
			if ch == "q" {
				if palIdx > 0 {
					palIdx -= 1
				} else {
					break
				}
			} else {
				palIdx = (palIdx + 1) % len(ledgrid.PaletteNames)
			}
			palName = ledgrid.PaletteNames[palIdx]
			fmt.Printf("New palette: %s\n", palName)
			// pal = ledgrid.PaletteMap[palName]
			pal.Fade(ledgrid.PaletteMap[palName], 2.0)
			// shader.FadePalette(pal, 2.0)
		case "e", "r":
			if ch == "e" {
				gammaValue -= 0.1
			} else {
				gammaValue += 0.1
			}
			client.SetGamma(gammaValue, gammaValue, gammaValue)
		case "x":
			break mainLoop
		default:
			fmt.Printf("command unknown: '%s'\n", ch)
		}

	}
	ticker.Stop()

	grid.Clear()
	client.Draw(grid)

	client.Close()
}
