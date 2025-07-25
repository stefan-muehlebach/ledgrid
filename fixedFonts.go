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

var (
	// ' ' ! " # $ % & ' ( ) * + , - . /  (16 Glyphs)
    // glyphRangeMarks1 = basicfont.Range{'\u0020', '\u0030', 0}
	// '0'-'9'                            (10 Glyphs: the decimal digits)
	// glyphRangeDigits = basicfont.Range{'\u0030', '\u003a', 16}
	// ':' ';' '<' '=' '>' '?' '@'        (7 Glyphs)
	// glyphRangeMarks2 = basicfont.Range{'\u003a', '\u0041', 26}
	// 'A'-'Z'                            (26 Glyphs: lowercase characters)
	// glyphRangeUpper = basicfont.Range{'\u0041', '\u005b', 33}
	// '[' '\' ']' '^' '_' '`'            (6 Glyphs)
	// glyphRangeMarks3 = basicfont.Range{'\u005b', '\u0061', 59}
	// 'a'-'z'                            (26 Glyphs: uppercase characters)
	// glyphRangeLower = basicfont.Range{'\u0061', '\u007b', 65}
	// '{' '|' '}' '~'                    (4 Glyphs)
	// glyphRangeMarks4 = basicfont.Range{'\u007b', '\u007f', 91}

    glyphRangeFull = []basicfont.Range{
        {'\u0020', '\u007f', 0},
    }
)

// This function can be used to produce a new fixed font by scaling an existing
// fixed font. Scaling factors can only be positive integers. The new font
// is
func ScaleFixedFont(face *basicfont.Face, factor int, newName string) {
	width := face.Width
	height := face.Ascent
	mask := face.Mask.(*image.Alpha)
	fileName := fmt.Sprintf("font%s.go", newName)
	fh, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Couldn't create file: %v", err)
	}
	defer fh.Close()
	fmt.Fprintf(fh, "package ledgrid\n\n")
	fmt.Fprintf(fh, "import (\n")
	fmt.Fprintf(fh, "    \"image\"\n")
	fmt.Fprintf(fh, "    \"golang.org/x/image/font/basicfont\"\n")
	fmt.Fprintf(fh, ")\n\n")
	fmt.Fprintf(fh, "var %s = &basicfont.Face{\n", newName)
	fmt.Fprintf(fh, "    Advance: %d,\n", factor*face.Advance)
	fmt.Fprintf(fh, "    Width:   %d,\n", factor*face.Width)
	fmt.Fprintf(fh, "    Height:  %d,\n", factor*face.Height)
	fmt.Fprintf(fh, "    Ascent:  %d,\n", factor*face.Ascent)
	fmt.Fprintf(fh, "    Descent: %d,\n", factor*face.Descent)
	fmt.Fprintf(fh, "    Mask:    mask%s,\n", newName)
	fmt.Fprintf(fh, "    Ranges:  glyphRange,\n")
	fmt.Fprintf(fh, "}\n\n")
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
