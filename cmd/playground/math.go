//go:build ignore
// +build ignore

// Der Versuch, das math-Package von Go auf float32 umzuschreiben. Aktuell
// auf Eis gelegt - wird ev. spaeter wieder aufgenommen.

package main

import "unsafe"

const (
	Pi  = 3.14159265358979323846264338327950288419716939937510582097494459 // https://oeis.org/A000796
)

const (
	MaxFloat32             = 0x1p127 * (1 + (1 - 0x1p-23)) // 3.40282346638528859811704183484516925440e+38
)

const (
	uvnan = 0x7FC00001
	mask  = 0x7F8
	bias  = 127
	shift = 32 - 8 - 1
)

func NaN() float32 { return Float32frombits(uvnan) }

func IsNaN(f float32) bool {
	return f != f
}

func IsInf(f float32, sign int) bool {
    return sign >= 0 && f > MaxFloat32 || sign <= 0 && f < -MaxFloat32
}

// Float32bits returns the IEEE 754 binary representation of f,
// with the sign bit of f and the result in the same bit position.
// Float32bits(Float32frombits(x)) == x.
func Float32bits(f float32) uint32 { return *(*uint32)(unsafe.Pointer(&f)) }

// Float32frombits returns the floating-point number corresponding
// to the IEEE 754 binary representation b, with the sign bit of b
// and the result in the same bit position.
// Float32frombits(Float32bits(x)) == x.
func Float32frombits(b uint32) float32 { return *(*float32)(unsafe.Pointer(&b)) }

func Abs(x float32) float32 {
    return Float32frombits(Float32bits(x) &^ (1 << 31))
}

func Modf(f float32) (int float32, frac float32) {
	return modf(f)
}

func modf(f float32) (int float32, frac float32) {
	if f < 1 {
		switch {
		case f < 0:
			int, frac = Modf(-f)
			return -int, -frac
		case f == 0:
			return f, f
		}
		return 0, f
	}

	x := Float32bits(f)
	e := uint(x>>shift)&mask - bias

	if e < 32-9 {
		x &^= 1<<(32-9-e) - 1
	}
	int = Float32frombits(x)
	frac = f - int
	return
}

// sin coefficients
var _sin = [...]float32{
	1.58962301576546568060e-10, // 0x3de5d8fd1fd19ccd
	-2.50507477628578072866e-8, // 0xbe5ae5e5a9291f5d
	2.75573136213857245213e-6,  // 0x3ec71de3567d48a1
	-1.98412698295895385996e-4, // 0xbf2a01a019bfdf03
	8.33333333332211858878e-3,  // 0x3f8111111110f7d0
	-1.66666666666666307295e-1, // 0xbfc5555555555548
}

// cos coefficients
var _cos = [...]float32{
	-1.13585365213876817300e-11, // 0xbda8fa49a0861a9b
	2.08757008419747316778e-9,   // 0x3e21ee9d7b4e3f05
	-2.75573141792967388112e-7,  // 0xbe927e4f7eac4bc6
	2.48015872888517045348e-5,   // 0x3efa01a019c844f5
	-1.38888888888730564116e-3,  // 0xbf56c16c16c14f91
	4.16666666666665929218e-2,   // 0x3fa555555555554b
}

func Cos(x float32) float32 {
	return cos(x)
}

func cos(x float32) float32 {
	const (
		PI4A = 7.85398125648498535156e-1  // 0x3fe921fb40000000, Pi/4 split into three parts
		PI4B = 3.77489470793079817668e-8  // 0x3e64442d00000000,
		PI4C = 2.69515142907905952645e-15 // 0x3ce8469898cc5170,
	)
	// special cases
	switch {
	case IsNaN(x) || IsInf(x, 0):
		return NaN()
	}

	// make argument positive
	sign := false
	x = Abs(x)

	var j uint32
	var y, z float32
	if x >= reduceThreshold {
		j, z = trigReduce(x)
	} else {
		j = uint32(x * (4 / Pi)) // integer part of x/(Pi/4), as integer for tests on the phase angle
		y = float32(j)           // integer part of x/(Pi/4), as float

		// map zeros to origin
		if j&1 == 1 {
			j++
			y++
		}
		j &= 7                               // octant modulo 2Pi radians (360 degrees)
		z = ((x - y*PI4A) - y*PI4B) - y*PI4C // Extended precision modular arithmetic
	}

	if j > 3 {
		j -= 4
		sign = !sign
	}
	if j > 1 {
		sign = !sign
	}

	zz := z * z
	if j == 1 || j == 2 {
		y = z + z*zz*((((((_sin[0]*zz)+_sin[1])*zz+_sin[2])*zz+_sin[3])*zz+_sin[4])*zz+_sin[5])
	} else {
		y = 1.0 - 0.5*zz + zz*zz*((((((_cos[0]*zz)+_cos[1])*zz+_cos[2])*zz+_cos[3])*zz+_cos[4])*zz+_cos[5])
	}
	if sign {
		y = -y
	}
	return y
}

func Sin(x float32) float32 {
	return sin(x)
}

func sin(x float32) float32 {
	const (
		PI4A = 7.85398125648498535156e-1
		PI4B = 3.77489470793079817668e-8
		PI4C = 2.69515142907905952645e-15
	)
	// special cases
	switch {
	case x == 0 || IsNaN(x):
		return x
	case IsInf(x, 0):
		return NaN()
	}

	sign := false
	if x < 0 {
		x = -x
		sign = true
	}

	var j uint32
	var y, z float32
	if x >= reduceThreshold {
		j, z = trigReduce(x)
	} else {
		j = uint32(x * (4 / Pi))
		y = float32(j)

		if j&1 == 1 {
			j++
			y++
		}
		j &= 7
		z = ((x - y*PI4A) - y*PI4B) - y*PI4C
	}

	if j > 3 {
		sign = !sign
		j -= 4
	}

	zz := z * z
	if j == 1 || j == 2 {
		y = 1.0 - 0.5*zz + zz*zz*((((((_cos[0]*zz)+_cos[1])*zz+_cos[2])*zz+_cos[3])*zz+_cos[4])*zz+_cos[5])
	} else {
		y = z + z*zz*((((((_sin[0]*zz)+_sin[1])*zz+_sin[2])*zz+_sin[3])*zz+_sin[4])*zz+_sin[5])
	}
	if sign {
		y = -y
	}
	return y
}
