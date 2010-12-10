// Copyright (c) 2010 Go Fightclub Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

// A small wrapper to handle the callback for portaudio output
// borrowed heavily from portaudio/test/patest_sine.c
// see: http://www.portaudio.com/trac/wiki/TutorialDir/Exploring

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <pthread.h>
#include <portaudio.h>

#define PROTECT(x) if((x) < 0) { perror(#x); return -1; }
#define LOCK(x) if((pthread_mutex_lock(x)) < 0) { perror(#x); return -1; }
#define UNLOCK(x) if((pthread_mutex_unlock(x)) < 0) { perror(#x); return -1; }

#define BUFFERS 4

/* context struct for buffering output audio, passed to context */
typedef struct {
	PaStream *stream;
	float **buffers;
	pthread_mutex_t queueing, writing, reading;
	int dirty[BUFFERS];
	int buffer_index;
	int fill_index;
	int buffer_size;
	int started;
	int stopped;
} pa_output_data;


int init_portaudio_output(int channels, int sample_rate, int frame_size, pa_output_data *output);
static int pa_output_callback(	const void *inputBuffer, void *outputBuffer,
	                            unsigned long framesPerBuffer,
	                            const PaStreamCallbackTimeInfo* timeInfo,
	                            PaStreamCallbackFlags statusFlags,
								void *userData );

/**
 * send the float data into the output buffer provided by pa_audio
 * TODO: just call memcpy()
 */
static int pa_output_callback(	const void *inputBuffer, void *outputBuffer,
	                            unsigned long framesPerBuffer,
	                            const PaStreamCallbackTimeInfo* timeInfo,
	                            PaStreamCallbackFlags statusFlags,
	                            void *userData ) {
	pa_output_data *data = (pa_output_data*)userData;
	float *out = (float*)outputBuffer;

	(void) timeInfo; /* Prevent unused variable warnings. */
	(void) statusFlags;
	(void) inputBuffer;
	
	while (1) {
		LOCK(&data->queueing);
		if (data->dirty[data->buffer_index] == 1 && data->stopped == 0)
			break;
		UNLOCK(&data->queueing);
	}
	UNLOCK(&data->queueing);

	memcpy(out, data->buffers[data->buffer_index], data->buffer_size * sizeof(float));

	/* this buffer is no longer dirty */
	data->dirty[data->buffer_index] = 0;
	data->buffer_index = (data->buffer_index + 1) % BUFFERS;

	return paContinue;
}

/**
 * send pasink output data
 */
int send_output_data(float *interleaved_float_samples, pa_output_data *data, int done) {
	PaError err;
	
	if (done != 0) {
		data->stopped = 1;
		
		err = Pa_StopStream( data->stream );
		if (err != 0) {
			fprintf(stderr, "Error with Pa_StopStream\n");
		    fprintf( stderr, "Error number: %d\n", err );
		    fprintf( stderr, "Error message: %s\n", Pa_GetErrorText( err ) );
		}
	}

	while (1) {
		LOCK(&data->queueing);
		if (data->dirty[data->buffer_index] == 0 || data->started == 0) {
			break;
		}
		UNLOCK(&data->queueing);
	}
	UNLOCK(&data->queueing);
	
	/* copy data into the output buffer */
	memcpy((void *)data->buffers[data->fill_index], (const void *)interleaved_float_samples, (size_t)(data->buffer_size * sizeof(float)));
	data->dirty[data->fill_index] = 1;
	data->fill_index = (data->fill_index + 1) % BUFFERS;

	/* start playing once we've filled all the BUFFERS */
	if (data->fill_index == 0 && data->started == 0) {
		err = Pa_StartStream( data->stream );
		if (err != 0) {
			fprintf(stderr, "Error with Pa_StartStream\n");
		    fprintf( stderr, "Error number: %d\n", err );
		    fprintf( stderr, "Error message: %s\n", Pa_GetErrorText( err ) );
		} else {
			fprintf(stderr, "Starting stream\n");
			data->started = 1;
		}
		return err;
	}
	
	return 0;
}

/**
 * set up port audio with the right number of channels, the sample rate, and frame size
 * configure callbacks for output & end of output
 */
int init_portaudio_output(int channels, int sample_rate, int frame_size, pa_output_data *data) {
	PaStreamParameters outputParameters;
	PaError err;
	int i;
	
	PROTECT(pthread_mutex_init(&data->queueing, NULL))
	PROTECT(pthread_mutex_init(&data->reading, NULL))
	PROTECT(pthread_mutex_init(&data->writing, NULL))
	
	if ((data->buffers = (float**)malloc(BUFFERS)) < 0) {
		fprintf(stderr,"Error: Not enough memory");
		return errno;
	}
	
	for (i = 0; i < BUFFERS; i++) {
		if ((data->buffers[i] = (float*)malloc(frame_size * channels * sizeof(float))) < 0) {
			fprintf(stderr,"Error: Not enough memory");
			return errno;
		}
		data->dirty[i] = 0;
	}
	
	data->started = 0;
	data->stopped = 0;
	data->buffer_index = 0;
	data->fill_index = 0;
	data->buffer_size = channels * frame_size;
	
	err = Pa_Initialize();
    if( err != paNoError ) goto error;
    
	outputParameters.device = Pa_GetDefaultOutputDevice(); /* default output device */
    if (outputParameters.device == paNoDevice) {
      fprintf(stderr,"Error: No default output device.\n");
      goto error;
    }
    
	outputParameters.channelCount = channels;
    outputParameters.sampleFormat = paFloat32; 		/* 32 bit floating point output */
    outputParameters.suggestedLatency = Pa_GetDeviceInfo( outputParameters.device )->defaultLowOutputLatency;
    outputParameters.hostApiSpecificStreamInfo = NULL;
    
	err = Pa_OpenStream(
              &(data->stream),
              NULL, /* no input */
              &outputParameters,
              sample_rate,
              frame_size,
              paClipOff,      /* we won't output out of range samples so don't bother clipping them */
              pa_output_callback,
              data );
    if( err != paNoError ) goto error;
    
	return 0;
	
error:
    Pa_Terminate();
    fprintf( stderr, "An error occured while using the portaudio stream\n" );
    fprintf( stderr, "Error number: %d\n", err );
    fprintf( stderr, "Error message: %s\n", Pa_GetErrorText( err ) );
    return err;
}

/**
 * close down portaudio
 */
int close_portaudio(pa_output_data *data) {
	int i;
	
	PaError err = Pa_CloseStream( data->stream );
    if( err != paNoError ) {
	    fprintf( stderr, "An error occured while using the portaudio stream\n" );
	    fprintf( stderr, "Error number: %d\n", err );
	    fprintf( stderr, "Error message: %s\n", Pa_GetErrorText( err ) );
	}

	PROTECT(pthread_mutex_destroy(&data->queueing));
	PROTECT(pthread_mutex_destroy(&data->reading));
	PROTECT(pthread_mutex_destroy(&data->writing));
	
	for (i = 0; i < BUFFERS; i++) {
		free(data->buffers[i]);
	}
	free(data->buffers);

    Pa_Terminate();
	return err;
}