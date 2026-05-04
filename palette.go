package ledgrid

import (
	"embed"
	"image"
	"image/color"
	"log"
	"math"
	"path"

	"github.com/stefan-muehlebach/gg/colors"
)

// Alles, was im Sinne einer Farbpalette Farben erzeugen kann, implementiert
// das ColorSource Interface.
type ColorSource interface {
	// Liefert in Abhaengigkeit des Parameters v eine Farbe aus der Palette
	// zurueck. v kann vielfaeltig verwendet, resp. interpretiert werden,
	// bsp. als Parameter im Intervall [0,1], als Index (natuerliche Zahl)
	// einer Farbenliste oder gar nicht, wenn die Farbquelle einfarbig ist.
	Color(v float64) colors.RGBA
	// Da alle Paletten noch einen Namen haben, der bspw. in einem GUI- oder
	// TUI-Element dargestellt werden kann, existiert diese Methode.
	Name() string
}

var (
	// Alle vorhandenen Paletten sind in diesem Slice aufgefuehrt. Falls
	// applikatorisch weitere Paletten erzeugt werden, ist es Aufgabe der
	// Applikation, diesen Slice nachzufuehren.
	PaletteNames []string
	// Im Gegensatz zu [PaletteList] sind hier die Paletten unter ihrem
	// Namen abgelegt. Siehe auch Kommentar bei [PaletteList] betr.
	// Nachfuehrung.
	PaletteMap map[string]ColorSource
)

//go:embed data/*.json
var dataFS embed.FS

// Mit der Initialisierung des ledgrid-Packages werden u.a. auch alle Paletten
// aus dem embedded Verzeichnis 'data' eingelesen. Aktuell ist dies bloss die
// Datei "palNew.json", welche die bisherigen Paletten in einem JSON-Format
// fuehrt.
func init() {
    var palMap map[string]colors.Palette
    var err error

    PaletteMap = make(map[string]ColorSource)

	fh, err := dataFS.Open(path.Join("data", "palNew.json"))
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer fh.Close()

	PaletteNames, palMap, err = colors.ReadPaletteData(fh)
	if err != nil {
		log.Fatalf("error reading palette data: %v", err.Error())
	}
	for _, name := range PaletteNames {
		PaletteMap[name] = palMap[name]
	}
}

// Damit auch einzelne Farben wie Paletten verwendet werden koennen,
// existiert der Typ UniformPalette. Die Ueberlegungen dazu sind analog zum
// Typ [image.Uniform].
type UniformPalette struct {
	color colors.RGBA
	name  string
}

// Erstellt eine neue einfarbige Farbquelle mit gegebenem namen.
func NewUniformPalette(name string, color colors.RGBA) *UniformPalette {
	p := &UniformPalette{}
	p.color = color
	p.name = name
	return p
}

// Damit wird das ColorSource-Interface implementiert. Der Parameter [v] hat
// bei dieser Farbquelle keine Bedeutung und wird ignoriert.
func (p *UniformPalette) Color(v float64) colors.RGBA {
	return p.color
}

func (p *UniformPalette) Name() string {
	return p.name
}

func (p *UniformPalette) ColorModel() color.Model {
	return colors.RGBAModel
}

func (p *UniformPalette) Bounds() image.Rectangle {
	return image.Rect(math.MinInt, math.MinInt, math.MaxInt, math.MaxInt)
}

func (p *UniformPalette) At(x, y int) color.Color {
	return p.color
}

func (p *UniformPalette) Set(x, y int, c colors.RGBA) {}

// Mit diesem Typ kann ein fliessender Uebergang von einer Palette zu einer
// anderen realisiert werden.
type PaletteFader struct {
	Pals  [2]ColorSource
	T     float64
	alpha uint8
}

// Initialisiert wird der Fader mit der aktuell anzuzeigenden Palette. Der
// PaletteFader wird anschliessend anstelle der ueblichen Palette verwendet.
func NewPaletteFader(pal ColorSource) *PaletteFader {
	p := &PaletteFader{}
	p.Pals[0] = pal
	p.Pals[1] = nil
	p.T = 0.0
	p.alpha = 0xff
	return p
}

// Mit dieser Methode wird der aktuelle Farbwert retourniert. Damit
// implementiert der Fader das ColorSource-Interface und kann als Farbquelle
// verwendet werden - genau wie anderen Paletten-, resp. Farbtypen.
func (p *PaletteFader) Color(v float64) colors.RGBA {
	c := p.Pals[0].Color(v)
	if p.T > 0 {
		c2 := p.Pals[1].Color(v)
		c = c.Interpolate(c2, p.T)
	}
	c.A = p.alpha
	return c
}

func (p *PaletteFader) Name() string {
	if p.T > 0 {
		return p.Pals[1].Name()
	}
	return p.Pals[0].Name()
}

func (p *PaletteFader) AlphaPtr() *uint8 {
	return &p.alpha
}
