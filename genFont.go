//go:build ignore

package ledgrid

import (
	"github.com/stefan-muehlebach/ledgrid"
)

func main() {
	ledgrid.ScaleFixedFont(ledgrid.Pico3x5, 2, "Pico6x10")
	ledgrid.ScaleFixedFont(ledgrid.Pico3x5, 3, "Pico9x15")
}
