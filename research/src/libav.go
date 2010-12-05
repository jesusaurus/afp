package libav;

/* 
#include <stdlib.h>
#include "libav.h"
#include "libav.c"
*/
import "C"

import "unsafe"

type AVDecodeContext struct {
	Context *C.AVDecodeContext
}

func InitDecoding() {	
	C.init_decoding()
}

func PrepareDecoding(infile string, context *AVDecodeContext) {
	p := C.CString(infile)
	context.Context = new(C.AVDecodeContext)
	C.prepare_decoding(p, context.Context)
	C.free(unsafe.Pointer(p))
}

func DecodePacket(context AVDecodeContext) int {
	return int(C.decode_packet(context.Context))
}
