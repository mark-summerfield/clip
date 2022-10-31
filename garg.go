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
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

type Parser struct {
	Positionals       []string
	HelpName          string
	VersionName       string
	shortVersionName  rune
	appName           string
	appVersion        string
	subCommands       map[string]*SubCommand
	mainSubCommand    *SubCommand
	positionalCount   PositionalCount
	positionalVarName string
	useLowerhForHelp  bool
}

// NewParser creates a new command line parser.
// If appname == "" the executable's basename will be used.
// If version == "" no version option will be available.
func NewParser(appname, version string) Parser {
	if appname == "" {
		appname = "<app>"
		if len(os.Args) > 0 {
			appname = path.Base(os.Args[0])
		}
	}
	mainSubCommand := newMainSubCommand()
	subcommands := make(map[string]*SubCommand)
	subcommands[mainSubCommandName] = mainSubCommand
	return Parser{appName: appname, appVersion: version,
		subCommands: subcommands, positionalCount: ZeroOrMorePositionals,
		positionalVarName: "FILENAME", HelpName: "help",
		VersionName: "version", useLowerhForHelp: true,
		mainSubCommand: mainSubCommand}
}

func (me *Parser) AppName() string {
	return me.appName
}

func (me *Parser) Version() string {
	return me.appVersion
}

func (me *Parser) SetPositionalCount(valueCount PositionalCount) {
	me.positionalCount = valueCount
}

func (me *Parser) SetPositionalVarName(name string) {
	me.positionalVarName = name
}

func (me *Parser) SubCommand(name, help string) *SubCommand {
	if name == "" {
		panic(fmt.Sprintf("#%d: can't have empty subcommand name",
			pEmptySubCommandName))
	}
	subcommand := newSubCommand(name, help)
	me.subCommands[name] = subcommand
	return subcommand
}

func (me *Parser) Flag(name, help string) *FlagOption {
	return me.mainSubCommand.Flag(name, help)
}

func (me *Parser) Int(name, help string, defaultValue int) *IntOption {
	return me.mainSubCommand.Int(name, help, defaultValue)
}

func (me *Parser) IntInRange(name, help string,
	minimum, maximum, defaultValue int) *IntOption {
	return me.mainSubCommand.IntInRange(name, help, minimum, maximum,
		defaultValue)
}

func (me *Parser) Real(name, help string, defaultValue float64) *RealOption {
	return me.mainSubCommand.Real(name, help, defaultValue)
}

func (me *Parser) RealInRange(name, help string,
	minimum, maximum, defaultValue float64) *RealOption {
	return me.mainSubCommand.RealInRange(name, help, minimum, maximum,
		defaultValue)
}

func (me *Parser) Str(name, help, defaultValue string) *StrOption {
	return me.mainSubCommand.Str(name, help, defaultValue)
}

func (me *Parser) Choice(name, help string, choices []string,
	defaultValue string) *StrOption {
	return me.mainSubCommand.Choice(name, help, choices, defaultValue)
}

func (me *Parser) Strs(name, help string) *StrsOption {
	return me.mainSubCommand.Strs(name, help)
}

func (me *Parser) Parse() error {
	return me.ParseArgs(os.Args[1:])
}

func (me *Parser) ParseLine(line string) error {
	return me.ParseArgs(strings.Fields(line))
}

func (me *Parser) ParseArgs(args []string) error {
	if err := me.prepareHelpAndVersionOptions(); err != nil {
		return err
	}
	subcommand, tokens, err := me.tokenize(args)
	if err != nil {
		return err
	}
	var currentOption optioner
	inPositionals := false
	for _, token := range tokens {
		if token.kind == positionalsFollowTokenKind {
			inPositionals = true
		} else if inPositionals {
			me.addPositional(token.text)
		} else if token.kind == helpTokenKind {
			me.onHelp(subcommand) // doesn't return
		} else if token.kind == nameTokenKind { // Option
			currentOption = token.option
			if me.isVersion(subcommand, currentOption) { // may not return
				return nil
			}
			if option, ok := currentOption.(*FlagOption); ok {
				option.value = true
			}
		} else { // Value
			if currentOption.wantsValue() {
				if msg := currentOption.addValue(token.text); msg != "" {
					return me.handleError(eInvalidValue, msg)
				}
			} else {
				inPositionals = me.addPositional(token.text)
			}
		}
	}
	if err := me.checkPositionals(); err != nil {
		return err
	}
	return me.checkValues(subcommand.options)
}

func (me *Parser) prepareHelpAndVersionOptions() error {
	usevForVersion := true
	useVForVersion := false
	seenV := false
	main := me.subCommands[mainSubCommandName]
	for _, option := range main.options {
		if option.LongName() == me.HelpName {
			return me.handleError(eInvalidHelpOption,
				"only auto-generated help is supported")
		} else if option.LongName() == me.VersionName {
			return me.handleError(eInvalidVersionOption,
				"only auto-generated version is supported")
		}
		if me.useLowerhForHelp && option.ShortName() == 'h' {
			me.useLowerhForHelp = false
		}
		if option.ShortName() == 'v' {
			usevForVersion = false
		} else if option.ShortName() == 'V' {
			seenV = true
		}
	}
	if me.VersionName != "" && !usevForVersion && !seenV {
		useVForVersion = true
	}
	if me.VersionName != "" && me.appVersion != "" {
		versionOpt := main.Flag(me.VersionName, "Show version and quit")
		if usevForVersion {
			versionOpt.SetShortName('v')
		} else if useVForVersion {
			versionOpt.SetShortName('V')
		}
		me.shortVersionName = versionOpt.ShortName()
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

func (me *Parser) isVersion(subcommand *SubCommand, option optioner) bool {
	if subcommand.longName == mainSubCommandName &&
		(option.LongName() == me.VersionName ||
			(me.shortVersionName != 0 && me.shortVersionName ==
				option.ShortName())) {
		me.onVersion() // doesn't return
		return true
	}
	return false
}

func (me *Parser) tokenize(args []string) (*SubCommand, []token, error) {
	var err error
	helpName := fmt.Sprintf("--%s", me.HelpName)
	state := me.initializeTokenState()
	tokens := make([]token, 0, len(args))
	for i, arg := range args {
		if me.isHelp(arg, helpName) {
			tokens = append(tokens, newHelpToken())
			continue
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
		subcommand:        me.subCommands[mainSubCommandName],
		subCommandForName: me.getSubCommandsForNames(),
		hasSubCommands:    len(me.subCommands) > 1,
		hadSubCommand:     false,
	}
	state.optionForLongName, state.optionForShortName =
		state.subcommand.optionsForNames()
	return state
}

func (me *Parser) isHelp(arg, helpName string) bool {
	if arg == helpName || (me.useLowerhForHelp && arg == "-h") {
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
			return tokens, me.handleError(eUnrecognizedLongOption,
				fmt.Sprintf("unrecognized option --%s", name))
		}
	} else { // --option
		option, ok := state.optionForLongName[name]
		if ok {
			tokens = append(tokens, newNameToken(name, option))
		} else {
			return tokens, me.handleError(eUnrecognizedShortOption1,
				fmt.Sprintf("unrecognized option --%s", name))
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
				return tokens, me.handleError(eUnexpectedValue, fmt.Sprintf(
					"unexpected value %s", rest))
			}
			break
		} else {
			return tokens, me.handleError(eUnrecognizedShortOption2,
				fmt.Sprintf("unrecognized option -%s", name))
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
		if long != mainSubCommandName {
			cmdForName[long] = command
			if command.ShortName() != 0 {
				cmdForName[string(command.ShortName())] = command
			}
		}
	}
	return cmdForName
}

func (me *Parser) onHelp(subcommand *SubCommand) {
	exitFunc(0, me.HelpText(subcommand.LongName()))
}

// HelpText is public only to aid testing
func (me *Parser) HelpText(name string) string {
	var text string
	if subcommand, ok := me.subCommands[name]; ok {
		if len(me.subCommands) == 1 { // No subcommands
			// show main help
			text = fmt.Sprintf("usage: %s TODO main", me.appName) // TODO
		} else { // Has subcommands
			if subcommand.longName == mainSubCommandName {
				// show main help with list of subcommands
				text = fmt.Sprintf("usage: %s TODO main + subs", me.appName) // TODO
			} else {
				// show this subcommand's help
				text = fmt.Sprintf("usage: %s TODO sub %s", me.appName, name) // TODO
			}
		}
		return text
	}
	panic(fmt.Sprintf("#%d: no main subcommand or subcommand", pBug))
}

func (me *Parser) onVersion() {
	exitFunc(0, fmt.Sprintf("%s v%s", me.appName, me.appVersion))
}

// VersionText is public only to aid testing

func (me *Parser) VersionText() string {
	return fmt.Sprintf("%s v%s", me.appName, me.appVersion)
}

func (me *Parser) checkPositionals() error {
	count := len(me.Positionals)
	ok := true
	switch me.positionalCount {
	case ZeroPositionals:
		if count > 0 {
			ok = false
		}
	case ZeroOrOnePositionals:
		if count > 1 {
			ok = false
		}
	case ZeroOrMorePositionals: // any count is valid
	case OnePositional:
		if count != 1 {
			ok = false
		}
	case OneOrMorePositionals:
		if count == 0 {
			ok = false
		}
	case TwoPositionals:
		if count != 2 {
			ok = false
		}
	case ThreePositionals:
		if count != 3 {
			ok = false
		}
	case FourPositionals:
		if count != 4 {
			ok = false
		}
	}
	if !ok {
		return me.handleError(eWrongPositionalCount,
			fmt.Sprintf("expected %s positional arguments, got %d",
				me.positionalCount, count))
	}
	return nil
}

func (me *Parser) checkValues(options []optioner) error {
	for _, option := range options {
		if msg := option.check(); msg != "" {
			return me.handleError(eInvalidOptionValue, msg)
		}
	}
	return nil
}

func (me *Parser) handleError(code int, msg string) error {
	exitFunc(2, fmt.Sprintf("error #%d: %s", code, msg))
	return nil // never returns
}

func (me *Parser) OnMissing(option optioner) error {
	if option.ShortName() != 0 {
		return me.handleError(eMissing,
			fmt.Sprintf("option -%c (or --%s) is required",
				option.ShortName(), option.LongName()))
	}
	return me.handleError(eMissing, fmt.Sprintf("option --%s is required",
		option.LongName()))
}

func (me *Parser) OnError(msg string) error {
	return me.handleError(eUser, msg)
}

func defaultExitFunc(exitCode int, msg string) {
	if exitCode == 0 {
		fmt.Println(msg)
	} else {
		fmt.Fprintln(os.Stderr, msg)
	}
	os.Exit(exitCode)
}

var exitFunc = defaultExitFunc
