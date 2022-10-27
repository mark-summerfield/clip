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
	AutoVersionOption bool
	SubCommands       map[string]*SubCommand
	PositionalCount   ValueCount
	PositionalVarName string
	Positionals       []string
}

func NewParser(appname, version string) Parser {
	subcommands := make(map[string]*SubCommand)
	subcommands[mainSubCommand] = newMainSubCommand()
	return Parser{AppName: appname, AppVersion: version,
		QuitOnError: true, AutoVersionOption: true,
		SubCommands: subcommands}
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
	me.maybeAddVersionOption()
	state := parserState{
		subcommand:        me.SubCommands[mainSubCommand],
		subCommandForName: me.getSubCommandsForNames(),
		hasSubCommands:    len(me.SubCommands) > 1,
		hadSubCommand:     false,
		args:              args,
	}
	state.optionForLongName, state.optionForShortName =
		state.subcommand.optionsForNames()
	for state.index < len(state.args) {
		arg := state.args[state.index]
		if arg == "--" { // end of options
			state.index++
			if err := me.checkPositionals(&state); err != nil {
				return err
			}
			break
		} else if strings.HasPrefix(arg, "--") { // long option
			if err := me.handleLongPrefix(arg, &state); err != nil {
				return err
			}
		} else if strings.HasPrefix(arg, "-") { // short option
			if err := me.handleShortPrefix(arg, &state); err != nil {
				return err
			}
		} else if state.hasSubCommands && !state.hadSubCommand { // subcmd?
			do_break, err := me.handlePossibleSubcommand(arg, &state)
			if err != nil {
				return err
			}
			if do_break {
				break
			}
		} else { // positionals
			if err := me.checkPositionals(&state); err != nil {
				return err
			}
			break
		}
		state.index++
	}
	return nil
}

func (me *Parser) maybeAddVersionOption() {
	if me.AutoVersionOption {
		seen_v := false
		seen_V := false
		main := me.SubCommands[mainSubCommand]
		for _, option := range main.Options {
			if option.ShortName == 'v' {
				seen_v = true
			} else if option.ShortName == 'V' {
				seen_V = true
			}
			if strings.EqualFold(option.LongName, "version") {
				return // user has added version option themselves
			}
		}
		option := me.newOption("version", "Print version and quit", Flag)
		option.ValueCount = Zero
		if seen_v {
			if seen_V {
				option.ShortName = noShortName
			} else {
				option.ShortName = 'V'
			}
		}
		main.Options = append(main.Options, option)
	}
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

func (me *Parser) handleLongPrefix(arg string, state *parserState) error {
	name := strings.TrimPrefix(arg, "--")
	option, ok := state.optionForLongName[name]
	if ok { // --option
		if err := me.handleOption(option, "", state); err != nil {
			return err
		}
	} else {
		parts := strings.SplitN(name, "=", 2)
		if len(parts) == 2 { // --option=value
			option, ok := state.optionForLongName[parts[0]]
			if ok {
				if err := me.handleOption(option, parts[1],
					state); err != nil {
					return err
				}
				return nil
			}
		}
		return me.handleError(20, fmt.Sprintf("unrecognized option %s",
			arg))
	}
	return nil
}

func (me *Parser) handleShortPrefix(arg string, state *parserState) error {
	name := strings.TrimPrefix(arg, "-")
	option, ok := state.optionForShortName[name]
	if ok { // -o
		if err := me.handleOption(option, "", state); err != nil {
			return err
		}
	} else {
		parts := strings.SplitN(name, "=", 2)
		if len(parts) == 2 { // -a=value or -abc=value
			flags := []rune(parts[0])
			if len(flags) == 1 { // -a=value
				name := string(flags[0])
				option, ok := state.optionForShortName[name]
				if ok {
					if err := me.handleOption(option, parts[1],
						state); err != nil {
						return err
					}
				}
			} else { // -abc=value
				for _, flag := range flags {
					name := string(flag)
					option, ok := state.optionForShortName[name]
					if ok {
						if err := me.handleOption(option, parts[1],
							state); err != nil {
							return err
						}
					}
				}
			}
			return nil
		} else { // -abc or -aValue
			for _, flag := range name {
				name = string(flag)
				option, ok := state.optionForShortName[name]
				if ok {
					if err := me.handleOption(option, "",
						state); err != nil {
						return err
					}
				} else {
					// TODO handle -aValue case
					return me.handleError(30,
						fmt.Sprintf("unrecognized option %s", arg))
				}
			}
		}
	}
	return nil
}

func (me *Parser) handlePossibleSubcommand(arg string,
	state *parserState) (bool, error) {
	// is it a subcommand? - only allow one subcommand (excl. main)
	state.hadSubCommand = true
	cmd, ok := state.subCommandForName[arg]
	if ok {
		state.subcommand = cmd
		state.optionForLongName, state.optionForShortName =
			state.subcommand.optionsForNames()
	} else { // must be positionals from now on
		if err := me.checkPositionals(state); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (me *Parser) handleOption(option *Option, value string,
	state *parserState) error {
	// TODO set the option's value & if necessary keep reading args (& inc
	// index) until next - or --
	// If the option accepts anything other than Zero & the args[index] item
	// doesn't start with - then that's a value, ..., and so on
	// NOTE should leave the index ready at the next item
	fmt.Printf("handleOption() %v %v %v\n", *option, value,
		state.args[state.index])
	return nil
}

func (me *Parser) handleError(code int, msg string) error {
	msg = fmt.Sprintf("error #%d: %s", code, msg)
	if me.QuitOnError {
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(2)
	}
	return errors.New(msg)
}
