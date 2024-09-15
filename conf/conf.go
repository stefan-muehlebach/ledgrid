// Package conf contains Types and Functions that help with the somehow
// weird configuration of modules and the cabeling. In this package, you'll
// get the mapping from pixel coordinates to index on the LED chain.

package conf

import (
	"fmt"
	"image"
	"log"
	"strings"
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

// This denotes the number of LEDs on a single module.
var (
	ModuleDim = image.Point{10, 10}
)

// There are two module types.
type ModuleType int

const (
	// Modules of this type start the cabeling at the top left corner and
	// end at the top right corner, therefore the cabeling runs from left to
	// right.
	ModLR ModuleType = iota
	// This module type start the cabeling at the top right corner and end at
	// the top left corner (running the cable from right to left).
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

func (m ModuleType) Index(pt image.Point) int {
	idx := 0
	if m == ModRL {
		pt.X = ModuleDim.X - 1 - pt.X
	}
	idx = ModuleDim.Y * pt.X
	if pt.X%2 == 0 {
		idx += pt.Y
	} else {
		idx += (ModuleDim.Y - 1 - pt.Y)
	}
	return idx
}

func (m ModuleType) Coord(idx int) image.Point {
	col, row := 0, 0
	x := idx / ModuleDim.Y
	y := idx % ModuleDim.Y
	if m == ModLR {
		col = x
	} else {
		col = (ModuleDim.X - 1 - x)
	}
	if x%2 == 0 {
		row = y
	} else {
		row = (ModuleDim.Y - 1 - y)
	}
	return image.Point{col, row}
}

// Each module can be rotated in steps of 90 degrees with respect to the base
// position (cable start/end point a the top row).
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

// With values of type Module, you describe a certain module type, rotated
// by a specific angle.
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

func (m Module) Index(pt image.Point) int {
	switch m.Rot {
	case Rot090:
		pt.X, pt.Y = (ModuleDim.X - 1 - pt.Y), pt.X
	case Rot180:
		pt.X, pt.Y = (ModuleDim.X - 1 - pt.X), (ModuleDim.Y - 1 - pt.Y)
	case Rot270:
		pt.X, pt.Y = pt.Y, (ModuleDim.Y - 1 - pt.X)
	}
	return m.Type.Index(pt)
}

func (m Module) Coord(idx int) image.Point {
	pt := m.Type.Coord(idx)
	switch m.Rot {
	case Rot090:
		pt.X, pt.Y = pt.Y, (ModuleDim.Y - 1 - pt.X)
	case Rot180:
		pt.X, pt.Y = (ModuleDim.X - 1 - pt.X), (ModuleDim.Y - 1 - pt.Y)
	case Rot270:
		pt.X, pt.Y = (ModuleDim.X - 1 - pt.Y), pt.X
	}
	return pt
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
// an welcher Stelle (Col, Row) innerhalb des gesamten Panels befindet und
// welches LED (Index innerhalb der gesamten Kette) am Anfang des Moduls steht.
type ModulePosition struct {
	Col, Row, Idx int
	Mod           Module
}

func (m ModulePosition) Bounds() image.Rectangle {
    x0, y0 := m.Col * ModuleDim.X, m.Row * ModuleDim.Y
    x1, y1 := x0 + ModuleDim.X, y0 + ModuleDim.Y
    return image.Rect(x0, y0, x1, y1)
}

func (m ModulePosition) Index(pt image.Point) int {
	pt.X, pt.Y = pt.X%ModuleDim.X, pt.Y%ModuleDim.Y
	return m.Idx + m.Mod.Index(pt)
}

func (m ModulePosition) Coord(idx int) image.Point {
	pt := image.Point{m.Col * ModuleDim.X, m.Row * ModuleDim.Y}
	return pt.Add(m.Mod.Coord(idx - m.Idx))
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
	// var err error

	if size.X < ModuleDim.X || size.Y < ModuleDim.Y ||
		size.X%ModuleDim.X != 0 || size.Y%ModuleDim.Y != 0 {
		log.Fatalf("Requested size of LED-Grid '%v' does not match with size of a module '%v'", size, ModuleDim)
	}
	cols, rows := size.X/ModuleDim.X, size.Y/ModuleDim.Y
	conf = make([]ModulePosition, 0)

	idx := 0
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
			conf = append(conf, ModulePosition{Col: col, Row: row, Idx: idx, Mod: mod})
			idx += ModuleDim.X * ModuleDim.Y
		}
	}
	return conf
}

// Bestimmt die Groesse des gesamten Panels in Anzahl Pixel
func (conf ModuleConfig) Size() image.Point {
	size := image.Point{}
	for _, modPos := range conf {
		size.X = max(size.X, ModuleDim.X * (modPos.Col+1))
		size.Y = max(size.Y, ModuleDim.Y * (modPos.Row+1))
	}
	return size
}

func (conf ModuleConfig) Index(pt image.Point) int {
    for _, modPos := range conf {
        if pt.In(modPos.Bounds()) {
            return modPos.Index(pt)
        }
    }
    return -1
}

func (conf ModuleConfig) Coord(idx int) image.Point {
    modPos := conf[idx / (ModuleDim.X * ModuleDim.Y)]
    return modPos.Coord(idx)
}

func (conf ModuleConfig) Contains(pt image.Point) bool {
    for _, modPos := range conf {
        if pt.In(modPos.Bounds()) {
            return true
        }
    }
    return false
}

func (conf ModuleConfig) Module(pt image.Point) ModulePosition {
    for _, modPos := range conf {
        if pt.In(modPos.Bounds()) {
            return modPos
        }
    }
    log.Fatalf("Coordinates not within this grid: %v", pt)
    return ModulePosition{}
}

// Mit diesem Typ koennen die Koordinaten der Pixel auf dem LEDGrid auf
// Indizes innerhalb der Lichterkette gemappt werden.
type IndexMap [][]int

// Mit dem Typ CoordMap koennen Indizes der Lichterkette auf Koordinaten
// auf dem LEDGrid gemapped werden.
type CoordMap []image.Point

// Erstellt ein Feld (Slice of Slices) fuer die direkte Uebersetzung von
// Pixel-Koordinaten zum Index innerhalb der Lichterkette. Das Feld ist
// Spalten-orientiert, damit ist die Verwendung der Koordinaten vergleichbar
// mit anderen Graphik-Funktionen.
func (conf ModuleConfig) IndexMap() IndexMap {
	var idxMap IndexMap

	size := conf.Size()
	idxMap = make([][]int, size.X)
	for col := range idxMap {
		idxMap[col] = make([]int, size.Y)
	}
	for row := range size.Y {
        for col := range size.X {
            idxMap[col][row] = conf.Index(image.Point{col, row})
        }
    }
	return idxMap
}

// Mit dieser Methode kann ausgehend von einem IndexMap der entsprechende
// CoordMap erstellt werden.
// TO DO: waere besser eine Methode von ModuleConfig.
func (conf ModuleConfig) CoordMap() CoordMap {
	coordMap := make([]image.Point, len(conf) * ModuleDim.X * ModuleDim.Y)
	for idx := range coordMap {
        coordMap[idx] = conf.Coord(idx)
	}
	return coordMap
}
