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

// Mit RandGroupColor kann der Zufall eine bestimmte Farbgruppe beschraenkt
// werden.
func RandGroupColor(group ColorGroup) ledgrid.LedColor {
	nameList, ok := Groups[group]
	if !ok {
		return ledgrid.LedColor{0, 0, 0, 1}
	}
	name := nameList[rand.Int()%len(nameList)]
	return Map[name]
}
