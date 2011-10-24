// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    mflagSimple.go
 *  Author:      Bryan Matsuo <bryan.matsuo@gmail.com>
 *  Created:     Sun Oct 23 22:31:04 PDT 2011
 *  Description: Main source file in simple
 */

import (
    "mflag"
    //"log"
    //"fmt"
    "os"
)

//  Tip: Give the option structure the same name as the program.
type mFlagSimple struct {
    Output    string `help:"Program output path"`
    Frequency int32  `help:"Repeatition delay (ms)"`
    Verbose   bool
}

var Opt = mFlagSimple{Output:"./output.log", Frequency:500}

func main() {
    m := mflag.New(&Opt)
    m.Parse(os.Args[1:])
    m.PrintHelp()
}
