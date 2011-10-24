package mflag

/*  Filename:    strconv_test.go
 *  Author:      Bryan Matsuo <bryan.matsuo@gmail.com>
 *  Created:     Sun Oct 23 23:14:13 PDT 2011
 *  Description: For testing strconv.go
 */

import (
    "testing"
    "reflect"
    "os"
)

func ToStringTest(v interface{}, expected string, T *testing.T) {
    T.Logf("Testing conversion of %#v to string", v)
    output := valueToString(reflect.ValueOf(v))
    if output != expected {
        T.Errorf("Error: output='%s' expected='%s'", output, expected)
    } else {
        T.Log("ok")
    }
}

func TestConversionToString(T *testing.T) {
    ToStringTest("astring", "astring", T)

    ToStringTest(true, "true", T)
    ToStringTest(false, "false", T)

    ToStringTest(int(127), "127", T)
    ToStringTest(int8(127), "127", T)
    ToStringTest(int16(127), "127", T)
    ToStringTest(int32(127), "127", T)
    ToStringTest(int64(127), "127", T)

    ToStringTest(uint(200), "200", T)
    ToStringTest(uint8(200), "200", T)
    ToStringTest(uint16(200), "200", T)
    ToStringTest(uint32(200), "200", T)
    ToStringTest(uint64(200), "200", T)

    ToStringTest(float32(1.25), "1.25", T)
    ToStringTest(float64(1.25), "1.25", T)
}

func TestConversionFromString(T *testing.T) {
    var err os.Error

    var s string
    err = stringToValue("astring", reflect.ValueOf(&s).Elem())
    if err != nil {
        T.Errorf("Error converting string to string; %s", err.String())
    }
    if s != "astring" {
        T.Errorf("Failed to covert string to string")
    }

    var b bool
    err = stringToValue("true", reflect.ValueOf(&b).Elem())
    if err != nil {
        T.Errorf("Error converting string to bool; %s", err.String())
    }
    if b != true {
        T.Errorf("Failed to covert string to bool")
    }

    var i int
    err = stringToValue("123", reflect.ValueOf(&i).Elem())
    if err != nil {
        T.Errorf("Error converting string to int; %s", err.String())
    }
    if i != int(123) {
        T.Errorf("Failed to covert string to int")
    }

    var i8 int8
    err = stringToValue("123", reflect.ValueOf(&i8).Elem())
    if err != nil {
        T.Errorf("Error converting string to int8; %s", err.String())
    }
    if i8 != int8(123) {
        T.Errorf("Failed to covert string to int8")
    }

    var i16 int16
    err = stringToValue("123", reflect.ValueOf(&i16).Elem())
    if err != nil {
        T.Errorf("Error converting string to int16; %s", err.String())
    }
    if i16 != int16(123) {
        T.Errorf("Failed to covert string to int16")
    }

    var i32 int32
    err = stringToValue("123", reflect.ValueOf(&i32).Elem())
    if err != nil {
        T.Errorf("Error converting string to int32; %s", err.String())
    }
    if i32 != int32(123) {
        T.Errorf("Failed to covert string to int32")
    }

    var i64 int64
    err = stringToValue("123", reflect.ValueOf(&i64).Elem())
    if err != nil {
        T.Errorf("Error converting string to int64; %s", err.String())
    }
    if i64 != int64(123) {
        T.Errorf("Failed to covert string to int64")
    }

    var ui uint
    err = stringToValue("123", reflect.ValueOf(&ui).Elem())
    if err != nil {
        T.Errorf("Error converting string to uint; %s", err.String())
    }
    if ui != uint(123) {
        T.Errorf("Failed to covert string to uint")
    }

    var ui8 uint8
    err = stringToValue("123", reflect.ValueOf(&ui8).Elem())
    if err != nil {
        T.Errorf("Error converting string to uint8; %s", err.String())
    }
    if ui8 != uint8(123) {
        T.Errorf("Failed to covert string to uint8")
    }

    var ui16 uint16
    err = stringToValue("123", reflect.ValueOf(&ui16).Elem())
    if err != nil {
        T.Errorf("Error converting string to uint16; %s", err.String())
    }
    if ui16 != uint16(123) {
        T.Errorf("Failed to covert string to uint16")
    }

    var ui32 uint32
    err = stringToValue("123", reflect.ValueOf(&ui32).Elem())
    if err != nil {
        T.Errorf("Error converting string to uint32; %s", err.String())
    }
    if ui32 != uint32(123) {
        T.Errorf("Failed to covert string to uint32")
    }

    var ui64 uint64
    err = stringToValue("123", reflect.ValueOf(&ui64).Elem())
    if err != nil {
        T.Errorf("Error converting string to uint64; %s", err.String())
    }
    if ui64 != uint64(123) {
        T.Errorf("Failed to covert string to uint64")
    }

    var f32 float32
    err = stringToValue("1.25", reflect.ValueOf(&f32).Elem())
    if err != nil {
        T.Errorf("Error converting string to float32; %s", err.String())
    }
    if f32 != float32(1.25) {
        T.Errorf("Failed to covert string to float32")
    }

    var f64 float64
    err = stringToValue("1.25", reflect.ValueOf(&f64).Elem())
    if err != nil {
        T.Errorf("Error converting string to float64; %s", err.String())
    }
    if f64 != float64(1.25) {
        T.Errorf("Failed to covert string to float64")
    }
}
