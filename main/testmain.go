// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package main

import (
	//	"afp"
	//	"os/signal"
	"log"
	"os"
)

var pipespec [][]string = [][]string{{"libavsource", "-i", "/tmp/test.mp3"}, {"stdoutsink"}}

const CHAN_BUF_LEN = 16

var (
	Pipeline []*FilterWrapper = make([]*FilterWrapper, 0, 100)

	errors  *log.Logger = log.New(os.Stderr, "[E] ", log.Ltime)
	info    *log.Logger = log.New(os.Stderr, "[I] ", log.Ltime)
	verbose bool        = false
)


func main() {
	InitPipeline(pipespec, verbose)
	StartPipeline()

	for _, filter := range Pipeline {
		<-filter.finished
		if verbose {
			info.Printf("Filter '%s' finished", filter.name)
		}
	}
}
