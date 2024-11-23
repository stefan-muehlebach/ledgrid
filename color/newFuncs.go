//go:build ignore

package color

import (
    "image/color"
    "math/rand/v2"
)

// Mit RandColor kann zufällig eine aus dem gesamten Sortiment der hier
// definierten Farben gewählt werden. Hilfreich für Tests, Beispielprogramme
// oder anderes.
func RandColor() color.RGBA {
	name := Names[rand.IntN(len(Names))]
	return Map[name]
}

// Mit RandGroupColor kann der Zufall eine bestimmte Farbgruppe beschraenkt
// werden.
func RandGroupColor(group ColorGroup) color.Color {
	nameList, ok := Groups[group]
	if !ok {
		return color.Black
	}
	name := nameList[rand.IntN(len(nameList))]
	return Map[name]
}
