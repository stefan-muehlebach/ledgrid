//go:build ignore

package main

import (
	"log"
	"os"
	"strings"
	"text/template"

	"golang.org/x/image/colornames"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	namesTemplate = `// Code generated  DO NOT EDIT.

package color

import (
    "image/color"
    	"golang.org/x/image/colornames"
)

// Dieses File wird automatisch durch genColor.go erstellt!
// Manuelle Aenderungen werden bei der naechsten Ausfuehrung
// automatisch ueberschrieben.

// Alle in der SVG 1.1 Spezifikation benannten Farben sind
// in diesem Package als Variablen definiert.
var (
{{- range $i, $row := .}}
    {{printf "%-24s = %s" $row.NewName $row.OldName}}
{{- end}}

    // Diese Farben tauchen im Style-Guide von Google zur Kommunikation von Go
    // auf und werden bspw. im GUI-Package 'adagui' fuer die Farben der
    // Bedienelemente verwendet.
	//GoGopherBlue             = color.RGBA{R:0.004, G:0.678, B:0.847, A:1}
	//GoLightBlue              = color.RGBA{R:0.369, G:0.788, B:0.890, A:1}
	//GoAqua                   = color.RGBA{R:0.000, G:0.635, B:0.622, A:1}
	//GoBlack                  = color.RGBA{R:0.000, G:0.000, B:0.000, A:1}
	//GoFuchsia                = color.RGBA{R:0.808, G:0.188, B:0.384, A:1}
	//GoYellow                 = color.RGBA{R:0.992, G:0.867, B:0.000, A:1}
	//GoTeal                   = color.RGBA{R:0.000, G:0.520, B:0.553, A:1}
	//GoDimGray                = color.RGBA{R:0.333, G:0.341, B:0.349, A:1}
	//GoIndigo                 = color.RGBA{R:0.251, G:0.169, B:0.337, A:1}
	//GoLightGray              = color.RGBA{R:0.859, G:0.851, B:0.839, A:1}

    // Map contains named colors defined in the SVG 1.1 spec.
    Map = map[string]color.RGBA{
    {{- range $i, $row := .}}
        {{printf "\"%s\": %[1]s," $row.NewName}}
    {{- end}}
    }

    // Der Slice 'Names' enthält die Namen aller Farben der SVG 1.1 Spezifikation.
    // Auf die Besonderheit betr. Gross-/Kleinschreibung ist weiter oben bereits
    // eingegangen worden. Jedes Element dieses Slices findet sich als Schlüssel
    // in der Variable 'Map'.
    Names = []string{
    {{- range $i, $row := .}}
        {{printf "\"%s\"," $row.NewName}}
    {{- end}}
    }
)
`
)

var (
	nameList = []string{
		"almond",
		"aquamarine",
		"blue",
		"blush",
		"brick",
		"brown",
		"chiffon",
		"coral",
		"cream",
		"cyan",
		"drab",
		"goldenrod",
		"gray",
		"green",
		"grey",
		"khaki",
		"lace",
		"magenta",
		"olive",
		"orange",
		"orchid",
		"pink",
		"puff",
		"purple",
		"red",
		"rose",
		"salmon",
		"salmon",
		"sea",
		"sky",
		"slate",
		"smoke",
		"spring",
		"steel",
		"turquoise",
		"violet",
		"whip",
		"white",
		"wood",
		"yellow",
	}
)

type TemplateData struct {
	NewName, OldName string
}

func main() {
	var replList []string
	var replacer *strings.Replacer
	var namesTempl *template.Template

	namesTempl = template.Must(template.New("names").Parse(namesTemplate))

	langTag := language.German
	titleCase := cases.Title(langTag)

	replList = make([]string, 2*len(nameList))
	for i, name := range nameList {
		replList[2*i] = name
		replList[2*i+1] = titleCase.String(name)
	}
	replacer = strings.NewReplacer(replList...)

	colorList := make([]TemplateData, len(colornames.Names))
	for i, name := range colornames.Names {
		capName := titleCase.String(name)
		oldName := "colornames." + capName
		newName := replacer.Replace(capName)
		colorList[i] = TemplateData{
			newName,
			oldName,
		}
	}

	fh, err := os.Create("color/newNames.go")
	if err != nil {
		log.Fatalf("creating file: %v", err)
	}
	defer fh.Close()
	err = namesTempl.Execute(fh, colorList)
	if err != nil {
		log.Fatalf("executing template: %v", err)
	}
}

// func main() {
//     fh, err := os.Create("color/names.go")
//     if err != nil {
//         log.Fatalf("Couldn't create file: %v", err)
//     }
//     defer fh.Close()

//     fmt.Fprintf(fh, "// Dieses File wird automatisch durch genColor.go erstellt!\n")
//     fmt.Fprintf(fh, "// Manuelle Aenderungen werden bei der naechsten Ausfuehrung\n")
//     fmt.Fprintf(fh, "// automatisch ueberschrieben.\n\n")
//     fmt.Fprintf(fh, "package color\n\n")
//     // fmt.Fprintf(fh, "import (\n")
//     // fmt.Fprintf(fh, "    \"github.com/stefan-muehlebach/ledgrid\"\n")
//     // fmt.Fprintf(fh, ")\n\n")
//     fmt.Fprintf(fh, "var (\n")
// 	for _, name := range color.Names {
// 		ledColor := ledcolor.LedColorModel.Convert(color.Map[name]).(ledcolor.LedColor)
// 		fmt.Fprintf(fh, "    %s = LedColor%+v\n", name, ledColor)
// 	}
//     fmt.Fprintf(fh, "\n    Map = map[string]LedColor {\n")
//     for _, name := range color.Names {
//         fmt.Fprintf(fh, "        \"%[1]s\": %[1]s,\n", name)
//     }
//     fmt.Fprintf(fh, "    }\n\n")
//     fmt.Fprintf(fh, "    Names = []string{\n")
//     for _, name := range color.Names {
//         fmt.Fprintf(fh, "        \"%s\",\n", name)
//     }
//     fmt.Fprintf(fh, "    }\n\n")
//     fmt.Fprintf(fh, ")\n")
// }
