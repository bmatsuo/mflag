// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mflag

/*  Filename:    flagset.go
 *  Author:      Bryan Matsuo <bryan.matsuo@gmail.com>
 *  Created:     Sun Oct 23 22:40:10 PDT 2011
 *  Description: 
 */

import (
    "template"
    "reflect"
    "strings"
    "bytes"
    "fmt"
    "os"
)

//  A flag represents a single command line flag. It correponds to a field
//  in a struct.
type Flag struct {
    name string
    def string
    help string
    field int
    typ reflect.Type
    setcount int
}

//  The name of the flag as it would be set in the command line.
func (f *Flag) Name() string { return f.name }

//  The help message for flag specified in the corresponding struct field's
//  help tag.
func (f *Flag) Help() string { return f.help }

//  The string default value of the flag.
func (f *Flag) Default() string {
    if f.typ.Kind() == reflect.String {
        return fmt.Sprintf(`"%s"`, f.def)
    }
    return f.def
}

//  The FlagSet object is responsible for parsing flags in the mflag package.
//  It holds a pointer to a struct and sets the struct's fields from the flags
//  set in command line arguments.
type FlagSet struct {
    usage    string
    flaghelp string
    name     string
    styles   []string
    handling ErrorHandling
    err      os.Error
    v        reflect.Value
    t        reflect.Type
    flags    []*Flag
    fcount   int
    args     []string
}

//  Create a new FlagSet using the object obj. The obj must be a pointer
//  to a struct. It's fields must be 'marshalable'. See FlagSet.Marhall
//  and FlagSet.UnMarshall for more information.
//
//  If there are errors creating a new FlagSet, a runtime panic occurs.
func New(v interface{}) *FlagSet {
    fs := new(FlagSet)
    val := reflect.ValueOf(v)

    // Ensure that the value is a non-nil pointer to a struct.
    if k := val.Kind(); k != reflect.Ptr {
        panic("non-ptr")
    }
    fs.v = val.Elem()
    if !fs.v.IsValid() {
        panic("nil ptr")
    }
    fs.t = fs.v.Type()
    if k := fs.t.Kind(); k != reflect.Struct {
        panic("non-struct ptr")
    }

    // Fill out the flags list.
    for i := 0; i < fs.t.NumField(); i++ {
        field := fs.t.Field(i)
        name := field.Tag.Get("flag")
        if name == "_skip" {
            continue
        }

        f := new(Flag)
        f.field = i
        if f.name = name; f.name == "" {
            f.name = strings.ToLower(field.Name)
        }
        f.help = field.Tag.Get("help")
        f.typ = field.Type
        f.def = field.Tag.Get("def")
        // Load the default value from the struct value if it is set non-zero.
        defaultFromStruct := f.def == "" &&
            !reflect.DeepEqual(fs.v.Field(i).Interface(), reflect.Zero(field.Type))
        if defaultFromStruct {
            f.def = valueToString(fs.v.Field(i))
        }
        fs.flags = append(fs.flags, f)
    }

    // Fill in the remaining default configuration if no errors have occurred.
    fs.Named(fs.t.Name())
    fs.name = fs.t.Name()
    fs.styles = []string{"[options] [Argument ...]"}
    fs.usage = "Usage: {{ with $x := . }}{{ range .ArgStyles }}{{ $x.Name }} {{ . }}\n{{ end }}{{ end }}"
    fs.flaghelp = `  -{{ .Name }}={{ .Default }}{{ if .Help }}: {{ .Help }}{{ end }}`
    return fs
}

//  Print fs' help message to standard error.
func (fs *FlagSet) UsageString() (usage string, err os.Error) {
    var t *template.Template
    t, err = template.New("UsageTemplate").Parse(fs.usage)
    if err != nil {
        return
    }
    b := new(bytes.Buffer)
    err = t.Execute(b, map[string]interface{}{"Name": fs.name, "ArgStyles": fs.styles})
    if err != nil {
        return
    }
    usage = string(b.Bytes())
    return
}

//  Print fs' help message to standard error.
func (fs *FlagSet) PrintHelp() (err os.Error) {
    var ftempl *template.Template
    // Attempt to compile the flag template before printing anything.
    ftempl, err = template.New("FlagTemplate").Parse(fs.flaghelp)
    if err != nil {
        return
    }

    // Print usage with a trailing newline.
    var usage string
    usage, err = fs.UsageString()
    if err != nil {
        return
    }
    if usage != "" {
        _, err = fmt.Fprint(os.Stderr, usage)
        if err != nil {
            return
        }
    }
    // Print the help message for each flag.
    for _, f := range fs.flags {
        err = ftempl.Execute(os.Stderr, f)
        if err != nil {
            return
        }
        _, err = fmt.Fprintln(os.Stderr)
    }
    return
}

//  The i-th flag in the FlagSet.
//func (fs *FlagSet)Flag(i int) *Flag { return fs.flags[i] }

func (fs *FlagSet) flagField(i int) int { return fs.flags[i].field }

func (fs *FlagSet) flagWithName(name string) int {
    for i := range fs.flags {
        if fs.flags[i].name == name {
            return i
        }
    }
    return -1
}

func (fs *FlagSet)setFlag(i int, value string) {
    f := fs.flags[i]
    stringToValue(value, fs.v.Field(f.field))
    f.setcount++
    fs.fcount++
}

//  Set the string value of a named flag.
func (fs *FlagSet) Set(flag string, value string) bool {
    i := fs.flagWithName(flag)
    if i < 0 {
        return false
    }
    fs.setFlag(i, value)
    return true
}

//  Set the name of the FlagSet for use in the PrintHelp method. The name
//  of the FlagSet defaults to the type name of the struct used to initialize
//  it. Returns the FlagSet fs so calls can be chained.
func (fs *FlagSet) Named(name string) *FlagSet {
    fs.name = name
    return fs
}

//  Describe the programs arguments with short string. This defaults to a
//  generic argument description string "[options] [ARGUMENT ...]".
//  Returns the FlagSet fs so calls can be chained.
func (fs *FlagSet) WithArgs(desc... string) *FlagSet {
    fs.styles = desc
    return fs
}

//  Supply a template with which to generate flag help text. Returns the
//  FlagSet fs so calls can be chained.
func (fs *FlagSet) WithFlags(helptempl string) *FlagSet {
    fs.flaghelp = helptempl
    return fs
}

//  Supply a template with which to generate the usage string. Returns the
//  FlagSet fs so calls can be chained.
func (fs *FlagSet) WithUsage(usagetempl string) *FlagSet {
    fs.usage = usagetempl
    return fs
}

//  Set the error handling strategy for FlagSet fs. Returns fs so the call can
//  be  chained.
func (fs *FlagSet) OnError(action ErrorHandling) *FlagSet {
    fs.handling = action
    return fs
}

func (fs *FlagSet) handleError(err os.Error) {
    fs.err = err
    switch fs.handling {
    case Continue:
        return
    case Exit:
        fs.PrintHelp()
        os.Exit(1)
    case Panic:
        fs.PrintHelp()
        panic(err)
    }
}

//  Parse a list of arguments, assigning values the fields of the object used
//  to create fs. It returns the FlagSet fs, any remaining non-flag arguments,
//  and any error that was encountered during processing.
func (fs *FlagSet) Parse(args []string) (*FlagSet, []string, os.Error) {
    // Clear previous parse data.
    fs.args = nil
    fs.err = nil
    fs.fcount = 0
    firstarg := 0
    for i := range fs.flags {
        if f := fs.flags[i]; f.def != "" {
            fs.setFlag(i, f.def)
        }
        fs.flags[i].setcount = 0
    }

    var needvalue bool
    var fi int
    for _, arg := range args {
        if !needvalue && !strings.HasPrefix(arg, "-") {
            break
        }
        firstarg++

        // Set the last flag seen if it is missing a value.
        if needvalue {
            fs.setFlag(fi, arg)
            continue
        }

        // Find the specified flag.
        flagval := strings.SplitN(arg[1:], "=", 2)
        fi = fs.flagWithName(flagval[0])
        if fi < 0 {
            if fname := flagval[0]; fname == "h" || fname == "help" {
                fs.handleError(ErrHelp)
            } else {
                fs.handleError(ErrNoFlag)
            }
        }
        // Try to set the flag without grabbing another argument.
        if len(flagval) > 1 {
            fs.setFlag(fi, flagval[1])
        } else {
            // Boolean options take on true when not specifically assigned.
            if fs.flags[fi].typ.Kind() == reflect.Bool {
                fs.setFlag(fi, "true")
            } else {
                fs.handleError(ErrHelp)
            }
        }
    }
    if needvalue {
        // The last string in args is a (non-bool) flag.
        fs.handleError(ErrMissingValue)
    }

    // Save the list of unparsed arguments.
    if firstarg < len(args) {
        fs.args = args[firstarg:]
    }
    return fs, fs.args, fs.err
}

//  Return the non-flag arguments.
func (fs *FlagSet) Args() []string { return fs.args }

//  Return the i-th argument, see FlagSet.NArg().
func (fs *FlagSet) Arg(i int) string { return fs.args[i] }

//  The number of (non-option) arguments remaining after parsing.
func (fs *FlagSet) NArg() int { return len(fs.args) }

//  The number of flags set since the last fs.Parse().
func (fs *FlagSet) NFlag() int { return len(fs.args) }

//  Call fn on all the flags in fs.
func (fs *FlagSet) VisitAll(fn func(*Flag)) {
    for i := range fs.flags {
        fn(fs.flags[i])
    }
}

// Call fn on all flags in fs which have been set.
func (fs *FlagSet) Visit(fn func(*Flag)) {
    fs.VisitAll(func(f *Flag) {
        if f.setcount > 0 {
            fn(f)
        }
    })
}
