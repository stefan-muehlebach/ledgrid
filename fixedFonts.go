// Diese Datei enthaelt zwei Pixel- resp. Bitmap-Schriften, die ich zum einen
// der PICO-8 Umgebung entliehen habe oder auf Pinterest gefunden. Die
// eigentlichen Font-Daten (die sog. Masken) befinden sich in den Dateien
// fixedMask3x5.go und fixedMask5x7.go

package ledgrid

import (
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

var Face3x5 = &basicfont.Face{
	Advance: 4,
	Width:   3,
	Height:  5,
	Ascent:  5,
	Descent: 0,
	Mask:    mask3x5,
	Ranges:  glyphRange,
}

var Face5x7 = &basicfont.Face{
	Advance: 6,
	Width:   5,
	Height:  7,
	Ascent:  7,
	Descent: 0,
	Mask:    mask5x7,
	Ranges:  glyphRange,
}
