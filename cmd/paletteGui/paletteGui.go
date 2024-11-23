package main

import (
	"github.com/stefan-muehlebach/ledgrid"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	_ "fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	Margin    = 10.0
	AppWidth  = 512.0
	AppHeight = 512.0
)

var (
	AppSize = fyne.NewSize(AppWidth, AppHeight)
)

var (
	App fyne.App
	Win fyne.Window
)

func main() {
	var paletteNameList []string

	var pal ledgrid.ColorSource = ledgrid.PaletteMap["HotAndCold"]
	paletteNameList = ledgrid.PaletteNames
	// paletteNameList = make([]string, len(ledgrid.PaletteList))
	// for i, palette := range ledgrid.PaletteList {
	// 	paletteNameList[i] = palette.Name()
	// }

	//------------------------------------------------------------------------
	//
	// Ab dieser Stelle wird das GUI aufgebaut
	//

	App = app.New()
	Win = App.NewWindow("Palette GUI")

	form := container.New(
		layout.NewFormLayout(),
	)
	label := widget.NewLabel("Palette/Color")
	label.Alignment = fyne.TextAlignTrailing
	label.TextStyle.Bold = true
	visPal := ledgrid.NewPaletteWidget(pal)
	selection := widget.NewSelect(paletteNameList, func(s string) {
		pal := ledgrid.PaletteMap[s]
		visPal.ColorSource = pal
		visPal.Refresh()
	})
	selection.Selected = pal.Name()
	form.Add(label)
	form.Add(selection)
	// label = widget.NewLabel("")
	// form.Add(label)
	// form.Add(visPal)

	card := widget.NewCard("Palette", "", form)

	quitBtn := widget.NewButton("Quit", App.Quit)
	btnBox := container.NewHBox(layout.NewSpacer(), quitBtn)

	root := container.NewVBox(
		card,
		visPal,
		layout.NewSpacer(),
		btnBox,
	)

	Win.Canvas().SetOnTypedKey(func(evt *fyne.KeyEvent) {
		switch evt.Name {
		case fyne.KeyEscape, fyne.KeyQ:
			App.Quit()
		}
	})

	Win.SetContent(root)
	Win.Resize(AppSize)
	Win.ShowAndRun()
}
