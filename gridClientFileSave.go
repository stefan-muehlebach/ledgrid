package ledgrid

import (
	"log"
	"os"

	"github.com/stefan-muehlebach/ledgrid/conf"
)

// Dieser Client-Typ schreibt alle Bilddaten in eine Datei, welche im
// Anschluss auf ein System mit echter Hardware transferiert und dort
// wie ein Film abgespielt wird.
type FileSaveClient struct {
	fh        *os.File
	modConf   conf.ModuleConfig
	sendWatch *Stopwatch
}

func NewFileSaveClient(fileName string, modConf conf.ModuleConfig) GridClient {
	var err error

	p := &FileSaveClient{}

	p.fh, err = os.Create(fileName)
	if err != nil {
		log.Fatalf("Couldn't create file: %v", err)
	}
	p.modConf = modConf
	p.sendWatch = NewStopwatch()

	return p
}

func (p *FileSaveClient) Send(buffer []byte) {
	if _, err := p.fh.Write(buffer); err != nil {
		log.Fatalf("Couldnt' write data to file: %v", err)
	}
}

func (p *FileSaveClient) NumLeds() int {
	return len(p.modConf) * (conf.ModuleDim.X * conf.ModuleDim.Y)
}

func (p *FileSaveClient) Gamma() (r, g, b float64) {
	return 1.0, 1.0, 1.0
}

func (p *FileSaveClient) SetGamma(r, g, b float64) {}

/*
func (p *FileSaveClient) MaxBright() (r, g, b uint8) {
	return 0xff, 0xff, 0xff
}

func (p *FileSaveClient) SetMaxBright(r, g, b uint8) {}
*/

func (p *FileSaveClient) ModuleConfig() conf.ModuleConfig {
	return p.modConf
}

func (p *FileSaveClient) Watch() *Stopwatch {
	return p.sendWatch
}

func (p *FileSaveClient) Close() {
	p.fh.Close()
}
