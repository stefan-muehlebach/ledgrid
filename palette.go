package ledgrid

import (
	"log"
)

//----------------------------------------------------------------------------

// Eine etwas 'weichere' Interpolationsfunktion.
func interpFunc(x float64) float64 {
	return 3.0*x*x - 2.0*x*x*x
}

//----------------------------------------------------------------------------

type Palette struct {
	PosList   []float64
	ColorList []LedColor
	darkFact  float64
}

func NewPalette() *Palette {
	p := &Palette{}
	p.PosList = []float64{0.0, 1.0}
	p.ColorList = []LedColor{Black, Black}
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

func (p *Palette) SetColorStops(colorList []LedColor) {
	posStep := 1.0 / (float64(len(colorList) - 1))
	p.ColorList = make([]LedColor, len(colorList))
	copy(p.ColorList, colorList)
	p.PosList = make([]float64, len(colorList))
	for i := range len(colorList) - 1 {
		p.PosList[i] = float64(i) * posStep
	}
	p.PosList[len(p.PosList)-1] = 1.0
}

func (p *Palette) Color(t float64) (c LedColor) {
	var lowerIndex int

	if t < 0.0 || t > 1.0 {
		log.Fatalf("Color: parameter t must be in [0, 1] instead of %f", t)
	}
	for i, pos := range p.PosList[1:] {
		if pos > t {
			lowerIndex = i
			break
		}
	}
	t = (t - p.PosList[lowerIndex]) / (p.PosList[lowerIndex+1] - p.PosList[lowerIndex])
	c = p.ColorList[lowerIndex].Interpolate(p.ColorList[lowerIndex+1], interpFunc(t))
	c = c.Interpolate(Black, p.darkFact)
	return c
}

func (p *Palette) DarkFactor() float64 {
	return p.darkFact
}

func (p *Palette) SetDarkFactor(f float64) {
	p.darkFact = max(min(f, 1.0), 0.0)

}
