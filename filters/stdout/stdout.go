package stdout;

import (
	"afp"
	"os"
	"encoding/binary"
)

type StdoutSink struct {
	context *afp.Context
	header afp.StreamHeader
}

func NewStdoutSink() afp.Filter {
	return &StdoutSink{}
}

func (self *StdoutSink) Init(ctx *afp.Context, args []string) os.Error {
    self.context = ctx
    self.header = <-self.context.HeaderSource

	return nil
}

func (self *StdoutSink) GetType() int {
    return afp.PIPE_SINK
}

func (self *StdoutSink) Start() {
    for buffer := range self.context.Source { //reading a [][]float32
        ibuf := make([]int16, int32(self.header.Channels) * self.header.FrameSize)
        length := len(buffer)
        chans := int(self.header.Channels)

		t := 0
        //interleave the channels
        for i := 0; i < length; i += chans {
            for j := 0; j < chans; j++ {
                ibuf[i+j] = int16(buffer[t][j] * 32767.0)
            }
			t += 1
        }

		binary.Write(os.Stdout, binary.LittleEndian, ibuf)
		
        //write some data to alsa
/*        written := C.snd_pcm_writei(self.playback, unsafe.Pointer(&cbuf[0]), C.snd_pcm_uframes_t(length))

        if int(written) < length {
            //not all the data was written
            panic(fmt.Sprintf("Could not write all data to ALSA device, wrote: ", written))
        }
*/    }

    return
}

func (self *StdoutSink) Stop() os.Error {
    return nil
}

