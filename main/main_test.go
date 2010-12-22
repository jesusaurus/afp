// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package main_test

import (
	"testing"
	// 	"main"
)


func TestEmptyPipeline(t *testing.T) {
	args := []string{"apl"}
	mainArgs, stages := main.ParsePipeline(args)
	if len(stages) != 0 {
		t.Errorf("Empty pipeline returned %d stages",
			len(stages))
	}
	if len(mainArgs) != 0 {
		t.Errorf("Empty pipeline returned %d main flags",
			len(mainArgs))
	}
}

func TestAll(t *testing.T) {
	args := []string{"apl", "-v", "filesrc", "-t", "flac",
		"file", "!", "filesink", "file"}
	targetStages := [][]string{
		{"filesrc", "-t", "flac", "file"},
		{"filesink", "file"}}
	mainArgs, stages := main.ParsePipeline(args)
	if mainArgs != []string{"-v"} {
		t.Error("main arguments incorrect")
	}

	for i, stage := range stages {
		if i > len(stages) {
			t.Error("stages != target stages")
		}
		for j, arg := range stage {
			if j > len(stage) {
				t.Error("stages != target stages")
			}
			if arg != targetStages[i][j] {
				t.Error("stages != target stages")
			}
		}
	}
}
