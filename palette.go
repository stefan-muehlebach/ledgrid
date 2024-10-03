package ledgrid

import (
	"image"
	gocolor "image/color"
	"log"
	"math"
	"slices"
    ledcolor "github.com/stefan-muehlebach/ledgrid/color"
)

// Alles, was im Sinne einer Farbpalette Farben erzeugen kann, implementiert
// das ColorSource Interface.
type ColorSource interface {
	// Liefert in Abhaengigkeit des Parameters v eine Farbe aus der Palette
	// zurueck. v kann vielfaeltig verwendet, resp. interpretiert werden,
	// bsp. als Parameter im Intervall [0,1], als Index (natuerliche Zahl)
	// einer Farbenliste oder gar nicht, wenn die Farbquelle einfarbig ist.
	Color(v float64) ledcolor.LedColor
}

var (
	// Alle vorhandenen Paletten sind in diesem Slice aufgefuehrt. Falls
	// applikatorisch weitere Paletten erzeugt werden, ist es Aufgabe der
	// Applikation, diesen Slice nachzufuehren.
	// PaletteList = []ColorSource{}
	PaletteNames = []string{}
	// Im Gegensatz zu [PaletteList] sind hier die Paletten unter ihrem
	// Namen abgelegt. Siehe auch Kommentar bei [PaletteList] betr.
	// Nachfuehrung.
	PaletteMap = map[string]ColorSource{}

	ColorNames = []string{}
	// ColorList = []ColorSource{}
	ColorMap = map[string]ColorSource{}
)

// Gradienten-Paletten basieren auf einer Anzahl Farben (Stuetzstellen)
// zwischen denen eine Farbe interpoliert werden kann. Jede Stuetzstelle
// besteht aus einer Position (Zahl im Intervall [0,1]) und einer dazu
// gehoerenden Farbe.
type GradientPalette struct {
	stops []ColorStop
	// Mit dieser Funktion wird die Interpolation zwischen den gesetzten
	// Farbwerten realisiert.
	// fnc InterpolFuncType
}

// Dieser (interne) Typ wird verwendet, um einen bestimmten Wert im Interval
// [0,1] mit einer Farbe zu assoziieren.
type ColorStop struct {
	Pos   float64
	Color ledcolor.LedColor
}

// Erzeugt eine neue Palette unter Verwendung der Stuetzwerte in cl. Die
// Stuetzwerte muessen nicht sortiert sein. Per Default ist 0.0 mit Schwarz
// und 1.0 mit Weiss vorbelegt - sofern in cl keine Stuetzwerte fuer 0.0 und
// 1.0 angegeben sind.
func NewGradientPalette(name string, cl ...ColorStop) *GradientPalette {
	p := &GradientPalette{}
	p.stops = []ColorStop{
		{0.0, ledcolor.NewLedColorHex(0x000000)},
		{1.0, ledcolor.NewLedColorHex(0xFFFFFF)},
	}
	for _, cs := range cl {
		p.SetColorStop(cs)
	}
	// p.fnc = LinearInterpol
	return p
}

// Erzeugt eine neue Gradienten-Palette unter Verwendung der Farben in cl,
// wobei die Farben aequidistant ueber das Interval [0,1] verteilt werden.
// Ist cycle true, dann wird die erste Farbe in cl auch als letzte Farbe
// verwendet.
func NewGradientPaletteByList(name string, cycle bool, cl ...ledcolor.LedColor) *GradientPalette {
	if len(cl) < 2 {
		log.Fatalf("At least two colors are required!")
	}
	if cycle {
		cl = append(cl, cl[0])
	}
	stops := make([]ColorStop, len(cl))
	posStep := 1.0 / (float64(len(cl) - 1))
	for i, c := range cl[:len(cl)-1] {
		stops[i] = ColorStop{float64(i) * posStep, ledcolor.LedColorModel.Convert(c).(ledcolor.LedColor)}
	}
	stops[len(cl)-1] = ColorStop{1.0, ledcolor.LedColorModel.Convert(cl[len(cl)-1]).(ledcolor.LedColor)}
	return NewGradientPalette(name, stops...)
}

// Setzt in der Palette einen neuen Stuetzwert. Existiert bereits eine Farbe
// an dieser Position, wird sie ueberschrieben.
func (p *GradientPalette) SetColorStop(colStop ColorStop) {
	var i int
	var stop ColorStop

	if colStop.Pos < 0.0 || colStop.Pos > 1.0 {
		log.Fatalf("Position must be in [0,1]; is: %f", colStop.Pos)
	}
	for i, stop = range p.stops {
		if stop.Pos == colStop.Pos {
			p.stops[i].Color = colStop.Color
			return
		}
		if stop.Pos > colStop.Pos {
			break
		}
	}
	p.stops = slices.Insert(p.stops, i, colStop)
}

// Retourniert den Slice mit den Stutzwerten (wird wohl eher fuer Debugging
// verwendet).
func (p *GradientPalette) ColorStops() []ColorStop {
	return p.stops
}

// Hier nun spielt die Musik: aufgrund des Wertes t (muss im Intervall [0,1]
// liegen) wird eine neue Farbe interpoliert.
func (p *GradientPalette) Color(t float64) (c ledcolor.LedColor) {
	var i int
	var stop ColorStop

	if t < 0.0 || t > 1.0 {
        t = max(0.0, min(1.0, t))
	}
	for i, stop = range p.stops[1:] {
		if stop.Pos > t {
			break
		}
	}
	t = (t - p.stops[i].Pos) / (p.stops[i+1].Pos - p.stops[i].Pos)
	c = p.stops[i].Color.Interpolate(p.stops[i+1].Color, t)
	return c
}

// Palette mit 256 einzelnen dedizierten Farbwerten - kein Fading oder
// sonstige Uebergaenge.
type SlicePalette struct {
	Colors []ledcolor.LedColor
}

func NewSlicePalette(name string, cl ...ledcolor.LedColor) *SlicePalette {
	p := &SlicePalette{}
	p.Colors = make([]ledcolor.LedColor, 256)
	for i, c := range cl {
		p.Colors[i] = c
	}
	return p
}

func (p *SlicePalette) Color(v float64) ledcolor.LedColor {
	return p.Colors[int(v)]
}

func (p *SlicePalette) SetColor(idx int, c ledcolor.LedColor) {
	p.Colors[idx] = c
}

// Damit auch einzelne Farben wie Paletten verwendet werden koennen,
// existiert der Typ UniformPalette. Die Ueberlegungen dazu sind analog zum
// Typ [image.Uniform].
type UniformPalette struct {
	Col ledcolor.LedColor
}

// Erstellt eine neue einfarbige Farbquelle mit gegebenem namen.
func NewUniformPalette(name string, color ledcolor.LedColor) *UniformPalette {
	p := &UniformPalette{}
	p.Col = color
	return p
}

// Damit wird das ColorSource-Interface implementiert. Der Parameter [v] hat
// bei dieser Farbquelle keine Bedeutung und wird ignoriert.
func (p *UniformPalette) Color(v float64) ledcolor.LedColor {
	return p.Col
}

func (p *UniformPalette) ColorModel() gocolor.Model {
	return ledcolor.LedColorModel
}

func (p *UniformPalette) Bounds() image.Rectangle {
	return image.Rect(math.MinInt, math.MinInt, math.MaxInt, math.MaxInt)
}

func (p *UniformPalette) At(x, y int) ledcolor.LedColor {
	return p.Col
}

func (p *UniformPalette) Set(x, y int, c ledcolor.LedColor) {}

// Mit diesem Typ kann ein fliessender Uebergang von einer Palette zu einer
// anderen realisiert werden.
type PaletteFader struct {
	Pals [2]ColorSource
	T    float64
}

// Initialisiert wird der Fader mit der aktuell anzuzeigenden Palette. Der
// PaletteFader wird anschliessend anstelle der ueblichen Palette verwendet.
func NewPaletteFader(pal ColorSource) *PaletteFader {
	p := &PaletteFader{}
	p.Pals[0] = pal
	p.Pals[1] = nil
	p.T = 0.0
	return p
}

// Mit dieser Methode wird der aktuelle Farbwert retourniert. Damit
// implementiert der Fader das ColorSource-Interface und kann als Farbquelle
// verwendet werden - genau wie anderen Paletten-, resp. Farbtypen.
func (p *PaletteFader) Color(v float64) (c ledcolor.LedColor) {
	c = p.Pals[0].Color(v)
	if p.T > 0 {
		c2 := p.Pals[1].Color(v)
		c = c.Interpolate(c2, p.T)
	}
	return c
}
