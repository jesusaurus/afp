// Copyright (c) 2010 Go Fightclub Authors

package alsa

import (
  "fmt"
  "unsafe"
  "afp"
)

// #include <alsa/asoundlib.h>
import "C"

type AlsaFilter struct {
    ctx *afp.Context
    header StreamHeader
}

func (self *AlsaFilter) GetType() int {
    return afp.PIPE_SINK
}

func (self *AlsaFilter) Init(ctx *afp.Context, args []string) os.Error {
    self.ctx = ctx

    header <-ctx.HeaderSource

    return self.prepare()
}

func (self *AlsaFilter) Start() {
    for buffer, ok := <-ctx.Source; ok {
        length := len(buffer)
        errno := C.snd_pcm_writen(playback, unsafe.Pointer(buffer), length)

        if errno < length {
            return //not all the data was written
        }
    }
}

// Ugly bastardized C code follows

func (self *AlsaFilter) prepare() {

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
        return os.NewError( fmt.Printf("Could not initialize hardware parameter structure. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_set_access (playback, params, C.SND_PCM_ACCESS_RW_INTERLEAVED); errno < 0 {
        return os.NewError( fmt.Printf("Could not set access type. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_set_format(playback, params, C.SND_PCM_FORMAT_FLOAT); errno < 0 {
        return os.NewError( fmt.Printf("Could not set sample format. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_set_rate(playback, params, header.SampleRate, 0); errno < 0 {
        return os.NewError( fmt.Printf("Could not set sample rate. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_set_channels(playback, params, header.Channels); errno < 0 {
        return os.NewError( fmt.Printf("Could not set channel count. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params(playback, params); errno < 0 {
        return os.NewError( fmt.Printf("Could not set hardware parameters. Error %d", errno) )
    }

    C.snd_pcm_hw_params_free(params)

    if errno := C.snd_pcm_prepare(playback); errno < 0 {
        return os.NewError( fmt.Printf("Could not prepare audio device for use. Error %d", errno) )
    }

}
