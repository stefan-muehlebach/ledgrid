package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"slices"
	"sync"
	"time"

	"golang.org/x/image/font"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/fonts"
	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/colornames"
	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
	"golang.org/x/image/draw"
)

var (
	// Damit wird die Groesse der Queues dimensioniert, welche zu und von
	// den Hintergrundprozessen fuehren.
	queueSize = 400

	// Alle Positionsdaten werden bei der Konvertierung um diesem Wert
	// verschoben. Da gg mit Fliesskommazahlen arbeitet, treffen Koordinaten
	// wie (1,5) nie direkt auf ein Pixel, sondern immer dazwischen.
	displ = geom.Point{0.5, 0.5}

	// Mit oversize wird ein Vergroesserungsfaktor beschrieben, der fuer alle
	// Zeichenoperationen verwendet wird. Damit wird ein insgesamt weicheres
	// Bild erzielt.
	oversize = 10.0

	// Ueber die NewXXX-Funktionen koennen die Objekte einfacher erzeugt
	// werden. Die Fuellfarbe ist gleich der Randfarbe, hat aber einen
	// niedrigeren Alpha-Wert, der mit dieser Konstante definiert werden
	// kann.
	fillAlpha = 0.4

	doLog = false
)

// Ein Canvas ist eine animierbare Zeichenflaeche. Ihr koennen eine beliebige
// Anzahl von zeichenbaren Objekten (Interface CanvasObject) hinzugefuegt
// werden.
type Canvas struct {
	ObjList                          []CanvasObject
	AnimList                         []Animation
	objMutex                         *sync.RWMutex
	animMutex                        *sync.RWMutex
	pixCtrl                          ledgrid.PixelClient
	ledGrid                          *ledgrid.LedGrid
	canvas                           *image.RGBA
	gc                               *gg.Context
	scaler                           draw.Scaler
	ticker                           *time.Ticker
	quit                             bool
	animPit                          time.Time
	logFile                          io.Writer
	animWatch, paintWatch, sendWatch *ledgrid.Stopwatch
	numThreads                       int
}

func NewCanvas(pixCtrl ledgrid.PixelClient, ledGrid *ledgrid.LedGrid) *Canvas {
	var err error

	c := &Canvas{}
	c.pixCtrl = pixCtrl
	c.ledGrid = ledGrid
	c.canvas = image.NewRGBA(image.Rectangle{Max: c.ledGrid.Rect.Max.Mul(int(oversize))})
	c.gc = gg.NewContextForRGBA(c.canvas)
	c.ObjList = make([]CanvasObject, 0)
	c.objMutex = &sync.RWMutex{}
	c.scaler = draw.CatmullRom.NewScaler(c.ledGrid.Rect.Dx(), c.ledGrid.Rect.Dy(), c.canvas.Rect.Dx(), c.canvas.Rect.Dy())
	c.ticker = time.NewTicker(refreshRate)
	c.AnimList = make([]Animation, 0)
	c.animMutex = &sync.RWMutex{}
	if doLog {
		c.logFile, err = os.Create("canvas.log")
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

func (c *Canvas) Close() {
	c.DelAllAnim()
	c.DelAll()
	c.quit = true
}

// Fuegt der Zeichenflaeche weitere Objekte hinzu. Der Zufgriff auf den
// entsprechenden Slice wird nicht synchronisiert.
func (c *Canvas) Add(objs ...CanvasObject) {
	c.objMutex.Lock()
	c.ObjList = append(c.ObjList, objs...)
	c.objMutex.Unlock()
}

// Loescht alle Objekte von der Zeichenflaeche.
func (c *Canvas) DelAll() {
	c.objMutex.Lock()
	c.ObjList = c.ObjList[:0]
	c.objMutex.Unlock()
}

// Fuegt weitere Animationen hinzu. Der Zugriff auf den entsprechenden Slice
// wird synchronisiert, da die Bearbeitung der Animationen durch den
// Background-Thread ebenfalls relativ haeufig auf den Slice zugreift.
func (c *Canvas) AddAnim(anims ...Animation) {
	c.animMutex.Lock()
	c.AnimList = append(c.AnimList, anims...)
	c.animMutex.Unlock()
}

// Loescht eine einzelne Animation.
func (c *Canvas) DelAnim(anim Animation) {
	c.animMutex.Lock()
	defer c.animMutex.Unlock()
	for idx, obj := range c.AnimList {
		if obj == anim {
			c.AnimList = slices.Delete(c.AnimList, idx, idx+1)
			return
		}
	}
}

// Loescht alle Animationen.
func (c *Canvas) DelAllAnim() {
	c.animMutex.Lock()
	c.AnimList = c.AnimList[:0]
	c.animMutex.Unlock()
}

// Mit Stop koennen die Animationen und die Darstellung auf der Hardware
// unterbunden werden.
func (c *Canvas) Stop() {
	c.ticker.Stop()
}

// Setzt die Animationen wieder fort.
// TO DO: Die Fortsetzung sollte fuer eine:n Beobachter:in nahtlos erfolgen.
// Im Moment tut es das nicht - man muesste sich bei den Methoden und Ideen
// von AnimationEmbed bedienen.
func (c *Canvas) Continue() {
	c.ticker.Reset(refreshRate)
}

// Mit den folgenden 4 Methoden verfolge ich das ambitionierte Ziel, die
// Animationen in irgendeiner Form serialisierbar zu machen, damit in ferner
// Zukunft die Animationen vollstaendig auf den Rechner des Pixelcontrollers
// verlegt werden koennen und das netzwerkbedingte "Ruckeln" der
// Vergangenheit angehoert.

func (c *Canvas) Save(fileName string) {
	fh, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Couldn't create file: %v", err)
	}
	c.Write(fh)
	fh.Close()
}

func (c *Canvas) Write(w io.Writer) {
	gobEncoder := gob.NewEncoder(w)
	err := gobEncoder.Encode(c)
	if err != nil {
		log.Fatalf("Couldn't encode data: %v", err)
	}
}

func (c *Canvas) Load(fileName string) {
	fh, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Couldn't create file: %v", err)
	}
	c.Read(fh)
	fh.Close()
}

func (c *Canvas) Read(r io.Reader) {
	gobDecoder := gob.NewDecoder(r)
	err := gobDecoder.Decode(c)
	if err != nil {
		log.Fatalf("Couldn't decode data: %v", err)
	}
}

// Hier sind wichtige aber private Methoden, darum in Kleinbuchstaben und
// darum noch sehr wenig Kommentare.
func (c *Canvas) backgroundThread() {
	// backColor := colornames.Black
	c.numThreads = runtime.NumCPU()
	startChan := make(chan int) //, queueSize)
	doneChan := make(chan bool) //, queueSize)
	// numCores := runtime.NumCPU()
	// animChan := make(chan int, queueSize)
	// doneChan := make(chan bool, queueSize)

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
		c.gc.SetFillColor(colornames.Black)
		c.gc.Clear()
		c.objMutex.RLock()
		for _, obj := range c.ObjList {
			obj.Draw(c)
		}
		c.objMutex.RUnlock()
		c.scaler.Scale(c.ledGrid, c.ledGrid.Rect, c.canvas, c.canvas.Rect, draw.Over, nil)
		c.paintWatch.Stop()

		c.sendWatch.Start()
		c.pixCtrl.Draw(c.ledGrid)
		c.sendWatch.Stop()
	}
	close(doneChan)
	close(startChan)
}

func (c *Canvas) animationUpdater(startChan <-chan int, doneChan chan<- bool) {
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

// Damit werden die jeweiligen Graphik-Objekte beim Package gob registriert,
// um sie binaer zu exportieren.
func init() {
	gob.Register(ledgrid.LedColor{})

	gob.Register(&Ellipse{})
	gob.Register(&Rectangle{})
	gob.Register(&Line{})
	gob.Register(&Pixel{})
}

// Mit ConvertPos muessen alle Positionsdaten konvertiert werden.
func ConvertPos(p geom.Point) geom.Point {
	return p.Add(displ).Mul(oversize)
}

// ConvertSize dagegen wird fuer die Konvertierung aller Groessenangaben
// verwendet.
func ConvertSize(s geom.Point) geom.Point {
	return s.Mul(oversize)
}

// Einzelne Laengen werden mit ConvertLen konvertiert.
func ConvertLen(l float64) float64 {
	return l * oversize
}

type ColorConvertFunc func(ledgrid.LedColor) ledgrid.LedColor

func ApplyAlpha(c ledgrid.LedColor) ledgrid.LedColor {
	alpha := float64(c.A) / 255.0
	return c.Alpha(alpha * fillAlpha)
}

// Alle Objekte, die durch den Controller auf dem LED-Grid dargestellt werden
// sollen, muessen im Minimum die Methode Draw implementieren, durch welche
// sie auf einem gg-Kontext gezeichnet werden.
type CanvasObject interface {
	Draw(c *Canvas)
}

// Dient dazu, ein Live-Bild ab einer beliebigen, aber ansprechbaren Kamera
// auf dem LED-Grid darzustellen.
const (
	camDevName    = "/dev/video0"
	camWidth      = 320
	camHeight     = 240
	camFrameRate  = 30
	camBufferSize = 4
)

type Camera struct {
	Pos, Size geom.Point
	dev       *device.Device
	img       image.Image
	cut       image.Rectangle
	cancel    context.CancelFunc
	running   bool
}

func NewCamera(pos, size geom.Point) *Camera {
	c := &Camera{Pos: pos, Size: size, cut: image.Rect(0, 80, 320, 160)}
	AnimCtrl.AddAnim(c)
	return c
}

func (c *Camera) Duration() time.Duration {
	return time.Duration(0)
}

func (c *Camera) SetDuration(dur time.Duration) {}

func (c *Camera) Start() {
	var ctx context.Context
	var err error

	if c.running {
		return
	}
	c.dev, err = device.Open(camDevName,
		device.WithIOType(v4l2.IOTypeMMAP),
		device.WithPixFormat(v4l2.PixFormat{
			PixelFormat: v4l2.PixelFmtMJPEG,
			Width:       camWidth,
			Height:      camHeight,
		}),
		device.WithFPS(camFrameRate),
		device.WithBufferSize(camBufferSize),
	)
	if err != nil {
		log.Fatalf("failed to open device: %v", err)
	}
	ctx, c.cancel = context.WithCancel(context.TODO())
	if err = c.dev.Start(ctx); err != nil {
		log.Fatalf("failed to start stream: %v", err)
	}
	c.running = true
}

func (c *Camera) Stop() {
	var err error

	if !c.running {
		return
	}
	c.cancel()
	if err = c.dev.Close(); err != nil {
		log.Fatalf("failed to close device: %v", err)
	}
	c.dev = nil
	c.running = false
}

func (c *Camera) Continue() {}

func (c *Camera) IsStopped() bool {
	return !c.running
}

func (c *Camera) Update(pit time.Time) bool {
	var err error
	var frame []byte
	var ok bool

	if frame, ok = <-c.dev.GetOutput(); !ok {
		log.Printf("no frame to process")
		return true
	}
	reader := bytes.NewReader(frame)
	c.img, err = jpeg.Decode(reader)
	if err != nil {
		log.Fatalf("failed to decode data: %v", err)
	}
	return true
}

func (c *Camera) Draw(canv *Canvas) {
	if c.img == nil {
		return
	}
	rect := geom.Rectangle{Max: c.Size}
	refPt := c.Pos.Sub(c.Size.Div(2.0))
	draw.CatmullRom.Scale(canv.canvas, rect.Add(refPt).Int(), c.img, c.cut, draw.Over, nil)
}

// Zur Darstellung von beliebigen Bildern (JPEG, PNG, etc) auf dem LED-Panel
// Da es nur wenige LEDs zur Darstellung hat, werden die Bilder gnadenlos
// skaliert und herunter gerechnet - manchmal bis der Arzt kommt... ;-)
type Image struct {
	Pos, Size geom.Point
	Angle     float64
	img       image.Image
}

func NewImage(pos geom.Point, fileName string) *Image {
	i := &Image{Pos: pos, Angle: 0.0}
	fh, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Couldn't open file: %v", err)
	}
	defer fh.Close()
	i.img, _, err = image.Decode(fh)
	if err != nil {
		log.Fatalf("Couldn't decode image: %v", err)
	}
	i.Size = geom.NewPointIMG(i.img.Bounds().Size().Mul(int(oversize)))
	return i
}

func (i *Image) Draw(c *Canvas) {
	rect := geom.Rectangle{Max: i.Size}
	refPt := i.Pos.Sub(i.Size.Div(2.0))
	draw.CatmullRom.Scale(c.canvas, rect.Add(refPt).Int(), i.img, i.img.Bounds(), draw.Over, nil)
}

// Zur Darstellung von beliebigem Text.
var (
	defFont     = fonts.SeafordBold
	defFontSize = ConvertLen(12.0)
)

type Text struct {
	Pos      geom.Point
	Angle    float64
	Color    ledgrid.LedColor
	Font     *fonts.Font
	FontSize float64
	Text     string
	fontFace font.Face
}

func NewText(pos geom.Point, text string, color ledgrid.LedColor) *Text {
	t := &Text{Pos: pos, Color: color, Font: defFont, FontSize: defFontSize,
		Text: text}
	t.fontFace = fonts.NewFace(t.Font, t.FontSize)
	return t
}

func (t *Text) Draw(c *Canvas) {
	if t.Angle != 0.0 {
		c.gc.Push()
		c.gc.RotateAbout(t.Angle, t.Pos.X, t.Pos.Y)
		defer c.gc.Pop()
	}
	c.gc.SetStrokeColor(t.Color)
	c.gc.SetFontFace(t.fontFace)
	c.gc.DrawStringAnchored(t.Text, t.Pos.X, t.Pos.Y, 0.5, 0.5)
}

// Mit Ellipse sind alle kreisartigen Objekte abgedeckt. Pos bezeichnet die
// Position des Mittelpunktes und mit Size ist die Breite und Hoehe des
// gesamten Objektes gemeint. Falls ein Rand gezeichnet werden soll, muss
// BorderWith einen Wert >0 enthalten und FillColor, resp. BorderColor
// enthalten die Farben fuer Rand und Flaeche.
type Ellipse struct {
	Pos, Size              geom.Point
	Angle                  float64
	BorderWidth            float64
	BorderColor, FillColor ledgrid.LedColor
	FillColorFnc           ColorConvertFunc
}

// Erzeugt eine 'klassische' Ellipse mit einer Randbreite von einem Pixel und
// setzt die Fuellfarbe gleich Randfarbe mit Alpha-Wert von 0.3.
// Will man die einzelnen Werte flexibler verwenden, empfiehlt sich die
// Erzeugung mittels &Ellipse{...}.
func NewEllipse(pos, size geom.Point, borderColor ledgrid.LedColor) *Ellipse {
	e := &Ellipse{Pos: pos, Size: size, BorderWidth: ConvertLen(1.0),
		BorderColor: borderColor, FillColorFnc: ApplyAlpha}
	e.FillColor = ledgrid.Transparent
	return e
}

func (e *Ellipse) Draw(c *Canvas) {
	if e.Angle != 0.0 {
		c.gc.Push()
		c.gc.RotateAbout(e.Angle, e.Pos.X, e.Pos.Y)
		defer c.gc.Pop()
	}
	c.gc.DrawEllipse(e.Pos.X, e.Pos.Y, e.Size.X/2, e.Size.Y/2)
	c.gc.SetStrokeWidth(e.BorderWidth)
	c.gc.SetStrokeColor(e.BorderColor)
	if e.FillColor == ledgrid.Transparent {
		c.gc.SetFillColor(e.FillColorFnc(e.BorderColor))
	} else {
		c.gc.SetFillColor(e.FillColor)
	}
	c.gc.FillStroke()
}

// Rectangle ist fuer alle rechteckigen Objekte vorgesehen. Pos bezeichnet
// den Mittelpunkt des Objektes und Size die Breite, rsep. Hoehe.
type Rectangle struct {
	Pos, Size              geom.Point
	Angle                  float64
	BorderWidth            float64
	BorderColor, FillColor ledgrid.LedColor
	FillColorFnc           ColorConvertFunc
}

func NewRectangle(pos, size geom.Point, borderColor ledgrid.LedColor) *Rectangle {
	r := &Rectangle{Pos: pos, Size: size, BorderWidth: ConvertLen(1.0),
		BorderColor: borderColor, FillColorFnc: ApplyAlpha}
	return r
}

func (r *Rectangle) Draw(c *Canvas) {
	if r.Angle != 0.0 {
		c.gc.Push()
		c.gc.RotateAbout(r.Angle, r.Pos.X, r.Pos.Y)
		defer c.gc.Pop()
	}
	c.gc.DrawRectangle(r.Pos.X-r.Size.X/2, r.Pos.Y-r.Size.Y/2, r.Size.X, r.Size.Y)
	c.gc.SetStrokeWidth(r.BorderWidth)
	c.gc.SetStrokeColor(r.BorderColor)
	if r.FillColorFnc != nil {
		c.gc.SetFillColor(r.FillColorFnc(r.BorderColor))
	} else {
		c.gc.SetFillColor(r.FillColor)
	}
	c.gc.FillStroke()
}

// Auch gleichmaessige Polygone duerfen nicht fehlen.
type RegularPolygon struct {
	Pos, Size              geom.Point
	Angle                  float64
	BorderWidth            float64
	BorderColor, FillColor ledgrid.LedColor
	FillColorFnc           ColorConvertFunc
	numPoints              int
}

func NewRegularPolygon(numPoints int, pos, size geom.Point, borderColor ledgrid.LedColor) *RegularPolygon {
	p := &RegularPolygon{Pos: pos, Size: size, Angle: 0.0, BorderWidth: ConvertLen(1.0),
		BorderColor: borderColor, FillColorFnc: ApplyAlpha, numPoints: numPoints}
	return p
}

func (p *RegularPolygon) Draw(c *Canvas) {
	c.gc.DrawRegularPolygon(p.numPoints, p.Pos.X, p.Pos.Y, p.Size.X/2.0, p.Angle)
	c.gc.SetStrokeWidth(p.BorderWidth)
	c.gc.SetStrokeColor(p.BorderColor)
	if p.FillColorFnc != nil {
		c.gc.SetFillColor(p.FillColorFnc(p.BorderColor))
	} else {
		c.gc.SetFillColor(p.FillColor)
	}
	c.gc.FillStroke()
}

// Fuer Geraden ist dieser Datentyp vorgesehen, der von Pos1 nach Pos2
// verlaeuft.
type Line struct {
	Pos1, Pos2 geom.Point
	Width      float64
	Color      ledgrid.LedColor
}

func NewLine(pos1, pos2 geom.Point, col ledgrid.LedColor) *Line {
	l := &Line{Pos1: pos1, Pos2: pos2, Width: ConvertLen(1.0), Color: col}
	return l
}

func (l *Line) Draw(c *Canvas) {
	c.gc.SetStrokeWidth(l.Width)
	c.gc.SetStrokeColor(l.Color)
	c.gc.DrawLine(l.Pos1.X, l.Pos1.Y, l.Pos2.X, l.Pos2.Y)
	c.gc.Stroke()
}

// Will man ein einzelnes Pixel zeichnen, so eignet sich dieser Typ. Er wird
// ueber die Zeichenfunktion DrawPoint im gg-Kontext realisiert und hat einen
// Radius von 0.5*sqrt(2).
type Pixel struct {
	Pos   geom.Point
	Color ledgrid.LedColor
}

func NewPixel(pos geom.Point, col ledgrid.LedColor) *Pixel {
	p := &Pixel{Pos: pos, Color: col}
	return p
}

func (p *Pixel) Draw(c *Canvas) {
	c.gc.SetFillColor(p.Color)
	c.gc.DrawPoint(p.Pos.X, p.Pos.Y, ConvertLen(0.5*math.Sqrt2))
	c.gc.Fill()
}
