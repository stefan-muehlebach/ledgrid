//go:build ignore

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/stefan-muehlebach/gg/color"
	ledcolor "github.com/stefan-muehlebach/ledgrid/color"
)

func main() {
	fh, err := os.Create("names.go")
	if err != nil {
		log.Fatalf("Couldn't create file: %v", err)
	}
	defer fh.Close()

	fmt.Fprintf(fh, "// Dieses File wird automatisch durch gen.go erstellt!\n")
	fmt.Fprintf(fh, "// Manuelle Aenderungen werden bei der naechsten Ausfuehrung\n")
	fmt.Fprintf(fh, "// automatisch ueberschrieben.\n\n")
	fmt.Fprintf(fh, "package color\n\n")
	fmt.Fprintf(fh, "var (\n")
	for _, name := range color.Names {
		ledColor := ledcolor.LedColorModel.Convert(color.Map[name]).(ledcolor.LedColor)
		fmt.Fprintf(fh, "    %s = LedColor%+v\n", name, ledColor)
	}
	fmt.Fprintf(fh, "\n    Map = map[string]LedColor {\n")
	for _, name := range color.Names {
		fmt.Fprintf(fh, "        \"%[1]s\": %[1]s,\n", name)
	}
	fmt.Fprintf(fh, "    }\n\n")
	fmt.Fprintf(fh, "    Names = []string{\n")
	for _, name := range color.Names {
		fmt.Fprintf(fh, "        \"%s\",\n", name)
	}
	fmt.Fprintf(fh, "    }\n\n")
	fmt.Fprintf(fh, ")\n")
}
