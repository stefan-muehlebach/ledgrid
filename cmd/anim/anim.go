package main

import (
	"flag"
	"fmt"
	"image"
	"math"
	"os"
	"os/signal"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/color"
	"github.com/stefan-muehlebach/ledgrid/conf"
)

const (
	defHost = "raspi-3"
	defPort = 5333
)

var (
	width, height int
	gridSize      image.Point
	backAlpha     = 1.0
	ledGrid       *ledgrid.LedGrid
	canvas        *ledgrid.Canvas
)

//----------------------------------------------------------------------------

type LedGridProgram interface {
	Name() string
	Run(c *ledgrid.Canvas)
}

func NewLedGridProgram(name string, runFunc func(c *ledgrid.Canvas)) LedGridProgram {
	return &simpleProgram{name, runFunc}
}

type simpleProgram struct {
	name    string
	runFunc func(c *ledgrid.Canvas)
}

func (p *simpleProgram) Name() string {
	return p.name
}

func (p *simpleProgram) Run(c *ledgrid.Canvas) {
	p.runFunc(c)
}

// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------

var (
	StopContShowHideTest = NewLedGridProgram("Hide/Show vs. Suspend/Continue by Tasks",
		func(c *ledgrid.Canvas) {
			rPos1 := geom.Point{5.0, float64(height) / 2.0}
			rPos2 := geom.Point{float64(width) - 5.0, float64(height) / 2.0}
			rSize1 := geom.Point{7.0, 7.0}
			rSize2 := geom.Point{6.0, 6.0}
			rColor1 := color.SkyBlue
			rColor2 := color.GreenYellow

			r := ledgrid.NewRectangle(rPos1, rSize1, rColor1)
			c.Add(r)

			aPos := ledgrid.NewPositionAnim(r, rPos2, 4*time.Second)
			aPos.AutoReverse = true
			aPos.RepeatCount = ledgrid.AnimationRepeatForever
			aSize := ledgrid.NewSizeAnim(r, rSize2, 4*time.Second)
			aSize.AutoReverse = true
			aSize.RepeatCount = ledgrid.AnimationRepeatForever
			aColor := ledgrid.NewColorAnim(r, rColor2, 4*time.Second)
			aColor.AutoReverse = true
			aColor.RepeatCount = ledgrid.AnimationRepeatForever
			aAngle := ledgrid.NewAngleAnim(r, math.Pi, 4*time.Second)
			aAngle.AutoReverse = true
			aAngle.RepeatCount = ledgrid.AnimationRepeatForever

			aGroup := ledgrid.NewGroup(aPos, aSize, aColor, aAngle)

			aTimeline := ledgrid.NewTimeline(4 * time.Second)
			aTimeline.RepeatCount = ledgrid.AnimationRepeatForever
			aTimeline.Add(1000*time.Millisecond, ledgrid.NewSuspContAnimation(aColor))
			aTimeline.Add(1500*time.Millisecond, ledgrid.NewSuspContAnimation(aColor))
			aTimeline.Add(2500*time.Millisecond, ledgrid.NewSuspContAnimation(aAngle))
			aTimeline.Add(3000*time.Millisecond, ledgrid.NewSuspContAnimation(aAngle))
			aTimeline.Add(1900*time.Millisecond, ledgrid.NewHideShowAnimation(r))
			aTimeline.Add(2100*time.Millisecond, ledgrid.NewHideShowAnimation(r))

			aGroup.Start()
			aTimeline.Start()
		})

	SpecialCamera = NewLedGridProgram("Special camera",
		func(c *ledgrid.Canvas) {
			pos := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
			size := geom.Point{float64(width), float64(height)}

			cam := NewHistCamera(pos, size, 100)
			c.Add(cam)
			cam.Start()
		})
)


//----------------------------------------------------------------------------

func SignalHandler() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}

//----------------------------------------------------------------------------

var (
	programList = []LedGridProgram{
		FarewellGery,
		GroupTest,
		SequenceTest,
		TimelineTest,
		StopContShowHideTest,
		PathTest,
		PolygonPathTest,
		RandomWalk,
		CirclingCircles,
		ChasingCircles,
		PushingRectangles,
		RegularPolygonTest,
		GlowingPixels,
		MovingPixels,
		// ImageFilterTest,
		EffectFaderTest,
		SpecialCamera,
		BlinkenAnimation,
		MovingText,
		FixedFontTest,
		SlideTheShow,
		ShowTheShader,
		ColorFields,
		SingleImageAlign,
		FirePlace,
	}
)

func main() {
	var host string
	var port uint
	var input string
	var ch byte
	var progId int
	var runInteractive bool
	var progList string
	var gR, gG, gB float64
	var customConfName string
	var modConf conf.ModuleConfig

	for i, prog := range programList {
		id := 'a' + i
		progList += fmt.Sprintf("\n%c - %s", id, prog.Name())
	}

	flag.IntVar(&width, "width", 40, "Width of LedGrid")
	flag.IntVar(&height, "height", 10, "Height of LedGrid")
	flag.StringVar(&host, "host", defHost, "Controller hostname")
	flag.UintVar(&port, "port", defPort, "Controller port")
	flag.StringVar(&input, "prog", input, "Play one single program"+progList)
	flag.StringVar(&customConfName, "custom", "", "Use a non standard module configuration")
	flag.Parse()

	if len(input) > 0 {
		runInteractive = false
		ch = input[0]
	} else {
		runInteractive = true
	}

	if customConfName != "" {
		modConf.Load(customConfName + ".json")
		ledGrid = ledgrid.NewLedGrid(host, port, modConf)
	} else {
		ledGrid = ledgrid.NewLedGridBySize(host, port, image.Pt(width, height))
	}
	gR, gG, gB = ledGrid.Client.Gamma()

	gridSize = ledGrid.Rect.Size()
	width = gridSize.X
	height = gridSize.Y

	canvas = ledgrid.NewCanvas(gridSize)
	ledgrid.NewAnimationController(canvas, ledGrid)

	if runInteractive {
		progId = -1
		for {
			fmt.Printf("---------------------------------------\n")
			fmt.Printf("  Program\n")
			fmt.Printf("---------------------------------------\n")
			for i, prog := range programList {
				if ch >= 'a' && ch <= 'z' && i == progId {
					fmt.Printf("> ")
				} else {
					fmt.Printf("  ")
				}
				fmt.Printf("[%c] %s\n", 'a'+i, prog.Name())
			}
			fmt.Printf("---------------------------------------\n")
			fmt.Printf("  Gamma values: %.1f, %.1f, %.1f\n", gR, gG, gB)
			fmt.Printf("   +/-: increase/decreases by 0.1\n")
			fmt.Printf("---------------------------------------\n")

			fmt.Printf("Enter a character (or '0' for quit): ")

			fmt.Scanf("%s\n", &input)
			ch = input[0]
			if ch == '0' {
				break
			}

			if ch >= 'a' && ch <= 'z' {
				progId = int(ch - 'a')
				if progId < 0 || progId >= len(programList) {
					continue
				}
				// ledgrid.AnimCtrl.Stop()
				fmt.Printf("Program statistics:\n")
				fmt.Printf("  animation: %v\n", ledgrid.AnimCtrl.Watch())
				fmt.Printf("  painting : %v\n", canvas.Watch())
				fmt.Printf("  sending  : %v\n", ledGrid.Client.Watch())
				ledgrid.AnimCtrl.Purge()
				// ledgrid.AnimCtrl.Continue()
				canvas.Purge()
				ledgrid.AnimCtrl.Watch().Reset()
				canvas.Watch().Reset()
				ledGrid.Client.Watch().Reset()
				programList[progId].Run(canvas)
			}
			if ch == 'S' {
				ledgrid.AnimCtrl.Save("gobs/program01.gob")
			}
			if ch == 'L' {
				ledgrid.AnimCtrl.Suspend()
				ledgrid.AnimCtrl.Purge()
				ledgrid.AnimCtrl.Watch().Reset()
				canvas.Purge()
				canvas.Watch().Reset()
				time.Sleep(60 * time.Millisecond)
				ledgrid.AnimCtrl.Load("gobs/program01.gob")
				ledgrid.AnimCtrl.Continue()
				// fmt.Printf("canvas  >>> %+v\n", canvas)
				// for i, obj := range canvas.ObjList {
				i := 0
				for ele := canvas.ObjList.Front(); ele != nil; ele = ele.Next() {
					obj := ele.Value.(ledgrid.CanvasObject)
					if obj == nil {
						continue
					}
					fmt.Printf(">>> obj[%d] : %[2]T %+[2]v\n", i, obj)
					i++
				}
				// fmt.Printf("animCtrl>>> %+v\n", ledgrid.AnimCtrl)
				for i, anim := range ledgrid.AnimCtrl.AnimList {
					if anim == nil {
						continue
					}
					fmt.Printf(">>> anim[%d]: %[2]T %+[2]v\n", i, anim)
				}
			}
			if ch == '+' {
				gR += 0.1
				gG += 0.1
				gB += 0.1
				ledGrid.Client.SetGamma(gR, gG, gB)
			}
			if ch == '-' {
				if gR > 0.1 {
					gR -= 0.1
					gG -= 0.1
					gB -= 0.1
					ledGrid.Client.SetGamma(gR, gG, gB)
				}
			}
		}
	} else {
		if ch >= 'a' && ch <= 'z' {
			progId = int(ch - 'a')
			if progId >= 0 && progId < len(programList) {
				programList[progId].Run(canvas)
			}
		}
		fmt.Printf("Quit by Ctrl-C\n")
		SignalHandler()
	}

	ledgrid.AnimCtrl.Suspend()
	ledGrid.Clear(color.Black)
	ledGrid.Close()

	fmt.Printf("Program statistics:\n")
	fmt.Printf("  animation: %v\n", ledgrid.AnimCtrl.Watch())
	fmt.Printf("  painting : %v\n", canvas.Watch())
	fmt.Printf("  sending  : %v\n", ledGrid.Client.Watch())
}
