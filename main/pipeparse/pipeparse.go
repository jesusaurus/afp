// Copyright (c) 2010 Go Fightclub Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package main

import (
	"io/ioutil"
	"fmt"
	"strings"
	"os"
	"regexp"
)

const INITIAL_STAGE_SIZE = 30


func ParsePipeline(args []string) ([]string, [][]string) {
	var pipelineStart int
	mainArgs := make([]string, 0, 3)
	stages := make([][]string, 0, INITIAL_STAGE_SIZE)
	for i, arg := range args[1:] {
		if !strings.HasPrefix(arg, "-") {
			pipelineStart = i + 1
			break
		} else {
			mainArgs = append(mainArgs, arg)
		}
	}

	currentStage := make([]string, 0, 10)
	for _, arg := range args[pipelineStart:] {
		if arg == "!" {
			stages = append(stages, currentStage)
			currentStage = make([]string, 0, 10)
		} else {
			currentStage = append(currentStage, arg)
		}
	}

	return mainArgs, stages
}

//If the pipeline is being pulled from a file, we'll need to split it
func getSpecFromFile(path string) []string {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open '%s': %s", path, err.String())
		os.Exit(1)
	}

	//The file may have newlines or odd whitespace patterns
	//Replace them by single spaces before we split
	strSpec := regexp.MustCompile(`[ \t\n\r]+`).ReplaceAllString(string(bytes), " ")

	return strings.Split(strSpec, " ", -1)
}