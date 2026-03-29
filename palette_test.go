package ledgrid

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/fonts"
)

const (
	FieldWidth       = 600
	FieldHeight      = 100
	TextHeight       = 30.0
	TileRows         = 4
	TileCols         = 256 / TileRows
	TileWidth        = FieldWidth / TileCols
	TileHeight       = FieldHeight / TileRows
	CornerRadius     = 20
	Margin           = 10.0
	ColorPaddingHori = 20.0
	ColorPaddingVert = 20.0
	FontSize         = 18.0
	FontSizeSmall    = 14.0
)

var (
	NumCols  = 7
	NumRows  = 0
	NameFont = fonts.LucidaBrightDemibold
	TypeFont = fonts.LucidaBrightItalic
)

func TestPalette(t *testing.T) {

	fh, err := os.Open(path.Join("data", "palNew.json"))
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer fh.Close()

	palNames, palMap, err := colors.ReadPaletteData(fh)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	NumRows = len(palNames) / NumCols
	if len(palNames)%NumCols != 0 {
		NumRows += 1
	}
	nameFace, _ := fonts.NewFace(NameFont, FontSize)
	typeFace, _ := fonts.NewFace(TypeFont, FontSizeSmall)
	gc := gg.NewContext(NumCols*(FieldWidth)+(NumCols-1)*ColorPaddingHori,
		NumRows*FieldHeight+NumRows*(ColorPaddingVert+TextHeight))
	gc.SetTextColor(colors.Black)
	gc.SetFillColor(colors.WhiteSmoke)
	gc.Clear()
	for i, name := range palNames {
		col, row := i/NumRows, i%NumRows
		x0 := col * (FieldWidth + ColorPaddingHori)
		y0 := row * (FieldHeight + ColorPaddingVert + TextHeight)
		pal := palMap[name]
		for w := range FieldWidth + 1 {
			t := float64(w) / float64(FieldWidth)
			color := pal.Color(t)
			x := x0 + w
			for h := range FieldHeight {
				y := y0 + h
				gc.SetPixel(x, y, color)
			}
		}

		gc.SetFontFace(nameFace)
		gc.DrawStringAnchored(name, float64(x0), float64(y0)+FieldHeight+TextHeight, 0.0, 0.0)
		gc.SetFontFace(typeFace)
		gc.DrawStringAnchored(pal.Type().String(), float64(x0)+FieldWidth, float64(y0)+FieldHeight+TextHeight, 1.0, 0.0)
	}
	gc.SavePNG("data/palNew.png")
}
