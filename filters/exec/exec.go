// Copyright (c) 2010 AFP Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package fexec

import (
	"afp"
	"os"
	"exec"
	"encoding/binary"
	"syscall"
	"bufio"
)

type execFilter struct {
	context *afp.Context
	filter *exec.Cmd
	header *afp.StreamHeader
	endianness binary.ByteOrder
	commErrors chan os.Error
	finished chan int
}

func (self *execFilter) Stop() os.Error {
	syscall.Kill(self.filter.Pid, syscall.SIGTERM)
	self.finished <- 1

	return nil
}

func (self *execFilter) Init(ctx *afp.Context, args []string) os.Error {
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
	
	self.endianness = binary.LittleEndian //Move this to a flag in the future 
	self.context = ctx
	self.commErrors = make(chan os.Error)
	self.finished = make(chan int)
	return nil
}

func (self *execFilter) write(v interface{}) {
	err := binary.Write(self.filter.Stdin, self.endianness, v)
	if err != nil {
		//On a write error, pass the error back to the calling thread
		//Then block permanently.
		//The afp infrastructure will receive the error and shutdown
		//The pipeline, killing the cmd in the process.
		self.commErrors <- err
		select{}
	}
}

func (self *execFilter) encoder() {
	defer self.filter.Close()

	self.write(afp.HEADER_LENGTH)
	self.write(self.header.Version)
	self.write(self.header.Channels)
	self.write(self.header.SampleSize)
	self.write(self.header.SampleRate)
	self.write(self.header.FrameSize)
	self.write(self.header.ContentLength)

	for frame := range self.context.Source {
		for _, slice := range frame {
			self.write(slice)
		}
	}
}

func (self *execFilter) read(v interface{}) {
	err := binary.Read(self.filter.Stdin, self.endianness, v)
	if err == os.EOF {
		self.commErrors <- err
		select{}
	}
}

func (self *execFilter) decoder() {
	OutHeader := afp.StreamHeader{}

	var length int32

	self.read(&length)
	if length != afp.HEADER_LENGTH {
		self.commErrors <- os.NewError("Header not the right length")
		select{}
	}
	self.read(&OutHeader.Version)
	self.read(&OutHeader.Channels)
	self.read(&OutHeader.SampleSize)
	self.read(&OutHeader.SampleRate)
	self.read(&OutHeader.FrameSize)
	self.read(&OutHeader.ContentLength)

	self.context.HeaderSink <- OutHeader

	frame := make([][]float32, OutHeader.FrameSize)

	chans := int32(OutHeader.Channels)

	for {
		rawFrame := make([]float32, chans * OutHeader.FrameSize)
		self.read(rawFrame)

		for i, slice := int32(0), 0; i < OutHeader.FrameSize / chans; slice++ {
			frame[slice] = rawFrame[i:i + chans]
			i += chans
		}
	}
}

func (self *execFilter) errors() {
	errs := bufio.NewReader(self.filter.Stderr)

	for str, _ := errs.ReadString('\n'); errs == nil; str, _ = errs.ReadString('\n') {
		self.context.Info.Print(str)
	}
}

func (self *execFilter) wait() {
	select {
	case <-self.finished :
		if self.context.Sink != nil {
			close(self.context.Sink)
		}
		self.filter.Wait(0)
		return
	case err := <-self.commErrors :
		syscall.Kill(self.filter.Pid, syscall.SIGTERM)
		self.filter.Wait(0)
		panic(err)
	}
}