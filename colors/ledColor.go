// Das Package colors enthaelt einerseits einen eigenen Farbtyp fuer die
// darstellbaren Farben eines eines LEDs (resp. NeoPixel). Im wesentlichen
// ist dies eine Kopie von color.NRGBA mit zusaetzlichen Methoden um bspw.
// zwischen zwei Farben zu interpolieren oder um Farben aufzuhellen, resp.
// abzudunkeln.

package colors

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"strconv"
)

var (
    // Transparent ist ein komplett durchsichtiges Schwarz.
	Transparent = LedColor{0x00, 0x00, 0x00, 0x00}
)

// Dieser Typ wird fuer die Farbwerte verwendet, welche via SPI zu den LED's
// gesendet werden. Die Daten sind _nicht_ gamma-korrigiert (dies wird erst
// auf dem Panel-Empfaenger gemacht) und entsprechen dem Typ color.NRGBA
// von Go. LedColor implementiert das color.Color Interface.
type LedColor struct {
    R, G, B, A uint8
}

// RGBA ist Teil des color.Color Interfaces und retourniert die Farbwerte
// als Alpha-korrigierte uint16-Werte.
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

// Hilfsvariable fuer die Methode HSL.
var (
	cm = []float64{
		1, -0.5, -0.5,
		0, math.Sqrt(3) / 2, -math.Sqrt(3) / 2,
		math.Sqrt2 / 2, math.Sqrt2 / 2, math.Sqrt2 / 2,
	}
)

// Mit HSL() koennen die Werte fuer Hue [0, 360], Saturation [0, 1] und
// Lightness [0, 1] aus dem gleichnamigen Farbmodell ermittelt werden.
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
    u := 1.0 - t
	r := u*float64(c1.R) + t*float64(c2.R)
	g := u*float64(c1.G) + t*float64(c2.G)
	b := u*float64(c1.B) + t*float64(c2.B)
	a := u*float64(c1.A) + t*float64(c2.A)
	return LedColor{uint8(r), uint8(g), uint8(b), uint8(a)}
}

// Auch wenn man es fast nicht glauben mag: die Variante mit den float64 Werten
// ist im Benchmark schneller als die unten auskommentierte Variante mit uint32
// oder uint16 Werten.
// func (c1 LedColor) Interpolate(c2 LedColor, t float64) LedColor {
//     t = max(min(t, 1.0), 0.0)
//
//     T := uint32(t * 255.0)
// 	if T == 0x00 {
// 		return c1
// 	}
// 	if T == 0xFF {
// 		return c2
// 	}
//
//     U := 0xFF - T
// 	r := U*uint32(c1.R)/0xFF + T*uint32(c2.R)/0xFF
// 	g := U*uint32(c1.G)/0xFF + T*uint32(c2.G)/0xFF
// 	b := U*uint32(c1.B)/0xFF + T*uint32(c2.B)/0xFF
// 	a := U*uint32(c1.A)/0xFF + T*uint32(c2.A)/0xFF
// 	return LedColor{uint8(r), uint8(g), uint8(b), uint8(a)}
// }

// Retourniert eine neue Farbe, basierend auf c, jedoch mit dem hier
// angegebenen Alpha-Wert (als Fliesskommazahl in [0, 1]).
func (c LedColor) Alpha(a float64) LedColor {
	a = max(min(a, 1.0), 0.0)
	return LedColor{c.R, c.G, c.B, uint8(255.0*a)}
}

// Retourniert eine neue Farbe, welche eine Interpolation zwischen c und Weiss
// ist. t ist ein Wert in [0, 1] und bestimmt die Position der Interpolation.
// t=0 retourniert c, t=1 retourniert Weiss.
func (c LedColor) Bright(t float64) LedColor {
	t = max(min(t, 1.0), 0.0)
	return c.Interpolate(White, t)
}

// Retourniert eine neue Farbe, welche eine Interpolation zwischen c und
// Schwarz ist. t ist ein Wert in [0, 1] und bestimmt die Position der
// Interpolation. t=0 retourniert c, t=1 retourniert Schwarz.
func (c LedColor) Dark(t float64) LedColor {
	t = max(min(t, 1.0), 0.0)
	return c.Interpolate(Black, t)
}

// Erzeugt eine druckbare Variante der Farbe. Im Wesentlichen werden die Werte
// fuer Rot, Gruen, Blau und Alpha als Hex-Zahlen ausgegeben.
func (c LedColor) String() string {
	return fmt.Sprintf("{0x%02X, 0x%02X, 0x%02X, 0x%02X}", c.R, c.G, c.B, c.A)
}

// Damit koennen Farbwerte im Format 0xRRGGBB eingelesen werden, wie sie bspw.
// in JSON-Files verwendet werden.
func (c *LedColor) UnmarshalText(text []byte) error {
	hexStr := string(text[2:])
	hexVal, err := strconv.ParseUint(hexStr, 16, 32)
	if err != nil {
		log.Fatal(err)
	}
	c.R = uint8((hexVal & 0xFF0000) >> 16)
	c.G = uint8((hexVal & 0x00FF00) >> 8)
	c.B = uint8((hexVal & 0x0000FF))
	c.A = 0xFF
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
// Verfahren, welches in mix spezifiziert ist. Siehe auch ColorMixType.
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
		return LedColor{uint8(r), uint8(g), uint8(b), uint8(255.0*a)}
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
		return LedColor{uint8(r>>8), uint8(g>>8), uint8(b>>8), 0xff}
	}
	if a == 0x0000 {
		return LedColor{}
	}
	r = (r * 0xFFFF) / a
	g = (g * 0xFFFF) / a
	b = (b * 0xFFFF) / a
	return LedColor{uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)}
}
