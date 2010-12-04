// Copyright (c) 2010 Go Fightclub Authors

package alsa

import (
  "fmt"
  "unsafe"
  "afp"
)

// #include <alsa/asoundlib.h>
import "C"

/////
// Alsa Source
// Listens to a microphone
type AlsaSource struct {
    ctx *afp.Context
    header StreamHeader
}

func (self *AlsaSource) GetType() int {
    return afp.PIPE_SOURCE
}

func (self *AlsaSource) Init(ctx *afp.Context, args []string) os.Error {
    self.ctx = ctx

    header := StreamHeader {
        Version: 1
        Channels: 1
        SampleSize: 32
        SampleRate: 44100
    }

    return self.prepare()
}

func (self *AlsaSource) Start() {
    var buf *[]float32

    for {
        errno := C.snd_pcm_readn(capture, unsafe.Pointer(buf), 512)
        if errno < 512 {
            errtwo := C.snd_pcm_recover(capture, C.int(errno), 0);
            if errtwo < 0 {
                fmt.Println( "While reading from ALSA device, failed to recover from error: ", errtwo)
                panic
            }
        }
    }
}

/////
// Alsa Sink
// Outputs to speakers via ALSA
type AlsaSink struct {
    ctx *afp.Context
    header StreamHeader
}

func (self *AlsaSink) GetType() int {
    return afp.PIPE_SINK
}

func (self *AlsaSink) Init(ctx *afp.Context, args []string) os.Error {
    self.ctx = ctx
    header <-ctx.HeaderSource
    return self.prepare()
}

func (self *AlsaSink) Start() {
    for buffer, ok := <-ctx.Source; ok {
        length := len(buffer)
        errno := C.snd_pcm_writen(playback, unsafe.Pointer(buffer), length)

        if errno < length {
            panic //not all the data was written
        }
    }
}

// Ugly bastardized C code follows
func (self *AlsaSink) prepare() {

    var playback *C.snd_pcm_t
    var params *C.snd_pcm_hw_params_t

    if errno := C.snd_pcm_open(&playback, C.CString("default"), C.SND_PCM_STREAM_PLAYBACK, 0); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not open device. Error %d", errno) )
    }

    defer C.snd_pcm_close(playback)

    if errno := C.snd_pcm_hw_params_malloc(&params); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not allocate hardware parameter structure. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_any(playback, params); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not initialize hardware parameter structure. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_set_access (playback, params, C.SND_PCM_ACCESS_RW_INTERLEAVED); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set access type. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_set_format(playback, params, C.SND_PCM_FORMAT_FLOAT); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set sample format. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_set_rate(playback, params, header.SampleRate, 0); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set sample rate. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_set_channels(playback, params, header.Channels); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set channel count. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params(playback, params); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set hardware parameters. Error %d", errno) )
    }

    C.snd_pcm_hw_params_free(params)

    if errno := C.snd_pcm_prepare(playback); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not prepare audio device for use. Error %d", errno) )
    }

}

//this one is slightly different
//note the change in scope
func (self *AlsaSource) prepare() {

    var capture *C.snd_pcm_t
    var params *C.snd_pcm_hw_params_t

    if errno := C.snd_pcm_open(&capture, C.CString("default"), C.SND_PCM_STREAM_CAPTURE, 0); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not open device. Error %d", errno) )
    }

    defer C.snd_pcm_close(capture)

    if errno := C.snd_pcm_hw_params_malloc(&params); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not allocate hardware parameters. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_any(capture, params); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not initialize hardware parameter structure. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_set_access(capture, params, C.SND_PCM_ACCESS_RW_INTERLEAVED); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set access. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_set_format(capture, params, C.SND_PCM_FORMAT_FLOAT); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set sample format. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_set_rate(capture, params, header.SampleRate, 0); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set sample rate. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_set_channels(capture, params, 1); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set channel count. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params(capture, params); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set parameters. Error %d", errno) )
    }

    C.snd_pcm_hw_params_free(params)

    if errno := C.snd_pcm_prepare(capture); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not prepare audio interface for use. Error %d", errno) )
    }

    return nil
}
