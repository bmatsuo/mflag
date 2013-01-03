package mflag

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

type flagTreeIndex struct {
	i int
	t *flagTree
}

type flagTree struct {
	indirect bool
	prefix   string
	name     string
	kind     reflect.Kind
	typ      reflect.Type
	ts       map[string]flagTreeIndex
	o        *flagObject
}

func newFlagTree(v interface{}) (*flagTree, error) {
	val := reflect.ValueOf(v)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("flag tree root must be a pointer; found %v", typ.Kind())
	}
	elem := typ.Elem()
	if elem.Kind() != reflect.Struct {
		return nil, fmt.Errorf("flag tree root must be pointer to a struct; found %v", elem.Kind())
	}
	tree, err := _newFlagTree(elem, "", "", true)
	if err != nil {
		return nil, err
	}
	tree.o = &flagObject{
		kind: tree.kind,
		typ:  tree.typ,
		v:    reflect.Indirect(val),
	}
	return tree, nil
}

func _newFlagTree(typ reflect.Type, prefix string, name string, indirect bool) (*flagTree, error) {
	t := &flagTree{
		indirect: indirect,
		prefix:   prefix,
		name:     name,
		kind:     typ.Kind(),
		typ:      typ,
	}
	if t.kind == reflect.Struct {
		numField := t.typ.NumField()
		t.ts = make(map[string]flagTreeIndex, numField)
		for i, n := 0, numField; i < n; i++ {
			ftyp := t.typ.Field(i)
			_ftyp := ftyp.Type

			fprefix := name
			if prefix != "" {
				fprefix = prefix + "." + name
			}

			fname := ftyp.Name

			if ftag := ftyp.Tag.Get("mflag"); ftag != "" {
				ftags := strings.Split(ftag, ",")

				fnametag := strings.TrimFunc(ftags[0], unicode.IsSpace)
				if fnametag != "" {
					fname = fnametag
				}
			}

			// XXX bugs here?
			findirect := false
			if _ftyp.Kind() == reflect.Ptr {
				findirect = true
				_ftyp = _ftyp.Elem()
			}
			tt, err := _newFlagTree(_ftyp, fprefix, fname, findirect)
			if err != nil {
				return nil, err
			}
			t.ts[fname] = flagTreeIndex{i, tt}
		}
	}
	return t, nil
}

func (t *flagTree) String() string {
	if t.prefix != "" {
		return t.prefix + "." + t.name
	}
	return t.name
}

func (t *flagTree) Assign(flag, value string) error {
	ft, err := t.Lookup(flag)
	if err != nil {
		return errUnknownFlag(flag)
	}
	if err := ft.o.Assign(value); err != nil {
		return errInvalidFlagValue{flag, value}
	}
	return nil
}

func (t *flagTree) Lookup(flag string) (*flagTree, error) {
	if flag == t.name {
		return t, nil
	}

	fflag := flag
	if t.name != "" {
		if !strings.HasPrefix(flag, t.name+".") {
			return nil, fmt.Errorf("%v %#v impossible flag: %v", t.typ, t, flag)
		}
		fflag = fflag[len(t.name)+1:]
	}

	fname := fflag
	if i := strings.Index(fname, "."); i > 0 {
		fname = fname[:i]
	}

	if t.ts != nil {
		ti, ok := t.ts[fname]
		if !ok {
			return nil, fmt.Errorf("%v %#v non-existent sub flag: %v", t.typ, t, flag)
		}
		if ti.t.o == nil {
			fval := t.o.v.Field(ti.i)
			ti.t.o = &flagObject{
				prefix: ti.t.prefix,
				name:   ti.t.name,
				kind:   ti.t.kind,
				typ:    ti.t.typ,
				v:      fval,
			}
			if ti.t.indirect {
				if fval.IsNil() {
					fval.Set(reflect.New(ti.t.typ))
				}
				ti.t.o.v = reflect.Indirect(ti.t.o.v)
			}
		}
		return ti.t.Lookup(fflag)
	}
	return nil, fmt.Errorf("%v %#v non-existent flag: %v", t.typ, t, flag)
}
