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

	gc "github.com/gbin/goncurses"
	// gc "github.com/rthornton128/goncurses"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/colors"
	"github.com/stefan-muehlebach/ledgrid/conf"
)

const (
	KEY_SUP       = 0x151 /* Shifted up arrow */
	KEY_SDOWN     = 0x150 /* Shifted down arrow */
	KEY_CLEFT     = 0x222 /* Ctrl-left arrow */
	KEY_CRIGHT    = 0x231 /* Ctrl-right arrow */
	KEY_CUP       = 0x237 /* Ctrl-up arrow */
	KEY_CDOWN     = 0x20e /* Ctrl-down arrow */
	KEY_ALEFT     = 0x228 /* Alt-left arrow */
	KEY_ARIGHT    = 0x237 /* Alt-right arrow */
	KEY_AUP       = 0x23d /* Alt-up arrow */
	KEY_ADOWN     = 0x214 /* Alt-down arrow */
	KEY_AINS      = 0x223 /* Alt-Insert */
	KEY_ADEL      = 0x20e /* Alt-Delete */
	KEY_AHOME     = 0x21e /* Alt-Home */
	KEY_AEND      = 0x219 /* Alt-End */
	KEY_APAGEUP   = 0x232 /* Alt-PageUp */
	KEY_APAGEDOWN = 0x22d /* Alt-PageDown */
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

func LogMouseEvent(log *os.File, event *gc.MouseEvent) {
	fmt.Fprintf(log, "Mouse event (state: %032b) ", event.State)
	switch {

	case event.State&gc.M_B1_CLICKED != 0:
		fmt.Fprintf(log, "B1 clicked\n")
	case event.State&gc.M_B1_DBL_CLICKED != 0:
		fmt.Fprintf(log, "B1 double clicked\n")
	case event.State&gc.M_B1_TPL_CLICKED != 0:
		fmt.Fprintf(log, "B1 tripple clicked\n")
	case event.State&gc.M_B1_PRESSED != 0:
		fmt.Fprintf(log, "B1 pressed\n")
	case event.State&gc.M_B1_RELEASED != 0:
		fmt.Fprintf(log, "B1 released\n")

	case event.State&gc.M_B2_CLICKED != 0:
		fmt.Fprintf(log, "B2 clicked\n")
	case event.State&gc.M_B2_DBL_CLICKED != 0:
		fmt.Fprintf(log, "B2 double clicked\n")
	case event.State&gc.M_B2_TPL_CLICKED != 0:
		fmt.Fprintf(log, "B2 tripple clicked\n")
	case event.State&gc.M_B2_PRESSED != 0:
		fmt.Fprintf(log, "B2 pressed\n")
	case event.State&gc.M_B2_RELEASED != 0:
		fmt.Fprintf(log, "B2 released\n")

	default:
		fmt.Fprintf(log, "\n")
	}
	fmt.Fprintf(log, "  at (%d,%d)\n", event.X, event.Y)
}

func main() {
	var host string = "raspi-3"
	var width, height int
	var dataPort, rpcPort uint
	var useTCP bool
	var network string

	var logFile *os.File
	var stdscr *gc.Window
	var winGrid, winHelp *gc.Window
	var ch gc.Key
	var curRow, selRow, curCol, selCol int
	var curColor int

	var err error
	var gridClient ledgrid.GridClient
	var ledGrid *ledgrid.LedGrid
	var ledColor colors.LedColor
	var palMode bool
	var palIdx int
	var curColorChanged bool
	var redrawGrid bool
	var cursorMoved bool
	var colorValues []uint8
	var incr uint8
	var newColorValue uint8
	var enterNewValue bool
	var newValueDigit int
	var gammaValues [3]float64
	var gridXMin, gridYMin, gridWidth, gridHeight int
	var helpXMin, helpYMin, helpWidth, helpHeight int
	var termWidth, termHeight int
	var gridSize image.Point
	var selRect image.Rectangle
	var clipRect image.Rectangle
	var clipData []colors.LedColor
	var modConf conf.ModuleConfig

	flag.StringVar(&host, "host", host, "Controller hostname")
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

	gridSize = ledGrid.Rect.Size()
	width = gridSize.X
	height = gridSize.Y

	// gridWidth und gridHeight sind die nCurses Abmessungen des Fensters,
	// welches die Farbwerte des Grids enthaelt.
	gridXMin, gridYMin = 4, 2
	gridWidth = 7*width + 10
	gridHeight = height + 8

	helpXMin, helpYMin = 4, height+10
	helpWidth = 70
	helpHeight = 22

	// helpHeight, helpWidth := 22, 70
	// y, x := height+10, 4

	termWidth = gridWidth + 10
	termHeight = gridHeight + 40

	selRect = image.Rect(0, 0, 1, 1)

	logFile, err = os.Create("colorEdit.log")
	if err != nil {
		log.Fatalf("Couldn't create log file: %v", err)
	}
	defer logFile.Close()

	// if modConf != nil {
	// ledGrid = ledgrid.NewLedGrid(host, port, modConf)
	// } else {
	// ledGrid = ledgrid.NewLedGridBySize(host, port, gridSize)
	// }
	gammaValues[0], gammaValues[1], gammaValues[2] = ledGrid.Client.Gamma()

	stdscr, err = gc.Init()
	if err != nil {
		log.Fatalf("Couldn't Init ncurses: %v", err)
	}
	defer gc.End()

	termHeight, termWidth = stdscr.MaxYX()
	fmt.Fprintf(logFile, "window size: %d, %d\n", termWidth, termHeight)

	gridWidth = min(gridWidth, termWidth-8)
	// gridHeight = min(gridHeight, termHeight-40)

	// err = gc.ResizeTerm(termHeight, termWidth)
	// if err != nil {
	// 	log.Fatalf("Couldn't resize terminal: %v", err)
	// }

	gc.StartColor()
	gc.Echo(false)
	gc.CBreak(true)
	gc.Cursor(0)
	gc.Raw(true)

	stdscr.Keypad(true)

	black := int16(0)
	red := int16(1)
	green := int16(2)
	// yellow := int16(3)
	blue := int16(4)
	// magenta := int16(5)
	// cyan := int16(6)
	white := int16(7)

	toM := func(v uint8) int16 { return int16(float64(v) * 1000.0 / 255.0) }

	// Setup dark colors (R, G, B)
	dark := func(colIdx int16) int16 { return 16 + colIdx }
	gc.InitColor(dark(red), toM(192), toM(28), toM(40))
	gc.InitColor(dark(green), 149, 635, 412)
	gc.InitColor(dark(blue), 71, 282, 545)

	// Setup bright colors
	bright := func(colIdx int16) int16 { return 20 + colIdx }
	gc.InitColor(bright(red), 1000, 110, 157)
	gc.InitColor(bright(green), 86, 894, 455)
	gc.InitColor(bright(blue), toM(41), toM(173), toM(255))
	// gc.InitColor(bright(blue), 184, 506, 1000)

	// Setup light colors
	light := func(colIdx int16) int16 { return 24 + colIdx }
	gc.InitColor(light(red), 918, 424, 459)
	gc.InitColor(light(green), 353, 847, 619)
	gc.InitColor(light(blue), 255, 545, 902)

	gc.InitPair(1, gc.C_WHITE, gc.C_BLACK)

	comb := func(fg, bg int16) int16 { return 2 + (fg - 16) + 16*bg }
	InitPair := func(fg, bg int16) { gc.InitPair(comb(fg, bg), fg, bg) }
	InitPair(dark(red), black)
	InitPair(dark(green), black)
	InitPair(dark(blue), black)
	InitPair(bright(red), black)
	InitPair(bright(green), black)
	InitPair(bright(blue), black)
	InitPair(light(red), black)
	InitPair(light(green), black)
	InitPair(light(blue), black)

	InitPair(dark(red), white)
	InitPair(dark(green), white)
	InitPair(dark(blue), white)
	InitPair(bright(red), white)
	InitPair(bright(green), white)
	InitPair(bright(blue), white)
	InitPair(light(red), white)
	InitPair(light(green), white)
	InitPair(light(blue), white)

	InitPair(black, dark(red))
	InitPair(black, dark(green))
	InitPair(black, dark(blue))
	InitPair(black, bright(red))
	InitPair(black, bright(green))
	InitPair(black, bright(blue))
	InitPair(black, light(red))
	InitPair(black, light(green))
	InitPair(black, light(blue))

	InitPair(white, dark(red))
	InitPair(white, dark(green))
	InitPair(white, dark(blue))
	InitPair(white, bright(red))
	InitPair(white, bright(green))
	InitPair(white, bright(blue))
	InitPair(white, light(red))
	InitPair(white, light(green))
	InitPair(white, light(blue))

	// gc.MouseInterval(50)
	gc.MouseMask(gc.M_ALL, nil)

	// y, x := gridY0, gridX0

	winGrid, err = gc.NewWindow(gridHeight, gridWidth, gridYMin, gridXMin)
	if err != nil {
		log.Fatalf("Couldn't create window: %v", err)
	}
	winGrid.Keypad(true)
	winGrid.Box(0, 0)

	winHelp, err = gc.NewWindow(helpHeight, helpWidth, helpYMin, helpXMin)
	if err != nil {
		log.Fatalf("Couldn't create window: %v", err)
	}
	winHelp.Box(0, 0)

	winHelp.MovePrintf(1, 2, "|   R   |    G   |    B   |")
	winHelp.MovePrintf(2, 2, "+-------+--------+--------+")
	winHelp.MovePrintf(3, 2, "| [Ins] | [Home] | [PgUp] | increase value by 1 ([Alt]: by 10)")
	winHelp.MovePrintf(4, 2, "| [Del] | [End]  | [PgDn] | decrease value by 1 ([Alt]: by 10)")
	winHelp.MovePrintf(6, 2, "[Cursor]        : Move selector")
	winHelp.MovePrintf(7, 2, "[Alt]-Left/Right: Move color selector")
	winHelp.MovePrintf(8, 2, "[Shift]-[Cursor]: Select range")
	winHelp.MovePrintf(9, 2, "[Ctrl]-a        : Select all pixels")
	winHelp.MovePrintf(10, 2, "[Ctrl]-c/x/v    : Copy/Cut/Paste selected pixels")
	winHelp.MovePrintf(11, 2, "[Backspace]     : Clear selected pixels")
	winHelp.MovePrintf(12, 2, "F               : Interpolate over selected range")
	winHelp.MovePrintf(13, 2, "0-9a-f          : Enter new hex value for selected pixel")
	winHelp.MovePrintf(14, 2, "g/G             : Decrease/increase gamma values by 0.1")
	winHelp.MovePrintf(15, 2, "i               : Invert selected pixels")
	winHelp.MovePrintf(16, 2, "+/-             : Darken/Brighten selected pixels")
	winHelp.MovePrintf(17, 2, "[Ctrl]-l/s      : Load/Save pixel data from/to file 'ledgrid.png'")
	winHelp.MovePrintf(18, 2, "[Ctrl]-p        : Toggle palette mode on/off")
	winHelp.MovePrintf(19, 2, "q               : Quit")

	winGrid.Refresh()
	winHelp.Refresh()

	ledGrid.Client.Send(ledGrid.Pix)

main:
	for {
		// Draw the upper part of the TUI (the big grid with the pixel values.
		//
		// The column headers.
		for col := 0; col < ledGrid.Rect.Dx(); col++ {
			if col == curCol {
				winGrid.AttrOn(gc.A_BOLD)
			}
			winGrid.MovePrintf(1, 10+col*7, "[%02x]", col)
			winGrid.AttrOff(gc.A_BOLD)
		}
		// The line after the column headers.
		winGrid.MoveAddChar(2, 0, gc.ACS_LTEE)
		winGrid.HLine(2, 1, gc.ACS_HLINE, gridWidth-2)
		winGrid.MoveAddChar(2, gridWidth-1, gc.ACS_RTEE)
		// The row 'headers' on the left side.
		for row := 0; row < ledGrid.Rect.Dy(); row++ {
			if row == curRow {
				winGrid.AttrOn(gc.A_BOLD)
			}
			winGrid.MovePrintf(3+row, 2, "[%02x]", row)
			winGrid.AttrOff(gc.A_BOLD)
		}
		// The separation line on the left side.
		row := height + 3
		winGrid.MoveAddChar(row, 0, gc.ACS_LTEE)
		winGrid.HLine(row, 1, gc.ACS_HLINE, gridWidth-2)
		winGrid.MoveAddChar(row, gridWidth-1, gc.ACS_RTEE)
		// The bottom separation line.
		winGrid.VLine(1, 7, gc.ACS_VLINE, height+2)
		winGrid.MoveAddChar(0, 7, gc.ACS_TTEE)
		winGrid.MoveAddChar(height+3, 7, gc.ACS_BTEE)
		winGrid.MoveAddChar(2, 7, gc.ACS_PLUS)

		// Draw now all the pixel values.
		for row := 0; row < ledGrid.Rect.Dy(); row++ {
			for col := 0; col < ledGrid.Rect.Dx(); col++ {
				pt := image.Point{col, row}
				ledColor = ledGrid.LedColorAt(col, row)
				colorValues = []uint8{ledColor.R, ledColor.G, ledColor.B}
				winGrid.Move(3+row, 9+(col*7))
				for id, colorIdx := range []int16{red, green, blue} {
					if (row == curRow) && (col == curCol) {
						winGrid.AttrOn(gc.A_REVERSE)
						if id == curColor {
							winGrid.AttrOn(gc.A_BOLD)
							winGrid.ColorOn(comb(bright(colorIdx), black))
						} else {
							winGrid.ColorOn(comb(dark(colorIdx), white))
						}
					} else if pt.In(selRect) {
						winGrid.ColorOn(comb(light(colorIdx), white))
					} else {
						winGrid.ColorOn(comb(light(colorIdx), black))
					}
					winGrid.Printf("%02X", colorValues[id])
					winGrid.AttrOff(gc.A_REVERSE)
					winGrid.AttrOff(gc.A_BOLD)
				}
			}
		}

		winGrid.ColorOn(1)
		winGrid.MovePrintf(row+1, 2, "New hex value for this color: ")
		if enterNewValue {
			winGrid.AttrOn(gc.A_REVERSE)
		}
		winGrid.Printf("%02x", newColorValue)
		winGrid.AttrOff(gc.A_REVERSE)
		winGrid.MovePrintf(row+2, 2, "Gamma values (R,G,B)        : %.2f, %.2f, %.2f",
			gammaValues[0], gammaValues[1], gammaValues[2])
		if palMode {
			winGrid.MovePrintf(row+3, 2, "Current palette shown       : ")
			winGrid.AttrOn(gc.A_BOLD)
			winGrid.Printf("%-20s", ledgrid.PaletteNames[palIdx])
			winGrid.AttrOff(gc.A_BOLD)
		} else {
			winGrid.MovePrintf(row+3, 2, "                                               ")
		}
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

		case gc.KEY_MOUSE:
			rowOff, colOff := 3, 9

			me := gc.GetMouse()
			// LogMouseEvent(logFile, me)
			if me.State&gc.M_B1_CLICKED != 0 {
				if !between(me.Y, gridYMin+rowOff, gridYMin+rowOff+10-1) ||
					!between(me.X, gridXMin+colOff, gridXMin+colOff+40*7-2) {
					break
				}
				dy, dx := me.Y-(gridYMin+rowOff), me.X-(gridXMin+colOff)
				row, col, field := dy, dx/7, (dx%7)/2
				if field > 2 {
					break
				}
				fmt.Fprintf(logFile, "col: %d, row: %d, field: %d\n", col, row, field)
				curCol, selCol = col, col
				curRow, selRow = row, row
				curColor = field
				cursorMoved = true
			}

		case 'C':
			ledGrid.Clear(colors.Transparent)
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
			ledGrid.Client.SetGamma(gammaValues[0], gammaValues[1], gammaValues[2])
			ledGrid.Client.Send(ledGrid.Pix)

		case 'G':
			gammaValues[0] += 0.1
			gammaValues[1] += 0.1
			gammaValues[2] += 0.1
			ledGrid.Client.SetGamma(gammaValues[0], gammaValues[1], gammaValues[2])
			ledGrid.Client.Send(ledGrid.Pix)

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
			clipData = make([]colors.LedColor, clipRect.Dx()*clipRect.Dy())
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
						ledGrid.SetLedColor(col, row, colors.Transparent)
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

		case Ctrl('l'):
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

		case Ctrl('p'):
			palMode = !palMode

		case gc.KEY_BACKSPACE:
			for row := selRect.Min.Y; row < selRect.Max.Y; row++ {
				for col := selRect.Min.X; col < selRect.Max.X; col++ {
					ledGrid.SetLedColor(col, row, colors.Transparent)
				}
			}
			redrawGrid = true

		case gc.KEY_LEFT:
			if palMode {
				if palIdx > 0 {
					palIdx--
				}
			} else {
				if curCol > 0 {
					curCol -= 1
				}
				selCol = curCol
				selRow = curRow
				cursorMoved = true
			}
		case gc.KEY_RIGHT:
			if palMode {
				if palIdx < len(ledgrid.PaletteNames)-1 {
					palIdx++
				}
			} else {
				if curCol < ledGrid.Rect.Dx()-1 {
					curCol += 1
				}
				selCol = curCol
				selRow = curRow
				cursorMoved = true
			}
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
			fmt.Fprintf(logFile, "Unhandled key: 0x%02x, '%s'\n", ch, gc.KeyString(ch))
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

		if palMode {
			palName := ledgrid.PaletteNames[palIdx]
			pal := ledgrid.PaletteMap[palName]
			for row := 0; row < ledGrid.Rect.Dy(); row++ {
				for col := 0; col < ledGrid.Rect.Dx(); col++ {
					t := float64(col) / float64(ledGrid.Rect.Dx()-1)
					c := pal.Color(t)
					ledGrid.SetLedColor(col, row, c)
				}
			}
			redrawGrid = true
		}

		if redrawGrid {
			ledGrid.Client.Send(ledGrid.Pix)
			redrawGrid = false
		}
	}
	winGrid.Delete()
	ledGrid.Client.Close()
}
