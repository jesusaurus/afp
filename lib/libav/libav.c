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
int decode_packet(AVDecodeContext *context);
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
}


/**
 * set up a decoding context
 *
 * @param the filename to decode
 * @param a pointer to an AVDecodeContext, will be initialized
 * @return 0 on success, -1 on error
 */
int prepare_decoding(char *filename, AVDecodeContext *context) {
	int err, len, sample_size;
    
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

	char buf[256];
	AVStream *st = context->Fctx->streams[0];
    
	avcodec_string(buf, sizeof(buf), st->codec, 0);
	fprintf(stderr, "codec string: %s", buf);
	// dump_format(context->Fctx, 0, filename, 0);

	/* initialize the AVPacket */
    av_init_packet(&context->Packet);
    
	AVOutputFormat *guessed_format = av_guess_format(NULL, filename, NULL);
	enum CodecID guessed_codec = av_guess_codec(guessed_format, NULL, filename, NULL, AVMEDIA_TYPE_AUDIO);

	/* find the audio decoder */
    context->Codec = avcodec_find_decoder(guessed_codec);
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
	context->Outbuf = (int16_t *)malloc(AVCODEC_MAX_AUDIO_FRAME_SIZE);
	context->Buf_size = AVCODEC_MAX_AUDIO_FRAME_SIZE;
	context->first_frame_used = -1;
	
	/* decode the first packet, so that Info.Frame_size can be set ... :\ */
	len = decode_packet(context);
	if (len < 0) {
        fprintf(stderr, "decode_packet: error %d\n", err);
		return -1;
	}

    /* set up Info context */
	context->Info.Content_length = 0;

	// FIXME: this will only work for mp3 files :(
		int bit_rot;
		uint32_t head = AV_RB32(context->Packet.data);
		err = ff_mpa_decode_header(context->Cctx, 
			head, 
			&context->Info.Sample_rate, 
			&context->Info.Channels,
			&context->Info.Frame_size,
			&bit_rot);
	
		MPADecodeHeader s1, *s = &s1;
		if (err < 0) {
			fprintf(stderr, "Error from ff_mpa_decode_header: %d\n", err);
			fprintf(stderr, "ff_mpa_check_header: %d\n", ff_mpa_check_header(head));
			fprintf(stderr, "ff_mpegaudio_decode_header: %d\n", ff_mpegaudio_decode_header(s, head));

			return -1;
		}
	// FIXME
	
	// context->Info.Channels = context->Cctx->channels;
	// context->Info.Sample_rate = context->Cctx->sample_rate;
	switch(context->Cctx->sample_fmt) {
	    case AV_SAMPLE_FMT_U8:          ///< unsigned 8 bits
			sample_size = 1;
			break;
	    case AV_SAMPLE_FMT_S16:         ///< signed 16 bits
			sample_size = 2;
			break;
	    case AV_SAMPLE_FMT_S32:         ///< signed 32 bits
			sample_size = 4;
			break;
		default:
			fprintf(stderr, "Unsupported sample format: %d\n", context->Cctx->sample_fmt);
			return -1;
	}
	context->Info.Sample_size = sample_size;
	
	context->first_frame_used = 0;
	context->Info.Frame_size = len / context->Info.Channels / context->Info.Sample_size;

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
	
	if (context->first_frame_used == 0) {
		context->first_frame_used = 1;
		return context->Info.Frame_size;
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

		len = avcodec_decode_audio3(context->Cctx, (short *)context->Outbuf, &out_size, &context->Packet);

        if (out_size < 0) {
            fprintf(stderr, "Error while decoding, len: %d, out_size: %d\n", len, out_size);
			return -1;
        }

		return out_size;
	}
	
	return 0;
}


/**
 * return the AVStreamInfo struct from an AVDecodeContext
 *
 * @param the AVDecodeContext
 * @return the AVStreamInfo thereof
 */
AVStreamInfo get_stream_info(AVDecodeContext *context) {
	return context->Info;
}

/**
 * check for presence of ID3 Tag in an AVDecodeContext
 *
 * @param the AVDecodeContext
 * @return 1 if the current packet is an ID3 tag, 0 otherwise
 */
int is_id3_tag(AVDecodeContext *context) {
	/* FIXME: seek to stream_len - 128, and if it's there, subtract 128 from stream_len */
	if ((context->Packet.data[0] == 'T') && (context->Packet.data[1] == 'A') && (context->Packet.data[2] == 'G')) {
		return -1;
	}
	
	return 0;
}
 
