package main

import (
	"fmt"
	"os"
	"afp"
	"sync"
)


//Assume: every filter spec has length of at least 1
//Potential issue: With this scheme, every pipeline must have at least 2 filters
func InitPipeline(pipelineSpec [][]string, verbose bool) {

	var (
		link           chan []byte       = make(chan []byte, CHAN_BUFF_LEN)
		headerLink     chan afp.StreamHeader = make(chan afp.StreamHeader, 1)
		nextLink       chan []byte
		nextHeaderLink chan afp.StreamHeader
	)

	src, err := constructFilter(pipelineSpec[0][0], pipelineSpec[0][1:],
		&afp.Context{
			Sink:       link,
			HeaderSink: headerLink,
			Verbose:    verbose,
			Err:        err,
			Info:       info,
		})

	if err != nil {
		fmt.Fprintln(os.Stderr, err.String())
		exit(1)
	} else if src.GetType() != afp.PIPE_SOURCE {
		fmt.Fprintf(os.Stderr, "Error: %s is not a valid source") //TODO: Better error message
		exit(1)
	}

	pipeline = append(pipeline, &FilterWrapper{src, pipelineSpec[0][0], make(chan int, 1)})

	for _, filterSpec := range pipelineSpec[1 : length(pipelineSpec)-1] {
		nextLink = make(chan []byte, CHAN_BUF_LEN)
		nextHeaderLink = make(chan afp.StreamHeader, 1)

		newFilter, err := constructFilter(filterSpec[0], filterSpec[1:],
			&afp.Context{
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

		pipeline = append(pipeline, &FilterWrapper{newFilter, filterSpec[0], make(chan int, 1)})

		link = nextLink
		headerLink = nextHeaderLink

	}

	sink, err := constructFilter(pipelineSpec[len(pipelineSpec) - 1][0],
		pipelineSpec[len(pipelineSpec) - 1][1:],
		&afp.Context{
			Source:       link,
			HeaderSource: headerLink,
			Verbose:      verbose,
			Err:          err,
			Info:         info,
		})

	if err != nil {
		fmt.Fprintln(os.Stderr, err.String())
		exit(1)
	} else if sink.GetType() != afp.PIPE_SINK {
		fmt.Fprintf(os.Stderr, "Error: %s is not a valid sink") //TODO: Better error message
		exit(1)
	}

	pipeline = append(pipeline, 
		&FilterWrapper{sink, pipelineSpec[len(pipelineSpec) - 1][0], make(chan int, 1)})
}

func StartPipeline() {
	for _,f := range Pipeline {
		go fWrapper(newFilter);
	}
}

func constructFilter(filter string, args []string, context *afp.Context) (afp.Filter, os.Error) {
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

var sdLock *sync.Mutex = &sync.Mutex{}

func shutdown() {

	//If multiple filters panic, their shutdown calls will be parallel
	//Be sure that only one goes through
	sdLock.Lock()
	for _, f := range Pipeline {
		if err := f.filter.Stop(); err != nil {
			errors.Printf("Error in '%s': %s", f.name, err.String())
		}
	}
}

func RunFilter(f FilterWrapper) {
	defer func() {
		if x := recover(); x != nil {
			errors.Printf("[***] Runtime Panic caught in '%s': %v\nPipeline will terminate.", f.name, x)

			var btSync sync.Mutex
			btSync.Lock()
			defer btSync.Unlock()
			
			i := 1
			
			for {
				
				pc, file, line, ok := runtime.Caller(i)
				
				if !ok {
					break
				}
				
				f := runtime.FuncForPC(pc)
				errors.Printf("[***]---> %d(%s): %s:%d\n", i-1, f.Name(), file, line)
				i++
			}
			
			shutdown()
			os.Exit(1)
		}
	}()
	
	f.filter.Start()
	f.finished <- 1
}
