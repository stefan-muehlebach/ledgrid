package main

import (
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"runtime"
	"slices"
	"sync"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/colornames"
)

type Grid struct {
	ObjList                          []GridObject
	AnimList                         []Animation
	width, height                    int
	pixCtrl                          ledgrid.PixelClient
	ledGrid                          *ledgrid.LedGrid
	ticker                           *time.Ticker
	quit                             bool
	objMutex                         *sync.RWMutex
	animMutex                        *sync.RWMutex
	animPit                          time.Time
	logFile                          io.Writer
	animWatch, paintWatch, sendWatch *ledgrid.Stopwatch
	numThreads                       int
}

func NewGrid(pixCtrl ledgrid.PixelClient, ledGrid *ledgrid.LedGrid) *Grid {
	var err error

	c := &Grid{}
	c.pixCtrl = pixCtrl
	c.ledGrid = ledGrid
	c.width = ledGrid.Rect.Dx()
	c.height = ledGrid.Rect.Dy()
	c.ObjList = make([]GridObject, 0)
	c.objMutex = &sync.RWMutex{}
	c.ticker = time.NewTicker(refreshRate)
	c.AnimList = make([]Animation, 0)
	c.animMutex = &sync.RWMutex{}
	if doLog {
		c.logFile, err = os.Create("grid.log")
		if err != nil {
			log.Fatalf("Couldn't create logfile: %v", err)
		}
	}
	c.animWatch = ledgrid.NewStopwatch()
	c.paintWatch = ledgrid.NewStopwatch()
	c.sendWatch = ledgrid.NewStopwatch()
	go c.backgroundThread()
	return c
}

func (c *Grid) Close() {
	c.DelAllAnim()
	c.DelAll()
	c.quit = true
}

// Fuegt der Zeichenflaeche weitere Objekte hinzu. Der Zufgriff auf den
// entsprechenden Slice wird nicht synchronisiert.
func (c *Grid) Add(objs ...GridObject) {
	c.objMutex.Lock()
	c.ObjList = append(c.ObjList, objs...)
	c.objMutex.Unlock()
}

// Loescht alle Objekte von der Zeichenflaeche.
func (c *Grid) DelAll() {
	c.objMutex.Lock()
	c.ObjList = c.ObjList[:0]
	c.objMutex.Unlock()
}

// Fuegt weitere Animationen hinzu. Der Zugriff auf den entsprechenden Slice
// wird synchronisiert, da die Bearbeitung der Animationen durch den
// Background-Thread ebenfalls relativ haeufig auf den Slice zugreift.
func (c *Grid) AddAnim(anims ...Animation) {
	c.animMutex.Lock()
	c.AnimList = append(c.AnimList, anims...)
	c.animMutex.Unlock()
}

// Loescht alle Animationen.
func (c *Grid) DelAllAnim() {
	c.animMutex.Lock()
	c.AnimList = c.AnimList[:0]
	c.animMutex.Unlock()
}

// Loescht eine einzelne Animation.
func (c *Grid) DelAnim(anim Animation) {
	c.animMutex.Lock()
	defer c.animMutex.Unlock()
	for idx, obj := range c.AnimList {
		if obj == anim {
			c.AnimList = slices.Delete(c.AnimList, idx, idx+1)
			return
		}
	}
}

// Mit Stop koennen die Animationen und die Darstellung auf der Hardware
// unterbunden werden.
func (c *Grid) Stop() {
	c.ticker.Stop()
}

// Setzt die Animationen wieder fort.
// TO DO: Die Fortsetzung sollte fuer eine:n Beobachter:in nahtlos erfolgen.
// Im Moment tut es das nicht - man muesste sich bei den Methoden und Ideen
// von AnimationEmbed bedienen.
func (c *Grid) Continue() {
	c.ticker.Reset(refreshRate)
}

// Hier sind wichtige aber private Methoden, darum in Kleinbuchstaben und
// darum noch sehr wenig Kommentare.
func (c *Grid) backgroundThread() {
	c.numThreads = runtime.NumCPU()
	startChan := make(chan int)
	doneChan := make(chan bool)

	for range c.numThreads {
		go c.animationUpdater(startChan, doneChan)
	}

	lastPit := time.Now()
	for c.animPit = range c.ticker.C {
		if doLog {
			delay := c.animPit.Sub(lastPit)
			lastPit = c.animPit
			fmt.Fprintf(c.logFile, "delay: %v\n", delay)
		}
		if c.quit {
			break
		}

		c.animWatch.Start()
		for id := range c.numThreads {
			startChan <- id
		}
		for range c.numThreads {
			<-doneChan
		}
		c.animWatch.Stop()

		c.paintWatch.Start()
		c.ledGrid.Clear(colornames.Black)
		c.objMutex.RLock()
		for _, obj := range c.ObjList {
			obj.Draw(c.ledGrid)
		}
		c.objMutex.RUnlock()
		c.paintWatch.Stop()

		c.sendWatch.Start()
		c.pixCtrl.Draw(c.ledGrid)
		c.sendWatch.Stop()
	}
	close(doneChan)
	close(startChan)
}

func (c *Grid) animationUpdater(startChan <-chan int, doneChan chan<- bool) {
	for id := range startChan {
		c.animMutex.RLock()
		for i := id; i < len(c.AnimList); i += c.numThreads {
			anim := c.AnimList[i]
			if anim == nil || anim.IsStopped() {
				continue
			}
			anim.Update(c.animPit)
		}
		c.animMutex.RUnlock()
		doneChan <- true
	}
}

// Alle Objekte, die durch den Controller auf dem LED-Grid dargestellt werden
// sollen, muessen im Minimum die Methode Draw implementieren, durch welche
// sie auf einem gg-Kontext gezeichnet werden.
type GridObject interface {
	Draw(lg *ledgrid.LedGrid)
}

// Will man ein einzelnes Pixel zeichnen, so eignet sich dieser Typ. Er wird
// ueber die Zeichenfunktion DrawPoint im gg-Kontext realisiert und hat einen
// Radius von 0.5*sqrt(2).
type GridPixel struct {
	Pos   image.Point
	Color ledgrid.LedColor
}

func NewGridPixel(pos image.Point, col ledgrid.LedColor) *GridPixel {
	p := &GridPixel{Pos: pos, Color: col}
	return p
}

func (p *GridPixel) Draw(lg *ledgrid.LedGrid) {
	lg.SetLedColor(p.Pos.X, p.Pos.Y, p.Color.Mix(lg.LedColorAt(p.Pos.X, p.Pos.Y), ledgrid.Blend))
}

// Fuer das direkte Zeichnen von Text auf dem LED-Grid, existieren einige
// 'fixed size' Bitmap-Schriften, die ohne Rastern und Rendern sehr schnell
// dargestellt werden koennen.
type GridText struct {
	// Pos   image.Point
	Pos    fixed.Point26_6
	Color  ledgrid.LedColor
	text   string
	drawer *font.Drawer
	rect   fixed.Rectangle26_6
	dp     fixed.Point26_6
}

func NewGridText(pos fixed.Point26_6, col ledgrid.LedColor, text string) *GridText {
	t := &GridText{Pos: pos, Color: col}
	t.drawer = &font.Drawer{
		Face: ledgrid.Face3x5,
	}
	t.SetText(text)
	return t
}

func (t *GridText) Text() string {
	return t.text
}

func (t *GridText) SetText(text string) {
	t.text = text
	t.rect, _ = t.drawer.BoundString(text)
	t.dp = t.rect.Min.Add(t.rect.Max).Div(fixed.I(2))
}

func (t *GridText) Draw(lg *ledgrid.LedGrid) {
	t.drawer.Dst = lg
	t.drawer.Src = image.NewUniform(t.Color)
	t.drawer.Dot = t.Pos.Sub(t.dp)
	t.drawer.DrawString(t.text)
}
