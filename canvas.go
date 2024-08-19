package ledgrid

import (
	"image"
	"log"
	"os"
	"sync"

	"golang.org/x/image/math/f64"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/fonts"
	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid/color"
	"golang.org/x/image/draw"
)

var (
	// Alle Positionsdaten werden bei der Konvertierung um diesem Wert
	// verschoben. Da gg mit Fliesskommazahlen arbeitet, treffen Koordinaten
	// wie (1,5) nie direkt auf ein Pixel, sondern immer dazwischen.
	// displ = geom.Point{0.5, 0.5}

	// Mit oversize wird ein Vergroesserungsfaktor beschrieben, der fuer alle
	// Zeichenoperationen verwendet wird. Damit wird ein insgesamt weicheres
	// Bild erzielt.
	// oversize = 1.0

	// Ueber die NewXXX-Funktionen koennen die Objekte einfacher erzeugt
	// werden. Die Fuellfarbe ist gleich der Randfarbe, hat aber einen
	// niedrigeren Alpha-Wert, der mit dieser Konstante definiert werden
	// kann.
	fillAlpha = 0.4
)

// Ein Canvas ist eine animierbare Zeichenflaeche. Ihr koennen eine beliebige
// Anzahl von zeichenbaren Objekten (Interface CanvasObject) hinzugefuegt
// werden.
type Canvas struct {
	ObjList    []CanvasObject
	objMutex   *sync.RWMutex
	Rect       image.Rectangle
	img        *image.RGBA
	gc         *gg.Context
	paintWatch *Stopwatch
}

func NewCanvas(size image.Point) *Canvas {
	c := &Canvas{}
	c.Rect = image.Rectangle{Max: size}
	c.img = image.NewRGBA(c.Rect)
	c.gc = gg.NewContextForRGBA(c.img)
	c.ObjList = make([]CanvasObject, 0)
	c.objMutex = &sync.RWMutex{}
	c.paintWatch = NewStopwatch()
	return c
}

func (c *Canvas) Close() {
	c.Purge()
}

// Fuegt der Zeichenflaeche weitere Objekte hinzu. Der Zufgriff auf den
// entsprechenden Slice wird nicht synchronisiert.
func (c *Canvas) Add(objs ...CanvasObject) {
	c.objMutex.Lock()
	c.ObjList = append(c.ObjList, objs...)
	c.objMutex.Unlock()
}

// Loescht alle Objekte von der Zeichenflaeche.
func (c *Canvas) Purge() {
	c.objMutex.Lock()
	c.ObjList = c.ObjList[:0]
	c.objMutex.Unlock()
}

func (c *Canvas) Draw(lg draw.Image) {
	c.paintWatch.Start()
	c.gc.SetFillColor(color.Black)
	c.gc.Clear()
	c.objMutex.RLock()
	for _, obj := range c.ObjList {
		if obj.IsHidden() {
			continue
		}
		obj.Draw(c)
	}
	c.objMutex.RUnlock()
	draw.Draw(lg, lg.Bounds(), c.img, image.Point{}, draw.Over)
	c.paintWatch.Stop()
}

func (c *Canvas) Watch() *Stopwatch {
	return c.paintWatch
}

// Alle Objekte, die durch den Controller auf dem LED-Grid dargestellt werden
// sollen, muessen das CanvasObject-Interface implementieren. Dieses
// enthaelt einerseits Methoden zum Ein-/Ausblenden von Objekten und
// andererseits die Methode Draw, mit welcher das CanvasObject auf einer
// Zeichenflaeche gezeichnet werden kann.
type CanvasObject interface {
	Show()
	Hide()
	IsHidden() bool
	Draw(c *Canvas)
}

// Wie bei den Animationen gibt es für die darstellbaren Objekte (CanvasObject)
// ein entsprechendes Embedable, welche die für die meisten Objekte
// brauchbaren Methoden enthält.
type CanvasObjectEmbed struct {
	wrapper CanvasObject
	visible bool
}

func (c *CanvasObjectEmbed) Extend(wrapper CanvasObject) {
	c.visible = true
	c.wrapper = wrapper
}

func (c *CanvasObjectEmbed) Show() {
	if !c.visible {
		c.visible = true
	}
}

func (c *CanvasObjectEmbed) Hide() {
	if c.visible {
		c.visible = false
	}
}

func (c *CanvasObjectEmbed) IsHidden() bool {
	return !c.visible
}

type ColorConvertFunc func(color.LedColor) color.LedColor

func ApplyAlpha(c color.LedColor) color.LedColor {
	alpha := float64(c.A) / 255.0
	return c.Alpha(alpha * fillAlpha)
}

//
// Basic geometric shapes
//

// Mit Ellipse sind alle kreisartigen Objekte abgedeckt. Pos bezeichnet die
// Position des Mittelpunktes und mit Size ist die Breite und Hoehe des
// gesamten Objektes gemeint. Falls ein Rand gezeichnet werden soll, muss
// BorderWith einen Wert >0 enthalten und FillColor, resp. BorderColor
// enthalten die Farben fuer Rand und Flaeche.
type Ellipse struct {
	CanvasObjectEmbed
	Pos, Size              geom.Point
	Angle                  float64
	BorderWidth            float64
	BorderColor, FillColor color.LedColor
	FillColorFnc           ColorConvertFunc
}

// Erzeugt eine 'klassische' Ellipse mit einer Randbreite von einem Pixel und
// setzt die Fuellfarbe gleich Randfarbe mit Alpha-Wert von 0.3.
// Will man die einzelnen Werte flexibler verwenden, empfiehlt sich die
// Erzeugung mittels &Ellipse{...}.
func NewEllipse(pos, size geom.Point, borderColor color.LedColor) *Ellipse {
	e := &Ellipse{Pos: pos, Size: size, BorderWidth: 1.0,
		BorderColor: borderColor, FillColorFnc: ApplyAlpha}
	e.CanvasObjectEmbed.Extend(e)
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
	if e.FillColor == color.Transparent {
		c.gc.SetFillColor(e.FillColorFnc(e.BorderColor))
	} else {
		c.gc.SetFillColor(e.FillColor)
	}
	c.gc.FillStroke()
}

// Rectangle ist fuer alle rechteckigen Objekte vorgesehen. Pos bezeichnet
// den Mittelpunkt des Objektes und Size die Breite, rsep. Hoehe.
type Rectangle struct {
	CanvasObjectEmbed
	Pos, Size              geom.Point
	Angle                  float64
	BorderWidth            float64
	BorderColor, FillColor color.LedColor
	FillColorFnc           ColorConvertFunc
}

func NewRectangle(pos, size geom.Point, borderColor color.LedColor) *Rectangle {
	r := &Rectangle{Pos: pos, Size: size, BorderWidth: 1.0,
		BorderColor: borderColor, FillColorFnc: ApplyAlpha}
	r.CanvasObjectEmbed.Extend(r)
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
	if r.FillColor == color.Transparent {
		c.gc.SetFillColor(r.FillColorFnc(r.BorderColor))
	} else {
		c.gc.SetFillColor(r.FillColor)
	}
	c.gc.FillStroke()
}

// Auch gleichmaessige Polygone duerfen nicht fehlen.
type RegularPolygon struct {
	CanvasObjectEmbed
	Pos, Size              geom.Point
	Angle                  float64
	BorderWidth            float64
	BorderColor, FillColor color.LedColor
	FillColorFnc           ColorConvertFunc
	numPoints              int
}

func NewRegularPolygon(numPoints int, pos, size geom.Point, borderColor color.LedColor) *RegularPolygon {
	p := &RegularPolygon{Pos: pos, Size: size, Angle: 0.0, BorderWidth: 1.0,
		BorderColor: borderColor, FillColorFnc: ApplyAlpha, numPoints: numPoints}
	p.CanvasObjectEmbed.Extend(p)
	return p
}

func (p *RegularPolygon) Draw(c *Canvas) {
	c.gc.DrawRegularPolygon(p.numPoints, p.Pos.X, p.Pos.Y, p.Size.X/2.0, p.Angle)
	c.gc.SetStrokeWidth(p.BorderWidth)
	c.gc.SetStrokeColor(p.BorderColor)
	if p.FillColor == color.Transparent {
		c.gc.SetFillColor(p.FillColorFnc(p.BorderColor))
	} else {
		c.gc.SetFillColor(p.FillColor)
	}
	c.gc.FillStroke()
}

// Fuer Geraden ist dieser Datentyp vorgesehen, der von Pos1 nach Pos2
// verlaeuft.
type Line struct {
	CanvasObjectEmbed
	Pos1, Pos2 geom.Point
	Width      float64
	Color      color.LedColor
}

func NewLine(pos1, pos2 geom.Point, col color.LedColor) *Line {
	l := &Line{Pos1: pos1, Pos2: pos2, Width: 1.0, Color: col}
	l.CanvasObjectEmbed.Extend(l)
	return l
}

func (l *Line) Draw(c *Canvas) {
	c.gc.SetStrokeWidth(l.Width)
	c.gc.SetStrokeColor(l.Color)
	c.gc.DrawLine(l.Pos1.X, l.Pos1.Y, l.Pos2.X, l.Pos2.Y)
	c.gc.Stroke()
}

// Will man ein einzelnes Pixel exakt an einer LED-Position zeichnen, so
// eignet sich dieser Typ. Im Gegensatz zu den obigen Typen sind die
// Koordinaten eines Pixels ganze Zahlen und das Zeichnen erfolgt direkt
// in die draw.Image Struktur und nicht in gg.Context.
// Der Typ eignet sich nicht, wenn man beabsichtigt die Pixel wandern zu
// lassen! Da er immer auf ganze Koordinaten springt, sind seine Bewegungen
// ziemlich "hoelzern".
//
//   - BorderWidth: 0.0
//   - Size       : sqrt(2), sqrt(2)
//   - FillColor  : (manuell setzen)
type Pixel struct {
	CanvasObjectEmbed
	Pos   image.Point
	Color color.LedColor
}

func NewPixel(pos image.Point, col color.LedColor) *Pixel {
	p := &Pixel{Pos: pos, Color: col}
	p.CanvasObjectEmbed.Extend(p)
	return p
}

func (p *Pixel) Draw(c *Canvas) {
	bgColor := color.LedColorModel.Convert(c.img.At(p.Pos.X, p.Pos.Y)).(color.LedColor)
	c.img.Set(p.Pos.X, p.Pos.Y, p.Color.Mix(bgColor, color.Blend))
}

// Zur Darstellung von beliebigen Bildern (JPEG, PNG, etc) auf dem LED-Panel
// Da es nur wenige LEDs zur Darstellung hat, werden die Bilder gnadenlos
// skaliert und herunter gerechnet - manchmal bis der Arzt kommt... ;-)
type Image struct {
	CanvasObjectEmbed
	Pos, Size geom.Point
	Angle     float64
	Img       draw.Image
}

func NewImage(pos geom.Point, fileName string) *Image {
	i := &Image{Pos: pos, Angle: 0.0}
	i.CanvasObjectEmbed.Extend(i)
	i.Img = DecodeImageFile(fileName)
	i.Size = geom.NewPointIMG(i.Img.Bounds().Size())
	return i
}

func (i *Image) Read(fileName string) {
	i.Img = DecodeImageFile(fileName)
	if i.Size.X > 0 || i.Size.Y > 0 {
		return
	}
	i.Size = geom.NewPointIMG(i.Img.Bounds().Size())
}

func (i *Image) Draw(c *Canvas) {
	c.gc.Push()
	defer c.gc.Pop()
	if i.Angle != 0.0 {
		c.gc.RotateAbout(i.Angle, i.Pos.X, i.Pos.Y)
	}
	sx := i.Size.X / float64(i.Img.Bounds().Dx())
	sy := i.Size.Y / float64(i.Img.Bounds().Dy())
	c.gc.ScaleAbout(sx, sy, i.Pos.X, i.Pos.Y)
	c.gc.DrawImageAnchored(i.Img, i.Pos.X, i.Pos.Y, 0.5, 0.5)
	// sx := i.Size.X / float64(i.Img.Bounds().Dx())
	// sy := i.Size.Y / float64(i.Img.Bounds().Dy())
	// m := f64.Aff3{sx, 0.0, i.Pos.X - i.Size.X/2.0, 0.0, sy, i.Pos.Y - i.Size.Y/2.0}
	// draw.BiLinear.Transform(c.img, m, i.Img, i.Img.Bounds(), draw.Over, nil)
}

func DecodeImageFile(fileName string) draw.Image {
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
	return img.(draw.Image)
}

// Zur Darstellung von beliebigen Bildern (JPEG, PNG, etc) auf dem LED-Panel
// Da es nur wenige LEDs zur Darstellung hat, werden die Bilder gnadenlos
// skaliert und herunter gerechnet - manchmal bis der Arzt kommt... ;-)
type ImageList struct {
	CanvasObjectEmbed
	Pos, Size geom.Point
	Angle     float64
	ImgIdx    int
	imgs      []draw.Image
	imgBounds image.Rectangle
}

func NewImageList(pos geom.Point) *ImageList {
	i := &ImageList{Pos: pos, Angle: 0.0}
	i.CanvasObjectEmbed.Extend(i)
	i.imgs = make([]draw.Image, 0)
	i.ImgIdx = 0
	return i
}

func (i *ImageList) Add(img draw.Image) {
	i.imgs = append(i.imgs, img)
	i.imgBounds = img.Bounds()
	i.Size = geom.NewPointIMG(img.Bounds().Size())
}

func (i *ImageList) AddBlinkenLight(b *BlinkenFile) {
	i.imgs = i.imgs[:0]
	for idx := range b.NumFrames() {
		i.Add(b.Decode(idx))
	}
}

func (i *ImageList) Draw(c *Canvas) {
	sx := i.Size.X / float64(i.imgBounds.Dx())
	sy := i.Size.Y / float64(i.imgBounds.Dy())
	m := f64.Aff3{sx, 0.0, i.Pos.X - i.Size.X/2.0, 0.0, sy, i.Pos.Y - i.Size.Y/2.0}
	draw.BiLinear.Transform(c.img, m, i.imgs[i.ImgIdx], i.imgBounds, draw.Over, nil)
}

// Zur Darstellung von beliebigem Text.
var (
	defFont     = fonts.GoMedium
	defFontSize = 10.0
)

type Text struct {
	CanvasObjectEmbed
	Pos      geom.Point
	AX, AY   float64
	Angle    float64
	Color    color.LedColor
	Font     *fonts.Font
	FontSize float64
	Text     string
	fontFace font.Face
}

func NewText(pos geom.Point, text string, color color.LedColor) *Text {
	t := &Text{Pos: pos, Color: color, Font: defFont, FontSize: defFontSize,
		Text: text}
	t.CanvasObjectEmbed.Extend(t)
	t.AX, t.AY = 0.5, 0.5
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
	c.gc.DrawStringAnchored(t.Text, t.Pos.X, t.Pos.Y, t.AX, t.AY)
}

// Fuer das direkte Zeichnen von Text auf dem LED-Grid, existieren einige
// 'fixed size' Bitmap-Schriften, die ohne Rastern und Rendern sehr schnell
// dargestellt werden koennen.
var (
	defFixedFontFace = Face3x5
)

type FixedText struct {
	CanvasObjectEmbed
	Pos    fixed.Point26_6
	Color  color.LedColor
	text   string
	drawer *font.Drawer
	rect   fixed.Rectangle26_6
	dp     fixed.Point26_6
}

func NewFixedText(pos fixed.Point26_6, col color.LedColor, text string) *FixedText {
	t := &FixedText{Pos: pos, Color: col}
	t.CanvasObjectEmbed.Extend(t)
	t.drawer = &font.Drawer{
		Face: defFixedFontFace,
	}
	t.SetText(text)
	return t
}

func (t *FixedText) Text() string {
	return t.text
}

func (t *FixedText) SetText(text string) {
	t.text = text
	t.rect, _ = t.drawer.BoundString(text)
	t.dp = t.rect.Min.Add(t.rect.Max).Div(fixed.I(2))
}

func (t *FixedText) Draw(c *Canvas) {
	t.drawer.Dst = c.img
	t.drawer.Src = image.NewUniform(t.Color)
	t.drawer.Dot = t.Pos.Sub(t.dp)
	t.drawer.DrawString(t.text)
}
