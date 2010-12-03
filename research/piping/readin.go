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

