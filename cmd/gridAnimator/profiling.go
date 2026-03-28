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
	cpuProfFile, memProfFile, traceFile string
	fhCpu, fhTrace             *os.File
)

func init() {
	flag.StringVar(&cpuProfFile, "cpuprofile", "", "write cpu profile to 'file'")
	flag.StringVar(&memProfFile, "memprofile", "", "write memory profile to 'file'")
	flag.StringVar(&traceFile, "trace", "", "write trace data to 'file'")
}

func StartProfiling() {
	if cpuProfFile != "" {
		fhCpu, err := os.Create(cpuProfFile)
		if err != nil {
			log.Fatal("couldn't create cpu profile: ", err)
		}
		if err := pprof.StartCPUProfile(fhCpu); err != nil {
			log.Fatal("couldn't start cpu profiling: ", err)
		}
	}

	if traceFile != "" {
		fhTrace, err := os.Create(traceFile)
		if err != nil {
			log.Fatal("couldn't create tracefile: ", err)
		}
		if err := trace.Start(fhTrace); err != nil {
			log.Fatal("couldn't start trace: ", err)
		}
	}
}

func StopProfiling() {
	if cpuProfFile != "" {
		pprof.StopCPUProfile()
		fhCpu.Close()
	}

	if memProfFile != "" {
		fh, err := os.Create(memProfFile)
		if err != nil {
			log.Fatal("couldn't create memory profile: ", err)
		}
		defer fh.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(fh); err != nil {
			log.Fatal("couldn't write memory profile: ", err)
		}
	}

	if traceFile != "" {
		trace.Stop()
		fhTrace.Close()
	}
}
