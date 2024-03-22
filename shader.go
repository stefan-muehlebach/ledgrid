package ledgrid

import (
	"image"
	"math"
)

// Der Shader verwendet zur Berechnung der darzustellenden Farben
// math. Funktionen. Dazu wird gedanklich ueber das gesamte LedGrid ein
// Koordinatensystem gelegt, welches math. korrekt ist, seinen Ursprung in der
// Mitte des LedGrid hat und so dimensioniert ist, dass der Betrag der
// groessten Koordinaten immer 1.0 ist. Mit Hilfe einer Funktion des Typs
// AnimFuncType werden dann die Farben berechnet. Die Parameter x und y sind
// Koordinaten im erwaehnten math. Koordinatensystem und mit t wird ein
// Zeitwert (in Sekunden und Bruchteilen davon) an die Funktion uebergeben.
// Der Rueckgabewert ist eine Fliesskommazahl in [0,1] und wird verwendet,
// um aus einer Palette einen Farbwert zu erhalten.

type ShaderFuncType func(x, y, t float64) float64

type Shader struct {
	field              [][]float64
	dPixel, xMin, yMax float64
	ShaderFunc         ShaderFuncType
	pal                Colorable
}

func NewShader(size image.Point, pal Colorable, fnc ShaderFuncType) *Shader {
	s := &Shader{}
	s.field = make([][]float64, size.Y)
	for i := range size.Y {
		s.field[i] = make([]float64, size.X)
	}
	s.dPixel = 2.0 / float64(max(size.X, size.Y)-1)
	ratio := float64(size.X) / float64(size.Y)
	if ratio > 1.0 {
		s.xMin = -1.0
		s.yMax = ratio * 1.0
	} else if ratio < 1.0 {
		s.xMin = ratio * -1.0
		s.yMax = 1.0
	} else {
		s.xMin = -1.0
		s.yMax = 1.0
	}
	s.ShaderFunc = fnc
	s.pal = pal
	return s
}

func (s *Shader) Update(t float64) bool {
	var col, row int
	var x, y float64

	y = s.yMax
	for row = range s.field {
		x = s.xMin
		for col = range s.field[row] {
			s.field[row][col] = s.ShaderFunc(x, y, t)
			x += s.dPixel
		}
		y -= s.dPixel
	}
	return true
}

func (s *Shader) Draw(grid *LedGrid) {
	var col, row int
	var v float64

	for row = range s.field {
		for col, v = range s.field[row] {
			c1 := s.pal.Color(v)
			c2 := grid.LedColorAt(col, row)
			grid.SetLedColor(col, row, c1.Mix(c2, Replace))
		}
	}
}

// Im folgenden Abschnitt sind ein paar Shader zusammengestellt.

var (
	p1, p2, p3, p4, p5, p6 float64
)

// Die beruehmt/beruechtigte Plasma-Animation. Die Parameter p1 - p6 sind eher
// als Konstanten zu verstehen und eignen sich nicht, um live veraendert
// zu werden.
func PlasmaShader(x, y, t float64) float64 {
	p1 = 1.2
	p2, p3, p4 = 1.2, 2.0, 3.0
	p5, p6 = 5.0, 3.0

	v1 := f1(x, y, t, p1)
	v2 := f2(x, y, t, p2, p3, p4)
	v3 := f3(x, y, t, p5, p6)
	v := (v1+v2+v3)/6.0 + 0.5
	return v
}
func f1(x, y, t, p1 float64) float64 {
	return math.Sin(x*p1 + t)
}
func f2(x, y, t, p1, p2, p3 float64) float64 {
	return math.Sin(p1*(x*math.Sin(t/p2)+y*math.Cos(t/p3)) + t)
}
func f3(x, y, t, p1, p2 float64) float64 {
	cx := 0.125*x + 0.5*math.Sin(t/p1)
	cy := 0.125*y + 0.5*math.Cos(t/p2)
	return math.Sin(math.Sqrt(100.0*(cx*cx+cy*cy)+1.0) + t)
}

// Zeichnet verschachtelte Kreisflaechen. Mit p1 kann die Geschw. und die
// Richtung der Anim. beeinflusst werden.
func CircleShader(x, y, t float64) float64 {
	p1 = 1.0 / 3.0

	_, f := math.Modf(math.Sqrt(x*x+y*y)/(2.0*math.Sqrt2) + p1*t)
	if f < 0.0 {
		f += 1.0
	}
	return f
}

// Zeichnet verschachtelte Karomuster. Mit p1 kann die Geschw. und die
// Richtung der Anim. beeinflusst werden.
func KaroShader(x, y, t float64) float64 {
	p1 = 1.0 / 3.0

	_, f := math.Modf((math.Abs(x)+math.Abs(y))/2.0 + p1*t)
	if f < 0.0 {
		f += 1.0
	}
	return f
}

// Allgemeine Funktion fuer einen animierten Farbverlauf. Mit p1 steuert man
// die Geschwindigkeit der Animation und mit p2/p3 kann festgelegt werden,
// in welche Richtung (x oder y) der Verlauf erfolgen soll.
func FadeShader(x, y, t float64) float64 {
	p1 = 1.0 / 3.0
	p2 = 0.2
	p3 = 0.5

	_, f := math.Modf(p2*x + p3*y + p1*t)
	if f < 0.0 {
		f += 1.0
	}
	return f
}
