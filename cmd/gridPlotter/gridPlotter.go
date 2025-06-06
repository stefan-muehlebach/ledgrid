package main

import (
	"strings"
	"fmt"
	"image"
	"flag"
	"github.com/stefan-muehlebach/ledgrid/conf"
)

func main() {
    var width, height int
    var customConfName string
	var modConf conf.ModuleConfig
	var outFileName string
    var showList bool

    flag.IntVar(&width, "width", 0, "Width of panel")
    flag.IntVar(&height, "height", 0, "Height of panel")
    flag.StringVar(&customConfName, "custom", "", "Use a non standard module configuration")
    flag.BoolVar(&showList, "list", false, "list all custom configuration files")
    flag.Parse()

    if showList {
        for _, name := range conf.AllCustomFiles() {
            name = strings.TrimSuffix(name, ".json")
            fmt.Printf("%s\n", name)
        }
    } else {
        if width > 0 && height > 0 {
		    gridSize := image.Point{width, height}
            outFileName = fmt.Sprintf("default%dx%d.png", gridSize.X, gridSize.Y)
		    modConf = conf.DefaultModuleConfig(gridSize)
	    } else if customConfName != "" {
            fileName := "data/" + customConfName + ".json"
            outFileName = customConfName + ".png"
	        modConf = conf.Load(fileName)
        } else {
            fmt.Printf("either width/height or custom must be specified!")
            return
        }
        modConf.Plot(outFileName)
    }
}

