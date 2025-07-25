//go:build !tinygo

package main

import (
	"context"
	"flag"
	"fmt"
	"image"
	"math"
	"os"
	"os/signal"
	"time"

	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/colors"
	"github.com/stefan-muehlebach/ledgrid/conf"
)

const (
	defHost   = "raspi-3"
	defWidth  = 40
	defHeight = 10
)

var (
	width, height int
	gridSize      image.Point
	gridClient    ledgrid.GridClient
	ledGrid       *ledgrid.LedGrid
	animCtrl      *ledgrid.AnimationController
	canvas        *ledgrid.Canvas
)

//----------------------------------------------------------------------------

type ProgramList []LedGridProgram

type StartFunc func(ctx context.Context, c *ledgrid.Canvas)

func (pl *ProgramList) Add(name, group string, startFunc StartFunc) {
	prog := NewProgram(name, group, startFunc)
	*pl = append(*pl, prog)
}

type LedGridProgram interface {
	Name() string
	Group() string
	Start(ctx context.Context, c *ledgrid.Canvas)
	Stop()
}

func NewProgram(name, group string, startFunc StartFunc) LedGridProgram {
	return &simpleProgram{
		name:      name,
		group:     group,
		startFunc: startFunc,
	}
}

type simpleProgram struct {
	name, group string
	startFunc   StartFunc
	ctx         context.Context
	cancel      context.CancelFunc
}

func (p *simpleProgram) Name() string {
	return p.name
}

func (p *simpleProgram) Group() string {
	return p.group
}

func (p *simpleProgram) Start(ctx context.Context, c *ledgrid.Canvas) {
	p.ctx, p.cancel = context.WithCancel(ctx)
	p.startFunc(p.ctx, c)
}

func (p *simpleProgram) Stop() {
	fmt.Printf("Stop(): stopping context\n")
	p.cancel()
	fmt.Printf("Stop(): context is stopped\n")
}

// ---------------------------------------------------------------------------

func SignalHandler(timeout time.Duration) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	if timeout == 0 {
		timeout = time.Duration(math.MaxInt64)
	}
	timer := time.NewTimer(timeout)
	select {
	case <-sigChan:
	case <-timer.C:
	}
}

//----------------------------------------------------------------------------

var (
	programList ProgramList = make([]LedGridProgram, 0)
)

func main() {
	var host string
	var dataPort, rpcPort uint
	var useTCP bool
	var network string
	var progChar string
	var input string
	var ch byte
	var progId, prevProgId int
	// var runInteractive bool
	var progList string
	// var gR, gG, gB float64
	var modConf conf.ModuleConfig
	var timeout time.Duration
	var outFile string

	for i, prog := range programList {
		var id byte
		switch prog.(type) {
		case *simpleProgram:
			if i < 26 {
				id = byte('a' + i)
			} else {
				id = byte('A' + (i - 26))
			}
			progList += fmt.Sprintf("\n%c - %s", id, prog.Name())
		}
	}

	flag.IntVar(&width, "width", defWidth, "Width (for 'out' option only)")
	flag.IntVar(&height, "height", defHeight, "Height (for 'out' option only)")
	flag.StringVar(&outFile, "out", "", "Send all data to this file")

	flag.StringVar(&host, "host", defHost, "Controller hostname")
	flag.BoolVar(&useTCP, "tcp", false, "Use TCP for data")
	flag.UintVar(&dataPort, "data", ledgrid.DefDataPort, "Data Port")
	flag.UintVar(&rpcPort, "rpc", ledgrid.DefRPCPort, "RPC Port")
	flag.StringVar(&progChar, "prog", "", "Play one single program"+progList)
	flag.DurationVar(&timeout, "timeout", 0, "Timeout in non interactive mode")
	flag.Parse()

	StartProfiling()
	defer StopProfiling()

	if useTCP {
		network = "tcp"
	} else {
		network = "udp"
	}

	if outFile != "" {
		gridClient = ledgrid.NewFileSaveClient(outFile, conf.DefaultModuleConfig(image.Point{width, height}))
	} else {
		gridClient = ledgrid.NewNetGridClient(host, network, dataPort, rpcPort)
	}
	modConf = gridClient.ModuleConfig()
	ledGrid = ledgrid.NewLedGrid(gridClient, modConf)
	// gR, gG, gB = ledGrid.Client.Gamma()

	gridSize = ledGrid.Rect.Size()
	width = gridSize.X
	height = gridSize.Y

	canvas = ledGrid.Canvas(0)
	animCtrl = ledGrid.AnimCtrl

	ledGrid.StartRefresh()

	progId, prevProgId = -1, -1

	for {
		if len(progChar) > 0 {
            time.Sleep(500 * time.Millisecond)
			ch = progChar[0]
		} else {
			fmt.Printf("---------------------------------------------------------------------\n")
			fmt.Printf("  Program\n")
			fmt.Printf("---------------------------------------------------------------------\n")
			for i, prog := range programList {
				var id byte

				if i == progId {
					fmt.Printf("> ")
				} else {
					fmt.Printf("  ")
				}

				if i < 26 {
					id = byte('a' + i)
				} else {
					id = byte('A' + (i - 26))
				}

				fmt.Printf("[%c] %s\n", id, prog.Name())
			}
			fmt.Printf("---------------------------------------------------------------------\n")
			// fmt.Printf("  Gamma values: %.1f, %.1f, %.1f\n", gR, gG, gB)
			// fmt.Printf("   +/-: increase/decreases by 0.1\n")
			// fmt.Printf("---------------------------------------------------------------------\n")

			fmt.Printf("Enter a character (or '0' for quit): ")

			fmt.Scanf("%s\n", &input)
			ch = input[0]
		}

		if ch == '0' {
			break
		}

		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			id := -1
			if ch >= 'a' {
				id = int(ch - 'a')
			} else {
				id = int(ch - 'A' + 26)
			}
			if id < 0 || id >= len(programList) {
				break
			}
			progId = id

			if prevProgId != -1 {
				fmt.Printf("Program statistics:\n")
				fmt.Printf("  animation: %v\n", ledgrid.AnimCtrl.Stopwatch())
				fmt.Printf("  painting : %v\n", canvas.Stopwatch())
				fmt.Printf("  sending  : %v\n", ledGrid.Client.Stopwatch())
				programList[prevProgId].Stop()
			}
			ledGrid.Reset()
			ledgrid.AnimCtrl.Stopwatch().Reset()
			canvas.Stopwatch().Reset()
			ledGrid.Client.Stopwatch().Reset()
			programList[progId].Start(context.Background(), canvas)
			prevProgId = progId

			if len(progChar) > 0 {
				fmt.Printf("Quit by Ctrl-C\n")
				SignalHandler(timeout)
				programList[progId].Stop()
				break
			}
		}
	}

	ledgrid.AnimCtrl.Suspend()
	ledGrid.Clear(colors.Black)
	ledGrid.Show()
	ledGrid.Close()

	fmt.Printf("Program statistics:\n")
	fmt.Printf("  animation: %v\n", ledgrid.AnimCtrl.Stopwatch())
	fmt.Printf("  painting : %v\n", canvas.Stopwatch())
	fmt.Printf("  sending  : %v\n", ledGrid.Client.Stopwatch())
}
