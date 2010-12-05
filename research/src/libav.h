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

/**
 * General information about a stream
 *
 * FIXME: elaborate
 */

typedef struct {
	int32_t Channels;
	int32_t Sample_size;
	int32_t Sample_rate;
	int64_t Content_length;
	int32_t Frame_size; 
} AVStreamInfo;

/**
 * A bunch of state to be passed around for the duration of a decode
 *
 * FIXME: elaborate
 */

typedef struct {
	AVCodec *Codec;
	AVFormatContext *Fctx;
	AVCodecContext *Cctx;
	int16_t *Outbuf;
	int32_t Buf_size, Buf_len;
	int8_t first_frame_used;
	AVStreamInfo Info;
	AVPacket Packet;
} AVDecodeContext;

/**
 * modified from libavcodec/mpegaudio.h
 *
 * fast header check for resync
 */
static inline int ff_mpa_check_header(uint32_t header){
    /* header */
    if ((header & 0xffe00000) != 0xffe00000)
        return -1;
    /* layer check */
    if ((header & (3<<17)) == 0)
        return -1;
    /* bit rate */
    if ((header & (0xf<<12)) == 0xf<<12)
        return -1;
    /* frequency */
    if ((header & (3<<10)) == 3<<10)
        return -1;
    return 0;
}

/**
 * modified from libavutil/intreadwrite.h
 *
 * does "fun" byt order conversions
 */

#ifndef AV_RB32
#   define AV_RB32(x)                           \
    ((((const uint8_t*)(x))[0] << 24) |         \
     (((const uint8_t*)(x))[1] << 16) |         \
     (((const uint8_t*)(x))[2] <<  8) |         \
      ((const uint8_t*)(x))[3])
#endif

#endif
