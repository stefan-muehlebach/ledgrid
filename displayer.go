package ledgrid

import (
	"log"
	"math"
)

// The output device can be realized in several ways. For example a SPI bus
// with a chain of WS2801 chips, driving combined LED's - this for sure is
// the default. But think of development: there may be an emulation of this
// (quite expensive) hardware. Every "thing" which can be used as a display
// must implement this interface.
type Displayer interface {
	// Returns the width and height of the display. Is mainly used by the
	// GridServer to determine the size of the receiving buffer.
	Size() int
	// Blah...
	SetPixelStatus(idx int, stat LedStatusType)
	// Returns the gamma values that should by default be used on this
	// type of hardware.
	DefaultGamma() (r, g, b float64)
	// Returns the currently set values for the gamma correction
	Gamma() (r, g, b float64)
	// Set new values for the gamma correction.
	SetGamma(r, g, b float64)
	//
	Display(buffer []byte)
	// Sends all the bytes in buffer to the displaying hardware.
	// The correct order of the pixel values as well as the order of the
	// colors is up to the sending part.
	Send(buffer []byte)
	// Closes the connection to the displaying hardware and releases any
	// allocated ressources.
	Close()
}

// This type has been introduced in order to mark some LEDs on the chain as
// 'ok', defect or missing (see constants PixelOK, PixelDefect, PixelMissing)
type LedStatusType byte

const (
	// This is the default state of a LED
	LedOK LedStatusType = iota
	// LEDs with this status will be blacked out, this mean we send color data
	// for this LED, but we send (0,0,0). This status can be used if a NeoPixel
	// receives data but does not display them correctly.
	LedDefect
	// This status can be used, if a NeoPixel does not even transmit the data
	// to the NeoPixels further down the chain. Such a pixel needs to be cut
	// out of the chain and for the time till a replacement Pixel is organized
	// and soldered in, the pixel has status missing.
	LedMissing
)

type DisplayEmbed struct {
	impl       Displayer
	size       int
	buffer     []byte
	gammaVal   [3]float64
	gammaTbl   [3][256]byte
	statusList []LedStatusType
}

func (d *DisplayEmbed) Init(impl Displayer, size int) {
	d.impl = impl
	d.size = size
	d.buffer = make([]byte, 3*size)
	d.statusList = make([]LedStatusType, size)
	d.SetGamma(impl.DefaultGamma())
}

func (d *DisplayEmbed) Size() int {
	return d.size
}

func (d *DisplayEmbed) SetPixelStatus(idx int, stat LedStatusType) {
	d.statusList[idx] = stat
}

func (d *DisplayEmbed) Gamma() (r, g, b float64) {
	return d.gammaVal[0], d.gammaVal[1], d.gammaVal[2]
}

func (d *DisplayEmbed) SetGamma(r, g, b float64) {
	d.gammaVal[0], d.gammaVal[1], d.gammaVal[2] = r, g, b
	for colorIdx, val := range d.gammaVal {
		for i := range 256 {
			d.gammaTbl[colorIdx][i] = byte(255.0 * math.Pow(float64(i)/255.0, val))
		}
	}
}

func (d *DisplayEmbed) Display(buffer []byte) {
	var srcIdx, dstIdx int
	var src, dst []byte
	var bufLen int

	bufLen = len(buffer)
	for srcIdx, dstIdx = 0, 0; srcIdx < len(buffer)/3; srcIdx++ {
		if d.statusList[srcIdx] == LedMissing {
			continue
		}
		dst = d.buffer[3*dstIdx : 3*dstIdx+3 : 3*dstIdx+3]
		if d.statusList[srcIdx] == LedDefect {
			dst[0] = 0x00
			dst[1] = 0x00
			dst[2] = 0x00
		} else {
			src = buffer[3*srcIdx : 3*srcIdx+3 : 3*srcIdx+3]
			dst[0] = d.gammaTbl[0][src[0]]
			dst[1] = d.gammaTbl[1][src[1]]
			dst[2] = d.gammaTbl[2][src[2]]
		}
		dstIdx++
	}
	d.impl.Send(d.buffer)
}
