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
	width                = 10
	height               = 10
	defHost              = "raspi-2"
	defPort         uint = 5333
	framesPerSecond      = 50
	frameDelayMs         = 1000 / framesPerSecond
	frameDelaySec        = float64(frameDelayMs) / 1000.0
	gammaValue           = 3.0
)

//-----------------------------------------------------------------------------

type AnimFuncType func(x, y, t float64) float64

func verticalBar(x, y, t float64) float64 {
    _, f := math.Modf(x/0.25 + 0.5 + t/2)
    if f > 0.8 {
        return 1.0
    } else {
        return 0.0
    }
}

func horizontalFade(x, y, t float64) float64 {
    _, f := math.Modf(x/0.25 + 0.5 + t/2)
    return f
}

func plasmaFunc(x, y, t float64) float64 {
	v1 := f1(x, y, t, 10.0)
	v2 := f2(x, y, t, 10.0, 2.0, 3.0)
	v3 := f3(x, y, t, 5.0, 3.0)
	v := (v1+v2+v3)/6.0 + 0.5
	return v
}
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


type Animator struct {
	grid                  *ledgrid.LedGrid
	width, height, dx, dy float64
    AnimFunc  AnimFuncType
	pals                  [2]*ledgrid.Palette
	ft, dft               float64
}

func NewAnimator(grid *ledgrid.LedGrid, pal *ledgrid.Palette, animFunc AnimFuncType) *Animator {
	p := &Animator{}
	p.grid = grid
	p.width, p.height = 0.25, 0.25
	p.dx = p.width / float64(grid.Rect.Dx()-1)
	p.dy = p.height / float64(grid.Rect.Dy()-1)
    p.AnimFunc = animFunc
	p.pals[0] = pal
	p.pals[1] = pal
	p.ft = 0.0
	p.dft = 0.03
	return p
}

func (p *Animator) Animate(t float64) {
	var col, row int
	var x, y float64

	if p.ft > 0.0 {
		p.ft -= frameDelaySec
		if p.ft < 0.0 {
			p.ft = 0.0
		}
	}
	y = p.height / 2.0
	for row = range p.grid.Rect.Dy() {
		x = -p.width / 2.0
		for col = range p.grid.Rect.Dx() {
			v := p.AnimFunc(x, y, t)
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

func (p *Animator) SetPalette(pal *ledgrid.Palette) {
	if p.ft > 0.0 {
		return
	}
	p.pals[0], p.pals[1] = pal, p.pals[0]
	p.ft = 1.0
}

func main() {
	var host string
	var port uint

	var client *ledgrid.PixelClient
	var grid *ledgrid.LedGrid
	var pal *ledgrid.Palette
	var anim *Animator
	var ch string
	var palIdx int = 0
	var palName string

	flag.StringVar(&host, "host", defHost, "Controller hostname")
	flag.UintVar(&port, "port", defPort, "Controller port")
	flag.Parse()

	client = ledgrid.NewPixelClient(host, port)
	client.SetGamma(gammaValue, gammaValue, gammaValue)
	grid = ledgrid.NewLedGrid(image.Rect(0, 0, width, height))

	palName = "FadeRed"
	pal = ledgrid.PaletteMap[palName]
	anim = NewAnimator(grid, pal, horizontalFade)

	ticker := time.NewTicker(time.Duration(frameDelayMs) * time.Millisecond)
	go func() {
		for t := range ticker.C {
			ti := float64(t.UnixMilli()) / 1000.0
			anim.Animate(ti)
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
			pal = ledgrid.PaletteMap[palName]
			anim.SetPalette(pal)
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
