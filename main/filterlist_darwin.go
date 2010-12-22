// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package main


//In order to rebuild afp to include your filter, import it below..
import (
	"afp"
	"afp/filters/null"
	"afp/filters/fexec"
	"afp/filters/libav"
	"afp/filters/stdout"
	"afp/filters/tone"
	"afp/filters/delay"
	"afp/filters/portaudio"
	"afp/filters/distort"
)

//And add a key : value pair to the map below, where the key is a string 
//by which your filter should be invoked, and the value is a function
//which constructs a ready to use instance of your filter.
var filters map[string]func() afp.Filter = map[string]func() afp.Filter{
	"execsink":    fexec.NewExecSink,
	"execlink":    fexec.NewExecLink,
	"execsource":  fexec.NewExecSource,
	"nullsource":  null.NewNullSource,
	"nulllink":    null.NewNullLink,
	"nullsink":    null.NewNullSink,
	"stdoutsink":  stdout.NewStdoutSink,
	"libavsource": libavfilter.NewLibAVSource,
	"tonesource":  tonefilter.NewToneSource,
	"pasink":      portaudio.NewPASink,
	"delay":       delay.NewDelayFilter,
	"distort":     distort.NewFilter,
}
