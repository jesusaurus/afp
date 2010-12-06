// Copyright (c) 2010 Go Fightclub Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package main

import (
	"io/ioutil"
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
			if len(currentStage) < 1 {
				errors.Println("Empty pipeline stages are not allowed")
				os.Exit(1)
			}
			
			stages = append(stages, currentStage)
			currentStage = make([]string, 0, 10)
		} else {
			currentStage = append(currentStage, arg)
		}
	}
	if len(currentStage) < 1 {
		errors.Println("Empty pipeline stages are not allowed")
		os.Exit(1)
	}
	stages = append(stages, currentStage)

	return mainArgs, stages
}

//If the pipeline is being pulled from a file, we'll need to split it
func GetPipelineFromFile(path string) ([]string, os.Error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	//The file may have newlines or odd whitespace patterns
	//Replace them by single spaces before we split
	strSpec := regexp.MustCompile(`[ \t\n\r]+`).ReplaceAllString(string(bytes), " ")

	return strings.Split(strSpec, " ", -1), nil
}
