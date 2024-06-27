package ledgrid

import (
	"fmt"
	"image"
	"log"
	"strings"
)

// Die Anzahl LED's pro Modul in Breite und Hoehe.
var (
	ModuleSize = image.Point{10, 10}
)

// Es zeichnet sich ab, dass groessere LED-Panels aus modularen Teilen
// (z.B. 10x10 grossen Quadraten) aufgebaut werden. Da sich bei 100 LEDs
// Anfang und Ende der Licherkette auf der selben Quadratseite befinden, muss
// man fuer groessere LED-Panels ein wenig "proebeln". Man benoetigt zwei
// Modultypen: "LR" und "RL" bezeichnet: bei "LR" beginnt die Verkabelung
// links oben und endet rechts oben, bei "RL" beginnt die Verkabelung rechts
// oben und endet links oben. Jedes Modul kann in 4 Positionen (0, 90, 180,
// 270 Grad) ausgerichtet werden. Grundzustand ist 0, Anfang und Ende der
// Lichterkette sind oben. Die Drehungen sind in math. positiver Richtung
// zu verstehen. Damit laesst sich die Konfiguration der Module beschreiben:
//
//    LR:0   LR:0   LR:0   RL:90
//    LR:180 LR:180 LR:180 RL:90
//    LR:0   LR:0   LR:0   RL:90
//
// Die Software prueft, ob diese Konfiguration sinnvoll ist, d.h. ob damit
// eine durchgaengige Verkabelung realisiert werden kann.

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
	Rot000 RotationType = iota
	Rot090
	Rot180
	Rot270
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

// Ein konkretes Modul besteht aus einem Modul-Typ (LR oder RL) der sich in
// einer bestimmten Lage (Rotation) befindet.
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

// Module implementier das Scanner-Interface, ist also in der Lage, via
// Funktion aus der Scanf-Familie eine konkrete Modul-Spezifikation zu lesen.
func (m *Module) Scan(state fmt.ScanState, verb rune) error {
	//log.Printf("In Module.Scan()")
	tok, err := state.Token(true, nil)
	if err != nil {
		return err
	}
	//log.Printf("   tokens: %s", tok)
	slc := strings.Split(string(tok), ":")
	m.Type.Set(slc[0])
	m.Rot.Set(slc[1])
	return nil
}

func (m *Module) AppendIdxMap(idxMap [][]int, basePt image.Point, baseIdx int) int {
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
					idxMap[x][y] = baseIdx + 3 * idx
				}
			}
		case Rot180:
			for row := range ModuleSize.Y {
				y := basePt.Y + row
				for col := range ModuleSize.X {
					x := basePt.X + col
					if col%2 == 0 {
						idx = ((ModuleSize.X - col - 1) * ModuleSize.Y) + (ModuleSize.Y - row - 1)
					} else {
						idx = ((ModuleSize.X - col - 1) * ModuleSize.Y) + row
					}
					idxMap[x][y] = baseIdx + 3 * idx
				}
			}
		default:
			log.Fatalf("Module type '%s' is only configured with a rotation of 0 and 180 degrees", m.Type)
		}
	case ModRL:
		switch m.Rot {
		case Rot090:
			for row := range ModuleSize.Y {
				y := basePt.Y + row
				for col := range ModuleSize.X {
					x := basePt.X + col
					if row%2 == 0 {
						idx = (row * ModuleSize.X) + col
					} else {
						idx = (row * ModuleSize.X) + (ModuleSize.X - col - 1)
					}
					idxMap[x][y] = baseIdx + 3 * idx
				}
			}
		default:
			log.Fatalf("Module type '%s' is only configured with a rotation of 90 degrees", m.Type)
		}
	}
	return baseIdx + 3*ModuleSize.X*ModuleSize.Y
}

// Das Zusamenfuehren von mehreren Modulen zu einem LED-Grid wird ueber den
// Typ ModuleLayout realisiert. Die Organisation der Module ist zeilenbasiert,
// d.h. ueber eine Variable gridLay des Typs GridLayout und der Zeile, resp.
// Spalte (row, col) wird mit gridLay[row][col] auf das entsprechende Modul
// zugegriffen.
type ModuleLayout [][]Module

// Erstellt aufgrund der vorgegebenen Groesse der LED-Grids eine Modul-
// Konfiguration, welche die Flaeche lueckenlos deckt, beim Pixel (0,0) (d.h.
// links oben beginnt) und eine minimale Kabellaenge aufweist. Zentral ist
// die Variable ModuleSize, welche die Groesse eines einzelnen Moduls angibt.
// Es koennen nur LED-Grids erstellt werden, deren Groesse ein Vielfaches der
// Modul-Groesse ist.
func NewModuleLayout(size image.Point) ModuleLayout {
    if size.X < ModuleSize.X || size.Y < ModuleSize.Y || size.X % ModuleSize.X != 0 || size.Y % ModuleSize.Y != 0 {
        log.Fatalf("Requested size of LED-Grid '%v' does not match with size of a module '%v'", size, ModuleSize)
    }
    cols, rows := size.X / ModuleSize.X, size.Y / ModuleSize.Y
    lay := make([][]Module, rows)
    for row := range rows {
        lay[row] = make([]Module, cols)
    }

    for row := range rows {
        for i := range cols {
            col := cols - i - 1
            if i == 0 {
                lay[row][col] = Module{ModRL, Rot090}
            } else {
                if row % 2 == 0 {
                    lay[row][col] = Module{ModLR, Rot000}
                } else {
                    lay[row][col] = Module{ModLR, Rot180}
                }
            }
        }
    }

    return lay
}
