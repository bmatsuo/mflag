// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mflag

/*  Filename:    errors.go
 *  Author:      Bryan Matsuo <bryan.matsuo@gmail.com>
 *  Created:     Mon Oct 24 03:31:54 PDT 2011
 *  Description: FlagSet error handling.
 */

import (
    "os"
)

//  Possible errors that can be returned from FlagSet.Parse when its
//  ErrorHandling is Continue.
var (
    ErrHelp = os.NewError("mflag: help requested")
    ErrNoFlag = os.NewError("mflag: unknown flag specified")
    ErrMissingValue = os.NewError("mflag: flag missing value")
)

//  A type that specifies how the FlagSet object handles errors during the
//  Parse method.
type ErrorHandling uint

//  2 bits for control flow.
const (
    Continue ErrorHandling = iota
    Exit
    Panic
)

//  Mutually exclusive flags
const (
    Help ErrorHandling = 1 << 2 + iota
    Debug
)

//  Returns one of the control flow values; Continue, Exit, or Panic.
func (h ErrorHandling) Control() ErrorHandling { return h&3 }
//  Returns true if the Help bit is set.
func (h ErrorHandling) Help() bool { return h&Help != 0 }
//  Returns true if the Debug bit is set.
func (h ErrorHandling) Debug() bool { return h&Debug != 0 }
