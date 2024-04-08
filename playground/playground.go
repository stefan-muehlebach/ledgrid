package main

import (
	"flag"
	"image"
	"log"
	"math"
	"os"
	"time"

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

// type Polygon struct {
// 	lg                        *ledgrid.LedGrid
// 	p0, p1, p2, dp0, dp1, dp2 fixed.Point26_6
// 	col                       ledgrid.LedColor
// }

// func NewPolygon(lg *ledgrid.LedGrid, p0, p1, p2 image.Point, col ledgrid.LedColor) *Polygon {
// 	p := &Polygon{}
// 	p.lg = lg
// 	p.p0 = fixed.P(p0.X, p0.Y)
// 	p.p1 = fixed.P(p1.X, p1.Y)
// 	p.p2 = fixed.P(p2.X, p2.Y)
// 	p.dp0 = fixp(+0.01, +0.02)
// 	p.dp1 = fixp(+0.02, -0.01)
// 	p.dp2 = fixp(-0.01, -0.02)
// 	p.col = col
// 	return p
// }

// func (p *Polygon) Update(t float64) bool {
// 	r := fixr(p.lg.Bounds())

// 	p.p0 = p.p0.Add(p.dp0)
// 	p.p1 = p.p1.Add(p.dp1)
// 	p.p2 = p.p2.Add(p.dp2)
// 	if !p.p0.In(r) {
// 		if p.p0.X < r.Min.X || p.p0.X >= r.Max.X {
// 			p.dp0.X = -p.dp0.X
// 		} else {
// 			p.dp0.Y = -p.dp0.Y
// 		}
// 	}
// 	if !p.p1.In(r) {
// 		if p.p1.X < r.Min.X || p.p1.X >= r.Max.X {
// 			p.dp1.X = -p.dp1.X
// 		} else {
// 			p.dp1.Y = -p.dp1.Y
// 		}
// 	}
// 	if !p.p2.In(r) {
// 		if p.p2.X < r.Min.X || p.p2.X >= r.Max.X {
// 			p.dp2.X = -p.dp2.X
// 		} else {
// 			p.dp2.Y = -p.dp2.Y
// 		}
// 	}
// 	return true
// }

// func (p *Polygon) Draw() {
// 	DrawLine(p.lg, p.p0, p.p1, p.col)
// 	DrawLine(p.lg, p.p1, p.p2, p.col)
// 	DrawLine(p.lg, p.p2, p.p0, p.col)
// }

// type Line struct {
// 	lg               *ledgrid.LedGrid
// 	p0, p1, dp0, dp1 fixed.Point26_6
// 	col              ledgrid.LedColor
// }

// func NewLine(lg *ledgrid.LedGrid, p0, p1 image.Point, col ledgrid.LedColor) *Line {
// 	l := &Line{}
// 	l.lg = lg
// 	l.p0 = fixed.P(p0.X, p0.Y)
// 	l.p1 = fixed.P(p1.X, p1.Y)
// 	l.dp0 = fixp(+0.05, 0.0)
// 	l.dp1 = fixp(-0.05, 0.0)
// 	l.col = col
// 	return l
// }

// func (l *Line) Update(t float64) bool {
// 	r := fixr(l.lg.Bounds())

// 	l.p0 = l.p0.Add(l.dp0)
// 	l.p1 = l.p1.Add(l.dp1)
// 	if !l.p0.In(r) {
// 		if l.p0.X < r.Min.X || l.p0.X >= r.Max.X {
// 			l.dp0.X = -l.dp0.X
// 		} else {
// 			l.dp0.Y = -l.dp0.Y
// 		}
// 	}
// 	if !l.p1.In(r) {
// 		if l.p1.X < r.Min.X || l.p1.X >= r.Max.X {
// 			l.dp1.X = -l.dp1.X
// 		} else {
// 			l.dp1.Y = -l.dp1.Y
// 		}
// 	}
// 	return true
// }

// func (l *Line) Draw() {
// 	DrawLine(l.lg, l.p0, l.p1, l.col)
// }

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

	var client ledgrid.PixelClient
	var grid *ledgrid.LedGrid
	var pal *ledgrid.PaletteFader
	var ch gc.Key
	var palIdx *ledgrid.Bounded[int]
	var palName string
	var palNext ledgrid.Colorable
	var palFadeTime *ledgrid.Bounded[float64]
	// var gridSize image.Point = image.Point{width, height}
	var anim *ledgrid.Animator
	var paramIdx *ledgrid.Bounded[int]
	var params []*ledgrid.Bounded[float64]
	// var speedup *ledgrid.Bounded[float64]
	// var shaders []*ledgrid.Shader
	// var shader *ledgrid.Shader
	var object ledgrid.Visualizable
	var shaderList = []ledgrid.ShaderRecord{
		ledgrid.PlasmaShader,
		ledgrid.CircleShader,
		ledgrid.KaroShader,
		ledgrid.LinearShader,
		ledgrid.LinearShader,
	}

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

	client = ledgrid.NewNetPixelClient(host, port)
	grid = ledgrid.NewLedGrid(image.Rect(0, 0, width, height))
	anim = ledgrid.NewAnimator(grid, client)

	gammaValue = ledgrid.NewBounded("gamma", defGammaValue, 1.0, 5.0, 0.1)
	gammaValue.SetCallback(func(oldVal, newVal float64) {
		client.SetGamma(newVal, newVal, newVal)
	})

	palIdx = ledgrid.NewBounded("pal idx", 0, 0, len(ledgrid.PaletteList)-1, 1)
	palIdx.Cycle = true
	palFadeTime = ledgrid.NewBounded("fade time", 1.5, 0.0, 5.0, 0.1)

	for _, shaderData := range shaderList {
		pal := ledgrid.NewPaletteFader(ledgrid.DefaultPalette)
		pal.SetAlive(true)
		shader := ledgrid.NewShader(grid, shaderData, pal)
		anim.AddObjects(shader, pal)
	}

	txt := ledgrid.NewText(grid, "Stefan Mühlebach", ledgrid.BlackColor)

	fire := ledgrid.NewFire(grid)

	pict := ledgrid.NewPicture(grid, "testbild.png")

	cam := ledgrid.NewCamera(grid)

	blinken := ledgrid.OpenBlinkenFile("icons.bml")
	pixPal := ledgrid.NewSlicePalette("Pico08", ledgrid.Pico08Colors...)
	pixAnim := blinken.MakePixelAnimation(grid, pixPal)

    img := ledgrid.NewGeomImage(grid)

	anim.AddObjects(cam, pict, fire, pixAnim, txt, img)

	objectIdx := ledgrid.NewBounded("obj idx", 0, 0, len(anim.Objects())-1, 1)
	objectIdx.Cycle = true
	objectIdx.SetCallback(func(oldVal, newVal int) {
		object = anim.Objects()[newVal]
		if shader, ok := object.(*ledgrid.Shader); ok {
			pal = shader.Pal.(*ledgrid.PaletteFader)
			params = shader.ParamList()
			paramIdx = ledgrid.NewBounded("param idx", 0, 0, len(params)-1, 1)
			paramIdx.Cycle = true
		} else {
			params = nil
		}
	})

	anim.Start()

mainLoop:
	for {
		win.Clear()
		win.Printf("+---------------------------------------------------+\n")
		win.Printf("|                      Welcome                      |\n")
		win.Printf("|                        to                         |\n")
		win.Printf("|                      LedGrid                      |\n")
		win.Printf("+---------------------------------------------------+\n")
		win.Printf("\nGlobal Keys:\n")
		win.Printf("  [Enter]: stop/continue animation\n")
		win.Printf("  [q]    : quit\n")
		win.Printf("\nGamma value(s)   : ")
		win.AttrOn(gc.A_STANDOUT)
		win.Printf(" %.3f ", gammaValue.Val())
		win.AttrOff(gc.A_STANDOUT)
		win.Printf("\n  [Home]/[End]: incr/decr value for all colors\n")
		palName = ledgrid.PaletteList[palIdx.Val()].Name()
		win.Printf("\nPalette to set   : ")
		win.AttrOn(gc.A_STANDOUT)
		win.Printf(" %s ", palName)
		win.AttrOff(gc.A_STANDOUT)
		win.Printf("\n  [PgUp]/[PgDown]: select next/prev palette\n")
		win.Printf("  [Tab]          : start fade\n")
		win.Printf("\nPalette fade time: ")
		win.AttrOn(gc.A_STANDOUT)
		win.Printf(" %.1f sec ", palFadeTime.Val())
		win.AttrOff(gc.A_STANDOUT)
		win.Printf("\n  [Shift-PgUp]/[Shift-PgDown]: incr/decr fade time\n")
		win.Printf("\nObjects:\n")
		win.Printf("  Cursor keys to navigate\n")
		win.Printf("  v: toggle visibility\n")
		win.Printf("  a: toggle animatable\n")
		win.Printf("  [Insert]/[Delete]: incr/decr the speedup factor\n")
		win.Printf("\n  id | Name       |   V   |   A   | Spd\n")
		win.Printf("-----+------------+-------+-------+-----\n")
		for i, o := range anim.Objects() {
			if i == objectIdx.Val() {
				win.Printf("> ")
			} else {
				win.Printf("  ")
			}
			switch obj := o.(type) {
			case ledgrid.Visualizable:
				win.Printf("%2d | %-10s | %-5v | %-5v | %.1f\n", i, obj.Name(), obj.Visible(), obj.Alive(), obj.Speedup().Val())
			case ledgrid.Drawable:
				win.Printf("%2d | %-10T | %-5v |       | \n", i, obj, obj.Visible())
			case ledgrid.Animatable:
				win.Printf("%2d | %-10T |       | %-5v | %.1f\n", i, obj, obj.Alive(), obj.Speedup().Val())
			}
		}
		win.Printf("-----+------------+-------+-------+-----\n")
		win.Printf("\nShader parameters:\n")
		win.Printf("  Shift-Cursor to navigate\n")
		win.Printf("  [+]/[-]: incr/decr parameter value\n")
		win.Printf("  r: reset parameter to default value\n\n")
		for i := range params {
			if i == paramIdx.Val() {
				win.Printf("> ")
			} else {
				win.Printf("  ")
			}
			win.Printf("%-5s: %5.2f\n", params[i].Name(), params[i].Val())
		}
		gc.Update()

		ch = win.GetChar()

		switch ch {
		case gc.KEY_PAGEUP:
			palIdx.Inc()
		case gc.KEY_PAGEDOWN:
			palIdx.Dec()
		case KEY_SPAGEUP:
			palFadeTime.Inc()
		case KEY_SPAGEDOWN:
			palFadeTime.Dec()
		case gc.KEY_HOME:
			gammaValue.Inc()
		case gc.KEY_END:
			gammaValue.Dec()
		case gc.KEY_TAB:
			log.Printf("Change palette\n")
			palNext = ledgrid.PaletteMap[palName]
			pal.StartFade(palNext, time.Duration(palFadeTime.Val()*float64(time.Second)))
		case gc.KEY_ENTER, gc.KEY_RETURN:
			if anim.IsRunning() {
				anim.Stop()
			} else {
				anim.Start()
			}
		case 'r':
			for _, p := range params {
				p.Reset()
			}
		case gc.KEY_UP:
			objectIdx.Dec()
		case gc.KEY_DOWN:
			objectIdx.Inc()
		case ' ':
			object.SetActive(!object.Active())
		case 'a':
			object.SetAlive(!object.Alive())
		case 'v':
			object.SetVisible(!object.Visible())
		case gc.KEY_IC:
			object.Speedup().Inc()
		case gc.KEY_DC:
			object.Speedup().Dec()
		case KEY_SUP:
			paramIdx.Dec()
		case KEY_SDOWN:
			paramIdx.Inc()
		case '+':
			params[paramIdx.Val()].Inc()
		case '-':
			params[paramIdx.Val()].Dec()
		case 'q':
			break mainLoop
		default:
			log.Printf("command unknown: [0x%x] '%s'\n", ch, ch)
		}

	}
	anim.Stop()

	grid.Clear(ledgrid.BlackColor)
	client.Draw(grid)

	client.Close()
	// trace.Stop()
	// fhTrace.Close()
}
