package mflag

/*  Filename:    flagset_test.go
 *  Author:      Bryan Matsuo <bryan.matsuo@gmail.com>
 *  Created:     Sun Oct 23 22:40:10 PDT 2011
 *  Description: For testing flagset.go
 */

import (
    "testing"
)

func FlagWithNameTest(fs *FlagSet, flag string, T *testing.T) (i int) {
    i = fs.flagWithName(flag)
    if i < 0 {
        T.Errorf("Can't find flag %s", flag)
    }
    return
}

type SampleFlags struct {
    Name     string   `help:"User name"`
    NumTimes int      `def:"10" flag:"n" help:"# of runs"`
    Paths    []string `flag:"_skip"`
    Verbose  bool     `help:"Verbose output"`
}

func TestFlagDefaults(T *testing.T) {
    testobj := SampleFlags{Name:"user"}
    _, _, err := New(&testobj).Parse([]string{})
    if err != nil {
        T.Fatalf("Error parsing flags; %s", err.String())
    }
    if testobj.NumTimes != 10 {
        T.Errorf("Tag specified default %d not set properly %d", 10, testobj.NumTimes)
    }
}

func TestFlagNames(T *testing.T) {
    testobj := SampleFlags{}
    testflags := New(&testobj)
    okFlagName := func(field, flag string) {
        i := FlagWithNameTest(testflags, flag, T)
        if i < 0 {
            return
        }
        if testflags.t.Field(testflags.flagField(i)).Name != field {
            T.Errorf("Wrong field %s for flag %s", field, flag)
        }
    }
    okFlagName("Name", "name")
    okFlagName("NumTimes", "n")
    okFlagName("Verbose", "verbose")
}

func TestFlagHelpTag(T *testing.T) {
    testobj := SampleFlags{}
    testflags := New(&testobj)
    okHelpTag := func(flag, help string) {
        i := FlagWithNameTest(testflags, flag, T)
        if i < 0 {
            return
        }
        if h := testflags.flags[i].help; h != help {
            T.Errorf("Wrong help message '%s' for flag %s (not %s)",
                h, flag, help)
        }
    }
    okHelpTag("name", "User name")
    okHelpTag("n", "# of runs")
    okHelpTag("verbose", "Verbose output")
}

func TestFlagDefaultsStrings(T *testing.T) {
    testobj := SampleFlags{Name:"user", NumTimes:10}
    testflags := New(&testobj)
    // Parse a command with no options.
    okDefault := func(flag, def string) {
        i := FlagWithNameTest(testflags, flag, T)
        if i < 0 {
            return
        }
        if d := testflags.flags[i].def; d != def {
            T.Errorf("Wrong default %s for flag %s (not %s)", d, flag, def)
        }
    }
    okDefault("name", "user")
    okDefault("n", "10")
    okDefault("verbose", "false")
}

func TestUsage(T *testing.T) {
    testobj := SampleFlags{Name:"user", NumTimes:10}
    testflags := New(&testobj)
    if testflags.name != "SampleFlags" {
        T.Errorf("FlagSet name is not initialized correctly %s", testflags.name)
    }
    usage, err := testflags.UsageString()
    if err != nil {
        T.Errorf("Error generating usage string; %s", err.String())
    }
    if usage != "Usage: SampleFlags [options] [Argument ...]\n" {
        T.Errorf("FlagSet usage string is not correct '%s'", usage)
    }
    usage, err = testflags.WithArgs("a b", "c d").UsageString()
    if err != nil {
        T.Errorf("Error generating usage string; %s", err.String())
    }
    if usage != "Usage: SampleFlags a b\n       SampleFlags c d\n" {
        T.Errorf("FlagSet usage string is not correct '%s'", usage)
    }
}

func TestFlagSkip(T *testing.T) {
    testobj := SampleFlags{Name:"user", NumTimes:10}
    _, _, err := New(&testobj).Parse([]string{})
    if err != nil {
        T.Fatalf("Error parsing flags; %s", err.String())
    }
}
