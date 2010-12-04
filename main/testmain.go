package main

import (
	"afp"
	"os/signal"
	)

type FilterWrapper {
	filter afp.Filter
	name string
	finished chan int
}

var pipespec [][]string = [][]string{{"source"},{"sink"}}

var (
	Pipeline []*FilterWrapper = make(*FilterWrapper, 0, 100)
	
	errors log.Logger = log.New(os.Stderr, "[E] ", log.Ltime)
	info log.Logger = log.New(os.Stderr, "[I] ", log.Ltime)
	verbose bool = false
	)


func main() {
	InitPipeline(pipespec)
	StartPipeline()

	for _, filter := range Pipeline {
		<-filter.finished
		if verbose {
			info.Printf("Filter '%s' finished", filter.name)
		}
	}
}

