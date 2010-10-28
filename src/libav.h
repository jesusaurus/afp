/*
 * libav.h
 * Eric O'Connell
 *
 * Reduce the surface area of libav* to something more suited to our purposes
 *
 */

#ifndef __LIBAV_H__
#define __LIBAV_H__

#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>

typedef struct {
	AVCodec *Codec;
	AVFormatContext *Fctx;
	AVCodecContext *Cctx;
	uint8_t *Outbuf;
	long Buf_size, Buf_len;
	AVPacket Packet;
} AVDecodeContext;

#endif
