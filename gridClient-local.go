package ledgrid

import (
	"github.com/stefan-muehlebach/ledgrid/conf"
)

type DirectGridClient struct {
	Disp      Displayer
	stopwatch *Stopwatch
}

func NewDirectGridClient(Disp Displayer) GridClient {
	c := &DirectGridClient{}

	c.Disp = Disp
	c.stopwatch = NewStopwatch()
	return c
}

func (c *DirectGridClient) Send(buffer []byte) {
	c.Disp.Display(buffer)
}

func (c *DirectGridClient) NumLeds() int {
	return c.Disp.NumLeds()
}

func (c *DirectGridClient) Gamma() (r, g, b float64) {
	return c.Disp.Gamma()
}

func (c *DirectGridClient) SetGamma(r, g, b float64) {
	c.Disp.SetGamma(r, g, b)
}

func (c *DirectGridClient) ModuleConfig() conf.ModuleConfig {
	return c.Disp.ModuleConfig()
}

func (c *DirectGridClient) Stopwatch() *Stopwatch {
	return c.stopwatch
}

func (c *DirectGridClient) Close() {
	c.Disp.Close()
}
