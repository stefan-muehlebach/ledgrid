package ledgrid

import (
	"image"
	"math"
	"time"

	"github.com/stefan-muehlebach/gg/color"

	"github.com/stefan-muehlebach/gg/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

//----------------------------------------------------------------------------

var (
	textFont     = fonts.GoBold
	textFontSize = 11.0
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
	pattern               image.Image
	drawer                font.Drawer
	params                []*Bounded[float64]
	anim                  *NormAnimation
}

func NewText(lg *LedGrid, txt string, col color.Color) *Text {
	t := &Text{}
	t.VisualEmbed.Init("Scrolling Text")
	t.lg = lg

	t.params = make([]*Bounded[float64], 2)
	t.params[0] = NewBounded("Font Size", textFontSize, 5.0, 15.0, 1.0)
	t.params[1] = NewBounded("Baseline", 9.0, 0.0, 20.0, 0.5)

	t.txt = txt
	t.startPos = coord2fix(float64(lg.Bounds().Max.X)+1.0, t.params[1].Val())
	t.endPos = t.startPos
	t.pos = t.startPos
	t.pattern = image.NewUniform(col)

	t.drawer = font.Drawer{
		Dst: lg,
		Src: t.pattern,
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

	t.anim = NewNormAnimation(10*time.Second, t.Update)
	t.anim.AutoReverse = true
	t.anim.RepeatCount = AnimationRepeatForever

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
	fixP := float2fix(p)
	fixQ := float2fix(1 - p)
	t.pos = t.startPos.Mul(fixQ).Add(t.endPos.Mul(fixP))
}

func (t *Text) Draw() {
	t.drawer.Dot = t.pos
	t.drawer.DrawString(t.txt)
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
