package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"math"
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
	width             = 10
	height            = 10
	defHost           = "raspi-2"
	defPort      uint = 5333
	frameRefresh      = 20 * time.Millisecond
	gammaValue        = 3.0
)

//-----------------------------------------------------------------------------

func f1(x, y, t, p1 float64) float64 {
	return math.Sin(x*p1 + t)
}

func f2(x, y, t, p1, p2, p3 float64) float64 {
	return math.Sin(p1*(x*math.Sin(t/p2)+y*math.Cos(t/p3)) + t)
}

func f3(x, y, t, p1, p2 float64) float64 {
	cx := x + 0.5*math.Sin(t/p1)
	cy := y + 0.5*math.Cos(t/p2)
	return math.Sin(math.Sqrt(100.0*(cx*cx+cy*cy)+1.0) + t)
}

type Plasma struct {
	grid                  *ledgrid.LedGrid
	width, height, dx, dy float64
	t, dt                 float64
	pals                  [2]*ledgrid.Palette
	ft, dft               float64
}

func NewPlasma(grid *ledgrid.LedGrid, pal *ledgrid.Palette) *Plasma {
	p := &Plasma{}
	p.grid = grid
	p.width, p.height = 0.25, 0.25
	p.t = 0.0
	p.dt = 0.05
	p.dx = p.width / float64(grid.Rect.Dx()-1)
	p.dy = p.height / float64(grid.Rect.Dy()-1)
	p.pals[0] = pal
	p.pals[1] = nil
	p.ft = 0.0
	p.dft = 0.03
	return p
}

func (p *Plasma) Animate() {
	p.t += p.dt
	if p.ft > 0.0 {
		p.ft -= p.dft
		if p.ft < 0.0 {
			p.ft = 0.0
		}
	}
}

func (p *Plasma) Draw() {
	var col, row int
	var x, y float64

	y = p.height / 2.0
	for row = range p.grid.Rect.Dy() {
		x = -p.width / 2.0
		for col = range p.grid.Rect.Dx() {
			v1 := f1(x, y, p.t, 10.0)
			v2 := f2(x, y, p.t, 10.0, 2.0, 3.0)
			v3 := f3(x, y, p.t, 5.0, 3.0)
			v := (v1+v2+v3)/6.0 + 0.5
			// v := (x / p.width) + 0.5
			c1 := p.pals[0].Color(v)
			if p.ft > 0.0 {
				c2 := p.pals[1].Color(v)
				c1 = c1.Interpolate(c2, p.ft)
			}
			p.grid.SetLedColor(col, row, c1)
			x += p.dx
		}
		y -= p.dy
	}
}

func (p *Plasma) FadePalette(pal *ledgrid.Palette) {
	if p.ft > 0.0 {
		return
	}
	p.pals[0], p.pals[1] = pal, p.pals[0]
	p.ft = 1.0
}

func plasma(client *ledgrid.PixelClient, grid *ledgrid.LedGrid) {
	var pal *ledgrid.Palette
	var plas *Plasma
	var ch string
	var doAnim bool = true
	var palIdx int = 0
	var palName string

	palName = ledgrid.PaletteNames[palIdx]
	pal = ledgrid.PaletteMap[palName]
	plas = NewPlasma(grid, pal)

	for t := 0.0; t <= 1.0; t += 0.05 {
		c := pal.Color(t)
		fmt.Printf("%.2f: %v\n", t, c)
	}

	ticker := time.NewTicker(frameRefresh)
	go func() {
		for range ticker.C {
			if doAnim {
				plas.Animate()
			}
			plas.Draw()
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
		case "s":
			doAnim = false
		case "r":
			doAnim = true
		case "n":
			palIdx = (palIdx + 1) % len(ledgrid.PaletteNames)
			palName = ledgrid.PaletteNames[palIdx]
			fmt.Printf("New palette: %s\n", palName)
			pal = ledgrid.PaletteMap[palName]
			plas.FadePalette(pal)
		case "1":
			pal.Func = ledgrid.LinearInterpol
		case "2":
			pal.Func = ledgrid.PolynomInterpol
		case "3":
			ledgrid.ColorInterpol = ledgrid.LinearInterpol
		case "4":
			ledgrid.ColorInterpol = ledgrid.SqrtInterpol
		case "q":
			break mainLoop
		default:
			fmt.Printf("command unknown: '%s'\n", ch)
		}

	}
	ticker.Stop()
}

func main() {
	var host string
	var port uint

	var ledGrid *ledgrid.LedGrid
	var pixelClient *ledgrid.PixelClient

	flag.StringVar(&host, "host", defHost, "Controller hostname")
	flag.UintVar(&port, "port", defPort, "Controller port")
	flag.Parse()

	pixelClient = ledgrid.NewPixelClient(host, port)
	pixelClient.SetGamma(0, gammaValue)
	pixelClient.SetGamma(1, gammaValue)
	pixelClient.SetGamma(2, gammaValue)
	ledGrid = ledgrid.NewLedGrid(image.Rect(0, 0, width, height))

	plasma(pixelClient, ledGrid)

	ledGrid.Clear()
	pixelClient.Draw(ledGrid)

	pixelClient.Close()
}
