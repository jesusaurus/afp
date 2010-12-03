package main

import (
	"io/ioutil"
	"fmt"
	"strings"
	"os"
)

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
