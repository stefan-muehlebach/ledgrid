package ledgrid

import (
	"image"
	"time"
)

// Diese Datei enthaelt viele Interfaces aber auch sog. Embedables, welche
// als Standard- oder Default-Implementation des entsprechenden Interfaces
// genutzt werden koennen - aber nicht muessen.

// Verschiedene Objekte sollten Namen haben koennen, die man beispielsweise in
// GUIs oder TUIs anzeigen kann. Dieses Interface implementieren also alle
// benennbaren Objekte.
type Nameable interface {
	Name() string
	SetName(name string)
}

// Das entsprechende Embedable.
type NameableEmbed struct {
	name string
}

func (n *NameableEmbed) Init(name string) {
	n.SetName(name)
}

func (n *NameableEmbed) Name() string {
	return n.name
}

func (n *NameableEmbed) SetName(name string) {
	n.name = name
}

// Alles, was sich auf dem LedGrid darstellen (d.h. zeichnen laesst),
// implementiert das Visual-Interface.
type Visual interface {
	Nameable
	// Mit diesen Methoden kann ermittelt, resp. festgelegt werden, ob das
	// Objekt dargestellt (d.h. gezeichnet) werden soll.
	Visible() bool
	SetVisible(v bool)
	// Zeichnet das Objekt auf dem LedGrid.
	image.Image
}

// Dieses Embedable kann fuer eine Default-Implementation des Drawable-
// Interfaces genutzt werden.
type VisualEmbed struct {
	NameableEmbed
	visible bool
}

func (v *VisualEmbed) Init(name string) {
	v.NameableEmbed.Init(name)
	v.visible = false
}

func (v *VisualEmbed) Visible() bool {
	return v.visible
}

func (v *VisualEmbed) SetVisible(vis bool) {
	v.visible = vis
}

// Einige der Objekte (wie beispielsweise Shader) koennen zusaetzlich mit
// Parametern gesteuert werden. Damit diese Steuerung so generisch wie
// moeglich ist, haben alle parametrisierbaren Typen dieses Interface zu
// implementieren.
type Parametrizable interface {
	ParamList() []Parameter
}

// Dieses Interface word von allen Objekten implementiert, die sich
// einfaerben lassen.
type Paintable interface {
	// Retourniert die Palette, welche zur Faerbung des Objektes hinterlegt
	// ist.
    PaletteParam() PaletteParameter
	Palette() ColorSource
	SetPalette(pal ColorSource, fadeTime time.Duration)
}

// Alles, was im Sinne einer Farbpalette Farben erzeugen kann, implementiert
// das ColorSource Interface.
type ColorSource interface {
	// Da diese Objekte auch oft in GUI angezeigt werden, muessen sie das
	// Nameable-Interface implementieren, d.h. einen Namen haben.
	Nameable
	// Liefert in Abhaengigkeit des Parameters v eine Farbe aus der Palette
	// zurueck. v kann vielfaeltig verwendet werden, bsp. als Parameter im
	// Intervall [0,1] oder als Index (natuerliche Zahl) einer Farbenliste
	// oder gar nicht, wenn die Farbquelle bspw. einfarbig ist.
	Color(v float64) LedColor
}

//----------------------------------------------------------------------------

type Text interface {
	Nameable
	Visual
	Parametrizable
	Paintable
	String() string
	SetString(s string)
}
