package ledgrid

import (
	"log"
	"slices"
	"time"

	"github.com/stefan-muehlebach/gg/color"
)

// Dieser (interne) Typ wird verwendet, um einen bestimmten Wert im Interval
// [0, 1] mit einer Farbe zu assoziieren.
type ColorStop struct {
	Pos float64
	Col LedColor
}

// Mit Paletten lassen sich anspruchsvolle Farbverlaeufe realisieren. Jeder
// Palette liegt eine Liste von Farben (die sog. Stuetzstellen) und ihre
// jeweilige Position auf dem Intervall [0, 1] zugrunde.
type GradientPalette struct {
	NameableEmbed
	stops []ColorStop
	// Mit dieser Funktion wird die Interpolation zwischen den gesetzten
	// Farbwerten realisiert.
	fnc InterpolFuncType
}

func NewGradientPalette(name string, cl ...ColorStop) *GradientPalette {
	p := &GradientPalette{}
	p.NameableEmbed.Init(name)
	p.stops = []ColorStop{
		{0.0, NewLedColor(0x000000)},
		{1.0, NewLedColor(0xFFFFFF)},
	}
	for _, cs := range cl {
		p.SetColorStop(cs)
	}
	p.fnc = LinearInterpol
	return p
}

func NewGradientPaletteByList(name string, cycle bool, cl ...LedColor) *GradientPalette {
	p := NewGradientPalette(name)
	if cycle {
		cl = append(cl, cl[0])
	}
	p.SetColorList(cl)
	return p
}

// Setzt in der Palette einen neuen Stuetzwert. Existiert bereits eine Farbe
// an dieser Position, wird sie ueberschrieben.
func (p *GradientPalette) SetColorStop(colStop ColorStop) {
	var i int
	var stop ColorStop

	if colStop.Pos < 0.0 || colStop.Pos > 1.0 {
		log.Fatalf("Positino must be in [0,1]; is: %f", colStop.Pos)
	}
	for i, stop = range p.stops {
		if stop.Pos == colStop.Pos {
			p.stops[i].Col = colStop.Col
			return
		}
		if stop.Pos > colStop.Pos {
			break
		}
	}
	p.stops = slices.Insert(p.stops, i, colStop)
}

// Verwendet die Eintraege in cl als neue Stuetzwerte der Palette.
func (p *GradientPalette) SetColorStops(cl []ColorStop) {
	for _, c := range cl {
		p.SetColorStop(c)
	}
}

func (p *GradientPalette) SetColorList(cl []LedColor) {
	if len(cl) < 2 {
		log.Fatalf("At least two colors are required!")
	}
	posStep := 1.0 / (float64(len(cl) - 1))
	for i, c := range cl[:len(cl)-1] {
		p.SetColorStop(ColorStop{float64(i) * posStep, c})
	}
	p.SetColorStop(ColorStop{1.0, cl[len(cl)-1]})

}

// Hier nun spielt die Musik: aufgrund des Wertes t (muss im Intervall [0,1]
// liegen) wird eine neue Farbe interpoliert.
func (p *GradientPalette) Color(t float64) (c LedColor) {
	var i int
	var stop ColorStop

	if t < 0.0 || t > 1.0 {
		log.Fatalf("Color: parameter t must be in [0,1] instead of %f", t)
	}
	for i, stop = range p.stops[1:] {
		if stop.Pos > t {
			break
		}
	}
	t = (t - p.stops[i].Pos) / (p.stops[i+1].Pos - p.stops[i].Pos)
	c = p.stops[i].Col.Interpolate(p.stops[i+1].Col, p.fnc(0, 1, t)).(LedColor)
	return c
}

// Palette mit 256 einzelnen dedizierten Farbwerten - kein Fading oder
// sonstige Uebergaenge.
type SlicePalette struct {
	NameableEmbed
	Colors []LedColor
}

func NewSlicePalette(name string, cl ...LedColor) *SlicePalette {
	p := &SlicePalette{}
	p.NameableEmbed.Init(name)
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

// Pseudo-Palette: liefert immer den gleichen Farbwert zurueck.
type UniformPalette struct {
	NameableEmbed
	Col LedColor
}

func NewUniformPalette(name string, color color.Color) *UniformPalette {
	p := &UniformPalette{}
	p.NameableEmbed.Init(name)
	p.Col = LedColorModel.Convert(color).(LedColor)
	return p
}

func (p *UniformPalette) Color(v float64) LedColor {
	return p.Col
}

// Mit diesem Typ kann ein fliessender Uebergang von einer Palette zu einer
// anderen realisiert werden.
type PaletteFader struct {
	Pals [2]ColorSource
	t    float64
}

func NewPaletteFader(pal ColorSource) *PaletteFader {
	p := &PaletteFader{}
	p.Pals[0] = pal
	return p
}

// Der PaletteFader implementiert das Nameable-Interface nur zur Haelfte: man
// kann den Namen der aktuell verwendeten Palette abfragen, setzen jedoch
// nicht - der PaletteFader verwendet ja bloss bestehende Paletten.
func (p *PaletteFader) Name() string {
	return p.Pals[0].Name()
}

func (p *PaletteFader) SetName(name string) {}

func (p *PaletteFader) StartFade(nextPal ColorSource, fadeTime time.Duration) bool {
	p.Pals[1] = nextPal
	anim := NewNormAnimation(fadeTime, p.Update)
	anim.Start()
	theAnimator.AddAnimations(anim)
	return true
}

func (p *PaletteFader) Update(t float64) {
	if t == 1.0 {
		p.Pals[0], p.Pals[1] = p.Pals[1], nil
		p.t = 0.0
	} else {
		p.t = t
	}
}

func (p *PaletteFader) Color(v float64) LedColor {
	c1 := p.Pals[0].Color(v)
	if p.t > 0 {
		c2 := p.Pals[1].Color(v)
		c1 = c1.Interpolate(c2, p.t).(LedColor)
	}
	return c1
}
