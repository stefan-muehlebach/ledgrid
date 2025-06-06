// Package conf contains Types and Functions that help with the somehow
// weird configuration of modules and the cabeling. In this package, you'll
// get the mapping from pixel coordinates to index on the LED chain.

package conf

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"log"
	"os"
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

// There are two module types (with respect to the cabeling).
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

// Returns the index of a pixel within a certain type of module. pt is a
// point with coordinates in the range of ModuleDim. The returned value
// is a number in the interval [0,99] or -1 if the coordinates are outside
// the module.
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

// Returns the coordinates of a pixel with index idx. idx must be a number
// in the interval [0,99] and the returned coordinate is within a rectangle
// of size ModuleDim. If idx is outside the valid interval, the returned
// coordinates are (-1,-1).
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

// Returns the correct description of the module ("LR" or "RL").
func (m ModuleType) String() string {
	switch m {
	case ModLR:
		return "LR"
	case ModRL:
		return "RL"
	}
	return "(unknown)"
}

// This method is used to set the type of a module according to a description
// (a string) in a configuration file.
func (m *ModuleType) Set(v string) error {
	switch v {
	case "LR":
		*m = ModLR
	case "RL":
		*m = ModRL
	}
	return nil
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

func (r RotationType) String() string {
	switch r {
	case Rot000:
		return "0"
	case Rot090:
		return "90"
	case Rot180:
		return "180"
	case Rot270:
		return "270"
	}
	return "(unknown)"
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

// This type denotes a specific implementation of module, consisting of a
// module type and a rotation.
type Module struct {
	Type ModuleType
	Rot  RotationType
}

// Given two module types and four rotations, there are 8 kind of modules
// which can be used to form a LedGrid. There is a constant for each of them.
var (
	ModLR000 = Module{ModLR, Rot000}
	ModLR090 = Module{ModLR, Rot090}
	ModLR180 = Module{ModLR, Rot180}
	ModLR270 = Module{ModLR, Rot270}
	ModRL000 = Module{ModRL, Rot000}
	ModRL090 = Module{ModRL, Rot090}
	ModRL180 = Module{ModRL, Rot180}
	ModRL270 = Module{ModRL, Rot270}
)

// Returns the index of a point pt within this module. If pt is outside the
// range ModuleDim, the method returns -1. See [ModuleType.Index] for more
// information.
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

// Returns the coordinate of a pixel given its index within the chain, forming
// this module. See [ModuleType.Coord] for more information.
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

// Since a module consists of a module type and a rotation, the string
// representation of a module contains the type and the rotation, separated
// by a colon, like in "RL:180" or "LR:000":.
func (m Module) String() string {
	return fmt.Sprintf("%v:%v", m.Type, m.Rot)
}

// With this method, Module implements the Scanner interface of package fmt
// and allows module configurations to be read with functions like Scan,
// Scanf or Scanln.
// func (m *Module) Scan(state fmt.ScanState, verb rune) error {
// 	tok, err := state.Token(true, nil)
// 	if err != nil {
// 		return err
// 	}
// 	slc := strings.Split(string(tok), ":")
// 	m.Type.Set(slc[0])
// 	m.Rot.Set(slc[1])
// 	return nil
// }

// MarshalText implement the TextMarshaler interface of package encoding.
// This allows to configure modules using JSON files.
func (m Module) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

// UnmarshalText implement the TextUnmarshaler interface of package encoding.
// This allows to configure modules using JSON files.
func (m *Module) UnmarshalText(text []byte) error {
	slc := strings.Split(string(text), ":")
	m.Type.Set(slc[0])
	m.Rot.Set(slc[1])
	return nil
}

// This type finally puts all information together: what module is placed at
// column Col and row Row of the LedGrid. Col and Row denote the position of
// the module (not a pixel) within the LedGrid. Idx finally is the index of
// the first pixel of this module within the whole chain of pixels, therefore
// this number must be a multiple of ModuleDim.X*ModuleDim.Y.
type ModulePosition struct {
	Col int    `json:"Col"`
	Row int    `json:"Row"`
	Mod Module `json:"Mod"`
	Idx int    `json:"-"`
}

// Returns the enclosing rectangle of this module, specified in pixel
// coordinates.
func (m ModulePosition) Bounds() image.Rectangle {
	x0, y0 := m.Col*ModuleDim.X, m.Row*ModuleDim.Y
	x1, y1 := x0+ModuleDim.X, y0+ModuleDim.Y
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

	if size.X <= 0 || size.Y <= 0 ||
		size.X%ModuleDim.X != 0 || size.Y%ModuleDim.Y != 0 {
		log.Fatalf("Requested size of LED-Grid '%v' does not match with size of a module '%v'", size, ModuleDim)
	}
	cols, rows := size.X/ModuleDim.X, size.Y/ModuleDim.Y

	for row = range rows {
		for i := range cols {
			if row%2 == 0 {
				col = i
				mod = ModLR000
			} else {
				col = cols - i - 1
				mod = ModLR180
			}
			if col == cols-1 {
				mod = ModRL090
			}
			conf.AddModule(col, row, mod)
		}
	}
	return conf
}

//go:embed data/*.json
var customFiles embed.FS

func AllCustomFiles() (fileList []string) {
    entryList, _ := customFiles.ReadDir("data")
    for _, entry := range entryList {
        fileList = append(fileList, entry.Name())
    }
    return fileList
}

func Load(fileName string) ModuleConfig {
	var conf ModuleConfig

	data, err := customFiles.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Couldn't read file: %v", err)
	}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatalf("Couldn't decode json data: %v", err)
	}
	for i := range conf {
		conf[i].Idx = i * ModuleDim.X * ModuleDim.Y
	}
	return conf
}

// Speichert die Konfiguration in conf in der Datei fileName ab.
func (conf ModuleConfig) Save(fileName string) {
	data, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		log.Fatalf("Couldn't encode data: %v", err)
	}
	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		log.Fatalf("Couldn't write to file: %v", err)
	}
}

// Helps to build up a module configuration. Important: the Add's must
// be done along the LED chain. The configuration will be verified after
// each add.
func (conf *ModuleConfig) AddModule(col, row int, mod Module) {
	modPos := ModulePosition{Col: col, Row: row, Mod: mod, Idx: len(*conf) * ModuleDim.X * ModuleDim.Y}
	*conf = append(*conf, modPos)
	if err := conf.VerifyModule(len(*conf) - 1); err != nil {
		log.Fatalf("Can't add module: %v", err)
	}
}

func (conf ModuleConfig) VerifyModule(i int) error {
	if i == 0 {
		return nil
	}
	if i >= len(conf) {
		return errors.New(fmt.Sprintf("no module with index %d", i))
	}
	idxA := i*ModuleDim.X*ModuleDim.Y - 1
	idxB := idxA + 1
	ptA := conf.Coord(idxA)
	ptB := conf.Coord(idxB)
	dx := abs(ptA.X - ptB.X)
	dy := abs(ptA.Y - ptB.Y)
	if dx > 1 || dy > 1 {
		return errors.New(fmt.Sprintf("from module %d to %d: %v and %v are not adjacent", i-1, i, ptA, ptB))
	}
	return nil
}

func (conf ModuleConfig) Verify() error {
	for i := range conf[1:] {
		err := conf.VerifyModule(i + 1)
		if err != nil {
			return err
		}
	}
	return nil
}

// Returns the size of the LEDGrid as number of pixels.
func (conf ModuleConfig) Size() image.Point {
	size := image.Point{}
	for _, modPos := range conf {
		size.X = max(size.X, ModuleDim.X*(modPos.Col+1))
		size.Y = max(size.Y, ModuleDim.Y*(modPos.Row+1))
	}
	return size
}

// Returns the index of the position pt within the LED chain or -1 if this
// position is not on a module.
func (conf ModuleConfig) Index(pt image.Point) int {
	for _, modPos := range conf {
		if pt.In(modPos.Bounds()) {
			return modPos.Index(pt)
		}
	}
	return -1
}

// Returns the coordinates of the LED with index idx.
func (conf ModuleConfig) Coord(idx int) image.Point {
	modPos := conf[idx/(ModuleDim.X*ModuleDim.Y)]
	return modPos.Coord(idx)
}

// Returns true if position pt is really visible on this LEDGrid.
func (conf ModuleConfig) Contains(pt image.Point) bool {
	for _, modPos := range conf {
		if pt.In(modPos.Bounds()) {
			return true
		}
	}
	return false
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

// Mit dieser Methode kann der entsprechende CoordMap erstellt werden.
func (conf ModuleConfig) CoordMap() CoordMap {
	coordMap := make([]image.Point, len(conf)*ModuleDim.X*ModuleDim.Y)
	for idx := range coordMap {
		coordMap[idx] = conf.Coord(idx)
	}
	return coordMap
}
