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
	reGradientListName = regexp.MustCompile(`^[[:space:]]*(([[:alpha:]]*)Gradient)`)

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
    PaletteNames = []string{
{{- range $i, $row := .}}
        {{printf "\"%s\"" (index $row 0)}},
{{- end}}
    }

    // PaletteMap ist die Verbindung zwischen Palettenname und einer Variable
    // vom Typ Palette.
    PaletteMap = map[string]Colorable{
{{- range $i, $row := .}}
        {{printf "\"%[1]s\": %[1]s" (index $row 0)}},
{{- end}}
    }

    // In diesem Block werden die Paletten konkret erstellt. Im Moment
    // koennen so nur Paletten mit aequidistanten Farbstuetzstellen
    // erzeugt werden.
{{- range $i, $row := .}}
    {{printf "%-20s = NewGradientPalette(%s...)" (index $row 0) (index $row 1)}}
{{- end}}
)
`
)

func main() {
	var nameList [][]string

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
			colorListName := matches[1]
			paletteName := fmt.Sprintf("%c%s", matches[2][0]-('a'-'A'), matches[2][1:])
			nameList = append(nameList, []string{paletteName, colorListName})
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
	t := template.Must(template.New("names").Parse(namesTempl))
	err = t.Execute(fh, nameList)
	if err != nil {
		log.Fatalf("executing template: %v", err)
	}
	fh.Close()

    fmt.Printf(">>> \u2502 <<<\n")
}
