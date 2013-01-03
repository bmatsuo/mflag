package mflag

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

var errInvalidArgument = errors.New("invalid arguments")

type errMissingFlagValue string // flag valued
type errUnknownFlag string      // flag valued
type errInvalidFlagValue struct{ Flag, Value string }

func (err errMissingFlagValue) Error() string {
	return "missing value: " + string(err)
}
func (err errUnknownFlag) Error() string {
	return "unknown flag: " + string(err)
}
func (err errInvalidFlagValue) Error() string {
	return fmt.Sprintf("invalid value: %s (%s)", err.Value, err.Flag)
}

var DefaultFlagSet = NewFlagSet(os.Args[0])

// Parse command line arguments set values in dst.
func Unmarshal(dst interface{}) ([]string, error) {
	if DefaultFlagSet == nil {
		panic("nil DefaultFlagSet")
	}
	_, err := DefaultFlagSet.Parse(os.Args[1:], dst)
	return DefaultFlagSet.Args(), err
}

type FlagSet struct {
	name string
	v    interface{}
	m    FlagMode
	t    *flagTree
	s    *flagScanner
}

func NewFlagSet(name string) *FlagSet {
	fs := &FlagSet{
		name: name,
	}
	return fs
}

func (fs *FlagSet) Mode(m FlagMode) *FlagSet {
	fs.m = m
	return fs
}

func (fs *FlagSet) Parse(args []string, dst interface{}) (*FlagSet, error) {
	fs.v = dst
	t, err := newFlagTree(fs.v)
	if err != nil {
		return fs, err
	}
	fs.t = t

	fs.s = newFlagScanner(fs.m, args)
	// This shit needs to be abstracted better
	for flag := fs.s.Next(); flag != ""; flag = fs.s.Next() {
		name, value, err := fs.nextFlag()
		if err != nil {
			return fs, err
		}
		err = fs.t.Assign(name, value)
		if err != nil {
			return fs, err
		}
	}
	return fs, nil
}

func (fs *FlagSet) nextFlag() (name, value string, err error) {
	flag := fs.s.Last()
	name = strings.TrimLeft(flag, "-")
	valsep := strings.Index(name, "=")
	switch valsep {
	case -1: // read a value from the arguments
		var tree *flagTree

		// bools are a special case
		if tree, err = fs.t.Lookup(name); err != nil {
			err = errUnknownFlag(flag)
			break
		}
		if tree.o.kind == reflect.Bool {
			value = "true"
			break
		}

		if value, err = fs.s.GetValue(); err == errNoFlagValue {
			err = errMissingFlagValue(flag)
			break
		}
	case 0:
		err = errors.New("missing flag: -=")
	default:
		name = name[:valsep]
		value = name[valsep+1:]
	}
	return
}

func (fs *FlagSet) NArg() int {
	if fs.s == nil {
		panic("unparsed arguments")
	}
	return len(fs.s.args)
}
func (fs *FlagSet) Args() []string {
	if fs.s == nil {
		panic("unparsed arguments")
	}
	return append(make([]string, 0, len(fs.s.args)), fs.s.args...)
}
func (fs *FlagSet) Arg(i int) string {
	if fs.s == nil {
		panic("unparsed arguments")
	}
	return fs.s.args[i]
}
