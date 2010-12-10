// Copyright (c) 2010 AFP Authors

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

func NewAlsaSource() afp.Filter {
    return &AlsaSource{}
}

func (self *AlsaSource) GetType() int {
    return afp.PIPE_SOURCE
}

func (self *AlsaSource) Init(ctx *afp.Context, args []string) os.Error {
    self.ctx = ctx
    return nil
}

func (self *AlsaSource) Start() {

    self.header = afp.StreamHeader {
        Version: 1,
        Channels: 1,
        SampleSize: 32,
        SampleRate: 44100,
        FrameSize: 4096,
    }

    self.ctx.HeaderSink <- self.header

    retval := self.prepare()
    if ( retval != nil) {
        panic(retval)
    }

    for {
		cbuf := make([]float32, int32(self.header.Channels) * self.header.FrameSize)
		buff := make([][]float32, self.header.FrameSize)
        length := len(cbuf)

        //first off, grab some data from alsa
        read := C.snd_pcm_readi(self.capture, unsafe.Pointer(&cbuf[0]), C.snd_pcm_uframes_t(length))
        if read < C.snd_pcm_sframes_t(length) {
            errno := C.snd_pcm_recover(self.capture, C.int(read), 0)
            if errno < 0 {
                panic(fmt.Sprint( "While reading from ALSA device, failed to recover from error: ", errno))
            }
        }

        // snd_pcm_readi gives us a one dimensional array of interleaved data
        // but what we want is a two dimensional array of samples
        chans := int(self.header.Channels)
		for slice, i := 0, 0; i < length; slice, i = slice + 1, i + chans {
			buff[slice] = make([]float32, chans)
			buff[slice] = cbuf[i : i + chans]
        }

        //send it on down the line
        self.ctx.Sink <- buff
    }

}

func (self *AlsaSource) Stop() os.Error {
    C.snd_pcm_close(self.capture)
    close(self.ctx.Sink)
    return nil
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

func NewAlsaSink() afp.Filter {
    return &AlsaSink{}
}

func (self *AlsaSink) GetType() int {
    return afp.PIPE_SINK
}

func (self *AlsaSink) Init(ctx *afp.Context, args []string) os.Error {
    self.ctx = ctx
    return nil
}

func(self *AlsaSink) Start() {

    const FRAMES int = 256 //the double buffer holds 256 frames
    self.header = <-self.ctx.HeaderSource

    retval := self.prepare()
    if (retval != nil) {
        panic(retval)
    }

    //almost a do..while
    var cbuf []float32 //C buffer
    var written chan C.snd_pcm_sframes_t = make(chan C.snd_pcm_sframes_t, 2)
    chans := int(self.header.Channels)
    double := make([][]float32, 0, self.header.FrameSize * 1024)
    double = append(double, <-self.ctx.Source...)
    length := len(double)
    oldLength := length
    cbuf = make([]float32, int32(self.header.Channels) * self.header.FrameSize * int32(length))
    streamOffset := 0
    for i := 0; i < length; i++ {
        for j := 0; j < chans; j++ {
            cbuf[streamOffset] = double[i][j]
            streamOffset++
        }
    }
    written <- C.snd_pcm_sframes_t(length)

    for buffer := range self.ctx.Source { //blocking

        select {

        //wait for the previous write to finish
        case error := <-written:

            //fmt.Printf(".")

            if int(error) < oldLength {
                //not all the data was written to the device
                panic(fmt.Sprintf("Could not write all data to ALSA device, wrote: ", written))
            }

            //write to the speaker in another thread
            go func() {
                oldLength = length
                written <- C.snd_pcm_writei(self.playback, unsafe.Pointer(&cbuf[0]), C.snd_pcm_uframes_t(length))
            }()

            double = buffer

        default:

            //fmt.Print("-")
            double = append(double, buffer...)
            length = len(double)

            //cbuf WILL be a new pointer, so we can gaurantee that the address space given to snd_pcm_writei
            //will not be modified once it is passed; but on the down side, we rely on the garbage collector
            //to take care of many stale slices created while filling our buffer
            cbuf = make([]float32, int32(self.header.Channels) * self.header.FrameSize * int32(length))

            //interleave the channels
            streamOffset = 0
            for i := 0; i < length; i++ {
                for j := 0; j < chans; j++ {
                    cbuf[streamOffset] = double[i][j]
                    streamOffset++
                }
            }

        }//end select

    }

    return
}

func (self *AlsaSink) Stop() os.Error {
    close(self.ctx.Source)
    C.snd_pcm_close(self.playback)
    return nil
}

// Ugly bastardized C code follows
func (self *AlsaSink) prepare() os.Error {

    if errno := C.snd_pcm_open(&self.playback, C.CString("default"), C.SND_PCM_STREAM_PLAYBACK, 0); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not open device. Error %d", errno) )
    }

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

    if errno := C.snd_pcm_hw_params_set_rate(self.playback, self.params, C.uint(self.header.SampleRate), 0); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set sample rate. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params_set_channels(self.playback, self.params, C.uint(self.header.Channels)); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set channel count. Error %d", errno) )
    }

    if errno := C.snd_pcm_hw_params(self.playback, self.params); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not set hardware parameters. Error %d", errno) )
    }

    C.snd_pcm_hw_params_free(self.params)

    if errno := C.snd_pcm_prepare(self.playback); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not prepare audio device for use. Error %d", errno) )
    }

    return nil
}

//this one is slightly different
//note the change in scope
func (self *AlsaSource) prepare() os.Error {

    if errno := C.snd_pcm_open(&self.capture, C.CString("default"), C.SND_PCM_STREAM_CAPTURE, 0); errno < 0 {
        return os.NewError( fmt.Sprintf("Could not open device. Error %d", errno) )
    }

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

    if errno := C.snd_pcm_hw_params_set_rate(self.capture, self.params, C.uint(self.header.SampleRate), 0); errno < 0 {
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
