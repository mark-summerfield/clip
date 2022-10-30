// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

// Package garg “gee arg” provides yet another Go command line argument
// parser.
//
// # Overview
//
// garg can handle flags, single argument options, multiple argument
// options, subcommands, and positional options.
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
//	parser := NewParser()
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
//		parser := NewParser()
//		verboseOpt := parser.Int("verbose", "how much output to show", 1)
//	 verboseOpt.AllowImplicit() // implicitly use the default so -v → -v1
//		parser.ParseLine("")
//		verbose := 0 // assume no verbosity
//		if verboseOpt.Given() {
//			verbose = verboseOpt.Value()
//		}
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
package garg

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

type Parser struct {
	DontExit          bool
	Positionals       []string
	HelpName          string
	VersionName       string
	appName           string
	appVersion        string
	subCommands       map[string]*SubCommand
	positionalCount   ValueCount
	positionalVarName string
	use_h_for_help    bool
	use_v_for_version bool
	use_V_for_version bool
}

func NewParser() Parser {
	appName := "<app>"
	if len(os.Args) > 0 {
		appName = path.Base(os.Args[0])
	}
	subcommands := make(map[string]*SubCommand)
	subcommands[mainSubCommand] = newMainSubCommand()
	return Parser{appName: appName, appVersion: "",
		subCommands: subcommands, positionalCount: ZeroOrMore,
		positionalVarName: "FILENAME", HelpName: "help",
		use_h_for_help: true,
	}
}

func (me *Parser) AppName() string {
	return me.appName
}

func (me *Parser) SetAppName(name string) {
	if name == "" {
		panic("#100: can't have empty appname")
	}
	me.appName = name
}

func (me *Parser) Version() string {
	return me.appVersion
}

func (me *Parser) SetVersion(version string) {
	me.appVersion = version
	if me.VersionName == "" {
		me.VersionName = "version"
		me.use_v_for_version = true
	}
}

func (me *Parser) SetPositionalCount(valueCount ValueCount) {
	me.positionalCount = valueCount
}

func (me *Parser) SetPositionalVarName(name string) {
	me.positionalVarName = name
}

func (me *Parser) SubCommand(name, help string) *SubCommand {
	if name == "" {
		panic("#110: can't have empty subcommand name")
	}
	subcommand := newSubCommand(name, help)
	me.subCommands[name] = subcommand
	return subcommand
}

func (me *Parser) Flag(name, help string) *FlagOption {
	option := newFlagOption(name, help)
	me.registerNewOption(option)
	return option
}

func (me *Parser) Int(name, help string, defaultValue int) *IntOption {
	option := newIntOption(name, help, defaultValue)
	me.registerNewOption(option)
	return option
}

func (me *Parser) IntInRange(name, help string,
	minimum, maximum, defaultValue int) *IntOption {
	option := newIntOption(name, help, defaultValue)
	option.validator = makeIntRangeValidator(minimum, maximum)
	me.registerNewOption(option)
	return option
}

func (me *Parser) Real(name, help string, defaultValue float64) *RealOption {
	option := newRealOption(name, help, defaultValue)
	me.registerNewOption(option)
	return option
}

func (me *Parser) RealInRange(name, help string,
	minimum, maximum, defaultValue float64) *RealOption {
	option := newRealOption(name, help, defaultValue)
	option.validator = makeRealRangeValidator(minimum, maximum)
	me.registerNewOption(option)
	return option
}

func (me *Parser) Str(name, help, defaultValue string) *StrOption {
	option := newStrOption(name, help, defaultValue)
	me.registerNewOption(option)
	return option
}

func (me *Parser) Choice(name, help string, choices []string,
	defaultValue string) *StrOption {
	option := newStrOption(name, help, defaultValue)
	option.validator = makeChoiceValidator(choices)
	me.registerNewOption(option)
	return option
}

func (me *Parser) Strs(name, help string) *StrsOption {
	option := newStrsOption(name, help)
	me.registerNewOption(option)
	return option
}

func (me *Parser) registerNewOption(option Optioner) {
	me.subCommands[mainSubCommand].options = append(
		me.subCommands[mainSubCommand].options, option)
}

func (me *Parser) Parse() error {
	return me.ParseArgs(os.Args[1:])
}

func (me *Parser) ParseLine(line string) error {
	return me.ParseArgs(strings.Fields(line))
}

func (me *Parser) ParseArgs(args []string) error {
	var err error
	me.prepareHelpAndVersionOptions()
	subcommand, tokens, err := me.tokenize(args)
	if err != nil {
		return err
	}
	//fmt.Printf("TOKENS: %q %s\n", strings.Join(args, " "), tokens)
	var currentOption Optioner
	inPositionals := false
	expect := Zero // ValueCount - how many values we expect to follow opt
	for _, token := range tokens {
		if token.positionalsFollow {
			inPositionals = true
		} else if inPositionals {
			me.addPositional(token.text)
		} else if !token.isValue() { // Option
			if currentOption != nil {
				if err = me.checkValue(currentOption); err != nil {
					return err
				}
			}
			currentOption = token.option
			expect = currentOption.ValueCount()
			if option, ok := currentOption.(*FlagOption); ok {
				option.value = true
			}
		} else { // Value
			switch expect {
			case Zero:
				inPositionals = me.addPositional(token.text)
			case ZeroOrOne:
				if currentOption.Count() == 1 {
					inPositionals = me.addPositional(token.text)
				} else {
					err = me.addValue(currentOption, token.text)
				}
			case ZeroOrMore:
				err = me.addValue(currentOption, token.text)
			case One:
				if currentOption.Count() == 0 {
					err = me.addValue(currentOption, token.text)
				} else {
					inPositionals = me.addPositional(token.text)
				}
			case OneOrMore:
				err = me.addValue(currentOption, token.text)
			case Two, Three:
				panic("#200: invalid ValueCount")
			}
		}
		if err != nil {
			return err
		}
	}
	if err := me.checkPositionals(); err != nil {
		return err
	}
	return me.checkValues(subcommand.options)
}

func (me *Parser) addValue(option Optioner, value string) error {
	var err error
	switch option := option.(type) {
	case *IntOption:
		err = option.addValue(value)
	case *RealOption:
		err = option.addValue(value)
	case *StrOption:
		err = option.addValue(value)
	case *StrsOption:
		err = option.addValue(value)
	default:
		panic("#300: missed type case")
	}
	if err != nil {
		return me.handleError(8, fmt.Sprintf("invalid value for %s: %s",
			option.LongName(), err))
	}
	return nil
}

func (me *Parser) addPositional(value string) bool {
	if me.Positionals == nil {
		me.Positionals = make([]string, 0, 1)
	}
	me.Positionals = append(me.Positionals, value)
	return true
}

func (me *Parser) prepareHelpAndVersionOptions() {
	seen_V := false
	main := me.subCommands[mainSubCommand]
	for _, option := range main.options {
		if option.LongName() == me.HelpName {
			panic("#400: only auto-generated help is supported")
		} else if option.LongName() == me.VersionName {
			panic("#402: only auto-generated version is supported")
		}
		if me.use_h_for_help && option.ShortName() == 'h' {
			me.use_h_for_help = false
		}
		if option.ShortName() == 'v' {
			me.use_v_for_version = false
		} else if option.ShortName() == 'V' {
			seen_V = true
		}
	}
	if me.VersionName != "" && !me.use_v_for_version && !seen_V {
		me.use_V_for_version = true
	}
}

func (me *Parser) tokenize(args []string) (*SubCommand, []token, error) {
	var err error
	state := me.initializeTokenState()
	helpName, versionName := me.getHelpAndVerboseNames()
	tokens := make([]token, 0, len(args))
	for i, arg := range args {
		if me.handledHelpOrVerbose(arg, helpName, versionName) {
			return nil, nil, nil
		}
		if arg == "--" { // --
			tokens = append(tokens, newPositionalsFollowToken())
			for _, v := range args[i+1:] {
				tokens = append(tokens, newValueToken(v))
			}
			break
		}
		if strings.HasPrefix(arg, "--") { // --option --option=value
			tokens, err = me.handleLongOption(arg, tokens, &state)
			if err != nil {
				return state.subcommand, tokens, err
			}
		} else if strings.HasPrefix(arg, "-") {
			if _, err := strconv.ParseFloat(arg, 64); err == nil {
				tokens = append(tokens, newValueToken(arg)) // -int | -real
			} else {
				tokens, err = me.handleShortOption(arg, tokens, &state)
				if err != nil {
					return state.subcommand, tokens, err
				}
			}
		} else if state.hasSubCommands && !state.hadSubCommand { // subcmd?
			tokens = me.handlePossibleSubcommand(arg, tokens, &state)
		} else {
			tokens = append(tokens, newValueToken(arg))
		}
	}
	return state.subcommand, tokens, nil
}

func (me *Parser) initializeTokenState() tokenState {
	state := tokenState{
		subcommand:        me.subCommands[mainSubCommand],
		subCommandForName: me.getSubCommandsForNames(),
		hasSubCommands:    len(me.subCommands) > 1,
		hadSubCommand:     false,
	}
	state.optionForLongName, state.optionForShortName =
		state.subcommand.optionsForNames()
	return state
}

func (me *Parser) getHelpAndVerboseNames() (string, string) {
	helpName := fmt.Sprintf("--%s", me.HelpName)
	versionName := ""
	if me.VersionName != "" {
		versionName = fmt.Sprintf("--%s", me.VersionName)
	}
	return helpName, versionName
}

func (me *Parser) handledHelpOrVerbose(arg, helpName,
	versionName string) bool {
	if arg == helpName || (me.use_h_for_help && arg == "-h") {
		me.onHelp() // doesn't return
		return true
	}
	if arg == versionName || (me.VersionName != "" &&
		(me.use_v_for_version && arg == "-v") ||
		(me.use_V_for_version && arg == "-V")) {
		me.onVersion() // doesn't return
		return true
	}
	return false
}

func (me *Parser) handleLongOption(arg string, tokens []token,
	state *tokenState) ([]token, error) {
	name := strings.TrimPrefix(arg, "--")
	parts := strings.SplitN(name, "=", 2)
	if len(parts) == 2 { // --option=value
		name := parts[0]
		option, ok := state.optionForLongName[name]
		if ok {
			tokens = append(tokens, newNameToken(name, option))
			tokens = append(tokens, newValueToken(parts[1]))
		} else {
			return tokens, me.handleError(10, fmt.Sprintf(
				"unrecognized option --%s", name))
		}
	} else { // --option
		option, ok := state.optionForLongName[name]
		if ok {
			tokens = append(tokens, newNameToken(name, option))
		} else {
			return tokens, me.handleError(12, fmt.Sprintf(
				"unrecognized option --%s", name))
		}
	}
	return tokens, nil
}

func (me *Parser) handleShortOption(arg string, tokens []token,
	state *tokenState) ([]token, error) {
	// -a -ab -abcValue -c=value -abc=value
	text := strings.TrimPrefix(arg, "-")
	parts := strings.SplitN(text, "=", 2)
	var pendingValue string
	if len(parts) == 2 { // -a=value -abc=value
		text = parts[0]
		pendingValue = parts[1]
	}
	for i, c := range text {
		name := string(c)
		option, ok := state.optionForShortName[name]
		if ok {
			tokens = append(tokens, newNameToken(name, option))
			if _, isFlag := option.(*FlagOption); !isFlag &&
				i+1 < len(text) {
				value := text[i+1:] // -aValue -abcValue
				tokens = append(tokens, newValueToken(value))
			}
		} else if pendingValue == "" {
			last := len(tokens) - 1
			rest := text[i:]
			if last >= 0 && rest != tokens[last].text {
				return tokens, me.handleError(16, fmt.Sprintf(
					"unexpected value %s", rest))
			}
			break
		} else {
			return tokens, me.handleError(18, fmt.Sprintf(
				"unrecognized option -%s", name))
		}
	}
	if pendingValue != "" {
		tokens = append(tokens, newValueToken(pendingValue))
	}
	return tokens, nil
}

// is it a subcommand? - only allow one subcommand (excl. main)
func (me *Parser) handlePossibleSubcommand(arg string, tokens []token,
	state *tokenState) []token {
	state.hadSubCommand = true
	cmd, ok := state.subCommandForName[arg]
	if ok {
		state.subcommand = cmd
		state.optionForLongName, state.optionForShortName =
			state.subcommand.optionsForNames()
	} else { // value
		tokens = append(tokens, newValueToken(arg))
	}
	return tokens
}

func (me *Parser) getSubCommandsForNames() map[string]*SubCommand {
	cmdForName := make(map[string]*SubCommand, len(me.subCommands)*2)
	for long, command := range me.subCommands {
		if long != mainSubCommand {
			cmdForName[long] = command
			if command.ShortName() != 0 {
				cmdForName[string(command.ShortName())] = command
			}
		}
	}
	return cmdForName
}

func (me *Parser) onHelp() {
	fmt.Printf("usage: %s TODO", me.appName) // TODO
	if !me.DontExit {
		os.Exit(0)
	}
}

func (me *Parser) onVersion() {
	fmt.Printf("%s v%s", me.appName, me.appVersion)
	if !me.DontExit {
		os.Exit(0)
	}
}

func (me *Parser) checkPositionals() error {
	count := len(me.Positionals)
	switch me.positionalCount {
	case Zero:
		if count > 0 {
			return me.handleError(20,
				fmt.Sprintf("expected no positional arguments, got %d",
					count))
		}
	case ZeroOrOne:
		if count > 1 {
			return me.handleError(21,
				fmt.Sprintf(
					"expected zero or one positional arguments, got %d",
					count))
		}
	case ZeroOrMore: // any count is valid
	case One:
		if count != 1 {
			return me.handleError(22,
				fmt.Sprintf(
					"expected exactly one positional argument, got %d",
					count))
		}
	case OneOrMore:
		if count == 0 {
			return me.handleError(23,
				fmt.Sprintf(
					"expected at least one positional argument, got %d",
					count))
		}
	case Two:
		if count != 2 {
			return me.handleError(24,
				fmt.Sprintf(
					"expected exactly two positional arguments, got %d",
					count))
		}
	case Three:
		if count != 3 {
			return me.handleError(25,
				fmt.Sprintf(
					"expected exactly three positional arguments, got %d",
					count))
		}
	}
	return nil
}

func (me *Parser) checkValues(options []Optioner) error {
	for _, option := range options {
		if err := me.checkValue(option); err != nil {
			return err
		}
	}
	return nil
}

func (me *Parser) checkValue(option Optioner) error {
	count := option.Count()
	switch option.ValueCount() {
	case Zero:
		if _, isFlag := option.(*FlagOption); !isFlag {
			panic(fmt.Sprintf(
				"#600: nonflag option %s with zero ValueCount",
				option.LongName()))
		}
	case ZeroOrOne:
		if count > 1 {
			return me.handleError(32,
				fmt.Sprintf(
					"expected zero or one values for %s, got %d",
					option.LongName(), count))
		}
	case ZeroOrMore: // any count is valid
	case One:
		if count != 1 {
			return me.handleError(34,
				fmt.Sprintf("expected exactly one value for %s, got %d",
					option.LongName(), count))
		}
	case OneOrMore:
		if count == 0 {
			return me.handleError(36,
				fmt.Sprintf(
					"expected at least one value for %s, got %d",
					option.LongName(), count))
		}
	case Two, Three:
		panic("#602: invalid ValueCount")
	}
	return nil
}

func (me *Parser) handleError(code int, msg string) error {
	msg = fmt.Sprintf("error #%d: %s", code, msg)
	if !me.DontExit {
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(2)
	}
	return errors.New(msg)
}

func (me *Parser) OnMissing(option Optioner) error {
	if option.ShortName() != 0 {
		return me.handleError(0,
			fmt.Sprintf("option -%c (or --%s) is required",
				option.ShortName(), option.LongName()))
	}
	return me.handleError(0, fmt.Sprintf("option --%s is required",
		option.LongName()))
}

func (me *Parser) OnError(msg string) error {
	return me.handleError(1, msg)
}
