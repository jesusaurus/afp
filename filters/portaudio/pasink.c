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
#include "portaudio.h"

#define PROTECT(x) if((x) < 0) { perror(#x); return -1; }
#define LOCK(x) if((pthread_mutex_lock(x)) < 0) { perror(#x); return -1; }
#define UNLOCK(x) if((pthread_mutex_unlock(x)) < 0) { perror(#x); return -1; }

#define BUFFERS 4

/* context struct for buffering output audio, passed to context */
typedef struct {
	PaStream *stream;

	float **buffers;
	int read_index;
	int write_index;

	pthread_mutex_t reading[BUFFERS];
	pthread_mutex_t writing[BUFFERS];

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
 */
static int pa_output_callback(	const void *inputBuffer, void *outputBuffer,
	                            unsigned long framesPerBuffer,
	                            const PaStreamCallbackTimeInfo* timeInfo,
	                            PaStreamCallbackFlags statusFlags,
	                            void *userData ) {
	pa_output_data *data = (pa_output_data*)userData;
	float *out = (float*)outputBuffer;
	int locked = data->read_index;

	(void) timeInfo; /* Prevent unused variable warnings. */
	(void) statusFlags;
	(void) inputBuffer;
	
	/* acquire lock for reading from the current output buffer */
	LOCK(&data->reading[locked]);

	memcpy(out, data->buffers[data->read_index], data->buffer_size * sizeof(float));

	/* increment our read position in the buffer ring */
	data->read_index = (data->read_index + 1) % BUFFERS;

	/* once we're done reading, unlock both reading and writing so this buffer can be filled again */
	UNLOCK(&data->reading[locked]);
	UNLOCK(&data->writing[locked]);
	
	return paContinue;
}

/**
 * send pasink output data
 */
int send_output_data(float *interleaved_float_samples, pa_output_data *data, int done) {
	PaError err = 0;
	int locked = data->write_index;
	
	/* check to see if our source is done playing */
	if (done != 0) {
		data->stopped = 1;
		
		/* Pa_StopStream will inform PortAudio we're done, and let it play any remaining buffers available */
		err = Pa_StopStream( data->stream );
		if (err != 0) {
			fprintf(stderr, "Error with Pa_StopStream\n");
		    fprintf( stderr, "Error number: %d\n", err );
		    fprintf( stderr, "Error message: %s\n", Pa_GetErrorText( err ) );
			return err;
		}

		return 0;
	}

	/* once we can write to the current buffer, prevent the callback from reading it */
	LOCK(&data->writing[locked]);
	LOCK(&data->reading[locked]);
	
	/* copy data into the output buffer */
	memcpy((void *)data->buffers[data->write_index], (const void *)interleaved_float_samples, (size_t)(data->buffer_size * sizeof(float)));

	/* increment our write position in the buffer ring */
	data->write_index = (data->write_index + 1) % BUFFERS;

	/* start playing once we've filled all the BUFFERS */
	if (data->write_index == 0 && data->started == 0) {
		err = Pa_StartStream( data->stream );
		if (err != 0) {
			fprintf(stderr, "Error with Pa_StartStream\n");
		    fprintf( stderr, "Error number: %d\n", err );
		    fprintf( stderr, "Error message: %s\n", Pa_GetErrorText( err ) );

			/* let go of this mutex, too */
			UNLOCK(&data->writing[locked]);
		} else {
			data->started = 1;
		}
	}

	/* we're done writing so this buffer is available for reading */
	UNLOCK(&data->reading[locked]);

	return err;
}

/**
 * set up port audio with the right number of channels, the sample rate, and frame size
 * configure callbacks for output & end of output
 */
int init_portaudio_output(int channels, int sample_rate, int frame_size, pa_output_data *data) {
	PaStreamParameters outputParameters;
	PaError err;
	int i;
	
	if ((data->buffers = (float**)malloc(BUFFERS)) < 0) {
		fprintf(stderr,"Error: Not enough memory");
		return errno;
	}
	
	for (i = 0; i < BUFFERS; i++) {
		if ((data->buffers[i] = (float*)malloc(frame_size * channels * sizeof(float))) < 0) {
			fprintf(stderr,"Error: Not enough memory");
			return errno;
		}
		PROTECT(pthread_mutex_init(&data->reading[i], NULL));
		PROTECT(pthread_mutex_init(&data->writing[i], NULL));
	}
	
	data->started = 0;
	data->stopped = 0;
	data->read_index = 0;
	data->write_index = 0;
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

	for (i = 0; i < BUFFERS; i++) {
		PROTECT(pthread_mutex_destroy(&data->reading[i]));
		PROTECT(pthread_mutex_destroy(&data->writing[i]));
		free(data->buffers[i]);
	}
	free(data->buffers);

    Pa_Terminate();
	return err;
}