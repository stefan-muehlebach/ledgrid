package main

import (
	"fmt"
	"image"
	"log"
	"os"

	gc "github.com/gbin/goncurses"
	"github.com/stefan-muehlebach/ledgrid"
)

const (
	width     = 10
	height    = 10
	defHost   = "raspi-2"
	defPort   = 5333
	KEY_SUP   = 0x151 /* Shifted up arrow */
	KEY_SDOWN = 0x150 /* Shifted down arrow */
)

func between(x, a, b int) bool {
	if a > b {
		a, b = b, a
	}
	return x >= a && x <= b
}

func main() {
	var stdscr, win *gc.Window
	var ch gc.Key
	var curRow, selRow, curCol, selCol int
	var curColor int
	var err error
	var ledGrid *ledgrid.LedGrid
	var ledColor ledgrid.LedColor
	var pixelClient *ledgrid.PixelClient
	var colorChanged bool
	var colors []uint8
	var incr uint8
	// var gammaValues []float64 = []float64{1.0, 1.0, 1.0}

	ledGrid = ledgrid.NewLedGrid(image.Rect(0, 0, width, height))
	pixelClient = ledgrid.NewPixelClient(defHost, defPort)

	stdscr, err = gc.Init()
	if err != nil {
		log.Fatalf("Couldn't Init ncurses: %v", err)
	}
	defer gc.End()

	gc.StartColor()
	gc.Echo(false)
	gc.CBreak(true)
	gc.Cursor(0)

	gc.InitPair(1, gc.C_RED, gc.C_BLACK)
	gc.InitPair(2, gc.C_GREEN, gc.C_BLACK)
	gc.InitPair(3, gc.C_BLUE, gc.C_WHITE)

	rows, cols := stdscr.MaxYX()
	height, width := 16, 89
	y, x := (rows-height)/2, (cols-width)/2
	y, x = 2, 4

	win, err = gc.NewWindow(height, width, y, x)
	if err != nil {
		log.Fatalf("Couldn't create window: %v", err)
	}
	win.Keypad(true)
	win.Box(0, 0)

    fields := make([]*gc.Field, 3)
    fields[0], _ = gc.NewField(1, 3, 14, 16, 0, 0)
    defer fields[0].Free()
	fields[0].SetBackground(gc.A_UNDERLINE)

    fields[1], _ = gc.NewField(1, 3, 14, 21, 0, 0)
    defer fields[1].Free()
	fields[1].SetBackground(gc.A_UNDERLINE)

    fields[2], _ = gc.NewField(1, 3, 14, 21, 0, 0)
    defer fields[2].Free()
	fields[2].SetBackground(gc.A_UNDERLINE)

    form, _ := gc.NewForm(fields)
    defer form.UnPost()
    defer form.Free()

    form.SetWindow(win)
    form.SetSub(win.Derived(1, 11, 14, 16))
    form.Post()

    win.Refresh()

main:
	for {
		for col := 0; col < ledGrid.Rect.Dx(); col++ {
			if col == curCol {
				win.AttrOn(gc.A_BOLD)
			}
			win.MovePrintf(1, 10+col*8, "[%02d]", col)
			win.AttrOff(gc.A_BOLD)
		}
		win.HLine(2, 1, gc.ACS_HLINE, width-2)
		for row := 0; row < ledGrid.Rect.Dy(); row++ {
			if row == curRow {
				win.AttrOn(gc.A_BOLD)
			}
			win.MovePrintf(3+row, 2, "[%02d]", row)
			win.AttrOff(gc.A_BOLD)
		}
		win.VLine(1, 7, gc.ACS_VLINE, 12)
		for row := 0; row < ledGrid.Rect.Dy(); row++ {
			for col := 0; col < ledGrid.Rect.Dx(); col++ {
				if between(row, curRow, selRow) && between(col, curCol, selCol) {
					win.AttrOn(gc.A_REVERSE)
					// win.MovePrintf(3+row, 8+(col*8), "        ")
				}
				ledColor = ledGrid.LedColorAt(col, row)
				colors = []uint8{ledColor.R, ledColor.G, ledColor.B}
				for k := 0; k < 3; k++ {
					if (row == curRow) && (col == curCol) && (k == curColor) {
						win.ColorOn(int16(k + 1))
						win.AttrOn(gc.A_BOLD)
					}
					win.MovePrintf(3+row, 9+(col*8)+(k*2), "%02x", colors[k])
					win.AttrOff(gc.A_BOLD)
					win.ColorOff(int16(k + 1))
				}
				win.AttrOff(gc.A_REVERSE)
			}
		}
		win.HLine(13, 1, gc.ACS_HLINE, width-2)

		win.MovePrintf(14, 2, "Gamma values:")
		win.Refresh()

		gc.Update()

		ch = win.GetChar()

		ledColor = ledGrid.LedColorAt(curCol, curRow)

		switch ch {

		case 'R':
			curColor = 0
		case 'G':
			curColor = 1
		case 'B':
			curColor = 2

		case 'C':
			for i := range ledGrid.Pix {
				ledGrid.Pix[i] = 0x00
			}
			ledColor = ledGrid.LedColorAt(curCol, curRow)
			colorChanged = true

		case 'f':
			r := image.Rect(min(selCol, curCol), min(selRow, curRow),
				max(selCol, curCol), max(selRow, curRow))
			if r.Dy() >= 2 {
				col := r.Min.X
				color0 := ledGrid.LedColorAt(col, r.Min.Y)
				color1 := ledGrid.LedColorAt(col, r.Max.Y)
				for row := r.Min.Y; row <= r.Max.Y; row++ {
					t := float64(row-r.Min.Y) / float64(r.Dy())
					color := color0.Interpolate(color1, t)
					ledGrid.SetLedColor(col, row, color)
				}
				col = r.Max.X
				color0 = ledGrid.LedColorAt(col, r.Min.Y)
				color1 = ledGrid.LedColorAt(col, r.Max.Y)
				for row := r.Min.Y; row <= r.Max.Y; row++ {
					t := float64(row-r.Min.Y) / float64(r.Dy())
					color := color0.Interpolate(color1, t)
					ledGrid.SetLedColor(col, row, color)
				}
			}
			if r.Dx() >= 2 {
				for row := r.Min.Y; row <= r.Max.Y; row++ {
					color0 := ledGrid.LedColorAt(r.Min.X, row)
					color1 := ledGrid.LedColorAt(r.Max.X, row)
					for col := r.Min.X; col <= r.Max.X; col++ {
						t := float64(col-r.Min.X) / float64(r.Dx())
						color := color0.Interpolate(color1, t)
						ledGrid.SetLedColor(col, row, color)
					}
				}
			}
			ledColor = ledGrid.LedColorAt(curCol, curRow)
			colorChanged = true

		case 'q':
			break main

        case gc.KEY_TAB:
            form.Driver(gc.REQ_NEXT_FIELD)
            form.Driver(gc.REQ_END_LINE)

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

		case gc.KEY_IC, gc.KEY_SIC, gc.KEY_DC, gc.KEY_SDC:
			curColor = 0
			incr = 1
			if ch == gc.KEY_SIC || ch == gc.KEY_SDC {
				incr = 16
			}
			if ch == gc.KEY_IC || ch == gc.KEY_SIC {
				ledColor.R += incr
			} else {
				ledColor.R -= incr
			}
			colorChanged = true

		case gc.KEY_HOME, gc.KEY_SHOME, gc.KEY_END, gc.KEY_SEND:
			curColor = 1
			incr = 1
			if ch == gc.KEY_SHOME || ch == gc.KEY_SEND {
				incr = 16
			}
			if ch == gc.KEY_HOME || ch == gc.KEY_SHOME {
				ledColor.G += incr
			} else {
				ledColor.G -= incr
			}
			colorChanged = true

		case gc.KEY_PAGEUP, gc.KEY_SPREVIOUS, gc.KEY_PAGEDOWN, gc.KEY_SNEXT:
			curColor = 2
			incr = 1
			if ch == gc.KEY_SNEXT || ch == gc.KEY_SPREVIOUS {
				incr = 16
			}
			if ch == gc.KEY_PAGEUP || ch == gc.KEY_SPREVIOUS {
				ledColor.B += incr
			} else {
				ledColor.B -= incr
			}
			colorChanged = true

		default:
			fmt.Fprintf(os.Stderr, "Unhandled key: 0x%02x, '%s'\n", ch, gc.KeyString(ch))
		}

		if colorChanged {
			ledGrid.SetLedColor(curCol, curRow, ledColor)
			pixelClient.Draw(ledGrid)
			colorChanged = false
		}
	}
	win.Delete()
	pixelClient.Close()
}
