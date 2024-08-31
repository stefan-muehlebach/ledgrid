package ledgrid

import (
	"image"
	"image/color"
    ledcolor "github.com/stefan-muehlebach/ledgrid/color"
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
	// Mit dieser Struktur (slice of slices) werden Pixel-Koordinaten in
	// Indizes uebersetzt.
	idxMap IndexMap
}

// Erstellt ein neues LED-Panel. size enthaelt die Dimension des (gesamten)
// Panels. Wird bei modConf nil uebergeben, so wird eine Default-Konfiguration
// der Module angenommen, welche bei DefaultModuleConfig naeher beschrieben
// wird.
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
	return ledcolor.LedColorModel
}

func (g *LedGrid) Bounds() image.Rectangle {
	return g.Rect
}

func (g *LedGrid) At(x, y int) color.Color {
	return g.LedColorAt(x, y)
}

func (g *LedGrid) Set(x, y int, c color.Color) {
	c1 := ledcolor.LedColorModel.Convert(c).(ledcolor.LedColor)
	g.SetLedColor(x, y, c1)
}

// Dient dem schnelleren Zugriff auf den Farbwert einer bestimmten Zelle, resp.
// einer bestimmten LED. Analog zu At(), retourniert den Farbwert jedoch als
// LedColor-Typ.
func (g *LedGrid) LedColorAt(x, y int) ledcolor.LedColor {
	if !(image.Point{x, y}.In(g.Rect)) {
		return ledcolor.LedColor{}
	}
	idx := g.PixOffset(x, y)
	src := g.Pix[idx : idx+3 : idx+3]
	return ledcolor.NewLedColorRGB(src[0], src[1], src[2])
}

// Analoge Methode zu Set(), jedoch ohne zeitaufwaendige Konvertierung.
func (g *LedGrid) SetLedColor(x, y int, c ledcolor.LedColor) {
	if !(image.Point{x, y}.In(g.Rect)) {
		return
	}
	idx := g.PixOffset(x, y)
	dst := g.Pix[idx : idx+3 : idx+3]
	dst[0] = c.R
	dst[1] = c.G
	dst[2] = c.B
}

// Damit wird der Offset eines bestimmten Farbwerts innerhalb des Slices
// Pix berechnet. Dabei wird beruecksichtigt, dass das die LED's im LedGrid
// schlangenfoermig angeordnet sind, und der Beginn der LED-Kette frei
// waehlbar in einer Ecke des Panels liegen kann.
func (g *LedGrid) PixOffset(x, y int) int {
	return g.idxMap[x][y]
}

// Mit Clear kann das ganze Grid geloescht, resp. alle LEDs auf die gleiche
// Farbe gebracht werden.
func (g *LedGrid) Clear(c ledcolor.LedColor) {
	for idx := 0; idx < len(g.Pix); idx += 3 {
		dst := g.Pix[idx : idx+3 : idx+3]
		dst[0] = c.R
		dst[1] = c.G
		dst[2] = c.B
	}
}
