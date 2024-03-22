package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"math"
	"os"

	gc "github.com/rthornton128/goncurses"
	"github.com/stefan-muehlebach/ledgrid"
)

type BoundedInt struct {
	Val    int
	Lb, Ub int
	Cycle  bool
}

func NewBoundedInt(val, lb, ub int) *BoundedInt {
	i := &BoundedInt{}
	i.Val = val
	i.Lb = lb
	i.Ub = ub
	i.Cycle = false
	return i
}

func (i *BoundedInt) Incr(s int) {
	i.Val += s
	if i.Val > i.Ub {
		if i.Cycle {
			i.Val = i.Lb
		} else {
			i.Val = i.Ub
		}
	}
}

func (i *BoundedInt) Decr(s int) {
	i.Val -= s
	if i.Val < i.Lb {
		if i.Cycle {
			i.Val = i.Ub
		} else {
			i.Val = i.Lb
		}
	}
}

type BoundedFloat struct {
	Val    float64
	Lb, Ub float64
	Cycle  bool
}

func NewBoundedFloat(val, lb, ub float64) *BoundedFloat {
	i := &BoundedFloat{}
	i.Val = val
	i.Lb = lb
	i.Ub = ub
	i.Cycle = false
	return i
}

func (i *BoundedFloat) Incr(s float64) {
	i.Val += s
	if i.Val > i.Ub {
		if i.Cycle {
			i.Val = i.Lb
		} else {
			i.Val = i.Ub
		}
	}
}

func (i *BoundedFloat) Decr(s float64) {
	i.Val -= s
	if i.Val < i.Lb {
		if i.Cycle {
			i.Val = i.Ub
		} else {
			i.Val = i.Lb
		}
	}
}

type ColorType int

const (
	Red ColorType = iota
	Green
	Blue
	NumColors
)

var (
	width                = 10
	height               = 10
	defHost              = "raspi-2"
	defPort         uint = 5333
	defGammaValue        = 3.0
	framesPerSecond      = 50
	frameRefreshMs       = 1000 / framesPerSecond
	frameRefreshSec      = float64(frameRefreshMs) / 1000.0
)

//----------------------------------------------------------------------------

type Counter struct {
	size  image.Point
	bits  []bool
	color ledgrid.LedColor
}

func NewCounter(size image.Point, color ledgrid.LedColor) *Counter {
	c := &Counter{}
	c.size = size
	c.bits = make([]bool, c.size.X*c.size.Y)
	c.color = color
	return c
}

func (c *Counter) Update(t float64) {
	for i, b := range c.bits {
		if !b {
			c.bits[i] = true
			break
		} else {
			c.bits[i] = false
		}
	}
}

func (c *Counter) Draw(grid *ledgrid.LedGrid) {
	for i, b := range c.bits {
		if !b {
			continue
		}
		row := i / c.size.X
		col := i % c.size.X
		grid.SetLedColor(col, row, c.color)
	}
}

//----------------------------------------------------------------------------

type Image struct {
	size image.Point
	pal  ledgrid.Colorable
	img  []int
}

func NewImage(size image.Point, pal ledgrid.Colorable) *Image {
	i := &Image{}
	i.size = size
	i.pal = pal
	i.img = make([]int, i.size.X*i.size.Y)
	return i
}

func (i *Image) Draw(grid *ledgrid.LedGrid) {
	for idx, v := range i.img {
		row := idx / i.size.X
		col := idx % i.size.X
		fg := i.pal.Color(float64(v))
		bg := grid.LedColorAt(col, row)
		grid.SetLedColor(col, row, fg.Mix(bg, ledgrid.Blend))
	}
}

func (i *Image) SetPixels(pix [][]byte) {
	for row, data := range pix {
		for col, v := range data {
			i.img[row*i.size.X+col] = int(v)
		}
	}
}

//----------------------------------------------------------------------------

type ImageAnimation struct {
	imageList []*Image
	timeList  []float64
	Idx       int
	Cycle     bool
}

func NewImageAnimation() *ImageAnimation {
	i := &ImageAnimation{}
	i.imageList = make([]*Image, 0)
	i.timeList = make([]float64, 0)
	i.Idx = 0
	i.Cycle = true
	return i
}

func (i *ImageAnimation) AddImage(img *Image, dur float64) {
	i.imageList = append(i.imageList, img)
	if len(i.timeList) > 0 {
		dur += i.timeList[len(i.timeList)-1]
	}
	i.timeList = append(i.timeList, dur)
}

func (i *ImageAnimation) Update(t float64) bool {
	t = math.Mod(t, i.timeList[len(i.timeList)-1])
	for idx, v := range i.timeList {
		if t < v {
			i.Idx = idx
			return true
		}
	}
	return true
}

func (i *ImageAnimation) Draw(grid *ledgrid.LedGrid) {
	i.imageList[i.Idx].Draw(grid)
}

//----------------------------------------------------------------------------

type BlinkenFile struct {
	XMLName xml.Name       `xml:"blm"`
	Width   int            `xml:"width,attr"`
	Height  int            `xml:"height,attr"`
	Bits    int            `xml:"bits,attr"`
	Header  BlinkenHeader  `xml:"header"`
	Frames  []BlinkenFrame `xml:"frame"`
}

type BlinkenHeader struct {
	XMLName  xml.Name `xml:"header"`
	Title    string   `xml:"title"`
	Author   string   `xml:"author"`
	Email    string   `xml:"email"`
	Duration int      `xml:"duration"`
}

type BlinkenFrame struct {
	XMLName  xml.Name `xml:"frame"`
	Duration int      `xml:"duration,attr"`
	Rows     [][]byte `xml:"row"`
}

func ReadBlinkenFile(fileName string) *BlinkenFile {
	b := &BlinkenFile{}

	xmlFile, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	xml.Unmarshal(byteValue, b)

	for i, frame := range b.Frames {
		for j, row := range frame.Rows {
			for k, val := range row {
				if val >= '0' && val <= '9' {
					val = val - '0'
				} else {
					val = val - 'a' + 10
				}
				b.Frames[i].Rows[j][k] = val
			}
		}
	}
	return b
}

func (b *BlinkenFile) MakeImageAnimation(size image.Point, pal ledgrid.Colorable) *ImageAnimation {
	i := NewImageAnimation()

	for _, frame := range b.Frames {
		img := NewImage(size, pal)
		img.SetPixels(frame.Rows)
		i.AddImage(img, float64(frame.Duration)/1000.0)
	}

	return i
}

//----------------------------------------------------------------------------

func Main() {
	var blinkenFile *BlinkenFile
	var pal *ledgrid.DiscPalette

	blinkenFile = ReadBlinkenFile("alien.bml")

	frame := blinkenFile.Frames[20]
	fmt.Printf("%v\n", frame)

	pal = ledgrid.NewDiscPalette()
	i := NewImage(image.Point{10, 10}, pal)
	i.SetPixels(frame.Rows)
}

func main() {
	var host string
	var port uint
	var gammaValue *BoundedFloat

	var client *ledgrid.PixelClient
	var grid *ledgrid.LedGrid
	var pal *ledgrid.PaletteFader
	// var discrPal *ledgrid.DiscPalette
	var shader *ledgrid.Shader
	var ch gc.Key
	var palIdx *BoundedInt
	var palName string
	var palFadeTime *BoundedFloat
	var gridSize image.Point = image.Point{width, height}
	// var imgAnim *ImageAnimation
	// var blinkenFile *BlinkenFile
	var anim *ledgrid.Animator

	gammaValue = NewBoundedFloat(defGammaValue, 1.0, 5.0)

	flag.StringVar(&host, "host", defHost, "Controller hostname")
	flag.UintVar(&port, "port", defPort, "Controller port")
	flag.Float64Var(&gammaValue.Val, "gamma", defGammaValue, "Gamma value")
	flag.Parse()

	win, err := gc.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer gc.End()

	gc.Echo(false)
	gc.CBreak(false)
	gc.Raw(true)

	client = ledgrid.NewPixelClient(host, port)
	client.SetGamma(gammaValue.Val, gammaValue.Val, gammaValue.Val)
	grid = ledgrid.NewLedGrid(image.Rect(0, 0, width, height))

	palIdx = NewBoundedInt(0, 0, len(ledgrid.PaletteNames)-1)
	palFadeTime = NewBoundedFloat(1.5, 0.0, 5.0)
	palName = ledgrid.PaletteNames[palIdx.Val]
	pal = ledgrid.NewPaletteFader(ledgrid.PaletteMap[palName])
	shader = ledgrid.NewShader(gridSize, pal, ledgrid.KaroShader)

	// discrPal = ledgrid.NewDiscPaletteWithColors([]ledgrid.LedColor{
	// 	{0x00, 0x00, 0x00, 0x00},
	// 	{0x00, 0x00, 0x00, 0xff},
	// 	{0xff, 0x00, 0x00, 0xff},
	// 	{0x00, 0xff, 0x00, 0xff},
	// 	{0x00, 0x00, 0xff, 0xff},
	// })

	// blinkenFile = ReadBlinkenFile("alien.bml")
	// imgAnim = blinkenFile.MakeImageAnimation(gridSize, discrPal)

	anim = ledgrid.NewAnimator(grid, client)
	anim.AddObject(pal)
	anim.AddObject(shader)
	// anim.AddObject(imgAnim)

mainLoop:
	for {
		win.Clear()
		win.Printf("Current palette: %s\n", palName)
		win.Printf("  q: next; a: prev\n")
		win.Printf("Fade time      : %.1f\n", palFadeTime.Val)
		win.Printf("  w: incr; s: decr\n")
		win.Printf("Gamma value(s) : %.3f\n", gammaValue.Val)
		win.Printf("  e: incr; d: decr\n")
		win.Printf("\n")
		win.Printf("  z/x: stop/continue animation\n")
		win.Printf("  ESC: quit\n")
		gc.Update()

		ch = win.GetChar()

		switch ch {
		case 'q', 'a':
			if ch == 'q' {
				palIdx.Incr(1)
			} else {
				palIdx.Decr(1)
			}
			palName = ledgrid.PaletteNames[palIdx.Val]
			pal.StartFade(ledgrid.PaletteMap[palName], palFadeTime.Val)
		case 'w', 's':
			if ch == 'w' {
				palFadeTime.Incr(0.1)
			} else {
				palFadeTime.Decr(0.1)
			}
		case 'e', 'd':
			if ch == 'e' {
				gammaValue.Incr(0.1)
			} else {
				gammaValue.Decr(0.1)
			}
			client.SetGamma(gammaValue.Val, gammaValue.Val, gammaValue.Val)
		case 'z':
			anim.Stop()
		case 'x':
			anim.Reset()
		case gc.KEY_ESC:
			break mainLoop
		default:
			fmt.Printf("command unknown: '%s'\n", ch)
		}

	}
	anim.Stop()

	grid.Clear(ledgrid.Black)
	client.Draw(grid)

	client.Close()
}
