package ledgrid

import (
	"log"
	"fmt"
	"image/color"
	"math"
)

var (
	Black = LedColor{0x00, 0x00, 0x00, 0xFF}
	White = LedColor{0xFF, 0xFF, 0xFF, 0xFF}
	Red   = LedColor{0xFF, 0x00, 0x00, 0xFF}
	Green = LedColor{0x00, 0xFF, 0x00, 0xFF}
	Blue  = LedColor{0x00, 0x00, 0xFF, 0xFF}
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
	return (1-t)*a + t*b
}

func SqrtInterpol(a, b, t float64) float64 {
	return math.Sqrt((1-t)*a*a + t*b*b)
}

type ColorMixType int

const (
    Replace ColorMixType = iota
	Blend
	Add
)

// Dieser Typ wird fuer die Farbwerte verwendet, welche via SPI zu den LED's
// gesendet werden. Die Daten sind _nicht_ gamma-korrigiert, dies wird erst
// auf dem Panel-Empfaenger gemacht (pixelcontroller-slave).
// LedColor implementiert das color.Color Interface.
type LedColor struct {
	R, G, B, A uint8
}

// RGBA ist Teil des color.Color Interfaces.
func (c LedColor) RGBA() (r, g, b, a uint32) {
	r, g, b, a = uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
	r |= r << 8
	r *= a
	r /= 0xFF
	g |= g << 8
	g *= a
	g /= 0xFF
	b |= b << 8
	b *= a
	b /= 0xFF
	a |= a << 8
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
	a := ColorInterpol(float64(c.A), float64(d.A), t)
	return LedColor{uint8(r), uint8(g), uint8(b), uint8(a)}
}

// Mischt die Farben c (Vordergrundfarbe) und d (Hintergrundfarbe) nach einem
// Verfahren, welches in typ spezifiziert ist. Aktuell stehen 'Blend' (Ueber-
// blendung von d durch c anhand des Alpha-Wertes von c) und 'Add' (nimm
// jeweils das Maximum pro Farbwert zwischen c und d) zur Verfuegung.
func (c LedColor) Mix(bg LedColor, typ ColorMixType) LedColor {
	switch typ {
    case Replace:
        return c
	case Blend:
		ca := float64(c.A) / 255.0
		da := float64(bg.A) / 255.0
		a := 1.0 - (1.0-ca)*(1.0-da)
		t1 := ca / a
		t2 := da * (1.0 - ca) / a
		r := float64(c.R)*t1 + float64(bg.R)*t2
		g := float64(c.G)*t1 + float64(bg.G)*t2
		b := float64(c.B)*t1 + float64(bg.B)*t2
		return LedColor{uint8(r), uint8(g), uint8(b), uint8(255.0 * a)}
	case Add:
		r := max(c.R, bg.R)
		g := max(c.G, bg.G)
		b := max(c.B, bg.B)
		a := max(c.A, bg.A)
		return LedColor{r, g, b, a}
    default:
        log.Fatalf("Unknown mixing function: '%d'", typ)
	}
    return LedColor{}
}

func (c LedColor) String() string {
	return fmt.Sprintf("{0x%02X, 0x%02X, 0x%02X, 0x%02X}", c.R, c.G, c.B, c.A)
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
	r, g, b, a := c.RGBA()
	if a == 0xFFFF {
		return LedColor{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 0xFF}
	}
	if a == 0 {
		return LedColor{0, 0, 0, 0}
	}
	r = (r * 0xFFFF) / a
	g = (g * 0xFFFF) / a
	b = (b * 0xFFFF) / a
	return LedColor{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}
