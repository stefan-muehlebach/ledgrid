package ledgrid

import (
	"container/list"
	"image"
	gocolor "image/color"
	"log"
	"math"
	"os"
	"sync"
	"time"

	"golang.org/x/image/math/f64"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/fonts"
	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid/color"
	"golang.org/x/image/draw"
)

// Ein Canvas ist eine animierbare Zeichenflaeche. Ihr koennen eine beliebige
// Anzahl von zeichenbaren Objekten (Interface CanvasObject) hinzugefuegt
// werden.
type Canvas struct {
	ObjList            *list.List
	BackColor          color.LedColor
	Rect               image.Rectangle
	Img                draw.Image
	GC                 *gg.Context
	objMutex           *sync.RWMutex
	paintWatch         *Stopwatch
	syncAnim, syncSend chan bool
}

func NewCanvas(size image.Point) *Canvas {
	c := &Canvas{}
	c.ObjList = list.New()
	c.BackColor = color.Black
	c.Rect = image.Rectangle{Max: size}
	c.Img = image.NewRGBA(c.Rect)
	c.GC = gg.NewContextForRGBA(c.Img.(*image.RGBA))
	c.objMutex = &sync.RWMutex{}
	c.paintWatch = NewStopwatch()
	return c
}

func (c *Canvas) Close() {
	c.Purge()
}

func (c *Canvas) ColorModel() gocolor.Model {
	return gocolor.RGBAModel
}

func (c *Canvas) Bounds() image.Rectangle {
	return c.Rect
}

func (c *Canvas) At(x, y int) gocolor.Color {
	return c.Img.At(x, y)
}

func (c *Canvas) Set(x, y int, col gocolor.Color) {
	c.Img.Set(x, y, col)
}

// Fuegt der Zeichenflaeche weitere Objekte hinzu. Der Zufgriff auf den
// entsprechenden Slice wird nicht synchronisiert.
func (c *Canvas) Add(objs ...CanvasObject) {
	c.objMutex.Lock()
	for _, obj := range objs {
		c.ObjList.PushBack(obj)
	}
	c.objMutex.Unlock()
}

// Loescht alle Objekte von der Zeichenflaeche.
func (c *Canvas) Purge() {
	c.objMutex.Lock()
	c.ObjList.Init()
	c.objMutex.Unlock()
}

func (c *Canvas) Refresh() {
	c.GC.SetFillColor(c.BackColor)
	c.GC.Clear()
	c.objMutex.RLock()
	for ele := c.ObjList.Front(); ele != nil; ele = ele.Next() {
		obj := ele.Value.(CanvasObject)
        if obj.IsDead() {
            c.ObjList.Remove(ele)
            continue
        }
		if !obj.IsVisible() {
			continue
		}
		obj.Draw(c)
	}
	c.objMutex.RUnlock()
}

func (c *Canvas) StartRefresh(syncAnim, syncSend chan bool) {
	c.syncAnim = syncAnim
	c.syncSend = syncSend
	go c.refreshThread()
}

func (c *Canvas) refreshThread() {
	for {
		<-c.syncSend
		<-c.syncAnim
		c.paintWatch.Start()
		c.Refresh()
		c.syncAnim <- true
		draw.Draw(AnimCtrl.ledGrid, AnimCtrl.ledGrid.Bounds(), c.Img,
			image.Point{}, draw.Over)
		c.paintWatch.Stop()
		c.syncSend <- true
	}
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
	IsVisible() bool
    Kill()
    IsDead() bool
	Draw(c *Canvas)
}

// Wie bei den Animationen gibt es für die darstellbaren Objekte (CanvasObject)
// ein entsprechendes Embedable, welche die für die meisten Objekte
// brauchbaren Methoden enthält.
type CanvasObjectEmbed struct {
	wrapper CanvasObject
	isVisible bool
    isDead bool
}

func (c *CanvasObjectEmbed) Extend(wrapper CanvasObject) {
	c.wrapper = wrapper
	c.isVisible = true
    c.isDead = false
}

func (c *CanvasObjectEmbed) Show() {
	if !c.isVisible {
		c.isVisible = true
	}
}

func (c *CanvasObjectEmbed) Hide() {
	if c.isVisible {
		c.isVisible = false
	}
}

func (c *CanvasObjectEmbed) IsVisible() bool {
	return c.isVisible
}

func (c *CanvasObjectEmbed) Kill() {
    c.isDead = true
}

func (c *CanvasObjectEmbed) IsDead() bool {
	return c.isDead
}

// Und hier etwas mit Farben...
var (
	// Ueber die NewXXX-Funktionen koennen die Objekte einfacher erzeugt
	// werden. Die Fuellfarbe ist gleich der Randfarbe, hat aber einen
	// niedrigeren Alpha-Wert, der mit dieser Konstante definiert werden
	// kann.
	fillAlpha = 0.4
)

type ColorFunc func(color.LedColor) color.LedColor

func ApplyAlpha(c color.LedColor) color.LedColor {
	alpha := float64(c.A) / 255.0
	return c.Alpha(alpha * fillAlpha)
}

var (
	colorFncMap = map[string]ColorFunc{
		"ApplyAlpha": ApplyAlpha,
	}
)

// Mit Ellipse sind alle kreisartigen Objekte abgedeckt. Pos bezeichnet die
// Position des Mittelpunktes und mit Size ist die Breite und Hoehe des
// gesamten Objektes gemeint. Falls ein Rand gezeichnet werden soll, muss
// BorderWith einen Wert >0 enthalten und FillColor, resp. BorderColor
// enthalten die Farben fuer Rand und Flaeche.
type Ellipse struct {
	CanvasObjectEmbed
	Pos, Size              geom.Point
	Angle                  float64
	StrokeWidth            float64
	StrokeColor, FillColor color.LedColor
	FillColorFnc           string
}

// Erzeugt eine 'klassische' Ellipse mit einer Randbreite von einem Pixel und
// setzt die Fuellfarbe gleich Randfarbe mit Alpha-Wert von 0.3.
// Will man die einzelnen Werte flexibler verwenden, empfiehlt sich die
// Erzeugung mittels &Ellipse{...}.
func NewEllipse(pos, size geom.Point, borderColor color.LedColor) *Ellipse {
	e := &Ellipse{Pos: pos, Size: size, StrokeWidth: 1.0,
		StrokeColor: borderColor, FillColorFnc: "ApplyAlpha"}
	e.CanvasObjectEmbed.Extend(e)
	return e
}

func (e *Ellipse) Draw(c *Canvas) {
	if e.Angle != 0.0 {
		c.GC.Push()
		c.GC.RotateAbout(e.Angle, e.Pos.X, e.Pos.Y)
		defer c.GC.Pop()
	}
	c.GC.DrawEllipse(e.Pos.X, e.Pos.Y, e.Size.X/2, e.Size.Y/2)
	c.GC.SetStrokeWidth(e.StrokeWidth)
	c.GC.SetStrokeColor(e.StrokeColor)
	if e.FillColor == color.Transparent && e.FillColorFnc != "" {
		c.GC.SetFillColor(colorFncMap[e.FillColorFnc](e.StrokeColor))
	} else {
		c.GC.SetFillColor(e.FillColor)
	}
	c.GC.FillStroke()
}

// Rectangle ist fuer alle rechteckigen Objekte vorgesehen. Pos bezeichnet
// den Mittelpunkt des Objektes und Size die Breite, rsep. Hoehe.
type Rectangle struct {
	CanvasObjectEmbed
	Pos, Size              geom.Point
	Angle                  float64
	StrokeWidth            float64
	StrokeColor, FillColor color.LedColor
	FillColorFnc           string
}

func NewRectangle(pos, size geom.Point, borderColor color.LedColor) *Rectangle {
	r := &Rectangle{Pos: pos, Size: size, StrokeWidth: 1.0,
		StrokeColor: borderColor, FillColorFnc: "ApplyAlpha"}
	r.CanvasObjectEmbed.Extend(r)
	return r
}

func (r *Rectangle) Draw(c *Canvas) {
	if r.Angle != 0.0 {
		c.GC.Push()
		c.GC.RotateAbout(r.Angle, r.Pos.X, r.Pos.Y)
		defer c.GC.Pop()
	}
	c.GC.DrawRectangle(r.Pos.X-r.Size.X/2, r.Pos.Y-r.Size.Y/2, r.Size.X, r.Size.Y)
	c.GC.SetStrokeWidth(r.StrokeWidth)
	c.GC.SetStrokeColor(r.StrokeColor)
	if r.FillColor == color.Transparent && r.FillColorFnc != "" {
		c.GC.SetFillColor(colorFncMap[r.FillColorFnc](r.StrokeColor))
	} else {
		c.GC.SetFillColor(r.FillColor)
	}
	c.GC.FillStroke()
}

// Auch gleichmaessige Polygone duerfen nicht fehlen.
type RegularPolygon struct {
	CanvasObjectEmbed
	Pos, Size              geom.Point
	Angle                  float64
	StrokeWidth            float64
	StrokeColor, FillColor color.LedColor
	FillColorFnc           string
	numPoints              int
}

func NewRegularPolygon(numPoints int, pos, size geom.Point, borderColor color.LedColor) *RegularPolygon {
	p := &RegularPolygon{Pos: pos, Size: size, Angle: 0.0, StrokeWidth: 1.0,
		StrokeColor: borderColor, FillColorFnc: "ApplyAlpha", numPoints: numPoints}
	p.CanvasObjectEmbed.Extend(p)
	return p
}

func (p *RegularPolygon) Draw(c *Canvas) {
	c.GC.DrawRegularPolygon(p.numPoints, p.Pos.X, p.Pos.Y, p.Size.X/2.0, p.Angle)
	c.GC.SetStrokeWidth(p.StrokeWidth)
	c.GC.SetStrokeColor(p.StrokeColor)
	if p.FillColor == color.Transparent && p.FillColorFnc != "" {
		c.GC.SetFillColor(colorFncMap[p.FillColorFnc](p.StrokeColor))
	} else {
		c.GC.SetFillColor(p.FillColor)
	}
	c.GC.FillStroke()
}

// Fuer Geraden ist dieser Datentyp vorgesehen, der von Pos1 nach Pos2
// verlaeuft.
type Line struct {
	CanvasObjectEmbed
	Pos1, Pos2  geom.Point
	StrokeWidth float64
	StrokeColor color.LedColor
}

func NewLine(pos1, pos2 geom.Point, col color.LedColor) *Line {
	l := &Line{Pos1: pos1, Pos2: pos2, StrokeWidth: 1.0, StrokeColor: col}
	l.CanvasObjectEmbed.Extend(l)
	return l
}

func (l *Line) Draw(c *Canvas) {
	c.GC.SetStrokeWidth(l.StrokeWidth)
	c.GC.SetStrokeColor(l.StrokeColor)
	c.GC.DrawLine(l.Pos1.X, l.Pos1.Y, l.Pos2.X, l.Pos2.Y)
	c.GC.Stroke()
}

// Will man ein einzelnes Pixel exakt an einer LED-Position zeichnen, so
// eignet sich dieser Typ. Im Gegensatz zu den obigen Typen sind die
// Koordinaten eines Pixels ganze Zahlen und das Zeichnen erfolgt direkt
// in die draw.Image Struktur und nicht in gg.Context. Es ist zu beachten,
// dass bei diesem Typ die Koordinaten von pos als Spalten-, resp. Zeilenindex
// des Led-Grids interpretiert werden!
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
	bgColor := color.LedColorModel.Convert(c.Img.At(p.Pos.X, p.Pos.Y)).(color.LedColor)
	c.Img.Set(p.Pos.X, p.Pos.Y, p.Color.Mix(bgColor, color.Blend))
}

// Ein einzelnes Pixel, dessen Bewegungen weicher (smooth) animiert werden
// koennen, ist der Typ Dot. Da er grosse Aehnlichkeit zum Typ Pixel aufweist,
// werden auch hier die Koordinaten als Spalten-, resp. Zeilenindex
// interpretiert.
type Dot struct {
	CanvasObjectEmbed
	Pos   geom.Point
	Color color.LedColor
}

func NewDot(pos geom.Point, col color.LedColor) *Dot {
	d := &Dot{Pos: pos, Color: col}
	d.CanvasObjectEmbed.Extend(d)
	return d
}

func (d *Dot) Draw(c *Canvas) {
	// c.GC.DrawEllipse(d.Pos.X+0.5, d.Pos.Y+0.5, 0.5, 0.5)
	c.GC.DrawEllipse(d.Pos.X+0.5, d.Pos.Y+0.5, math.Sqrt2/2.0, math.Sqrt2/2.0)
	c.GC.SetFillColor(d.Color)
	c.GC.Fill()
}

// Text laesst sich auf zwei Arten darstellen: mittels TrueType fonts and
// beliebig skaliert oder mit FixedSize Pixel-Schriften
var (
	defFont     = fonts.GoMedium
	defFontSize = 10.0
)

// Die folgenden Typen, Records (embeddable), Konstanten, etc. sind fuer
// das Ausrichten von Textbloecken in vertikaler und horizontaler Richtung
// vorgesehen.
type Align int

const (
	AlignLeft Align = 1 << iota
	AlignCenter
	AlignRight
	AlignBottom
	AlignMiddle
	AlignTop

	alignHMask = 0b000111
	alignVMask = 0b111000
)

type alignEmbed struct {
	ax, ay float64
}

func (a *alignEmbed) SetAlign(align Align) {
	hAlign := align & alignHMask
	vAlign := align & alignVMask
	switch hAlign {
	case AlignLeft:
		a.ax = 0.0
	case AlignCenter:
		a.ax = 0.5
	case AlignRight:
		a.ax = 1.0
	}
	switch vAlign {
	case AlignBottom:
		a.ay = 0.0
	case AlignMiddle:
		a.ay = 0.5
	case AlignTop:
		a.ay = 1.0
	}
}

// Zur Darstellung von Text mit TrueType-Schriften
type Text struct {
	CanvasObjectEmbed
	alignEmbed
	Pos   geom.Point
	Angle float64
	// AX, AY   float64
	Color    color.LedColor
	Text     string
	font     *fonts.Font
	fontSize float64
	fontFace font.Face
}

func NewText(pos geom.Point, text string, col color.LedColor) *Text {
	t := &Text{Pos: pos, Color: col, Text: text}
	t.CanvasObjectEmbed.Extend(t)
	t.SetFont(defFont, defFontSize)
	return t
}

func (t *Text) SetFont(font *fonts.Font, size float64) {
	t.font = font
	t.fontSize = size
	t.fontFace = fonts.NewFace(t.font, t.fontSize)
}

func (t *Text) Draw(c *Canvas) {
	if t.Angle != 0.0 {
		c.GC.Push()
		c.GC.RotateAbout(t.Angle, t.Pos.X, t.Pos.Y)
		defer c.GC.Pop()
	}
	c.GC.SetTextColor(t.Color)
	c.GC.SetFontFace(t.fontFace)
	c.GC.DrawStringAnchored(t.Text, t.Pos.X, t.Pos.Y, t.ax, t.ay)
}

// Fuer das direkte Zeichnen von Text auf dem LED-Grid, existieren einige
// 'fixed size' Bitmap-Schriften, die ohne Rastern und Rendern sehr schnell
// dargestellt werden koennen.
var (
	defFixedFontFace = Pico3x5
)

type FixedText struct {
	CanvasObjectEmbed
	alignEmbed
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
	t.drawer = &font.Drawer{}
	t.drawer.Face = defFixedFontFace
	t.text = text
	t.updateSize()
	return t
}

func (t *FixedText) SetFont(font font.Face) {
	t.drawer.Face = font
	t.updateSize()
}

func (t *FixedText) SetAlign(align Align) {
	t.alignEmbed.SetAlign(align)
	t.updateSize()
}

func (t *FixedText) Text() string {
	return t.text
}

func (t *FixedText) SetText(text string) {
	t.text = text
	t.updateSize()
}

func (t *FixedText) updateSize() {
	rect, _ := font.BoundString(t.drawer.Face, t.text)
	t.dp.X = (rect.Max.X - rect.Min.X).Mul(fixed.Int26_6(64 * t.ax))
	t.dp.Y = (rect.Min.Y - rect.Max.Y).Mul(fixed.Int26_6(64 * t.ay))
}

func (t *FixedText) Draw(c *Canvas) {
	t.drawer.Dst = c.Img
	t.drawer.Src = image.NewUniform(t.Color)
	t.drawer.Dot = t.Pos.Sub(t.dp)
	t.drawer.DrawString(t.text)
}

// Zur Darstellung von beliebigen Bildern (JPEG, PNG, etc). Wie Pos genau
// interpretiert wird, ist vom Alignment (wie beim Text) abhaengig. Size
// ist die Zielgroesse des Bildes auf dem LedGrid und muss nicht mit
// Img.Bounds() uebereinstimmen.
type Image struct {
	CanvasObjectEmbed
	alignEmbed
	Pos, Size geom.Point
	Angle     float64
	Img       draw.Image
}

func NewImage(pos geom.Point, fileName string) *Image {
	i := &Image{Pos: pos, Angle: 0.0}
	i.CanvasObjectEmbed.Extend(i)
	i.Img = LoadImage(fileName)
	i.Size = geom.NewPointIMG(i.Img.Bounds().Size())
	i.ax, i.ay = 0.5, 0.5
	return i
}

func (i *Image) Read(fileName string) {
	i.Img = LoadImage(fileName)
	if i.Size.X > 0 || i.Size.Y > 0 {
		return
	}
	i.Size = geom.NewPointIMG(i.Img.Bounds().Size())
}

func (i *Image) Draw(c *Canvas) {
	sx := i.Size.X / float64(i.Img.Bounds().Dx())
	sy := i.Size.Y / float64(i.Img.Bounds().Dy())
	cos := math.Cos(i.Angle)
	sin := math.Sin(i.Angle)
	m := f64.Aff3{cos * sx, -sin * sy, -cos*i.ax*i.Size.X + sin*(1-i.ay)*i.Size.Y + i.Pos.X,
		sin * sx, cos * sy, -sin*i.ax*i.Size.X - cos*(1-i.ay)*i.Size.Y + i.Pos.Y}
	draw.BiLinear.Transform(c.Img, m, i.Img, i.Img.Bounds(), draw.Over, nil)
}

func LoadImage(fileName string) draw.Image {
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
	Image
    NormAnimationEmbed
	imgList []draw.Image
	durList []time.Duration
}

func NewImageList(pos geom.Point) *ImageList {
	i := &ImageList{}
	i.CanvasObjectEmbed.Extend(i)
	i.NormAnimationEmbed.Extend(i)
    i.Pos = pos
	i.Curve = AnimationLinear
    i.ax, i.ay = 0.5, 0.5
	i.imgList = make([]draw.Image, 0)
    i.durList = make([]time.Duration, 0)
	return i
}

func (i *ImageList) Add(img draw.Image, dur time.Duration) {
	i.imgList = append(i.imgList, img)
    if len(i.imgList) == 1 {
        i.Img = i.imgList[0]
    }
	i.Size = geom.NewPointIMG(img.Bounds().Size())
	i.duration += dur
	i.durList = append(i.durList, i.duration)
}

func (i *ImageList) AddBlinkenLight(b *BlinkenFile) {
	i.imgList = i.imgList[:0]
	i.durList = i.durList[:0]
	for idx := range b.NumFrames() {
		i.Add(b.Decode(idx), b.Duration(idx))
	}
}

func (i *ImageList) Init() {
	i.Img = i.imgList[0]
}

func (i *ImageList) Tick(t float64) {
	var idx int

	ts := time.Duration(t * float64(i.duration))
	for idx = 0; idx < len(i.durList); idx++ {
		if i.durList[idx] >= ts {
			break
		}
	}
	i.Img = i.imgList[idx]
}


/*
// Zur Darstellung von beliebigen Bildern (JPEG, PNG, etc) auf dem LED-Panel
// Da es nur wenige LEDs zur Darstellung hat, werden die Bilder gnadenlos
// skaliert und herunter gerechnet - manchmal bis der Arzt kommt... ;-)
type ImageList struct {
	CanvasObjectEmbed
	NormAnimationEmbed
	Pos, Size geom.Point
	Angle     float64
	ImgIdx    int
	Imgs      []draw.Image
	Durs      []time.Duration
	imgBounds image.Rectangle
}

func NewImageList(pos geom.Point) *ImageList {
	i := &ImageList{Pos: pos, Angle: 0.0, ImgIdx: 0}
	i.CanvasObjectEmbed.Extend(i)
	i.NormAnimationEmbed.Extend(i)
	i.Curve = AnimationLinear
	i.Imgs = make([]draw.Image, 0)
	return i
}

func (i *ImageList) Add(img draw.Image, dur time.Duration) {
	i.Imgs = append(i.Imgs, img)
	i.imgBounds = img.Bounds()
	i.Size = geom.NewPointIMG(img.Bounds().Size())
	i.duration += dur
	i.Durs = append(i.Durs, i.duration)
}

func (i *ImageList) AddBlinkenLight(b *BlinkenFile) {
	i.Imgs = i.Imgs[:0]
	i.Durs = i.Durs[:0]
	for idx := range b.NumFrames() {
		i.Add(b.Decode(idx), b.Duration(idx))
	}
	// i.Size = geom.Point{float64(b.Width), float64(b.Height)}
}

func (i *ImageList) Draw(c *Canvas) {
	sx := i.Size.X / float64(i.imgBounds.Dx())
	sy := i.Size.Y / float64(i.imgBounds.Dy())
	m := f64.Aff3{sx, 0.0, i.Pos.X - i.Size.X/2.0, 0.0, sy, i.Pos.Y - i.Size.Y/2.0}
	draw.BiLinear.Transform(c.Img, m, i.Imgs[i.ImgIdx], i.imgBounds, draw.Over, nil)
}

func (i *ImageList) Init() {
	i.ImgIdx = 0
}

func (i *ImageList) Tick(t float64) {
	var idx int

	ts := time.Duration(t * float64(i.duration))
	for idx = 0; idx < len(i.Durs); idx++ {
		if i.Durs[idx] >= ts {
			break
		}
	}
	i.ImgIdx = idx
}

*/
