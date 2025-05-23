package colors

import (
	"math/rand/v2"
)

// Mit RandColor kann zufällig eine aus dem gesamten Sortiment der hier
// definierten Farben gewählt werden. Hilfreich für Tests, Beispielprogramme
// oder anderes.
func RandColor() LedColor {
	name := Names[rand.IntN(len(Names))]
	return Map[name]
}

// Mit RandGroupColor wird der Zufall auf eine bestimmte Farbgruppe beschraenkt.
func RandGroupColor(group ColorGroup) LedColor {
	nameList, ok := Groups[group]
	if !ok {
		return LedColor{A: 0xff}
	}
	name := nameList[rand.IntN(len(nameList))]
	return Map[name]
}
