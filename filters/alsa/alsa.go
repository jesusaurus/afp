// Copyright (c) 2010 Go Fightclub Authors

package alsa

import (
  "fmt"
  "unsafe"
  "afp"
  "os"
)

// #include <alsa/asoundlib.h>
import "C"

/////
// Alsa Source
// Listens to a microphone
type AlsaSource struct {
    ctx *afp.Context
    header afp.StreamHeader
    capture *C.snd_pcm_t
    params *C.snd_pcm_hw_params_t
}

func (self *AlsaSource) GetType() int {
    return afp.PIPE_SOURCE
}

func (self *AlsaSource) Init(ctx *afp.Context, args []string) os.Error {
    self.ctx = ctx

    header := afp.StreamHeader {
        Version: 1,
        Channels: 1,
        SampleSize: 32,
        SampleRate: 44100,
    }

    self.ctx.HeaderSink <- header

    retval := self.prepare()
    return retval
}

func (self *AlsaSource) Start() {
    var buf [512]float32

    for {
        errno := C.snd_pcm_readn(self.capture, unsafe.Pointer(&buf[0]), 512)
        if errno < 512 {
            errtwo := C.snd_pcm_recover(self.capture, C.int(errno), 0);
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
    header afp.StreamHeader
    playback *C.snd_pcm_t
    params *C.snd_pcm_hw_params_t

}

func (self *AlsaSink) GetType() int {
    return afp.PIPE_SINK
}

func (self *AlsaSink) Init(ctx *afp.Context, args []string) os.Error {
    self.ctx = ctx
    header <-self.ctx.HeaderSource
    retval := self.prepare()
    return retval
}

func (self *AlsaSink) Start() {
    buffer, ok := <-self.ctx.Source
    for ok {
        length := len(buffer)
        errno := C.snd_pcm_writen(playback, unsafe.Pointer(buffer), length)

        if errno < length {
            panic //not all the data was written
        }

        buffer, ok := <-self.ctx.Source
    }
}

// Ugly bastardized C code follows
func (self *AlsaSink) prepare() os.Error {

    if errno := C.snd_pcm_open(&self.playback, C.CString("default"), C.SND_PCM_STREAM_PLAYBACK, 0); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not open device. Error %d", errno) )
    }

    defer C.snd_pcm_close(self.playback)

    if errno := C.snd_pcm_hw_params_malloc(&self.params); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not allocate hardware parameter structure. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_any(self.playback, self.params); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not initialize hardware parameter structure. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_set_access(self.playback, self.params, C.SND_PCM_ACCESS_RW_INTERLEAVED); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set access type. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_set_format(self.playback, self.params, C.SND_PCM_FORMAT_FLOAT); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set sample format. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_set_rate(self.playback, self.params, self.header.SampleRate, 0); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set sample rate. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_set_channels(self.playback, self.params, self.header.Channels); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set channel count. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params(self.playback, self.params); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set hardware parameters. Error %d", errno) )
    }

    C.snd_pcm_hw_params_free(self.params)

    if errno := C.snd_pcm_prepare(self.playback); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not prepare audio device for use. Error %d", errno) )
    }

}

//this one is slightly different
//note the change in scope
func (self *AlsaSource) prepare() os.Error {

    if errno := C.snd_pcm_open(&self.capture, C.CString("default"), C.SND_PCM_STREAM_CAPTURE, 0); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not open device. Error %d", errno) )
    }

    defer C.snd_pcm_close(self.capture)

    if errno := C.snd_pcm_hw_params_malloc(&self.params); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not allocate hardware parameters. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_any(self.capture, self.params); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not initialize hardware parameter structure. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_set_access(self.capture, self.params, C.SND_PCM_ACCESS_RW_INTERLEAVED); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set access. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_set_format(self.capture, self.params, C.SND_PCM_FORMAT_FLOAT); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set sample format. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_set_rate(self.capture, self.params, self.header.SampleRate, 0); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set sample rate. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_set_channels(self.capture, self.params, 1); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set channel count. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params(self.capture, self.params); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set parameters. Error %d", errno) )
    }

    C.snd_pcm_hw_params_free(self.params)

    if errno := C.snd_pcm_prepare(self.capture); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not prepare audio interface for use. Error %d", errno) )
    }

    return nil
}
