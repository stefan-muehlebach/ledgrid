package ledgrid

import (
	"log"
	"slices"
	"time"
)

type colorStop struct {
	pos float64
	col LedColor
}

// Mit Paletten lassen sich anspruchsvolle Farbverlaeufe realisieren. Jeder
// Palette liegt eine Liste von Farben (die sog. Stuetzstellen) und ihre
// jeweilige Position auf dem Intervall [0, 1] zugrunde.
type GradientPalette struct {
	stops []colorStop
	// Mit dieser Funktion wird die Interpolation zwischen den gesetzten
	// Farbwerten realisiert.
	Func InterpolFuncType
}

// Erzeugt eine neue Palette und verwendet die Farben in cl als Stuetzwerte.
// In diesem Fall werden die Farben in cl gleichmaessig (aequidistant) auf
// dem Intervall [0,1] verteilt.
func NewGradientPalette(cycle bool, cl ...LedColor) *GradientPalette {
	p := &GradientPalette{}
	if cl == nil {
		cl = []LedColor{Black, White}
	}
    if cycle {
        cl = append(cl, cl[0])
    }
	p.SetColorStops(cl)
	p.Func = PolynomInterpol
	return p
}

// Setzt die Farbe c als neuen Stuetzwert bei Position t. Existiert bereits
// eine Farbe mit dieser Position, wird sie ueberschrieben.
func (p *GradientPalette) SetColorStop(t float64, c LedColor) {
	var i int
	var stop colorStop

	if t < 0.0 || t > 1.0 {
		log.Fatalf("Parameter t must be in [0, 1] instead of %f", t)
	}
	for i, stop = range p.stops {
		if stop.pos == t {
			p.stops[i].col = c
			return
		}
		if stop.pos > t {
			break
		}
	}
	p.stops = slices.Insert(p.stops, i, colorStop{t, c})
}

// Farbwerte in cl werden als Stuetzstellen der Palett verwendet. Die
// Stuetzstellen sind gleichmaessig ueber das Intervall [0,1] verteilt.
func (p *GradientPalette) SetColorStops(cl []LedColor) {
	if len(cl) < 2 {
		log.Fatalf("At least two colors are required!")
	}
	posStep := 1.0 / (float64(len(cl)-1))
	p.stops = make([]colorStop, len(cl))
	for i, c := range cl {
		p.stops[i] = colorStop{float64(i) * posStep, c}
    }
}

// Hier nun spielt die Musik: aufgrund des Wertes t (muss im Intervall [0,1]
// liegen) wird eine neue Farbe interpoliert.
func (p *GradientPalette) Color(t float64) (c LedColor) {
	var i int
	var stop colorStop

	if t < 0.0 || t > 1.0 {
		log.Fatalf("Color: parameter t must be in [0,1] instead of %f", t)
	}
	for i, stop = range p.stops[1:] {
		if stop.pos > t {
			break
		}
	}
	t = (t - p.stops[i].pos) / (p.stops[i+1].pos - p.stops[i].pos)
	c = p.stops[i].col.Interpolate(p.stops[i+1].col, p.Func(0, 1, t)).(LedColor)
	return c
}

// Palette mit 256 einzelnen dedizierten Farbwerten - kein Fading oder
// sonstige Uebergaenge.
type SlicePalette struct {
	Colors []LedColor
}

func NewSlicePalette(cl ...LedColor) *SlicePalette {
	p := &SlicePalette{}
	p.Colors = make([]LedColor, 256)
	for i, c := range cl {
		p.Colors[i] = c
	}
	return p
}

func (p *SlicePalette) Color(v float64) LedColor {
	return p.Colors[int(v)]
}

func (p *SlicePalette) SetColor(idx int, c LedColor) {
	p.Colors[idx] = c
}

// Mit diesem Typ kann ein fliessender Uebergang von einer Palette zu einer
// anderen realisiert werden.
type PaletteFader struct {
	AnimatableEmbed
	Pals                 [2]Colorable
	FadeTime, RemainTime time.Duration
}

func NewPaletteFader(pal Colorable) *PaletteFader {
	p := &PaletteFader{}
	p.AnimatableEmbed.Init()
	p.Pals[0] = pal
	p.FadeTime = 0
	p.RemainTime = 0
	return p
}

func (p *PaletteFader) StartFade(nextPal Colorable, fadeTime time.Duration) bool {
	if p.RemainTime > 0 {
		return false
	}
	p.Pals[0], p.Pals[1] = nextPal, p.Pals[0]
	if fadeTime > 0 {
		p.FadeTime = fadeTime
		p.RemainTime = fadeTime
	}
	return true
}

func (p *PaletteFader) Update(dt time.Duration) bool {
	dt = p.AnimatableEmbed.Update(dt)
	if p.RemainTime > 0 && dt > 0 {
		p.RemainTime -= dt
		if p.RemainTime < 0 {
			p.RemainTime = 0
		}
	}
	return true
}

func (p *PaletteFader) Color(v float64) LedColor {
	c1 := p.Pals[0].Color(v)
	if p.RemainTime > 0 {
		t := p.RemainTime.Seconds() / p.FadeTime.Seconds()
		c2 := p.Pals[1].Color(v)
		c1 = c1.Interpolate(c2, t).(LedColor)
	}
	return c1
}
