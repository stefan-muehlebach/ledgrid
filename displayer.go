package ledgrid

import (
	"image"
	"math"

	"github.com/stefan-muehlebach/ledgrid/conf"
)

// The output device can be realized in several ways. For example a SPI bus
// with a chain of NeoPixels, driven by WS2801 - this for sure is the default.
// But think of development: there may be an emulation of this (quite
// expensive) hardware. Every"thing" which can be used as a display must
// implement this interface.
type Displayer interface {
	// Returns the number of addressable NeoPixels in this implementation.
	// Will be used by the Server to allocate memory for buffer.
	NumLeds() int
	// Returns a Rectangle enclosing the whole display. Because the
	// configuration of the modules allow empty space within this
	// rectangle, you cannot derive the number of NeoPixels from this value.
	Bounds() image.Rectangle
	// Returns the module configuration of this displayer.
	ModuleConfig() conf.ModuleConfig
	// Sets a new module configuration
	SetModuleConfig(cnf conf.ModuleConfig)
	// Sets the status of NeoPixel with index idx to stat.
	SetPixelStatus(idx int, stat LedStatusType)
	// Returns the gamma values that should be used by default on this
	// displayer.
	DefaultGamma() (r, g, b float64)
	// Returns the currently set values for the gamma correction
	Gamma() (r, g, b float64)
	// Set new values for the gamma correction.
	SetGamma(r, g, b float64)
	// Display is used by a Server to show the image data in buffer. The bytes
	// in buffer must already be in a suitable order for this specific device.
	// The order of RGB has to be in device order as well.
	Display(buffer []byte)
	// Send is called by Display and must not be called from other parts of
	// the software
	Send(buffer []byte)
	// Closes the connection to the displaying hardware and releases any
	// allocated ressources.
	Close()
}

// This type has been introduced in order to mark some NeoPixels on the chain
// as 'ok' (the default), 'defect' or 'missing' (see constants PixelOK,
// PixelDefect or PixelMissing for more information).
type LedStatusType byte

const (
	// LedOK is the default state of a NeoPixel.
	LedOK LedStatusType = iota
	// NeoPixels with status LedDefect will be blacked out. This means that
    // the sent color data for this pixel will always be (0,0,0), i.e. black.
	// This status can be used if a NeoPixel propagates data as expected
    // but is not able to correctly display its own color.
	LedDefect
	// LedMissing can be used, if a NeoPixel does not even propagate data
	// to NeoPixels further down the chain. Such a pixel has to be cut out
	// of the chain and the wires need to be shortened. For the time till a
	// replacement pixel is organized and soldered in, the pixel must have
	// status LedMissing.
	LedMissing
)

// Each implementation of a Displayer should embed this embeddable. It
// provides default implementations for a number of general methods.
type DisplayEmbed struct {
	ModConf    conf.ModuleConfig
	impl       Displayer
	numLeds    int
	size       image.Point
	buffer     []byte
	gammaVal   [3]float64
	gammaTbl   [3][256]byte
	statusList []LedStatusType
}

// An embedding type needs to call this method once in order to set initial
// values correctly.
func (d *DisplayEmbed) Init(impl Displayer, numLeds int) {
	d.impl = impl
	d.numLeds = numLeds
	d.buffer = make([]byte, 3*numLeds)
	d.statusList = make([]LedStatusType, numLeds)
	d.SetGamma(impl.DefaultGamma())
}

// See NumLeds in interface Displayer.
func (d *DisplayEmbed) NumLeds() int {
	return d.numLeds
}

// See Bounds in interface Displayer.
func (d *DisplayEmbed) Bounds() image.Rectangle {
	return image.Rectangle{image.Point{}, d.size}
}

// See SetPixelStatus in interface Displayer.
func (d *DisplayEmbed) SetPixelStatus(idx int, stat LedStatusType) {
	d.statusList[idx] = stat
}

// See ModuleConfig in interface Displayer.
func (d *DisplayEmbed) ModuleConfig() conf.ModuleConfig {
	return d.ModConf
}

// See SetModuleConfig in interface Displayer.
func (d *DisplayEmbed) SetModuleConfig(cnf conf.ModuleConfig) {
	d.ModConf = cnf
	d.size = cnf.Size()
}

// See Gamma in interface Displayer.
func (d *DisplayEmbed) Gamma() (r, g, b float64) {
	return d.gammaVal[0], d.gammaVal[1], d.gammaVal[2]
}

// See SetGamma in interface Displayer.
func (d *DisplayEmbed) SetGamma(r, g, b float64) {
	d.gammaVal[0], d.gammaVal[1], d.gammaVal[2] = r, g, b
	for colorIdx, val := range d.gammaVal {
		for i := range 256 {
			d.gammaTbl[colorIdx][i] = byte(255.0 * math.Pow(float64(i)/255.0, val))
		}
	}
}

// See Display in interface Displayer.
func (d *DisplayEmbed) Display(buffer []byte) {
	var srcIdx, dstIdx int
	var src, dst []byte

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
