package mflag

import (
	"reflect"
	"testing"
)

type testFlagTreeSource struct {
	Foo struct {
		Bar *struct {
			Baz int
		}
	}
	X int32
	S []string
}

func TestFlagTreeSliceAppend(t *testing.T) {
	source := new(testFlagTreeSource)
	tree, err := newFlagTree(source)
	if err != nil {
		t.Errorf("initialization error: %v", err)
	}
	if err = tree.Assign("S", "1"); err != nil {
		t.Errorf("append error: %v", err)
	}
	if err = tree.Assign("S", "2"); err != nil {
		t.Errorf("append error: %v", err)
	}
	if err = tree.Assign("S", "3"); err != nil {
		t.Errorf("append error: %v", err)
	}
	if len(source.S) != 3 {
		t.Errorf("length error: %v != %v", err, len(source.S), 3)
	}
	if source.S[0] != "1" {
		t.Errorf("value error: %v != %v", err, source.S[0], "1")
	}
	if source.S[1] != "2" {
		t.Errorf("value error: %v != %v", err, source.S[1], "2")
	}
	if source.S[2] != "3" {
		t.Errorf("value error: %v != %v", err, source.S[2], "3")
	}
}

type testFlagTreeSourceFn func(*testFlagTreeSource) bool

func TestFlagTreeAssign(t *testing.T) {
	for i, test := range []*struct {
		err   bool
		flag  string
		value string
		fn    testFlagTreeSourceFn
	}{
		{
			false, "X", "10",
			func(tree *testFlagTreeSource) bool { return tree.X == 10 },
		},
		{
			false, "Foo.Bar.Baz", "-10",
			func(tree *testFlagTreeSource) bool { return tree.Foo.Bar.Baz == -10 },
		},
		{
			false, "S", "abc",
			func(tree *testFlagTreeSource) bool { return len(tree.S) > 0 && tree.S[0] == "abc" },
		},
	} {
		t.Logf("[%d] %s=%s (error: %v)", i, test.flag, test.value, test.err)
		source := new(testFlagTreeSource)
		tree, err := newFlagTree(source) // XXX tested elsewhere
		if err != nil {
			t.Fatalf("[%d] tree initialization: %v", err)
		}
		err = tree.Assign(test.flag, test.value)
		if err != nil && !test.err {
			t.Errorf("[%d] %s=%s unexpected error: %v", i, test.flag, test.value, err)
		} else if err == nil && test.err {
			t.Errorf("[%d] %s=%s expected error: %v", i, test.flag, test.value)
		} else if err == nil {
			if !test.fn(source) {
				t.Errorf("[%d] %s=%s test function failure", i, test.flag, test.value)
			}
		}
	}
}

func TestFlagTreeLookup(t *testing.T) {
	for i, test := range []*struct {
		err  bool
		flag string
		kind reflect.Kind
		typ  string
		v    interface{}
	}{
		{
			err:  false,
			flag: "X",
			kind: reflect.Int32,
			typ:  "int32",
			v:    new(testFlagTreeSource),
		},
		{
			err:  false,
			flag: "S",
			kind: reflect.Slice,
			typ:  "[]string",
			v:    new(testFlagTreeSource),
		},
		{
			err:  false,
			flag: "Foo.Bar.Baz",
			kind: reflect.Int,
			typ:  "int",
			v:    new(testFlagTreeSource),
		},
	} {
		tree, err := newFlagTree(test.v)
		if err != nil {
			t.Errorf("[%d] error creating tree: %v", i, err)
		}
		subtree, err := tree.Lookup(test.flag)
		if err != nil && !test.err {
			t.Errorf("[%d] unexpected error: %v", i, err)
		} else if err == nil && test.err {
			t.Errorf("[%d] expected error", i)
		} else if err == nil {
			if subtree.kind != test.kind {
				t.Errorf("[%d] lookup has kind %v (expected %v)", i, subtree.kind, test.kind)
			}
			if subtree.typ.String() != test.typ {
				t.Errorf("[%d] lookup has type %v (expected %v)", i, subtree.typ, test.typ)
			}
		}
	}
}
