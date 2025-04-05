//go:build ignore

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/stefan-muehlebach/gg/colors"
	ledcolors "github.com/stefan-muehlebach/ledgrid/colors"
)

func main() {
	fh, err := os.Create("names.go")
	if err != nil {
		log.Fatalf("Couldn't create file: %v", err)
	}
	defer fh.Close()

	fmt.Fprintf(fh, "// Code generated  DO NOT EDIT.\n")
	fmt.Fprintf(fh, "// Dieses File wird automatisch durch gen.go erstellt!\n")
	fmt.Fprintf(fh, "// Manuelle Aenderungen werden bei der naechsten Ausfuehrung\n")
	fmt.Fprintf(fh, "// automatisch ueberschrieben.\n\n")
	fmt.Fprintf(fh, "package colors\n\n")
	fmt.Fprintf(fh, "var (\n")
	for _, name := range colors.Names {
		ledColor := ledcolors.LedColorModel.Convert(colors.Map[name]).(ledcolors.LedColor)
		fmt.Fprintf(fh, "    %-20s = LedColor%+v\n", name, ledColor)
	}
	fmt.Fprintf(fh, "\n    Map = map[string]LedColor {\n")
	for _, name := range colors.Names {
		fmt.Fprintf(fh, "        \"%[1]s\": %[1]s,\n", name)
	}
	fmt.Fprintf(fh, "    }\n\n")
	fmt.Fprintf(fh, "    Names = []string{\n")
	for _, name := range colors.Names {
		fmt.Fprintf(fh, "        \"%s\",\n", name)
	}
	fmt.Fprintf(fh, "    }\n\n")
	fmt.Fprintf(fh, ")\n")
}
