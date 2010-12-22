// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package main

import (
	"fmt"
	"os"
	"afp"
)

type FilterWrapper struct {
	filter   afp.Filter
	ctx      *afp.Context
	name     string
	finished chan int
}

//InitPipeline takes a parsed pipeline spec and attempts to create the appropriate 
//filters, wrap them, and insert them into the pipeline.  It checks that the first 
//filter is a source, the last a sink, and that any between are links
func InitPipeline(pipelineSpec [][]string, verbose bool) {
	if len(pipelineSpec) < 2 {
		errors.Println("Pipeline specification must have at least a Source and Sink")
		os.Exit(1)
	}

	var (
		link           chan [][]float32      = make(chan [][]float32, afp.CHAN_BUF_LEN)
		headerLink     chan afp.StreamHeader = make(chan afp.StreamHeader, 1)
		nextLink       chan [][]float32
		nextHeaderLink chan afp.StreamHeader
		ctx            *afp.Context
	)

	ctx = &afp.Context{
		Sink:       link,
		HeaderSink: headerLink,
		Verbose:    verbose,
		Err:        errors,
		Info:       info,
	}

	src, err := constructFilter(pipelineSpec[0][0], pipelineSpec[0][:], ctx)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.String())
		os.Exit(1)
	} else if src.GetType() != afp.PIPE_SOURCE {
		fmt.Fprintf(os.Stderr, "Error: %s is not a valid source") //TODO: Better error message
		os.Exit(1)
	}

	//In case of an error elsewhere, or a external signal, shutdown may be called
	//in the middle of constructing the pipeline.  Be sure that we do not modify
	//the pipeline at the same time.
	pipelineLock.Lock()
	Pipeline = append(Pipeline, &FilterWrapper{src, ctx, pipelineSpec[0][0], make(chan int, 1)})
	pipelineLock.Unlock()

	for _, filterSpec := range pipelineSpec[1 : len(pipelineSpec)-1] {
		nextLink = make(chan [][]float32, afp.CHAN_BUF_LEN)
		nextHeaderLink = make(chan afp.StreamHeader, 1)

		ctx = &afp.Context{
			Source:       link,
			HeaderSource: headerLink,
			Sink:         nextLink,
			HeaderSink:   nextHeaderLink,
			Verbose:      verbose,
			Err:          errors,
			Info:         info,
		}

		newFilter, err := constructFilter(filterSpec[0], filterSpec[:], ctx)

		if err != nil {
			fmt.Fprintln(os.Stderr, err.String())
			os.Exit(1)
		}

		pipelineLock.Lock()
		Pipeline = append(Pipeline, &FilterWrapper{newFilter, ctx, filterSpec[0], make(chan int, 1)})
		pipelineLock.Unlock()

		link = nextLink
		headerLink = nextHeaderLink
	}

	ctx = &afp.Context{
		Source:       link,
		HeaderSource: headerLink,
		Verbose:      verbose,
		Err:          errors,
		Info:         info,
	}

	sink, err := constructFilter(pipelineSpec[len(pipelineSpec)-1][0],
		pipelineSpec[len(pipelineSpec)-1][:], ctx)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.String())
		os.Exit(1)
	} else if sink.GetType() != afp.PIPE_SINK {
		fmt.Fprintf(os.Stderr, "Error: %s is not a valid sink") //TODO: Better error message
		os.Exit(1)
	}

	pipelineLock.Lock()
	Pipeline = append(Pipeline,
		&FilterWrapper{sink, ctx, pipelineSpec[len(pipelineSpec)-1][0], make(chan int, 1)})
	pipelineLock.Unlock()
}

func StartPipeline() {
	for _, f := range Pipeline {
		go RunFilter(f)
	}
}

func constructFilter(name string, args []string, context *afp.Context) (afp.Filter, os.Error) {
	//Is the filter in the list of known filters?
	ctor, ok := filters[name]
	if !ok {
		return nil, os.NewError(fmt.Sprintf("Error: %s: filter not found.", name))
	}

	newFilter := ctor()
	if newFilter == nil {
		return nil, os.NewError(fmt.Sprintf("Error: %s: Attempt to create filter failed.", name))
	}

	err := newFilter.Init(context, args)
	if err != nil {
		return nil, err
	}

	return newFilter, nil
}


func shutdown() {
	//If multiple filters panic, their shutdown calls will be parallel
	//Be sure that only one goes through
	pipelineLock.Lock()
	for _, f := range Pipeline {
		//This is ugly.  We want to catch any panics thrown in Stop methods
		//and continue shutting down the other filters
		func() {
			if !debugging {
				defer func() {
					if x := recover(); x != nil {
						errors.Printf("Panic caught in '%s': %v", f.name, x)
					}
				}()
			}
			if err := f.filter.Stop(); err != nil {
				errors.Printf("Error in '%s': %s", f.name, err.String())
			}
		}()
	}
	os.Exit(1)
}

func RunFilter(f *FilterWrapper) {
	if !debugging {
		defer func() {
			if x := recover(); x != nil {
				errors.Printf("Runtime Panic caught in '%s': %v\nPipeline will terminate.", f.name, x)
				shutdown()
				os.Exit(1)
			}
		}()
	}

	f.filter.Start()

	if f.ctx.Sink != nil && !closed(f.ctx.Sink) {
		close(f.ctx.Sink)
	}

	f.finished <- 1
}
