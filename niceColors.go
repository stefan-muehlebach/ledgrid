package ledgrid

import (
	"embed"
	"encoding/json"
	"log"
	"path/filepath"
)

//go:embed data/*.json
var colorFiles embed.FS

type jsonPaletteRecord struct {
	ID       int
	Name     string    `json:"Title,Name"`
	IsCyclic bool
	Colors   []LedColor
	Stops    []ColorStop
}

type StopsPalette struct {
	ID    int
	Name  string
	Stops []ColorStop
}

type ListPalette struct {
	ID       int
	Name     string
	IsCyclic bool
	Colors   []LedColor
}

func ReadJsonPalette(fileName string) {
	var colorListJson []jsonPaletteRecord

	data, err := colorFiles.ReadFile(filepath.Join("data", fileName))
	if err != nil {
		log.Fatalf("ReadFile failed: %v", err)
	}
	err = json.Unmarshal(data, &colorListJson)
	if err != nil {
        if err, ok := err.(*json.SyntaxError); ok {
		    log.Fatalf("Unmarshal failed in %s: %+v at offset %d", fileName, err, err.Offset)
        } else {
            log.Fatalf("Unmarshal failed in %s: %+v (%T)", fileName, err, err)
        }
	}

    // log.Printf("ReadPixelPalettes(): %d entries unmarshalled", len(colorListJson))
	for _, rec := range colorListJson {
        // log.Printf("%+v", rec)
        if len(rec.Colors) > 0 {
            pal := NewGradientPaletteByList(rec.Name, rec.IsCyclic, rec.Colors...)
            PaletteMap[rec.Name] = pal
            PaletteList = append(PaletteList, pal)
        } else if len(rec.Stops) > 0 {
            pal := NewGradientPalette(rec.Name, rec.Stops...)
            PaletteMap[rec.Name] = pal
            PaletteList = append(PaletteList, pal)
        } else {
            log.Printf("Palette '%s' has no colors", rec.Name)
        }
	}
}

func init() {
	ReadJsonPalette("palettes.json")
	ReadJsonPalette("colourlovers.json")
}
