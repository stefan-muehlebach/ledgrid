package ledgrid

import (
	"image/color"
	"testing"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/fonts"
)

const (
    NumCols          = 10
    NumRows          = 15
	FieldWidth       = 100
	FieldHeight      = 100
	CornerRadius     = 20
	Margin           = 10.0
	ColorPaddingHori = 40.0
	ColorPaddingVert = 60.0
	FieldPadding     = 0.0
	FontSize         = 24.0
)

var (
	Font            = fonts.GoBold
	PaletteFileList = []string{
		"colormaps.json",
		"palettes.json",
	}
)

func TestReadJsonData(t *testing.T) {
	for _, fileName := range PaletteFileList {
		t.Logf("Reading palette file '%s'", fileName)
		ReadJsonData(fileName)
	}
}

func TestPaletteSamples(t *testing.T) {
	face := fonts.NewFace(Font, FontSize)
	gc := gg.NewContext(2*Margin+NumCols*(5*FieldWidth+4*FieldPadding)+(NumCols-1)*ColorPaddingHori,
		2*Margin+NumRows*FieldHeight+NumRows*ColorPaddingVert)
	gc.SetFontFace(face)
	gc.SetFillColor(color.White)
	gc.Clear()
	for i, name := range PaletteNames {
		col, row := i/NumRows, i%NumRows
		pal := PaletteMap[name]
		x := Margin + float64(col)*(5*FieldWidth+4*FieldPadding+ColorPaddingHori)
		y := Margin + float64(row)*(FieldHeight+ColorPaddingVert)
		for n := range 501 {
			t := float64(n) / float64(500)
			x := x + float64(n)
			gc.SetFillColor(pal.Color(t))
			gc.SetStrokeWidth(0.0)
			gc.DrawRectangle(x, y, 1.0, FieldHeight)
			gc.Fill()
		}
		gc.SetStrokeColor(color.Black)
		gc.DrawStringAnchored(pal.Name(), x, y+FieldHeight+Margin/2, 0.0, 1.0)
	}
	gc.SavePNG("palette_test.png")
}
