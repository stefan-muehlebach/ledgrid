//go:build ignore

package main

import (
	"log"
	"os"
	"fmt"

	"github.com/stefan-muehlebach/gg/color"
	ledcolor "github.com/stefan-muehlebach/ledgrid/color"
)

func main() {
    fh, err := os.Create("color/names.go")
    if err != nil {
        log.Fatalf("Couldn't create file: %v", err)
    }
    defer fh.Close()

    fmt.Fprintf(fh, "// Dieses File wird automatisch erstellt!\n")
    fmt.Fprintf(fh, "// Manuelle Aenderungen werden automatisch ueberschrieben.\n\n")
    fmt.Fprintf(fh, "package color\n\n")
    // fmt.Fprintf(fh, "import (\n")
    // fmt.Fprintf(fh, "    \"github.com/stefan-muehlebach/ledgrid\"\n")
    // fmt.Fprintf(fh, ")\n\n")
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
