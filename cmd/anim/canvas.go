package main

import (
	"bytes"
	"context"
	"encoding/xml"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"gocv.io/x/gocv"

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

type DrawingArea interface {
	Add(objs ...DrawableObject)
	DelAll()
	Refresh()
}

type DrawableObject interface {
	Draw(d DrawingArea)
}

// Ein Canvas ist eine animierbare Zeichenflaeche. Ihr koennen eine beliebige
// Anzahl von zeichenbaren Objekten (Interface CanvasObject) hinzugefuegt
// werden.
type Canvas struct {
	ObjList               []CanvasObject
	objMutex              *sync.RWMutex
	pixCtrl               ledgrid.PixelClient
	ledGrid               *ledgrid.LedGrid
	canvas                *image.RGBA
	gc                    *gg.Context
	scaler                draw.Scaler
	logFile               io.Writer
	paintWatch, sendWatch *ledgrid.Stopwatch
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
	if doLog {
		c.logFile, err = os.Create("canvas.log")
		if err != nil {
			log.Fatalf("Couldn't create logfile: %v", err)
		}
	}
	c.paintWatch = ledgrid.NewStopwatch()
	c.sendWatch = ledgrid.NewStopwatch()
	return c
}

func (c *Canvas) Close() {
	c.DelAll()
}

// Fuegt der Zeichenflaeche weitere Objekte hinzu. Der Zufgriff auf den
// entsprechenden Slice wird nicht synchronisiert.
func (c *Canvas) Add(objs ...DrawableObject) {
	c.objMutex.Lock()
	for _, tmp := range objs {
		obj := tmp.(CanvasObject)
		c.ObjList = append(c.ObjList, obj)
	}
	c.objMutex.Unlock()
}

// Loescht alle Objekte von der Zeichenflaeche.
func (c *Canvas) DelAll() {
	c.objMutex.Lock()
	c.ObjList = c.ObjList[:0]
	c.objMutex.Unlock()
}

func (c *Canvas) Refresh() {
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
	Draw(d DrawingArea)
}

// Dient dazu, ein Live-Bild ab einer beliebigen, aber ansprechbaren Kamera
// auf dem LED-Grid darzustellen. Als erstes eine Implementation mit Hilfe
// der Video4Linux-Umgebung.
const (
	camDevName    = "/dev/video0"
	camDevId      = 0
	camWidth      = 320
	camHeight     = 240
	camFrameRate  = 30
	camBufferSize = 4
)

type CameraV4L struct {
	Pos, Size geom.Point
	dev       *device.Device
	img       image.Image
	cut       image.Rectangle
	cancel    context.CancelFunc
	running   bool
}

func NewCameraV4L(pos, size geom.Point) *CameraV4L {
	c := &CameraV4L{Pos: pos, Size: size, cut: image.Rect(0, 80, 320, 160)}
	animCtrl.AddAnim(c)
	return c
}

func (c *CameraV4L) Duration() time.Duration {
	return time.Duration(0)
}

func (c *CameraV4L) SetDuration(dur time.Duration) {}

func (c *CameraV4L) Start() {
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

func (c *CameraV4L) Stop() {
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

func (c *CameraV4L) Continue() {}

func (c *CameraV4L) IsStopped() bool {
	return !c.running
}

func (c *CameraV4L) Update(pit time.Time) bool {
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

func (c *CameraV4L) Draw(d DrawingArea) {
	canv := d.(*Canvas)
	if c.img == nil {
		return
	}
	rect := geom.Rectangle{Max: c.Size}
	refPt := c.Pos.Sub(c.Size.Div(2.0))
	draw.CatmullRom.Scale(canv.canvas, rect.Add(refPt).Int(), c.img, c.cut, draw.Over, nil)
}

type CameraCV struct {
	Pos, Size geom.Point
	dev       *gocv.VideoCapture
	img       image.Image
	cut       image.Rectangle
	mat       gocv.Mat
	running   bool
}

func NewCameraCV(pos, size geom.Point) *CameraCV {
	c := &CameraCV{Pos: pos, Size: size, cut: image.Rect(0, 80, 320, 160)}
	c.mat = gocv.NewMat()
	animCtrl.AddAnim(c)
	return c
}

func (c *CameraCV) Duration() time.Duration {
	return time.Duration(0)
}

func (c *CameraCV) SetDuration(dur time.Duration) {}

func (c *CameraCV) Start() {
	var err error

	if c.running {
		return
	}
	c.dev, err = gocv.VideoCaptureDevice(camDevId)
	if err != nil {
		log.Fatalf("Couldn't open device: %v", err)
	}
	c.dev.Set(gocv.VideoCaptureFrameWidth, camWidth)
	c.dev.Set(gocv.VideoCaptureFrameHeight, camHeight)
	c.running = true
}

func (c *CameraCV) Stop() {
	var err error

	if !c.running {
		return
	}
	err = c.dev.Close()
	if err != nil {
		log.Fatalf("failed to close device: %v", err)
	}
	c.dev = nil
	c.running = false
}

func (c *CameraCV) Continue() {}

func (c *CameraCV) IsStopped() bool {
	return !c.running
}

func (c *CameraCV) Update(pit time.Time) bool {
	var err error

	if !c.dev.Read(&c.mat) {
		log.Fatal("Device closed")
	}
	c.img, err = c.mat.ToImage()
	if err != nil {
		log.Fatalf("Couldn't convert image: %v", err)
	}
	return true
}

func (c *CameraCV) Draw(d DrawingArea) {
	canv := d.(*Canvas)
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
	Img       draw.Image
}

func NewImageFromFile(pos geom.Point, fileName string) *Image {
	var tmp image.Image

	i := &Image{Pos: pos, Angle: 0.0}
	fh, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Couldn't open file: %v", err)
	}
	defer fh.Close()
	tmp, _, err = image.Decode(fh)
	if err != nil {
		log.Fatalf("Couldn't decode image: %v", err)
	}
	i.Img = tmp.(draw.Image)
	i.Size = geom.NewPointIMG(i.Img.Bounds().Size().Mul(int(oversize)))
	return i
}

func (i *Image) Read(fileName string) {
	var img image.Image

	fh, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Couldn't open file: %v", err)
	}
	defer fh.Close()
	img, _, err = image.Decode(fh)
	if err != nil {
		log.Fatalf("Couldn't decode image: %v", err)
	}
	i.Img = img.(draw.Image)
	if i.Size.X > 0 || i.Size.Y > 0 {
		return
	}
	i.Size = geom.NewPointIMG(i.Img.Bounds().Size().Mul(int(oversize)))
}

func (i *Image) Draw(d DrawingArea) {
	c := d.(*Canvas)
	draw.CatmullRom.Scale(c.canvas, geom.NewRectangleWH(i.Pos.X, i.Pos.Y, i.Size.X, i.Size.Y).Int(), i.Img, i.Img.Bounds(), draw.Over, nil)
}

type BlinkenFile struct {
	XMLName  xml.Name       `xml:"blm"`
	Width    int            `xml:"width,attr"`
	Height   int            `xml:"height,attr"`
	Bits     int            `xml:"bits,attr"`
	Channels int            `xml:"channels,attr"`
	Header   BlinkenHeader  `xml:"header"`
	Frames   []BlinkenFrame `xml:"frame"`
}

type BlinkenHeader struct {
	XMLName  xml.Name `xml:"header"`
	Title    string   `xml:"title"`
	Author   string   `xml:"author"`
	Email    string   `xml:"email"`
	Creator  string   `xml:"creator"`
	Duration int      `xml:"duration,omitempty"`
}

type BlinkenFrame struct {
	XMLName  xml.Name  `xml:"frame"`
	Duration int       `xml:"duration,attr"`
	Rows     [][]byte  `xml:"row"`
	Values   [][]uint8 `xml:"-"`
}

func ReadBlinkenFile(fileName string) *BlinkenFile {
	b := &BlinkenFile{Channels: 1}

	xmlFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Couldn't open file '%s': %v", fileName, err)
	}
	defer xmlFile.Close()

	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Fatalf("Couldn't read content of file: %v", err)
	}

	err = xml.Unmarshal(byteValue, b)
	if err != nil {
		log.Fatal(err)
	}

	numberWidth := b.Bits / 4
	if b.Bits%4 != 0 {
		numberWidth++
	}
	for i, frame := range b.Frames {
		b.Frames[i].Values = make([][]uint8, b.Height)
		for j, row := range frame.Rows {
			b.Frames[i].Values[j] = make([]uint8, b.Width*b.Channels)
			for k := 0; k < b.Width; k++ {
				for l := range b.Channels {
					idx := k*numberWidth*b.Channels + l*numberWidth
					val := row[idx : idx+numberWidth]
					v, err := strconv.ParseUint(string(val), 16, b.Bits)
					if err != nil {
						log.Fatalf("Cannot parse '%s': %v", string(val), err)
					}
					idx = k*b.Channels + l
					b.Frames[i].Values[j][idx] = uint8(v)
				}
			}
		}
	}
	return b
}

func (b *BlinkenFile) Image(idx int) *Image {
	var c color.Color

	i := &Image{}
	i.Img = image.NewRGBA(image.Rect(0, 0, b.Width, b.Height))
	colorScale := uint8(255 / ((1 << b.Bits) - 1))
	for row := range b.Height {
		for col := range b.Width {
			idxFrom := col * b.Channels
			idxTo := idxFrom + b.Channels
			src := b.Frames[idx].Values[row][idxFrom:idxTo:idxTo]
			switch b.Channels {
			case 1:
				v := colorScale * src[0]
				if v == 0 {
					c = color.RGBA{0, 0, 0, 0}
				} else {
					c = color.RGBA{v, v, v, 0xff}
				}
			case 3:
				r, g, b := colorScale*src[0], colorScale*src[1], colorScale*src[2]
				if r == 0 && g == 0 && b == 0 {
					c = color.RGBA{0, 0, 0, 0}
				} else {
					c = color.RGBA{r, g, b, 0xff}
				}
			}
			i.Img.Set(col, row, c)
		}
	}
	i.Size = ConvertSize(geom.NewPointIMG(i.Img.Bounds().Size()))
	return i
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

func (t *Text) Draw(draw DrawingArea) {
	c := draw.(*Canvas)
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

func (e *Ellipse) Draw(draw DrawingArea) {
	c := draw.(*Canvas)
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

func (r *Rectangle) Draw(d DrawingArea) {
	c := d.(*Canvas)
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

func (p *RegularPolygon) Draw(draw DrawingArea) {
	c := draw.(*Canvas)
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

func (l *Line) Draw(draw DrawingArea) {
	c := draw.(*Canvas)
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

func (p *Pixel) Draw(draw DrawingArea) {
	c := draw.(*Canvas)
	c.gc.SetFillColor(p.Color)
	c.gc.DrawPoint(p.Pos.X, p.Pos.Y, ConvertLen(0.5*math.Sqrt2))
	c.gc.Fill()
}
