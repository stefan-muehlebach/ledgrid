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

var (
    	fireGradient = []ColorStop{
		{0.00, NewLedColorAlpha(0x00000000)},
		{0.10, NewLedColorAlpha(0x5f080900)},
		{0.14, NewLedColorAlpha(0x5f0809e5)},
		{0.29, NewLedColor(0xbe1013)},
		{0.43, NewLedColor(0xd23008)},
		{0.57, NewLedColor(0xe45323)},
		{0.71, NewLedColor(0xee771c)},
		{0.86, NewLedColor(0xf6960e)},
		{1.00, NewLedColor(0xffcd06)},
	}
)

type Fire struct {
	VisualEmbed
	lg                *LedGrid
	rect              image.Rectangle
	heat              [][]float64
	pal               ColorSource
	cooling, sparking float64
	params            []Parameter
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
	f.pal = NewGradientPalette("Fire", fireGradient...)

	f.params = make([]Parameter, 2)
	f.params[0] = NewFloatParameter("Cooling factor", fireDefCooling, 0.08, 1.00, 0.05)
	f.params[0].SetCallback(func (p Parameter) {
        v := f.params[0].(FloatParameter).Val()
        f.cooling = v
    })
    // f.params[0].BindVar(&f.cooling)

	f.params[1] = NewFloatParameter("Sparking factor", fireDefSparking, 0.19, 0.78, 0.05)
	f.params[1].SetCallback(func (p Parameter) {
        v := f.params[1].(FloatParameter).Val()
        f.sparking = v
    })
	// f.params[1].BindVar(&f.sparking)

	f.anim = NewInfAnimation(f.Update)

	return f
}

func (f *Fire) ParamList() []Parameter {
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
