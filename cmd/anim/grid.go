package main

import (
	"image"
	"image/draw"
	"io"
	"log"
	"os"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/colornames"
)

type Grid struct {
	ObjList []GridObject
	width, height int
	pixCtrl       ledgrid.PixelClient
	ledGrid       *ledgrid.LedGrid
	objMutex *sync.RWMutex
	logFile io.Writer
	paintWatch, sendWatch *ledgrid.Stopwatch
}

func NewGrid(pixCtrl ledgrid.PixelClient, ledGrid *ledgrid.LedGrid) *Grid {
	var err error

	c := &Grid{}
	c.pixCtrl = pixCtrl
	c.ledGrid = ledGrid
	c.width = ledGrid.Rect.Dx()
	c.height = ledGrid.Rect.Dy()
	c.ObjList = make([]GridObject, 0)
	c.objMutex = &sync.RWMutex{}
	if doLog {
		c.logFile, err = os.Create("grid.log")
		if err != nil {
			log.Fatalf("Couldn't create logfile: %v", err)
		}
	}
	c.paintWatch = ledgrid.NewStopwatch()
	c.sendWatch = ledgrid.NewStopwatch()
	return c
}

func (c *Grid) Close() {
	c.DelAll()
}

// Fuegt der Zeichenflaeche weitere Objekte hinzu. Der Zufgriff auf den
// entsprechenden Slice wird nicht synchronisiert.
func (c *Grid) Add(objs ...DrawableObject) {
	c.objMutex.Lock()
	for _, tmp := range objs {
		obj := tmp.(GridObject)
		c.ObjList = append(c.ObjList, obj)
	}
	c.objMutex.Unlock()
}

// Loescht alle Objekte von der Zeichenflaeche.
func (c *Grid) DelAll() {
	c.objMutex.Lock()
	c.ObjList = c.ObjList[:0]
	c.objMutex.Unlock()
}

func (c *Grid) Refresh() {
	c.paintWatch.Start()
	c.ledGrid.Clear(colornames.Black)
	c.objMutex.RLock()
	for _, obj := range c.ObjList {
		obj.Draw(c)
	}
	c.objMutex.RUnlock()
	c.paintWatch.Stop()

	c.sendWatch.Start()
	c.pixCtrl.Draw(c.ledGrid)
	c.sendWatch.Stop()
}

// Alle Objekte, die durch den Controller auf dem LED-Grid dargestellt werden
// sollen, muessen im Minimum die Methode Draw implementieren, durch welche
// sie auf einem gg-Kontext gezeichnet werden.
type GridObject interface {
	Draw(d DrawingArea)
	// Draw(lg *ledgrid.LedGrid)
}

// Grid-Objekt fuer Images
type GridImage struct {
	Pos image.Point
	Img *image.RGBA
}

func NewGridImage(pos image.Point, size image.Point) *GridImage {
	i := &GridImage{Pos: pos}
	i.Img = image.NewRGBA(image.Rectangle{Max: size})
	return i
}

func (i *GridImage) Draw(d DrawingArea) {
    g := d.(*Grid)
	draw.Draw(g.ledGrid, i.Img.Bounds().Add(i.Pos), i.Img, image.Point{0, 0}, draw.Over)
}

// Will man ein einzelnes Pixel zeichnen, so eignet sich dieser Typ. Er wird
// ueber die Zeichenfunktion DrawPoint im gg-Kontext realisiert und hat einen
// Radius von 0.5*sqrt(2).
type GridPixel struct {
	Pos   image.Point
	Color ledgrid.LedColor
}

func NewGridPixel(pos image.Point, col ledgrid.LedColor) *GridPixel {
	p := &GridPixel{Pos: pos, Color: col}
	return p
}

func (p *GridPixel) Draw(d DrawingArea) {
    g := d.(*Grid)
	g.ledGrid.SetLedColor(p.Pos.X, p.Pos.Y, p.Color.Mix(g.ledGrid.LedColorAt(p.Pos.X, p.Pos.Y), ledgrid.Blend))
}

// Fuer das direkte Zeichnen von Text auf dem LED-Grid, existieren einige
// 'fixed size' Bitmap-Schriften, die ohne Rastern und Rendern sehr schnell
// dargestellt werden koennen.
type GridText struct {
	Pos    fixed.Point26_6
	Color  ledgrid.LedColor
	text   string
	drawer *font.Drawer
	rect   fixed.Rectangle26_6
	dp     fixed.Point26_6
}

func NewGridText(pos fixed.Point26_6, col ledgrid.LedColor, text string) *GridText {
	t := &GridText{Pos: pos, Color: col}
	t.drawer = &font.Drawer{
		Face: ledgrid.Face3x5,
	}
	t.SetText(text)
	return t
}

func (t *GridText) Text() string {
	return t.text
}

func (t *GridText) SetText(text string) {
	t.text = text
	t.rect, _ = t.drawer.BoundString(text)
	t.dp = t.rect.Min.Add(t.rect.Max).Div(fixed.I(2))
}

func (t *GridText) Draw(d DrawingArea) {
    g := d.(*Grid)
	t.drawer.Dst = g.ledGrid
	t.drawer.Src = image.NewUniform(t.Color)
	t.drawer.Dot = t.Pos.Sub(t.dp)
	t.drawer.DrawString(t.text)
}
