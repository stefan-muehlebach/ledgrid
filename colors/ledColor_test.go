package colors

import (
	"fmt"
	"testing"
    "math/rand/v2"
)

var (
    c1 = LedColor{0x00, 0x00, 0x00, 0xFF}
    c2 = LedColor{0xFF, 0xFF, 0xFF, 0xFF}
    c3 LedColor
)

func TestLedColorHue(t *testing.T) {
	colorList := []LedColor{
		Black, White, LedColor{127, 127, 127, 255},
		Red, LedColor{0, 255, 0, 255}, Blue,
		Cyan, Magenta, Yellow,
	}

	for _, col := range colorList {
		fmt.Printf("%v\n", col)
		h, s, l := col.HSL()
		fmt.Printf("  H: %f\n", h)
		fmt.Printf("  S: %f\n", s)
		fmt.Printf("  L: %f\n", l)
	}
}


func BenchmarkInterpolate(b *testing.B) {
    for b.Loop() {
        c3 = c1.Interpolate(c2, rand.Float64())
    }
}
