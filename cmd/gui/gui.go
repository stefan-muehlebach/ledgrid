//go:generate fyne bundle -o data.go Icon.ico

package main

import (
	"flag"
	"image"
	"time"

	"fyne.io/fyne/v2/data/binding"

	"fyne.io/fyne/v2/layout"

	"github.com/stefan-muehlebach/gg/colornames"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	_ "fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/stefan-muehlebach/ledgrid"
)

const (
	Margin    = 10.0
	AppWidth  = 480.0
	AppHeight = 480.0
)

var (
	AppSize = fyne.NewSize(AppWidth, AppHeight)
)

var (
	width              = 10
	height             = 10
	defLocal           = false
	defHost            = "raspi-2"
	defPort       uint = 5333
	defGammaValue      = 3.0
)

func main() {
	var local bool
	var host string
	var port uint
	var gammaValue *ledgrid.Bounded[float64]

	var pixCtrl ledgrid.PixelClient
	var pixGrid *ledgrid.LedGrid
	var pixAnim *ledgrid.Animator

	var animList []ledgrid.Visual
	var animNameList, paletteNameList []string
	// var curAnim ledgrid.Visual

	var pal *ledgrid.PaletteFader
	var palFadeTime, backFadeTime *ledgrid.Bounded[float64]

	var blinken *ledgrid.BlinkenFile
	var flatterAnim, torusAnim, lemmingAnim, marioAnim *ledgrid.ImageAnimation

	var animSelect, paletteSelect *widget.Select
	var animLabel, paletteLabel *widget.Label

	var paramForm *fyne.Container

	flag.BoolVar(&local, "local", defLocal, "PixelController is local")
	flag.StringVar(&host, "host", defHost, "Controller hostname")
	flag.UintVar(&port, "port", defPort, "Controller port")
	flag.Parse()

	if local {
		pixCtrl = ledgrid.NewPixelServer(5333, "/dev/spidev0.0", 2_000_000)
	} else {
		pixCtrl = ledgrid.NewNetPixelClient(host, port)
	}
	pixGrid = ledgrid.NewLedGrid(image.Rect(0, 0, width, height))
	pixAnim = ledgrid.NewAnimator(pixGrid, pixCtrl)

	gammaValue = ledgrid.NewBounded("Gamma", defGammaValue, 1.0, 5.0, 0.1)
	gammaValue.SetCallback(func(oldVal, newVal float64) {
		pixCtrl.SetGamma(newVal, newVal, newVal)
	})

	pal = ledgrid.NewPaletteFader(ledgrid.HipsterPalette)
	palFadeTime = ledgrid.NewBounded("Fade Time", 1.5, 0.0, 5.0, 0.1)

	backFadeTime = ledgrid.NewBounded("Fade Time", 1.5, 0.0, 5.0, 0.1)

	blinken = ledgrid.ReadBlinkenFile("flatter.bml")
	flatterAnim = blinken.NewImageAnimation(pixGrid)

	blinken = ledgrid.ReadBlinkenFile("torus.bml")
	torusAnim = blinken.NewImageAnimation(pixGrid)

	blinken = ledgrid.ReadBlinkenFile("lemmingWalk.bml")
	lemmingAnim = blinken.NewImageAnimation(pixGrid)

	blinken = ledgrid.ReadBlinkenFile("mario.bml")
	marioAnim = blinken.NewImageAnimation(pixGrid)

	animList = []ledgrid.Visual{
		ledgrid.NewShader(pixGrid, ledgrid.PlasmaShader, pal),
		ledgrid.NewShader(pixGrid, ledgrid.CircleShader, pal),
		ledgrid.NewShader(pixGrid, ledgrid.KaroShader, pal),
		ledgrid.NewShader(pixGrid, ledgrid.LinearShader, pal),
		ledgrid.NewFire(pixGrid),
		ledgrid.NewCamera(pixGrid),
		ledgrid.NewText(pixGrid, "Lochbach", colornames.Crimson),
		ledgrid.NewImageFromFile(pixGrid, "image.png"),
		flatterAnim,
		torusAnim,
		lemmingAnim,
		marioAnim,
		// txtAnim,dir
	}
	animNameList = make([]string, len(animList))
	for i, anim := range animList {
		animNameList[i] = anim.Name()
		//pixAnim.AddObjects(anim)
	}

	paletteNameList = make([]string, len(ledgrid.PaletteList))
	for i, palette := range ledgrid.PaletteList {
		paletteNameList[i] = palette.Name()
	}

	// Ab dieser Stelle wird das GUI aufgebaut
	myApp := app.New()
	myApp.SetIcon(resourceIconIco)
	myWindow := myApp.NewWindow("PixelGui")

	animLabel = widget.NewLabel("Type:")
	animSelect = widget.NewSelect(animNameList, func(s string) {
        newBack := animList[animSelect.SelectedIndex()]
        pixAnim.SetBackground(newBack, time.Duration(backFadeTime.Val() * float64(time.Second)))
        if obj, ok := newBack.(ledgrid.Paintable); ok {
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
		if obj, ok := newBack.(ledgrid.Parametrizable); ok {
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

	backFadeTimeLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(backFadeTime, backFadeTime.Name()+": %.1f"))
	backFadeTimeSlider := widget.NewSliderWithData(backFadeTime.Min(), backFadeTime.Max(), backFadeTime)
	backFadeTimeSlider.Step = backFadeTime.Step()
	backFadeTimeSlider.SetValue(backFadeTime.Val())

	paletteLabel = widget.NewLabel("Color:")
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

	backForm := container.New(
		layout.NewFormLayout(),
		animLabel, animSelect,
		backFadeTimeLabel, backFadeTimeSlider,
		paletteLabel, paletteSelect,
		palFadeTimeLabel, palFadeTimeSlider,
		gammaLabel, gammaSlider,
	)
	backCard := widget.NewCard("Background", "", backForm)

	paramForm = container.New(
		layout.NewFormLayout(),
	)
	paramCard := widget.NewCard("Parameters", "", paramForm)

	quitBtn := widget.NewButton("Quit", myApp.Quit)
	btnBox := container.NewHBox(layout.NewSpacer(), quitBtn)

	root := container.NewVBox(
		backCard,
		widget.NewSeparator(),
		paramCard,
		layout.NewSpacer(),
		btnBox,
	)

	myWindow.Canvas().SetOnTypedKey(func(evt *fyne.KeyEvent) {
		switch evt.Name {
		case fyne.KeyEscape, fyne.KeyQ:
			myApp.Quit()
		}
	})

	myWindow.SetContent(root)
	myWindow.Resize(AppSize)
	myWindow.ShowAndRun()

	pixGrid.Clear(ledgrid.Black)
	pixCtrl.Draw(pixGrid)
	pixCtrl.Close()

}
