package ledgrid

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
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
	textImgFactor    = 4
	textBaseFontSize = 11.0
	textFontSize     = float64(textImgFactor) * textBaseFontSize
	textFont         = fonts.GoBold
	freeFont, _      = truetype.Parse(gobold.TTF)
	// textFont = fonts.LucidaConsole
	// textFont = fonts.LucidaSansTypewriterBold
	// textFont = fonts.LucidaBright
)

//----------------------------------------------------------------------------

type Text struct {
	VisualEmbed
	lg                    *LedGrid
	txt                   string
	size                  fixed.Point26_6
	startPos, endPos, pos fixed.Point26_6
	col                   color.NRGBA
	img                   draw.Image
	drawer                *font.Drawer
	params                []*Bounded[float64]
	anim                  Animation
}

func NewText(lg *LedGrid, txt string, col color.Color) *Text {
	t := &Text{}
	t.VisualEmbed.Init("Scrolling Text (Go builtin)")
	t.lg = lg
	t.img = image.NewNRGBA(image.Rectangle{
		Min: lg.Rect.Min.Mul(textImgFactor),
		Max: lg.Rect.Max.Mul(textImgFactor),
	})

	t.params = make([]*Bounded[float64], 2)
	t.params[0] = NewBounded("Font Size", textFontSize, textFontSize/2.0, 2.0*textFontSize, 1.0)
	t.params[1] = NewBounded("Baseline", float64(t.img.Bounds().Max.Y),
        0.0, float64(2*t.img.Bounds().Max.Y), 1.0)

	t.txt = txt
	t.startPos = coord2fix(float64(t.img.Bounds().Max.X)+1.0, t.params[1].Val())
	t.endPos = t.startPos
	t.pos = t.startPos
	t.col = color.NRGBAModel.Convert(col).(color.NRGBA)

	t.drawer = &font.Drawer{
		Dst: t.img,
		Src: image.NewUniform(t.col),
	}

	t.params[0].SetCallback(func(oldVal, newVal float64) {
		face := fonts.NewFace(textFont, newVal)
		rect, _ := font.BoundString(face, t.txt)
		t.size = fixed.Point26_6{rect.Max.X - rect.Min.X, rect.Max.Y - rect.Min.Y}
		t.drawer.Face = face
		t.endPos.X = -(t.size.X + float2fix(1.0))
	})

	t.params[1].SetCallback(func(oldVal, newVal float64) {
		t.startPos.Y = float2fix(newVal)
		t.endPos.Y = float2fix(newVal)
	})

	anim := NewNormAnimation(10*time.Second, t.Update)
	anim.AutoReverse = true
	anim.RepeatCount = AnimationRepeatForever
	t.anim = anim

	return t
}

func (t *Text) ParamList() []*Bounded[float64] {
	return t.params
}

func (t *Text) SetVisible(vis bool) {
	if vis {
		t.anim.Start()
	} else {
		t.anim.Stop()
	}
	t.VisualEmbed.SetVisible(vis)
}

func (t *Text) Update(p float64) {
	draw.Draw(t.img, t.img.Bounds(), NewUniform(Transparent), image.Point{}, draw.Src)
	fixP := float2fix(p)
	fixQ := float2fix(1 - p)
	t.drawer.Dot = t.startPos.Mul(fixQ).Add(t.endPos.Mul(fixP))
	t.drawer.DrawString(t.txt)
}

func (t *Text) ColorModel() color.Model {
	return LedColorModel
}

func (t *Text) Bounds() image.Rectangle {
	return t.img.Bounds()
}

func (t *Text) At(x, y int) color.Color {
	return t.img.At(x, y)
}

// Dieser Typ macht grundsaetzlich genau das gleiche wie der oben gezeigte,
// ausser, dass er fuer die Rasterung die FreeType-Library verwendet.
type TextFT struct {
	VisualEmbed
	lg                    *LedGrid
	txt                   string
	size                  fixed.Point26_6
	startPos, endPos, pos fixed.Point26_6
	col                   color.NRGBA
	img                   draw.Image
	ctx                   *freetype.Context
	params                []*Bounded[float64]
	anim                  Animation
}

func NewTextFT(lg *LedGrid, txt string, col color.Color) *TextFT {
	t := &TextFT{}
	t.VisualEmbed.Init("Scrolling Text (FreeType)")
	t.lg = lg
	t.img = image.NewNRGBA(image.Rectangle{
		Min: lg.Rect.Min.Mul(textImgFactor),
		Max: lg.Rect.Max.Mul(textImgFactor),
	})
	t.ctx = freetype.NewContext()

	t.params = make([]*Bounded[float64], 2)
	t.params[0] = NewBounded("Font Size", textFontSize, textFontSize/2.0, 2.0*textFontSize, 1.0)
	t.params[1] = NewBounded("Baseline", float64(t.img.Bounds().Max.Y),
        0.0, float64(2*t.img.Bounds().Max.Y), 1.0)

	t.txt = txt
	t.startPos = coord2fix(float64(t.img.Bounds().Max.X)+1.0, t.params[1].Val())
	t.endPos = t.startPos
	t.pos = t.startPos
	t.col = color.NRGBAModel.Convert(col).(color.NRGBA)

	t.params[0].SetCallback(func(oldVal, newVal float64) {
		face := fonts.NewFace(textFont, newVal)
		rect, _ := font.BoundString(face, t.txt)
		t.size = rect.Max.Sub(rect.Min)
		t.ctx.SetFontSize(newVal)
		t.ctx.SetClip(t.img.Bounds().Inset(int(fix2float(-max(t.size.X, t.size.Y)))))
		t.endPos.X = -(t.size.X + float2fix(1.0))
	})

	t.params[1].SetCallback(func(oldVal, newVal float64) {
		t.startPos.Y = float2fix(newVal)
		t.endPos.Y = float2fix(newVal)
	})

	t.ctx.SetSrc(image.NewUniform(t.col))
	t.ctx.SetDst(t.img)
	t.ctx.SetFont(freeFont)
	t.ctx.SetDPI(72.0)
	t.ctx.SetHinting(font.HintingNone)

	anim := NewNormAnimation(10*time.Second, t.Update)
	anim.AutoReverse = true
	anim.RepeatCount = AnimationRepeatForever
	t.anim = anim

	return t
}

func (t *TextFT) ParamList() []*Bounded[float64] {
	return t.params
}

func (t *TextFT) SetVisible(vis bool) {
	if vis {
		t.anim.Start()
	} else {
		t.anim.Stop()
	}
	t.VisualEmbed.SetVisible(vis)
}

func (t *TextFT) Update(p float64) {
	draw.Draw(t.img, t.img.Bounds(), NewUniform(Transparent), image.Point{}, draw.Src)
	fixP := float2fix(p)
	fixQ := float2fix(1 - p)
	t.pos = t.startPos.Mul(fixQ).Add(t.endPos.Mul(fixP))
	_, err := t.ctx.DrawString(t.txt, t.pos)
	if err != nil {
		log.Fatalf("DrawString: %v", err)
	}
}

func (t *TextFT) ColorModel() color.Model {
	return LedColorModel
}

func (t *TextFT) Bounds() image.Rectangle {
	return t.img.Bounds()
}

func (t *TextFT) At(x, y int) color.Color {
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

func float2fix(x float64) fixed.Int26_6 {
	return fixed.Int26_6(math.Round(x * 64))
}

func fix2float(x fixed.Int26_6) float64 {
	return float64(x) / 64.0
}
