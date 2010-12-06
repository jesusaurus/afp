// Copyright (c) 2010 Go Fightclub Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package main

import (
	"afp"
	"afp/filters/null"
	"afp/filters/fexec"
	"afp/filters/libav"
    "afp/filters/alsa"
)

var filters map[string]func() afp.Filter = map[string]func() afp.Filter {
	"execsink"		: fexec.NewExecSink,
	"execlink"		: fexec.NewExecLink,
	"execsource"	: fexec.NewExecSource,*/
	"nullsource"	: null.NewNullSource,
	"nulllink"		: null.NewNullLink,
	"nullsink"		: null.NewNullSink,
	"stdoutsink"	: stdout.NewStdoutSink,
	"libavsource"	: libavfilter.NewLibAVSource,
    "alsasource"    : alsa.NewAlsaSource,
    "alsasink"      : alsa.NewAlsaSink,
}
