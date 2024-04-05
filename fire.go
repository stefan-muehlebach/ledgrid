//go:build ignore
// +build ignore

package ledgrid

import (
	"time"
	"image"
	"math/rand"
)

const (
    cooling = 55.0
    sparking = 120.0

type Fire struct {
    VisualizableEmbed
     lg *LedGrid
    numLeds int
    heat [][]float64
}

func NewFire(lg *LedGrid) *Fire {
    f := &Fire{}
    f.VisualizableEmbed.Init()
    f.lg = lg
    f.numLeds = lg.Rect.Dx() * lg.Rect.Dy()
    f.heat = make([][]float64, lg.Rect.Dy())
    for i := range f.heat {
        f.heat[i] = make([]float64, lg.Rect.Dx())
    }
    return f
}

func (f *Fire) Update(dt time.Duration) bool {
    dt = f.AnimatableEmbed.Update(dt)

    for row := range f.heat {
        for col := range f.heat[row] {


    return true
}
