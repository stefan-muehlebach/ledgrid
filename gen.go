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
	reListName = regexp.MustCompile(`^[[:space:]]*(([[:alpha:]]*)(Gradient|ColorList)(NonCyc)?)`)

	paletteNamesTemplate = `
// ACHTUNG: dieses File wird automatisch durch das Tool 'gen' in diesem
// Verzeichnis erzeugt! Manuelle Anpassungen koennen verloren gehen.

package ledgrid

var (
    // PaletteNames ist ein Slice mit den Namen aller vorhandenen Paletten.
    PaletteList = []Colorable{
{{- range $i, $row := .}}
        {{printf "%sPalette" $row.PaletteName}},
{{- end}}
    }

    // In diesem Block werden die Paletten konkret erstellt. Im Moment
    // koennen so nur Paletten mit aequidistanten Farbstuetzstellen
    // erzeugt werden.
{{- range $i, $row := .}}
    {{- if $row.GradientType}}
    {{printf "%sPalette = NewGradientPalette(\"%[1]s\", %s...)" $row.PaletteName $row.ColorListName}}
    {{- else}}
    {{printf "%sPalette = NewGradientPaletteByList(\"%[1]s\", %v, %s...)" $row.PaletteName $row.Cycle $row.ColorListName}}
    {{- end}}
{{- end}}
)
`
	templ = template.Must(template.New("paletteNames").Parse(paletteNamesTemplate))
)

type Record struct {
    PaletteName string
    ColorListName string
    GradientType bool
    Cycle bool
}

func main() {
	var nameList []Record

	fh, err := os.Open("paletteColors.go")
	if err != nil {
		log.Fatalf("opening file: %v", err)
	}

	nameList = make([]Record, 0)
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := scanner.Text()
		if reListName.MatchString(line) {
			matches := reListName.FindStringSubmatch(line)
            record := Record{}
			record.ColorListName = matches[1]
			record.PaletteName = fmt.Sprintf("%c%s", matches[2][0]-('a'-'A'), matches[2][1:])
            if matches[3] == "Gradient" {
                record.GradientType = true
            } else {
                record.GradientType = false
            }
            if matches[4] == "NonCyc" {
                record.Cycle = false
            } else {
                record.Cycle = true
            }
			nameList = append(nameList, record)
		}
    }
	fh.Close()

    slices.SortFunc(nameList, func(a Record, b Record) int {
        return strings.Compare(a.PaletteName, b.PaletteName)
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
