package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"log"
	"math"
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
	width        = 10
	height       = 10
	defHost      = "raspi-2"
	defPort      = 5333
	defGroup     = colornames.Reds
	frameRefresh = 30 * time.Millisecond
)

type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
)

type MovingBar struct {
	Grid               *ledgrid.LedGrid
	Color              ledgrid.LedColor
	Orient             Orientation
	Pos, DirSpeed, Max float64
}

func NewMovingBar(grid *ledgrid.LedGrid, orient Orientation, dirSpeed float64, color ledgrid.LedColor) *MovingBar {
	b := &MovingBar{}
	b.Grid = grid
	b.Color = color
	b.Orient = orient
	b.Pos = 0.0
	b.DirSpeed = dirSpeed
	switch b.Orient {
	case Horizontal:
		b.Max = float64(b.Grid.Rect.Dy())
	case Vertical:
		b.Max = float64(b.Grid.Rect.Dx())
	}
	return b
}

func (b *MovingBar) Animate() {
	if (b.DirSpeed > 0 && b.Pos+b.DirSpeed > b.Max-1) || (b.DirSpeed < 0 && b.Pos+b.DirSpeed < 0) {
		b.DirSpeed = -b.DirSpeed
	}
	b.Pos += b.DirSpeed
}

func (b *MovingBar) Draw() {
	var size, x, y int
	var posInt, posPart float64
	var color ledgrid.LedColor

	switch b.Orient {
	case Horizontal:
		size = b.Grid.Rect.Dx()
	case Vertical:
		size = b.Grid.Rect.Dy()
	}
	// posInt = math.Round(b.Pos)
	// posPart = 0.0
	posInt, posPart = math.Modf(b.Pos)
	color = b.Color.Interpolate(ledgrid.LedColor{}, posPart)
	// color = b.Color
	for i := range size {
		switch b.Orient {
		case Horizontal:
			x, y = i, int(posInt)
		case Vertical:
			x, y = int(posInt), i
		}
		b.Grid.SetLedColor(x, y, b.Grid.LedColorAt(x, y).Mix(color))
	}
	if posPart != 0.0 && posInt < b.Max-1 {
		color = b.Color.Interpolate(ledgrid.LedColor{}, 1-posPart)
		for i := range size {
			switch b.Orient {
			case Horizontal:
				x, y = i, int(posInt+1)
			case Vertical:
				x, y = int(posInt+1), i
			}
			b.Grid.SetLedColor(x, y, b.Grid.LedColorAt(x, y).Mix(color))
		}
	}
}

type TimedObject struct {
	Obj          *Shape
	Active       bool
	SlotStart    time.Time
	DurationList []time.Duration
	StatusList   []bool
	Cycle        bool
	currSlot     int
}

func NewTimedObject(obj *Shape, durations []int) *TimedObject {
	t := &TimedObject{}
	t.Obj = obj
	t.Active = true
	t.SlotStart = time.Now()
	t.DurationList = make([]time.Duration, 0)
	t.StatusList = make([]bool, 0)
	t.Cycle = false
	for _, dur := range durations {
		switch {
		case dur == 0:
			t.Cycle = true
		case dur < 0:
			t.StatusList = append(t.StatusList, false)
			t.DurationList = append(t.DurationList, time.Duration(-dur)*time.Millisecond)
		case dur > 0:
			t.StatusList = append(t.StatusList, true)
			t.DurationList = append(t.DurationList, time.Duration(dur)*time.Millisecond)
		}
	}
	t.currSlot = 0
	return t
}

func (t *TimedObject) Reset() {
	t.Active = true
	t.SlotStart = time.Now()
	t.currSlot = 0
}

func (t *TimedObject) Animate() {
	if !t.Active {
		return
	}
	if time.Since(t.SlotStart) >= t.DurationList[t.currSlot] {
		t.SlotStart = time.Now()
		t.currSlot++
		if t.currSlot == len(t.DurationList) {
			t.currSlot = 0
			if !t.Cycle {
				t.Active = false
			}
		}
	}
}

func (t *TimedObject) Draw() {
	if !t.Active || !t.StatusList[t.currSlot] {
		return
	}
	t.Obj.Draw()
}

type Shape struct {
	Grid    *ledgrid.LedGrid
	Pos     image.Point
	Color   ledgrid.LedColor
	Pattern []image.Point
}

func NewShape(grid *ledgrid.LedGrid, pos image.Point, c ledgrid.LedColor) *Shape {
	s := &Shape{}
	s.Grid = grid
	s.Pos = pos
	s.Color = c
	s.Pattern = []image.Point{
		image.Point{0, 0},
		image.Point{1, 0},
		image.Point{0, 1},
		image.Point{1, 1},
	}
	return s
}

func (s *Shape) Draw() {
	for _, pt := range s.Pattern {
		pt = pt.Add(s.Pos)
		s.Grid.SetLedColor(pt.X, pt.Y, s.Color)
	}
}

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

func plasma(client *ledgrid.PixelClient, grid *ledgrid.LedGrid) {
	var col, row int
	var x, y, dx, dy float64
	var v1, v2, v3, v float64
	var t, dt float64
	var mathWidth, mathHeight float64
	var pal *ledgrid.Palette

	t = 0.0
	dt = 0.05
	// dDarkFactor = 0.01
	mathWidth = 0.75
	mathHeight = 0.75
	dx = mathWidth / float64(grid.Rect.Dx()-1)
	dy = mathHeight / float64(grid.Rect.Dy()-1)
	pal = ledgrid.NewPalette()
	pal.SetColorStops(ledgrid.Nightspell)

	ticker := time.NewTicker(frameRefresh)
	go func() {
		for range ticker.C {
			y = mathHeight / 2.0
			for row = range grid.Rect.Dy() {
				x = -mathWidth / 2.0
				for col = range grid.Rect.Dx() {
					v1 = f1(x, y, t, 10.0)
					v2 = f2(x, y, t, 10.0, 2.0, 3.0)
					v3 = f3(x, y, t, 5.0, 3.0)
					v = (v1+v2+v3)/6.0 + 0.5
					grid.SetLedColor(col, row, pal.Color(v))
					x += dx
				}
				y -= dy
			}
			t += dt
			client.Draw(grid)
			//pal.SetDarkFactor(pal.DarkFactor() + dDarkFactor)
		}
	}()
	fmt.Scanln()
	ticker.Stop()
}

func movingBar(client *ledgrid.PixelClient, grid *ledgrid.LedGrid) {
	var b1, b2, b3, b4 *MovingBar
	var speed float64 = 0.25

	b1 = NewMovingBar(grid, Horizontal, +speed, ledgrid.LedColor{0, 0, 255})
	b2 = NewMovingBar(grid, Horizontal, -speed, ledgrid.LedColor{255, 0, 255})
	b2.Pos = float64(grid.Rect.Max.Y - 1)

	b3 = NewMovingBar(grid, Vertical, +speed, ledgrid.LedColor{255, 0, 0})
	b4 = NewMovingBar(grid, Vertical, -speed, ledgrid.LedColor{255, 255, 0})
	b4.Pos = float64(grid.Rect.Max.X - 1)

	ticker := time.NewTicker(frameRefresh)
	go func() {
		for range ticker.C {
			grid.Fade()
			b1.Draw()
			b2.Draw()
			b3.Draw()
			b4.Draw()
			client.Draw(grid)
			b1.Animate()
			b2.Animate()
			b3.Animate()
			b4.Animate()
		}
	}()
	fmt.Scanln()
	ticker.Stop()
}

func blurring(client *ledgrid.PixelClient, grid *ledgrid.LedGrid) {
	// var rnd float64
	var shape1, shape2, shape3, shape4 *Shape
	var obj1, obj2, obj3, obj4 *TimedObject

	shape1 = NewShape(grid, image.Point{1, 1}, ledgrid.LedColor{240, 0, 0})
	shape2 = NewShape(grid, image.Point{3, 1}, ledgrid.LedColor{0, 240, 0})
	shape3 = NewShape(grid, image.Point{5, 1}, ledgrid.LedColor{0, 0, 240})
	shape4 = NewShape(grid, image.Point{7, 1}, ledgrid.LedColor{240, 0, 240})

	obj1 = NewTimedObject(shape1, []int{500, -500, 0})
	obj2 = NewTimedObject(shape2, []int{1000, -1000, 0})
	obj3 = NewTimedObject(shape3, []int{2000, -2000, 0})
	obj4 = NewTimedObject(shape4, []int{4000, -4000, 0})
	// obj3 = NewTimedObject(blk3, []int{2000, 300, 2000, 300, 2000, 1000, 300, 300, 300, 0})

	ticker := time.NewTicker(frameRefresh)
	go func() {
		for range ticker.C {
			grid.Clear()
			// grid.Blur()
			obj1.Draw()
			obj2.Draw()
			obj3.Draw()
			obj4.Draw()
			client.Draw(grid)
			obj1.Animate()
			obj2.Animate()
			obj3.Animate()
			obj4.Animate()
		}
	}()
	fmt.Scanln()
	ticker.Stop()
}

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
	pixelClient = ledgrid.NewPixelClient(host, port)
	pixelClient.SetGamma(0, 3.0)
	pixelClient.SetGamma(1, 3.0)
	pixelClient.SetGamma(2, 3.0)
	ledGrid = ledgrid.NewLedGrid(image.Rect(0, 0, width, height))
	prevColor = color.Black
	uniColor = image.NewUniform(prevColor)

	// movingBar(pixelClient, ledGrid)
	// blurring(pixelClient, ledGrid)
	plasma(pixelClient, ledGrid)

	ledGrid.Clear()
	pixelClient.Draw(ledGrid)

	return

	for _, colorName := range colornames.Groups[colorGroup] {
		log.Printf("[%s]", colorName)
		nextColor = colornames.Map[colorName]
		for t := 0.0; t <= 1.0; t += 0.05 {
			currColor = nextColor.Alpha(1 - f(t))
			radius = f(t) * 10
			gc.SetFillColor(color.Black)
			gc.Clear()
			// gc.SetStrokeColor(currColor)
			gc.SetFillColor(currColor)
			gc.DrawCircle(5, 5, radius)
			gc.Fill()
			draw.Draw(ledGrid, ledGrid.Bounds(), gc.Image(), image.Point{}, draw.Src)
			pixelClient.Draw(ledGrid)
			time.Sleep(80 * time.Millisecond)
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
