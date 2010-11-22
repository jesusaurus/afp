package afp

//This file is still just a sketch, does not compile

import (
	"fmt"
	"os"
)

var pipeline []Filter = make([]Filter, 0, 100)
var CHAN_BUFF_LEN int = 16 //Move this elsewhere

//Assume: every filter spec has length of at least 1
//Potential issue: With this scheme, every pipeline must have at least 2 filters
func InitPipeline(pipelineSpec [][]string, verbose bool, err, info log.Logger) {

	var (
		link           chan []byte       = make(chan []byte, CHAN_BUFF_LEN)
		headerLink     chan StreamHeader = make(chan StreamHeader, 1)
		nextLink       chan []byte
		nextHeaderLink chan StreamHeader
	)

	src, err := constructFilter(pipelineSpec[0][0], pipelineSpec[0][1:],
		&Context{
			Sink:       link,
			HeaderSink: headerLink,
			Verbose:    verbose,
			Err:        err,
			Info:       info,
		})

	if err != nil {
		fmt.Fprintln(os.Stderr, err.String())
		exit(1)
	} else if src.GetType() != PIPE_SOURCE {
		fmt.Fprintf(os.Stderr, "Error: %s is not a valid source") //TODO: Better error message
		exit(1)
	}

	for _, filterSpec := range pipelineSpec[1 : length(pipelineSpec)-1] {
		nextLink = make(chan []byte, CHAN_BUF_LEN)
		nextHeaderLink = make(chan StreamHeader, 1)

		newFilter, err := constructFilter(filterSpec[0], filterSpec[1:],
			&Context{
				Source:       link,
				HeaderSource: headerLink,
				Sink:         nextLink,
				HeaderSink:   nextHeaderLink,
				Verbose:      verbose,
				Err:          err,
				Info:         info,
			})

		if err != nil {
			fmt.Fprintln(os.Stderr, err.String())
			exit(1)
		}

		pipeline = append(pipeline, newFilter)

		link = nextLink
		headerLink = nextHeaderLink
	}

	sink, err := constructFilter(filterSpec[0], filterSpec[1:],
		&Context{
			Source:       link,
			HeaderSource: headerLink,
			Verbose:      verbose,
			Err:          err,
			Info:         info,
		})

	if err != nil {
		fmt.Fprintln(os.Stderr, err.String())
		exit(1)
	} else if src.GetType() != PIPE_SINK {
		fmt.Fprintf(os.Stderr, "Error: %s is not a valid source") //TODO: Better error message
		exit(1)
	}
}

func constructFilter(filter string, args []string, context *Context) (Filter, os.Error) {
	//Is the filter in the list of known filters?
	ctor, ok := filters[filterSpec[0]]
	if !ok {
		return nil, os.NewError(fmt.Sprintf("Error: %s: filter not found.", filterSpec[0]))
	}

	newFilter := ctor()
	if newFilter == nil {
		return nil, os.NewError(fmt.Sprintf("Error: %s: Attempt to create filter failed.", filterSpec[0]))
	}

	err := newFilter.Init(context, args)
	if err != nil {
		return nil, err
	}

	return newFilter, nil
}

func StartPipeline() {
	for _, f := range pipeline {
		f.Start()
	}
}
