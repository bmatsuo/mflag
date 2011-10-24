// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    mflagChaining.go
 *  Author:      Bryan Matsuo <bryan.matsuo@gmail.com>
 *  Created:     Mon Oct 24 04:44:45 PDT 2011
 *  Description: Main source file in chaining
 */

import (
    "mflag"
    //"log"
    //"fmt"
    "os"
)

type options struct {
    Verbose bool `flag:"v" help:"verbose output"`
}

func main() {
    opt := options{}
    mflag.New(&opt).
        Named("mflagChaining").
        WithArgs("command [argument ...]").
        OnError(mflag.Panic).
        Parse(os.Args[1:])
}
