package ledgrid

import (
	"testing"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/fonts"
)

const (
	FieldWidth       = 500
	FieldHeight      = 100
	CornerRadius     = 20
	Margin           = 10.0
	ColorPaddingHori = 40.0
	ColorPaddingVert = 40.0
	FontSize         = 24.0
)

var (
	NumCols         = 4
	NumRows         = 15
	Font            = fonts.GoBold
	PaletteFileList = []string{
		"palettes.json",
	}
)

func TestReadJsonData(t *testing.T) {
	for _, fileName := range PaletteFileList {
		t.Logf("Reading palette file '%s'", fileName)
		ReadJsonData(fileName)
	}
}

func TestPaletteOverview(t *testing.T) {
	NumRows = len(PaletteNames) / NumCols
	if len(PaletteNames)%NumCols != 0 {
		NumRows += 1
	}
	face := fonts.NewFace(Font, FontSize)
	gc := gg.NewContext(2*Margin+NumCols*(FieldWidth)+(NumCols-1)*ColorPaddingHori,
		2*Margin+NumRows*FieldHeight+NumRows*ColorPaddingVert)
	gc.SetFontFace(face)
	gc.SetTextColor(color.Black)
	gc.SetFillColor(color.WhiteSmoke)
	gc.Clear()
	for i, name := range PaletteNames {
		col, row := i%NumCols, i/NumCols
		pal := PaletteMap[name]
		x := Margin + float64(col)*(FieldWidth+ColorPaddingHori)
		y := Margin + float64(row)*(FieldHeight+ColorPaddingVert)
		for n := range 501 {
			t := float64(n) / float64(500)
			x := x + float64(n)
			gc.SetFillColor(pal.Color(t))
			gc.SetStrokeWidth(0.0)
			gc.DrawRectangle(x, y, 1.0, FieldHeight)
			gc.Fill()
		}
		gc.DrawStringAnchored(name, x, y+FieldHeight+Margin/2, 0.0, 1.0)
	}
	gc.SavePNG("data/palOverview.png")
}
