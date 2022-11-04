// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

// Package clip “gee arg” provides yet another Go command line argument
// parser.
//
// # Overview
//
// clip can handle flags, single argument options, multiple argument
// options, and positional arguments.
//
// # Flags
//
// A flag is either present or absent.
//
// Examples:
//
//	myapp -v
//	myapp --verbose
//
// If the flag is present, the option's value is true; otherwise it is
// false.
//
// Flags support short and long names. For example, a flag name of "version"
// can be set with `--version` or `-v`. If you don't want a short name, or
// want a different one (e.g., `-V`), use [Option.SetShortName].
//
//	parser := NewParserVersion("1.0.0") # AppName is os.Base(os.Args[0]))
//	verboseOpt := parser.Flag("verbose", "whether to show more output")
//	parser.ParseLine("")
//	verbose := verboseOpt.Value() // verbose == false
//	// -or-
//	verbose = verboseOpt.Given() // verbose == false
//
//	parser.ParseLine("-v")
//	verbose = verboseOpt.Value() // verbose == true
//	// -or-
//	verbose = verboseOpt.Given() // verbose == true
//
// If you want the user to be able to optionally specify how verbose to be
// then use an Int value option: see [Parser.Int].
//
// Multiple flags can be grouped together if their short names are used,
// e.g., given flags `-v`, `-x`, and `-c`, they can be set individually, or
// together, i.e., `-v -x -c` or `-vxc`. The last option in such a group may
// be a single- or multi-value option. For example, if option `o` takes a
// string argument, we could write any of these:
//
//	myapp -v -x -c -o outfile.dat
//	myapp -v -x -c -o=outfile.dat
//	myapp -vxcooutfile.dat
//	myapp -vxco outfile.dat
//	myapp -vxco=outfile.dat
//
// And if we are happy with `-o`'s default value, we can use these:
//
//	myapp -v -x -c -o
//	myapp -v -x -c
//	myapp -vxco
//	myapp -vxc
//
// All of which set the `v`, `x`, and `c` flags as before and set the `-o`
// option to its default value.
//
// # Single Value Options
//
// A single value option is either present—either with a value or without
// (in which case its default is used)—or absent, in which case its default
// is its value.
//
// Examples:
//
//	myapp
//	myapp -v
//	myapp --verbose
//	myapp -v1
//	myapp -v=2
//	myapp -v 3
//	myapp --verbose=4
//	myapp --verbose 5
//
// If the option is absent, the option's value is the default that was set.
// If the option is present, the option's value is the default if no value
// is given, otherwise the given value.
//
// If you need to distinguish between whether a value was given at all
// (i.e., between the first two examples, assuming the default was set to
// 1), then use [Option.Given].
//
//	parser := NewParser()
//	verboseOpt := parser.Int("verbose", "how much output to show", 1)
//	verboseOpt.AllowImplicit = true // implicitly use the default so -v → -v1
//	parser.ParseLine("")
//	verbose := 0 // assume no verbosity
//	if verboseOpt.Given() {
//		verbose = verboseOpt.Value()
//	}
//
// Here, verbose == 0 (since we started at 0 and checked whether it was
// given and it wasn't)
//
//	// first three lines as before
//	parser.ParseLine("-v")
//	verbose := 0 // assume no verbosity
//	if verboseOpt.Given() {
//		verbose = verboseOpt.Value()
//	}
//
// Here, verbose == 1 (since it was given with no value, but due to
// AllowImplicit, the default was used for its value)
//
//	// first three lines as before
//	parser.ParseLine("-v2")
//	verbose := 0 // assume no verbosity
//	if verboseOpt.Given() {
//		verbose = verboseOpt.Value()
//	}
//
// Here, verbose == 2 (as given)
//
// TODO IntInRange eg + test
// TODO Real eg + test + note RealInRange
// TODO Choice eg + test
// TODO Str eg + test
//
// # Multi-Value Options TODO text + tests
//
// TODO Strs eg + test
// TODO Ints eg + test
// TODO Reals eg + test
//
// # Setting a Validator # TODO
//
// # Post-Parsing Validation TODO test
//
// If some post-parsing validation finds invalid data it is possible to
// treat it as a parser error by calling [Parser.OnError] with a message
// string.
//
// # Required Options TODO tests
//
// This is a contradiction in terms, but if we really want to require an
// option then handle it like this:
//
//	parser := NewParser() // below: name, help, minimum, maximum, default
//	countOpt := parser.IntInRange("count", "how many are wanted", 0, 100, 0)
//	parser.ParseLine("")
//	if !countOpt.Given() {
//		parser.OnMissing(countOpt) // won't return (calls os.Exit)
//	}
//	count := countOpt.Value() // if we got here the user set it
package clip
