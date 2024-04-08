//go:build ignore
// +build ignore

package main

import (
	"strings"
	"slices"
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"text/template"
)

var (
	reGradientListName = regexp.MustCompile(`^[[:space:]]*(([[:alpha:]]*)Gradient([[:alpha:]]*))`)

	namesTempl = `
//----------------------------------------------------------------------------
//
//   ACHTUNG: dieses File wird automatisch durch das Tool 'gen' in diesem
//   Verzeichnis erzeugt! Manuelle Anpassungen koennen verloren gehen.
//
//----------------------------------------------------------------------------

package ledgrid

var (
    // PaletteNames ist ein Slice mit den Namen aller vorhandenen Paletten.
    PaletteList = []Colorable{
{{- range $i, $row := .}}
        {{printf "%sPalette" (index $row 0)}},
{{- end}}
    }

    // PaletteMap ist die Verbindung zwischen Palettenname und einer Variable
    // vom Typ Palette.
    PaletteMap = map[string]Colorable{
{{- range $i, $row := .}}
        {{printf "\"%[1]s\": %[1]sPalette" (index $row 0)}},
{{- end}}
    }

    // In diesem Block werden die Paletten konkret erstellt. Im Moment
    // koennen so nur Paletten mit aequidistanten Farbstuetzstellen
    // erzeugt werden.
{{- range $i, $row := .}}
    {{printf "%sPalette = NewGradientPalette(\"%[1]s\", %s, %s...)" (index $row 0) (index $row 2) (index $row 1)}}
{{- end}}
)
`
	templ = template.Must(template.New("names").Parse(namesTempl))
)

func main() {
	var nameList [][]string
    var colorListName, paletteName, cycleFlag string

	fh, err := os.Open("paletteColors.go")
	if err != nil {
		log.Fatalf("opening file: %v", err)
	}

	nameList = make([][]string, 0)
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := scanner.Text()
		if reGradientListName.MatchString(line) {
			matches := reGradientListName.FindStringSubmatch(line)
			colorListName = matches[1]
			paletteName = fmt.Sprintf("%c%s", matches[2][0]-('a'-'A'), matches[2][1:])
            if matches[3] == "NoCycle" {
                cycleFlag = "false"
            } else {
                cycleFlag = "true"
            }
			nameList = append(nameList, []string{paletteName, colorListName, cycleFlag})
		}
	}
	fh.Close()

    slices.SortFunc(nameList, func(a []string, b []string) int {
        return strings.Compare(a[0], b[0])
    })

	fh, err = os.Create("paletteNames.go")
	if err != nil {
		log.Fatalf("creating file: %v", err)
	}
	err = templ.Execute(fh, nameList)
	if err != nil {
		log.Fatalf("executing template: %v", err)
	}
	fh.Close()
}
