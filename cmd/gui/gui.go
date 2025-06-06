//go:generate fyne bundle -o data.go Icon.ico

package main

import (
	"context"
	"flag"
	"image"
	"math/rand"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/colors"
	"github.com/stefan-muehlebach/ledgrid/conf"
	"golang.org/x/image/math/fixed"
)

const (
	Margin         = 10.0
	AppWidth       = 512.0
	AppHeight      = 1024.0
	defHost        = "raspi-3"
	defWidth       = 40
	defHeight      = 10
	defPaletteName = "BackPinkBlue"
	defTextColor   = "White"
)

// Global variables, most of them related to the GUI and the LEDGrid.
var (
	App           fyne.App
	Win           fyne.Window
	AppSize       = fyne.NewSize(AppWidth, AppHeight)
	width, height int
	gridSize      image.Point
	gridClient    ledgrid.GridClient
	ledGrid       *ledgrid.LedGrid
	animCtrl      *ledgrid.AnimationController
	canvas        *ledgrid.Canvas
)

// Global variables associated with the (only) animation function [ColorWaves].
var (
	palFader = ledgrid.NewPaletteFader(ledgrid.PaletteMap[defPaletteName])
	txtColor = colors.Map[defTextColor]
	wordList []string
	wordIdx  int
	animTime = 3 * time.Second
)

// The animation function - so far the only one taken over from the [anim]
// CLI program.
func ColorWaves(ctx context.Context, c *ledgrid.Canvas) {
	aGrpLedColor := ledgrid.NewGroup()

	for y := range c.Rect.Dy() {
		// ty := float64(y) / float64(c.Rect.Dy()-1)
		for x := range c.Rect.Dx() {
			// tx := float64(x) / float64(c.Rect.Dx()-1)
			pt := image.Point{x, y}
			pix := ledgrid.NewPixel(pt, colors.Black)

			c.Add(pix)

			aColorPal := ledgrid.NewPaletteAnim(pix, palFader, animTime)
			aColorPal.AutoReverse = true
			aColorPal.RepeatCount = ledgrid.AnimationRepeatForever
			aColorPal.Curve = ledgrid.AnimationLinear
			aColorPal.Pos = rand.Float64()
			aGrpLedColor.Add(aColorPal)
		}
	}

	txt := ledgrid.NewFixedText(fixed.P(width/2, height/2), "", txtColor.Alpha(0))
	txt.SetAlign(ledgrid.AlignCenter | ledgrid.AlignMiddle)

	txtFadeIn := ledgrid.NewFadeAnim(txt, ledgrid.FadeIn, 500*time.Millisecond)
	txtColorOut := ledgrid.NewColorAnim(txt, colors.Black, 1000*time.Millisecond)
	txtFadeOut := ledgrid.NewFadeAnim(txt, ledgrid.FadeOut, 2000*time.Millisecond)

	txtNextWord := ledgrid.NewTask(func() {
		if len(wordList) == 0 {
			return
		}
		txt.SetText(wordList[wordIdx])
		wordIdx = (wordIdx + 1) % len(wordList)
		txt.Color = txtColor.Alpha(0)
	})
	txtSeq := ledgrid.NewSequence(txtNextWord, txtFadeIn,
		ledgrid.NewDelay(time.Second), txtColorOut, txtFadeOut)
	txtSeq.RepeatCount = ledgrid.AnimationRepeatForever
	c.Add(txt)

	txtSeq.Start()
	aGrpLedColor.Start()
}

// The function called when the user quits the program.
func Quit() {
	dialog.ShowConfirm("Quit", "Do you really want to quit the application?",
		func(b bool) {
			if b {
				App.Quit()
			}
		}, Win)
}

// ----------------------------------------------------------------------------
func main() {
	var host string
	var dataPort, rpcPort uint
	var useTCP bool
	var network string
	var gR, gG, gB float64
	var modConf conf.ModuleConfig

	flag.IntVar(&width, "width", defWidth, "Width (for 'out' option only)")
	flag.IntVar(&height, "height", defHeight, "Height (for 'out' option only)")
	flag.StringVar(&host, "host", defHost, "Controller hostname")
	flag.BoolVar(&useTCP, "tcp", false, "Use TCP for data")
	flag.UintVar(&dataPort, "data", ledgrid.DefDataPort, "Data Port")
	flag.UintVar(&rpcPort, "rpc", ledgrid.DefRPCPort, "RPC Port")
	flag.Parse()

	if useTCP {
		network = "tcp"
	} else {
		network = "udp"
	}
	gridClient = ledgrid.NewNetGridClient(host, network, dataPort, rpcPort)
	modConf = gridClient.ModuleConfig()
	ledGrid = ledgrid.NewLedGrid(gridClient, modConf)
	gR, gG, gB = gridClient.Gamma()

	gridSize = ledGrid.Rect.Size()
	width = gridSize.X
	height = gridSize.Y

	canvas = ledGrid.Canvas(0)
	animCtrl = ledGrid.AnimCtrl

	//------------------------------------------------------------------------
	//
	// Create the Fyne.io gui elements.
	//
	App = app.New()
	App.SetIcon(resourceIcon)
	Win = App.NewWindow("LedGrid GUI")

	//------------------------------------------------------------------------
	//
	// BEGIN of the GUI creation.
	//
	// This should be moved to a separate function in order to keep the
	// main function small and handy.
	//
	// Create the application form for displaying animations.
	//
	txtLabel := widget.NewLabel("Text")
	txtLabel.TextStyle.Bold = true
	txtEntry := widget.NewEntry()
	txtEntry.TextStyle.Monospace = true
	txtEntry.OnSubmitted = func(s string) {
		wordList = strings.Split(s, " ")
		wordIdx = 0
	}

	colorLabel := widget.NewLabel("Text Color")
	colorLabel.TextStyle.Bold = true
	colorSelect := widget.NewSelect(colors.Names, func(colorName string) {
		txtColor = colors.Map[colorName]
	})
	colorSelect.SetSelected(defTextColor)

	backLabel := widget.NewLabel("Palette")
	backLabel.TextStyle.Bold = true
	backSelect := widget.NewSelect(ledgrid.PaletteNames, func(backName string) {
		backAnim := ledgrid.NewPaletteFadeAnimation(palFader,
			ledgrid.PaletteMap[backName], animTime)
		backAnim.Start()
	})
	backSelect.SetSelected(defPaletteName)

    animLabel := widget.NewLabel("Animation")
    animLabel.TextStyle.Bold = true
    	animRadio := widget.NewRadioGroup([]string{"On", "Off"}, func(s string) {
        if s == "On" {
            ledgrid.AnimCtrl.Continue()
        } else if s == "Off" {
            	ledgrid.AnimCtrl.Suspend()
        }
    })
    animRadio.SetSelected("On")
    animRadio.Horizontal = true

	animForm := container.New(
		layout.NewFormLayout(),
		txtLabel, txtEntry,
		colorLabel, colorSelect,
		backLabel, backSelect,
        animLabel, animRadio,
	)
	animCard := widget.NewCard("Animations", "On this form, you specify which animation you want to run.", animForm)
	animTab := container.NewVBox(
		animCard,
	)

	// Create the application form for preferences.
	//
	gammaRed := binding.BindFloat(&gR)
	gammaRedLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(gammaRed,
		"Red (%.1f)"))
	// gammaRedLabel.Alignment = fyne.TextAlignTrailing
	gammaRedLabel.TextStyle.Bold = true
	gammaRedSlider := widget.NewSliderWithData(1.0, 3.0, gammaRed)
	gammaRedSlider.Step = 0.1
	gammaRedSlider.OnChangeEnded = func(v float64) {
		gR = v
		gridClient.SetGamma(gR, gG, gB)
	}

	gammaGreen := binding.BindFloat(&gG)
	gammaGreenLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(gammaGreen,
		"Green (%.1f)"))
	// gammaGreenLabel.Alignment = fyne.TextAlignTrailing
	gammaGreenLabel.TextStyle.Bold = true
	gammaGreenSlider := widget.NewSliderWithData(1.0, 3.0, gammaGreen)
	gammaGreenSlider.Step = 0.1
	gammaGreenSlider.OnChangeEnded = func(v float64) {
		gG = v
		gridClient.SetGamma(gR, gG, gB)
	}

	gammaBlue := binding.BindFloat(&gB)
	gammaBlueLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(gammaBlue,
		"Blue (%.1f)"))
	// gammaBlueLabel.Alignment = fyne.TextAlignTrailing
	gammaBlueLabel.TextStyle.Bold = true
	gammaBlueSlider := widget.NewSliderWithData(1.0, 3.0, gammaBlue)
	gammaBlueSlider.Step = 0.1
	gammaBlueSlider.OnChangeEnded = func(v float64) {
		gB = v
		gridClient.SetGamma(gR, gG, gB)
	}

	prefForm := container.New(
		layout.NewFormLayout(),
		gammaRedLabel, gammaRedSlider,
		gammaGreenLabel, gammaGreenSlider,
		gammaBlueLabel, gammaBlueSlider,
	)
	prefCard := widget.NewCard("Preferences", "Here you find all the settings for configuring the LEDGrid.", prefForm)
	prefTab := container.NewVBox(
		prefCard,
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Animations", animTab),
		container.NewTabItem("Preferences", prefTab),
	)

	quitBtn := widget.NewButton("Quit", Quit)
	btnBox := container.NewHBox(layout.NewSpacer(), quitBtn)

	root := container.NewVBox(
		tabs,
		layout.NewSpacer(),
		btnBox,
	)

	// END of the GUI creation.
	//------------------------------------------------------------------------

	Win.Canvas().SetOnTypedKey(func(evt *fyne.KeyEvent) {
		switch evt.Name {
		case fyne.KeyEscape, fyne.KeyQ:
			Quit()
		}
	})

	ledGrid.StartRefresh()
	ColorWaves(context.Background(), canvas)

	Win.SetContent(root)
	Win.Resize(AppSize)
	Win.ShowAndRun()

	ledgrid.AnimCtrl.Suspend()
	ledGrid.Clear(colors.Black)
	ledGrid.Show()
	ledGrid.Close()
}
