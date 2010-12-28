package libav

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

type AVStreamInfo struct {
	Channels      int32
	SampleSize    int32
	SampleRate    int32
	ContentLength int64
	FrameSize     int32
}

func InitDecoding() {
	C.init_decoding()
}

func PrepareDecoding(infile string, context *AVDecodeContext) int {
	p := C.CString(infile)

	context.Context = new(C.AVDecodeContext)
	err := int(C.prepare_decoding(p, context.Context))

	C.free(unsafe.Pointer(p))

	return err
}

func DecodePacket(context AVDecodeContext) int {
	var ignore C.int
	return int(C.decode_packet(context.Context, &ignore))
}

func StreamInfo(context AVDecodeContext) AVStreamInfo {
	var info AVStreamInfo
	sourceStreamInfo := C.get_stream_info(context.Context)

	info.Channels = int32(sourceStreamInfo.Channels)
	info.SampleSize = int32(sourceStreamInfo.Sample_size)
	info.SampleRate = int32(sourceStreamInfo.Sample_rate)
	info.ContentLength = int64(sourceStreamInfo.Content_length)
	info.FrameSize = int32(sourceStreamInfo.Frame_size)

	return info
}
