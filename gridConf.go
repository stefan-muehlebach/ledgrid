package ledgrid

import (
	"errors"
	"fmt"
	"image"
	"log"
	"math"
	"strings"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/fonts"
	"github.com/stefan-muehlebach/gg/geom"
)

// Es zeichnet sich ab, dass groessere LED-Panels aus quadratischen Modulen
// mit 100 LEDs (10x10) aufgebaut werden. Da sich bei 100 LEDs Anfang und Ende
// der Lichterkette auf der selben Quadratseite befinden, muss man fuer die
// korrekte Verkabelung von groesseren LED-Panels ein wenig "proebeln".
// Urspruenglich war vorgesehen, diese Konfiguration durch ein aufrufendes
// Programm, resp. einen Benutzer vornehmen zu lassen. Mittlerweile ist klar,
// dass die optimale Anordnung der Module automatisch erstellt werden sollte.
// Trotzdem muss ein Anwender die Idee hinter dieser Konfiguration verstehen -
// schliesslich ist er für die Erstellung der Module und deren Verkabelung
// zustaendig.
//
//   +--------+--------+--------+
//   |I      O|I      O|I       |
//   |   LR   |   LR   |   RL   |
//   |        |        |O       |
//   +--------+--------+--------+
//   |        |        |I       |
//   |   LR   |   LR   |   RL   |
//   |O      I|O      I|O       |
//   +--------+--------+--------+
//   |I      O|I      O|I       |
//   |   LR   |   LR   |   RL   |
//   |        |        |O       |
//   +--------+--------+--------+
//
// Man benoetigt zwei Modultypen, die mit "LR", resp. "RL" bezeichnet werden.
// Bei LR" beginnt die Verkabelung links oben und endet rechts oben; bei "RL"
// beginnt die Verkabelung rechts oben und endet links oben. Mit 'I' und 'O'
// sind Input und Output der Lichterkette gekennzeichnet.
//
//   +--------+   +--------+
//   |I      O|   |O      I|
//   |   LR   |   |   RL   |
//   |        |   |        |
//   +--------+   +--------+
//
// Jedes Modul kann in 4 Positionen (0, 90, 180, 270 Grad) ausgerichtet werden.
// Grundzustand ist 0 Grad, d.h. Anfang und Ende der Lichterkette sind oben.
// Die Drehungen sind in math. positiver Richtung zu verstehen.
//
//   +--------+   +--------+
//   |O       |   |        |
//   |  LR:90 |   | RL:180 |
//   |I       |   |I      O|
//   +--------+   +--------+
//
// Eine rechteckige, flaechig deckende und beliebig vergroesserbare
// Konfiguration der Module erhaelt man nach folgendem Muster.
//
//   +--------+--------+...+--------+
//   |I      O|I      O|   |I       |
//   |  LR:0  |  LR:0  |   |  RL:90 |
//   |        |        |   |O       |
//   +--------+--------+...+--------+
//   |        |        |   |I       |
//   | LR:180 | LR:180 |   |  RL:90 |
//   |O      I|O      I|   |O       |
//   +--------+--------+...+--------+
//   |I      O|I      O|   |I       |
//   |  LR:0  |  LR:0  |   |  RL:90 |
//   |        |        |   |O       |
//   +--------+--------+...+--------+
//                         .        .

// Die Anzahl LED's pro Modul in Breite und Hoehe (nicht quadratische Module
// sind theoretisch moeglich, wurden jedoch noch nie gebaut und folglich nicht
// getestet).
var (
	ModuleSize = image.Point{10, 10}
)

// Die Modul-Typen werden mit diesem Datentyp spezifiziert.
type ModuleType int

const (
	ModLR ModuleType = iota
	ModRL
)

func (m ModuleType) String() (s string) {
	switch m {
	case ModLR:
		s = "LR"
	case ModRL:
		s = "RL"
	}
	return
}

// Die Set-Methode ist eigentlich etwas aus dem Setter-Interface des Packages
// 'flag', wird hier aber verwendet um bestimmte Konfigurationen einfacher
// verarbeiten zu koennen.
func (m *ModuleType) Set(v string) error {
	switch v {
	case "LR":
		*m = ModLR
	case "RL":
		*m = ModRL
	}
	return nil
}

// In welcher Position sich ein Modul befindet, wird ueber diesen Typen
// festgeleg.
type RotationType int

const (
	Rot000 RotationType = 0
	Rot090 RotationType = 90
	Rot180 RotationType = 180
	Rot270 RotationType = 270
)

func (r RotationType) String() (s string) {
	switch r {
	case Rot000:
		s = "0"
	case Rot090:
		s = "90"
	case Rot180:
		s = "180"
	case Rot270:
		s = "270"
	}
	return
}

// (Siehe Kommentar zu Set von ModuleType)
func (r *RotationType) Set(v string) error {
	switch v {
	case "0":
		*r = Rot000
	case "90":
		*r = Rot090
	case "180":
		*r = Rot180
	case "270":
		*r = Rot270
	}
	return nil
}

// Ein konkretes Modul wird durch den Modul-Typ (LR oder RL) und die Rotation
// beschrieben.
type Module struct {
	Type ModuleType
	Rot  RotationType
}

// Die textuelle Darstellung eines Moduls ist in der Einleitung am Anfang des
// Packages zu sehen: Modul-Typ und Rotationsart werden mit Doppelpunkt
// getrennt als zusammenhaengende Zeichenkette dargestell.
func (m Module) String() string {
	return fmt.Sprintf("%v:%v", m.Type, m.Rot)
}

// Module implementiert das Scanner-Interface, ist also in der Lage, via
// Funktion aus der Scanf-Familie eine konkrete Modul-Spezifikation zu lesen.
func (m *Module) Scan(state fmt.ScanState, verb rune) error {
	tok, err := state.Token(true, nil)
	if err != nil {
		return err
	}
	slc := strings.Split(string(tok), ":")
	m.Type.Set(slc[0])
	m.Rot.Set(slc[1])
	return nil
}

// Mit diesem Typ wird festgehalten, welches Modul (Typ und Ausrichtung) sich
// an welcher Stelle (Col, Row) innerhalb des gesamten Panels befindet.
type ModulePosition struct {
	Col, Row int
	Mod      Module
}

// Der Typ ModuleConfig schliesslich dient dazu, eine komplette
// Modul-Konfiguration zu speichern. Die Reihenfolge der Module ist relevant
// und entspricht der Verkabelung (d.h. die Einspeisung beginnt beim Modul
// an Position [0], geht dann zum Modul an Position [1] weiter, etc.)
type ModuleConfig []ModulePosition

func DefaultModuleConfig(size image.Point) ModuleConfig {
	var col, row int
	var conf ModuleConfig
	var mod Module
	var err error

	if size.X < ModuleSize.X || size.Y < ModuleSize.Y ||
		size.X%ModuleSize.X != 0 || size.Y%ModuleSize.Y != 0 {
		log.Fatalf("Requested size of LED-Grid '%v' does not match with size of a module '%v'", size, ModuleSize)
	}
	cols, rows := size.X/ModuleSize.X, size.Y/ModuleSize.Y
	conf = make([]ModulePosition, 0)

	for row = range rows {
		for i := range cols {
			if row%2 == 0 {
				col = i
				mod = Module{ModLR, Rot000}
			} else {
				col = cols - i - 1
				mod = Module{ModLR, Rot180}
			}
			if col == cols-1 {
				mod = Module{ModRL, Rot090}
			}
			conf, err = conf.Append(col, row, mod)
			if err != nil {
				log.Fatalf("Couldn't append module: %v", err)
			}
		}
	}
	return conf
}

// Fuegt der Modul-Konfiguration hinter c ein weiteres Modul hinzu und (falls
// erfolgreich) retourniert die neue Modul-Konfiguration. Falls ein Fehler
// erkannt wird, wird die alte Modul-Konfiguration retourniert und der zweite
// Rueckgabewert beschreibt die Art des Fehlers.
func (conf ModuleConfig) Append(col, row int, mod Module) (ModuleConfig, error) {
	modPos := ModulePosition{col, row, mod}
	if len(conf) != 0 {
		for _, pos := range conf {
			if pos.Col == col && pos.Row == row {
				return conf, errors.New(fmt.Sprintf("Position (%d,%d) is already occupied", col, row))
			}
		}
		lastModPos := conf[len(conf)-1]
		if abs(lastModPos.Col-col)+abs(lastModPos.Row-row) != 1 {
			return conf, errors.New(fmt.Sprintf("Module %v is not adjacent to last module %v", modPos, lastModPos))
		}
	}
	return append(conf, modPos), nil
}

// Bestimmt die Groesse des gesamten Panels in Anzahl Modulen in X-, resp.
// Y-Richtung.
func (conf ModuleConfig) Size() image.Point {
	size := image.Point{}
	for _, pos := range conf {
		size.X = max(size.X, pos.Col+1)
		size.Y = max(size.Y, pos.Row+1)
	}
	return size
}

// Mit dieser Struktur wird die Koordinate einer LED auf dem Panel auf den
// Index dieser LED innerhalb der Lichterkette gemappt.
type IndexMap [][]int

type CoordMap []image.Point

// Erstellt ein Feld (Slice of slice) fuer die direkte Uebersetzung von
// Pixel-Koordinaten zur Position (Index) innerhalb der Lichterkette.
func (conf ModuleConfig) IndexMap() IndexMap {
	var idxMap IndexMap

	size := conf.Size()
	idxMap = make([][]int, size.X*ModuleSize.X)
	for col := range idxMap {
		idxMap[col] = make([]int, size.Y*ModuleSize.Y)
	}
	idx := 0
	for _, pos := range conf {
		pt := image.Point{pos.Col * ModuleSize.X, pos.Row * ModuleSize.Y}
		mod := pos.Mod
		idx = idxMap.Append(mod, pt, idx)
	}
	return idxMap
}

func (idxMap IndexMap) CoordMap() CoordMap {
    coordMap := make([]image.Point, len(idxMap)*len(idxMap[0]))
    for col, idxColumn := range idxMap {
        for row, idx := range idxColumn {
            coordMap[idx/3] = image.Point{col, row}
        }
    }
    return coordMap
}

// Fuer die graphische Ausgabe des Verkabelungsplanes werden viele Konstanten
// verwendet.

var (
	marginSize = 100.0

	moduleSize        = 400.0
	moduleBorderWidth = 3.0
	moduleBorderColor = color.Black
	moduleFillColor   = color.Beige

	moduleTextFont  = fonts.GoBold
	moduleTextSize  = 60.0
	moduleTextColor = color.Gainsboro

	ledFieldSize        = moduleSize / float64(ModuleSize.X)
	ledSize             = ledFieldSize - 15.0
	ledInputFieldColor  = color.OrangeRed.Alpha(0.3)
	ledOutputFieldColor = color.Teal.Alpha(0.3)
	ledBorderWidth      = 1.0
	ledBorderColor      = color.Black
	ledFillColor        = color.White

	arrowSize  = ledSize / 4.0
	arrowColor = color.Black
)

func (mod Module) Draw(gc *gg.Context) {
	// Referenzpunkt links oben...
	p0 := geom.Point{-moduleSize / 2.0, -moduleSize / 2.0}
	// ... und rechts oben
	p1 := p0.Add(geom.Point{moduleSize, 0})

	mp := geom.Point{}
	dp := geom.Point{}

	// Feldfuellung fuer das Modul
	gc.DrawRectangle(p0.X, p0.Y, moduleSize, moduleSize)
	gc.SetFillColor(moduleFillColor)
	gc.Fill()

	// Index in der Mitte des Feldes plus Bezeichnung des Modules und
	// seiner Ausrichtung.
	gc.SetFontFace(fonts.NewFace(moduleTextFont, moduleTextSize))
	gc.SetStrokeColor(moduleTextColor)
	gc.DrawStringAnchored(fmt.Sprintf("%v", mod), 0.0, 0.0, 0.5, 0.5)

	// Farbliche Hervorhebung des Feldes links oben
	gc.DrawRectangle(p0.X, p0.Y, ledFieldSize, ledFieldSize)
	if mod.Type == ModLR {
		gc.SetFillColor(ledInputFieldColor)
	} else {
		gc.SetFillColor(ledOutputFieldColor)
	}
	gc.Fill()

	// Farbliche Hervorhebung des Feldes rechts oben
	gc.DrawRectangle(p1.X-ledFieldSize, p1.Y, ledFieldSize, ledFieldSize)
	if mod.Type == ModRL {
		gc.SetFillColor(ledInputFieldColor)
	} else {
		gc.SetFillColor(ledOutputFieldColor)
	}
	gc.Fill()

	// Erste Spalte mit LED
	for i := range 2 {
		if mod.Type == ModLR {
			mp = p0.AddXY(ledFieldSize/2.0, ledFieldSize/2.0)
			if i == 1 {
				mp = mp.AddXY(moduleSize-ledFieldSize, moduleSize-ledFieldSize)
			}
		} else {
			mp = p1.AddXY(-ledFieldSize/2.0, ledFieldSize/2.0)
			if i == 1 {
				mp = mp.AddXY(-(moduleSize - ledFieldSize), moduleSize-ledFieldSize)
			}
		}
		if i == 0 {
			dp = geom.Point{0, ledFieldSize}
		} else {
			dp = geom.Point{0, -ledFieldSize}
		}
		for j := range ModuleSize.Y {
			gc.DrawCircle(mp.X, mp.Y, ledSize/2.0)
			gc.SetStrokeWidth(ledBorderWidth)
			gc.SetStrokeColor(ledBorderColor)
			gc.SetFillColor(ledFillColor)
			gc.FillStroke()
			if j > 0 {
				prevPt := mp.Sub(dp)
				arrowPt := mp.Interpolate(prevPt, 0.53)
				angle := prevPt.Sub(mp).Angle() - math.Pi/2.0
				gc.DrawRegularPolygon(3, arrowPt.X, arrowPt.Y, arrowSize, angle)
				gc.SetFillColor(arrowColor)
				gc.Fill()
			}
			mp = mp.Add(dp)
		}
	}

	// Abschliessend die Modul-Umrahmung.
	gc.DrawRectangle(p0.X, p0.Y, moduleSize, moduleSize)
	gc.SetStrokeWidth(moduleBorderWidth)
	gc.SetStrokeColor(moduleBorderColor)
	gc.Stroke()
}

func (conf ModuleConfig) Draw(gc *gg.Context) {
	p0 := geom.Point{marginSize, marginSize}.AddXY(moduleSize/2.0, moduleSize/2.0)
	for _, modPos := range conf {
		pt := p0.Add(geom.Point{float64(modPos.Col), float64(modPos.Row)}.Mul(moduleSize))

		gc.Push()
		gc.Translate(pt.X, pt.Y)
		gc.Rotate(math.Pi * float64(-modPos.Mod.Rot) / 180.0)
		modPos.Mod.Draw(gc)
		gc.Pop()
	}
}

func (conf ModuleConfig) Plot(fileName string) {
	size := conf.Size()
	gc := gg.NewContext(size.X*int(moduleSize)+2*int(marginSize),
		size.Y*int(moduleSize)+2*int(marginSize))
	gc.SetFillColor(color.White)
	gc.Clear()

    conf.Draw(gc)

	err := gc.SavePNG(fileName)
	if err != nil {
		log.Fatalf("Couldn't save configuration: %v", err)
	}
}

// Hilfsfunktioenchen (sogar generisch!)
func abs[T ~int | ~float64](i T) T {
	if i < 0 {
		return -i
	} else {
		return i
	}
}

// Mit dieser Methode kann eine bestimmte Position im LED-Panel als 'defekt'
// markiert werden. In der Lichterkette muss diese Position ueberbrueckt,
// d.h. die entsprechende LED entfernt und die Anschlusskabel direkt
// miteinander verbunden werden.
// func (idxMap IndexMap) MarkDefect(pos image.Point) {
// 	idxDefect := idxMap[pos.X][pos.Y]
// 	cols := len(idxMap)
// 	rows := len(idxMap[0])
// 	for col := range cols {
// 		for row := range rows {
// 			if idxMap[col][row] > idxDefect {
// 				idxMap[col][row] -= 3
// 			}
// 		}
// 	}
// 	idxSpare := 3 * (cols*rows - 1)
// 	idxMap[pos.X][pos.Y] = idxSpare
// }

// Diese Methode ergaenzt den Slice idxMap um die Koordinaten und Indizes des
// Modules m. basePt sind die Pixel-Koordinaten der linken oberen Ecke des
// Moduls und baseIdx ist der Index der ersten LED des Moduls.
// Der Rueckgabewert ist der Index der ersten LED des nachfolgenden Moduls.
func (idxMap IndexMap) Append(m Module, basePt image.Point, baseIdx int) int {
	var idx int

	switch m.Type {
	case ModLR:
		switch m.Rot {
		case Rot000:
			for row := range ModuleSize.Y {
				y := basePt.Y + row
				for col := range ModuleSize.X {
					x := basePt.X + col
					if col%2 == 0 {
						idx = (col * ModuleSize.Y) + row
					} else {
						idx = (col * ModuleSize.Y) + (ModuleSize.Y - row - 1)
					}
					idxMap[x][y] = baseIdx + 3*idx
				}
			}
		case Rot090:
			for row := range ModuleSize.Y {
				y := basePt.Y + row
				for col := range ModuleSize.X {
					x := basePt.X + col
					if row%2 == 0 {
						idx = ((ModuleSize.X - row - 1) * ModuleSize.Y) + (ModuleSize.Y - col - 1)
					} else {
						idx = ((ModuleSize.X - row - 1) * ModuleSize.Y) + col
					}
					idxMap[x][y] = baseIdx + 3*idx
				}
			}
		case Rot180:
			for row := range ModuleSize.Y {
				y := basePt.Y + row
				for col := range ModuleSize.X {
					x := basePt.X + col
					if col%2 == 0 {
						idx = ((ModuleSize.X - col - 1) * ModuleSize.Y) + row
					} else {
						idx = ((ModuleSize.X - col - 1) * ModuleSize.Y) + (ModuleSize.Y - row - 1)
					}
					idxMap[x][y] = baseIdx + 3*idx
				}
			}
		case Rot270:
			for row := range ModuleSize.Y {
				y := basePt.Y + row
				for col := range ModuleSize.X {
					x := basePt.X + col
					if row%2 == 0 {
						idx = (row * ModuleSize.Y) + (ModuleSize.Y - col - 1)
					} else {
						idx = (row * ModuleSize.Y) + col
					}
					idxMap[x][y] = baseIdx + 3*idx
				}
			}
			// default:
			// 	log.Fatalf("Module type '%s' is only configured with a rotation of 0 and 180 degrees", m.Type)
		}
	case ModRL:
		switch m.Rot {
		case Rot000:
			for row := range ModuleSize.Y {
				y := basePt.Y + row
				for col := range ModuleSize.X {
					x := basePt.X + col
					if col%2 == 0 {
						idx = ((ModuleSize.X - col - 1) * ModuleSize.Y) + (ModuleSize.Y - row - 1)
					} else {
						idx = ((ModuleSize.X - col - 1) * ModuleSize.Y) + row
					}
					idxMap[x][y] = baseIdx + 3*idx
				}
			}
		case Rot090:
			for row := range ModuleSize.Y {
				y := basePt.Y + row
				for col := range ModuleSize.X {
					x := basePt.X + col
					if row%2 == 0 {
						idx = (row * ModuleSize.Y) + col
					} else {
						idx = (row * ModuleSize.Y) + (ModuleSize.X - col - 1)
					}
					idxMap[x][y] = baseIdx + 3*idx
				}
			}
		case Rot180:
			for row := range ModuleSize.Y {
				y := basePt.Y + row
				for col := range ModuleSize.X {
					x := basePt.X + col
					if col%2 == 0 {
						idx = (col * ModuleSize.Y) + (ModuleSize.Y - row - 1)
					} else {
						idx = (col * ModuleSize.Y) + row
					}
					idxMap[x][y] = baseIdx + 3*idx
				}
			}
		case Rot270:
			for row := range ModuleSize.Y {
				y := basePt.Y + row
				for col := range ModuleSize.X {
					x := basePt.X + col
					if row%2 == 0 {
						idx = ((ModuleSize.X - row - 1) * ModuleSize.Y) + col
					} else {
						idx = ((ModuleSize.X - row - 1) * ModuleSize.Y) + (ModuleSize.X - col - 1)
					}
					idxMap[x][y] = baseIdx + 3*idx
				}
			}
			// default:
			// 	log.Fatalf("Module type '%s' is only configured with a rotation of 90 degrees", m.Type)
		}
	}
	return baseIdx + 3*ModuleSize.X*ModuleSize.Y
}
