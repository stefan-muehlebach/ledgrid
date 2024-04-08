package ledgrid

import (
	"math/rand"
	"time"
)

const (
	cooling  = 70.0 // Range: 20 - 100;    Default: 50
	sparking = 0.4 // Range: 0.07 -  0.7; Default: 0.4
	scaling  = 4
)

type Fire struct {
	VisualizableEmbed
	lg      *LedGrid
	numLeds int
	heat    [][]float64
	pal     Colorable
}

func NewFire(lg *LedGrid) *Fire {
	f := &Fire{}
	f.VisualizableEmbed.Init("Fire")
	f.lg = lg
	f.numLeds = scaling * lg.Rect.Dy()
	f.heat = make([][]float64, lg.Rect.Dx())
	for i := range f.heat {
		f.heat[i] = make([]float64, f.numLeds)
	}
	f.pal = FirePalette
	return f
}

func (f *Fire) Update(dt time.Duration) bool {
	dt = f.AnimatableEmbed.Update(dt)

	// Cool down all heat points
	maxCooling := ((10.0*cooling)/float64(f.numLeds) + 2.0) / 255.0
	for col := range f.heat {
		for row := range f.heat[col] {
			f.heat[col][row] -= rand.Float64() * maxCooling
			if f.heat[col][row] < 0.0 {
				f.heat[col][row] = 0.0
			}
		}
	}

	// Diffuse the heat
	for col := range f.heat {
		for row := f.numLeds - 1; row >= 2; row-- {
			f.heat[col][row] = (f.heat[col][row-1] + 2*f.heat[col][row-2]) / 3.0
		}
	}

	// Random create new heat cells
	for col := range f.heat {
		if rand.Float64() < sparking {
			row := rand.Intn(2)
			f.heat[col][row] += 0.627 + 0.373*rand.Float64()
			if f.heat[col][row] > 1.0 {
				f.heat[col][row] = 1.0
			}
		}
	}

	return true
}

func (f *Fire) Draw() {
	var col, row int

	for row = range f.lg.Bounds().Dy() {
		heatRow := scaling * (f.lg.Bounds().Dy() - row - 1)
		for col = range f.lg.Bounds().Dx() {
			fireColor := f.pal.Color(f.heat[col][heatRow])
			f.lg.MixLedColor(col, row, fireColor, Max)
		}
	}
}
