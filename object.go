package mflag

import (
	"reflect"
	"strconv"
)

type flagObject struct {
	prefix string        // dot-separated prefix XXX unused
	name   string        // flag string
	kind   reflect.Kind  // the kind of v
	typ    reflect.Type  // generally the type of v. can be the elem type of v
	v      reflect.Value // a value to be assigned to.
}

func (o *flagObject) String() string {
	if o.prefix == "" {
		return o.name
	}
	return o.prefix + "." + o.name
}

func (o *flagObject) Assign(val string) error {
	var v interface{}
	var err error
	switch o.kind {
	// TODO map?
	case reflect.String:
		v = val
	case reflect.Bool:
		v, err = strconv.ParseBool(val)
	case reflect.Float32:
		v, err = strconv.ParseFloat(val, 32)
		v = float32(v.(float64))
	case reflect.Float64:
		v, err = strconv.ParseFloat(val, 64)
	case reflect.Uint:
		v, err = strconv.ParseUint(val, 10, 0)
		v = uint(v.(uint64))
	case reflect.Uint8:
		v, err = strconv.ParseUint(val, 10, 8)
		v = uint8(v.(uint64))
	case reflect.Uint16:
		v, err = strconv.ParseUint(val, 10, 16)
		v = uint16(v.(uint64))
	case reflect.Uint32:
		v, err = strconv.ParseUint(val, 10, 32)
		v = uint32(v.(uint64))
	case reflect.Uint64:
		v, err = strconv.ParseUint(val, 10, 64)
	case reflect.Int:
		v, err = strconv.ParseInt(val, 10, 0)
		v = int(v.(int64))
	case reflect.Int8:
		v, err = strconv.ParseInt(val, 10, 8)
		v = int8(v.(int64))
	case reflect.Int16:
		v, err = strconv.ParseInt(val, 10, 16)
		v = int16(v.(int64))
	case reflect.Int32:
		v, err = strconv.ParseInt(val, 10, 32)
		v = int32(v.(int64))
	case reflect.Int64:
		v, err = strconv.ParseInt(val, 10, 64)
	case reflect.Slice:
		// here o.typ is o.v.Type().Elem()
		_o := &flagObject{
			prefix: o.prefix,
			name:   o.name,
			kind:   o.typ.Elem().Kind(),
			typ:    o.typ.Elem(),
			v:      reflect.Indirect(reflect.New(o.typ.Elem())),
		}
		err = _o.Assign(val)
		v = _o.v
	default:
		err = errInvalidArgument
	}
	if err != nil {
		return err
	}

	switch o.kind {
	case reflect.Slice:
		o.v.Set(reflect.Append(o.v, v.(reflect.Value)))
	default:
		o.v.Set(reflect.ValueOf(v))
	}
	return nil
}
