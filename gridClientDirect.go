package ledgrid

import (
	"github.com/stefan-muehlebach/ledgrid/conf"
)

type DirectGridClient struct {
	Disp      Displayer
	sendWatch *Stopwatch
}

func NewDirectGridClient(Disp Displayer) GridClient {
	c := &DirectGridClient{}

	c.Disp = Disp
	c.sendWatch = NewStopwatch()
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

func (c *DirectGridClient) Watch() *Stopwatch {
	return c.sendWatch
}

func (c *DirectGridClient) Close() {
	c.Disp.Close()
}
