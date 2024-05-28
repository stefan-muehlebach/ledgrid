package ledgrid

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"sync"
	"time"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/gobold"

	"github.com/golang/freetype"

	"github.com/stefan-muehlebach/gg/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

//----------------------------------------------------------------------------

var (
	textImgFactor      = 20
	textImgFactorFloat = float64(textImgFactor)
	textBaseFontSize   = 11.0
	textFontSize       = textImgFactorFloat * textBaseFontSize
	textFont           = fonts.GoBold
	freeFont, _        = truetype.Parse(gobold.TTF)
	transpPattern      = NewUniform(Transparent)
)

//----------------------------------------------------------------------------

type TextNative struct {
	VisualEmbed
	lg                    *LedGrid
	txt                   string
	size                  fixed.Point26_6
	startPos, endPos, pos fixed.Point26_6
	pal                   *PaletteFader
	img                   draw.Image
	drawer                *font.Drawer
	params                []*Bounded[float64]
	anim                  Animation
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

	t.params = make([]*Bounded[float64], 2)
	t.params[0] = NewBounded("Font Size", textBaseFontSize, textBaseFontSize/2.0, 2.0*textBaseFontSize, 1.0)
	t.params[1] = NewBounded("Baseline", float64(lg.Bounds().Max.Y), 0.0, float64(2*lg.Bounds().Max.Y), 1.0)

	t.txt = txt
	t.startPos = coord2fix(textImgFactorFloat*(float64(lg.Bounds().Max.X)+1.0), textImgFactorFloat*t.params[1].Val())
	t.endPos = t.startPos
	t.pos = t.startPos
	t.pal = NewPaletteFader(pal)
	t.drawer = &font.Drawer{
		Dst: t.img,
	}

	t.params[0].SetCallback(func(oldVal, newVal float64) {
		t.updateSize(newVal)
	})

	t.params[1].SetCallback(func(oldVal, newVal float64) {
		t.startPos.Y = float2fix(textImgFactorFloat * newVal)
		t.endPos.Y = float2fix(textImgFactorFloat * newVal)
	})

	anim := NewNormAnimation(10*time.Second, t.Update)
	anim.AutoReverse = true
	anim.RepeatCount = AnimationRepeatForever
	anim.Curve = CubicAnimationCurve
	t.anim = anim

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

func (t *TextNative) ParamList() []*Bounded[float64] {
	return t.params
}

func (t *TextNative) Palette() ColorSource {
	return t.pal
}

func (t *TextNative) SetPalette(pal ColorSource, fadeTime time.Duration) {
	t.pal.StartFade(pal, fadeTime)
}

func (t *TextNative) String() string {
	return t.txt
}

func (t *TextNative) SetString(txt string) {
	t.txt = txt
	t.updateSize(t.params[0].Val())
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
	t.drawer.Src = image.NewUniform(t.pal.Color(0))
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

// Dieser Typ macht grundsaetzlich genau das gleiche wie der oben gezeigte,
// ausser, dass er fuer die Rasterung die FreeType-Library verwendet.
type TextFreeType struct {
	VisualEmbed
	lg                    *LedGrid
	txt                   string
	size                  fixed.Point26_6
	startPos, endPos, pos fixed.Point26_6
	pal                   *PaletteFader
	img                   draw.Image
	ctx                   *freetype.Context
	params                []*Bounded[float64]
	anim                  Animation
	mutex                 sync.Mutex
}

func NewTextFreeType(lg *LedGrid, txt string, pal ColorSource) *TextFreeType {
	t := &TextFreeType{}
	t.VisualEmbed.Init("Scrolling Text (FreeType)")
	t.lg = lg
	t.img = image.NewNRGBA(image.Rectangle{
		Min: lg.Rect.Min.Mul(textImgFactor),
		Max: lg.Rect.Max.Mul(textImgFactor),
	})
	t.ctx = freetype.NewContext()

	t.params = make([]*Bounded[float64], 2)
	t.params[0] = NewBounded("Font Size", textBaseFontSize, textBaseFontSize/2.0, 2.0*textBaseFontSize, 1.0)
	t.params[1] = NewBounded("Baseline", float64(lg.Bounds().Max.Y), 0.0, float64(2*lg.Bounds().Max.Y), 1.0)

	t.txt = txt
	t.startPos = coord2fix(textImgFactorFloat*(float64(lg.Bounds().Max.X)+1.0), textImgFactorFloat*t.params[1].Val())
	t.endPos = t.startPos
	t.pos = t.startPos
	t.pal = NewPaletteFader(pal)

	t.params[0].SetCallback(func(oldVal, newVal float64) {
		t.updateSize(newVal)
	})

	t.params[1].SetCallback(func(oldVal, newVal float64) {
		t.startPos.Y = float2fix(textImgFactorFloat * newVal)
		t.endPos.Y = float2fix(textImgFactorFloat * newVal)
	})

	t.ctx.SetDst(t.img)
	t.ctx.SetFont(freeFont)
	t.ctx.SetDPI(72.0)
	t.ctx.SetHinting(font.HintingNone)

	anim := NewNormAnimation(10*time.Second, t.Update)
	anim.AutoReverse = true
	anim.RepeatCount = AnimationRepeatForever
	anim.Curve = CubicAnimationCurve
	t.anim = anim

	return t
}

func (t *TextFreeType) updateSize(fontSize float64) {
	face := fonts.NewFace(textFont, textImgFactorFloat*fontSize)
	rect, _ := font.BoundString(face, t.txt)
	t.mutex.Lock()
	t.size = rect.Max.Sub(rect.Min)
	t.ctx.SetFontSize(textImgFactorFloat * fontSize)
	t.ctx.SetClip(t.img.Bounds().Inset(int(fix2float(-max(t.size.X, t.size.Y)))))
	t.endPos.X = -(t.size.X + float2fix(1.0))
	t.mutex.Unlock()
}

func (t *TextFreeType) ParamList() []*Bounded[float64] {
	return t.params
}

func (t *TextFreeType) Palette() ColorSource {
	return t.pal
}

func (t *TextFreeType) SetPalette(pal ColorSource, fadeTime time.Duration) {
	t.pal.StartFade(pal, fadeTime)
}

func (t *TextFreeType) String() string {
	return t.txt
}

func (t *TextFreeType) SetString(txt string) {
	t.txt = txt
	t.updateSize(t.params[0].Val())
}

func (t *TextFreeType) SetVisible(vis bool) {
	if vis {
		t.anim.Start()
	} else {
		t.anim.Stop()
	}
	t.VisualEmbed.SetVisible(vis)
}

func (t *TextFreeType) Update(p float64) {
	draw.Draw(t.img, t.img.Bounds(), transpPattern, image.Point{}, draw.Src)
	t.mutex.Lock()
	x0, y0 := fix2coord(t.startPos)
	x1, y1 := fix2coord(t.endPos)
	x, y := (1-p)*x0+p*x1, (1-p)*y0+p*y1
	t.pos = coord2fix(x, y)
	t.ctx.SetSrc(image.NewUniform(t.pal.Color(0)))
	_, err := t.ctx.DrawString(t.txt, t.pos)
	t.mutex.Unlock()
	if err != nil {
		log.Fatalf("DrawString: %v", err)
	}
}

func (t *TextFreeType) ColorModel() color.Model {
	return LedColorModel
}

func (t *TextFreeType) Bounds() image.Rectangle {
	return t.img.Bounds()
}

func (t *TextFreeType) At(x, y int) color.Color {
	return t.img.At(x, y)
}

// Es folgen Hilfsfunktionen fuer die schnelle Umrechnung zwischen Fliess-
// und Fixkommazahlen sowie den geometrischen Typen, die auf Fixkommazahlen
// aufgebaut sind.
func rect2fix(r image.Rectangle) fixed.Rectangle26_6 {
	return fixed.R(r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)
}

func coord2fix(x, y float64) fixed.Point26_6 {
	return fixed.Point26_6{X: float2fix(x), Y: float2fix(y)}
}

func fix2coord(p fixed.Point26_6) (x, y float64) {
	return fix2float(p.X), fix2float(p.Y)
}

func float2fix(x float64) fixed.Int26_6 {
	return fixed.Int26_6(math.Round(x * 64))
}

func fix2float(x fixed.Int26_6) float64 {
	return float64(x) / 64.0
}
