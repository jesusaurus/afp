package libavfilter


import (
	"afp"
	"afp/flags"
	"afp/libav"
	"os"
	"unsafe"
)


/**
 * Filter: LibAVSource
 *
 * Flags:
 * -i inputFile
 *
 * We carry a main afp.Context as all filters
 * And also an AVDecodeContext for interfacing with libAV
 * the streamInfo contains general information about the input file; but will only work for MP3 files :(
 * floatSamples contain the input data converted into -1.0 to 1.0 float32 data arranged in 
 *   a 2D array indexed by sample then bychannel 
 * currBuffer holds the index of the current buffer
 * inFile is the path of the input file
 */

type LibAVSource struct {
	actx *afp.Context
	dctx libav.AVDecodeContext
	streamInfo libav.AVStreamInfo
	floatSamples [][][]float32
	currBuffer int8
	
	inFile string
}


/**
 * initialize the filter storage
 */

func NewLibAVSource() afp.Filter {
	return &LibAVSource{}
}


/**
 * find our input file and set up libav structures
 */

func (self *LibAVSource) Init(ctx *afp.Context, args []string) os.Error {
	self.actx = ctx
println(args[0],args[1])
	parser := flags.FlagParser(args)
	var i *string = parser.String("i", "", "The input file")
	parser.Parse()
println(i, *i)
	if (i != nil) {
		self.inFile = *i;
	} else {
		return os.NewError("Please specify an input file, good sir")
	}
	
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
	
	return nil
}


/**
 * LibAVSource is unsurprisingly a source
 */

func (self *LibAVSource) GetType() int {
	return afp.PIPE_SOURCE
}


/**
 * send the StreamHeader down the pipe,
 * then successively decode each packet in the stream until there is no more data
 * bouncing the decoded data between the two main floatSamples buffers
 */

func (self *LibAVSource) Start() {
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

/**
 * convert intSamples into floatSamples (-1.0 .. 1.0)
 */
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

/**
 * close the Sink
 */
func (self *LibAVSource) Stop() os.Error {
	close(self.actx.Sink)
	return nil
}
