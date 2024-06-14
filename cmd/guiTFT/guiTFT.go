package main

import (
	"flag"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"os/signal"

	"github.com/stefan-muehlebach/adagui"
	"github.com/stefan-muehlebach/adagui/binding"
	"github.com/stefan-muehlebach/adagui/props"
	"github.com/stefan-muehlebach/adagui/touch"
	"github.com/stefan-muehlebach/adatft"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/colornames"
	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
)

//-----------------------------------------------------------------------

var (
	PickerProps = props.NewProperties(props.PropsMap["Default"])
)

func init() {
	PickerProps.SetColor(props.Color, colornames.Black)
	PickerProps.SetColor(props.BorderColor, colornames.WhiteSmoke)
	PickerProps.SetSize(props.BorderWidth, 2.0)
}

type ColorPicker struct {
	adagui.ContainerEmbed
	orient           adagui.Orientation
	colorList        [][]ledgrid.LedColor
	fieldSize        float64
	numCols, numRows int
	selIdx0, selIdx1 int
	value            binding.Untyped
}

func NewColorPicker(fieldSize float64, orient adagui.Orientation,
	numColsRows int) *ColorPicker {
	c := &ColorPicker{}
	c.Wrapper = c
	c.Init()
	c.PropertyEmbed.Init(PickerProps)
	c.orient = orient
	c.colorList = make([][]ledgrid.LedColor, numColsRows)
	c.fieldSize = fieldSize
	switch orient {
	case adagui.Horizontal:
		c.numCols = 1
		c.numRows = numColsRows
	case adagui.Vertical:
		c.numCols = numColsRows
		c.numRows = 1
	}
	c.SetMinSize(geom.Point{c.fieldSize * float64(c.numCols),
		c.fieldSize * float64(c.numRows)})
	c.selIdx0, c.selIdx1 = 0, 0
	c.value = binding.NewUntyped()
	return c
}

func NewColorPickerWithCallback(fieldSize float64, orient adagui.Orientation,
	numColsRows int, callback func(color.Color)) *ColorPicker {
	c := NewColorPicker(fieldSize, orient, numColsRows)
	c.value.AddCallback(func(data binding.DataItem) {
		callback(data.(binding.Untyped).Get().(color.Color))
	})
	return c
}

func NewColorPickerWithData(fieldSize float64, orient adagui.Orientation,
	numColsRows int, data binding.Untyped) *ColorPicker {
	c := NewColorPicker(fieldSize, orient, numColsRows)
	c.value = data
	return c
}

func (c *ColorPicker) SetColors(colRowIdx int, cl []ledgrid.LedColor) {
	c.colorList[colRowIdx] = make([]ledgrid.LedColor, len(cl))
	copy(c.colorList[colRowIdx], cl)
	switch c.orient {
	case adagui.Horizontal:
		c.numCols = max(c.numCols, len(cl))
	case adagui.Vertical:
		c.numRows = max(c.numRows, len(cl))
	}
	c.SetMinSize(geom.Point{c.fieldSize * float64(c.numCols),
		c.fieldSize * float64(c.numRows)})
}

func (c *ColorPicker) Color() color.Color {
	return c.value.Get().(color.Color)
}

func (c *ColorPicker) SetColor(col color.Color) {
	ledColor := col.(ledgrid.LedColor)
	for c.selIdx0 = range c.colorList {
		for c.selIdx1 = range c.colorList[c.selIdx0] {
			if ledColor == c.colorList[c.selIdx0][c.selIdx1] {
				c.value.Set(ledColor)
				return
			}
		}
	}
}

func (c *ColorPicker) OnInputEvent(evt touch.Event) {
	//log.Printf("ColorPicker: %v", evt.Pos)
	switch evt.Type {
	case touch.TypeRelease:
		selPt := c.Bounds().SetInside(evt.Pos)
		dx, dy := selPt.AsCoord()
		col := int(dx / c.fieldSize)
		row := int(dy / c.fieldSize)
		switch c.orient {
		case adagui.Horizontal:
			c.selIdx0 = row
			c.selIdx1 = col
		case adagui.Vertical:
			c.selIdx0 = col
			c.selIdx1 = row
		}
		if c.selIdx0 >= len(c.colorList) || c.selIdx1 >= len(c.colorList[c.selIdx0]) {
			return
		}
		c.value.Set(c.colorList[c.selIdx0][c.selIdx1])
		c.Mark(adagui.MarkNeedsPaint)
	}
}

func (c *ColorPicker) Paint(gc *gg.Context) {
	var col, row float64
	var r geom.Rectangle

	gc.SetStrokeWidth(0.0)
	for i, colorSlice := range c.colorList {
		for j, color := range colorSlice {
			gc.SetFillColor(color)
			switch c.orient {
			case adagui.Horizontal:
				col = float64(j) * c.fieldSize
				row = float64(i) * c.fieldSize
			case adagui.Vertical:
				col = float64(i) * c.fieldSize
				row = float64(j) * c.fieldSize
			}
			gc.DrawRectangle(col, row, c.fieldSize, c.fieldSize)
			gc.FillStroke()
		}
	}

	gc.SetFillColor(c.Color())
	gc.SetStrokeColor(c.BorderColor())
	gc.SetStrokeWidth(c.BorderWidth())
	gc.DrawRectangle(c.Bounds().AsCoord())
	gc.Stroke()

	switch c.orient {
	case adagui.Horizontal:
		r = geom.NewRectangleWH(float64(c.selIdx1)*c.fieldSize,
			float64(c.selIdx0)*c.fieldSize, c.fieldSize, c.fieldSize)
	case adagui.Vertical:
		r = geom.NewRectangleWH(float64(c.selIdx0)*c.fieldSize,
			float64(c.selIdx1)*c.fieldSize, c.fieldSize, c.fieldSize)
	}
	r = r.Inset(-3, -3)
	color := c.value.Get().(ledgrid.LedColor)
	gc.SetFillColor(color)
	gc.SetStrokeWidth(2.0)
	gc.SetStrokeColor(c.BorderColor())
	gc.DrawRectangle(r.AsCoord())
	gc.FillStroke()
}

//-----------------------------------------------------------------------

var (
	GridProps = props.NewProperties(props.PropsMap["Default"])
)

func init() {
	GridProps.SetColor(props.Color, colornames.Black)
	GridProps.SetColor(props.BorderColor, colornames.WhiteSmoke)
	GridProps.SetColor(props.LineColor, colornames.WhiteSmoke)
	GridProps.SetSize(props.BorderWidth, 2.0)
	GridProps.SetSize(props.LineWidth, 1.0)
}

type LedGrid struct {
	adagui.ContainerEmbed
	pixelSize float64
	DrawColor ledgrid.LedColor
	quitQ     chan bool
	grid      *ledgrid.LedGrid
	ctrl      *ledgrid.PixelClient
}

func NewLedGrid(pixelSize float64, host string, port uint) *LedGrid {
	g := &LedGrid{}
	g.Wrapper = g
	g.Init()
	g.PropertyEmbed.Init(GridProps)
	g.SetMinSize(geom.Point{10 * pixelSize, 10 * pixelSize})
	g.pixelSize = pixelSize
	g.DrawColor = ledgrid.LedColor{0x00, 0x00, 0x00, 0xFF}
	g.quitQ = make(chan bool)
	g.grid = ledgrid.NewLedGrid(image.Rect(0, 0, 10, 10))
	g.ctrl = ledgrid.NewPixelClient(host, port)
	return g
}

func (g *LedGrid) OnInputEvent(evt touch.Event) {
	//log.Printf("LedGrid: %v", evt.Pos)
	//log.Printf("    Pos: %v", g.Pos())
	//log.Printf("    Pos: %v", g.Pos())
	switch evt.Type {
	case touch.TypeDrag:
		//if !evt.Pos.In(g.Rect()) {
		//	break
		//}
		fx, fy := g.Bounds().PosRel(evt.Pos)
		col, row := int(10*fx), int(10*fy)
		oldColor := g.grid.LedColorAt(col, row)
		newColor := g.DrawColor.Mix(oldColor, ledgrid.Blend)
		g.grid.SetLedColor(col, row, newColor)
		g.Mark(adagui.MarkNeedsPaint)
	}
}

func (g *LedGrid) SetSize(size geom.Point) {
	s := max(size.X, size.Y)
	g.pixelSize = s / 10.0
	size.X, size.Y = s, s
	g.ContainerEmbed.SetSize(size)
}

func (g *LedGrid) Clear(c ledgrid.LedColor) {
	for idx := 0; idx < len(g.grid.Pix); idx += 3 {
		g.grid.Pix[idx+0] = c.R
		g.grid.Pix[idx+1] = c.G
		g.grid.Pix[idx+2] = c.B
	}
	g.Mark(adagui.MarkNeedsPaint)
}

func (g *LedGrid) Paint(gc *gg.Context) {
	for row := 0; row < 10; row++ {
		y0 := float64(row) * g.pixelSize
		for col := 0; col < 10; col++ {
			x0 := float64(col) * g.pixelSize
			c := g.grid.At(col, row)
			gc.SetFillColor(c)
			gc.DrawRectangle(x0, y0, pixelSize, pixelSize)
			gc.Fill()
		}
	}

	gc.SetStrokeWidth(g.LineWidth())
	gc.SetStrokeColor(g.LineColor())
	for t := g.pixelSize; t < g.Size().Y; t += g.pixelSize {
		gc.DrawLine(0.0, t, g.Size().Y, t)
		gc.DrawLine(t, 0.0, t, g.Size().Y)
	}
	gc.Stroke()

	gc.SetStrokeWidth(g.BorderWidth())
	gc.SetStrokeColor(g.BorderColor())
	gc.DrawRectangle(g.Bounds().AsCoord())
	gc.Stroke()

	g.ctrl.Draw(g.grid)
}

func (g *LedGrid) Save(fileName string) {
	fh, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	if err := png.Encode(fh, g.grid); err != nil {
		fh.Close()
		log.Fatal(err)
	}
	if err := fh.Close(); err != nil {
		log.Fatal(err)
	}
}

func (g *LedGrid) Load(fileName string) {
	fh, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	img, err := png.Decode(fh)
	if err != nil {
		fh.Close()
		log.Fatal(err)
	}
	draw.Draw(g.grid, img.Bounds(), img, image.Point{}, draw.Over)
	g.Mark(adagui.MarkNeedsPaint)
}

//-----------------------------------------------------------------------

const (
	host           = "raspi-2"
	port           = 5333
	pixelSize      = 31.0
	colorFieldSize = 25.0
)

func f1(t float64) float64 {
	return 3*t*t - 2*t*t*t
}

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lmsgprefix)
	log.SetPrefix(": ")
}

func SignalHandler() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	screen.Quit()
}

var (
	screen *adagui.Screen
	win    *adagui.Window
)

func main() {
	flag.Parse()
	adagui.StartProfiling()

	screen = adagui.NewScreen(adatft.Rotate000)
	win = screen.NewWindow()

	root := adagui.NewGroup()
	root.Layout = adagui.NewPaddedLayout()

	mainGrp := adagui.NewGroupPL(root, adagui.NewVBoxLayout(10))
	//drawGrp := adagui.NewGroupPL(mainGrp, adagui.NewHBoxLayout())

	drawColor := ledgrid.LedColor{}
	colorValue := binding.BindUntyped(&drawColor)
	colorPicker := NewColorPickerWithData(colorFieldSize, adagui.Horizontal,
		2, colorValue)

	ledPanel := NewLedGrid(pixelSize, host, port)

	colorPicker.SetColors(0, ledgrid.Pico08Colors[:8])
	colorPicker.SetColors(1, ledgrid.Pico08Colors[8:16])
	colorValue.AddCallback(func(data binding.DataItem) {
		if data == nil {
			return
		}
		ledPanel.DrawColor = data.(binding.Untyped).Get().(ledgrid.LedColor)
	})
	colorPicker.SetColor(ledgrid.Pico08Colors[0])
	mainGrp.Add(ledPanel)
	mainGrp.Add(colorPicker)

	toolGrp := adagui.NewGroupPL(mainGrp, adagui.NewHBoxLayout())
	btnClear := adagui.NewTextButton("Clear")
	btnClear.SetOnTap(func(evt touch.Event) {
		ledPanel.Clear(ledgrid.LedColor{0, 0, 0, 0xff})
	})
	toolGrp.Add(btnClear)

	mainGrp.Add(adagui.NewSpacer())

	buttonGrp := adagui.NewGroupPL(mainGrp, adagui.NewHBoxLayout())
	btnQuit := adagui.NewTextButton("Quit")
	btnQuit.SetOnTap(func(evt touch.Event) {
		ledPanel.Clear(ledgrid.LedColor{0, 0, 0, 0xff})
		screen.Quit()
	})
	btnSave := adagui.NewTextButton("Save")
	btnSave.SetOnTap(func(evt touch.Event) {
		ledPanel.Save("icon.png")
	})
	btnLoad := adagui.NewTextButton("Load")
	btnLoad.SetOnTap(func(evt touch.Event) {
		ledPanel.Load("icon.png")
	})
	buttonGrp.Add(btnSave, btnLoad, adagui.NewSpacer(), btnQuit)

	win.SetRoot(root)
	screen.SetWindow(win)
	screen.Run()

	adagui.StopProfiling()
}
