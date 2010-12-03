package main

import (
	"./filters/delay"
	"./filters/fexec"
	"./filters/ospipe"
	"./filters/demo"
)

var filters map[string]func() Filter = map[string]func() Filter {
	"delay"  : delay.NewFilter,
	"exec"   : fexec.NewFilter,
	"stdin"  : ospipe.StdinSource,
	"stdout" : ospipe.StdoutSink,
	"nop"    : demo.NopFilter,
}