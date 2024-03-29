// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

/**
 * libav.c
 * (c) Eric O'Connell 2010
 *
 * provide a minimal (possibly braindead) interface into audio decoding using
 * libavformat to demux input files & libavcodec to produce PCM data
 * 
 * with help from http://qtdvd.com/guides/ffmpeg.html for ffmpeg .5.1 compatibility
 * originally based on a combination of ffmpeg/tools/pktdumper.c and ffmpeg/libavcodec/api-sample.c
 *
 */


// #define __FFMPEG5__
// #undef AVCORE_SAMPLEFMT_H

#include "libav.h"
#include <string.h>

#define MPA_DECODE_HEADER \
    int frame_size; \
    int error_protection; \
    int layer; \
    int sample_rate; \
    int sample_rate_index; /* between 0 and 8 */ \
    int bit_rate; \
    int nb_channels; \
    int mode; \
    int mode_ext; \
    int lsf;

typedef struct MPADecodeHeader {
  MPA_DECODE_HEADER
} MPADecodeHeader;

void init_decoding(void);
int prepare_decoding(char *filename, AVDecodeContext *context);
int decode_packet(AVDecodeContext *context, int *out_size);
AVStreamInfo get_stream_info(AVDecodeContext *context);
int is_id3_tag(AVDecodeContext *context);
int ff_mpa_decode_header(AVCodecContext *avctx, uint32_t head, int *sample_rate, int *channels, int *frame_size, int *bit_rate);
int ff_mpegaudio_decode_header(MPADecodeHeader *s, uint32_t header);
void avcodec_string(char *buf, int buf_size, AVCodecContext *enc, int encode);

/**
 * initialize all of the fun libavformat & libavdecode stuff
 *
 */
void init_decoding(void) {
	/* must be called before using avformat lib */
    av_register_all();

    /* must be called before using avcodec lib */
    avcodec_init();

    /* register all the codecs */
    avcodec_register_all();

	/* hush! */
	av_log_set_level(AV_LOG_QUIET);
}


/**
 * set up a decoding context
 *
 * @param filename the filename to decode
 * @param context a pointer to an AVDecodeContext, will be initialized
 * @return 0 on success, -1 on error
 */
int prepare_decoding(char *filename, AVDecodeContext *context) {
	int err, len;

	/* initialize storage for the AVPacket */
    av_init_packet(&context->Packet);

	/* use libavformat to open the source file; initializes AVFormatContext */
    if ((err = av_open_input_file(&context->Fctx, filename, NULL, 0, NULL)) < 0) {
        fprintf(stderr, "av_open_input_file: file: %s, error %d\n", filename, err);
		return -1;
	}

	/* find the input file's stream - assume this further sets up info w/in the AVFormatContext */
    err = av_find_stream_info(context->Fctx);
    if (err < 0) {
        fprintf(stderr, "av_find_stream_info: error %d\n", err);
		return -1;
    }

	AVCodecContext *decoder = context->Fctx->streams[0]->codec;
	context->Cctx = decoder;

    /* set up Info context */
	context->Info.Content_length = 0;
	context->Info.Sample_rate = decoder->sample_rate;
	context->Info.Channels = decoder->channels;
	context->Info.Frame_size = decoder->frame_size;
	context->Info.Sample_size = av_get_bits_per_sample_fmt(decoder->sample_fmt) >> 3;

	/* find the audio decoder */
    context->Codec = avcodec_find_decoder(decoder->codec_id);
    if (!context->Codec) {
        fprintf(stderr, "codec not found\n");
		return -1;
    }
	
    /* open the codec */
    if (avcodec_open(context->Cctx, context->Codec) < 0) {
        fprintf(stderr, "could not open codec\n");
        exit(1);
    }

	/* alloc the output buffer for libavcodec */
	context->Outbuf = (int16_t *)malloc(AVCODEC_MAX_AUDIO_FRAME_SIZE);
	context->Buf_size = AVCODEC_MAX_AUDIO_FRAME_SIZE;
	context->first_frame_used = -1;

	/* decode the first packet, so that Info.Frame_size can be set ... :\ */
	decode_packet(context, &len);
	if (len < 0) {
        fprintf(stderr, "decode_packet: error %d\n", err);
		return -1;
	}
	
	if (context->Info.Frame_size == 0) {
		context->Info.Frame_size = len / (context->Info.Channels * context->Info.Sample_size);
	}

	context->first_frame_used = 0;

	return 0;
}


/**
 * decode a packet
 *
 * @param context the AVDecodeContext
 * @param decoded_bytes a place to store the number of decoded bytes, ignore if null
 * @return size of output buffer on success, -1 on error
 */
int decode_packet(AVDecodeContext *context, int *decoded_bytes) {
	int err, out_size, len;

	if (context->first_frame_used == 0) {
		context->first_frame_used = 1;
	}

	/* try to read a frame from the context */
	err = av_read_frame(context->Fctx, &context->Packet);

	if (err < 0) {
		fprintf(stderr, "av_read_frame: error %d\n", err);
		return err;
	}

	/* check for id3 tag */
	if (!is_id3_tag(context)) {
		/* we'll take as much as you can give us */
        out_size = AVCODEC_MAX_AUDIO_FRAME_SIZE;

        // av_pkt_dump_log(NULL, AV_LOG_DEBUG, &pkt, 1);

#ifdef AVCORE_SAMPLEFMT_H
		len = avcodec_decode_audio3(context->Cctx, (short *)context->Outbuf, &out_size, &context->Packet);

		// fprintf(stderr, "len: %d, dts: %lld pts: %lld pkt size: %d out size: %d duration: %d pos %lld\n", len,
		// 			context->Packet.dts, context->Packet.pts, context->Packet.size, out_size, context->Packet.duration, context->Packet.pos);
#else
		AVPacket *packet = &context->Packet;
		uint8_t *packetData = packet->data;
		int packetSize = packet->size;
		len = avcodec_decode_audio2(context->Cctx, (short *)context->Outbuf, &out_size, packetData, packetSize);
#endif
        if (len < 0) {
            fprintf(stderr, "Error while decoding, len: %d, out_size: %d\n", len, out_size);
			return -1;
        }

		if (decoded_bytes != NULL) {
			*decoded_bytes = out_size;
		}

		return context->Info.Frame_size;
	}

	return 0;
}


/**
 * return the AVStreamInfo struct from an AVDecodeContext
 *
 * @param context the AVDecodeContext
 * @return the AVStreamInfo thereof
 */
AVStreamInfo get_stream_info(AVDecodeContext *context) {
	return context->Info;
}

/**
 * check for presence of ID3 Tag in an AVDecodeContext
 *
 * @param context the AVDecodeContext
 * @return 1 if the current packet is an ID3 tag, 0 otherwise
 */
int is_id3_tag(AVDecodeContext *context) {
	/* FIXME: seek to stream_len - 128, and if it's there, subtract 128 from stream_len */
	if ((context->Packet.data[0] == 'T') && (context->Packet.data[1] == 'A') && (context->Packet.data[2] == 'G')) {
		return -1;
	}

	return 0;
}

