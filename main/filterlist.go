// Copyright (c) 2010 Go Fightclub Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package main

import (
	"afp"
//	"./filters/delay"
//	"./filters/fexec"
//	"./filters/ospipe"
//	"./filters/demo"
	"../filters/null"
)

var filters map[string]func() afp.Filter = map[string]func() afp.Filter {
//	"delay"  : delay.NewFilter,
//	"exec"   : fexec.NewFilter,
//	"stdin"  : ospipe.StdinSource,
//	"stdout" : ospipe.StdoutSink,
//	"nop"    : demo.NopFilter,
	"nullsource" : null.NewNullSource,
	"nulllink"   : null.NewNullLink,
	"nullsink" : null.NewNullSink,
}