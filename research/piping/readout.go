package fightclub

// #include "readout.c"
import "C"

import (
    "fmt"
    "os"
)

func Write(toWrite string) {

    var x = C.WriteOut( C.CString(toWrite) )
    if x == 0 {
        fmt.Fprintln(os.Stderr, ":)\n")
    } else {
        fmt.Fprintln(os.Stderr, ":(\n")
    }
}

