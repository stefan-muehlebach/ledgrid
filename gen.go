//go:build ignore

package main

import (
	"log"
	"os"
	"fmt"

	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/ledgrid"
)

func main() {
    fh, err := os.Create("colornames/colornames.go")
    if err != nil {
        log.Fatalf("Couldn't create file: %v", err)
    }
    defer fh.Close()

    fmt.Fprintf(fh, "package colornames\n\n")
    fmt.Fprintf(fh, "import (\n")
    fmt.Fprintf(fh, "    \"github.com/stefan-muehlebach/ledgrid\"\n")
    fmt.Fprintf(fh, ")\n\n")
    fmt.Fprintf(fh, "var (\n")
	for _, name := range color.Names {
		ledColor := ledgrid.LedColorModel.Convert(color.Map[name]).(ledgrid.LedColor)
		fmt.Fprintf(fh, "    %s = ledgrid.LedColor%+v\n", name, ledColor)
	}
    fmt.Fprintf(fh, "\n    Map = map[string]ledgrid.LedColor {\n")
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
