package ledgrid

import (
	"embed"
	"encoding/json"
	"log"
	"path"

	"github.com/stefan-muehlebach/ledgrid/color"
)

//go:embed data/*.json
var colorFiles embed.FS

type JsonPalette struct {
	ID       int
	Name     string `json:"Title"`
	IsCyclic bool
	IsSlice  bool
	Colors   []color.LedColor
	Stops    []ColorStop
}

func ReadJsonData(fileName string) []JsonPalette {
	var jsonPaletteList []JsonPalette

	data, err := colorFiles.ReadFile(path.Join("data", fileName))
	if err != nil {
		log.Fatalf("ReadFile failed: %v", err)
	}
	err = json.Unmarshal(data, &jsonPaletteList)
	if err != nil {
		if err, ok := err.(*json.SyntaxError); ok {
			log.Fatalf("Unmarshal failed in %s: %+v at offset %d", fileName, err, err.Offset)
		} else {
			log.Fatalf("Unmarshal failed in %s: %+v (%T)", fileName, err, err)
		}
	}
	return jsonPaletteList
}

func ReadJsonPalette(fileName string) {
	var jsonPaletteList []JsonPalette

    jsonPaletteList = ReadJsonData(fileName)

	// log.Printf("ReadPixelPalettes(): %d entries unmarshalled", len(colorListJson))
	for _, rec := range jsonPaletteList {
		// log.Printf("%+v", rec)
		if len(rec.Colors) > 0 {
			if rec.IsSlice {
				pal := NewSlicePalette(rec.Name, rec.Colors...)
				PaletteNames = append(PaletteNames, rec.Name)
				PaletteMap[rec.Name] = pal
			} else {
				pal := NewGradientPaletteByList(rec.Name, rec.IsCyclic, rec.Colors...)
				PaletteNames = append(PaletteNames, rec.Name)
				PaletteMap[rec.Name] = pal
			}
			// PaletteList = append(PaletteList, pal)
		} else if len(rec.Stops) > 0 {
			pal := NewGradientPalette(rec.Name, rec.Stops...)
			PaletteNames = append(PaletteNames, rec.Name)
			PaletteMap[rec.Name] = pal
			// PaletteList = append(PaletteList, pal)
		} else {
			log.Printf("Palette '%s' has no colors", rec.Name)
		}
	}
}

func ReadNamedColors() {
	for _, colorName := range color.Names {
		ColorNames = append(ColorNames, colorName)
		pal := NewUniformPalette(colorName, color.Map[colorName])
		ColorMap[colorName] = pal
	}
	colorName := "Transparent"
	ColorNames = append(ColorNames, colorName)
	pal := NewUniformPalette(colorName, color.Transparent)
	ColorMap[colorName] = pal
}

func init() {
	ReadJsonPalette("palettes.json")
	ReadNamedColors()
}
