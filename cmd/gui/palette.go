package main

import (
	"image/color"
	"image"
	"log"

	"github.com/stefan-muehlebach/ledgrid"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const (
	minHeight   = 30
	borderWidth = 1.0
	colStopDia  = 20
)

var (
	colStopSize   = fyne.NewSize(colStopDia, colStopDia)
	colStopPosOff = fyne.NewSize(-colStopDia/2, -colStopDia/2)
)

// -----------------------------------------------------------------------------
//
// Objekt-Himmel
type Palette struct {
	widget.BaseWidget
	Orientation widget.Orientation
	ColorSource ledgrid.ColorSource
}

func NewPalette(colSource ledgrid.ColorSource) *Palette {
	var cs ledgrid.ColorSource

	switch pal := colSource.(type) {
	case *ledgrid.PaletteFader:
		cs = pal.Pals[0]
	default:
		cs = pal
	}
	p := &Palette{
		Orientation: widget.Horizontal,
		ColorSource: cs,
	}
	p.ExtendBaseWidget(p)
	return p
}

func (p *Palette) Tapped(evt *fyne.PointEvent) {
	log.Printf("tapped at %v\n", evt)
	log.Printf("  size: %v\n", p.Size())
	log.Printf("  fraction within: %.4f\n", evt.Position.X/p.Size().Width)
    colorPicker := dialog.NewColorPicker("Title", "Message", func(c color.Color){}, Win)
    colorPicker.Advanced = true
    colorPicker.Show()
}

func (p *Palette) CreateRenderer() fyne.WidgetRenderer {
	p.ExtendBaseWidget(p)

	renderer := &paletteRenderer{pal: p}

	gradient := canvas.NewRaster(renderer.generator)
    renderer.gradient = gradient

	rect := canvas.NewRectangle(ledgrid.Transparent)
	rect.StrokeColor = ledgrid.White
	rect.StrokeWidth = borderWidth
	rect.SetMinSize(fyne.NewSize(0, minHeight))
    renderer.rect = rect

	renderer.Refresh()
	return renderer
}

// type ColorStop struct {
// 	widget.BaseWidget
// 	ColStop *ledgrid.ColorStop
// 	Pal     *Palette
// }

// func NewColorStop(pal *Palette, colStop *ledgrid.ColorStop) *ColorStop {
// 	c := &ColorStop{
// 		Pal: pal, ColStop: colStop,
// 	}
// 	c.ExtendBaseWidget(c)
// 	return c
// }

// func (c *ColorStop) Tapped(evt *fyne.PointEvent) {
// 	log.Printf("tapped at %v\n", evt)
// 	log.Printf("  size: %v\n", c.Size())
// }

// func (c *ColorStop) CreateRenderer() fyne.WidgetRenderer {
// 	return nil
// 	// return newColorStopRenderer(c)
// }

// -----------------------------------------------------------------------------
//
// Render-Keller
type paletteRenderer struct {
	// objects    []fyne.CanvasObject
	gradient   *canvas.Raster
	rect       *canvas.Rectangle
	colorStops []*canvas.Circle
	pal        *Palette
}

func (r *paletteRenderer) Destroy() {
}

func (r *paletteRenderer) Layout(s fyne.Size) {
	r.gradient.Resize(s)
	r.rect.Resize(s)

	switch pal := r.pal.ColorSource.(type) {
    case *ledgrid.GradientPalette:
		for i, cs := range pal.ColorStops() {
			pos := r.rect.Position().AddXY(float32(cs.Pos)*r.rect.Size().Width, r.rect.Size().Height/2)
			pos = pos.Add(colStopPosOff)
			r.colorStops[i].Move(pos)
		}
	}
}

func (r *paletteRenderer) MinSize() fyne.Size {
	return r.rect.MinSize()
}

func (r *paletteRenderer) Objects() []fyne.CanvasObject {
    objects := []fyne.CanvasObject{r.gradient, r.rect}
	for _, cs := range r.colorStops {
		objects = append(objects, cs)
	}
	return objects
}

func (r *paletteRenderer) Refresh() {
	r.gradient.Refresh()
	r.rect.Refresh()
    r.updateColorStops()
	r.Layout(r.pal.Size())
}

func (r *paletteRenderer) generator(w, h int) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for x := range w {
		f := float64(x) / float64(w)
		for y := range h {
			img.Set(x, y, r.pal.ColorSource.Color(f))
		}
	}
	return img
}

func (r *paletteRenderer) updateColorStops() {
    switch pal := r.pal.ColorSource.(type) {
    case *ledgrid.GradientPalette:
        colorStops := pal.ColorStops()
        	if len(colorStops) != len(r.colorStops) {
		    r.colorStops = make([]*canvas.Circle, len(colorStops))
        }
        	for i, cs := range colorStops {
        		stop := canvas.NewCircle(cs.Col)
        		stop.Resize(colStopSize)
        		stop.StrokeColor = ledgrid.White
        		stop.StrokeWidth = 2 * borderWidth
        		r.colorStops[i] = stop
        	}
    case *ledgrid.UniformPalette:
        r.colorStops = make([]*canvas.Circle, 0)
    }
}
