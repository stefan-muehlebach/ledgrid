package main

import (
	"flag"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
)

var (
	cpuprofile, memprofile, tracefile string
	fhCpu, fhTrace             *os.File
)

func init() {
	flag.StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile to 'file'")
	flag.StringVar(&memprofile, "memprofile", "", "write memory profile to 'file'")
	flag.StringVar(&tracefile, "trace", "", "write trace data to 'file'")
}

func StartProfiling() {
	if cpuprofile != "" {
		fhCpu, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal("couldn't create cpu profile: ", err)
		}
		if err := pprof.StartCPUProfile(fhCpu); err != nil {
			log.Fatal("couldn't start cpu profiling: ", err)
		}
	}

	if tracefile != "" {
		fhTrace, err := os.Create(tracefile)
		if err != nil {
			log.Fatal("couldn't create tracefile: ", err)
		}
		if err := trace.Start(fhTrace); err != nil {
			log.Fatal("couldn't start trace: ", err)
		}
	}
}

func StopProfiling() {
	if cpuprofile != "" {
		pprof.StopCPUProfile()
		fhCpu.Close()
	}

	if memprofile != "" {
		fh, err := os.Create(memprofile)
		if err != nil {
			log.Fatal("couldn't create memory profile: ", err)
		}
		defer fh.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(fh); err != nil {
			log.Fatal("couldn't write memory profile: ", err)
		}
	}

	if tracefile != "" {
		trace.Stop()
		fhTrace.Close()
	}
}
