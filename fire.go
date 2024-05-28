package ledgrid

import (
	"image"
	"image/color"
	"math/rand"
)

const (
	fireDefCooling  = 0.27 // Range: 20 - 100; Default:  70
	fireDefSparking = 0.47 // Range: 50 - 200; Default: 120
	fireYScaling    = 6
)

type Fire struct {
	VisualEmbed
	lg                *LedGrid
	rect              image.Rectangle
	heat              [][]float64
	pal               ColorSource
	cooling, sparking float64
	params            []*Bounded[float64]
	anim              Animation
}

func NewFire(lg *LedGrid) *Fire {
	f := &Fire{}
	f.VisualEmbed.Init("Fire")
	f.lg = lg
	f.rect = image.Rect(0, 0, lg.Bounds().Dx(), fireYScaling*lg.Bounds().Dy())
	f.heat = make([][]float64, f.rect.Dx())
	for i := range f.heat {
		f.heat[i] = make([]float64, f.rect.Dy())
	}
	f.pal = FirePalette

	f.params = make([]*Bounded[float64], 2)
	f.params[0] = NewBounded("Cooling factor", fireDefCooling, 0.08, 1.00, 0.05)
	f.params[0].BindVar(&f.cooling)
	f.params[1] = NewBounded("Sparking factor", fireDefSparking, 0.19, 0.78, 0.05)
	f.params[1].BindVar(&f.sparking)
	f.anim = NewInfAnimation(f.Update)

	return f
}

func (f *Fire) ParamList() []*Bounded[float64] {
	return f.params
}

func (f *Fire) SetVisible(vis bool) {
	if vis {
		f.anim.Start()
	} else {
		f.anim.Stop()
	}
	f.VisualEmbed.SetVisible(vis)
}

func (f *Fire) Update(t float64) {

	// Cool down all heat points
	maxCooling := ((10.0 * f.cooling) / float64(f.rect.Dy())) + 0.0078
	for col := range f.heat {
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
		for row := f.rect.Dy() - 1; row >= 2; row-- {
			f.heat[col][row] = (f.heat[col][row-1] + 2.0*f.heat[col][row-2]) / 3.0
		}
	}

	// Random create new heat cells
	for col := range f.heat {
		if rand.Float64() < f.sparking {
			row := rand.Intn(4)
			heat := f.heat[col][row]
			spark := 0.625 + 0.375*rand.Float64()
			if spark >= 1.0-heat {
				f.heat[col][row] = 1.0
			} else {
				f.heat[col][row] = heat + spark
			}
		}
	}
}

func (f *Fire) ColorModel() color.Model {
	return LedColorModel
}

func (f *Fire) Bounds() image.Rectangle {
	return f.rect
}

func (f *Fire) At(x, y int) color.Color {
	y = f.rect.Dy() - 1 - y
	return f.pal.Color(f.heat[x][y])
}
