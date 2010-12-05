// Copyright (c) 2010 Go Fightclub Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

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
/*		println("Packet decode: ", l, " bytes")*/
		l = libav.DecodePacket(context)
		sample := (*(*[1 << 31 - 1]int16)(unsafe.Pointer(context.Context.Outbuf)))[:l/2]

/*		for i := 0; i < l/2; i++ {
			print(sample[i], " ")
		}
*/		for j,s := range sample {
			if (j % 2 == 0) {
				widt := 136
				half := int16(width/2)
				offset := int16(float(s) * float(width) / 65535.0)

				for i := 0; i < int(half + offset); i++ {
					print(" ")
				}
				println("#")
			}
		}
/*		println("Decoded ", l, " bytes, oh joy!")
		println("Context shows a packet at: ", []int16(reflect.MakeSlice(unsafe.Pointer(context.Context.Outbuf))))*/
	}

}
