package libavfilter

import (
	"afp"
	"afp/flags"
	"libav"
	"os"
	"unsafe"
)


/* 
 * LibAVSource
 *
 * Flags:
 * -i inputFile
 */

type LibAVSource struct {
	actx *afp.Context
	dctx libav.AVDecodeContext
	streamInfo libav.AVStreamInfo
	floatSamples [][][]float32
	currBuffer int8
	
	inFile string
}

func NewLibAVSource() afp.Filter {
	return &LibAVSource{}
}

func (self *LibAVSource) Init(ctx *afp.Context, args []string) os.Error {
	self.actx = ctx

	parser := flags.FlagParser(args)
	var i *string = parser.String("i", "", "The input file")
	parser.Parse()

	if (i != nil) {
		self.inFile = *i;
	} else {
		return os.NewError("Please specify an input file, good sir")
	}	

	return nil
}

func (self *LibAVSource) GetType() int {
	return afp.PIPE_SOURCE
}

func (self *LibAVSource) Start() {
	libav.InitDecoding()
	libav.PrepareDecoding(self.inFile, &self.dctx)

	self.streamInfo = libav.StreamInfo(self.dctx)

	/* initialize floatSamples buffer */
	self.floatSamples = make([][][]float32, 2)
	for i,_ := range(self.floatSamples) {
		self.floatSamples[i] = make([][]float32, self.streamInfo.FrameSize)

		for j,_ := range(self.floatSamples[i]) {
			self.floatSamples[i][j] = make([]float32, self.streamInfo.Channels)
		}
	}
	
	self.currBuffer = 0

	self.actx.HeaderSink <- afp.StreamHeader{
		Version : 1,
		Channels : int8(self.streamInfo.Channels),
		SampleSize : int8(self.streamInfo.SampleSize),
		SampleRate : self.streamInfo.SampleRate,
		FrameSize : self.streamInfo.FrameSize,
		ContentLength : 0,
	}
	
	l := int32(libav.DecodePacket(self.dctx))
	for l > 0 {
		numberOfSamples := l / self.streamInfo.SampleSize
		decodedSamples := (*(*[1 << 31 - 1]int16)(unsafe.Pointer(self.dctx.Context.Outbuf)))[:numberOfSamples]

		self.int16ToFloat32(decodedSamples)
		self.actx.Sink <- self.floatSamples[self.currBuffer]
		self.currBuffer = 1 - self.currBuffer
	}
}

func (self *LibAVSource) int16ToFloat32(intSamples []int16) {	
	var (
		streamOffset int32 = 0
		i int32
		j int32
	)
	
	for i = 0; i < int32(len(intSamples)); i+=self.streamInfo.Channels {
		for j = 0; j < self.streamInfo.Channels; j++ {
			self.floatSamples[self.currBuffer][i][j] = float32(intSamples[streamOffset]) / float32(1 << 31)
			streamOffset += 1
		}
	}
}

func (self *LibAVSource) Stop() os.Error {
	close(self.actx.Sink)
	return nil
}
