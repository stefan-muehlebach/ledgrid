package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"os/signal"
	"time"

	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/conf"
)

const (
	defHost       = "raspi-3"
	defWidth      = 40
	defHeight     = 20
	defClientType = NetClient
	defBaud       = 2_000_000
	spiDevFile    = "/dev/spidev0.0"
	defWordFile   = "Faust.txt"
	asciLine      = "---------------------------------------------------------------------"
)

var (
	width, height int
	gridSize      image.Point
	gridClient    ledgrid.GridClient
	ledGrid       *ledgrid.LedGrid
	animCtrl      *ledgrid.AnimationController
	canvas        *ledgrid.Canvas
	hostName      string
	wordFile      string
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

type ClientType byte

const (
	NetClient ClientType = iota
	FileClient
	DirectClient
)

func (c ClientType) String() string {
	switch c {
	case NetClient:
		return "Net"
	case FileClient:
		return "File"
	case DirectClient:
		return "Direct"
	default:
		return ""
	}
}

func (c *ClientType) Set(s string) error {
	switch s {
	case "Net", "net":
		*c = NetClient
	case "File", "file":
		*c = FileClient
	case "Direct", "direct":
		*c = DirectClient
	default:
		return errors.New("Unknown client type")
	}
	return nil
}

//----------------------------------------------------------------------------

var (
	programList ProgramList = make([]LedGridProgram, 0)
	modConf     conf.ModuleConfig
)

func main() {
	var host string
	var dataPort, rpcPort uint
	var clientType ClientType = defClientType
	var customConfName string
	var progNum int
	var input int
	var progId, prevProgId int
	var baud int
	var progList string
	var timeout time.Duration
	var outFile string
	var ws2801 ledgrid.Displayer
	var err error

	for i, prog := range programList {
		switch prog.(type) {
		case *simpleProgram:
			progList += fmt.Sprintf("\n%02d - %s", i+1, prog.Name())
		}
	}
	flag.Var(&clientType, "type", "Type of client; 'net' (default), 'file' or " +
		"'direct'")
	flag.StringVar(&customConfName, "custom", "", "Use a non standard module " +
		"configuration (Type: 'file' or 'direct')")
	flag.IntVar(&width, "width", defWidth, "Width (Type: 'file' or 'direct')")
	flag.IntVar(&height, "height", defHeight, "Height (Type: 'file' or 'direct')")
	flag.StringVar(&host, "host", defHost, "Controller hostname (Type: 'net')")
	flag.UintVar(&dataPort, "tcp", ledgrid.DefTCPPort, "TCP Port (Type: 'net')")
	flag.UintVar(&rpcPort, "rpc", ledgrid.DefRPCPort, "RPC Port (Type: 'net')")
	flag.StringVar(&outFile, "out", "", "Send all data to this file (Type: 'file')")
	flag.IntVar(&baud, "baud", defBaud, "SPI baudrate in Hz (Type: 'direct')")
	flag.IntVar(&progNum, "prog", 0, "Play one single program" + progList)
	flag.DurationVar(&timeout, "timeout", 0, "Timeout in non interactive mode")
	flag.StringVar(&wordFile, "words", defWordFile, "File with space " +
		"separated words")
	flag.Parse()

	StartProfiling()
	defer StopProfiling()

	switch clientType {
	case NetClient:
		gridClient = ledgrid.NewNetGridClient(host, dataPort, rpcPort)
		hostName = gridClient.(*ledgrid.NetGridClient).Address()
		modConf = gridClient.ModuleConfig()
	case FileClient:
		if outFile == "" {
			log.Fatalf("Must specify 'out' when using File client type")
		}
		modConf = conf.DefaultModuleConfig(image.Point{width, height})
		gridClient = ledgrid.NewFileSaveClient(outFile, modConf)
	case DirectClient:
		modConf = conf.DefaultModuleConfig(image.Point{width, height})
		ws2801 = ledgrid.NewWS2801(spiDevFile, baud, modConf)
		gridClient = ledgrid.NewDirectGridClient(ws2801)
	default:
		log.Fatalf("Client type %d not defined (expected 0..2)")
	}
	ledGrid = ledgrid.NewLedGrid(gridClient, modConf)

	gridSize = ledGrid.Rect.Size()
	width = gridSize.X
	height = gridSize.Y

	canvas = ledGrid.Canvas(0)
	animCtrl = ledGrid.AnimCtrl

	ledGrid.StartRefresh()

	progId, prevProgId = -1, -1

	for {
		if progNum >= 1 && progNum <= len(programList) {
			time.Sleep(500 * time.Millisecond)
			input = progNum
		} else {
			groupName := ""
			for i, prog := range programList {
				if groupName != prog.Group() {
					groupName = prog.Group()
					fmt.Printf("%.*s %s\n", len(asciLine)-len(groupName)-1, asciLine, groupName)
				}
				if i == progId {
					fmt.Printf("> ")
				} else {
					fmt.Printf("  ")
				}
				fmt.Printf("[%02d] %s\n", i+1, prog.Name())
			}
			fmt.Println(asciLine)
			fmt.Printf("Enter a number (or '0' for quit): ")

			_, err = fmt.Scanf("%d", &input)
			if err != nil {
				log.Fatal(err)
			}
		}

		if input == 0 {
			break
		}
		if input >= 1 && input <= len(programList) {
			progId = input - 1
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

			if progNum >= 1 && progNum <= len(programList) {
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
