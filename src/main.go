package main;

import (
	"libav"
	"unsafe"
	"os"
)

const (
	samplesPerSecond = 44100
	samplesPerMillisecond = samplesPerSecond / 1000
	delayTimeInMs = 150
	channels = 2
	bytesPerSample = 2
)

func main() {
	var context libav.AVDecodeContext
/*	var t0 int = 0*/
	var currBuffer = 0
	var buffer [2][]int16
	var extraSamples int = samplesPerMillisecond * delayTimeInMs * channels

	buffer[0] = make([]int16, extraSamples, extraSamples)
	
	libav.InitDecoding()
	libav.PrepareDecoding("/tmp/test.mp3", &context)
	
	for l := libav.DecodePacket(context); l > 0; {
		l = libav.DecodePacket(context)
		numberOfSamples := l / bytesPerSample
		decodedSamples := (*(*[1 << 31 - 1]int16)(unsafe.Pointer(context.Context.Outbuf)))[:numberOfSamples]

/*		os.Stdout.Write((*(*[1 << 31 - 1]uint8)(unsafe.Pointer(context.Context.Outbuf)))[:(numberOfSamples*bytesPerSample)])*/
		buffer[1 - currBuffer] = append(buffer[currBuffer], decodedSamples...)
		currBuffer = 1 - currBuffer;
	}
	
	for t0,_ := range buffer[currBuffer] {
		if t0 > extraSamples {
			buffer[1 - currBuffer][t0] = buffer[currBuffer][t0] + (buffer[currBuffer][t0 - extraSamples] / 2)
		}
	}
	currBuffer = 1 - currBuffer;
	
/*	println(t0, len(buffer[currBuffer]))*/

/*	for _,s := range buffer[currBuffer] {
		os.Stdout.Write()
	}
*/	os.Stdout.Write(*(*[]uint8)(unsafe.Pointer(&buffer[currBuffer])));
}
