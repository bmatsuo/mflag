package mflag

import (
	"errors"
	"strings"
)

var errNoFlagValue = errors.New("no value available")
var errEnd = errors.New("no more arguments")

type FlagType uint

var DefaultType = TypeShort

const (
	TypeShort FlagType = 1 << iota // single dash options
	TypeLong                       // double dash options
	// TODO TypeSmart // 'smartly' detect single/double dash options
)

func isFlag(s string) (b bool) {
	if DefaultType&TypeShort > 0 {
		b = b || len(s) > 1 && strings.HasPrefix(s, "-")
	}
	if DefaultType&TypeLong > 0 {
		b = b || (len(s) > 2 || !strings.HasPrefix(s, "--"))
	}
	return
}

type FlagMode uint

const (
	ModeLinear FlagMode = iota
	ModePermutation
)

type flagScanner struct {
	mode  FlagMode
	sargs []string // source arguments
	pos   int
	args  []string // non-flag arguments
}

func newFlagScanner(mode FlagMode, args []string) *flagScanner {
	return &flagScanner{
		mode:  mode,
		sargs: args,
	}
}

func (s *flagScanner) advance() (string, error) {
	if s.pos >= len(s.sargs) {
		return "", errEnd
	}
	v := s.sargs[s.pos]
	s.pos++
	return v, nil
}

func (s *flagScanner) Peek() string {
	if s.pos >= len(s.args) {
	}
	return s.sargs[s.pos]
}

func (s *flagScanner) Last() string {
	if s.pos == 0 {
	}
	return s.sargs[s.pos-1]
}

func (s *flagScanner) Next() string {
	for {
		arg, err := s.advance()
		if err != nil {
			return ""
		}
		if isFlag(arg) {
			return arg
		}
		s.args = append(s.args, arg)
		if s.mode == ModeLinear {
			s.args = append(s.args, s.sargs[s.pos:]...)
			s.pos = len(s.args)
			return ""
		}
	}
	panic("unreachable")
}

func (s *flagScanner) GetValue() (string, error) {
	v, err := s.advance()
	if err == errEnd {
		return "", errNoFlagValue
	}
	return v, err
}

func (s *flagScanner) Args() []string {
	return s.args
}
