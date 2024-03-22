package ledgrid

import (
	"math"
	"log"
)

type Index interface {
    int | float64
}

// Alles, was im Sinne einer Farbpalette Farben erzeugen kann, implementiert
// das Colorable Interface.
type Colorable interface {
    Color(v float64) LedColor
}

// Mit Paletten lassen sich anspruchsvolle Farbverlaeufe realisieren. Jeder
// Palette liegt eine Liste von Farben (die sog. Stuetzstellen) und ihre
// jeweilige Position auf dem Intervall [0, 1] zugrunde.
type Palette struct {
    // Die Liste der Positionen. Dabei muss PosList[0] = 0.0 und
    // PosList[len(PosList)-1] = 1.0 sein. Ausserdem muessen die Positionen
    // in aufsteigender Reihenfolge sortiert sein.
	PosList   []float64
    // Dies sind die Farbwerte an den jeweiligen Positionen. Es muss gelten
    // len(PosList) = len(ColorList).
	ColorList []LedColor
    // Mit dieser Funktion wird die Interpolation zwischen den gesetzten
    // Farbwerten realisiert.
	Func      InterpolFuncType
}

// Erzeugt eine neue Palette, welche einen Farbverlauf von Schwarz (0.0) nach
// Weiss (1.0) beinhaltet. Mit der Funktion SetColorStop koennen bestehende
// Stuetzstellen ersetzt oder neue hinzugefuegt werden.
func NewPalette() *Palette {
	p := &Palette{}
	p.PosList = []float64{0.0, 1.0}
	p.ColorList = []LedColor{Black, White}
	p.Func = LinearInterpol
	return p
}

// Erzeugt eine neue Palette und verwendet die Farben in cl als Stuetzwerte.
// In diesem Fall werden die Farben in cl gleichmaessig (aequidistant) auf
// dem Intervall [0,1] verteilt.
func NewPaletteWithColors(cl []LedColor) *Palette {
    if len(cl) < 2 {
        log.Fatalf("At least two colors are required!")
    }
    p := NewPalette()
    if cl[0] != cl[len(cl)-1] {
        cl = append(cl, cl[0])
    }
    p.SetColorStops(cl)
    return p
}

// Setzt die Farbe c als neuen Stuetzwert bei Position t. Existiert bereits
// eine Farbe mit dieser Position, wird sie ueberschrieben.
func (p *Palette) SetColorStop(t float64, c LedColor) {
	var i int
    var pos float64

	if t < 0.0 || t > 1.0 {
		log.Fatalf("SetColor: parameter t must be in [0, 1] instead of %f", t)
	}
	for i, pos = range p.PosList {
		if pos == t {
			p.ColorList[i] = c
			return
		}
		if pos > t {
			break
		}
	}
	p.PosList = append(p.PosList, 0.0)
	copy(p.PosList[i+1:], p.PosList[i:])
	p.PosList[i] = t

	p.ColorList = append(p.ColorList, LedColor{})
	copy(p.ColorList[i+1:], p.ColorList[i:])
	p.ColorList[i] = c
}

// Farbwerte in cl werden als Stuetzstellen der Palett verwendet. Die
// Stuetzstellen sind gleichmaessig ueber das Intervall [0,1] verteilt.
func (p *Palette) SetColorStops(cl []LedColor) {
	posStep := 1.0 / (float64(len(cl) - 1))
	p.ColorList = make([]LedColor, len(cl))
	copy(p.ColorList, cl)
	p.PosList = make([]float64, len(cl))
	for i := range len(cl) - 1 {
		p.PosList[i] = float64(i) * posStep
	}
    // Dies muss sein, da durch Rundungsfehler der letzte Positionswert nicht
    // immer gleich 1.0 ist.
	p.PosList[len(p.PosList)-1] = 1.0
}

// Hier nun spielt die Musik: aufgrund des Wertes t (muss im Intervall [0,1]
// liegen) wird eine neue Farbe interpoliert.
func (p *Palette) Color(t float64) (c LedColor) {
	var i int
    var pos float64

	if t < 0.0 || t > 1.0 {
		log.Fatalf("Color: parameter t must be in [0,1] instead of %f", t)
	}
	for i, pos = range p.PosList[1:] {
		if pos > t {
			break
		}
	}
	t = (t - p.PosList[i]) / (p.PosList[i+1] - p.PosList[i])
	c = p.ColorList[i].Interpolate(p.ColorList[i+1], p.Func(0, 1, t))
	return c
}

// Palette mit 256 einzelnen Farbwerten
type DiscPalette struct {
    Colors []LedColor
}

func NewDiscPalette() *DiscPalette {
    p := &DiscPalette{}
    p.Colors = make([]LedColor, 256)
    for idx := range p.Colors {
        p.Colors[idx] = LedColor{}
    }
    return p
}

func NewDiscPaletteWithColors(cl []LedColor) *DiscPalette {
    p := NewDiscPalette()
    for i, c := range cl {
        p.SetColor(i, c)
    }
    return p
}

func (p *DiscPalette) Color(v float64) LedColor {
    _, f := math.Modf(v)
    if f != 0.0 {
        v = v * 255.0
    }
    return p.Colors[int(v)]
}

func (p *DiscPalette) SetColor(idx int, c LedColor) {
    p.Colors[idx] = c
}

// Mit diesem Typ kann ein fliessender Uebergang von einer Palette zu einer
// anderen realisiert werden.
type PaletteFader struct {
    Pals [2]Colorable
    FadePos, FadeStep float64
}

func NewPaletteFader(pal Colorable) *PaletteFader {
    p := &PaletteFader{}
    p.Pals[0] = pal
    p.FadePos = 0.0
    p.FadeStep = 0.0
    return p
}

func (p *PaletteFader) StartFade(pal Colorable, fadeTimeSec float64) {
    // Solange noch ein Uebergang am Laufen ist, kann kein neuer gestartet
    // werden -- oder doch? TO DO!
    if p.FadePos > 0.0 {
        return
    }
    p.Pals[0], p.Pals[1] = pal, p.Pals[0]
    if fadeTimeSec > 0.0 {
        p.FadePos = 1.0
        p.FadeStep = 1.0 / (fadeTimeSec / frameRefreshSec)
    }
}

func (p *PaletteFader) Update(t float64) bool {
	if p.FadePos > 0.0 {
		p.FadePos -= p.FadeStep
		if p.FadePos < 0.0 {
			p.FadePos = 0.0
		}
	}
    return true
}

func (p *PaletteFader) Color(v float64) (LedColor) {
    c1 := p.Pals[0].Color(v)
	if p.FadePos > 0.0 {
		c2 := p.Pals[1].Color(v)
		c1 = c1.Interpolate(c2, p.FadePos)
	}
    return c1
}
