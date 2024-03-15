package ledgrid

import (
	"log"
	"image"
	"image/color"
)

var (
	ClearColor   = LedColor{0, 0, 0}
	OutsideColor = LedColor{0, 0, 0}
)

// Entspricht dem Bild, welches auf einem LED-Panel angezeigt werden kann.
// Implementiert die Interfaces image.Image und draw.Image, also die Methoden
// ColorModel, Bounds, At und Set.
type LedGrid struct {
	// Groesse des LedGrids. Falls dieses LedGrid Teil eines groesseren
	// Panels sein sollte, dann muss Rect.Min nicht unbedingt {0, 0} sein.
	Rect image.Rectangle
	// Enthaelt die Farbwerte red, green, blue (RGB) fuer jede LED, welche
	// das LedGrid ausmachen. Die Reihenfolge entspricht dabei der
	// Verkabelung, d.h. sie beginnt links oben mit der LED Nr. 0,
	// geht dann nach rechts und auf der zweiten Zeile wieder nach links und
	// so schlangenfoermig weiter.
	Pix []uint8
}

func NewLedGrid(r image.Rectangle) *LedGrid {
	g := &LedGrid{}
	g.Rect = r
	g.Pix = make([]uint8, 3*r.Dx()*r.Dy())
	return g
}

func (g *LedGrid) ColorModel() color.Model {
	return LedColorModel
}

func (g *LedGrid) Bounds() image.Rectangle {
	return g.Rect
}

func (g *LedGrid) At(x, y int) color.Color {
	return g.LedColorAt(x, y)
}

func (g *LedGrid) Set(x, y int, c color.Color) {
	c1 := LedColorModel.Convert(c).(LedColor)
	g.SetLedColor(x, y, c1)
}

// Dient dem schnelleren Zugriff auf den Farbwert einer bestimmten Zelle, resp.
// einer bestimmten LED. Analog zu At(), retourniert den Farbwert jedoch als
// LedColor-Typ.
func (g *LedGrid) LedColorAt(x, y int) LedColor {
	if !(image.Point{x, y}.In(g.Rect)) {
		log.Printf("LedColorAt(): point outside LedGrid: %d, %d\n", x, y)
		return LedColor{}
	}
	idx := g.PixOffset(x, y)
	slc := g.Pix[idx : idx+3 : idx+3]
	return LedColor{slc[0], slc[1], slc[2]}
}

// Analoge Methode zu Set(), jedoch ohne zeitaufwaendige Konvertierung.
func (g *LedGrid) SetLedColor(x, y int, c LedColor) {
	if !(image.Point{x, y}.In(g.Rect)) {
		log.Printf("SetLedColor(): point outside LedGrid: %d, %d\n", x, y)
        return
	}
	idx := g.PixOffset(x, y)
	slc := g.Pix[idx : idx+3 : idx+3]
	slc[0] = c.R
	slc[1] = c.G
	slc[2] = c.B
}

// Damit wird der Offset eines bestimmten Farbwerts innerhalb des Slices
// Pix berechnet. Dabei wird beruecksichtigt, dass das die LED's im LedGrid
// schlangenfoermig angeordnet sind, wobei die Reihe mit der LED links oben
// beginnt.
func (g *LedGrid) PixOffset(x, y int) int {
	var idx int

	idx = y * g.Rect.Dx()
	if y%2 == 0 {
		idx += x
	} else {
		idx += (g.Rect.Dx() - x - 1)
	}
	return 3 * idx
}

// Hier kommen nun die fuer das LedGrid spezifischen Funktionen.
func (g *LedGrid) Clear() {
	for idx := 0; idx < len(g.Pix); idx += 3 {
		slc := g.Pix[idx : idx+3 : idx+3]
		slc[0] = ClearColor.R
		slc[1] = ClearColor.G
		slc[2] = ClearColor.B
	}
}

// func (g *LedGrid) FadeOld() {
// 	for row := 0; row < g.Rect.Dy(); row++ {
// 		for col := 0; col < g.Rect.Dx(); col++ {
// 			c := g.LedColorAt(col, row).Interpolate(FadeColor, FadeFactor)
// 			g.SetLedColor(col, row, c)
// 		}
// 	}
// }

// EdgeType bestimmt, wie bei einer Faltung der Rand behandelt wird.
type EdgeType int

const (
    // Mit Constant werden Pixel ausserhalb des darstellbaren Bereiches als
    // Pixel mit einer bestimmten, konstanten Farbe (normalerweise Schwarz)
    // interpretiert.
	Constant EdgeType = iota
    // Mit Crop werden die Randpixel fuer eine Faltung zwar beruecksichtig,
    // werden aber nie neu berechnet, da die Faltungsmatrix nur innerhalb des
    // darstellbaren Bereiches zu liegen kommt.
	Crop
    // Mit Wrap werden die ausserhalb liegenden Pixel mit Farbwerten der gegen-
    // ueberliegenden Seite belegt.
	Wrap
    // Mit Mirror werden ausserhalb liegende Pixel mit Farbwerten jener Pixel
    // belegt, die sich direkt vis-a-vis des Randes befinden.
	Mirror
    // Mit Extend schliesslich werden die Farben der Randpixel linear nach
    // aussen projeziert.
	Extend
)

type KernelType struct {
    Kernel [][]float64
    Edge EdgeType
    Color LedColor
}

var (
    blurVal = 1.0 / 9.0
	BlurKernel = KernelType{
        Kernel: [][]float64{
		    {blurVal, blurVal, blurVal},
        		{blurVal, blurVal, blurVal},
        		{blurVal, blurVal, blurVal},
        },
        Edge: Constant,
        Color: Black,
	}

	FadeKernel = KernelType{
        Kernel: [][]float64{
		    {0.0, 0.0, 0.0},
		    {0.0, 0.9, 0.0},
		    {0.0, 0.0, 0.0},
        },
        Edge: Constant,
        Color: Black,
	}

    SharpKernel = KernelType{
        Kernel: [][]float64{
            { 0.0, -1.0,  0.0},
            {-1.0,  5.0, -1.0},
            { 0.0, -1.0,  0.0},
        },
        Edge: Constant,
        Color: Black,
    }

    EdgeKernel = KernelType{
        Kernel: [][]float64{
            {0.0,  1.0, 0.0},
            {1.0, -4.0, 1.0},
            {0.0,  1.0, 0.0},
        },
        Edge: Constant,
        Color: Black,
    }
)

func (lg *LedGrid) convolute(op KernelType) {
	var pix []byte
	var r, g, b float64
	var c LedColor
	var x, y int

	pix = make([]byte, len(lg.Pix))
	copy(pix, lg.Pix)
	for row := range lg.Rect.Dy() {
		for col := range lg.Rect.Dx() {
			r, g, b = 0.0, 0.0, 0.0
		matrixLoop:
			for i, opValues := range op.Kernel {
				for j, opVal := range opValues {
					if opVal == 0.0 {
						continue
					}
					x, y = col-1+j, row-1+i
					if !(image.Point{x, y}.In(lg.Rect)) {
						switch op.Edge {
						case Constant:
							c = op.Color
						case Crop:
							r, g, b = 0.0, 0.0, 0.0
							break matrixLoop
						case Wrap:
							x = (x + lg.Rect.Dx()) % lg.Rect.Dx()
							y = (y + lg.Rect.Dy()) % lg.Rect.Dy()
							c = lg.LedColorAt(x, y)
						case Mirror:
							if x < 0 {
								x += 2
							} else if x >= lg.Rect.Dx() {
								x -= 2
							}
							if y < 0 {
								y += 2
							} else if y >= lg.Rect.Dy() {
								y -= 2
							}
							c = lg.LedColorAt(x, y)
						case Extend:
							x = max(lg.Rect.Min.X, min(lg.Rect.Max.X-1, x))
							y = max(lg.Rect.Min.Y, min(lg.Rect.Max.Y-1, y))
							c = lg.LedColorAt(x, y)
						}
					} else {
						c = lg.LedColorAt(x, y)
					}

					r += float64(c.R) * opVal
					g += float64(c.G) * opVal
					b += float64(c.B) * opVal
				}
			}
			lg.SetLedColor(col, row, LedColor{uint8(r), uint8(g), uint8(b)})
		}
	}
}

func (g *LedGrid) Blur() {
	g.convolute(BlurKernel)
}

func (g *LedGrid) Fade() {
	g.convolute(FadeKernel)
}

func (g *LedGrid) Sharpen() {
    g.convolute(SharpKernel)
}

func (g *LedGrid) Edges() {
    g.convolute(EdgeKernel)
}
