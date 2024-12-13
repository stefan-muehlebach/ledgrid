package main

import (
	"image"
	"testing"
)

var (
    img *image.RGBA
    flt *ShuffleFilter
    x, y int
)

func init() {
    img = image.NewRGBA(image.Rect(0, 0, 80, 20))
    flt = NewShuffleFilter(img)
}

func TestIndex2Coord(t *testing.T) {
    index2coord(135)
}

func BenchmarkFilterRandFF(b *testing.B) {
    x, y = 20, 10
    for i:=0; i<b.N; i++ {
        x, y = flt.FF(x, y)
    }
}

func BenchmarkFilterRandShuffle(b *testing.B) {
    for i:=0; i<b.N; i++ {
        flt.Shuffle()
    }
}

func BenchmarkFilterRandOrderOne(b *testing.B) {
    flt.Shuffle()
    b.ResetTimer()
    for i:=0; i<b.N; i++ {
        flt.OrderOne()
    }
}



