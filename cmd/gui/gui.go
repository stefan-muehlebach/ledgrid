//go:generate fyne bundle -o data.go Icon.ico

package main

import (
	"flag"
	"image"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/stefan-muehlebach/ledgrid"
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
	width              = 40
	height             = 10
	gridSize           = image.Point{width, height}
	defLocal           = false
	defDummy           = false
	defHost            = "raspi-3"
	defPort       uint = 5333
	defGammaValue      = 3.0
	blinkenFiles       = []string{
		"bml/flatter.bml",
		"bml/torus.bml",
		"bml/cube.bml",
		"bml/kreise.bml",
		"bml/benedictus.bml",
		"bml/lemming.bml",
		"bml/marioWalkRight.bml",
		"bml/marioRunRight.bml",
	}
	gradientImageFile = "gradient.png"

	App fyne.App
	Win fyne.Window
)

func Quit() {
	dialog.ShowConfirm("Quit", "Wollen sie die Applikation beenden?", func(b bool) {
		if b {
			App.Quit()
		}
	}, Win)
}

func main() {
	var local, dummy bool
	var host string
	var port uint
	var gammaValue, maxBrightValue, fadeTime ledgrid.FloatParameter

	var pixCtrl ledgrid.PixelClient
	var pixGrid *ledgrid.LedGrid
	var pixAnim *ledgrid.Animator

	var bgList, fgList []ledgrid.Visual
	var bgNameList, fgNameList []string
	// var paletteNameList, colorNameList []string

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
	pixGrid = ledgrid.NewLedGrid(gridSize, nil)
	pixAnim = ledgrid.NewAnimator(pixGrid, pixCtrl)

	gammaValue = ledgrid.NewFloatParameter("Gamma", defGammaValue, 1.0, 5.0, 0.1)
	gammaValue.SetCallback(func(p ledgrid.Parameter) {
		v := gammaValue.Val()
		pixCtrl.SetGamma(v, v, v)
	})
	maxBrightValue = ledgrid.NewFloatParameter("Brightness", 255, 1, 255, 1)
	maxBrightValue.SetCallback(func(p ledgrid.Parameter) {
		v := uint8(maxBrightValue.Val())
		pixCtrl.SetMaxBright(v, v, v)
	})

	fadeTime = ledgrid.NewFloatParameter("Fade Time", 2.0, 0.0, 5.0, 0.1)

	transpVisual := ledgrid.NewUniform(pixGrid, ledgrid.ColorMap["Transparent"])
	transpVisual.SetName("Uniform Color")

	bgList = []ledgrid.Visual{
		transpVisual,
		ledgrid.NewShader(pixGrid, ledgrid.ExperimentalShader, ledgrid.PaletteMap["Hipster"]),
		ledgrid.NewShader(pixGrid, ledgrid.PlasmaShader, ledgrid.PaletteMap["Nightspell"]),
		ledgrid.NewShader(pixGrid, ledgrid.CircleShader, ledgrid.PaletteMap["Hipster"]),
		ledgrid.NewShader(pixGrid, ledgrid.KaroShader, ledgrid.PaletteMap["Hipster"]),
		ledgrid.NewShader(pixGrid, ledgrid.LinearShader, ledgrid.PaletteMap["Hipster"]),
		ledgrid.NewFire(pixGrid),
		ledgrid.NewCamera(pixGrid),
	}
	bgNameList = make([]string, len(bgList))
	for i, anim := range bgList {
		bgNameList[i] = anim.Name()
	}

	fgList = []ledgrid.Visual{
		transpVisual,
		//ledgrid.NewTextNative(pixGrid, "Beni und Stefan haben Ferien", ledgrid.ColorMap["GreenYellow"]),
		ledgrid.NewTextFreeType(pixGrid, "Benedict", ledgrid.ColorMap["SkyBlue"]),
		ledgrid.NewImageFromFile(pixGrid, "image.png"),
		ledgrid.NewImageFromFile(pixGrid, "gradient.png"),
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

	// paletteNameList = ledgrid.PaletteNames
	// paletteNameList = make([]string, len(ledgrid.PaletteList))
	// for i, palette := range ledgrid.PaletteList {
	// 	paletteNameList[i] = palette.Name()
	// }

	// colorNameList = ledgrid.ColorNames
	// colorNameList = make([]string, len(ledgrid.ColorList))
	// for i, palette := range ledgrid.ColorList {
	// 	colorNameList[i] = palette.Name()
	// }

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
			var selection *widget.Select
			var pal ledgrid.ColorSource

			palParam := obj.PaletteParam()
			label := widget.NewLabel(palParam.Name())
			label.Alignment = fyne.TextAlignTrailing
			label.TextStyle.Bold = true
			if p, ok := palParam.Val().(*ledgrid.PaletteFader); ok {
				pal = p.Pals[0]
			}
			switch pal.(type) {
			case *ledgrid.UniformPalette:
				selection = widget.NewSelect(ledgrid.ColorNames, func(s string) {
					pal := ledgrid.ColorMap[s]
					obj.SetPalette(pal, time.Duration(fadeTime.Val()*float64(time.Second)))
				})
			default:
				selection = widget.NewSelect(ledgrid.PaletteNames, func(s string) {
					pal := ledgrid.PaletteMap[s]
					obj.SetPalette(pal, time.Duration(fadeTime.Val()*float64(time.Second)))
				})
			}
			selection.Selected = palParam.Val().Name()
			form.Add(label)
			form.Add(selection)
		}
		if obj, ok := vis.(ledgrid.Parametrizable); ok {
			for _, p := range obj.ParamList() {
				switch param := p.(type) {
				case ledgrid.FloatParameter:
					label := widget.NewLabelWithData(binding.FloatToStringWithFormat(param, param.Name()+" (%.1f)"))
					label.Alignment = fyne.TextAlignTrailing
					label.TextStyle.Bold = true
					slider := widget.NewSliderWithData(param.Min(), param.Max(), param)
					slider.Step = param.Step()
					slider.SetValue(param.Val())
					form.Add(label)
					form.Add(slider)

				case ledgrid.StringParameter:
					label := widget.NewLabel(param.Name())
					label.Alignment = fyne.TextAlignTrailing
					label.TextStyle.Bold = true
					entry := widget.NewEntry()
					entry.Text = param.Val()
					button := widget.NewButton("Apply", func() {
						param.SetVal(entry.Text)
					})
					form.Add(label)
					form.Add(entry)
					form.Add(layout.NewSpacer())
					form.Add(button)
				}
			}
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

	gammaLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(gammaValue, gammaValue.Name()+" (%.1f)"))
	gammaLabel.Alignment = fyne.TextAlignTrailing
	gammaLabel.TextStyle.Bold = true
	gammaSlider := widget.NewSliderWithData(gammaValue.Min(), gammaValue.Max(), gammaValue)
	gammaSlider.Step = gammaValue.Step()
	gammaSlider.SetValue(gammaValue.Val())

	maxBrightLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(maxBrightValue, maxBrightValue.Name()+" (%.0f)"))
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
	prefTab := container.NewVBox(
		prefCard,
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

	quitBtn := widget.NewButton("Quit", Quit)
	btnBox := container.NewHBox(layout.NewSpacer(), quitBtn)

	root := container.NewVBox(
		tabs,
		layout.NewSpacer(),
		btnBox,
	)

	Win.Canvas().SetOnTypedKey(func(evt *fyne.KeyEvent) {
		switch evt.Name {
		case fyne.KeyEscape, fyne.KeyQ:
			Quit()
		}
	})

	Win.SetContent(root)
	Win.Resize(AppSize)
	Win.ShowAndRun()

	pixAnim.Stop()
	pixGrid.Clear(ledgrid.Black)
	pixCtrl.Draw(pixGrid)
	pixCtrl.Close()
}
