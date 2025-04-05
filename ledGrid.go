//go:build !tinygo

package ledgrid

import (
	"container/list"
	"image"
	"image/color"
	"log"
	"sync"

	"github.com/stefan-muehlebach/ledgrid/colors"
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
	// Es koennen eine ganze Reihe von Canvas'es verwendet werden - so um bspw.
	// mit mehreren Layern oder Ueberblendungen zu arbeiten. Die Canvas'es
	// werden in einer dynamischen Liste verwaltet. Die Darstellung beginnt
	// mit dem hintersten Canvas und stellt zuletzt (d.h. zuvorderst) das
	// Canvas am Anfang der Liste dar.
	CanvasList *list.List
	canvMutex  *sync.RWMutex

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
	g.CanvasList = list.New()
	g.canvMutex = &sync.RWMutex{}

	g.NewCanvas()

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
	return colors.LedColorModel
}

func (g *LedGrid) Bounds() image.Rectangle {
	return g.Rect
}

func (g *LedGrid) At(x, y int) color.Color {
	return g.LedColorAt(x, y)
}

func (g *LedGrid) Set(x, y int, c color.Color) {
	c1 := colors.LedColorModel.Convert(c).(colors.LedColor)
	g.SetLedColor(x, y, c1)
}

// Dient dem schnelleren Zugriff auf den Farbwert einer bestimmten Zelle, resp.
// einer bestimmten LED. Analog zu At(), retourniert den Farbwert jedoch als
// LedColor-Typ.
func (g *LedGrid) LedColorAt(x, y int) colors.LedColor {
	if !(image.Point{x, y}.In(g.Rect)) {
		return colors.LedColor{}
	}
	idx := g.PixOffset(x, y)
	if idx < 0 {
		return colors.Black
	}
	src := g.Pix[idx : idx+3 : idx+3]
	return colors.LedColor{src[0], src[1], src[2], 0xff}
}

// Analoge Methode zu Set(), jedoch ohne zeitaufwaendige Konvertierung.
func (g *LedGrid) SetLedColor(x, y int, c colors.LedColor) {
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
func (g *LedGrid) Clear(c colors.LedColor) {
	draw.Draw(g, g.Rect, image.NewUniform(c), image.Point{}, draw.Src)
}

func (g *LedGrid) Reset() {
	g.AnimCtrl.Purge()
	for elem := g.CanvasList.Front(); elem != nil; elem = elem.Next() {
		canv, ok := elem.Value.(*Canvas)
		if !ok {
			continue
		}
		canv.Reset()
	}
}

// Zeigt den aktuellen Inhalt des Grid auf der beim Erstellen spezifizierten
// Hardware dar.
func (g *LedGrid) Show() {
	g.Client.Send(g.Pix)
}

// Erzeugt ein neues Canvas-Objekt und hanegt es an den Schluss der Liste.
// Beim Zeichnen geht LedGrid von hinten nach vorne
func (g *LedGrid) NewCanvas() (*Canvas, int) {
	canv := NewCanvas(g.Rect.Size())
	g.canvMutex.Lock()
	g.CanvasList.PushBack(canv)
	layer := g.CanvasList.Len() - 1
	g.canvMutex.Unlock()
	return canv, layer
}

func (g *LedGrid) DelCanvas(canv *Canvas) {
	for elem := g.CanvasList.Front().Next(); elem != nil; elem = elem.Next() {
		obj := elem.Value.(*Canvas)
		if obj == canv {
			g.CanvasList.Remove(elem)
			return
		}
	}
}

// Liefert das Canvas-Objekt zurueck, welches fuer den Layer layer definiert
// ist. Per Default ist nur Layer 0 vorhanden, die Methode bricht ab, wenn
// es kein Canvas-Objekt zum gewuenschten Layer gibt.
func (g *LedGrid) Canvas(layer int) *Canvas {
	var elem *list.Element
	var id int

	if layer >= g.CanvasList.Len() {
		log.Fatalf("Cannot access canvas[%d]; only %d in list", layer, g.CanvasList.Len())
		return nil
	}
	for elem, id = g.CanvasList.Front(), 0; elem != nil; elem, id = elem.Next(), id+1 {
		if id == layer {
			break
		}
	}
	return elem.Value.(*Canvas)
}

func (g *LedGrid) StartRefresh() {
	go g.refreshThread()
}

func (g *LedGrid) refreshThread() {
	var canv *Canvas
	var ok bool

	for {
		<-g.syncChan
		g.canvMutex.RLock()
		g.Clear(colors.Black)
		for ele := g.CanvasList.Back(); ele != nil; ele = ele.Prev() {
			if canv, ok = ele.Value.(*Canvas); !ok {
				log.Fatalf("Wrong data in canvas-list")
			}
			canv.Refresh()
			draw.DrawMask(g, g.Bounds(), canv.Img, image.Point{},
				canv.Mask, image.Point{}, draw.Over)
		}
		g.canvMutex.RUnlock()
		g.syncChan <- true
		g.Show()
	}
}
