//go:build ignore
// +build ignore

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"text/template"
)

var (
	reColorListName = regexp.MustCompile(`^[[:space:]]*(([[:alpha:]]*)Colors)`)

	namesTempl = `
// ACHTUNG: dieses File wird automatisch erzeugt

package ledgrid

var (
    // PaletteNames ist ein Slice mit den Namen aller vorhandenen Paletten.
    PaletteNames = []string{
{{range $i, $row := .}}        "{{(index $row 1)}}",
{{end}}    }

    // PaletteMap ist die Verbindung zwischen Palettenname und einer Variable
    // vom Typ Palette.
    PaletteMap = map[string]*Palette{
{{range $i, $row := .}}        "{{(index $row 1)}}": {{(index $row 1)}},
{{end}}    }

    // In diesem Block werden die Paletten konkret erstellt. Im Moment
    // koennen so nur Paletten mit aequidistanten Farbstuetzstellen
    // erzeugt werden.
{{range $i, $row := .}}    {{(index $row 1)}} = NewPaletteWithColors({{(index $row 0)}})
{{end}}
)
`
)

func main() {
	var nameList [][]string

	fh, err := os.Open("paletteColors.go")
	if err != nil {
		log.Fatal(err)
	}

	nameList = make([][]string, 0)
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := scanner.Text()
		if reColorListName.MatchString(line) {
			matches := reColorListName.FindStringSubmatch(line)
			colorListName := matches[1]
			paletteName := fmt.Sprintf("%c%s", matches[2][0]-('a'-'A'), matches[2][1:])
			// listName := fmt.Sprintf("%c%s", matches[1][0], matches[1][1:])
			// listName[0] = strings.ToUpper(listName[0])
			// fmt.Printf(">> %s, %s <<\n", colorListName, paletteName)
			nameList = append(nameList, []string{colorListName, paletteName})
		}
	}
    fh.Close()

    fh, err = os.Create("paletteNames.go")
	if err != nil {
		log.Fatal(err)
	}
	t := template.Must(template.New("names").Parse(namesTempl))
	err = t.Execute(fh, nameList)
	if err != nil {
		log.Fatalf("executing template: %v", err)
	}
    fh.Close()
}
