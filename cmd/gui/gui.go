//go:generate fyne bundle -o data.go Icon.ico

package main

import (
	"flag"
	"image"
	"time"

	"github.com/stefan-muehlebach/ledgrid"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	_ "fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	Margin    = 10.0
	AppWidth  = 512.0
	AppHeight = 1024.0
)

var (
	AppSize = fyne.NewSize(AppWidth, AppHeight)
)

var (
	width              = 10
	height             = 10
	defLocal           = false
	defDummy           = false
	defHost            = "raspi-2"
	defPort       uint = 5333
	defGammaValue      = 3.0
	blinkenFiles       = []string{
		"bml/flatter.bml",
		"bml/torus.bml",
		"bml/cube.bml",
		"bml/kreise.bml",
		"bml/benedictus.bml",
		"bml/lemming.bml",
		"bml/mario.bml",
	}

    App fyne.App
    Win fyne.Window
)

func main() {
	var local, dummy bool
	var host string
	var port uint
	var gammaValue, maxBrightValue *ledgrid.Bounded[float64]

	var pixCtrl ledgrid.PixelClient
	var pixGrid *ledgrid.LedGrid
	var pixAnim *ledgrid.Animator

	var bgList, fgList []ledgrid.Visual
	var bgNameList, fgNameList []string
	var paletteNameList []string

    var pal ledgrid.ColorSource
	// var pal *ledgrid.PaletteFader
	var fadeTime *ledgrid.Bounded[float64]

	var blinken *ledgrid.BlinkenFile
	var blinkenAnim *ledgrid.ImageAnimation

	var bgTypeSelect, fgTypeSelect *widget.Select

	var bgParamForm, fgParamForm *fyne.Container

	flag.BoolVar(&local, "local", defLocal, "PixelController is local")
	flag.BoolVar(&dummy, "dummy", defDummy, "Use dummy PixelController")
	flag.StringVar(&host, "host", defHost, "Controller hostname")
	flag.UintVar(&port, "port", defPort, "Controller port")
	flag.Parse()

	if dummy {
		pixCtrl = ledgrid.NewDummyPixelClient()
	} else {
		if local {
			pixCtrl = ledgrid.NewLocalPixelClient(5333, "/dev/spidev0.0", 2_000_000)
		} else {
			pixCtrl = ledgrid.NewNetPixelClient(host, port)
		}
	}
	pixGrid = ledgrid.NewLedGrid(image.Rect(0, 0, width, height))
	pixAnim = ledgrid.NewAnimator(pixGrid, pixCtrl)

	gammaValue = ledgrid.NewBounded("Gamma", defGammaValue, 1.0, 5.0, 0.1)
	gammaValue.SetCallback(func(oldVal, newVal float64) {
		pixCtrl.SetGamma(newVal, newVal, newVal)
	})
	maxBrightValue = ledgrid.NewBounded("MaxBright", 255.0, 1.0, 255.0, 1.0)
	maxBrightValue.SetCallback(func(oldVal, newVal float64) {
		val := uint8(newVal)
		pixCtrl.SetMaxBright(val, val, val)
	})

	// pal = ledgrid.NewPaletteFader(ledgrid.HipsterPalette)
    pal = ledgrid.HipsterPalette
	fadeTime = ledgrid.NewBounded("Fade Time", 2.0, 0.0, 5.0, 0.1)

	transpVisual := ledgrid.NewImageFromColor(pixGrid, ledgrid.Transparent)
	transpVisual.SetName("(Transparent)")

	bgList = []ledgrid.Visual{
		transpVisual,
		ledgrid.NewShader(pixGrid, ledgrid.PlasmaShader, pal),
		ledgrid.NewShader(pixGrid, ledgrid.CircleShader, pal),
		ledgrid.NewShader(pixGrid, ledgrid.KaroShader, pal),
		ledgrid.NewShader(pixGrid, ledgrid.LinearShader, pal),
		ledgrid.NewFire(pixGrid),
		ledgrid.NewCamera(pixGrid),
	}
	bgNameList = make([]string, len(bgList))
	for i, anim := range bgList {
		bgNameList[i] = anim.Name()
	}

	fgList = []ledgrid.Visual{
		transpVisual,
		ledgrid.NewTextNative(pixGrid, "Benedict", ledgrid.PaletteMap["GreenYellowColor"]),
		ledgrid.NewTextFreeType(pixGrid, "Benedict", ledgrid.PaletteMap["LightSeaGreenColor"]),
		ledgrid.NewImageFromFile(pixGrid, "image.png"),
	}
	for _, fileName := range blinkenFiles {
		blinken = ledgrid.ReadBlinkenFile(fileName)
		blinkenAnim = blinken.NewImageAnimation(pixGrid)
		fgList = append(fgList, blinkenAnim)
	}
	fgNameList = make([]string, len(fgList))
	for i, anim := range fgList {
		fgNameList[i] = anim.Name()
	}

	paletteNameList = make([]string, len(ledgrid.PaletteList))
	for i, palette := range ledgrid.PaletteList {
		paletteNameList[i] = palette.Name()
	}

	//------------------------------------------------------------------------
	//
	// Ab dieser Stelle wird das GUI aufgebaut
	//
	ShowParameter := func(vis ledgrid.Visual, form *fyne.Container) {
		for _, obj := range form.Objects {
			switch o := obj.(type) {
			case *widget.Label:
				o.Unbind()
			case *widget.Slider:
				o.Unbind()
			}
		}
		form.RemoveAll()
		if obj, ok := vis.(ledgrid.Paintable); ok {
			label := widget.NewLabel("Palette/Color")
			label.Alignment = fyne.TextAlignTrailing
			label.TextStyle.Bold = true
			visPal := ledgrid.NewPalette(obj.Palette())
			selection := widget.NewSelect(paletteNameList, func(s string) {
				pal := ledgrid.PaletteMap[s]
				obj.SetPalette(pal, time.Duration(fadeTime.Val()*float64(time.Second)))
				visPal.ColorSource = pal
				visPal.Refresh()
			})
			selection.Selected = obj.Palette().Name()
			form.Add(label)
			form.Add(selection)
			label = widget.NewLabel("")
			form.Add(label)
			form.Add(visPal)
		}
		if obj, ok := vis.(ledgrid.Parametrizable); ok {
			for _, param := range obj.ParamList() {
				label := widget.NewLabelWithData(binding.FloatToStringWithFormat(param, param.Name()+" (%.3f)"))
				label.Alignment = fyne.TextAlignTrailing
				label.TextStyle.Bold = true
				slider := widget.NewSliderWithData(param.Min(), param.Max(), param)
				slider.Step = param.Step()
				slider.SetValue(param.Val())
				form.Add(label)
				form.Add(slider)
			}
		}
		if obj, ok := vis.(ledgrid.Text); ok {
			label := widget.NewLabel("Message")
			label.Alignment = fyne.TextAlignTrailing
			label.TextStyle.Bold = true
			entry := widget.NewEntry()
			entry.Text = obj.String()
			button := widget.NewButton("Apply", func() {
				obj.SetString(entry.Text)
			})
			form.Add(label)
			form.Add(entry)
			form.Add(layout.NewSpacer())
			form.Add(button)
		}
	}

	App = app.New()
	App.SetIcon(resourceIconIco)
	Win = App.NewWindow("LedGrid GUI")

	bgTypeSelect = widget.NewSelect(bgNameList, func(s string) {
		newBg := bgList[bgTypeSelect.SelectedIndex()]
		pixAnim.SetBackground(newBg, time.Duration(fadeTime.Val()*float64(time.Second)))
		ShowParameter(newBg, bgParamForm)
	})

	fgTypeSelect = widget.NewSelect(fgNameList, func(s string) {
		newFg := fgList[fgTypeSelect.SelectedIndex()]
		pixAnim.SetForeground(newFg, time.Duration(fadeTime.Val()*float64(time.Second)))
		ShowParameter(newFg, fgParamForm)
	})

	visualForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Background", Widget: bgTypeSelect},
			{Text: "Foreground", Widget: fgTypeSelect},
		},
	}

	gammaLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(gammaValue, "Gamma (%.1f)"))
	gammaLabel.Alignment = fyne.TextAlignTrailing
	gammaLabel.TextStyle.Bold = true
	gammaSlider := widget.NewSliderWithData(gammaValue.Min(), gammaValue.Max(), gammaValue)
	gammaSlider.Step = gammaValue.Step()
	gammaSlider.SetValue(gammaValue.Val())

	maxBrightLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(maxBrightValue, "Brightness (%.0f)"))
	maxBrightLabel.Alignment = fyne.TextAlignTrailing
	maxBrightLabel.TextStyle.Bold = true
	maxBrightSlider := widget.NewSliderWithData(maxBrightValue.Min(), maxBrightValue.Max(), maxBrightValue)
	maxBrightSlider.Step = maxBrightValue.Step()
	maxBrightSlider.SetValue(maxBrightValue.Val())

	fadeTimeLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(fadeTime, fadeTime.Name()+" (%.1f)"))
	fadeTimeLabel.Alignment = fyne.TextAlignTrailing
	fadeTimeLabel.TextStyle.Bold = true
	fadeTimeSlider := widget.NewSliderWithData(fadeTime.Min(), fadeTime.Max(), fadeTime)
	fadeTimeSlider.Step = fadeTime.Step()
	fadeTimeSlider.SetValue(fadeTime.Val())

	prefForm := container.New(
		layout.NewFormLayout(),
		gammaLabel, gammaSlider,
		maxBrightLabel, maxBrightSlider,
		fadeTimeLabel, fadeTimeSlider,
	)
	prefCard := widget.NewCard("Preferences", "", prefForm)
	// testWidget := NewPalettePreview(ledgrid.NightspellPalette)
	prefTab := container.NewVBox(
		prefCard,
		// testWidget,
	)

	visualCard := widget.NewCard("Visuals", "There can be one background and one foreground", visualForm)

	bgParamForm = container.New(
		layout.NewFormLayout(),
	)
	bgParamCard := widget.NewCard("Background Parameters", "", bgParamForm)
	fgParamForm = container.New(
		layout.NewFormLayout(),
	)
	fgParamCard := widget.NewCard("Foreground Parameters", "", fgParamForm)

	effectTab := container.NewVBox(
		visualCard,
		bgParamCard,
		fgParamCard,
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Visuals", effectTab),
		container.NewTabItem("Preferences", prefTab),
	)

	bgTypeSelect.SetSelectedIndex(0)
	fgTypeSelect.SetSelectedIndex(0)

	quitBtn := widget.NewButton("Quit", App.Quit)
	btnBox := container.NewHBox(layout.NewSpacer(), quitBtn)

	root := container.NewVBox(
		tabs,
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

	pixGrid.Clear(ledgrid.Black)
	pixCtrl.Draw(pixGrid)
	pixCtrl.Close()

}
