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
static void audio_decode_example(const char *outfilename, const char *filename)
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

    /* find the mpeg audio decoder */
    codec = avcodec_find_decoder(CODEC_ID_MP3);
    if (!codec) {
        fprintf(stderr, "codec not found\n");
        exit(1);
    }

    c= avcodec_alloc_context();

    /* open the codec */
    if (avcodec_open(c, codec) < 0) {
        fprintf(stderr, "could not open codec\n");
        exit(1);
    }

    outbuf = malloc(AVCODEC_MAX_AUDIO_FRAME_SIZE);

	/* use libavformat to open the source file; initializes AVFormatContext */
    if ((err = av_open_input_file(&fctx, filename, NULL, 0, NULL)) < 0) {
        fprintf(stderr, "av_open_input_file: error %d\n", err);
        exit(1);
    }

	/* find the input file's stream - assume this further sets up info w/in the AVFormatContext */
    err = av_find_stream_info(fctx);
    if (err < 0) {
        fprintf(stderr, "av_find_stream_info: error %d\n", err);
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
	        if (out_size < 0) {
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

	/* must be called before using avformat lib */
    av_register_all();

    /* must be called before using avcodec lib */
    avcodec_init();

    /* register all the codecs */
    avcodec_register_all();

    if (argc <= 1) {
        // audio_encode_example("/tmp/test.mp2");
        audio_decode_example("/tmp/test.raw", "/tmp/test.mp3");
    } else {
        filename = argv[1];
    }

    return 0;
}
