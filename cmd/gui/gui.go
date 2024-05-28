//go:generate fyne bundle -o data.go Icon.ico

package main

import (
	"flag"
	"image"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/stefan-muehlebach/gg/colornames"
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
	AppWidth  = 480.0
	AppHeight = 640.0
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
	blinkenFiles       = []string{"flatter.bml", "torus.bml", "lemming.bml", "mario.bml"}
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

	var pal *ledgrid.PaletteFader
	var palFadeTime, bgFadeTime, fgFadeTime *ledgrid.Bounded[float64]

	var blinken *ledgrid.BlinkenFile
	var blinkenAnim *ledgrid.ImageAnimation

	var bgTypeLabel, fgTypeLabel, paletteLabel *widget.Label
	var bgTypeSelect, fgTypeSelect, paletteSelect *widget.Select

	var paramForm *fyne.Container

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

	pal = ledgrid.NewPaletteFader(ledgrid.HipsterPalette)
	palFadeTime = ledgrid.NewBounded("Fade Time", 1.5, 0.0, 5.0, 0.1)

	bgFadeTime = ledgrid.NewBounded("Fade Time", 2.0, 0.0, 5.0, 0.1)
	fgFadeTime = ledgrid.NewBounded("Fade Time", 2.0, 0.0, 5.0, 0.1)

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
		ledgrid.NewText(pixGrid, "Lochbach", colornames.Crimson),
		ledgrid.NewTextFT(pixGrid, "Lochbach", colornames.Teal),
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
	app := app.New()
	app.SetIcon(resourceIconIco)
	win := app.NewWindow("LedGrid GUI")

	bgTypeLabel = widget.NewLabel("Background:")
	bgTypeSelect = widget.NewSelect(bgNameList, func(s string) {
		newBg := bgList[bgTypeSelect.SelectedIndex()]
		pixAnim.SetBackground(newBg, time.Duration(bgFadeTime.Val()*float64(time.Second)))
		if obj, ok := newBg.(ledgrid.Paintable); ok {
			paletteSelect.Enable()
			paletteSelect.SetSelected(obj.Palette().Name())
		} else {
			paletteSelect.SetSelectedIndex(-1)
			paletteSelect.Disable()
		}
		for _, obj := range paramForm.Objects {
			switch o := obj.(type) {
			case *widget.Label:
				o.Unbind()
			case *widget.Slider:
				o.Unbind()
			}
		}
		paramForm.RemoveAll()
		if obj, ok := newBg.(ledgrid.Parametrizable); ok {
			for _, param := range obj.ParamList() {
				label := widget.NewLabelWithData(binding.FloatToStringWithFormat(param, param.Name()+": %.3f"))
				slider := widget.NewSliderWithData(param.Min(), param.Max(), param)
				slider.Step = param.Step()
				slider.SetValue(param.Val())
				paramForm.Add(label)
				paramForm.Add(slider)
			}
		}
	})

	bgFadeTimeLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(bgFadeTime, bgFadeTime.Name()+": %.1f"))
	bgFadeTimeSlider := widget.NewSliderWithData(bgFadeTime.Min(), bgFadeTime.Max(), bgFadeTime)
	bgFadeTimeSlider.Step = bgFadeTime.Step()
	bgFadeTimeSlider.SetValue(bgFadeTime.Val())

	fgTypeLabel = widget.NewLabel("Foreground:")
	fgTypeSelect = widget.NewSelect(fgNameList, func(s string) {
		newFg := fgList[fgTypeSelect.SelectedIndex()]
		pixAnim.SetForeground(newFg, time.Duration(fgFadeTime.Val()*float64(time.Second)))
	})
	fgFadeTimeLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(fgFadeTime, fgFadeTime.Name()+": %.1f"))
	fgFadeTimeSlider := widget.NewSliderWithData(fgFadeTime.Min(), fgFadeTime.Max(), fgFadeTime)
	fgFadeTimeSlider.Step = fgFadeTime.Step()
	fgFadeTimeSlider.SetValue(fgFadeTime.Val())

	paletteLabel = widget.NewLabel("Palette:")
	paletteSelect = widget.NewSelect(paletteNameList, func(s string) {
		id := paletteSelect.SelectedIndex()
		pal.StartFade(ledgrid.PaletteList[id], time.Duration(palFadeTime.Val()*float64(time.Second)))
	})
	paletteSelect.Disable()

	palFadeTimeLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(palFadeTime, palFadeTime.Name()+": %.1f"))
	palFadeTimeSlider := widget.NewSliderWithData(palFadeTime.Min(), palFadeTime.Max(), palFadeTime)
	palFadeTimeSlider.Step = palFadeTime.Step()
	palFadeTimeSlider.SetValue(palFadeTime.Val())

	gammaLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(gammaValue, "Gamma: %.1f"))
	gammaSlider := widget.NewSliderWithData(gammaValue.Min(), gammaValue.Max(), gammaValue)
	gammaSlider.Step = gammaValue.Step()
	gammaSlider.SetValue(gammaValue.Val())

	maxBrightLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(maxBrightValue, "Max. Bright: %.0f"))
	maxBrightSlider := widget.NewSliderWithData(maxBrightValue.Min(), maxBrightValue.Max(), maxBrightValue)
	maxBrightSlider.Step = maxBrightValue.Step()
	maxBrightSlider.SetValue(maxBrightValue.Val())

	visualForm := container.New(
		layout.NewFormLayout(),
		bgTypeLabel, bgTypeSelect,
		bgFadeTimeLabel, bgFadeTimeSlider,
		fgTypeLabel, fgTypeSelect,
		fgFadeTimeLabel, fgFadeTimeSlider,
	)
	visualCard := widget.NewCard("Visual Effects", "", visualForm)

	colorForm := container.New(
		layout.NewFormLayout(),
		paletteLabel, paletteSelect,
		palFadeTimeLabel, palFadeTimeSlider,
	)
	colorCard := widget.NewCard("Color / Palette", "", colorForm)

	paramForm = container.New(
		layout.NewFormLayout(),
	)
	paramCard := widget.NewCard("Parameters", "", paramForm)

	effectTab := container.NewVBox(
		visualCard,
		colorCard,
        paramCard,
	)

	prefTab := container.New(
        layout.NewFormLayout(),
    		gammaLabel, gammaSlider,
        maxBrightLabel, maxBrightSlider,
    )

	tabs := container.NewAppTabs(
		container.NewTabItem("Effects", effectTab),
		container.NewTabItem("Preferences", prefTab),
	)

	saveBtn := widget.NewButton("Save", func() {
		fh, err := os.Create("ledgrid.png")
		if err != nil {
			log.Fatalf("Couldn't create file 'ledgrid.png': %v", err)
		}
		png.Encode(fh, pixGrid)
		fh.Close()
	})
	quitBtn := widget.NewButton("Quit", app.Quit)
	btnBox := container.NewHBox(saveBtn, layout.NewSpacer(), quitBtn)

    root := container.NewVBox(
        tabs,
        layout.NewSpacer(),
        btnBox,
    )

	win.Canvas().SetOnTypedKey(func(evt *fyne.KeyEvent) {
		switch evt.Name {
		case fyne.KeyEscape, fyne.KeyQ:
			app.Quit()
		}
	})

	win.SetContent(root)
	win.Resize(AppSize)
	win.ShowAndRun()

	pixGrid.Clear(ledgrid.Black)
	pixCtrl.Draw(pixGrid)
	pixCtrl.Close()

}
