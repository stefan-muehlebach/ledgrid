package ledgrid

import (
	"fmt"
	gocolor "image/color"
	"log"
	"strconv"

	"github.com/stefan-muehlebach/gg/color"
)

var (
	Black       = LedColor{0x00, 0x00, 0x00, 0xFF}
	White       = LedColor{0xFF, 0xFF, 0xFF, 0xFF}
	Red         = LedColor{0xFF, 0x00, 0x00, 0xFF}
	Green       = LedColor{0x00, 0xFF, 0x00, 0xFF}
	Blue        = LedColor{0x00, 0x00, 0xFF, 0xFF}
	Transparent = LedColor{0x00, 0x00, 0x00, 0x00}
)

// // Damit verschiedene Interpolationsfunktionen verwendet werden koennen,
// // ist das Profil als Typ definiert. Jede Interp.funktion realisiert eine
// // (wie auch immer geartete) Interpolation zwischen den Werten a und b in
// // Abhaengigkeit von t, wobei t in [0, 1] ist. Es muss gelten:
// //
// //   - f(a, b, 0.0) = a
// //   - f(a, b, 1.0) = b
// //
// // TO DO: aktuell wird keine Fehlererkennung gemacht. Wenn beispielsweise
// // t nicht in [0, 1] liegt, erzeugen die Funktionen ggf. unsinnige Werte.
// type InterpolFuncType func(a, b, t float64) float64

// var (
// 	ColorInterpol = PolynomInterpol
// )

// // Realisiert die klassische lineare Interpolation zwischen den Werten a und b.
// func LinearInterpol(a, b, t float64) float64 {
// 	return (1-t)*a + t*b
// }

// // Realisiert eine kubische Interpolation.
// func PolynomInterpol(a, b, t float64) float64 {
// 	t = 3*t*t - 2*t*t*t
// 	return LinearInterpol(a, b, t)
// }

// // Woher diese Interpolation stammt, kann ich auch nicht mehr mit Bestimmtheit
// // sagen. Sicher ist, dass sie nicht mehr funktioniert, wenn a oder b negativ
// // sind: dann gelten die oben genannte Bedingungen nicht mehr.
// func SqrtInterpol(a, b, t float64) float64 {
// 	return math.Sqrt((1-t)*a*a + t*b*b)
// }

func interp(a, b, t float64) float64 {
	return (1-t)*a + t*b
}

// Dieser Typ wird fuer die Farbwerte verwendet, welche via SPI zu den LED's
// gesendet werden. Die Daten sind _nicht_ gamma-korrigiert, dies wird erst
// auf dem Panel-Empfaenger gemacht (pixelcontroller-slave). LedColor
// implementiert das color.Color Interface.
type LedColor struct {
	R, G, B, A uint8
}

func NewLedColor(hex uint32) LedColor {
	r := (hex & 0xff0000) >> 16
	g := (hex & 0x00ff00) >> 8
	b := (hex & 0x0000ff)
	return LedColor{uint8(r), uint8(g), uint8(b), 0xff}
}

func NewLedColorAlpha(hex uint64) LedColor {
	r := (hex & 0xff000000) >> 24
	g := (hex & 0x00ff0000) >> 16
	b := (hex & 0x0000ff00) >> 8
	a := (hex & 0x000000ff)
	return LedColor{uint8(r), uint8(g), uint8(b), uint8(a)}
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

// Berechnet eine RGB-Farbe, welche 'zwischen' den Farben c1 und c2 liegt,
// so dass bei t=0 der Farbwert c1 und bei t=1 der Farbwert c2 retourniert
// wird. t wird vorgaengig auf das Interval [0,1] eingeschraenkt.
func (c1 LedColor) Interpolate(c2 color.Color, t float64) color.Color {
	t = max(min(t, 1.0), 0.0)
	if t == 0.0 {
		return c1
	}
	if t == 1.0 {
		return c2
	}
	if c3, ok := c2.(LedColor); ok {
		r := interp(float64(c1.R), float64(c3.R), t)
		g := interp(float64(c1.G), float64(c3.G), t)
		b := interp(float64(c1.B), float64(c3.B), t)
		a := interp(float64(c1.A), float64(c3.A), t)
		return LedColor{uint8(r), uint8(g), uint8(b), uint8(a)}
	} else {
		return LedColor{}
	}
}

// Mit folgenden Konstanten kann das Verfahren bestimmt werden, welches beim
// Mischen von Farben verwendet werden soll (siehe auch Methode Mix).
type ColorMixType int

const (
	// Ersetzt die Hintergundfarbe durch die Vordergrundfarbe.
	Replace ColorMixType = iota
	// Ueberblendet die Hintergrundfarbe mit der Vordergrundfarbe anhand
	// des Alpha-Wertes.
	Blend
	// Bestimmt die neue Farbe durch das Maximum von RGB zwischen Hinter- und
	// Vordergrundfarbe.
	Max
	// Analog zu Max, nimmt jedoch den Mittelwert von jeweils R, G und B.
	Average
	// Analog zu Max, nimmt jedoch das Minimum von jeweils R, G und B.
	Min
)

// Mischt die Farben c (Vordergrundfarbe) und bg (Hintergrundfarbe) nach einem
// Verfahren, welches in typ spezifiziert ist. Aktuell stehen 'Blend' (Ueber-
// blendung von bg durch c anhand des Alpha-Wertes von c) und 'Add' (nimm
// jeweils das Maximum pro Farbwert von c und bg) zur Verfuegung.
func (c LedColor) Mix(bg LedColor, mix ColorMixType) LedColor {
	switch mix {
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
	case Max:
		r := max(c.R, bg.R)
		g := max(c.G, bg.G)
		b := max(c.B, bg.B)
		a := max(c.A, bg.A)
		return LedColor{r, g, b, a}
	case Average:
		r := c.R/2 + bg.R/2
		g := c.G/2 + bg.G/2
		b := c.B/2 + bg.B/2
		a := c.A/2 + bg.A/2
		return LedColor{r, g, b, a}
	case Min:
		r := min(c.R, bg.R)
		g := min(c.G, bg.G)
		b := min(c.B, bg.B)
		a := min(c.A, bg.A)
		return LedColor{r, g, b, a}
	default:
		log.Fatalf("Unknown mixing function: '%d'", mix)
	}
	return LedColor{}
}

func (c LedColor) Alpha(a float64) color.Color {
	return LedColor{c.R, c.G, c.B, uint8(255.0 * a)}
}

func (c LedColor) Bright(t float64) color.Color {

	return c
}

func (c LedColor) Dark(t float64) color.Color {
	return c
}

func (c LedColor) String() string {
	return fmt.Sprintf("{0x%02X, 0x%02X, 0x%02X, 0x%02X}", c.R, c.G, c.B, c.A)
}

func (c *LedColor) UnmarshalText(text []byte) error {
	hexStr := string(text)
	hexVal, err := strconv.ParseUint(hexStr, 16, 32)
	if err != nil {
		log.Fatal(err)
	}
	c.R = uint8((hexVal & 0xff0000) >> 16)
	c.G = uint8((hexVal & 0x00ff00) >> 8)
	c.B = uint8((hexVal & 0x0000ff))
	c.A = 0xff
	return nil
}

// Das zum Typ LedColor zugehoerende ColorModel.
var (
	LedColorModel gocolor.Model = gocolor.ModelFunc(ledColorModel)
)

// Wandelt einen beliebigen Farbwert c in einen LedColor-Typ um.
func ledColorModel(c gocolor.Color) gocolor.Color {
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
