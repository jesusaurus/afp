package libavfilter

import (
	"afp"
	"afp/flags"
	"libav"
	"os"
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

	streamInfo := libav.StreamInfo(self.dctx)

	self.actx.HeaderSink <- afp.StreamHeader{
		Version : 1,
		Channels : int8(streamInfo.Channels),
		SampleSize : int8(streamInfo.SampleSize),
		SampleRate : streamInfo.SampleRate,
		FrameSize : streamInfo.FrameSize,
		ContentLength : 0,
	}

}

func (self *LibAVSource) Stop() os.Error {
	close(self.actx.Sink)
	return nil
}
