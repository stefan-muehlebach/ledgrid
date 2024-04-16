package ledgrid

import (
	"image"
	"math"
	"time"

	"github.com/stefan-muehlebach/gg/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

//----------------------------------------------------------------------------

var (
	textFont = fonts.GoBold
    textFontSize = 11.0
	// textFont = fonts.LucidaConsole
	// textFont = fonts.LucidaSansTypewriterBold
    // textFont = fonts.LucidaBright
)

//----------------------------------------------------------------------------

type Text struct {
	VisualizableEmbed
	lg      *LedGrid
	txt     string
	size    fixed.Point26_6
	pos, dp fixed.Point26_6
	pattern image.Image
	drawer  font.Drawer
    	params  []*Bounded[float64]

}

func NewText(lg *LedGrid, txt string, col LedColor) *Text {
	face := fonts.NewFace(textFont, textFontSize)
	t := &Text{}
	t.VisualizableEmbed.Init("Text")
	t.lg = lg
	t.txt = txt
	rect, _ := font.BoundString(face, t.txt)
	t.size = fixed.Point26_6{rect.Max.X - rect.Min.X, rect.Max.Y - rect.Min.Y}
	t.pos = coord2fix(10.0, 9.0)
	t.dp = coord2fix(-0.1, 0.0)
	t.pattern = image.NewUniform(col)
	t.drawer = font.Drawer{
		Dst:  lg,
		Src:  t.pattern,
		Face: face,
	}
    t.params = make([]*Bounded[float64], 2)
    t.params[0] = NewBounded("Font Size", 10.0, 5.0, 15.0, 0.2)
    t.params[0].SetCallback(func(oldVal, newVal float64) {
        face := fonts.NewFace(textFont, newVal)
        rect, _ := font.BoundString(face, t.txt)
        t.size = fixed.Point26_6{rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y}
        t.drawer.Face = face
    })
    t.params[1] = NewBounded("Y-Coordinate of the Baseline", 9.0, 0.0, 20.0, 0.5)
    t.params[1].SetCallback(func(oldVal, newVal float64) {
        t.pos.Y = float2fix(newVal)
    })
	return t
}

func (t *Text) ParamList() []*Bounded[float64] {
	return t.params
}

func (o *Text) Update(dt time.Duration) bool {
    dt = o.VisualizableEmbed.Update(dt)
	o.pos = o.pos.Add(o.dp)
	if o.pos.X+o.size.X < 0 ||
		o.pos.X > fixed.I(o.lg.Bounds().Dx()) {
		o.dp.X *= -1.0
	}
	if o.pos.Y < 0 ||
		o.pos.Y > o.size.Y+fixed.I(o.lg.Bounds().Dy()) {
		o.dp.Y *= -1.0
	}
	return true
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
