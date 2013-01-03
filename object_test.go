package mflag

import (
	"reflect"
	"testing"
)

func TestFlagObjectAssign(t *testing.T) {
	for i, test := range []*struct {
		err    bool
		in     string
		out    interface{}
		object *flagObject
	}{
		{
			false,
			"abc",
			"abc",
			&flagObject{
				v:    reflect.Indirect(reflect.ValueOf(new(string))),
				typ:  reflect.TypeOf(new(string)),
				kind: reflect.String,
			},
		},
		{
			false,
			"12",
			uint8(12),
			&flagObject{
				v:    reflect.Indirect(reflect.ValueOf(new(uint8))),
				typ:  reflect.TypeOf(uint8(0)),
				kind: reflect.Uint8,
			},
		},
		{
			false,
			"12",
			int16(12),
			&flagObject{
				v:    reflect.Indirect(reflect.ValueOf(new(int16))),
				typ:  reflect.TypeOf(int16(0)),
				kind: reflect.Int16,
			},
		},
		{
			false,
			"12",
			float32(12),
			&flagObject{
				v:    reflect.Indirect(reflect.ValueOf(new(float32))),
				typ:  reflect.TypeOf(float32(0)),
				kind: reflect.Float32,
			},
		},
	} {
		t.Logf("[%d] %v = %v (%v)", i, test.object.typ, test.in, test.err)
		err := test.object.Assign(test.in)
		if err != nil && !test.err {
			t.Errorf("[%d] %v = %v unexpected error: %v", i, test.object.typ, test.in, err)
		} else if err == nil && test.err {
			t.Errorf("[%d] %v = %v error expected", i, test.object.typ, test.in)
		} else if err == nil && !reflect.DeepEqual(test.out, test.object.v.Interface()) {
			t.Errorf("[%d] %v = %v => %v (expected %v)", i, test.object.typ, test.in, test.object.v.Interface(), test.out)
		}
	}
}
