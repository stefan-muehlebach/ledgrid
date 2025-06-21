//go:build guiSDL2

package main

import (
	"image"
	"log"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/conf"
)

const (
	FrameRate        = 30
	EventQueueLength = 20
)

type Window struct {
	ledgrid.DisplayEmbed
	Win      *sdl.Window
	Renderer *sdl.Renderer
	size     image.Point
	indexMap conf.IndexMap
	coordMap conf.CoordMap
	ledChain []*Circle
}

// A new grid object must only know it's size in order to get the
// configuration of the emulated modules.
func NewWindowBySize(title string, pixelSize float64, size image.Point) *Window {
	modConf := conf.DefaultModuleConfig(size)
	return NewWindow(title, pixelSize, modConf)
}

func NewWindow(title string, pixelSize float64, modConf conf.ModuleConfig) *Window {
	var err error

	e := &Window{}
	e.DisplayEmbed.Init(e, len(modConf)*conf.ModuleDim.X*conf.ModuleDim.Y)
	e.ModConf = modConf
	e.size = e.ModConf.Size()

	e.Win, err = sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, int32(e.size.X)*int32(pixelSize),
		int32(e.size.Y)*int32(pixelSize), sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("Failed to create window: %s\n", err)
	}

	e.Renderer, err = sdl.CreateRenderer(e.Win, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("Failed to create renderer: %s\n", err)
	}
	e.Renderer.Clear()

	e.coordMap = e.ModConf.CoordMap()
	e.indexMap = e.ModConf.IndexMap()

	// fmt.Printf("Setup the slice with the NeoPixels...\n")
	e.ledChain = make([]*Circle, len(e.coordMap))
	ledSize := image.Point{int(pixelSize) - 2, int(pixelSize) - 2}
	for i, coord := range e.coordMap {
		col, row := coord.X, coord.Y
		ledPos := image.Point{int(pixelSize/2) + col*int(pixelSize), int(pixelSize/2) + row*int(pixelSize)}
		led := NewCircle(int32(ledPos.X), int32(ledPos.Y), int32(ledSize.X/2), sdl.RGB888{200, 200, 200})
		e.ledChain[i] = led
	}

	return e
}

func (e *Window) HandleEvents() {
	// fmt.Printf("Enter main event loop...\n")
	events := make([]sdl.Event, EventQueueLength)
	e.redraw()
	isRunning := true
	for isRunning {
		if event := sdl.WaitEvent(); event != nil {
			// fmt.Printf("  event received: %+v\n", event)
			switch evt := event.(type) {
			case *sdl.QuitEvent:
				// fmt.Printf("    in QuitEvent\n")
				return

			case *sdl.KeyboardEvent:
				if evt.Type != sdl.KEYDOWN {
					continue
				}
				// fmt.Printf("    in KeyboardEvent\n")
				// fmt.Printf("    %+v\n", evt)
				switch evt.Keysym.Sym {
				case sdl.K_ESCAPE, sdl.K_q:
					return

				case sdl.K_t:
					ToggleTests()
				case sdl.K_s:
					PrintStatistics()

				case sdl.K_h:
					println("Key commands")
					println("-----------------")
					println("h     : This help")
					println("t     : Run test programs")
					println("s     : Show statistics")
					println("q/ESC : Quit")
				}

			case *sdl.MouseMotionEvent:
				sdl.PeepEvents(events, sdl.GETEVENT, sdl.MOUSEMOTION, sdl.MOUSEMOTION+1)
				// if n > 0 {
				// 	fmt.Printf(">>> PeepEvents: %d events found\n", n)
				// }

			case *sdl.MouseButtonEvent:

			case *sdl.UserEvent:
				// fmt.Printf("    in UserEvent\n")
				sdl.PeepEvents(events, sdl.GETEVENT, sdl.USEREVENT, sdl.USEREVENT+1)
				e.Renderer.SetDrawColor(0x20, 0x20, 0x20, 0xff)
				e.Renderer.Clear()
				for _, led := range e.ledChain {
					led.Draw(e.Renderer)
				}
				e.Renderer.Present()

				// default:
				// 	fmt.Printf("    no handler for this event type: %+v (%T)\n", evt, evt)
			}
			sdl.Delay(1000 / FrameRate)
		}
	}
}

func (e *Window) DefaultGamma() (r, g, b float64) {
	return 1.0, 1.0, 1.0
}

func (e *Window) Close() {
	e.Win.Destroy()
	e.Renderer.Destroy()
}

// Takes the bytes in buffer, and uses them exactly as the real hardware
// would to recolor the individual LEDs (circle objects) of the emulation.
// Only if the colors really change, a refresh is issued.
func (e *Window) Send(buffer []byte) {
	var r, g, b uint8
	var i int
	var needsRefresh bool
	var newColor, oldColor sdl.RGB888

	needsRefresh = false
	for i = 0; i < len(buffer); i += 3 {
		ledIndex := i / 3
		// coord := e.coordMap[i/3]
		src := buffer[i : i+3 : i+3]
		r = src[0]
		g = src[1]
		b = src[2]
		newColor.R, newColor.G, newColor.B = r, g, b
		if !needsRefresh {
			oldColor = e.ledChain[ledIndex].col
			if newColor != oldColor {
				needsRefresh = true
			}
		}
		e.ledChain[ledIndex].col = newColor
	}
	if needsRefresh {
		e.redraw()
	}
}

func (e *Window) redraw() {
	win := sdl.GetKeyboardFocus()
	id, _ := win.GetID()
	evt := &sdl.UserEvent{sdl.USEREVENT, sdl.GetTicks(), id, 1331, nil, nil}
	sdl.PushEvent(evt)
}

//-----------------------------------------------------------------------------

type Circle struct {
	xm, ym, rad int32
	col         sdl.RGB888
}

func NewCircle(xm, ym, rad int32, col sdl.RGB888) *Circle {
	c := &Circle{}
	c.xm = xm
	c.ym = ym
	c.rad = rad
	c.col = col
	return c
}

func (c *Circle) Draw(r *sdl.Renderer) {
	var offsetx, offsety, d int32

	r.SetDrawColor(c.col.R, c.col.G, c.col.B, 0xff)
	offsetx = 0
	offsety = c.rad
	d = c.rad - 1

	for offsety >= offsetx {
		err := r.DrawLine(c.xm-offsety, c.ym+offsetx, c.xm+offsety, c.ym+offsetx)
		r.DrawLine(c.xm-offsetx, c.ym+offsety, c.xm+offsetx, c.ym+offsety)
		r.DrawLine(c.xm-offsetx, c.ym-offsety, c.xm+offsetx, c.ym-offsety)
		r.DrawLine(c.xm-offsety, c.ym-offsetx, c.xm+offsety, c.ym-offsetx)

		if err != nil {
			return
		}

		if d >= 2*offsetx {
			d -= 2*offsetx + 1
			offsetx += 1
		} else if d < 2*(c.rad-offsety) {
			d += 2*offsety - 1
			offsety -= 1
		} else {
			d += 2 * (offsety - offsetx - 1)
			offsety -= 1
			offsetx += 1
		}
	}
}

func (c *Circle) Draw2(r *sdl.Renderer) {
	var diameter, x, y, tx, ty, err int32

	diameter = 2 * c.rad
	x, y = (c.rad - 1), 0
	tx, ty = 1, 1
	err = (tx - diameter)
	for x >= y {
		r.DrawPoint(c.xm+x, c.ym-y)
		r.DrawPoint(c.xm+x, c.ym+y)
		r.DrawPoint(c.xm-x, c.ym-y)
		r.DrawPoint(c.xm-x, c.ym+y)
		r.DrawPoint(c.xm+y, c.ym-x)
		r.DrawPoint(c.xm+y, c.ym+x)
		r.DrawPoint(c.xm-y, c.ym-x)
		r.DrawPoint(c.xm-y, c.ym+x)

		if err <= 0 {
			y++
			err += ty
			ty += 2
		}
		if err > 0 {
			x--
			tx += 2
			err += (tx - diameter)
		}
	}
}
