package color

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"strconv"
)

var (
	Transparent = NewLedColorRGBA(0x00, 0x00, 0x00, 0x00)
)

var (
	cm = []float64{
		1, -0.5, -0.5,
		0, math.Sqrt(3) / 2, -math.Sqrt(3) / 2,
		math.Sqrt2 / 2, math.Sqrt2 / 2, math.Sqrt2 / 2,
	}
	// cm = []float64{
	// 	math.Sqrt(6) / 3, -math.Sqrt(6) / 6, -math.Sqrt(6) / 6,
	// 	0, math.Sqrt2 / 2, -math.Sqrt2 / 2,
	// 	math.Sqrt(3) / 3, math.Sqrt(3) / 3, math.Sqrt(3) / 3,
	// }
)

// Dieser Typ wird fuer die Farbwerte verwendet, welche via SPI zu den LED's
// gesendet werden. Die Daten sind _nicht_ gamma-korrigiert, dies wird erst
// auf dem Panel-Empfaenger gemacht (pixelcontroller-slave). LedColor
// implementiert das color.Color Interface.
type LedColor struct {
	R, G, B, A uint8
}

func NewLedColorRGB(r, g, b uint8) LedColor {
	return LedColor{R: r, G: g, B: b, A: 0xff}
}

func NewLedColorRGBA(r, g, b, a uint8) LedColor {
	return LedColor{R: r, G: g, B: b, A: a}
}

func NewLedColorHex(hex uint32) LedColor {
	r := (hex & 0xff0000) >> 16
	g := (hex & 0x00ff00) >> 8
	b := (hex & 0x0000ff)
	return NewLedColorRGB(uint8(r), uint8(g), uint8(b))
}

func NewLedColorHexA(hex uint64) LedColor {
	r := (hex & 0xff000000) >> 24
	g := (hex & 0x00ff0000) >> 16
	b := (hex & 0x0000ff00) >> 8
	a := (hex & 0x000000ff)
	return NewLedColorRGBA(uint8(r), uint8(g), uint8(b), uint8(a))
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

func (c LedColor) HSL() (h, s, l float64) {
	v1 := []float64{float64(c.R) / 255.0, float64(c.G) / 255.0, float64(c.B) / 255.0}
	v2 := []float64{
		cm[0]*v1[0] + cm[1]*v1[1] + cm[2]*v1[2],
		cm[3]*v1[0] + cm[4]*v1[1] + cm[5]*v1[2],
		cm[6]*v1[0] + cm[7]*v1[1] + cm[8]*v1[2],
	}
	h = 180.0 * math.Atan2(v2[1], v2[0]) / math.Pi
	s = math.Hypot(v2[0], v2[1])
	l = v2[2] * math.Sqrt2 / 3
	return
}

// Berechnet eine RGB-Farbe, welche 'zwischen' den Farben c1 und c2 liegt,
// so dass bei t=0 der Farbwert c1 und bei t=1 der Farbwert c2 retourniert
// wird. t wird vorgaengig auf das Interval [0,1] eingeschraenkt.
func (c1 LedColor) Interpolate(c2 LedColor, t float64) LedColor {
	t = max(min(t, 1.0), 0.0)
	if t == 0.0 {
		return c1
	}
	if t == 1.0 {
		return c2
	}
	r := (1.0-t)*float64(c1.R) + t*float64(c2.R)
	g := (1.0-t)*float64(c1.G) + t*float64(c2.G)
	b := (1.0-t)*float64(c1.B) + t*float64(c2.B)
	a := (1.0-t)*float64(c1.A) + t*float64(c2.A)
	return NewLedColorRGBA(uint8(r), uint8(g), uint8(b), uint8(a))
}

func (c LedColor) Alpha(a float64) LedColor {
	return NewLedColorRGBA(c.R, c.G, c.B, uint8(255.0*a))
}

func (c LedColor) Bright(t float64) LedColor {
	t = max(min(t, 1.0), 0.0)
	return c.Interpolate(White, t)
}

func (c LedColor) Dark(t float64) LedColor {
	t = max(min(t, 1.0), 0.0)
	return c.Interpolate(Black, t)
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
		return NewLedColorRGBA(uint8(r), uint8(g), uint8(b), uint8(255.0*a))
	case Max:
		r := max(c.R, bg.R)
		g := max(c.G, bg.G)
		b := max(c.B, bg.B)
		a := max(c.A, bg.A)
		return NewLedColorRGBA(r, g, b, a)
	case Average:
		r := c.R/2 + bg.R/2
		g := c.G/2 + bg.G/2
		b := c.B/2 + bg.B/2
		a := c.A/2 + bg.A/2
		return NewLedColorRGBA(r, g, b, a)
	case Min:
		r := min(c.R, bg.R)
		g := min(c.G, bg.G)
		b := min(c.B, bg.B)
		a := min(c.A, bg.A)
		return NewLedColorRGBA(r, g, b, a)
	default:
		log.Fatalf("Unknown mixing function: '%d'", mix)
	}
	return LedColor{}
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

	// return LedColor{color.NRGBAModel.Convert(c).(color.NRGBA)}

	r, g, b, a := c.RGBA()
	if a == 0xFFFF {
		return NewLedColorRGB(uint8(r>>8), uint8(g>>8), uint8(b>>8))
	}
	if a == 0 {
		return LedColor{}
	}
	r = (r * 0xFFFF) / a
	g = (g * 0xFFFF) / a
	b = (b * 0xFFFF) / a
	return NewLedColorRGBA(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8))
}
