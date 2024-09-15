package ledgrid

// The output device can be realized in several ways. For example a SPI bus
// with a chain of WS2801 chips, driving combined LED's - this for sure is
// the default. But think of development: there may be an emulation of this
// (quite expensive) hardware. Every "thing" which can be used as a display
// must implement this interface.
type Displayer interface {
	// Returns the width and height of the display. Is mainly used by the
	// GridServer to determine the size of the receiving buffer.
	Size() int
	// Returns the gamma values that should by default be used on this
	// type of hardware.
	DefaultGamma() (r, g, b float64)
	// Sends all the bytes in buffer to the displaying hardware.
	// The correct order of the pixel values as well as the order of the
	// colors is up to the sending part.
	Send(buffer []byte)
	// Closes the connection to the displaying hardware and releases any
	// allocated ressources.
	Close()
}
