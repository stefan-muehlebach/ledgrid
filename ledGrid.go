package ledgrid

import (
	"image"
	"image/color"

	ledcolor "github.com/stefan-muehlebach/ledgrid/color"
	"github.com/stefan-muehlebach/ledgrid/conf"
	"golang.org/x/image/draw"
)

// Entspricht dem Bild, welches auf einem LED-Panel angezeigt werden kann.
// Implementiert die Interfaces image.Image und draw.Image, also die Methoden
// ColorModel, Bounds, At und Set.
type LedGrid struct {
	// Groesse des LedGrids. Falls dieses LedGrid Teil eines groesseren
	// Panels sein sollte, dann muss Rect.Min nicht unbedingt {0, 0} sein.
	Rect image.Rectangle
	// Enthaelt die Farbwerte red, green, blue (RGB) fuer jede LED, welche
	// das LedGrid ausmachen. Die Reihenfolge entspricht dabei der
	// Verkabelung!
	Pix    []uint8
	Client GridClient
	// Neu ist das Objekt LedGrid die zentrale steuernde Instanz. Das bedeutet
	// dass sowohl AnimationController als auch Canvas(es) hier vermerkt
	// sein muessen.
	AnimCtrl *AnimationController
	Canvas   *Canvas

	// Mit dieser Struktur (slice of slices) werden Pixel-Koordinaten in
	// Indizes uebersetzt.
	idxMap   conf.IndexMap
	syncChan chan bool
}

// Erstellt ein neues LedGrid-Objekt, welches die Groesse size in Anzahl LEDs
// horizontal, resp. vertikal hat. Die Verkabelung wird vollflaechig und
// gem. Methode DefaultModuleConfig vorgenommen.
func NewLedGridBySize(client GridClient, size image.Point) *LedGrid {
	modConf := conf.DefaultModuleConfig(size)
	return NewLedGrid(client, modConf)
}

// Erstellt ein neues LedGrid-Objekt, welches als Verkabelung modConf hat.
func NewLedGrid(client GridClient, modConf conf.ModuleConfig) *LedGrid {
	g := &LedGrid{}
	g.Client = client
	g.Rect = image.Rectangle{Max: modConf.Size()}
	g.Pix = make([]uint8, 3*len(modConf)*conf.ModuleDim.X*conf.ModuleDim.Y)
	g.idxMap = modConf.IndexMap()
	g.syncChan = make(chan bool)
	g.AnimCtrl = NewAnimationController(g.syncChan)
	g.Canvas = NewCanvas(g.Rect.Size())
	return g
}

func (g *LedGrid) Close() {
	g.Client.Close()
}

// The following methods implement the draw.Image interface, LedGrid can
// therefore be used as the destination for all kind of drawings - especially
// for a call to draw.Draw() in order to compose the data from Canvas
// objects before sending the picture to a GridClient.
func (g *LedGrid) ColorModel() color.Model {
	return ledcolor.LedColorModel
}

func (g *LedGrid) Bounds() image.Rectangle {
	return g.Rect
}

func (g *LedGrid) At(x, y int) color.Color {
	return g.LedColorAt(x, y)
}

func (g *LedGrid) Set(x, y int, c color.Color) {
	c1 := ledcolor.LedColorModel.Convert(c).(ledcolor.LedColor)
	g.SetLedColor(x, y, c1)
}

// Dient dem schnelleren Zugriff auf den Farbwert einer bestimmten Zelle, resp.
// einer bestimmten LED. Analog zu At(), retourniert den Farbwert jedoch als
// LedColor-Typ.
func (g *LedGrid) LedColorAt(x, y int) ledcolor.LedColor {
	if !(image.Point{x, y}.In(g.Rect)) {
		return ledcolor.LedColor{}
	}
	idx := g.PixOffset(x, y)
	if idx < 0 {
		return ledcolor.Black
	}
	src := g.Pix[idx : idx+3 : idx+3]
	return ledcolor.NewLedColorRGB(src[0], src[1], src[2])
}

// Analoge Methode zu Set(), jedoch ohne zeitaufwaendige Konvertierung.
func (g *LedGrid) SetLedColor(x, y int, c ledcolor.LedColor) {
	if !(image.Point{x, y}.In(g.Rect)) {
		return
	}
	idx := g.PixOffset(x, y)
	if idx < 0 {
		return
	}
	dst := g.Pix[idx : idx+3 : idx+3]
	dst[0] = c.R
	dst[1] = c.G
	dst[2] = c.B
}

// Damit wird der Offset eines bestimmten Farbwerts innerhalb des Slices
// Pix berechnet. Dabei wird beruecksichtigt, dass das die LED's im LedGrid
// schlangenfoermig angeordnet sind, und der Beginn der LED-Kette frei
// waehlbar in einer Ecke des Panels liegen kann.
func (g *LedGrid) PixOffset(x, y int) int {
	return 3 * g.idxMap[x][y]
}

// Mit Clear kann das ganze Grid geloescht, resp. alle LEDs auf die gleiche
// Farbe gebracht werden. Das Anzeigen, resp. der Refresh des Panel ist
// Teil dieser Methode.
func (g *LedGrid) Clear(c ledcolor.LedColor) {
	for idx := 0; idx < len(g.Pix); idx += 3 {
		dst := g.Pix[idx : idx+3 : idx+3]
		dst[0] = c.R
		dst[1] = c.G
		dst[2] = c.B
	}
	g.Show()
}

// Zeigt den aktuellen Inhalt des Grid auf der beim Erstellen spezifizierten
// Hardware dar.
func (g *LedGrid) Show() {
	g.Client.Send(g.Pix)
}

func (g *LedGrid) StartRefresh() {
	go g.refreshThread()
}

func (g *LedGrid) refreshThread() {
	for {
		<-g.syncChan
		g.Canvas.Refresh()
		draw.DrawMask(g, g.Bounds(), g.Canvas.Img, image.Point{},
			g.Canvas.Mask, image.Point{}, draw.Over)
		g.syncChan <- true
		g.Show()
	}
}
