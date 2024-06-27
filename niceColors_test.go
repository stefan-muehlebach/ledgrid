package ledgrid

import (
	"testing"

	"github.com/stefan-muehlebach/gg/fonts"
)

const (
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
	Font = fonts.GoBold
)

func TestReadColors(t *testing.T) {
	// face := fonts.NewFace(Font, FontSize)
	// cl := ReadColourlovers()
	// t.Logf("#colors: %d", len(cl))
	// gc := gg.NewContext(2*Margin+10*(5*FieldWidth+4*FieldPadding)+9*ColorPaddingHori,
	// 	2*Margin+10*FieldHeight+10*ColorPaddingVert)
	// gc.SetFontFace(face)
	// gc.SetFillColor(color.White)
	// gc.Clear()
	// for i, niceColor := range cl {
	// 	col, row := i/10, i%10
	// 	x := Margin + float64(col)*(5*FieldWidth+4*FieldPadding+ColorPaddingHori)
	// 	y := Margin + float64(row)*(FieldHeight+ColorPaddingVert)
	// 	for j, color := range niceColor.Colors {
	// 		x := x + float64(j)*(FieldWidth+FieldPadding)
	// 		gc.SetFillColor(color)
	// 		gc.SetStrokeWidth(0.0)
	// 		gc.DrawRoundedRectangle(x, y, FieldWidth, FieldHeight, CornerRadius)
	// 		gc.Fill()
	// 	}
	// 	gc.SetStrokeColor(color.Black)
	// 	gc.DrawStringAnchored(niceColor.Name, x, y+FieldHeight+Margin/2, 0.0, 1.0)
	// }
	// gc.SavePNG("niceColors.png")
}
