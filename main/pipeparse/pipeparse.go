package main

import (
	"io/ioutil"
	"fmt"
	"strings"
	"os"
)

type ParsePipelineResult struct {
	StageArgs [][]string
	Flags *FlagParserType
}

func ParsePipeline(args []string) *ParsePipelineResult {
	stageArgs = make([][]string, 30)
	mainArgs = make([]string, 30)
	int stagesStart
	for i, arg := range args[1:] {
		if arg.startswith("-") { // change to Go equivalent


	result = make([][]string, 30)
	for i, arg := range args[1:] {

	append(StageArgs, )
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