package ledgrid

import (
	"image"
	"image/color"
)

// Entspricht dem Bild, welches auf einem LED-Panel angezeigt werden kann.
// Implementiert die Interfaces image.Image und draw.Image, also die Methoden
// ColorModel, Bounds, At und Set.
type LedGrid struct {
    // Groesse des LedGrids. Falls dieses LedGrid Teil eines groesseren
    // Panels sein sollte, dann muss Rect.Min nicht unbedingt {0, 0} sein.
	Rect image.Rectangle
    // Enthaelt die Farbwerte red, green, blue (RGB) fuer jede LED, welche
    // das LedGrid ausmachen. Die Reihenfolge entspricht dabei der
    // Verkabelung, d.h. sie beginnt links oben mit der LED Nr. 0,
    // geht dann nach rechts und auf der zweiten Zeile wieder nach links und
    // so schlangenfoermig weiter.
	Pix []uint8
}

func NewLedGrid(r image.Rectangle) *LedGrid {
	g := &LedGrid{}
	g.Rect = r
	g.Pix = make([]uint8, 3*r.Dx()*r.Dy())
	return g
}

func (g *LedGrid) ColorModel() color.Model {
	return LedColorModel
}

func (g *LedGrid) Bounds() image.Rectangle {
	return g.Rect
}

func (g *LedGrid) At(x, y int) color.Color {
	return g.LedColorAt(x, y)
}

func (g *LedGrid) Set(x, y int, c color.Color) {
	c1 := LedColorModel.Convert(c).(LedColor)
	g.SetLedColor(x, y, c1)
}

// Dient dem schnelleren Zugriff auf den Farbwert einer bestimmten Zelle, resp.
// einer bestimmten LED. Analog zu At(), retourniert den Farbwert jedoch als
// LedColor-Typ.
func (g *LedGrid) LedColorAt(x, y int) LedColor {
	if !(image.Point{x, y}.In(g.Rect)) {
		return LedColor{}
	}
	idx := g.PixOffset(x, y)
	slc := g.Pix[idx : idx+3 : idx+3]
	return LedColor{slc[0], slc[1], slc[2]}
}

// Analoge Methode zu Set(), jedoch ohne zeitaufwaendige Konvertierung.
func (g *LedGrid) SetLedColor(x, y int, c LedColor) {
	if !(image.Point{x, y}.In(g.Rect)) {
		return
	}
	idx := g.PixOffset(x, y)
	slc := g.Pix[idx : idx+3 : idx+3]
	slc[0] = c.R
	slc[1] = c.G
	slc[2] = c.B
}

// Damit wird der Offset eines bestimmten Farbwerts innerhalb des Slices
// Pix berechnet. Dabei wird beruecksichtigt, dass das die LED's im LedGrid
// schlangenfoermig angeordnet sind, wobei die Reihe mit der LED links oben
// beginnt.
func (g *LedGrid) PixOffset(x, y int) int {
	var idx int

	idx = y * g.Rect.Dx()
	if y%2 == 0 {
		idx += x
	} else {
		idx += (g.Rect.Dx() - x - 1)
	}
	return 3 * idx
}
