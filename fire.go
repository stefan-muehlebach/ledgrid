package ledgrid

import (
	"image"
	"math/rand"
	"time"

	"github.com/stefan-muehlebach/gg/colors"
)

// ---------------------------------------------------------------------------

var (
	fireGradient = []colors.RGBA{
		colors.RGBA{0x00, 0x00, 0x00, 0x00},
		colors.RGBA{0x5f, 0x08, 0x09, 0x00},
		colors.RGBA{0x5f, 0x08, 0x09, 0x80},
		colors.RGBA{0xbe, 0x10, 0x13, 0x80},
		colors.RGBA{0xd2, 0x30, 0x08, 0x80},
		colors.RGBA{0xe4, 0x53, 0x23, 0xcf},
		colors.RGBA{0xee, 0x77, 0x1c, 0xcf},
		colors.RGBA{0xf6, 0x96, 0x0e, 0xcf},
		colors.RGBA{0xff, 0xcd, 0x06, 0xcf},
		colors.RGBA{0xff, 0xdb, 0x5a, 0xef},
		colors.RGBA{0xff, 0xe6, 0x68, 0xff},
	}

	fireYScaling    = 10
	fireDefCooling  = 0.35
	fireDefSparking = 0.47
)

type Fire struct {
	CanvasObjectEmbed
	Pos, Size         image.Point
	ySize             int
	heat              [][]float64
	cooling, sparking float64
	pal               ColorSource
	running           bool
}

func NewFire(pos, size image.Point) *Fire {
	f := &Fire{Pos: pos, Size: size}
	f.ySize = fireYScaling * size.Y
	f.CanvasObjectEmbed.Extend(f)
	f.heat = make([][]float64, f.Size.X)
	for i := range f.heat {
		f.heat[i] = make([]float64, f.ySize)
	}
	f.cooling = fireDefCooling
	f.sparking = fireDefSparking
	// f.pal = NewGradientPalette("Fire", fireGradient...)
	f.pal = colors.NewPaletteByColors("Fire", fireGradient...)
	// AnimCtrl.Add(0, f)
	return f
}

func (f *Fire) Duration() time.Duration {
	return time.Duration(0)
}

func (f *Fire) SetDuration(dur time.Duration) {}

func (f *Fire) StartAt(t time.Time) {
	if f.running {
		return
	}
	// Would do starting things here.
	f.running = true
	AnimCtrl.Add(f)
}

func (f *Fire) Start() {
	f.StartAt(AnimCtrl.Now())
}

func (f *Fire) Stop() {
	if !f.running {
		return
	}
	// Would do the stopping things here.
	f.running = false
}

func (f *Fire) Suspend() {}

func (f *Fire) Continue() {}

func (f *Fire) IsRunning() bool {
	return f.running
}

func (f *Fire) Update(pit time.Time) bool {
	// Cool down all heat points
	maxCooling := ((10.0 * f.cooling) / float64(f.ySize)) + 0.0078
	for col := range f.Size.X {
		for row, heat := range f.heat[col] {
			cooling := maxCooling * rand.Float64()
			if cooling >= heat {
				f.heat[col][row] = 0.0
			} else {
				f.heat[col][row] = heat - cooling
			}
		}
	}

	// Diffuse the heat
	for col := range f.heat {
		for row := f.ySize - 1; row >= 2; row-- {
			f.heat[col][row] = (f.heat[col][row-1] + 2.0*f.heat[col][row-2]) / 3.0
		}
	}

	// Random create new heat cells
	for col := range f.Size.X {
		if rand.Float64() < f.sparking {
			row := rand.Intn(4)
			heat := f.heat[col][row]
			spark := 0.5 + 0.5*rand.Float64()
			if spark >= 1.0-heat {
				f.heat[col][row] = 1.0
			} else {
				f.heat[col][row] = heat + spark
			}
		}
	}
	return true
}

func (f *Fire) Draw(c *Canvas) {
	for col := range f.Size.X {
		for row := range f.Size.Y {
			fireRow := fireYScaling * (f.Size.Y - row - 1)
			heat := f.heat[col][fireRow]
			bgColor := colors.RGBAModel.Convert(c.Img.At(col, row)).(colors.RGBA)
			fgColor := f.pal.Color(heat)
			c.Img.Set(col, row, fgColor.Mix(bgColor, colors.Blend))
		}
	}
}
