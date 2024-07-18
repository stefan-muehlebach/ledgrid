package main

import (
	"github.com/stefan-muehlebach/ledgrid/colornames"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/stefan-muehlebach/ledgrid"
)

type Cell struct {
	Color ledgrid.LedColor
    Alpha float64
}

func (c *Cell) LedColor() ledgrid.LedColor {
    r := uint8(c.Alpha * float64(c.Color.R))
    g := uint8(c.Alpha * float64(c.Color.G))
    b := uint8(c.Alpha * float64(c.Color.B))
    a := uint8(c.Alpha * 255.0)
    return ledgrid.LedColor{r, g, b, a}
}

type Grid struct {
	Cells                            []*Cell
	AnimList                         []Animation
	width, height                    int
	pixCtrl                          ledgrid.PixelClient
	ledGrid                          *ledgrid.LedGrid
	ticker                           *time.Ticker
	quit                             bool
	animMutex                        *sync.Mutex
	animPit                          time.Time
	logFile                          io.Writer
	animWatch, paintWatch, sendWatch *Stopwatch
}

func NewGrid(pixCtrl ledgrid.PixelClient, ledGrid *ledgrid.LedGrid) *Grid {
	var err error

	c := &Grid{}
	c.pixCtrl = pixCtrl
	c.ledGrid = ledGrid
    c.width = ledGrid.Rect.Dx()
    c.height = ledGrid.Rect.Dy()
	c.Cells = make([]*Cell, c.width * c.height)
	for idx := range c.Cells {
		c.Cells[idx] = &Cell{Color: colornames.Red, Alpha: 1.0}
	}
	c.ticker = time.NewTicker(refreshRate)
	c.AnimList = make([]Animation, 0)
	c.animMutex = &sync.Mutex{}
	if doLog {
		c.logFile, err = os.Create("grid.log")
		if err != nil {
			log.Fatalf("Couldn't create logfile: %v", err)
		}
	}
	c.animWatch = NewStopwatch()
	c.paintWatch = NewStopwatch()
	c.sendWatch = NewStopwatch()
	go c.backgroundThread()
	return c
}

func (c *Grid) Close() {
	c.DelAllAnim()
	c.quit = true
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
	for i, a := range c.AnimList {
		if a == anim {
			c.AnimList[i] = nil
			return
		}
	}
	c.animMutex.Unlock()
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
	numCores := runtime.NumCPU()
	animChan := make(chan int, queueSize)
	doneChan := make(chan bool, queueSize)

	for range numCores {
		go c.animationThread(animChan, doneChan)
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
		numAnims := 0
		for i, anim := range c.AnimList {
			if anim == nil || anim.IsStopped() {
				continue
			}
			numAnims++
			animChan <- i
		}
		for range numAnims {
			<-doneChan
		}
		c.animWatch.Stop()

		c.paintWatch.Start()
		for y := range c.height {
			for x := range c.width {
                idx := y*c.width + x
				c.ledGrid.SetLedColor(x, y, c.Cells[idx].LedColor())
			}
		}
		c.paintWatch.Stop()

		c.sendWatch.Start()
		c.pixCtrl.Draw(c.ledGrid)
		c.sendWatch.Stop()
	}
	close(doneChan)
	close(animChan)
}

func (c *Grid) animationThread(animChan <-chan int, doneChan chan<- bool) {
	for animId := range animChan {
		c.AnimList[animId].Update(c.animPit)
		// if !c.AnimList[animId].Update(c.animPit) {
		//     c.AnimList[animId] = nil
		// }
		doneChan <- true
	}
}
