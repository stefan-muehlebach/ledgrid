package colornames

import (
    "math/rand/v2"
    "github.com/stefan-muehlebach/ledgrid"
)

// Mit RandColor kann zufällig eine aus dem gesamten Sortiment der hier
// definierten Farben gewählt werden. Hilfreich für Tests, Beispielprogramme
// oder anderes.
func RandColor() ledgrid.LedColor {
	name := Names[rand.Int()%len(Names)]
	return Map[name]
}
