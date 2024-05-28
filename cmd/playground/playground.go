package main

import (
	"flag"
	"fmt"
	"image"

	// "image/color"
	"log"
	"math"
	"os"
	"path"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"time"

	"golang.org/x/term"

	// gc "github.com/rthornton128/goncurses"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/colornames"
	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"golang.org/x/image/math/fixed"
)

const (
	KEY_TAB       = 0x009
	KEY_DOWN      = 0x102
	KEY_UP        = 0x103
	KEY_LEFT      = 0x104
	KEY_RIGHT     = 0x105
	KEY_HOME      = 0x106
	KEY_DC        = 0x14a
	KEY_IC        = 0x14b
	KEY_SDOWN     = 0x150
	KEY_SUP       = 0x151
	KEY_PAGEDOWN  = 0x152
	KEY_PAGEUP    = 0x153
	KEY_END       = 0x166
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

var (
	doCpuProf, doMemProf, doTrace       bool
	cpuProfFile, memProfFile, traceFile string
	fhCpu, fhMem, fhTrace               *os.File
)

func init() {
	cpuProfFile = fmt.Sprintf("%s.cpuprof", path.Base(os.Args[0]))
	memProfFile = fmt.Sprintf("%s.memprof", path.Base(os.Args[0]))
	traceFile = fmt.Sprintf("%s.trace", path.Base(os.Args[0]))

	flag.BoolVar(&doCpuProf, "cpuprof", false,
		"write cpu profile to "+cpuProfFile)
	flag.BoolVar(&doMemProf, "memprof", false,
		"write memory profile to "+memProfFile)
	flag.BoolVar(&doTrace, "trace", false,
		"write trace data to "+traceFile)
}

func StartProfiling() {
	var err error

	if doCpuProf {
		fhCpu, err = os.Create(cpuProfFile)
		if err != nil {
			log.Fatal("couldn't create cpu profile: ", err)
		}
		err = pprof.StartCPUProfile(fhCpu)
		if err != nil {
			log.Fatal("couldn't start cpu profiling: ", err)
		}
	}

	if doMemProf {
		fhMem, err = os.Create(memProfFile)
		if err != nil {
			log.Fatal("couldn't create memory profile: ", err)
		}
	}

	if doTrace {
		fhTrace, err = os.Create(traceFile)
		if err != nil {
			log.Fatal("couldn't create tracefile: ", err)
		}
		trace.Start(fhTrace)
	}
}

func StopProfiling() {
	if fhCpu != nil {
		pprof.StopCPUProfile()
		fhCpu.Close()
	}

	if fhMem != nil {
		runtime.GC()
		err := pprof.WriteHeapProfile(fhMem)
		if err != nil {
			log.Fatal("couldn't write memory profile: ", err)
		}
		fhMem.Close()
	}

	if fhTrace != nil {
		trace.Stop()
		fhTrace.Close()
	}
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
	defLocal           = false
	defHost            = "raspi-2"
	defPort       uint = 5333
	defGammaValue      = 3.0
)

type genericParam any

func main() {
	var local bool
	var host string
	var port uint
	var gammaValue *ledgrid.Bounded[float64]

	var client ledgrid.PixelClient
	var grid *ledgrid.LedGrid
	// var ch gc.Key
	// var palIdx *ledgrid.Bounded[int]
	// var palList *ledgrid.Listing[string]
	var palListNew *ledgrid.Listing[ledgrid.Colorable]
	var palFader *ledgrid.PaletteFader
	// var palName string
	// var palNext ledgrid.Colorable
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
		ledgrid.HorizontalShader,
		ledgrid.VerticalShader,
	}

	log.SetOutput(os.Stderr)

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Couldn't set terminal: %v", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	console := term.NewTerminal(os.Stdin, "> ")
	console.AutoCompleteCallback = func(line string, pos int, key rune) (string, int, bool) {
		log.Printf("In callback: %s, %d, %c", line, pos, key)
		return line, pos, false
	}

	// traceFile := fmt.Sprintf("%s.trace", path.Base(os.Args[0]))
	// fhTrace, err := os.Create(traceFile)
	// if err != nil {
	// 	log.Fatal("couldn't create tracefile: ", err)
	// }
	// trace.Start(fhTrace)

	flag.BoolVar(&local, "local", defLocal, "Local PixelController")
	flag.StringVar(&host, "host", defHost, "Controller hostname")
	flag.UintVar(&port, "port", defPort, "Controller port")
	flag.Parse()

	StartProfiling()

	// win, err := gc.Init()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer gc.End()

	// gc.Echo(false)
	// gc.CBreak(true)
	// gc.Raw(true)
	// gc.Cursor(0)
	// win.Keypad(true)

	if local {
		client = ledgrid.NewPixelServer(5333, "/dev/spidev0.0", 2_000_000)
	} else {
		client = ledgrid.NewNetPixelClient(host, port)
	}
	grid = ledgrid.NewLedGrid(image.Rect(0, 0, width, height))
	anim = ledgrid.NewAnimator(grid, client)

	gammaValue = ledgrid.NewBounded("gamma", defGammaValue, 1.0, 5.0, 0.1)
	gammaValue.SetCallback(func(oldVal, newVal float64) {
		client.SetGamma(newVal, newVal, newVal)
	})

	// palNameList := make([]string, 0)
	// for _, pal := range ledgrid.PaletteList {
	// 	palNameList = append(palNameList, pal.Name())
	// }
	// palList = ledgrid.NewListing("Paletten", palNameList)
	// palList.Cycle = true

	palListNew = ledgrid.NewListing("Paletten", ledgrid.PaletteList)
	palListNew.Cycle = true

	palFadeTime = ledgrid.NewBounded("Fade Time", 1.5, 0.0, 5.0, 0.1)

	for _, shaderData := range shaderList {
		pal := ledgrid.NewPaletteFader(ledgrid.HipsterPalette)
		pal.SetAlive(true)
		shader := ledgrid.NewShader(grid, shaderData, pal)
		anim.AddObjects(shader, pal)
	}

	txt := ledgrid.NewText(grid, "Benedict", ledgrid.White)
	txtPos := txt.ParamList()[1]
	txtAnim := ledgrid.NewAnimation(10*time.Second, func(t float64) {
		txtPos.SetVal((1-t)*txtPos.Min() + t*txtPos.Max())
	})
	txtAnim.AutoReverse = true
	txtAnim.RepeatCount = ledgrid.AnimationRepeatForever
	txtAnim.Start()

	fire := ledgrid.NewFire(grid)

	pictAnim := ledgrid.NewPictureAnimation(grid)
	for i := range 10 {
		fileName := fmt.Sprintf("pattern_%05d.png", i)
		pict := ledgrid.NewPicture(grid, fileName)
		pictAnim.AddPicture(pict, 100*time.Millisecond)
	}

	// cam := ledgrid.NewCamera(grid)

	blinken := ledgrid.OpenBlinkenFile("icons.bml")
	pixPal := ledgrid.NewSlicePalette("Pico08", ledgrid.Pico08Colors...)
	pixAnim := blinken.MakePixelAnimation(grid, pixPal)

	objAnim := ledgrid.NewCanvas(grid)
	// obj1 := &ledgrid.RotatingLine{
	// 	Pos:   geom.Point{25.0, 0.0},
	// 	Len:   50.0,
	// 	Speed: 1.0,
	// 	Angle: 0.0,
	// 	Color: colornames.Teal,
	// 	Width: 5.0,
	// }
	// obj2 := &ledgrid.RotatingLine{
	// 	Pos:   geom.Point{0.0, 25.0},
	// 	Len:   50.0,
	// 	Speed: 1.0,
	// 	Angle: -math.Pi / 2.0,
	// 	Color: colornames.Crimson,
	// 	Width: 5.0,
	// }
	// obj3 := &ledgrid.RotatingLine{
	// 	Pos:   geom.Point{25.0, 50.0},
	// 	Len:   50.0,
	// 	Speed: 1.0,
	// 	Angle: math.Pi,
	// 	Color: colornames.YellowGreen,
	// 	Width: 5.0,
	// }
	// obj4 := &ledgrid.RotatingLine{
	// 	Pos:   geom.Point{50.0, 25.0},
	// 	Len:   50.0,
	// 	Speed: 1.0,
	// 	Angle: math.Pi / 2.0,
	// 	Color: colornames.Purple,
	// 	Width: 5.0,
	// }
	obj5 := &ledgrid.GlowingCircle{
		Pos:         geom.Point{50.0, 50.0},
		Dir:         geom.Point{0.0, 0.0},
		Speed:       1.0,
		Radius:      []float64{0.0, 75.0},
		FillColor:   []color.Color{color.Transparent, color.Transparent},
		StrokeColor: []color.Color{colornames.Crimson, colornames.LemonChiffon},
		StrokeWidth: []float64{10.0, 10.0},
		GlowPeriod:  10 * time.Second,
	}
	// obj6 := &ledgrid.GlowingCircle{
	// 	Pos:         geom.Point{25.0, 25.0},
	// 	Dir:         geom.Point{0.3, 0.6},
	// 	Speed:       1.0,
	// 	Radius:      []float64{10.0, 15.0},
	// 	FillColor:   []color.Color{colornames.GreenYellow, colornames.Aquamarine},
	// 	StrokeColor: []color.Color{colornames.Crimson, color.Black},
	// 	StrokeWidth: []float64{5.0, 0.0},
	// 	GlowPeriod:  3 * time.Second,
	// }
	objAnim.AddObjects( /*obj1, obj2, obj3, obj4, */ obj5)

	anim.AddObjects( /*cam,*/ fire, pixAnim, pictAnim, txt, objAnim, txtAnim)

	objectIdx := ledgrid.NewBounded("obj idx", 0, 0, len(anim.Objects())-1, 1)
	objectIdx.Cycle = true
	objectIdx.SetCallback(func(oldVal, newVal int) {
		object = anim.Objects()[newVal]
		palFader = nil
		params = nil
		if obj, ok := object.(ledgrid.Palettable); ok {
			if pal, ok := obj.Palette().(*ledgrid.PaletteFader); ok {
				palFader = pal
				palListNew.SetVal(palFader.Pals[0])
			}
		}
		// if obj, ok := object.(*ledgrid.Shader); ok {
		// 	pal = obj.Pal.(*ledgrid.PaletteFader)
		// }
		if obj, ok := object.(ledgrid.Parametrizable); ok {
			params = obj.ParamList()
			paramIdx = ledgrid.NewBounded("param idx", 0, 0, len(params)-1, 1)
			paramIdx.Cycle = true
		}
	})

	// anim.Start()

	// b := make([]byte, 5)

mainLoop:
	for {
		/*
			// win.Clear()
			fmt.Printf("+----------------------------------------------------------------------------+\n")
			fmt.Printf("|                               Welcome  to  LedGrid                         |\n")
			fmt.Printf("+----------------------------------------------------------------------------+\n")
			fmt.Printf("\nGlobal Keys:\n")
			// fmt.Printf("  [Enter]               : stop/continue animation globally\n")
			fmt.Printf("  [q]                   : quit\n")
			fmt.Printf("\nGamma value(s): ")
			// win.AttrOn(gc.A_STANDOUT)
			fmt.Printf(" %.3f ", gammaValue.Val())
			// win.AttrOff(gc.A_STANDOUT)
			fmt.Printf("\n  [Home]/[End]        : incr/decr value for all colors\n")
			fmt.Printf("\nObjects:\n")
			fmt.Printf("  Cursor keys to navigate\n")
			fmt.Printf("  [Space]               : activate/deactivate object\n")
			fmt.Printf("  [PgUp]/[PgDown]       : select a new palette\n")
			fmt.Printf("  [Shift-PgUp]/[-PgDown]: incr/decr fade time ")
			// win.AttrOn(gc.A_STANDOUT)
			fmt.Printf(" %.1f sec ", palFadeTime.Val())
			// win.AttrOff(gc.A_STANDOUT)
			fmt.Printf("\n  [Tab]                 : switch to the new palette\n")
			fmt.Printf("  [Insert]/[Delete]     : incr/decr the animation speed\n\n")
			fmt.Printf(" id | Name       | Active | Spd | Palette      > Palette\n")
			fmt.Printf("----+------------+--------+-----+--------------+--------------\n")
			for i, o := range anim.Objects() {
				selObj := false
				if i == objectIdx.Val() {
					// win.AttrOn(gc.A_STANDOUT)
					selObj = true
				}
				if obj, ok := o.(ledgrid.Visualizable); ok {
					fmt.Printf(" %2d | %-10s | %-6v | %.1f |", i, obj.Name(), obj.Active(), obj.Speedup().Val())
				} else {
					continue
				}
				if obj, ok := o.(ledgrid.Palettable); ok {
					fmt.Printf(" %-12s |", obj.Palette().Name())
					if selObj {
						fmt.Printf(" %-12s\n", palListNew.Val().Name())
					} else {
						fmt.Printf("\n")
					}
				} else {
					fmt.Printf("\n")
				}
				// win.AttrOff(gc.A_STANDOUT)
			}
			fmt.Printf("----+------------+--------+-----+--------------+--------------\n")
			fmt.Printf("\nObject parameters:\n")
			fmt.Printf("  Shift-Cursor to navigate\n")
			fmt.Printf("  [+]/[-]: incr/decr parameter value\n")
			fmt.Printf("  [r]: reset parameter to default value\n\n")
			for i, p := range params {
				if i == paramIdx.Val() {
					fmt.Printf("> ")
				} else {
					fmt.Printf("  ")
				}
				fmt.Printf("%-30s: %7.2f  [%.2f - %.2f]\n",
					p.Name(), p.Val(), p.Min(), p.Max())
			}
			// gc.Update()
			// ch = win.GetChar()
		*/

		str, err := console.ReadLine()
		// len, err := os.Stdin.Read(b)
		if err != nil {
			log.Fatalf("Couldn't read from terminal: %v", err)
		}
		log.Printf("Got %d bytes from the terminal: '%s'", len(str), str)
		ch := int(str[0])

		switch ch {
		case KEY_PAGEUP:
			palListNew.Next()
		case KEY_PAGEDOWN:
			palListNew.Prev()
		case KEY_SPAGEUP:
			palFadeTime.Inc()
		case KEY_SPAGEDOWN:
			palFadeTime.Dec()
		case KEY_HOME:
			gammaValue.Inc()
		case KEY_END:
			gammaValue.Dec()
		case KEY_TAB:
			if palFader != nil {
				// palNext = ledgrid.PaletteList[palIdx.Val()]
				palFader.StartFade(palListNew.Val(), time.Duration(palFadeTime.Val()*float64(time.Second)))
			}
		// case gc.KEY_ENTER, gc.KEY_RETURN:
		// 	if anim.IsRunning() {
		// 		anim.Stop()
		// 	} else {
		// 		anim.Start()
		// 	}
		case 'r':
			for _, p := range params {
				p.Reset()
			}
		case KEY_UP:
			objectIdx.Dec()
		case KEY_DOWN:
			objectIdx.Inc()
		case ' ':
			object.SetActive(!object.Active())
		case 'a':
			object.SetAlive(!object.Alive())
		case 'v':
			object.SetVisible(!object.Visible())
		case KEY_IC:
			object.Speedup().Inc()
		case KEY_DC:
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
	// anim.Stop()

	grid.Clear(ledgrid.Black)
	client.Draw(grid)

	client.Close()
	StopProfiling()
}
