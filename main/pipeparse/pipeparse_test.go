package pipeparse_test

import (
	"testing"
	"main"
)

func TestEmptyPipeline(t *testing.T) {
	args = []string{"apl"}
	parsed := parsepipe.ParsePipeline(args)
	parsed.Flags.Parse()
	if len(parsed.StageArgs) != 0 {
		t.Error("Empty pipeline returned %d stages", len(parsed.StageArgs))
	}
}

func TestAll(t *testing.T) {
	args = []string{"apl", "-v", "filesrc", "-t", "flac", "file", "!", "filesink", "file"}
	mainArgs, stages := parsepipe.ParsePipeline(args)
}
// 	parsed.Flags.Bool("verbose"
// 	parsed.Flags.Parse()