// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*  Filename:    mflag.go
 *  Author:      Bryan Matsuo <bryan.matsuo@gmail.com>
 *  Created:     Sun Oct 23 22:20:52 PDT 2011
 *  Description: Main source file in mflags
 */


/*
Package mflags (un)marshals command line options with structs. The package
is meant to intuitively feel like the "flag" package. But, it takes much
of the hassle out of the process. It is inspired by DRY because the struct
that wraps the option values defines the options themselves. Let's look at
an example; A simple clone of Unix's `head` called `gohead`.

    package main
    import (
        "mflag"
        "container/list"
        "bufio"
        "io"
        "os"
    )
    type gohead struct {
        Pre    int  `flag:"n" help:"number of lines"`
        Invert bool `flag:"i" help:"same as supplying negative numbers"`
    }
    var opt = gohead{Pre: 25}

    func main() {
        fs := mflag.New(&opt).
            WithArgs("[-n lines] [-i] [file ...]").
            OnError(mflag.Exit)
        fs.Parse(os.Args[1:])

        if opt.Pre < 0 {
            opt.Invert = !opt.Invert
            opt.Pre *= -1
        }

        headfn := head
        if opt.Invert {
            headfn = inversehead
        }
        for i := range fs.Args() {
            f, _ := os.Open(fs.Arg(i))
            headfn(f, opt.Pre)
        }
        if fs.NArg() == 0 {
            headfn(os.Stdout, opt.Pre)
        }
    }
    ...

Flag type and help specifications are derived from the source object's
struct definition. Default flag values are derived from the source object.
There are chainable helper methods to customize the FlagSet's usage like the
WithArgs method seen above. The rest of the program can be located in the
subdirectory "./examples/gohead/".

The above code shows that the programmatic interface is similar to the "flag"
package. The mflag.FlagSet object has the same basic workflow as the
flag.FlagSet object mod the type-specific methods like String, BoolVar,
Int64Var, etc. It is a good practice in general to wrap options up into a
struct (for static type safety) or map. So defining all the options separately
is really extra work.
 */
package mflag

import ()
