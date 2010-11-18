// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Modifications Copyright 2010 Aaron DeVore
// MIT License

package flags_test

import (
	"./flags"
	"fmt"
	"testing"
)

func boolString(s string) string {
	if s == "0" {
		return "false"
	}
	return "true"
}

func TestEverything(t *testing.T) {
	m := make(map[string] *flags.Flag)
	args := []string{"command"}
	parser := flags.FlagParser(args)
	var (
		_  = parser.Bool("test_bool", false, "bool value")
		_  = parser.Int("test_int", 0, "int value")
		_  = parser.Int64("test_int64", 0, "int64 value")
		_  = parser.Uint("test_uint", 0, "uint value")
		_  = parser.Uint64("test_uint64", 0, "uint64 value")
		_  = parser.String("test_string", "0", "string value")
		_  = parser.Float("test_float", 0, "float value")
		_  = parser.Float("test_float64", 0, "float64 value")
	)
	desired := "0"
	visitor := func(f *flags.Flag) {
		if len(f.Name) > 5 && f.Name[0:5] == "test_" {
			m[f.Name] = f
			ok := false
			switch {
				case f.Value.String() == desired:
					ok = true
				case f.Name == "test_bool" && f.Value.String() == boolString(desired):
					ok = true
			}
			if !ok {
				t.Error("Visit: bad value", f.Value.String(), "for", f.Name)
			}
		}
	}
	parser.VisitAll(visitor)
	if len(m) != 8 {
		t.Error("VisitAll misses some flags")
		for k, v := range m {
			t.Log(k, *v)
		}
	}
}

func TestParse(t *testing.T) {
	extra := "one-extra-argument"
	args := []string{
		"a.out",
		"-bool",
		"-bool2=true",
		"--int", "22",
		"--int64", "23",
		"-uint", "24",
		"--uint64", "25",
		"-string", "hello",
		"--float", "3141.5",
		"-float64", "2718e28",
		extra,
	}
	parser := flags.FlagParser(args)
	boolFlag := parser.Bool("bool", false, "bool value")
	bool2Flag := parser.Bool("bool2", false, "bool2 value")
	intFlag := parser.Int("int", 0, "int value")
	int64Flag := parser.Int64("int64", 0, "int64 value")
	uintFlag := parser.Uint("uint", 0, "uint value")
	uint64Flag := parser.Uint64("uint64", 0, "uint64 value")
	stringFlag := parser.String("string", "0", "string value")
	floatFlag := parser.Float("float", 0, "float value")
	float64Flag := parser.Float("float64", 0, "float64 value")
	parser.Parse()
	if *boolFlag != true {
		t.Error("bool flag should be true, is ", *boolFlag)
	}
	if *bool2Flag != true {
		t.Error("bool2 flag should be true, is ", *bool2Flag)
	}
	if *intFlag != 22 {
		t.Error("int flag should be 22, is ", *intFlag)
	}
	if *int64Flag != 23 {
		t.Error("int64 flag should be 23, is ", *int64Flag)
	}
	if *uintFlag != 24 {
		t.Error("uint flag should be 24, is ", *uintFlag)
	}
	if *uint64Flag != 25 {
		t.Error("uint64 flag should be 25, is ", *uint64Flag)
	}
	if *stringFlag != "hello" {
		t.Error("string flag should be `hello`, is ", *stringFlag)
	}
	if *floatFlag != 3141.5 {
		t.Error("float flag should be 3141.5, is ", *floatFlag)
	}
	if *float64Flag != 2718e28 {
		t.Error("float64 flag should be 2718e28, is ", *float64Flag)
	}
	if len(parser.Args()) != 1 {
		t.Error("expected one argument, got", len(parser.Args()))
	} else if parser.Args()[0] != extra {
		t.Errorf("expected argument %q got %q", extra, parser.Args()[0])
	}
}

// Declare a user-defined flag.
type flagVar []string

func (f *flagVar) String() string {
	return fmt.Sprint([]string(*f))
}

func (f *flagVar) Set(value string) bool {
	n := make(flagVar, len(*f)+1)
	copy(n, *f)
	*f = n
	(*f)[len(*f)-1] = value
	return true
}

func TestUserDefined(t *testing.T) {
	var v flagVar
	parser := flags.FlagParser([]string{"a.out", "-v", "1", "-v", "2", "-v=3"})
	parser.Var(&v, "v", "usage")
	defer func() {
		if recover() != nil {
			t.Error("parse failed")
		}
	}()
	parser.Parse()
	if len(v) != 3 {
		t.Fatal("expected 3 args; got ", len(v))
	}
	expect := "[1 2 3]"
	if v.String() != expect {
		t.Errorf("expected value %q got %q", expect, v.String())
	}
}
