// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package main

import (
	"log"
	"os"
	"os/signal"
	"afp/flags"
	"sync"
	"syscall"
)

const CHAN_BUF_LEN = 16

var (
	Pipeline []*FilterWrapper = make([]*FilterWrapper, 0, 100)
	pipelineLock *sync.Mutex = &sync.Mutex{}


	errors   *log.Logger      = log.New(os.Stderr, "[E] ", log.Ltime)
	info     *log.Logger      = log.New(os.Stderr, "[I] ", log.Ltime)
	verbose  bool
	specFile string
)

func Init() {
	go SigHandler()
}

func main() {
	mainArgs, pipespec := ParsePipeline(os.Args)
	mainFlags := flags.FlagParser(mainArgs)
	mainFlags.BoolVar(&verbose, "v", false, "Verbose output")
	mainFlags.StringVar(&specFile, "f", "",
		"Pull pipeline spec from a file rather than command line")
	mainFlags.Parse()

	if specFile != "" {
		rawPipe, err := GetPipelineFromFile(specFile)
		if err != nil {
			errors.Println(err.String())
			os.Exit(1)
		}
		_, pipespec = ParsePipeline(rawPipe)
	}
	InitPipeline(pipespec, verbose)
	StartPipeline()

	for _, filter := range Pipeline {
		<-filter.finished
		if verbose {
			info.Printf("Filter '%s' finished", filter.name)
		}
	}
}

func SigHandler() {
	for sig := range signal.Incoming {
		usig, ok := sig.(signal.UnixSignal)

		if !ok {
			shutdown()
			errors.Printf("Process received unknown signal: %s" + sig.String())
			os.Exit(1)
		}

		switch usig {
		case syscall.SIGABRT, syscall.SIGFPE,  syscall.SIGILL, 
			 syscall.SIGINT,  syscall.SIGKILL, syscall.SIGQUIT, 
			 syscall.SIGSEGV, syscall.SIGSTOP, syscall.SIGTERM,
		     syscall.SIGTSTP :
			errors.Printf("Received signal: %v. Pipeline will terminate.", usig)
			shutdown()
			os.Exit(1)
		}
	}
}
