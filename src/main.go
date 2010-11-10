package main;

import (
	"libav"
	"unsafe"
)

func main() {
	var context libav.AVDecodeContext
	libav.InitDecoding()
	libav.PrepareDecoding("/tmp/test.mp3", &context)
	for l := libav.DecodePacket(context); l > 0; {
		l = libav.DecodePacket(context)
		sample := (*(*[1 << 31 - 1]int16)(unsafe.Pointer(context.Context.Outbuf)))[:l/2]

		for j,s := range sample {
			if (j % 32 == 0) {
				width := 236
				half := int16(width/2) 
				offset := int16(float(s) * float(width) / 65535.0)

				for i := 0; i < int(half + offset); i++ {
					print(" ")
				}
				println("#")
			}
		}
	}
	
}
