/*
 * copyright (c) 2001 Fabrice Bellard
 * bastardized (b) 2010 Eric O'Connell
 *
 * build this in the ffmpeg directory, and cp a song to /tmp/test.mp3
 * if all goes well, you will get /tmp/test.raw, which can be imported into Audacity
 * File > Import > Raw Data, default options seemed to work for me.
 *
 * This file is part of FFmpeg.
 *
 * FFmpeg is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2.1 of the License, or (at your option) any later version.
 *
 * FFmpeg is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public
 * License along with FFmpeg; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA
 */


/**
 * @file
 * avcodec API use example.
 *
 * Note that this library only handles codecs (mpeg, mpeg4, etc...),
 * not file formats (avi, vob, etc...). See library 'libavformat' for the
 * format handling
 */


#include <stdlib.h>
#include <stdio.h>
#include <string.h>


#ifdef HAVE_AV_CONFIG_H
#undef HAVE_AV_CONFIG_H
#endif


#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <libavutil/mathematics.h>


#define AUDIO_REFILL_THRESH 4096


/*
 * Audio decoding.
 */
static void audio_decode_example(const char *filename, const char *outfilename)
{
    AVCodec *codec;
    AVCodecContext *c= NULL;
    AVFormatContext *fctx;
	int out_size, len, err, i = 0;
    FILE *outfile;
    uint8_t *outbuf;
    AVPacket avpkt;

    av_init_packet(&avpkt);

    printf("Audio decoding\n");

    outbuf = malloc(AVCODEC_MAX_AUDIO_FRAME_SIZE);

	/* use libavformat to open the source file; initializes AVFormatContext */
    if ((err = av_open_input_file(&fctx, filename, NULL, 0, NULL)) < 0) {
        fprintf(stderr, "av_open_input_file: error %d opening %s\n", err, filename);
        exit(1);
    }

	/* find the input file's stream - assume this further sets up info w/in the AVFormatContext */
    err = av_find_stream_info(fctx);
    if (err < 0) {
        fprintf(stderr, "av_find_stream_info: error %d\n", err);
        exit(1);
    }

	enum CodecID guessed_codec;

	for (i = 0; i < fctx->nb_streams; i++) {
		AVStream *st = fctx->streams[i];
        AVCodecContext *dec = st->codec;
        
        int audio_channels    = dec->channels;
        int audio_sample_rate = dec->sample_rate;
        enum AVSampleFormat audio_sample_fmt  = dec->sample_fmt;
        
		fprintf(stderr, "channels: %d sample_rate: %d frame_size: %d\n", audio_channels, audio_sample_rate, dec->frame_size);
		guessed_codec = dec->codec_id;
	}

	dump_format(fctx, 0, filename, 0);
	
	// AVOutputFormat *guessed_format = av_guess_format(NULL, filename, NULL);
	// enum CodecID guessed_codec = av_guess_codec(guessed_format, NULL, filename, NULL, AVMEDIA_TYPE_AUDIO);
	
	/* find the audio decoder */
    codec = avcodec_find_decoder(guessed_codec);
    if (!codec) {
        fprintf(stderr, "codec not found\n");
		exit(1);
    }

	fprintf(stderr, "Guessed codec: %s\n", codec->name);

	c = fctx->streams[0]->codec;

    // c= avcodec_alloc_context();
    // 
    /* open the codec */
    if ((err = avcodec_open(c, codec)) < 0) {
        fprintf(stderr, "could not open codec: %d\n", err);
        exit(1);
    }

	/* open the output file */
    outfile = fopen(outfilename, "wb");
    if (!outfile) {
        av_free(c);
        exit(1);
    }

	fprintf(stderr, "Duration: %lld, file size: %lld\n", fctx->duration, fctx->file_size);

	/* read frames into the AVPacket, then decode them */
    while ((err = av_read_frame(fctx, &avpkt)) >= 0) {
		fprintf(stderr, "i: %d, dts: %lld pts: %lld len: %d pkt size: %d out size: %d duration: %d pos %lld\n", i++,
					avpkt.dts, avpkt.pts, len, avpkt.size, out_size, avpkt.duration, avpkt.pos);

        out_size = AVCODEC_MAX_AUDIO_FRAME_SIZE;

		/* did libavformat just hand us an ID3V1 tag? */
		if (!((avpkt.data[0] == 'T') && (avpkt.data[1] == 'A') && (avpkt.data[2] == 'G'))) {
	        len = avcodec_decode_audio3(c, (short *)outbuf, &out_size, &avpkt);
	
	        if (len < 0) {
	            fprintf(stderr, "Error while decoding, len: %d, out_size: %d\n", len, out_size);
	            exit(1);
	        }

	        if (out_size > 0) {
	            /* if a frame has been decoded, output it */
	            fwrite(outbuf, 1, out_size, outfile);
	        }	
		}

		/* release memory */
        av_free_packet(&avpkt);
    }

    fclose(outfile);
    free(outbuf);

    avcodec_close(c);
    av_free(c);
}


int main(int argc, char **argv)
{
    const char *filename;

	if (argc != 3) {
		fprintf(stderr, "Damn, yo. call me with input & output paths\n");
		exit(-1);
	}

	/* must be called before using avformat lib */
    av_register_all();

    /* must be called before using avcodec lib */
    avcodec_init();

    /* register all the codecs */
    avcodec_register_all();

    audio_decode_example(argv[1], argv[2]);

    return 0;
}
