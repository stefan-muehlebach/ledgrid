package ledgrid

import (
	"image"
	"image/color"
	"math"
	"slices"
)

// Der Shader verwendet zur Berechnung der darzustellenden Farben
// math. Funktionen. Dazu wird gedanklich ueber das gesamte LedGrid ein
// Koordinatensystem gelegt, welches math. korrekt ist, seinen Ursprung in der
// Mitte des LedGrid hat und so dimensioniert ist, dass der Betrag der
// groessten Koordinaten immer 1.0 ist. Mit Hilfe einer Funktion des Typs
// ShaderFuncType werden dann die Farben berechnet. Die Parameter x und y sind
// Koordinaten im erwaehnten math. Koordinatensystem und mit t wird ein
// Zeitwert (in Sekunden und Bruchteilen davon) an die Funktion uebergeben.
// Der Rueckgabewert ist eine Fliesskommazahl in [0,1] und wird verwendet,
// um aus einer Palette einen Farbwert zu erhalten.

// Jeder Shader basiert auf einer Funktion mit diesem Profil. x und y sind
// Koordinaten des darzustellenden Punktes (siehe Text oben für die
// Dimensionierung des Koord.system), t ist ein fortlaufender Zeitparameter
// und p ist ein Slice von Parametern, die für diesen Shader verwendet werden
// (siehe dazu auch den Typ ShaderParam).
type ShaderFuncType func(x, y, t float64, p []*Bounded[float64]) float64

type Shader struct {
	VisualEmbed
	lg                 *LedGrid
	rect               image.Rectangle
	field              [][]float64
	dPixel, xMin, yMax float64
	fnc                ShaderFuncType
	params             []*Bounded[float64]
	pal                ColorSource
	anim               Animation
}

func NewShader(lg *LedGrid, shr ShaderRecord, pal ColorSource) *Shader {
	s := &Shader{}
	s.VisualEmbed.Init("Shader")
	s.lg = lg
	s.rect = lg.Bounds()
	s.field = make([][]float64, s.rect.Dy())
	for i := range s.rect.Dy() {
		s.field[i] = make([]float64, s.rect.Dx())
	}
	s.dPixel = 2.0 / float64(max(s.rect.Dx(), s.rect.Dy())-1)
	ratio := float64(s.rect.Dx()) / float64(s.rect.Dy())
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
	s.pal = pal
	s.SetShaderData(shr)
	s.anim = NewInfAnimation(s.Update)
	return s
}

func (s *Shader) SetShaderData(shr ShaderRecord) {
	s.name = shr.n
	s.fnc = shr.f
	s.params = slices.Clone(shr.p)
}

func (s *Shader) ParamList() []*Bounded[float64] {
	return s.params
}

func (s *Shader) Param(name string) *Bounded[float64] {
	for _, param := range s.params {
		if param.Name() == name {
			return param
		}
	}
	return nil
}

func (s *Shader) Palette() ColorSource {
	return s.pal
}

func (s *Shader) SetVisible(vis bool) {
	if vis {
		s.anim.Start()
	} else {
		s.anim.Stop()
	}
	s.VisualEmbed.SetVisible(vis)
}

func (s *Shader) Update(t float64) {
	var col, row int
	var x, y float64

	y = s.yMax
	for row = range s.field {
		x = s.xMin
		for col = range s.field[row] {
			s.field[row][col] = s.fnc(x, y, t, s.params)
			x += s.dPixel
		}
		y -= s.dPixel
	}
}

func (s *Shader) ColorModel() color.Model {
	return LedColorModel
}

func (s *Shader) Bounds() image.Rectangle {
	return s.rect
}

func (s *Shader) At(x, y int) color.Color {
	return s.pal.Color(s.field[y][x])
}

// Im folgenden Abschnitt sind ein paar vordefinierte Shader zusammengestellt.

type ShaderRecord struct {
	n string
	f ShaderFuncType
	p []*Bounded[float64]
}

var (
	PlasmaShader = ShaderRecord{
		"Plasma (Shader)",
		PlasmaShaderFunc,
		[]*Bounded[float64]{
			NewBounded("Func1, P1", 1.2, 0.0, 10.0, 0.1),
			NewBounded("Func2, P1", 1.6, 0.0, 10.0, 0.1),
			NewBounded("Func2, P2", 3.0, 0.0, 10.0, 0.1),
			NewBounded("Func2, P3", 1.5, 0.0, 10.0, 0.1),
			NewBounded("Func3, P1", 5.0, 0.0, 10.0, 0.1),
			NewBounded("Func3, P2", 5.0, 0.0, 10.0, 0.1),
		},
	}
	CircleShader = ShaderRecord{
		"Circle (Shader)",
		CircleShaderFunc,
		[]*Bounded[float64]{
			NewBounded("X-Scale", 1.0, -2.0, 2.0, 0.1),
			NewBounded("Y-Scale", 1.0, -2.0, 2.0, 0.1),
		},
	}
	KaroShader = ShaderRecord{
		"Karo (Shader)",
		KaroShaderFunc,
		[]*Bounded[float64]{
			NewBounded("X-Scale", 1.0, -2.0, 2.0, 0.1),
			NewBounded("Y-Scale", 1.0, -2.0, 2.0, 0.1),
		},
	}
	LinearShader = ShaderRecord{
		"Linear (Shader)",
		LinearShaderFunc,
		[]*Bounded[float64]{
			NewBounded("X-Speed", 1.0, -2.0, 2.0, 0.1),
			NewBounded("Y-Speed", 0.0, -2.0, 2.0, 0.1),
		},
	}
)

// Die beruehmt/beruechtigte Plasma-Animation. Die Parameter p1 - p6 sind eher
// als Konstanten zu verstehen und eignen sich nicht, um live veraendert
// zu werden.
func PlasmaShaderFunc(x, y, t float64, p []*Bounded[float64]) float64 {
	v1 := f1(x, y, t, p[0].Val())
	v2 := f2(x, y, t, p[1].Val(), p[2].Val(), p[3].Val())
	v3 := f3(x, y, t, p[4].Val(), p[5].Val())
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
func CircleShaderFunc(x, y, t float64, p []*Bounded[float64]) float64 {
	x = p[0].Val() * x / 10.0
	y = p[1].Val() * y / 10.0
	t /= 5.0
	return math.Abs(math.Mod(math.Sqrt(x*x+y*y)-t, 1.0))
	// return math.Abs(math.Mod(math.Sqrt(x*x+y*y)-t, 1.0))
}

// Zeichnet verschachtelte Karomuster. Mit p1 kann die Geschw. und die
// Richtung der Anim. beeinflusst werden.
func KaroShaderFunc(x, y, t float64, p []*Bounded[float64]) float64 {
	x = p[0].Val() * x / 10.0
	y = p[1].Val() * y / 10.0
	t /= 5.0
	return math.Abs(math.Mod(math.Abs(x)+math.Abs(y)-t, 1.0))
}

// Allgemeine Funktion fuer einen animierten Farbverlauf. Mit p1 steuert man
// die Geschwindigkeit der Animation und mit p2/p3 kann festgelegt werden,
// in welche Richtung (x oder y) der Verlauf erfolgen soll.
func LinearShaderFunc(x, y, t float64, p []*Bounded[float64]) float64 {
	x = p[0].Val() * x / 10.0
	y = p[1].Val() * y / 10.0
	t /= 5.0
	return math.Abs(math.Mod(x+y-t, 1.0))
}
