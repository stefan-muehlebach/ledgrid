package main

import (
	"image"
	"log"
	"strconv"

	gc "github.com/rthornton128/goncurses"
	"github.com/stefan-muehlebach/ledgrid"
)

const (
	width         = 20
	height        = 10
    gridWidth     = 8 * width + 9
    gridHeight    = height + 7
	termWidth     = gridWidth + 10
	termHeight    = gridHeight + 40
	defHost       = "raspi-2"
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
	var ledColor ledgrid.LedColor
	var pixelClient ledgrid.PixelClient
	var colorChanged bool
	var colors []uint8
	var incr uint8
	var newColorValue uint8
	var gammaValues [3]float64

	ledGrid = ledgrid.NewLedGrid(image.Point{width, height})
	pixelClient = ledgrid.NewNetPixelClient(defHost, defPort)
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
	gc.InitPair(2, gc.C_GREEN, gc.C_BLACK)
	gc.InitPair(3, gc.C_BLUE, gc.C_WHITE)

	// rows, cols := stdscr.MaxYX()

	// gridHeight, gridWidth := height+7, 8*width+9
	y, x := 2, 4

	winGrid, err = gc.NewWindow(gridHeight, gridWidth, y, x)
	if err != nil {
		log.Fatalf("Couldn't create window: %v", err)
	}
	winGrid.Keypad(true)
	winGrid.Box(0, 0)

	helpHeight, helpWidth := 18, 55
	y, x = height+9, 4

	winHelp, err = gc.NewWindow(helpHeight, helpWidth, y, x)
	if err != nil {
		log.Fatalf("Couldn't create window: %v", err)
	}
	// winHelp.Keypad(true)
	winHelp.Box(0, 0)

	winHelp.MovePrintf(1, 2, "[Cursor]: Move selector")
	winHelp.MovePrintf(2, 2, "[Shift]-[Cursor]: Select range")
	winHelp.MovePrintf(6, 2, "  R   |    G   |    B   |")
	winHelp.MovePrintf(7, 2, "------+--------+--------+")
	winHelp.MovePrintf(8, 2, "[Ins] | [Home] | [PgUp] | increase color value")
	winHelp.MovePrintf(9, 2, "[Del] | [End]  | [PgDn] | decrease color value")
	winHelp.MovePrintf(11, 2, "C: clear panel")
	winHelp.MovePrintf(12, 2, "F: interpolate colors")
	winHelp.MovePrintf(13, 2, "0-9a-f: enter new hex value for selected color")
	winHelp.MovePrintf(14, 2, "g/G: decrease/increase gamma values by 0.1")
	winHelp.MovePrintf(16, 2, "q: Quit")

	// fields := make([]*gc.Field, 3)
	// fields[0], _ = gc.NewField(1, 3, 14, 16, 0, 0)
	// defer fields[0].Free()
	// fields[0].SetBackground(gc.A_UNDERLINE)

	// fields[1], _ = gc.NewField(1, 3, 14, 21, 0, 0)
	// defer fields[1].Free()
	// fields[1].SetBackground(gc.A_UNDERLINE)

	// fields[2], _ = gc.NewField(1, 3, 14, 21, 0, 0)
	// defer fields[2].Free()
	// fields[2].SetBackground(gc.A_UNDERLINE)

	// form, _ := gc.NewForm(fields)
	// defer form.UnPost()
	// defer form.Free()

	// form.SetWindow(winGrid)
	// form.SetSub(winGrid.Derived(1, 11, 14, 16))
	// form.Post()

	winGrid.Refresh()
	winHelp.Refresh()

	pixelClient.Draw(ledGrid)

main:
	for {
		for col := 0; col < ledGrid.Rect.Dx(); col++ {
			if col == curCol {
				winGrid.AttrOn(gc.A_BOLD)
			}
			winGrid.MovePrintf(1, 10+col*8, "[%02d]", col)
			winGrid.AttrOff(gc.A_BOLD)
		}
		winGrid.HLine(2, 1, gc.ACS_HLINE, gridWidth-2)
		for row := 0; row < ledGrid.Rect.Dy(); row++ {
			if row == curRow {
				winGrid.AttrOn(gc.A_BOLD)
			}
			winGrid.MovePrintf(3+row, 2, "[%02d]", row)
			winGrid.AttrOff(gc.A_BOLD)
		}
		winGrid.VLine(1, 7, gc.ACS_VLINE, height+2)
		for row := 0; row < ledGrid.Rect.Dy(); row++ {
			for col := 0; col < ledGrid.Rect.Dx(); col++ {
				if between(row, curRow, selRow) && between(col, curCol, selCol) {
					winGrid.AttrOn(gc.A_REVERSE)
					// win.MovePrintf(3+row, 8+(col*8), "        ")
				}
				ledColor = ledGrid.LedColorAt(col, row)
				colors = []uint8{ledColor.R, ledColor.G, ledColor.B}
				for k := 0; k < 3; k++ {
					if (row == curRow) && (col == curCol) && (k == curColor) {
						winGrid.ColorOn(int16(k + 1))
						winGrid.AttrOn(gc.A_BOLD)
					}
					winGrid.MovePrintf(3+row, 9+(col*8)+(k*2), "%02x", colors[k])
					winGrid.AttrOff(gc.A_BOLD)
					winGrid.ColorOff(int16(k + 1))
				}
				winGrid.AttrOff(gc.A_REVERSE)
			}
		}
        row := height + 3
		winGrid.HLine(row, 1, gc.ACS_HLINE, gridWidth-2)

		winGrid.MovePrintf(row+1, 2, "New hex value for this color: %02x", newColorValue)
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
			continue
		}

		switch ch {

		case 'C':
			for i := range ledGrid.Pix {
				ledGrid.Pix[i] = 0x00
			}
			ledColor = ledGrid.LedColorAt(curCol, curRow)
			colorChanged = true

		case 'F':
			r := image.Rect(selCol, selRow, curCol, curRow)
			if r.Dy() >= 2 {
				col := r.Min.X
				color0 := ledGrid.LedColorAt(col, r.Min.Y)
				color1 := ledGrid.LedColorAt(col, r.Max.Y)
				for row := r.Min.Y; row <= r.Max.Y; row++ {
					t := float64(row-r.Min.Y) / float64(r.Dy())
					color := color0.Interpolate(color1, t).(ledgrid.LedColor)
					ledGrid.SetLedColor(col, row, color)
				}
				col = r.Max.X
				color0 = ledGrid.LedColorAt(col, r.Min.Y)
				color1 = ledGrid.LedColorAt(col, r.Max.Y)
				for row := r.Min.Y; row <= r.Max.Y; row++ {
					t := float64(row-r.Min.Y) / float64(r.Dy())
					color := color0.Interpolate(color1, t).(ledgrid.LedColor)
					ledGrid.SetLedColor(col, row, color)
				}
			}
			if r.Dx() >= 2 {
				for row := r.Min.Y; row <= r.Max.Y; row++ {
					color0 := ledGrid.LedColorAt(r.Min.X, row)
					color1 := ledGrid.LedColorAt(r.Max.X, row)
					for col := r.Min.X; col <= r.Max.X; col++ {
						t := float64(col-r.Min.X) / float64(r.Dx())
						color := color0.Interpolate(color1, t).(ledgrid.LedColor)
						ledGrid.SetLedColor(col, row, color)
					}
				}
			}
			ledColor = ledGrid.LedColorAt(curCol, curRow)
			selCol = curCol
			selRow = curRow
			colorChanged = true

		case 'g':
			gammaValues[0] -= 0.1
			gammaValues[1] -= 0.1
			gammaValues[2] -= 0.1
			pixelClient.SetGamma(gammaValues[0], gammaValues[1], gammaValues[2])
			pixelClient.Draw(ledGrid)

		case 'G':
			gammaValues[0] += 0.1
			gammaValues[1] += 0.1
			gammaValues[2] += 0.1
			pixelClient.SetGamma(gammaValues[0], gammaValues[1], gammaValues[2])
			pixelClient.Draw(ledGrid)

		case 'q':
			break main

		// case gc.KEY_TAB:
		// 	form.Driver(gc.REQ_NEXT_FIELD)
		// 	form.Driver(gc.REQ_END_LINE)

		case gc.KEY_LEFT:
			if curCol > 0 {
				curCol -= 1
			}
			selCol = curCol
			selRow = curRow
		case gc.KEY_RIGHT:
			if curCol < ledGrid.Rect.Dx()-1 {
				curCol += 1
			}
			selCol = curCol
			selRow = curRow
		case gc.KEY_UP:
			if curRow > 0 {
				curRow -= 1
			}
			selCol = curCol
			selRow = curRow
		case gc.KEY_DOWN:
			if curRow < ledGrid.Rect.Dy()-1 {
				curRow += 1
			}
			selCol = curCol
			selRow = curRow

		case gc.KEY_SLEFT:
			if selCol > 0 {
				selCol -= 1
			}
		case gc.KEY_SRIGHT:
			if selCol < ledGrid.Rect.Dx()-1 {
				selCol += 1
			}
		case KEY_SUP:
			if selRow > 0 {
				selRow -= 1
			}
		case KEY_SDOWN:
			if selRow < ledGrid.Rect.Dy()-1 {
				selRow += 1
			}

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
			colorChanged = true

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
			colorChanged = true

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
			colorChanged = true

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
			colorChanged = true

			// default:
			// 	fmt.Fprintf(os.Stderr, "Unhandled key: 0x%02x, '%s'\n", ch, gc.KeyString(ch))
		}

		if colorChanged {
			ledGrid.SetLedColor(curCol, curRow, ledColor)
			pixelClient.Draw(ledGrid)
			colorChanged = false
		}
	}
	winGrid.Delete()
	pixelClient.Close()
}
