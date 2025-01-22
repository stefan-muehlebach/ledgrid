package ledgrid

import (
	"github.com/stefan-muehlebach/ledgrid/conf"
)

// Um den clientseitigen Code so generisch wie moeglich zu halten, ist der
// GridClient als Interface definiert. Aktuell stehen zwei Implementationen
// am Start:
//
//		NetGridClient  - Verbindet sich via UDP und RPC mit einem externen
//		                 gridController.
//	 FileSaveClient - Schreibt die Bilddaten in ein File, welches dann auf das
//	                  System mit dem Grid-Controller kopiert und dort direkt
//	                  abgespielt werden kann.
type GridClient interface {
	Send(buffer []byte)
	NumLeds() int
	Gamma() (r, g, b float64)
	SetGamma(r, g, b float64)
	ModuleConfig() conf.ModuleConfig
	Watch() *Stopwatch
	Close()
}
