package ledgrid

import "image/color"

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
func (c LedColor) Interpolate(d LedColor, t float64) (LedColor) {
    t = max(min(t, 1), 0)
    r := (1-t)*float64(c.R) + t*float64(d.R)
    g := (1-t)*float64(c.G) + t*float64(d.G)
    b := (1-t)*float64(c.B) + t*float64(d.B)
    return LedColor{uint8(r), uint8(g), uint8(b)}
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
