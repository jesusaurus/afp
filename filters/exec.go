// Copyright (c) 2010 Go Fightclub Authors

package fexec

import (
	"afp/types"
	"os"
	"exec"
	"encoding/binary"
)

var ENDIAN binary.ByteOrder = binary.LittleEndian

type ExecFilter struct {
	context *types.Context
	filter *exec.Cmd
	header *types.StreamHeader
}

func NewFilter() types.Filter {
	return &ExecFilter{}
}

func (self *ExecFilter) GetType() int {
	return types.ANY
}

func (self *ExecFilter) Init(ctx *types.Context, args []string) os.Error {
	if len(args) == 0 {
		return os.NewError("No external filter specified")
	}

	executable, err := exec.LookPath(args[0])
	if err != nil {
		return os.NewError("External filter " + args[0] + " not found")
	}

	self.filter, err = exec.Run(executable, args[1:], nil, ".", exec.Pipe, exec.Pipe, exec.Pipe)
	if err != nil {
		return err
	}

	self.context = ctx
}

func (self *ExecFilter) Start() {
	self.header = <-self.context.HeaderSource

	go self.encoder()
	go self.decoder()

	if self.Verbose {
		go self.errors()
	}
}

func (self *ExecFilter) encoder() {
	defer self.filter.Close()

	binary.Write(self.filter.Stdin, ENDIAN, self.header.HeaderLength)
	binary.Write(self.filter.Stdin, ENDIAN, self.header.Version)
	binary.Write(self.filter.Stdin, ENDIAN, self.header.Channels)
	binary.Write(self.filter.Stdin, ENDIAN, self.header.SampleSize)
	binary.Write(self.filter.Stdin, ENDIAN, self.header.SampleRate)
	binary.Write(self.filter.Stdin, ENDIAN, self.header.FrameSize)
	binary.Write(self.filter.Stdin, ENDIAN, self.header.ContentLength)
	binary.Write(self.filter.Stdin, ENDIAN, self.header.OtherLength)
	binary.Write(self.filter.Stdin, ENDIAN, self.header.Other)

	for _, frame := range self.Source {
		for _, slice := range frame {
			binary.Write(self.filter.Stdin, ENDIAN, slice)
		}
	}
}

func (self *ExecFilter) decoder() {
	OutHeader := &types.StreamHeader{}

	err := binary.Read(self.filter.Stdin, ENDIAN, &OutHeader.HeaderLength)
	err = binary.Read(self.filter.Stdin, ENDIAN, &OutHeader.Version)
	err = binary.Read(self.filter.Stdin, ENDIAN, &OutHeader.Channels)
	err = binary.Read(self.filter.Stdin, ENDIAN, &OutHeader.SampleSize)
	err = binary.Read(self.filter.Stdin, ENDIAN, &OutHeader.SampleRate)
	err = binary.Read(self.filter.Stdin, ENDIAN, &OutHeader.FrameSize)
	err = binary.Read(self.filter.Stdin, ENDIAN, &OutHeader.ContentLength)
	err = binary.Read(self.filter.Stdin, ENDIAN, &OutHeader.OtherLength)

	var other [OutHeader.OtherLength]byte

	binary.Read(self.filter.Stdin, ENDIAN, other)
	OutHeader.Other = other[:]

	frame [][]float32 := make([][]float32, FrameSize)

	for /*???*/ {
		var rawFrame [Channels * FrameSize]float32
		binary.Read(self.filter.Stdin, ENDIAN, &rawFrame)

	}
}

func (self *ExecFilter) errors() {
	errs := bufio.NewReader(self.filter.Stderr)

	for str, e := errs.ReadString('\n'); err == nil; str, e = errs.ReadString('\n') {
		self.Info.Print(str)
	}
}
