// Copyright (c) 2010 Go Fightclub Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package main

import (
	"log"
	"os"
	"afp/flags"
)

const CHAN_BUF_LEN = 16

var (
	Pipeline []*FilterWrapper = make([]*FilterWrapper, 0, 100)
	errors *log.Logger = log.New(os.Stderr, "[E] ", log.Ltime)
	verbose bool
	)

func main() {
	mainArgs, pipespec := ParsePipeline(os.Args)
	mainFlags := flags.FlagParser(mainArgs)	
	mainFlags.BoolVar(&verbose, "v", false, "Verbose output")
	
	InitPipeline(pipespec, verbose)
	StartPipeline()

	for _, filter := range Pipeline {
		<-filter.finished
		if verbose {
			info.Printf("Filter '%s' finished", filter.name)
		}
	}
}
