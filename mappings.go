package ledgrid

import "image"

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
	Edge   EdgeType
	Color  LedColor
}

var (
	blurVal    = 1.0 / 9.0
	BlurKernel = KernelType{
		Kernel: [][]float64{
			{blurVal, blurVal, blurVal},
			{blurVal, blurVal, blurVal},
			{blurVal, blurVal, blurVal},
		},
		Edge:  Constant,
		Color: Black,
	}

	FadeKernel = KernelType{
		Kernel: [][]float64{
			{0.0, 0.0, 0.0},
			{0.0, 0.9, 0.0},
			{0.0, 0.0, 0.0},
		},
		Edge:  Constant,
		Color: Black,
	}

	SharpKernel = KernelType{
		Kernel: [][]float64{
			{0.0, -1.0, 0.0},
			{-1.0, 5.0, -1.0},
			{0.0, -1.0, 0.0},
		},
		Edge:  Constant,
		Color: Black,
	}

	EdgeKernel = KernelType{
		Kernel: [][]float64{
			{0.0, 1.0, 0.0},
			{1.0, -4.0, 1.0},
			{0.0, 1.0, 0.0},
		},
		Edge:  Constant,
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
			lg.SetLedColor(col, row, LedColor{uint8(r), uint8(g), uint8(b), 0xFF})
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
