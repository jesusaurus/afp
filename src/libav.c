/**
 * libav.c
 * (c) Eric O'Connell 2010
 *
 * provide a minimal (possibly braindead) interface into audio decoding using
 * libavformat to demux input files & libavcodec to produce PCM data
 *
 */


#include "libav.h"
#include <string.h>

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
}


/**
 * set up a decoding context
 *
 * @param the filename to decode
 * @param a pointer to an AVDecodeContext, will be initialized
 * @return 0 on success, -1 on error
 */
int prepare_decoding(char *filename, AVDecodeContext *context) {
	int err;
	
	/* initialize the AVPacket */
    av_init_packet(&context->Packet);
    
	/* find the mp3 audio decoder */
    context->Codec = avcodec_find_decoder(CODEC_ID_MP3);
    if (!context->Codec) {
        fprintf(stderr, "codec not found\n");
		return -1;
    }

	/* set up the codec context */
    context->Cctx = avcodec_alloc_context();

    /* open the codec */
    if (avcodec_open(context->Cctx, context->Codec) < 0) {
        fprintf(stderr, "could not open codec\n");
        exit(1);
    }
    
	/* alloc the output buffer for libavcodec */
	context->Outbuf = malloc(AVCODEC_MAX_AUDIO_FRAME_SIZE);
	context->Buf_size = AVCODEC_MAX_AUDIO_FRAME_SIZE;
    
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

	return 0;
}


/**
 * decode a packet
 *
 * @param the AVDecodeContext
 * @return size of output buffer on success, -1 on error
 */
int decode_packet(AVDecodeContext *context) {
	int err, out_size, len;

	/* try to read a frame from the context */
	err = av_read_frame(context->Fctx, &context->Packet);
	
	if (err < 0) {
		fprintf(stderr, "av_read_frame: error %d\n", err);
		return err;
	}
	
	/* check for id3 tag */
	if (!((context->Packet.data[0] == 'T') && (context->Packet.data[1] == 'A') && (context->Packet.data[2] == 'G'))) {
		/* we'll take as much as you can give us */
        out_size = AVCODEC_MAX_AUDIO_FRAME_SIZE;

		len = avcodec_decode_audio3(context->Cctx, (short *)context->Outbuf, &out_size, &context->Packet);
        if (len < 0) {
            fprintf(stderr, "avcodec_decode_audio3: Error while decoding: %d\n", len);
			return -1;
        }

		return out_size;
	}
	
	return 0;
}
