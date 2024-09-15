package ledgrid

import (
	"image"
	"image/color"

	ledcolor "github.com/stefan-muehlebach/ledgrid/color"
	"github.com/stefan-muehlebach/ledgrid/conf"
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
	Pix []uint8
    Client GridClient
	// Mit dieser Struktur (slice of slices) werden Pixel-Koordinaten in
	// Indizes uebersetzt.
	idxMap conf.IndexMap
    syncChan chan bool
}

// Erstellt ein neues LED-Panel. size enthaelt die Dimension des (gesamten)
// Panels. Wird bei modConf nil uebergeben, so wird eine Default-Konfiguration
// der Module angenommen, welche bei DefaultModuleConfig naeher beschrieben
// wird.
// func NewLedGrid(size image.Point, modConf conf.ModuleConfig) *LedGrid {
// 	g := &LedGrid{}
// 	g.Rect = image.Rectangle{Max: size}
// 	g.Pix = make([]uint8, 3*g.Rect.Dx()*g.Rect.Dy())

// 	// Autom. Formatwahl
// 	if modConf == nil {
// 		modConf = conf.DefaultModuleConfig(g.Rect.Size())
// 	}

// 	g.idxMap = modConf.IndexMap()
//     g.syncChan = make(chan bool)
//     go g.refreshThread()
// 	return g
// }

// Dies ist die neue Art, ein LedGrid-Objekt zu erstellen. Der GridClient
// (d.h. der Client-seitige Teil fuer die Ansteuerung) ist dabei Teil von
// LedGrid. Mit Angaben von host und port kann die Groesse des Panels (oder
// Emulations-Fensters) selbst. ermittelt werden.
func NewLedGridBySize(host string, port uint, size image.Point) *LedGrid {
    modConf := conf.DefaultModuleConfig(size)
    return NewLedGrid(host, port, modConf)
}

func NewLedGrid(host string, port uint, modConf conf.ModuleConfig) *LedGrid {
    g := &LedGrid{}
    g.Client = NewNetGridClient(host, port)
    g.Rect = image.Rectangle{Max: modConf.Size()}
    g.Pix = make([]uint8, 3*len(modConf)*conf.ModuleDim.X*conf.ModuleDim.Y)
	g.idxMap = modConf.IndexMap()
    g.syncChan = make(chan bool)
    go g.refreshThread()
	return g
}

func (g *LedGrid) Close() {
    g.Client.Close()
}

// Die folgenden Methoden implementieren das image.Image Interface (resp.
// draw.Image).
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
	return 3*g.idxMap[x][y]
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

func (g *LedGrid) refreshThread() {
    g.syncChan <- true
    for {
        <- g.syncChan
        g.Client.Send(g.Pix)
        g.syncChan <- true
    }
}

