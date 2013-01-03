package mflag

import (
	"reflect"
	"testing"
)

type scannerTest struct {
	mode               FlagMode
	input, flags, args []string
	ops                func(*flagScanner) []string
}

func TestFlagScannerLinearFlagDetection(t *testing.T) {
	for i, test := range []scannerTest{
		{
			ModeLinear,
			[]string{"-a", "-b", "-c", "abc", "def", "ghi"},
			[]string{"-a", "-b", "-c", ""},
			[]string{"abc", "def", "ghi"},
			func(s *flagScanner) (flags []string) {
				flags = append(flags, s.Next())
				flags = append(flags, s.Next())
				flags = append(flags, s.Next())
				flags = append(flags, s.Next())
				return
			},
		},
		{
			ModeLinear,
			[]string{"-a", "-b", "-c", "abc", "def", "ghi"},
			[]string{"-a", "-b", "-c", "abc", "def", "ghi", ""},
			nil,
			func(s *flagScanner) (flags []string) {
				flags = append(flags, s.Next())
				flags = append(flags, s.Next())
				flags = append(flags, s.Next())
				if v, err := s.GetValue(); err == nil {
					flags = append(flags, v)
				}
				if v, err := s.GetValue(); err == nil {
					flags = append(flags, v)
				}
				if v, err := s.GetValue(); err == nil {
					flags = append(flags, v)
				}
				flags = append(flags, s.Next())
				return
			},
		},
		{
			ModePermutation,
			[]string{"-a", "abc", "def", "ghi", "-b", "-c"},
			[]string{"-a", "-b", "-c", ""},
			[]string{"abc", "def", "ghi"},
			func(s *flagScanner) (flags []string) {
				flags = append(flags, s.Next())
				flags = append(flags, s.Next())
				flags = append(flags, s.Next())
				flags = append(flags, s.Next())
				return
			},
		},
		{
			ModePermutation,
			[]string{"-a", "abc", "def", "ghi", "-b", "-c"},
			[]string{"-a", "abc", "-b", "-c", ""},
			[]string{"def", "ghi"},
			func(s *flagScanner) (flags []string) {
				flags = append(flags, s.Next())
				if v, err := s.GetValue(); err == nil {
					flags = append(flags, v)
				}
				flags = append(flags, s.Next())
				flags = append(flags, s.Next())
				flags = append(flags, s.Next())
				return
			},
		},
	} {
		s := newFlagScanner(test.mode, test.input)
		flags := test.ops(s)
		if !reflect.DeepEqual(test.flags, flags) {
			t.Errorf("[%d] flags %v not equal expected %v", i, flags, test.flags)
		}
		if !reflect.DeepEqual(test.args, s.Args()) {
			t.Errorf("[%d] args %v not equal expected %v", i, s.Args(), test.args)
		}
	}
}
