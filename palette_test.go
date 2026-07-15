package ledgrid

import (
	"log"
	"os"
	"testing"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/fonts"
)

const (
	ColorBoxWidth  = 600
	ColorBoxHeight = 100
	TextBoxHeight  = 30.0
	//TileRows         = 4
	//TileCols         = 256 / TileRows
	//TileWidth        = ColorBoxWidth / TileCols
	//TileHeight       = ColorBoxHeight / TileRows
	//CornerRadius     = 20
	//Margin           = 10.0
	ColorPaddingHori = 20.0
	ColorPaddingVert = 20.0
	FontSize         = 18.0
	FontSizeSmall    = 14.0
	JSONFileName     = "data/palNew.json"
	PNGFileName      = "data/palNew.png"
)

var (
	// Von den beiden Variablen NumCols und NumRows kann eine Null sein, dann
	// wird sie aufgrund der anderen Variable (die nicht Null sein darf)
	// berechnet.
	NumCols  = 4
	NumRows  = 0
	NameFont = fonts.LucidaBrightDemibold
	TypeFont = fonts.LucidaBrightItalic
)

func TestPalette(t *testing.T) {

	fh, err := os.Open(JSONFileName)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer fh.Close()

	palNames, palMap, err := colors.ReadPaletteData(fh)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	if NumCols == 0 {
		if NumRows == 0 {
			log.Fatal("NumCols and NumRows can't be both zero!")
		}
		NumCols = len(palNames) / NumRows
		if len(palNames) % NumRows != 0 {
			NumCols += 1
		}
	} else {
		if NumRows == 0 {
			NumRows = len(palNames) / NumCols
			if len(palNames)%NumCols != 0 {
				NumRows += 1
			}
		}
	nameFace, _ := fonts.NewFace(NameFont, FontSize)
	typeFace, _ := fonts.NewFace(TypeFont, FontSizeSmall)
	gc := gg.NewContext(NumCols*(ColorBoxWidth)+(NumCols-1)*ColorPaddingHori,
		NumRows*ColorBoxHeight+NumRows*(ColorPaddingVert+TextBoxHeight))
	gc.SetTextColor(colors.Black)
	gc.SetFillColor(colors.WhiteSmoke)
	gc.Clear()
	for i, name := range palNames {
		col, row := i/NumRows, i%NumRows
		x0 := col * (ColorBoxWidth + ColorPaddingHori)
		y0 := row * (ColorBoxHeight + ColorPaddingVert + TextBoxHeight)
		pal := palMap[name]
		for w := range ColorBoxWidth + 1 {
			t := float64(w) / float64(ColorBoxWidth)
			color := pal.Color(t)
			x := x0 + w
			for h := range ColorBoxHeight {
				y := y0 + h
				gc.SetPixel(x, y, color)
			}
		}

		gc.SetFontFace(nameFace)
		gc.DrawStringAnchored(name, float64(x0), float64(y0)+ColorBoxHeight+TextBoxHeight, 0.0, 0.0)
		gc.SetFontFace(typeFace)
		gc.DrawStringAnchored(pal.Type().String(), float64(x0)+ColorBoxWidth, float64(y0)+ColorBoxHeight+TextBoxHeight, 1.0, 0.0)
	}
	gc.SavePNG(PNGFileName)
}
