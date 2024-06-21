//go:build ignore

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

    "github.com/stefan-muehlebach/gg/color"
)

var (
	reListName = regexp.MustCompile(`^[[:space:]]*(([[:alpha:]]*)(Gradient|ColorList)(NonCyc)?)`)

	paletteNamesTemplate = `// Code generated  DO NOT EDIT.

// ACHTUNG: dieses File wird automatisch durch das Tool 'gen' in diesem
// Verzeichnis erzeugt! Manuelle Anpassungen koennen verloren gehen.

package ledgrid

import (
    "github.com/stefan-muehlebach/gg/color"
)

var (
    // PaletteList ist ein Slice mit allen vorhandenen Paletten.
    PaletteList = []ColorSource{
{{- range $i, $row := .}}
        {{printf "%s" $row.PaletteName}},
{{- end}}
    }

    // PaletteMap ist ein Map um Paletten mit ihrem Namen anzusprechen.
    PaletteMap = map[string]ColorSource{
{{- range $i, $row := .}}
        {{printf "\"%s\": %[1]s" $row.PaletteName}},
{{- end}}
    }

)

var (
    // In diesem Block werden die Paletten konkret erstellt.
{{- range $i, $row := .}}
    {{- if $row.IsGradient}}
    {{printf "%s = NewGradientPalette(\"%[1]s\", %s...)" $row.PaletteName $row.ColorListName}}
    {{- else}}
    {{printf "%s = NewGradientPaletteByList(\"%[1]s\", %v, %s...)" $row.PaletteName $row.Cycle $row.ColorListName}}
    {{- end}}
{{- end}}
)
`
	paletteTempl = template.Must(template.New("paletteNames").Parse(paletteNamesTemplate))

    colorNamesTemplate = `
var (
    // ColorList ist ein Slice mit allen vorhandenen Paletten.
    ColorList = []ColorSource{
{{- range $i, $row := .}}
        {{printf "%s" $row.PaletteName}},
{{- end}}
    }

    // ColorMap ist ein Map um Paletten mit ihrem Namen anzusprechen.
    ColorMap = map[string]ColorSource{
{{- range $i, $row := .}}
        {{printf "\"%s\": %[1]s" $row.PaletteName}},
{{- end}}
    }
)

var (
    // In diesem Block werden die uniformen Paletten erstellt.
{{- range $i, $row := .}}
    {{printf "%s = NewUniformPalette(\"%[1]s\", %s)" $row.PaletteName $row.ColorListName}}
{{- end}}
)
`
    colorTempl = template.Must(template.New("colorNames").Parse(colorNamesTemplate))
)

type Record struct {
    PaletteName string
    ColorListName string
    IsUniform bool
    IsGradient bool
    Cycle bool
}

func main() {
	var paletteNameList, colorNameList []Record

	fh, err := os.Open("paletteColors.go")
	if err != nil {
		log.Fatalf("opening file: %v", err)
	}

	paletteNameList = make([]Record, 0)
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := scanner.Text()
		if reListName.MatchString(line) {
			matches := reListName.FindStringSubmatch(line)
            record := Record{}
			record.ColorListName = matches[1]
			record.PaletteName = fmt.Sprintf("%c%sPalette", matches[2][0]-('a'-'A'), matches[2][1:])
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
			paletteNameList = append(paletteNameList, record)
		}
    }
	fh.Close()
    slices.SortFunc(paletteNameList, func(a Record, b Record) int {
        return strings.Compare(a.PaletteName, b.PaletteName)
    })

	colorNameList = make([]Record, 0)
    for _, colName := range color.Names {
        record := Record{}
        record.ColorListName = fmt.Sprintf("color.%s", colName)
        record.PaletteName = fmt.Sprintf("%sColor", colName)
        record.IsUniform = true
        colorNameList = append(colorNameList, record)
    }
    slices.SortFunc(colorNameList, func(a Record, b Record) int {
        return strings.Compare(a.PaletteName, b.PaletteName)
    })

    // paletteNameList = append(paletteNameList, colorNameList...)

	fh, err = os.Create("paletteNames.go")
	if err != nil {
		log.Fatalf("creating file: %v", err)
	}
    defer fh.Close()
	err = paletteTempl.Execute(fh, paletteNameList)
	if err != nil {
		log.Fatalf("executing template: %v", err)
	}
    err = colorTempl.Execute(fh, colorNameList)
	if err != nil {
		log.Fatalf("executing template: %v", err)
	}
}
