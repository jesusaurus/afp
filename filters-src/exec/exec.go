// Copyright (c) 2010 Go Fightclub Authors
// This source code is released under the terms of the
// MIT license. Please see the file LICENSE for license details.

package fexec

import (
	"afp"
	"afp/flags"
	"os"
	"exec"
	"encoding/binary"
)

type execFilter struct {
	context *types.Context
	filter *exec.Cmd
	header *types.StreamHeader
	endianness binary.ByteOrder
	commErrors chan os.Error
	finished chan int	
}

func (self *execFilter) Stop() os.Error {
	self.filter.Close(os.WNOHANG)
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

	self.context = ctx
	self.commErrors = make(chan os.Error)
	self.finished = make(chan int) 
	return nil
}

func (self *execFilter) write(v interface{}) {
	err := binary.Write(self.filter.Stdin, self.endianness, v)
	if err != nil {
		self.commErrors <- err
		select{}
	}
}

func (self *execFilter) encoder() {
	defer self.filter.Close()

	self.write(self.header.HeaderLength)
	self.write(self.header.Version)
	self.write(self.header.Channels)
	self.write(self.header.SampleSize)
	self.write(self.header.SampleRate)
	self.write(self.header.FrameSize)
	self.write(self.header.ContentLength)
	self.write(self.header.OtherLength)
	self.write(self.header.Other)

	for _, frame := range self.ctx.Source {
		for _, slice := range frame {
			self.write(self.filter.Stdin, self.endianness, slice)
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
	OutHeader := &afp.StreamHeader{}

	self.read(&OutHeader.HeaderLength)
	self.read(&OutHeader.Version)
	self.read(&OutHeader.Channels)
	self.read(&OutHeader.SampleSize)
	self.read(&OutHeader.SampleRate)
	self.read(&OutHeader.FrameSize)
	self.read(&OutHeader.ContentLength)
	self.read(&OutHeader.OtherLength)

	var other [OutHeader.OtherLength]byte

	self.read(&other)
	OutHeader.Other = other[:]

	self.context.HeaderSink <- OutHeader

	frame [][]float32 := make([][]float32, FrameSize)

	for {
		var rawFrame [OutHeader.Channels * OutHeader.FrameSize]float32
		err := self.read(&rawFrame)
		for i, slice := 0, 0; i < OutHeader.FrameSize / OutHeader.Channels; slice++ {
			frame[slice] = rawFrame[i:i + OutHeader.Channels]
			i +=  OutHeader.Channels
		}
	}
}

func (self *ExecFilter) errors() {
	errs := bufio.NewReader(self.filter.Stderr)

	for str, e := errs.ReadString('\n'); err == nil; str, e = errs.ReadString('\n') {
		self.context.info.Print(str)
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