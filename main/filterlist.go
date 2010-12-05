// Copyright (c) 2010 Go Fightclub Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package main

import (
	"afp"
	"afp/filters/null"
	"afp/filters/fexec"
//	"./filters/delay"
//	"./filters/fexec"
//	"./filters/ospipe"
//	"./filters/demo"
)

var filters map[string]func() afp.Filter = map[string]func() afp.Filter {
	"exec"   : fexec.NewFilter,
	"nullsource" : null.NewNullSource,
	"nulllink"   : null.NewNullLink,
	"nullsink" : null.NewNullSink,
//	"delay"  : delay.NewFilter,
//	"stdin"  : ospipe.StdinSource,
//	"stdout" : ospipe.StdoutSink,
//	"nop"    : demo.NopFilter,
}