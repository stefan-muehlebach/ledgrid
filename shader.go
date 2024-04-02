package ledgrid

import (
	"time"
	"image"
	"math"
	"slices"
)

// Mit dem ShaderController koennen mehrere Shaders auf ein LedGrid angewandt
// werden. Ueber den Controller werden die einzelnen Shaders aktualisiert
// (Methode Update()) aber auch gezeichnet (Methode Draw()).
type ShaderController struct {
	VisualizableEmbed
	lg      *LedGrid
	shaders []*Shader
}

// Erstellt einen neuen Controller, der auf das LedGrid lg operiert.
func NewShaderController(lg *LedGrid) *ShaderController {
	c := &ShaderController{}
    c.VisualizableEmbed.Init()
	c.lg = lg
	c.shaders = make([]*Shader, 0)
	return c
}

// Fuegt einen neuen Shader hinzu und verwendet pal als Palette fuer
// den Shader.
func (c *ShaderController) AddShader(shd ShaderRecord, pal Colorable) *Shader {
	s := newShader(c.lg.Bounds().Size(), shd, pal)
	c.shaders = append(c.shaders, s)
	return s
}

// Aktualisiert alle Shader, die diesem Controller angehaengt sind.
func (c *ShaderController) Update(dt time.Duration) bool {
    dt = c.AnimatableEmbed.Update(dt)
	for _, s := range c.shaders {
		s.Update(dt)
	}
	return true
}

// Zeichnet das Resultat aller Shader in das LedGrid. Als Methode fuer das
// Mischen der Farben wird die Strategie 'Max' verwendet.
func (c *ShaderController) Draw() {
	var col, row int

	for row = range c.lg.Bounds().Dy() {
		for col = range c.lg.Bounds().Dx() {
			for _, s := range c.shaders {
				if !s.Visible() {
					continue
				}
				shaderColor := s.Pal.Color(s.field[row][col])
				c.lg.MixLedColor(col, row, shaderColor, Max)
			}
		}
	}
}

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

type ShaderFuncType func(x, y, t float64, p []ShaderParam) float64

type ShaderParam struct {
	Name                         string
	Val                          float64
	LowerBound, UpperBound, Step float64
}

type Shader struct {
	VisualizableEmbed
	Name               string
	field              [][]float64
	dPixel, xMin, yMax float64
	fnc                ShaderFuncType
	Params             []ShaderParam
	Pal                Colorable
}

func newShader(size image.Point, shr ShaderRecord, pal Colorable) *Shader {
	s := &Shader{}
    s.VisualizableEmbed.Init()
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
	s.Pal = pal
    s.SetShaderData(shr)
    s.Update(0)
	return s
}

func (s *Shader) SetShaderData(shr ShaderRecord) {
    s.Name = shr.n
    s.fnc = shr.f
    s.Params = slices.Clone(shr.p)
}

func (s *Shader) Param(name string) float64 {
    for _, param := range s.Params {
        if param.Name == name {
            return param.Val
        }
    }
    return 0.0
}

func (s *Shader) SetParam(name string, val float64) {
    for _, param := range s.Params {
        if param.Name == name {
            param.Val = val
            return
        }
    }
}

func (s *Shader) Update(dt time.Duration) bool {
	var col, row int
	var x, y float64

    dt = s.AnimatableEmbed.Update(dt)
	y = s.yMax
	for row = range s.field {
		x = s.xMin
		for col = range s.field[row] {
			s.field[row][col] = s.fnc(x, y, s.t0.Seconds(), s.Params)
			x += s.dPixel
		}
		y -= s.dPixel
	}
	return true
}

// Im folgenden Abschnitt sind ein paar vordefinierte Shader zusammengestellt.
type ShaderRecord struct {
	n string
	f ShaderFuncType
	p []ShaderParam
}

var (
	PlasmaShader = ShaderRecord{
		"Plasma",
		PlasmaShaderFunc,
		[]ShaderParam{
			{Name: "p1", Val: 1.2, LowerBound: 0.0, UpperBound: 10.0, Step: 0.1},
			{Name: "p2", Val: 1.6, LowerBound: 0.0, UpperBound: 10.0, Step: 0.1},
			{Name: "p3", Val: 3.0, LowerBound: 0.0, UpperBound: 10.0, Step: 0.1},
			{Name: "p4", Val: 1.5, LowerBound: 0.0, UpperBound: 10.0, Step: 0.1},
			{Name: "p5", Val: 5.0, LowerBound: 0.0, UpperBound: 10.0, Step: 0.1},
			{Name: "p6", Val: 3.0, LowerBound: 0.0, UpperBound: 10.0, Step: 0.1},
		},
	}
	CircleShader = ShaderRecord{
		"Circle",
		CircleShaderFunc,
		[]ShaderParam{
			{Name: "x", Val: 1.0, LowerBound: 0.5, UpperBound: 2.0, Step: 0.1},
			{Name: "y", Val: 1.0, LowerBound: 0.5, UpperBound: 2.0, Step: 0.1},
			{Name: "dir", Val: 1.0, LowerBound: -1.0, UpperBound: 1.0, Step: 2.0},
		},
	}
	KaroShader = ShaderRecord{
		"Karo",
		KaroShaderFunc,
		[]ShaderParam{
			{Name: "x", Val: 1.0, LowerBound: 0.5, UpperBound: 2.0, Step: 0.1},
			{Name: "y", Val: 1.0, LowerBound: 0.5, UpperBound: 2.0, Step: 0.1},
			{Name: "dir", Val: 1.0, LowerBound: -1.0, UpperBound: 1.0, Step: 2.0},
		},
	}
	LinearShader = ShaderRecord{
		"Linear",
		LinearShaderFunc,
		[]ShaderParam{
			{Name: "x", Val: 1.0, LowerBound: 0.0, UpperBound: 2.0, Step: 0.1},
			{Name: "y", Val: 0.0, LowerBound: 0.0, UpperBound: 2.0, Step: 0.1},
			{Name: "dir", Val: 1.0, LowerBound: -1.0, UpperBound: 1.0, Step: 2.0},
		},
	}
)

// Die beruehmt/beruechtigte Plasma-Animation. Die Parameter p1 - p6 sind eher
// als Konstanten zu verstehen und eignen sich nicht, um live veraendert
// zu werden.
func PlasmaShaderFunc(x, y, t float64, p []ShaderParam) float64 {
	v1 := f1(x, y, t, p[0].Val)
	v2 := f2(x, y, t, p[1].Val, p[2].Val, p[3].Val)
	v3 := f3(x, y, t, p[4].Val, p[5].Val)
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
func CircleShaderFunc(x, y, t float64, p []ShaderParam) float64 {
	x /= p[0].Val
	y /= p[1].Val
	t *= p[2].Val
	return math.Abs(math.Mod(math.Sqrt(x*x+y*y)/(2.0*math.Sqrt2)-t, 1.0))
}

// func CircleShaderFunc(x, y, t float64, p []float64) float64 {
// 	_, f := math.Modf(math.Sqrt(x*x+y*y)/(2.0*math.Sqrt2) + t)
// 	if f < 0.0 {
// 		f += 1.0
// 	}
// 	return f
// }

// Zeichnet verschachtelte Karomuster. Mit p1 kann die Geschw. und die
// Richtung der Anim. beeinflusst werden.
func KaroShaderFunc(x, y, t float64, p []ShaderParam) float64 {
	x /= p[0].Val
	y /= p[1].Val
	t *= p[2].Val
	return math.Abs(math.Mod((math.Abs(x)+math.Abs(y))/2.0-t, 1.0))
}

// func KaroShaderFunc(x, y, t float64, p []float64) float64 {
// 	_, f := math.Modf((math.Abs(x)+math.Abs(y))/2.0 + t)
// 	if f < 0.0 {
// 		f += 1.0
// 	}
// 	return f
// }

// Allgemeine Funktion fuer einen animierten Farbverlauf. Mit p1 steuert man
// die Geschwindigkeit der Animation und mit p2/p3 kann festgelegt werden,
// in welche Richtung (x oder y) der Verlauf erfolgen soll.
func LinearShaderFunc(x, y, t float64, p []ShaderParam) float64 {
	x = p[0].Val * (x + 1.0) / 4.0
	y = p[1].Val * (y + 1.0) / 4.0
	t *= p[2].Val/4.0
	return math.Abs(math.Mod(x+y-t, 1.0))
}

// func LinearShaderFunc(x, y, t float64, p []float64) float64 {
// 	_, f := math.Modf(p[1]*x + p[2]*y + p[0]*t)
// 	if f < 0.0 {
// 		f += 1.0
// 	}
// 	return f
// }
