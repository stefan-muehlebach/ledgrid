package ledgrid

import (
	"fmt"
	"math"
	"image/color"
)

var (
	Black = LedColor{0x00, 0x00, 0x00}
	White = LedColor{0xFF, 0xFF, 0xFF}
	Red   = LedColor{0xFF, 0x00, 0x00}
	Green = LedColor{0x00, 0xFF, 0x00}
	Blue  = LedColor{0x00, 0x00, 0xFF}
)

type InterpolFuncType func(a, b, t float64) float64

var (
    ColorInterpol = LinearInterpol
)

func LinearInterpol(a, b, t float64) float64 {
    return (1-t)*a + t*b
}

func PolynomInterpol(a, b, t float64) float64 {
    t = 3*t*t - 2*t*t*t
    return LinearInterpol(a, b, t)
}

func SqrtInterpol(a, b, t float64) float64 {
    return math.Sqrt((1-t)*a*a + t*b*b)
}

// Dieser Typ wird fuer die Farbwerte verwendet, welche via SPI zu den LED's
// gesendet werden. Die Daten sind _nicht_ gamma-korrigiert, dies wird erst
// auf dem Panel-Empfaenger gemacht (pixelcontroller-slave).
// LedColor implementiert das color.Color Interface.
type LedColor struct {
	R, G, B uint8
}

// RGBA ist Teil des color.Color Interfaces.
func (c LedColor) RGBA() (r, g, b, a uint32) {
	r, g, b = uint32(c.R), uint32(c.G), uint32(c.B)
	r |= r << 8
	g |= g << 8
	b |= b << 8
	a = 65535
	return
}

// Dient dem schnelleren Zugriff auf das Trippel der drei Farbwerte.
func (c LedColor) RGB() (r, g, b uint8) {
	return c.R, c.G, c.B
}

// Berechnet eine RGB-Farbe, welche 'zwischen' den Farben c und d liegt, so
// dass bei t=0 der Farbwert c und bei t=1 der Farbwert d retourniert wird.
// t wird vorgaengig auf das Interval [0,1] eingeschraenkt.
func (c LedColor) Interpolate(d LedColor, t float64) LedColor {
	t = max(min(t, 1.0), 0.0)
	if t == 0.0 {
		return c
	}
	if t == 1.0 {
		return d
	}
    r := ColorInterpol(float64(c.R), float64(d.R), t)
    g := ColorInterpol(float64(c.G), float64(d.G), t)
    b := ColorInterpol(float64(c.B), float64(d.B), t)
	return LedColor{uint8(r), uint8(g), uint8(b)}
}

// Mischt die Farben c und d so, dass jeweils der maximale Farbwert pro
// R, G, B von c und d beruecksichtigt wird.
func (c LedColor) Mix(d LedColor) LedColor {
	r := max(c.R, d.R)
	g := max(c.G, d.G)
	b := max(c.B, d.B)
	return LedColor{r, g, b}
}

func (c LedColor) String() (string) {
    return fmt.Sprintf("{0x%02X, 0x%02X, 0x%02X}", c.R, c.G, c.B)
}

// Das zum Typ LedColor zugehoerende ColorModel.
var (
	LedColorModel color.Model = color.ModelFunc(ledColorModel)
)

// Wandelt einen beliebigen Farbwert c in einen LedColor-Typ um.
func ledColorModel(c color.Color) color.Color {
	if _, ok := c.(LedColor); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	return LedColor{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
}
