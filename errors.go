// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mflag

/*  Filename:    errors.go
 *  Author:      Bryan Matsuo <bryan.matsuo@gmail.com>
 *  Created:     Mon Oct 24 03:31:54 PDT 2011
 *  Description: 
 */

import (
    "os"
)

type ErrorHandling uint

const (
    Continue ErrorHandling = iota
    Exit
    Panic
)

var ErrHelp = os.NewError("mflag: help requested")
var ErrNoFlag = os.NewError("mflag: unknown flag specified")
var ErrMissingValue = os.NewError("mflag: flag missing value")
