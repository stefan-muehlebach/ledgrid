package ledgrid

import (
	"strings"
	"image"
	"image/color"
	"image/draw"
	"sync"
	"time"

	"github.com/stefan-muehlebach/gg/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Es gibt zwei brauchbare TrueType Rasterer in Go: einmal die nativ Variante
// im Package golang.org/x/image/font und einmal die FreeType-Variante im
// Package github.com/golang/freetype/truetype - fuer beide findet sich in
// diesem File eine Umsetzung.

//----------------------------------------------------------------------------

var (
	textImgFactor      = 20
	textImgFactorFloat = float64(textImgFactor)
	textBaseFontSize   = 11.0
	textFontSize       = textImgFactorFloat * textBaseFontSize
	textFont           = fonts.SeafordBold
	transpPattern      = image.NewUniform(color.Transparent)
	wps                = 3
)

//----------------------------------------------------------------------------

type TextNative struct {
	VisualEmbed
	lg                    *LedGrid
	txt                   string
	size                  fixed.Point26_6
	startPos, endPos, pos fixed.Point26_6
	pal                   PaletteParameter
	img                   draw.Image
	drawer                *font.Drawer
	params                []Parameter
	anim                  *NormAnimation
	mutex                 sync.Mutex
}

func NewTextNative(lg *LedGrid, txt string, pal ColorSource) *TextNative {
	t := &TextNative{}
	t.VisualEmbed.Init("Scrolling Text (Go builtin)")
	t.lg = lg
	t.img = image.NewNRGBA(image.Rectangle{
		Min: lg.Rect.Min.Mul(textImgFactor),
		Max: lg.Rect.Max.Mul(textImgFactor),
	})

	t.params = make([]Parameter, 3)
	t.params[0] = NewFloatParameter("Font Size", textBaseFontSize, textBaseFontSize/2.0, 2.0*textBaseFontSize, 0.1)
	t.params[1] = NewFloatParameter("Baseline", float64(lg.Bounds().Max.Y-1), 0.0, float64(2*lg.Bounds().Max.Y), 0.1)
	t.params[2] = NewStringParameter("Message", txt)

	t.startPos = coord2fix(textImgFactorFloat*(float64(lg.Bounds().Max.X)+1.0), textImgFactorFloat*(t.params[1].(FloatParameter).Val()))
	t.endPos = t.startPos
	t.pos = t.startPos
	t.pal = NewPaletteParameter("Color", NewPaletteFader(pal))
	t.drawer = &font.Drawer{
		Dst: t.img,
	}
    t.txt = txt

	t.params[0].SetCallback(func(p Parameter) {
		v := t.params[0].(FloatParameter).Val()
		t.updateSize(v)
	})

	t.params[1].SetCallback(func(p Parameter) {
		v := textImgFactorFloat * t.params[1].(FloatParameter).Val()
		t.startPos.Y = float2fix(v)
		t.endPos.Y = float2fix(v)
	})

	t.params[2].SetCallback(func(p Parameter) {
		v := t.params[2].(StringParameter).Val()
		t.SetString(v)
	})

	t.anim = NewNormAnimation(time.Second, t.Update)
	t.anim.RepeatCount = AnimationRepeatForever
    t.anim.Curve = LinearAnimationCurve

    t.SetString(txt)

	return t
}

func (t *TextNative) updateSize(fontSize float64) {
	face := fonts.NewFace(textFont, textImgFactorFloat*fontSize)
	rect, _ := font.BoundString(face, t.txt)
	t.mutex.Lock()
	t.size = fixed.Point26_6{rect.Max.X - rect.Min.X, rect.Max.Y - rect.Min.Y}
	t.drawer.Face = face
	t.endPos.X = -(t.size.X + float2fix(1.0))
	t.mutex.Unlock()
}

func (t *TextNative) ParamList() []Parameter {
	return t.params
}

func (t *TextNative) PaletteParam() PaletteParameter {
	return t.pal
}

func (t *TextNative) Palette() ColorSource {
	return t.pal.Val()
}

func (t *TextNative) SetPalette(pal ColorSource, fadeTime time.Duration) {
	t.pal.Val().(*PaletteFader).StartFade(pal, fadeTime)
}

func (t *TextNative) String() string {
	return t.txt
}

func (t *TextNative) SetString(txt string) {
	t.txt = txt
    numWords := len(strings.Fields(t.txt))
	t.anim.Duration = time.Duration(numWords) * (2 * time.Second)
	t.updateSize(t.params[0].(FloatParameter).Val())
}

func (t *TextNative) SetVisible(vis bool) {
	if vis {
		t.anim.Start()
	} else {
		t.anim.Stop()
	}
	t.VisualEmbed.SetVisible(vis)
}

func (t *TextNative) Update(p float64) {
	draw.Draw(t.img, t.img.Bounds(), transpPattern, image.Point{}, draw.Src)
	t.mutex.Lock()
	x0, y0 := fix2coord(t.startPos)
	x1, y1 := fix2coord(t.endPos)
	x, y := (1-p)*x0+p*x1, (1-p)*y0+p*y1
	t.drawer.Dot = coord2fix(x, y)
	t.drawer.Src = image.NewUniform(t.pal.Val().Color(0))
	t.drawer.DrawString(t.txt)
	t.mutex.Unlock()
}

func (t *TextNative) ColorModel() color.Model {
	return LedColorModel
}

func (t *TextNative) Bounds() image.Rectangle {
	return t.img.Bounds()
}

func (t *TextNative) At(x, y int) color.Color {
	return t.img.At(x, y)
}
