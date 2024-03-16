package ledgrid

import (
	"log"
)

//----------------------------------------------------------------------------

type Palette struct {
	PosList   []float64
	ColorList []LedColor
	Func      InterpolFuncType
	darkFact  float64
}

func NewPalette() *Palette {
	p := &Palette{}
	p.PosList = []float64{0.0, 1.0}
	p.ColorList = []LedColor{Black, Black}
	p.Func = LinearInterpol
	p.darkFact = 0.0
	return p
}

func (p *Palette) SetColorStop(t float64, c LedColor) {
	var newIndex int

	if t < 0.0 || t > 1.0 {
		log.Fatalf("SetColor: parameter t must be in [0, 1] instead of %f", t)
	}
	for i, pos := range p.PosList {
		if pos == t {
			p.ColorList[i] = c
			return
		}
		if pos > t {
			newIndex = i
			break
		}
	}
	p.PosList = append(p.PosList, 0.0)
	copy(p.PosList[newIndex+1:], p.PosList[newIndex:])
	p.PosList[newIndex] = t

	p.ColorList = append(p.ColorList, LedColor{})
	copy(p.ColorList[newIndex+1:], p.ColorList[newIndex:])
	p.ColorList[newIndex] = c
}

func (p *Palette) SetColorStops(cl []LedColor) {
	posStep := 1.0 / (float64(len(cl) - 1))
	p.ColorList = make([]LedColor, len(cl))
	copy(p.ColorList, cl)
	p.PosList = make([]float64, len(cl))
	for i := range len(cl) - 1 {
		p.PosList[i] = float64(i) * posStep
	}
	p.PosList[len(p.PosList)-1] = 1.0
}

func (p *Palette) Color(t float64) (c LedColor) {
	var i int
    var pos float64

	if t < 0.0 || t > 1.0 {
		log.Fatalf("Color: parameter t must be in [0, 1] instead of %f", t)
	}
	for i, pos = range p.PosList[1:] {
		if pos > t {
			break
		}
	}
	t = (t - p.PosList[i]) / (p.PosList[i+1] - p.PosList[i])
	c = p.ColorList[i].Interpolate(p.ColorList[i+1], p.Func(0, 1, t))
	c = c.Interpolate(Black, p.darkFact)
	return c
}

func (p *Palette) DarkFactor() float64 {
	return p.darkFact
}

func (p *Palette) SetDarkFactor(f float64) {
	p.darkFact = max(min(f, 1.0), 0.0)

}
