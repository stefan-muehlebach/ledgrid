package ledgrid

import (
	"math"
	"math/rand"
	"testing"
)

var (
	gammaValue [3]float64
	gamma      [3][256]byte
	buffer     []byte
	bufferSize = 3 * 40 * 10
)

func init() {
	buffer = make([]byte, bufferSize)
	gammaValue[0], gammaValue[1], gammaValue[2] = 3.0, 3.0, 3.0
	for color, val := range gammaValue {
		for i := range 256 {
			gamma[color][i] = byte(255.0 * math.Pow(float64(i)/255.0, val))
		}
	}
}

func randomizeBuffer() {
	rand.Seed(123_456_789)
	for i := range buffer {
		buffer[i] = byte(rand.Intn(256))
	}
}

// Was ist schneller: die Gamma-Anpassung mit jeweils drei Zeilen pro Farbe...
func BenchmarkGammaByThree(b *testing.B) {
	randomizeBuffer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(buffer); j += 3 {
			buffer[j] = gamma[0][buffer[j]]
			buffer[j+1] = gamma[1][buffer[j+1]]
			buffer[j+2] = gamma[2][buffer[j+2]]
		}
	}
}

// ... (wobei hier noch die Variante mit Subslices getestet wird)...
func BenchmarkGammaByThreeSubSlice(b *testing.B) {
	var dst, src []byte
	var i, j int

	randomizeBuffer()
	b.ResetTimer()
	for i = 0; i < b.N; i++ {
		for j = 0; j < len(buffer); j += 3 {
			src = buffer[j : j+3 : j+3]
			dst = buffer[j : j+3 : j+3]
			dst[0] = gamma[0][src[0]]
			dst[1] = gamma[1][src[1]]
			dst[2] = gamma[2][src[2]]
		}
	}
}

// ... oder einem einzelnen Durchlauf mit Modulo-Operation.
func BenchmarkGammaByMod(b *testing.B) {
	randomizeBuffer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j, ch := range buffer {
			buffer[j] = gamma[j%3][ch]
		}
	}
}
