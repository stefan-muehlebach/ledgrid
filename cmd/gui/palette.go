package main

import (
	"image"
	"log"

	"fyne.io/fyne/v2/driver/desktop"

	"github.com/stefan-muehlebach/ledgrid"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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
	log.Printf("tapped at %+v\n", evt)

	log.Printf("  size: %v\n", p.Size())
	log.Printf("  ratio within: %.4f\n", p.getRatio(evt))
	// colorPicker := dialog.NewColorPicker("Title", "Message", func(c color.Color){}, Win)
	// colorPicker.Advanced = true
	// colorPicker.Show()
}

func (p *Palette) getRatio(evt *fyne.PointEvent) float64 {
    margin := float32(colStopDia/2)

    x := evt.Position.X

    switch p.Orientation {
    case widget.Horizontal:
        if x > p.Size().Width-margin {
            return 1.0
        } else if x < margin {
            return 0.0
        } else {
            return float64(x-margin) / float64(p.Size().Width-2*margin)
        }
    }
    return 0.0
}


func (p *Palette) MouseIn(evt *desktop.MouseEvent) {
	// log.Printf("Mouse in!")
}

func (p *Palette) MouseMoved(evt *desktop.MouseEvent) {
	// log.Printf("Mouse moved: %v", evt)
}

func (p *Palette) MouseOut() {
	// log.Printf("Mouse out!")
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
    pos := fyne.NewPos(colStopDia/2, 0)
    s = s.SubtractWidthHeight(colStopDia, 0)
    r.gradient.Move(pos)
	r.gradient.Resize(s)

    r.rect.Move(pos)
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
