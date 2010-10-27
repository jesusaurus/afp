package main;

import (
	"libav"
)

func main() {
	var context libav.AVDecodeContext
	libav.InitDecoding()
	libav.PrepareDecoding("/tmp/test.mp3", &context)
	for l := libav.DecodePacket(context); l > 0; {
		l = libav.DecodePacket(context)
	}
	
}
