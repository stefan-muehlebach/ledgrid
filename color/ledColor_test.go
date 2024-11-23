package color

import (
	"fmt"
	"testing"
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
