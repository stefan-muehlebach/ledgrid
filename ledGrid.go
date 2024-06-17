package ledgrid

import (
	"image"
	"image/color"
	"time"
)

// Fuer das gesamte Package gueltige Variablen, resp. Konstanten.
var (
	OutsideColor = BlackColor
)

const (
	defFramesPerSec = 30
)

type CableConfig struct {
	start, dir image.Point
}

var (
	UpperLeft2Right = CableConfig{upperLeft, right}
	UpperLeft2Down  = CableConfig{upperLeft, down}
	UpperRight2Left = CableConfig{upperRight, left}
	UpperRight2Down = CableConfig{upperRight, down}
	LowerLeft2Right = CableConfig{lowerLeft, right}
	LowerLeft2Up    = CableConfig{lowerLeft, up}
	LowerRight2Left = CableConfig{lowerRight, left}
	LowerRight2Up   = CableConfig{lowerRight, up}

	upperLeft  = image.Point{0, 0}
	upperRight = image.Point{1, 0}
	lowerLeft  = image.Point{0, 1}
	lowerRight = image.Point{1, 1}
	right      = image.Point{1, 0}
	left       = image.Point{-1, 0}
	down       = image.Point{0, 1}
	up         = image.Point{0, -1}
)

var (
	framesPerSecond int
	frameRefresh    time.Duration
	frameRefreshMs  int64
	frameRefreshSec float64
)

func init() {
	framesPerSecond = defFramesPerSec
	frameRefresh = time.Second / time.Duration(framesPerSecond)
	frameRefreshMs = frameRefresh.Microseconds()
	frameRefreshSec = frameRefresh.Seconds()
}

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

	cabConf CableConfig
}

// Erstellt ein neues LED-Panel. r enthaelt die Dimension des (gesamten)
// Panels, und mit cabConf wird die Verkabelungskonfiguration bezeichnet.
func NewLedGrid(r image.Rectangle, cabConf CableConfig) *LedGrid {
	g := &LedGrid{}
	g.Rect = r
	g.Pix = make([]uint8, 3*r.Dx()*r.Dy())
	g.cabConf = cabConf
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
	return LedColor{slc[0], slc[1], slc[2], 0xFF}
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

// func (g *LedGrid) MixLedColor(x, y int, c LedColor, mixType ColorMixType) {
// 	if !(image.Point{x, y}.In(g.Rect)) {
// 		return
// 	}
// 	bgCol := g.LedColorAt(x, y)
// 	g.SetLedColor(x, y, c.Mix(bgCol, mixType))
// }

// Damit wird der Offset eines bestimmten Farbwerts innerhalb des Slices
// Pix berechnet. Dabei wird beruecksichtigt, dass das die LED's im LedGrid
// schlangenfoermig angeordnet sind, und der Beginn der LED-Kette frei
// waehlbar in einer Ecke des Panels liegen kann.
func (g *LedGrid) PixOffset(x, y int) int {
	var idx int

	if g.cabConf.start.X == 1 {
		x = (g.Rect.Dx() - 1) - x
	}
	if g.cabConf.start.Y == 1 {
		y = (g.Rect.Dy() - 1) - y
	}

	if g.cabConf.dir.X != 0 {
		idx = y * g.Rect.Dx()
		if y%2 == 0 {
			idx += x
		} else {
			idx += (g.Rect.Dx() - 1) - x
		}
	}
	if g.cabConf.dir.Y != 0 {
		idx = x * g.Rect.Dy()
		if x%2 == 0 {
			idx += y
		} else {
			idx += (g.Rect.Dy() - 1) - y
		}
	}
	return 3 * idx
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
