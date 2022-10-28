// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Parser struct {
	AppName           string
	AppVersion        string
	QuitOnError       bool
	SubCommands       map[string]*SubCommand
	PositionalCount   ValueCount
	PositionalVarName string
	Positionals       []string
	HelpName          string
	use_h_for_help    bool
	VersionName       string
	use_v_for_version bool
	use_V_for_version bool
}

func NewParser(appname, version string) Parser {
	subcommands := make(map[string]*SubCommand)
	subcommands[mainSubCommand] = newMainSubCommand()
	return Parser{AppName: appname, AppVersion: version,
		QuitOnError: true, SubCommands: subcommands,
		PositionalCount: ZeroOrMore, PositionalVarName: "FILENAME",
		HelpName: "help", use_h_for_help: true, VersionName: "version",
		use_v_for_version: true,
	}
}

func (me *Parser) SubCommand(name, help string) *SubCommand {
	subcommand := newSubCommand(name, help)
	me.SubCommands[name] = subcommand
	return subcommand
}

func (me *Parser) Flag(name, help string) *Option {
	option := me.newOption(name, help, Flag)
	option.ValueCount = Zero
	return option
}

func (me *Parser) Int(name, help string) *Option {
	return me.newOption(name, help, Int)
}

func (me *Parser) IntInRange(name, help string,
	minimum, maximum int) *Option {
	option := me.newOption(name, help, Int)
	option.Validator = makeIntRangeValidator(minimum, maximum)
	return option
}

func (me *Parser) Real(name, help string) *Option {
	return me.newOption(name, help, Real)
}

func (me *Parser) RealInRange(name, help string,
	minimum, maximum float64) *Option {
	option := me.newOption(name, help, Real)
	option.Validator = makeRealRangeValidator(minimum, maximum)
	return option
}

func (me *Parser) Str(name, help string) *Option {
	return me.newOption(name, help, Str)
}

func (me *Parser) Choice(name, help string, choices []string) *Option {
	option := me.newOption(name, help, Str)
	option.Validator = makeChoiceValidator(choices)
	return option
}

func (me *Parser) Strs(name, help string) *Option {
	option := me.newOption(name, help, Strs)
	option.ValueCount = OneOrMore
	return option
}

func (me *Parser) newOption(name, help string, valueType ValueType) *Option {
	option := newOption(name, help, valueType)
	me.SubCommands[mainSubCommand].Options = append(
		me.SubCommands[mainSubCommand].Options, option)
	return option
}

func (me *Parser) Parse() error {
	return me.ParseArgs(os.Args[1:])
}

func (me *Parser) ParseLine(line string) error {
	return me.ParseArgs(strings.Fields(line))
}

func (me *Parser) ParseArgs(args []string) error {
	me.prepareHelpAndVersionOptions()
	tokens, err := me.tokenize(args)
	if err != nil {
		return err
	} else { // TODO delete this else clause
		fmt.Printf("TOKENS for \"%s\" %s\n", strings.Join(args, " "),
			tokens)
	}
	var currentOption *Option
	expect := Zero // ValueCount - how many values we expect to follow opt
	index := 0
	for {
		if index >= len(tokens) {
			break
		}
		token := tokens[index]
		index++
		if !token.isValue() { // Option
			currentOption = token.option
			expect = currentOption.ValueCount
			if currentOption.ValueType == Flag {
				currentOption.Value = true
			} else {
				currentOption.Value = currentOption.DefaultValue
			}
		} else { // Value
			if expect == Zero {
				return me.handleError(20,
					fmt.Sprintf("unexpected value: %s", token.text))
			}
		}
		// TODO
	}
	// TODO check for absent Required options that don't have a DefaultValue
	return nil
}

func (me *Parser) prepareHelpAndVersionOptions() {
	seen_V := false
	main := me.SubCommands[mainSubCommand]
	for _, option := range main.Options {
		if option.LongName == me.HelpName {
			panic("only auto-generated help is supported")
		} else if option.LongName == me.VersionName {
			panic("only auto-generated version is supported")
		}
		if me.use_h_for_help && option.ShortName == 'h' {
			me.use_h_for_help = false
		}
		if option.ShortName == 'v' {
			me.use_v_for_version = false
		} else if option.ShortName == 'V' {
			seen_V = true
		}
	}
	if !me.use_v_for_version && !seen_V {
		me.use_V_for_version = true
	}
}

// TODO refactor into cases
func (me *Parser) tokenize(args []string) ([]token, error) {
	state := tokenState{
		subcommand:        me.SubCommands[mainSubCommand],
		subCommandForName: me.getSubCommandsForNames(),
		hasSubCommands:    len(me.SubCommands) > 1,
		hadSubCommand:     false,
	}
	state.optionForLongName, state.optionForShortName =
		state.subcommand.optionsForNames()
	tokens := make([]token, 0, len(args))
	for i, arg := range args {
		if arg == me.HelpName || (me.use_h_for_help && arg == "-h") {
			me.onHelp() // doesn't return
			return nil, nil
		}
		if arg == me.VersionName || (me.use_v_for_version && arg == "-v") ||
			(me.use_V_for_version && arg == "-V") {
			me.onVersion() // doesn't return
			return nil, nil
		}
		if arg == "--" { // --
			tokens = append(tokens, newPositionalsFollowToken())
			for _, v := range args[i+1:] {
				tokens = append(tokens, newValueToken(v))
			}
			break
		}
		if strings.HasPrefix(arg, "--") { // --option --option=value
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
		} else if strings.HasPrefix(arg, "-") {
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
					if option.ValueType != Flag && i+1 < len(text) {
						value := text[i+1:] // -aValue -abcValue
						tokens = append(tokens, newValueToken(value))
					}
				} else if pendingValue == "" {
					size := len(tokens)
					rest := text[i:]
					if size > 0 && rest != tokens[size-1].text {
						return tokens, me.handleError(14, fmt.Sprintf(
							"unexpected value %s", rest))
					}
					break
				} else {
					return tokens, me.handleError(16, fmt.Sprintf(
						"unrecognized option -%s", name))
				}
			}
			if pendingValue != "" {
				tokens = append(tokens, newValueToken(pendingValue))
			}
		} else if state.hasSubCommands && !state.hadSubCommand { // subcmd?
			// is it a subcommand? - only allow one subcommand (excl. main)
			state.hadSubCommand = true
			cmd, ok := state.subCommandForName[arg]
			if ok {
				state.subcommand = cmd
				state.optionForLongName, state.optionForShortName =
					state.subcommand.optionsForNames()
			} else { // value
				tokens = append(tokens, newValueToken(arg))
			}
		} else {
			tokens = append(tokens, newValueToken(arg))
		}
	}
	return tokens, nil
}

func (me *Parser) getSubCommandsForNames() map[string]*SubCommand {
	cmdForName := make(map[string]*SubCommand, len(me.SubCommands)*2)
	for long, command := range me.SubCommands {
		if long != mainSubCommand {
			cmdForName[long] = command
			if command.ShortName != 0 {
				cmdForName[string(command.ShortName)] = command
			}
		}
	}
	return cmdForName
}

func (me *Parser) onHelp() {
	fmt.Printf("usage: %s TODO", me.AppName)
	os.Exit(0)
}

func (me *Parser) onVersion() {
	fmt.Printf("%s v%s", me.AppName, me.AppVersion)
	os.Exit(0)
}

func (me *Parser) handleError(code int, msg string) error {
	msg = fmt.Sprintf("error #%d: %s", code, msg)
	if me.QuitOnError {
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(2)
	}
	return errors.New(msg)
}

/* TODO
func (me *Parser) checkPositionals(state *parserState) error {
	size := len(state.args)
	if size == 0 {
		if me.PositionalCount == One {
			return me.handleError(10,
				"expected one positional argument, got none")
		} else if me.PositionalCount == OneOrMore {
			return me.handleError(11,
				"expected at least one positional argument, got none")
		}
	} else if size == 1 && me.PositionalCount == Zero {
		return me.handleError(12,
			"no positional arguments expected, got one")
	} else if size > 1 {
		if me.PositionalCount == Zero {
			return me.handleError(13, fmt.Sprintf(
				"no positional arguments expected, got %d", size))
		} else if me.PositionalCount == ZeroOrOne {
			return me.handleError(14, fmt.Sprintf(
				"expected at most one positional argument, got %d", size))
		} else if me.PositionalCount == One {
			return me.handleError(15, fmt.Sprintf(
				"expected one positional argument, got %d", size))
		}
	}
	me.Positionals = state.args
	return nil
}
*/

/* TODO
func (me *Parser) handleOption(option *Option, value string,
	state *parserState) error {
	if option.LongName == "version" {
		me.onVersion() // never returns
		return nil
	}
	if option.ValueType == Flag {
		if value == "" {
			option.Value = true
		} else {
			return me.handleError(40, fmt.Sprintf(
				"unexpected value for flag %s: %s", option.LongName, value))
		}
	}

	// TODO set the option's value & if necessary keep reading args (& inc
	// index) until next - or --
	// If the option accepts anything other than Zero & the args[index] item
	// doesn't start with - then that's a value, ..., and so on
	// NOTE should leave the index ready at the next item

	// DEBUG
	var next string
	if state.index < len(state.args) {
		next = state.args[state.index]
	}
	fmt.Printf("handleOption() %v %#v %#v\n", *option, value, next)
	// END DEBUG

	return nil
}
*/
