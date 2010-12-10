// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package fightclub

// #include "readin.c"
import "C"

import (
    "fmt"
    "os"
)

func Open(toOpen string) {

    var x = C.ReadIn( C.CString(toOpen) )
    if x == 0 {
        fmt.Fprintln(os.Stderr, ":)\n")
    } else {
        fmt.Fprintln(os.Stderr, ":(\n")
    }
}

