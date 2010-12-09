// Copyright (c) 2010 Go Fightclub Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

/**
 * decode_wrapper.h
 * (c) Eric O'Connell 2010
 *
 * provide a minimal (possibly braindead) interface into audio decoding using
 * libavformat to demux input files & libavcodec to produce PCM data
 *
 */

#ifndef __DECODE_WRAPPER_H__
#define __DECODE_WRAPPER_H__

#include "libav.h"

/**
 * initialize all of the fun libavformat & libavdecode stuff
 *
 */
void init_decoding(void);

/**
 * prepare a decoding context for an input filename
 *
 */
int prepare_decoding(const char *filename, AVDecodeContext *context);

#endif
