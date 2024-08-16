package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strconv"

	gc "github.com/rthornton128/goncurses"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/color"
)

const (
	defWidth      = 40
	defHeight     = 10
	defHost       = "raspi-3"
	defPort       = 5333
	KEY_SUP       = 0x151 /* Shifted up arrow */
	KEY_SDOWN     = 0x150 /* Shifted down arrow */
	KEY_CLEFT     = 0x222 /* Ctrl-left arrow */
	KEY_CRIGHT    = 0x231 /* Ctrl-right arrow */
	KEY_CUP       = 0x237 /* Ctrl-up arrow */
	KEY_CDOWN     = 0x20e /* Ctrl-down arrow */
	KEY_ALEFT     = 0x220 /* Alt-left arrow */
	KEY_ARIGHT    = 0x22f /* Alt-right arrow */
	KEY_AINS      = 0x21b /* Alt-Insert */
	KEY_ADEL      = 0x206 /* Alt-Delete */
	KEY_AHOME     = 0x216 /* Alt-Home */
	KEY_AEND      = 0x211 /* Alt-End */
	KEY_APAGEUP   = 0x22a /* Alt-PageUp */
	KEY_APAGEDOWN = 0x225 /* Alt-PageDown */
)

func Ctrl(x gc.Key) gc.Key {
	return x & 0x1f
}

func between(x, a, b int) bool {
	if a > b {
		a, b = b, a
	}
	return x >= a && x <= b
}

func main() {
	var stdscr *gc.Window
	var winGrid, winHelp *gc.Window
	var ch gc.Key
	var curRow, selRow, curCol, selCol int
	var curColor int
	var err error
	var ledGrid *ledgrid.LedGrid
	var ledColor color.LedColor
	var pixelClient ledgrid.PixelClient
	var curColorChanged bool
	var redrawGrid bool
	var cursorMoved bool
	var colors []uint8
	var incr uint8
	var newColorValue uint8
	var enterNewValue bool
	var newValueDigit int
	var gammaValues [3]float64
	var host string
	var port uint
	var width, height int
	var gridWidth, gridHeight int
	var termWidth, termHeight int
	var gridSize image.Point
	var selRect image.Rectangle
	var clipRect image.Rectangle
	var clipData []color.LedColor

	flag.IntVar(&width, "width", defWidth, "Width of panel")
	flag.IntVar(&height, "height", defHeight, "Height of panel")
	flag.StringVar(&host, "host", defHost, "Controller hostname")
	flag.UintVar(&port, "port", defPort, "Controller port")
	flag.Parse()

	gridSize = image.Point{width, height}
	gridWidth = 7*width + 10
	gridHeight = height + 7
	termWidth = gridWidth + 10
	termHeight = gridHeight + 40

	selRect = image.Rect(0, 0, 1, 1)

	ledGrid = ledgrid.NewLedGrid(gridSize, nil)
	pixelClient = ledgrid.NewNetPixelClient(host, port)
	gammaValues[0], gammaValues[1], gammaValues[2] = pixelClient.Gamma()

	stdscr, err = gc.Init()
	if err != nil {
		log.Fatalf("Couldn't Init ncurses: %v", err)
	}
	defer gc.End()

	err = gc.ResizeTerm(termHeight, termWidth)
	if err != nil {
		log.Fatalf("Couldn't resize terminal: %v", err)
	}

	gc.StartColor()
	gc.Echo(false)
	gc.CBreak(true)
	gc.Cursor(0)
	gc.Raw(true)

	stdscr.Keypad(true)

	gc.InitPair(1, gc.C_RED, gc.C_BLACK)
	gc.InitPair(2, gc.C_RED, gc.C_BLACK)
	gc.InitPair(3, gc.C_GREEN, gc.C_BLACK)
	gc.InitPair(4, gc.C_GREEN, gc.C_BLACK)
	gc.InitPair(5, gc.C_BLUE, gc.C_WHITE)
	gc.InitPair(6, gc.C_BLUE, gc.C_BLACK)

	y, x := 2, 4

	winGrid, err = gc.NewWindow(gridHeight, gridWidth, y, x)
	if err != nil {
		log.Fatalf("Couldn't create window: %v", err)
	}
	winGrid.Keypad(true)
	winGrid.Box(0, 0)

	helpHeight, helpWidth := 20, 60
	y, x = height+9, 4

	winHelp, err = gc.NewWindow(helpHeight, helpWidth, y, x)
	if err != nil {
		log.Fatalf("Couldn't create window: %v", err)
	}
	winHelp.Box(0, 0)

	winHelp.MovePrintf(1, 2, "|   R   |    G   |    B   |")
	winHelp.MovePrintf(2, 2, "+-------+--------+--------+")
	winHelp.MovePrintf(3, 2, "| [Ins] | [Home] | [PgUp] | increase color value")
	winHelp.MovePrintf(4, 2, "| [Del] | [End]  | [PgDn] | decrease color value")
	winHelp.MovePrintf(6, 2, "[Cursor]        : Move selector")
	winHelp.MovePrintf(7, 2, "[Alt]-Left/Right: Move color selector")
	winHelp.MovePrintf(8, 2, "[Shift]-[Cursor]: Select range")
	winHelp.MovePrintf(9, 2, "[Ctrl]-A        : Select all pixels")
	winHelp.MovePrintf(10, 2, "[Ctrl]-C/X/V    : Copy/Cut/Paste selected color values")
	winHelp.MovePrintf(11, 2, "[Backspace]     : Clear selection")
	winHelp.MovePrintf(12, 2, "F               : Interpolate colors")
	winHelp.MovePrintf(13, 2, "0-9a-f          : Enter new hex value for selected color")
	winHelp.MovePrintf(14, 2, "g/G             : Decrease/increase gamma values by 0.1")
	winHelp.MovePrintf(15, 2, "i               : Invert selected colors")
	winHelp.MovePrintf(16, 2, "+/-             : Darken/Brighten selected colors")
	winHelp.MovePrintf(18, 2, "q               : Quit")

	winGrid.Refresh()
	winHelp.Refresh()

	pixelClient.Send(ledGrid)

main:
	for {
		for col := 0; col < ledGrid.Rect.Dx(); col++ {
			if col == curCol {
				winGrid.AttrOn(gc.A_BOLD)
			}
			winGrid.MovePrintf(1, 10+col*7, "[%02x]", col)
			winGrid.AttrOff(gc.A_BOLD)
		}
		winGrid.MoveAddChar(2, 0, gc.ACS_LTEE)
		winGrid.HLine(2, 1, gc.ACS_HLINE, gridWidth-2)
		winGrid.MoveAddChar(2, gridWidth-1, gc.ACS_RTEE)
		for row := 0; row < ledGrid.Rect.Dy(); row++ {
			if row == curRow {
				winGrid.AttrOn(gc.A_BOLD)
			}
			winGrid.MovePrintf(3+row, 2, "[%02x]", row)
			winGrid.AttrOff(gc.A_BOLD)
		}

		row := height + 3
		winGrid.MoveAddChar(row, 0, gc.ACS_LTEE)
		winGrid.HLine(row, 1, gc.ACS_HLINE, gridWidth-2)
		winGrid.MoveAddChar(row, gridWidth-1, gc.ACS_RTEE)

		winGrid.VLine(1, 7, gc.ACS_VLINE, height+2)
		winGrid.MoveAddChar(0, 7, gc.ACS_TTEE)
		winGrid.MoveAddChar(height+3, 7, gc.ACS_BTEE)
		winGrid.MoveAddChar(2, 7, gc.ACS_PLUS)

		for row := 0; row < ledGrid.Rect.Dy(); row++ {
			for col := 0; col < ledGrid.Rect.Dx(); col++ {
				pt := image.Point{col, row}
				if pt.In(selRect) {
					if !enterNewValue {
						winGrid.AttrOn(gc.A_REVERSE)
					} else {
						winGrid.AttrOn(gc.A_BOLD)
					}
				}
				ledColor = ledGrid.LedColorAt(col, row)
				colors = []uint8{ledColor.R, ledColor.G, ledColor.B}
				for k := 0; k < 3; k++ {
					if (row == curRow) && (col == curCol) && (k == curColor) {
						if !enterNewValue {
							winGrid.ColorOn(int16(2*k + 1))
						} else {
							winGrid.ColorOn(int16(2*k + 2))
						}
					}
					winGrid.MovePrintf(3+row, 9+(col*7)+(k*2), "%02x", colors[k])
					winGrid.ColorOff(int16(2*k + 1))
					winGrid.ColorOff(int16(2*k + 2))
				}
				winGrid.AttrOff(gc.A_REVERSE)
				winGrid.AttrOff(gc.A_BOLD)
			}
		}

		winGrid.MovePrintf(row+1, 2, "New hex value for this color: ")
		if enterNewValue {
			winGrid.AttrOn(gc.A_REVERSE)
		}
		winGrid.Printf("%02x", newColorValue)
		winGrid.AttrOff(gc.A_REVERSE)
		winGrid.MovePrintf(row+2, 2, "Current gamma values        : %.2f, %.2f, %.2f",
			gammaValues[0], gammaValues[1], gammaValues[2])
		winGrid.NoutRefresh()
		winHelp.NoutRefresh()

		gc.Update()

		ch = winGrid.GetChar()

		ledColor = ledGrid.LedColorAt(curCol, curRow)

		if (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') {
			val, err := strconv.ParseUint(string(ch), 16, 8)
			if err != nil {
				log.Fatalf("Couldn't convert value: %v", err)
			}
			newColorValue = (newColorValue << 4) | uint8(val)
			switch newValueDigit {
			case 0:
				enterNewValue = true
				newValueDigit++
			case 1:
				switch curColor {
				case 0:
					ledColor.R = newColorValue
				case 1:
					ledColor.G = newColorValue
				case 2:
					ledColor.B = newColorValue
				}
				ledGrid.SetLedColor(curCol, curRow, ledColor)
				if curColor < 2 {
					curColor++
				} else {
					curColor = 0
					if curCol < ledGrid.Bounds().Dx()-1 {
						curCol++
						selCol = curCol
						cursorMoved = true
					}
				}

				redrawGrid = true
				enterNewValue = false
				newValueDigit = 0
			}
		}

		switch ch {

		case 'C':
			ledGrid.Clear(color.Black)
			curCol, selCol = 0, 0
			curRow, selRow = 0, 0
			redrawGrid = true

		case 'F':
			if selRect.Dy() > 2 {
				col := selRect.Min.X
				color0 := ledGrid.LedColorAt(col, selRect.Min.Y)
				color1 := ledGrid.LedColorAt(col, selRect.Max.Y-1)
				for row := selRect.Min.Y; row < selRect.Max.Y; row++ {
					t := float64(row-selRect.Min.Y) / float64(selRect.Dy()-1)
					color2 := color0.Interpolate(color1, t)
					ledGrid.SetLedColor(col, row, color2)
				}
				col = selRect.Max.X - 1
				color0 = ledGrid.LedColorAt(col, selRect.Min.Y)
				color1 = ledGrid.LedColorAt(col, selRect.Max.Y-1)
				for row := selRect.Min.Y; row < selRect.Max.Y; row++ {
					t := float64(row-selRect.Min.Y) / float64(selRect.Dy()-1)
					color2 := color0.Interpolate(color1, t)
					ledGrid.SetLedColor(col, row, color2)
				}
			}
			if selRect.Dx() > 2 {
				for row := selRect.Min.Y; row < selRect.Max.Y; row++ {
					color0 := ledGrid.LedColorAt(selRect.Min.X, row)
					color1 := ledGrid.LedColorAt(selRect.Max.X-1, row)
					for col := selRect.Min.X; col < selRect.Max.X; col++ {
						t := float64(col-selRect.Min.X) / float64(selRect.Dx()-1)
						color2 := color0.Interpolate(color1, t)
						ledGrid.SetLedColor(col, row, color2)
					}
				}
			}
			redrawGrid = true

		case 'g':
			gammaValues[0] -= 0.1
			gammaValues[1] -= 0.1
			gammaValues[2] -= 0.1
			pixelClient.SetGamma(gammaValues[0], gammaValues[1], gammaValues[2])
			pixelClient.Send(ledGrid)

		case 'G':
			gammaValues[0] += 0.1
			gammaValues[1] += 0.1
			gammaValues[2] += 0.1
			pixelClient.SetGamma(gammaValues[0], gammaValues[1], gammaValues[2])
			pixelClient.Send(ledGrid)

		case 'i':
			for row := selRect.Min.Y; row < selRect.Max.Y; row++ {
				for col := selRect.Min.X; col < selRect.Max.X; col++ {
					c := ledGrid.LedColorAt(col, row)
					c.R = 255 - c.R
					c.G = 255 - c.G
					c.B = 255 - c.B
					ledGrid.SetLedColor(col, row, c)
				}
			}
			redrawGrid = true

		case '+':
			for row := selRect.Min.Y; row < selRect.Max.Y; row++ {
				for col := selRect.Min.X; col < selRect.Max.X; col++ {
					c := ledGrid.LedColorAt(col, row)
					ledGrid.SetLedColor(col, row, c.Bright(0.1))
				}
			}
			redrawGrid = true

		case '-':
			for row := selRect.Min.Y; row < selRect.Max.Y; row++ {
				for col := selRect.Min.X; col < selRect.Max.X; col++ {
					c := ledGrid.LedColorAt(col, row)
					ledGrid.SetLedColor(col, row, c.Dark(0.1))
				}
			}
			redrawGrid = true

		case 'q':
			break main

		case Ctrl('a'):
			curRow, curCol = 0, 0
			selRow, selCol = curRow, curCol
			selRect = ledGrid.Bounds()

		case Ctrl('c'), Ctrl('x'):
			clipRect = selRect
			clipData = make([]color.LedColor, clipRect.Dx()*clipRect.Dy())
			for y := range clipRect.Dy() {
				row := clipRect.Min.Y + y
				for x := range clipRect.Dx() {
					col := clipRect.Min.X + x
					idx := clipRect.Dx()*y + x
					clipData[idx] = ledGrid.LedColorAt(col, row)
				}
			}
			if ch == Ctrl('x') {
				for row := selRect.Min.Y; row < selRect.Max.Y; row++ {
					for col := selRect.Min.X; col < selRect.Max.X; col++ {
						ledGrid.SetLedColor(col, row, color.Black)
					}
				}
				redrawGrid = true
			}

		case Ctrl('v'):
			dstRect := selRect
			if dstRect.Dx() == 1 && dstRect.Dy() == 1 {
				dstRect = clipRect.Add(dstRect.Min.Sub(clipRect.Min))
			}
			for y := range dstRect.Dy() {
				clipY := y % clipRect.Dy()
				row := dstRect.Min.Y + y
				for x := range dstRect.Dx() {
					clipX := x % clipRect.Dx()
					col := dstRect.Min.X + x
					idx := clipRect.Dx()*clipY + clipX
					ledGrid.SetLedColor(col, row, clipData[idx])
					redrawGrid = true
				}
			}

		case Ctrl('s'):
			fh, err := os.Create("ledgrid.png")
			if err != nil {
				log.Fatalf("Couldn't create file: %v", err)
			}
			err = png.Encode(fh, ledGrid)
			if err != nil {
				log.Fatalf("Couldn't encode file: %v", err)
			}
			fh.Close()

		case Ctrl('o'):
			fh, err := os.Open("ledgrid.png")
			if err != nil {
				log.Fatalf("Couldn't open file: %v", err)
			}
			img, err := png.Decode(fh)
			if err != nil {
				log.Fatalf("Couldn't decode file: %v", err)
			}
			fh.Close()
			draw.Draw(ledGrid, ledGrid.Bounds(), img, image.Point{}, draw.Over)
			redrawGrid = true

		case gc.KEY_BACKSPACE:
			for row := selRect.Min.Y; row < selRect.Max.Y; row++ {
				for col := selRect.Min.X; col < selRect.Max.X; col++ {
					ledGrid.SetLedColor(col, row, color.Black)
				}
			}
			redrawGrid = true

		case gc.KEY_LEFT:
			if curCol > 0 {
				curCol -= 1
			}
			selCol = curCol
			selRow = curRow
			cursorMoved = true
		case gc.KEY_RIGHT:
			if curCol < ledGrid.Rect.Dx()-1 {
				curCol += 1
			}
			selCol = curCol
			selRow = curRow
			cursorMoved = true
		case gc.KEY_UP:
			if curRow > 0 {
				curRow -= 1
			}
			selCol = curCol
			selRow = curRow
			cursorMoved = true
		case gc.KEY_DOWN:
			if curRow < ledGrid.Rect.Dy()-1 {
				curRow += 1
			}
			selCol = curCol
			selRow = curRow
			cursorMoved = true

		case gc.KEY_SLEFT:
			if curCol > 0 {
				curCol -= 1
			}
			cursorMoved = true
		case gc.KEY_SRIGHT:
			if curCol < ledGrid.Rect.Dx()-1 {
				curCol += 1
			}
			cursorMoved = true
		case KEY_SUP:
			if curRow > 0 {
				curRow -= 1
			}
			cursorMoved = true
		case KEY_SDOWN:
			if curRow < ledGrid.Rect.Dy()-1 {
				curRow += 1
			}
			cursorMoved = true

		case KEY_ALEFT:
			if curColor > 0 {
				curColor--
			}
		case KEY_ARIGHT:
			if curColor < 2 {
				curColor++
			}

		case gc.KEY_RETURN:
			switch curColor {
			case 0:
				ledColor.R = newColorValue
			case 1:
				ledColor.G = newColorValue
			case 2:
				ledColor.B = newColorValue
			}
			enterNewValue = false
			curColorChanged = true

		case gc.KEY_IC, gc.KEY_DC, KEY_AINS, KEY_ADEL:
			curColor = 0
			if ch == gc.KEY_IC || ch == gc.KEY_DC {
				incr = 1
			} else {
				incr = 16
			}
			if ch == gc.KEY_IC || ch == KEY_AINS {
				ledColor.R += incr
			} else {
				ledColor.R -= incr
			}
			curColorChanged = true

		case gc.KEY_HOME, gc.KEY_END, KEY_AHOME, KEY_AEND:
			curColor = 1
			if ch == gc.KEY_HOME || ch == gc.KEY_END {
				incr = 1
			} else {
				incr = 16
			}
			if ch == gc.KEY_HOME || ch == KEY_AHOME {
				ledColor.G += incr
			} else {
				ledColor.G -= incr
			}
			curColorChanged = true

		case gc.KEY_PAGEUP, gc.KEY_PAGEDOWN, KEY_APAGEUP, KEY_APAGEDOWN:
			curColor = 2
			if ch == gc.KEY_PAGEUP || ch == gc.KEY_PAGEDOWN {
				incr = 1
			} else {
				incr = 16
			}
			if ch == gc.KEY_PAGEUP || ch == KEY_APAGEUP {
				ledColor.B += incr
			} else {
				ledColor.B -= incr
			}
			curColorChanged = true

		default:
			fmt.Fprintf(os.Stderr, "Unhandled key: 0x%02x, '%s'\n", ch, gc.KeyString(ch))
		}

		if cursorMoved {
			selRect = image.Rect(curCol, curRow, selCol, selRow)
			selRect.Max = selRect.Max.Add(image.Point{1, 1})
			cursorMoved = false
		}

		if curColorChanged {
			ledGrid.SetLedColor(curCol, curRow, ledColor)
			curColorChanged = false
			redrawGrid = true
		}

		if redrawGrid {
			pixelClient.Send(ledGrid)
			redrawGrid = false
		}
	}
	winGrid.Delete()
	pixelClient.Close()
}
