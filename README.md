About mflags
=============

Package mflags (un)marshals command line options with structs. The package
is meant to intuitively feel like the "flag" package. But, it takes much
of the hassle out of the process. It is inspired by DRY because the struct
that wraps the option values defines the options themselves.

Documentation
=============

Synopsis
--------

The package "mflag" consolidates flag parsing and reduces the amount of
repetition in command line flag declaration.

    type goserver struct {
        Port     int    `         help:"local http port"`
        Timeout  int64  `         help:"request timeout (ms)"`
        ProcName string `flag:"p" help:"process name"`
        BigCache bool   `flag:"C" help:"use a big cache for data"`
        Demonize bool   `flag:"d" help:"run as a deamon"`
        Seed     int64  `         help:"Random seed"`
    }
    var opt = goserver{
        Port: 6262, Timeout=300, ProcName:"goserver",
        BigCache:true, Seed:0xBEEF,
    }

    func main() {
        fs = mflag.New(&opt).
            WithArgs("[-Cd] [-p pname] [-timeout ms] [-port port] source [...]").
            OnError(mflag.Exit)
        fs.Parse(os.Args[1:])
        sources := fs.Args()
        ...

The help message from the above code would look something like

    Usage: goserver [-Cd] [-p pname] [-timeout ms] [-port port] source [...]
      -port=6262: local http port
      -timeout=300: request timeout (ms)
      -seed=3735928559: random seed
      -procname="goserver": process name
      -bigcache=false: use a big cache for data
      -demonize=true: run as a deamon

The logical behavior of the parser the package "mflag" is similar to the Go
standard "flag". The above example can be even more consise. Wrapping the
default values into struct definition.

    type goserver struct {
        Port     int    `         def:"6262"     help:"local http port"`
        Timeout  int64  `         def:"300"      help:"request timeout (ms)"`
        ProcName string `flag:"p" def:"goserver" help:"process name"`
        BigCache bool   `flag:"C" def:"true"     help:"use a big cache for data"`
        Demonize bool   `flag:"d"                help:"run as a deamon"`
        Seed     int64  `         def:"1337"     help:"Random seed"`
    }
    var opt = goserver{}

    func main() {
        fs := mflag.New(&opt).
            WithArgs("[-Cd] [-p pname] [-timeout ms] [-port port] source [...]").
            OnError(mflag.Exit)
        fs.Parse(os.Args[1:])
        sources := fs.Args()
        ...

This will do the same thing with a subtle difference. The default values in
this case are evaluated as strings and converted to their values when Parse
is called.

Installation
-------------

Use goinstall to install mflag

    goinstall github.com/bmatsuo/mflag

Otherwise clone the repository and install using `gomake`

    git clone https://github.com/bmatsuo/mflag

Examples
--------

There are working examples in the ./examples/ subdirectory. They can be
installed from the project directory with the command

    gomake ex.install

The example programs all have fairly silly executable names, so there is
little chance that they will overwrite another program. They can then be
uninstalled with the command

    gomake ex.nuke

General Documentation
---------------------

Use godoc to vew the documentation for mflags

    godoc github.com/bmatsuo/mflag

Or alternatively, use a godoc http server

    godoc -http=:6060

and view the url http://localhost:6060/pkg/github.com/bmatsuo/mflag/

Author
======

Bryan Matsuo <bryan.matsuo@gmail.com>

Copyright & License
===================

Copyright (c) 2011, Bryan Matsuo.
All rights reserved.

Use of this source code is governed by a BSD-style license that can be
found in the LICENSE file.
