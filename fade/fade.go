package main

import (
	"flag"
	"image"
	"log"
	"math"
	"os"

	"golang.org/x/image/math/fixed"

	gc "github.com/rthornton128/goncurses"
	"github.com/stefan-muehlebach/ledgrid"
)

const (
	KEY_SDOWN     = 0x150
	KEY_SUP       = 0x151
	KEY_SPAGEUP   = 0x18e
	KEY_SPAGEDOWN = 0x18c
)

func fixr(r image.Rectangle) fixed.Rectangle26_6 {
	return fixed.R(r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)
}

func fixp(x, y float64) fixed.Point26_6 {
	return fixed.Point26_6{X: float2fix(x), Y: float2fix(y)}
}

func float2fix(x float64) fixed.Int26_6 {
	return fixed.Int26_6(math.Round(x * 64))
}

func fix2float(x fixed.Int26_6) float64 {
	return float64(x) / 64.0
}

//----------------------------------------------------------------------------

type Polygon struct {
	lg                        *ledgrid.LedGrid
	p0, p1, p2, dp0, dp1, dp2 fixed.Point26_6
	col                       ledgrid.LedColor
}

func NewPolygon(lg *ledgrid.LedGrid, p0, p1, p2 image.Point, col ledgrid.LedColor) *Polygon {
	p := &Polygon{}
	p.lg = lg
	p.p0 = fixed.P(p0.X, p0.Y)
	p.p1 = fixed.P(p1.X, p1.Y)
	p.p2 = fixed.P(p2.X, p2.Y)
	p.dp0 = fixp(+0.01, +0.02)
	p.dp1 = fixp(+0.02, -0.01)
	p.dp2 = fixp(-0.01, -0.02)
	p.col = col
	return p
}

func (p *Polygon) Update(t float64) bool {
	r := fixr(p.lg.Bounds())

	p.p0 = p.p0.Add(p.dp0)
	p.p1 = p.p1.Add(p.dp1)
	p.p2 = p.p2.Add(p.dp2)
	if !p.p0.In(r) {
		if p.p0.X < r.Min.X || p.p0.X >= r.Max.X {
			p.dp0.X = -p.dp0.X
		} else {
			p.dp0.Y = -p.dp0.Y
		}
	}
	if !p.p1.In(r) {
		if p.p1.X < r.Min.X || p.p1.X >= r.Max.X {
			p.dp1.X = -p.dp1.X
		} else {
			p.dp1.Y = -p.dp1.Y
		}
	}
	if !p.p2.In(r) {
		if p.p2.X < r.Min.X || p.p2.X >= r.Max.X {
			p.dp2.X = -p.dp2.X
		} else {
			p.dp2.Y = -p.dp2.Y
		}
	}
	return true
}

func (p *Polygon) Draw() {
	DrawLine(p.lg, p.p0, p.p1, p.col)
	DrawLine(p.lg, p.p1, p.p2, p.col)
	DrawLine(p.lg, p.p2, p.p0, p.col)
}

type Line struct {
	lg               *ledgrid.LedGrid
	p0, p1, dp0, dp1 fixed.Point26_6
	col              ledgrid.LedColor
}

func NewLine(lg *ledgrid.LedGrid, p0, p1 image.Point, col ledgrid.LedColor) *Line {
	l := &Line{}
	l.lg = lg
	l.p0 = fixed.P(p0.X, p0.Y)
	l.p1 = fixed.P(p1.X, p1.Y)
	l.dp0 = fixp(+0.05, 0.0)
	l.dp1 = fixp(-0.05, 0.0)
	l.col = col
	return l
}

func (l *Line) Update(t float64) bool {
	r := fixr(l.lg.Bounds())

	l.p0 = l.p0.Add(l.dp0)
	l.p1 = l.p1.Add(l.dp1)
	if !l.p0.In(r) {
		if l.p0.X < r.Min.X || l.p0.X >= r.Max.X {
			l.dp0.X = -l.dp0.X
		} else {
			l.dp0.Y = -l.dp0.Y
		}
	}
	if !l.p1.In(r) {
		if l.p1.X < r.Min.X || l.p1.X >= r.Max.X {
			l.dp1.X = -l.dp1.X
		} else {
			l.dp1.Y = -l.dp1.Y
		}
	}
	return true
}

func (l *Line) Draw() {
	DrawLine(l.lg, l.p0, l.p1, l.col)
}

//----------------------------------------------------------------------------

type ColorType int

const (
	Red ColorType = iota
	Green
	Blue
	NumColors
)

var (
	width              = 10
	height             = 10
	defHost            = "raspi-2"
	defPort       uint = 5333
	defGammaValue      = 3.0
)

func main() {
	var host string
	var port uint
	var gammaValue *ledgrid.Bounded[float64]

	var client *ledgrid.PixelClient
	var grid *ledgrid.LedGrid
	var pal *ledgrid.PaletteFader
	var ch gc.Key
	var palIdx *ledgrid.Bounded[int]
	var palName string
	var palFadeTime *ledgrid.Bounded[float64]
	// var gridSize image.Point = image.Point{width, height}
	var anim *ledgrid.Animator
	var paramIdx *ledgrid.Bounded[int]
	var params []*ledgrid.Bounded[float64]
	// var speedup *ledgrid.Bounded[float64]
	var shaders []*ledgrid.Shader
	var shader *ledgrid.Shader

	log.SetOutput(os.Stderr)

	// traceFile := fmt.Sprintf("%s.trace", path.Base(os.Args[0]))
	// fhTrace, err := os.Create(traceFile)
	// if err != nil {
	// 	log.Fatal("couldn't create tracefile: ", err)
	// }
	// trace.Start(fhTrace)

	flag.StringVar(&host, "host", defHost, "Controller hostname")
	flag.UintVar(&port, "port", defPort, "Controller port")
	// flag.Float64Var(&gammaValue.Val(), "gamma", defGammaValue, "Gamma value")
	flag.Parse()

	win, err := gc.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer gc.End()

	gc.Echo(false)
	gc.CBreak(true)
	gc.Raw(true)
	gc.Cursor(0)
	win.Keypad(true)

	client = ledgrid.NewPixelClient(host, port)
	grid = ledgrid.NewLedGrid(image.Rect(0, 0, width, height))

	gammaValue = ledgrid.NewBounded(defGammaValue, 1.0, 5.0, 0.1)
	gammaValue.SetCallback(func(oldVal, newVal float64) {
		client.SetGamma(newVal, newVal, newVal)
	})

	palIdx = ledgrid.NewBounded(0, 0, len(ledgrid.PaletteNames)-1, 1)
	palIdx.Cycle = true
	palFadeTime = ledgrid.NewBounded(1.5, 0.0, 5.0, 0.1)
	palName = ledgrid.PaletteNames[palIdx.Val()]
	pal = ledgrid.NewPaletteFader(ledgrid.PaletteMap[palName])

	shaderCtrl := ledgrid.NewShaderController(grid)

	shaders = make([]*ledgrid.Shader, 5)
	shaders[0] = shaderCtrl.AddShader(ledgrid.PlasmaShader, pal)
	shaders[0].Disable()
	shaders[1] = shaderCtrl.AddShader(ledgrid.CircleShader, ledgrid.FadeBlue)
	shaders[1].Disable()
	shaders[2] = shaderCtrl.AddShader(ledgrid.CircleShader, ledgrid.FadeGreen)
	shaders[2].Disable()
	shaders[3] = shaderCtrl.AddShader(ledgrid.KaroShader, ledgrid.FadeGreen)
	shaders[3].Disable()
	shaders[4] = shaderCtrl.AddShader(ledgrid.FadeShader, ledgrid.FadeRed)
	shaders[4].Disable()

	shaderIdx := ledgrid.NewBounded(0, 0, 4, 1)
	shaderIdx.Cycle = true
	shaderIdx.SetCallback(func(oldVal, newVal int) {
		shader = shaders[newVal]
		params = make([]*ledgrid.Bounded[float64], len(shader.Params))
		for i, p := range shader.Params {
			params[i] = ledgrid.NewBounded(p.Val, p.LowerBound, p.UpperBound, p.Step)
			params[i].BindVar(&shader.Params[i].Val)
			params[i].Name = p.Name
		}
		paramIdx = ledgrid.NewBounded(0, 0, len(params)-1, 1)
		paramIdx.Cycle = true
	})

	txt := ledgrid.NewText(grid, "Stefan Mühlebach", ledgrid.White)
	// line := NewLine(grid, image.Point{0, 1}, image.Point{9, 8}, ledgrid.Blue)
	// poly := NewPolygon(grid, image.Point{0, 4}, image.Point{0, 9}, image.Point{9, 9}, ledgrid.Green)
	// speedup = ledgrid.NewBounded(1.0, 0.1, 10.0, 0.1)

	anim = ledgrid.NewAnimator(grid, client)
	// speedup.BindVar(&anim.Speedup)
	anim.AddObjects(pal, shaderCtrl, txt)
	// anim.AddObject(shader)
	// anim.AddObject(txt)
	// anim.AddObject(line)
	// anim.AddObject(poly)

mainLoop:
	for {
		win.Clear()
		win.Printf("Current palette: %s\n", palName)
		win.Printf("  q: next; a: prev\n")
		win.Printf("Fade time      : %.1f\n", palFadeTime.Val())
		win.Printf("  w: incr; s: decr\n")
		win.Printf("Gamma value(s) : %.3f\n", gammaValue.Val())
		win.Printf("  e: incr; d: decr\n")
		win.Printf("Shaders:\n")
		win.Printf("  PgUp: next; PgDn: prev:\n")
		win.Printf("  Speedup: r: incr; f: decr\n\n")
		win.Printf("  id | Name       |   V   |   A   | Spd\n")
		win.Printf("-----+------------+-------+-------+-----\n")
		for i, s := range shaders {
			if i == shaderIdx.Val() {
				win.Printf("> ")
			} else {
				win.Printf("  ")
			}
			win.Printf("%2d | %-10s | %-5v | %-5v | %.1f\n", i, s.Name, s.Visible(), s.Alive(), s.Speedup().Val())
		}
		win.Printf("-----+------------+-------+-------+-----\n\n")
		win.Printf(" shader parameters:\n")
		for i := range params {
			if i == paramIdx.Val() {
				win.Printf("> ")
			} else {
				win.Printf("  ")
			}
			win.Printf("%-5s: %5.2f\n", params[i].Name, params[i].Val())
		}
		win.Printf("\n  r: reset parameter\n")
		win.Printf("\n")
		win.Printf("  z/x: stop/continue animation\n")
		win.Printf("  ESC: quit\n")
		gc.Update()

		ch = win.GetChar()

		switch ch {
		case gc.KEY_PAGEUP:
			palIdx.Inc()
			palName = ledgrid.PaletteNames[palIdx.Val()]
			pal.StartFade(ledgrid.PaletteMap[palName], palFadeTime.Val())
		case gc.KEY_PAGEDOWN:
			palIdx.Dec()
			palName = ledgrid.PaletteNames[palIdx.Val()]
			pal.StartFade(ledgrid.PaletteMap[palName], palFadeTime.Val())
		case KEY_SPAGEUP:
			palFadeTime.Inc()
		case KEY_SPAGEDOWN:
			palFadeTime.Dec()
		case gc.KEY_HOME:
			gammaValue.Inc()
		case gc.KEY_END:
			gammaValue.Dec()
		case 'z':
			anim.Stop()
		case 'x':
			anim.Reset()
		case 'R':
			for _, p := range params {
				p.Reset()
			}
		case gc.KEY_UP:
			shaderIdx.Dec()
		case gc.KEY_DOWN:
			shaderIdx.Inc()
		case ' ':
			shader.SetAlive(!shader.Alive())
			shader.SetVisible(!shader.Visible())
		case gc.KEY_IC:
			shader.Speedup().Inc()
		case gc.KEY_DC:
			shader.Speedup().Dec()
		case KEY_SUP:
			paramIdx.Dec()
		case KEY_SDOWN:
			paramIdx.Inc()
		case '+':
			params[paramIdx.Val()].Inc()
		case '-':
			params[paramIdx.Val()].Dec()
		case gc.KEY_ESC:
			break mainLoop
		default:
			log.Printf("command unknown: [0x%x] '%s'\n", ch, ch)
		}

	}
	anim.Stop()

	grid.Clear(ledgrid.Black)
	client.Draw(grid)

	client.Close()
	// trace.Stop()
	// fhTrace.Close()
}
