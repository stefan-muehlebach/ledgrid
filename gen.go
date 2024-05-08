//go:build ignore
// +build ignore

// Dieses Programm dient dem automatischen Aufbau der Farbpaletten aufgrund
// der in paletteColors.go definierten Farblisten. Mehr Info dazu: siehe
// Kommentare in paletteColors.go.

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

    "github.com/stefan-muehlebach/gg/colornames"
)

var (
	reListName = regexp.MustCompile(`^[[:space:]]*(([[:alpha:]]*)(Gradient|ColorList)(NonCyc)?)`)

	paletteNamesTemplate = `
// ACHTUNG: dieses File wird automatisch durch das Tool 'gen' in diesem
// Verzeichnis erzeugt! Manuelle Anpassungen koennen verloren gehen.

package ledgrid

import (
    "github.com/stefan-muehlebach/gg/colornames"
)

var (
    // PaletteList ist ein Slice mit den Namen aller vorhandenen Paletten.
    PaletteList = []Colorable{
{{- range $i, $row := .}}
        {{printf "%sPalette" $row.PaletteName}},
{{- end}}
    }
)

var (
    // In diesem Block werden die Paletten konkret erstellt.
{{- range $i, $row := .}}
    {{- if $row.IsGradient}}
    {{printf "%sPalette = NewGradientPalette(\"%[1]s\", %s...)" $row.PaletteName $row.ColorListName}}
    {{- else}}
    {{printf "%sPalette = NewGradientPaletteByList(\"%[1]s\", %v, %s...)" $row.PaletteName $row.Cycle $row.ColorListName}}
    {{- end}}
{{- end}}
)
`
	paletteTempl = template.Must(template.New("paletteNames").Parse(paletteNamesTemplate))

    colorNamesTemplate = `
var (
    ColorList = []Colorable{
{{- range $i, $row := .}}
        {{printf "%sColor" $row}},
{{- end}}
    }
)

var (
    // In diesem Block werden die uniformen Paletten erstellt.
{{- range $i, $row := .}}
    {{printf "%sColor = NewUniformPalette(\"%[1]s\", colornames.%[1]s)" $row}}
{{- end}}
)
`
    colorTempl = template.Must(template.New("colorNames").Parse(colorNamesTemplate))
)

type Record struct {
    PaletteName string
    ColorListName string
    IsGradient bool
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
                record.IsGradient = true
            } else {
                record.IsGradient = false
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
    defer fh.Close()
	err = paletteTempl.Execute(fh, nameList)
	if err != nil {
		log.Fatalf("executing template: %v", err)
	}
    err = colorTempl.Execute(fh, colornames.Names)
	if err != nil {
		log.Fatalf("executing template: %v", err)
	}
}
