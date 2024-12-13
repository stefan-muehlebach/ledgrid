package conf

import (
	"fmt"
	"log"
	"math"

	"golang.org/x/image/font"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/fonts"
	"github.com/stefan-muehlebach/gg/geom"
)

// For the graphical output of the wiring diagram, many constants are used to
// configure the visual appearance.
var (
    // scaleFactor can be used to scale the whole plot. Use it, when tiny
    // detail are not visible anymore.
	scaleFactor = 1.0

	MarginLeft   = 30.0 * scaleFactor
	MarginRight  = 5.0 * scaleFactor
	MarginTop    = 30.0 * scaleFactor
	MarginBottom = 5.0 * scaleFactor

	AxesTextFont   = fonts.GoRegular
	AxesTextSize   = 12.0 * scaleFactor
	AxesTextColor  = color.Black
	AxesTickSep    = 5.0 * scaleFactor
	AxesTickHeight = 20.0 * scaleFactor
	AxesTickWidth  = 1.0 * scaleFactor
	AxesTickColor  = color.Black

	ModuleSize        = 400.0 * scaleFactor
	ModuleBorderWidth = 2.0 * scaleFactor
	ModuleBorderColor = color.Black
	ModuleFillColor   = color.AntiqueWhite
	ModuleTextFont    = fonts.GoBold
	ModuleTextSize    = 220.0 * scaleFactor
	ModuleTextColor   = color.DarkSlateGray

	LedFieldSize      = ModuleSize / float64(ModuleDim.X)
	LedSize           = LedFieldSize - 2.0
	LedBorderWidth    = 1.0 * scaleFactor
	LedBorderColor    = color.Black
	LedFillColor      = color.White.Alpha(0.7)
	LedStartFillColor = color.DarkGreen.Alpha(0.8)
	LedEndFillColor   = color.FireBrick.Alpha(0.8)
	LedTextFont       = fonts.GoRegular
	LedTextSize       = 16.0 * scaleFactor
	LedTextColor      = color.Black
	LedTextColorInv   = color.White

	// Trace is...
	TraceWidth     = 15.0 * scaleFactor
	TraceColor     = color.DarkSlateGray
	TraceBezierPos = 0.8

	sizeFact = 0.8

	axesFontFace, moduleFontFace font.Face
	ledFontFaces                 [7]font.Face

	pTL, pTR geom.Point
)

func init() {
	axesFontFace = fonts.NewFace(AxesTextFont, AxesTextSize)
	moduleFontFace = fonts.NewFace(ModuleTextFont, ModuleTextSize)
	size := LedTextSize
	for i := range ledFontFaces {
		if i >= 3 {
			size *= sizeFact
		}
		ledFontFaces[i] = fonts.NewFace(LedTextFont, size)
	}
	pTL = geom.Point{-ModuleSize / 2.0, -ModuleSize / 2.0}
	pTR = pTL.AddXY(ModuleSize, 0.0)
}

// With Plot, you can create a PNG file, showing the exact module config,
// the cabeling and the mapping between pixel coordinates and index on the
// LED chain.
func (conf ModuleConfig) Plot(fileName string) {
    err := conf.Verify()
	if err != nil {
        log.Fatalf("Cannot plot configuration: %v", err)
    }

	size := conf.Size()
	gc := gg.NewContext(size.X*int(LedFieldSize)+int(MarginLeft+MarginRight),
		size.Y*int(LedFieldSize)+int(MarginTop+MarginBottom))
	gc.SetFillColor(color.White)
	gc.Clear()

	conf.Draw(gc)

	err = gc.SavePNG(fileName)
	if err != nil {
		log.Fatalf("Couldn't save configuration: %v", err)
	}
}

// Draw is the workhorse of Plot. With Draw, you can create the configuration
// image into the graphical context, provided by the parameter gc.
func (conf ModuleConfig) Draw(gc *gg.Context) {

    conf.DrawAxes(gc)

	// And then draw the individual modules
	p0 := geom.Point{MarginLeft, MarginTop}.AddXY(ModuleSize/2.0, ModuleSize/2.0)
	for i, modPos := range conf {
		pt := p0.Add(geom.Point{float64(modPos.Col), float64(modPos.Row)}.Mul(ModuleSize))

		gc.Push()
		gc.Translate(pt.X, pt.Y)
		gc.Rotate(-gg.Radians(float64(modPos.Mod.Rot)))
		modPos.Mod.Draw(gc, i)
		gc.Pop()
	}
}

func (conf ModuleConfig) DrawAxes(gc *gg.Context) {
	// Label the columns and rows of the LEDs over the whole panel.
	p0 := geom.Point{MarginLeft, MarginTop}
	gc.SetStrokeWidth(AxesTickWidth)
    gc.SetStrokeColor(AxesTickColor)
	gc.SetTextColor(AxesTextColor)
	gc.SetFontFace(axesFontFace)

	p1 := p0.AddXY(0, -AxesTickSep)
	p2 := p1.AddXY(0, -AxesTickHeight)
	gc.DrawLine(p1.X, p1.Y, p2.X, p2.Y)
	gc.Stroke()
	for col := range conf.Size().X {
		p1 = p1.AddXY(float64(LedFieldSize), 0.0)
		p2 = p1.AddXY(0, -AxesTickHeight)
		gc.DrawLine(p1.X, p1.Y, p2.X, p2.Y)
		gc.Stroke()

		p3 := p1.Interpolate(p2, 0.5).AddXY(-LedFieldSize/2.0, 0.0)
		gc.DrawStringAnchored(fmt.Sprintf("%d", col), p3.X, p3.Y, 0.5, 0.5)

	}
	p1 = p0.AddXY(-AxesTickSep, 0)
	p2 = p1.AddXY(-AxesTickHeight, 0)
	gc.DrawLine(p1.X, p1.Y, p2.X, p2.Y)
	gc.Stroke()
	for row := range conf.Size().Y {
		p1 = p1.AddXY(0.0, float64(LedFieldSize))
		p2 = p1.AddXY(-AxesTickHeight, 0)
		gc.DrawLine(p1.X, p1.Y, p2.X, p2.Y)
		gc.Stroke()

		p3 := p1.Interpolate(p2, 0.5).AddXY(0.0, -LedFieldSize/2.0)
		gc.DrawStringAnchored(fmt.Sprintf("%d", row), p3.X, p3.Y, 0.5, 0.5)
	}
}

// This method draws a single module. The calling method/function must ensure,
// using translation and possibly rotation, that the origin of the coordinate
// system is positioned at the center of the module.
func (mod Module) Draw(gc *gg.Context, idxMod int) {
	// p0 und p1 are reference points which are placed on the top left and top
	// right of the module.
	// p0 := geom.Point{-ModuleSize / 2.0, -ModuleSize / 2.0}
	// p1 := p0.Add(geom.Point{ModuleSize, 0})

	// Draws the filling of the module - the border will be drawn later.
	gc.DrawRectangle(pTL.X, pTL.Y, ModuleSize, ModuleSize)
	gc.SetFillColor(ModuleFillColor)
	gc.Fill()

	// Draw the module type name in huge letters in the middle of the module.
	//gc.Push()
	//gc.Rotate(gg.Radians(float64(mod.Rot)))
	gc.SetFontFace(moduleFontFace)
	gc.SetTextColor(ModuleTextColor)
	gc.DrawStringAnchored(fmt.Sprintf("%v", mod.Type), 0.0, 0.0, 0.5, 0.5)
	//gc.Pop()

	mod.DrawTrace(gc)

	// Draw the individual LEDs or table tennis balls, respectively.
	// Set first some orientation dependent variables:
	// mp0: the midpoint of the first LED to be drawn.
	// dx:  step in x-direction between two LEDs
	// dy:  step in y-direction between two LEDs
	// turn: angle (in radians) to turn at the end of a column
	mp0 := pTL.AddXY(LedFieldSize/2.0, LedFieldSize/2.0)
	dx := geom.Point{LedFieldSize, 0.0}
	dy := geom.Point{0.0, LedFieldSize}
	turn := gg.Radians(-90.0)
	dir := math.Pi
	if mod.Type == ModRL {
		mp0 = pTR.AddXY(-LedFieldSize/2.0, LedFieldSize/2.0)
		dx = geom.Point{-LedFieldSize, 0.0}
		dy = geom.Point{0.0, LedFieldSize}
		turn = gg.Radians(90.0)
	}

	// mp is the midpoint of the next LED to be drawn.
	mp := mp0
	// mode denotes the part of the led chain we are currently drawing.
	// Values of mode are
	// 0: on an even column, next led will be drawn under the current led
	// 1: at the buttom of an even column
	// 2: on a odd column, next led will be drawn on top of the current led
	// 3: at the top of a odd column
	mode := 0
	// dp is the step to take in order to get to the next LED position.
	dp := dy

	// In this loop, we iterate over all of the 100 LEDs of the module.
	for idx := range ModuleDim.X * ModuleDim.Y {
		idxLed := 100*idxMod + idx
		switch mode {
		case 0:
			if idx%ModuleDim.Y == 9 {
				dir += turn
				dp = dx
				mode++
			}
		case 1:
			dir += turn
			dp = dy.Neg()
			mode++
		case 2:
			if idx%ModuleDim.Y == 9 {
				dir -= turn
				dp = dx
				mode++
			}
		case 3:
			dir -= turn
			dp = dy
			mode = 0
		}

		// Draw the LED (or table tennis ball) as a filled circle.
		gc.DrawCircle(mp.X, mp.Y, LedSize/2.0)
		gc.SetStrokeWidth(LedBorderWidth)
		gc.SetStrokeColor(LedBorderColor)
		if idx == 0 {
			gc.SetFillColor(LedStartFillColor)
		} else if idx == ModuleDim.X*ModuleDim.Y-1 {
			gc.SetFillColor(LedEndFillColor)
		} else {
			gc.SetFillColor(LedFillColor)
		}
		gc.FillStroke()

		// Label the LED with the index number of this LED within the whole
		// LED chain.
		gc.Push()
		gc.Translate(mp.X, mp.Y)
		gc.Rotate(gg.Radians(float64(mod.Rot)))
		i := int(math.Log10(float64(max(idxLed, 1))))
		gc.SetFontFace(ledFontFaces[i])
		if idx > 0 && idx < ModuleDim.X*ModuleDim.Y-1 {
			gc.SetTextColor(LedTextColor)
		} else {
			gc.SetTextColor(LedTextColorInv)
		}
		gc.DrawStringAnchored(fmt.Sprintf("%d", idxLed), 0.0, 0.0, 0.5, 0.5)
		gc.Pop()

		mp = mp.Add(dp)
	}

    mod.DrawBorder(gc)
}

func (mod Module) DrawTrace(gc *gg.Context) {

	// Verlauf der Verkabelung als graues, maeandrierendes Band darstellen.
	var mpA, mpB, mpANew, mpBNew, mpC1, mpC2 geom.Point

	// p0 := geom.Point{-ModuleSize / 2.0, -ModuleSize / 2.0}
	// p1 := p0.Add(geom.Point{ModuleSize, 0})
	mp0 := pTL.AddXY(LedFieldSize/2.0, LedFieldSize/2.0)
	dx := geom.Point{LedFieldSize, 0.0}
	dy := geom.Point{0.0, LedFieldSize}
	if mod.Type == ModRL {
		mp0 = pTR.AddXY(-LedFieldSize/2.0, LedFieldSize/2.0)
		dx = geom.Point{-LedFieldSize, 0.0}
		dy = geom.Point{0.0, LedFieldSize}
	}

	mpA = mp0.Add(dy.Div(2))
	mpB = mpA.Add(dy.Mul(8))
	gc.SetStrokeWidth(TraceWidth)
	gc.SetStrokeColor(TraceColor)
	gc.MoveTo(mp0.AsCoord())
	gc.LineTo(mpA.AsCoord())
	for i := range 10 {
		if i%2 == 0 {
			gc.LineTo(mpB.AsCoord())
		} else {
			gc.LineTo(mpA.AsCoord())
		}
		mpANew = mpA.Add(dx)
		mpBNew = mpB.Add(dx)
		if i < 9 {
			if i%2 == 0 {
				mpC1 = mpB.Add(dy.Mul(TraceBezierPos))
				mpC2 = mpBNew.Add(dy.Mul(TraceBezierPos))
				gc.CubicTo(mpC1.X, mpC1.Y, mpC2.X, mpC2.Y, mpBNew.X, mpBNew.Y)
			} else {
				mpC1 = mpA.Sub(dy.Mul(TraceBezierPos))
				mpC2 = mpANew.Sub(dy.Mul(TraceBezierPos))
				gc.CubicTo(mpC1.X, mpC1.Y, mpC2.X, mpC2.Y, mpANew.X, mpANew.Y)
			}
		}
		mpA = mpANew
		mpB = mpBNew
	}
	gc.LineTo(mp0.Add(dx.Mul(9)).AsCoord())
	gc.Stroke()
}

func (mod Module) DrawBorder(gc *gg.Context) {
    	// Draw the border of the module.
	gc.DrawRectangle(pTL.X, pTL.Y, ModuleSize, ModuleSize)
	gc.SetStrokeWidth(ModuleBorderWidth)
	gc.SetStrokeColor(ModuleBorderColor)
	gc.Stroke()
}

// Hilfsfunktioenchen (sogar generisch!)
func abs[T ~int | ~float64](i T) T {
	if i < 0 {
		return -i
	} else {
		return i
	}
}
