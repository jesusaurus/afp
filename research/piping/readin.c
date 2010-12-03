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

#ifndef FIGHTCLUB_READIN
#define FIGHTCLUB_READIN

//primarily based off of http://cekirdek.pardus.org.tr/~ismail/ffmpeg-docs/api-example_8c-source.html
//thanks to Eric O'Connell for discovering pktdumper.c and throwing it in the mix
int ReadIn(char* file) {
    AVCodec *codec;
    AVCodecContext *cContext = NULL;
    AVFormatContext *fContext = NULL;
    AVPacket packet;
    int out_size, len, err;
    //uint8_t inbuf[AUDIO_INBUF_SIZE + FF_INPUT_BUFFER_PADDING_SIZE];
    uint8_t *outbuf;
    //FILE *fin;
    //FILE *fout;

    //standard initializations
    av_register_all();
    avcodec_init();
    avcodec_register_all();
    
    av_init_packet(&packet);

    //load the codec
    codec = avcodec_find_decoder(CODEC_ID_MP3);
    if (!codec){
        fprintf(stderr, "CODEC_ID_MP3 not found.");
        return 1;
    }

    //load the codec context
    cContext = avcodec_alloc_context();
    if (avcodec_open(cContext, codec) < 0) {
        fprintf(stderr, "Failed to open codec.");
        return 1;
    }

    outbuf = malloc(AVCODEC_MAX_AUDIO_FRAME_SIZE);

    //open the source file
    if ((err = av_open_input_file(&fContext, file, NULL, 0, NULL)) < 0) {
        fprintf(stderr, "Failed to open input file: %d\n", err);
        return 1;
    }

    //open the input stream
    if ((err = av_find_stream_info(fContext)) < 0) {
        fprintf(stderr, "Failed to open stream: %d\n", err);
        return 1;
    }

    /*
    //open the output file
    fout = fopen("/tmp/test.mp3.raw", "wb");
    if (!fout) {
        fprintf(stderr, "Failed to open test.mp3.raw");
        av_free(cContext);
        return 1;
    }
    */

    fprintf(stderr, "Reading\n");

    while ((err = av_read_frame(fContext, &packet)) >= 0) {
        fprintf(stderr, ".");

        out_size = AVCODEC_MAX_AUDIO_FRAME_SIZE;

        //check for the final frame of an mp3 (id3v1 tag)
        if ( (packet.data[0] != 'T') && (packet.data[1] != 'A') && (packet.data[2] != 'G') ) {
            len = avcodec_decode_audio3(cContext, (short *)outbuf, &out_size, &packet);

            if (len < 0) {
                fprintf(stderr, "Error decoding\n");
                return 1;
            }

            if (out_size > 0) {
                fwrite(outbuf, 1, out_size, stdout);
            }
        }

        av_free_packet(&packet);
    }

    fprintf(stderr, "\n");
    fflush(stderr);

    //fclose(fout);
    free(outbuf);
    avcodec_close(cContext);
    av_free(cContext);

    return 0;
}


#endif
