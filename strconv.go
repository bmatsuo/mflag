// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mflag

/*  Filename:    strconv.go
 *  Author:      Bryan Matsuo <bryan.matsuo@gmail.com>
 *  Created:     Sun Oct 23 23:14:13 PDT 2011
 *  Description: 
 */

import (
	"errors"
	"reflect"
	"strconv"
)

//  Stringify the value in v.
func valueToString(v reflect.Value) string {
	switch k := v.Kind(); k {
	case reflect.String:
		return v.Interface().(string)
	case reflect.Bool:
		return strconv.Btoa(v.Interface().(bool))
	case reflect.Int:
		return strconv.Itoa64(int64(v.Interface().(int)))
	case reflect.Int8:
		return strconv.Itoa64(int64(v.Interface().(int8)))
	case reflect.Int16:
		return strconv.Itoa64(int64(v.Interface().(int16)))
	case reflect.Int32:
		return strconv.Itoa64(int64(v.Interface().(int32)))
	case reflect.Int64:
		return strconv.Itoa64(int64(v.Interface().(int64)))
	case reflect.Uint:
		return strconv.Uitoa64(uint64(v.Interface().(uint)))
	case reflect.Uint8:
		return strconv.Uitoa64(uint64(v.Interface().(uint8)))
	case reflect.Uint16:
		return strconv.Uitoa64(uint64(v.Interface().(uint16)))
	case reflect.Uint32:
		return strconv.Uitoa64(uint64(v.Interface().(uint32)))
	case reflect.Uint64:
		return strconv.Uitoa64(uint64(v.Interface().(uint64)))
	case reflect.Float32:
		return strconv.Ftoa32(v.Interface().(float32), 'f', -1)
	case reflect.Float64:
		return strconv.Ftoa64(v.Interface().(float64), 'f', -1)
	default:
		panic("unsupported string conversion")
	}
	return ""
}

//  Assign to v a value from string.
func stringToValue(s string, v reflect.Value) error {
	if !v.CanSet() {
		return errors.New("Value is not assignable")
	}

	wide, err := stringToWideType(s, v.Type())
	if err != nil {
		return err
	}

	switch wide.Kind() {
	case reflect.String:
		v.SetString(wide.Interface().(string))
	case reflect.Bool:
		v.SetBool(wide.Interface().(bool))
	case reflect.Int64:
		x := wide.Interface().(int64)
		v.SetInt(x)
		if v.OverflowInt(x) {
			err = errors.New("Assignment precision error")
		}
	case reflect.Uint64:
		x := wide.Interface().(uint64)
		v.SetUint(x)
		if v.OverflowUint(x) {
			err = errors.New("Assignment precision error")
		}
	case reflect.Float64:
		x := wide.Interface().(float64)
		v.SetFloat(x)
		if v.OverflowFloat(x) {
			err = errors.New("Assignment precision error")
		}
	}
	return nil
}

func stringToWideType(s string, t reflect.Type) (v reflect.Value, err error) {
	handle := func(i interface{}, e error) {
		v = reflect.ValueOf(i)
		err = e
	}
	switch k := t.Kind(); k {
	case reflect.String:
		return reflect.ValueOf(s), nil
	case reflect.Bool:
		handle(strconv.Atob(s))
	case reflect.Int:
		handle(strconv.Atoi64(s))
	case reflect.Int8:
		handle(strconv.Atoi64(s))
	case reflect.Int16:
		handle(strconv.Atoi64(s))
	case reflect.Int32:
		handle(strconv.Atoi64(s))
	case reflect.Int64:
		handle(strconv.Atoi64(s))
	case reflect.Uint:
		handle(strconv.Atoui64(s))
	case reflect.Uint8:
		handle(strconv.Atoui64(s))
	case reflect.Uint16:
		handle(strconv.Atoui64(s))
	case reflect.Uint32:
		handle(strconv.Atoui64(s))
	case reflect.Uint64:
		handle(strconv.Atoui64(s))
	case reflect.Float32:
		handle(strconv.Atof64(s))
	case reflect.Float64:
		handle(strconv.Atof64(s))
	default:
		panic("unsupported string conversion")
	}
	if err != nil {
		panic(err)
	}
	return
}
