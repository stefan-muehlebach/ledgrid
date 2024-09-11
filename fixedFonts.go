// Diese Datei enthaelt zwei Pixel- resp. Bitmap-Schriften, die ich zum einen
// der PICO-8 Umgebung entliehen habe oder auf Pinterest gefunden. Die
// eigentlichen Font-Daten (die sog. Masken) befinden sich in den Dateien
// fixedMask3x5.go und fixedMask5x7.go

package ledgrid

import (
	"fmt"
	"image"
	"log"
	"os"

	"golang.org/x/image/font/basicfont"
)

const (
	numGlyphs = 95
)

var glyphRange = []basicfont.Range{
	{'\u0020', '\u0030', 0},  // ' ' ! " # $ % & ' ( ) * + , - . /  (16 Glyphs)
	{'\u0030', '\u003a', 16}, // '0'-'9'                            (10 Glyphs)
	{'\u003a', '\u0041', 26}, // ':' ';' '<' '=' '>' '?' '@'        (7 Glyphs)
	{'\u0041', '\u005b', 33}, // 'A'-'Z'                            (26 Glyphs)
	{'\u005b', '\u0061', 59}, // '[' '\' ']' '^' '_' '`'            (6 Glyphs)
	{'\u0061', '\u007b', 65}, // 'a'-'z'                            (26 Glyphs)
	{'\u007b', '\u007f', 91}, // '{' '|' '}' '~'                    (4 Glyphs)
}

// Original Pico-8 font in the original size (3x5 Pixels!)
var Pico3x5 = &basicfont.Face{
	Advance: 4,
	Width:   3,
	Height:  8,
	Ascent:  5,
	Descent: 0,
	Mask:    maskPico3x5,
	Ranges:  glyphRange,
}

var Fixed5x7 = &basicfont.Face{
	Advance: 6,
	Width:   5, // Dies ist die Breite eines Buchstabens gem. Maske
	Height:  9,
	Ascent:  7, // Dies ist die Hoehe eines Buchstabens gem. Maske
	Descent: 0,
	Mask:    maskFixed5x7,
	Ranges:  glyphRange,
}

func ScaleFace(face *basicfont.Face, factor int, newName string) {
	width := face.Width
	height := face.Ascent
	mask := face.Mask.(*image.Alpha)
	fileName := fmt.Sprintf("mask%s.go", newName)
	fh, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Couldn't create file: %v", err)
	}
	defer fh.Close()
	fmt.Fprintf(fh, "package ledgrid\n\n")
	fmt.Fprintf(fh, "import (\n")
	fmt.Fprintf(fh, "    \"image\"\n")
	fmt.Fprintf(fh, ")\n\n")
	fmt.Fprintf(fh, "var mask%s = &image.Alpha{\n", newName)
	fmt.Fprintf(fh, "    Stride: %d,\n", factor*width)
	fmt.Fprintf(fh, "    Rect:   image.Rectangle{Max: image.Point{%d, %d}},\n",
		factor*mask.Rect.Max.X, factor*mask.Rect.Max.Y)
	fmt.Fprintf(fh, "    Pix: []byte{\n")
	idx := 0
	for idx < len(mask.Pix) {
		row := mask.Pix[idx : idx+mask.Stride]
		for range factor {
			fmt.Fprintf(fh, "        ")
			for _, val := range row {
				for range factor {
					fmt.Fprintf(fh, "0x%02x, ", val)
				}
			}
			fmt.Fprintf(fh, "\n")
		}
		idx += mask.Stride
		if idx%(height*width) == 0 {
			fmt.Fprintf(fh, "\n")
		}
	}
	fmt.Fprintf(fh, "    },\n")
	fmt.Fprintf(fh, "}\n")
}
