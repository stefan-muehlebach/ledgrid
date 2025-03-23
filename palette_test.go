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
	TileRows         = 4
	TileCols         = 256 / TileRows
	TileWidth        = FieldWidth / TileCols
	TileHeight       = FieldHeight / TileRows
	CornerRadius     = 20
	Margin           = 10.0
	ColorPaddingHori = 20.0
	ColorPaddingVert = 20.0
	FontSize         = 24.0
)

var (
	NumCols         = 4
	NumRows         = 15
	Font            = fonts.GoBold
	PaletteFileList = []string{
		"palGradient.json",
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
		2*Margin+NumRows*FieldHeight+NumRows*(ColorPaddingVert+FontSize))
	gc.SetFontFace(face)
	gc.SetTextColor(color.Black)
	gc.SetFillColor(color.WhiteSmoke)
	gc.SetStrokeWidth(0.0)
	gc.Clear()
	for i, name := range PaletteNames {
		col, row := i%NumCols, i/NumCols
		x := Margin + float64(col)*(FieldWidth+ColorPaddingHori)
		y := Margin + float64(row)*(FieldHeight+ColorPaddingVert+FontSize)
		switch pal := PaletteMap[name].(type) {
		case *GradientPalette:
			for n := range FieldWidth + 1 {
				t := float64(n) / float64(FieldWidth)
				x := x + float64(n)
				gc.SetFillColor(pal.Color(t))
				gc.DrawRectangle(x, y, 1.0, FieldHeight)
				gc.Fill()
			}

		case *SlicePalette:
			t.Logf("Slice palette...")
			for j := range TileRows {
				yNew := y + float64(j)*TileHeight
				for k := range TileCols {
					xNew := x + float64(k)*TileWidth
					gc.SetFillColor(pal.Color(float64(j*TileCols + k)))
					gc.DrawRectangle(xNew, yNew, TileWidth, TileHeight)
					gc.Fill()
				}
			}
		}
		gc.DrawStringAnchored(name, x, y+FieldHeight+FontSize, 0.0, 0.0)
	}
	gc.SavePNG("data/Overview.png")
}
