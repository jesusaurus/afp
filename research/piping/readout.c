// Copyright (c) 2010 Go Fightclub Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

#include <stdlib.h>
#include <stdio.h>

#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <libavutil/avutil.h>
#include <libavutil/mathematics.h>

#ifdef HAVE_AV_CONFIG_H
#undef HAVE_AV_CONFIG_H
#endif

#define AUDIO_INBUF_SIZE 20480
#define AUDIO_REFILL_THRESH 4096

#ifndef FIGHTCLUB_WRITEOUT
#define FIGHTCLUB_WRITEOUT

int WriteOut(char* file) {
    AVCodec *codec;
    AVCodecContext *cContext = NULL;
    //AVFormatContext *fContext = NULL;
    AVPacket packet;
    int frame_size, out_size;
    short *samples;
    //float t, tincr;
    uint8_t *outbuf;
    FILE *fout;

    //standard initializations
    av_register_all();
    avcodec_init();
    avcodec_register_all();

    av_init_packet(&packet);

    //load the codec
    codec = avcodec_find_encoder(CODEC_ID_MP3); //TODO: get ogg/flac/${opensourcecodec} working
    if (!codec){
        fprintf(stderr, "CODEC_ID_MP3 not found.");
        return 1;
    }

    //create the codec context
    cContext = avcodec_alloc_context();
    cContext->bit_rate = 128000;
    cContext->sample_rate = 44100;
    cContext->channels = 2;
    if (avcodec_open(cContext, codec) < 0) {
        fprintf(stderr, "Failed to open codec.");
        return 1;
    }

    frame_size = cContext->frame_size;
    samples = malloc(frame_size * 2 * cContext->channels);
    outbuf = malloc(AVCODEC_MAX_AUDIO_FRAME_SIZE);

    /*
    //create the format context
    fContext = avformat_alloc_context();
    */

    fout = fopen(file, "wb");
    if (!fout) {
        fprintf(stderr, "Error opening output file\n");
        return 1;
    }

    //read stdin to the buffer, then out to file
    do {
        fread(samples, 1, frame_size, stdin);
        out_size = avcodec_encode_audio(cContext, outbuf, AVCODEC_MAX_AUDIO_FRAME_SIZE, samples);
        fwrite(outbuf, 1, out_size, fout);
    } while (!feof(stdin) && !ferror(stdin));

    //clean up
    fclose(fout);
    free(outbuf);
    free(samples);

    avcodec_close(cContext);
    av_free(cContext);

    return 0;
}

#endif
