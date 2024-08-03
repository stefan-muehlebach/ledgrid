package ledgrid

import (
	"image"
	"image/color"
	"time"
)

const (
	defFramesPerSec = 30
)

var (
	framesPerSecond int           = defFramesPerSec
	frameRefresh    time.Duration = time.Second / time.Duration(framesPerSecond)
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
	// Verkabelung!
	Pix []uint8

	idxMap IndexMap
}

// Erstellt ein neues LED-Panel. size enthaelt die Dimension des (gesamten)
// Panels.
func NewLedGrid(size image.Point, modConf ModuleConfig) *LedGrid {
	g := &LedGrid{}
	g.Rect = image.Rectangle{Max: size}
	g.Pix = make([]uint8, 3*g.Rect.Dx()*g.Rect.Dy())

    // Autom. Formatwahl
    if modConf == nil {
        modConf = DefaultModuleConfig(g.Rect.Size())
    }

	g.idxMap = modConf.IndexMap()
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
	return LedColor{slc[0], slc[1], slc[2], 0xff}
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
// schlangenfoermig angeordnet sind, und der Beginn der LED-Kette frei
// waehlbar in einer Ecke des Panels liegen kann.
func (g *LedGrid) PixOffset(x, y int) int {
	return g.idxMap[x][y]
}

// Hier kommen nun die fuer das LedGrid spezifischen Funktionen.
func (g *LedGrid) Clear(c LedColor) {
	for idx := 0; idx < len(g.Pix); idx += 3 {
		slc := g.Pix[idx : idx+3 : idx+3]
		slc[0] = c.R
		slc[1] = c.G
		slc[2] = c.B
	}
}

// func (g *LedGrid) MarkDefect(pos image.Point) {
//     g.idxMap.MarkDefect(pos)
// }
