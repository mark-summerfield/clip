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

// TODO refactor into separate functions
func (me *Parser) ParseArgs(args []string) error {
	var err error
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
			if err = me.checkPositionals(&state); err != nil {
				return err
			}
			break
		} else if strings.HasPrefix(arg, "--") {
			name := strings.TrimPrefix(arg, "--")
			option, ok := state.optionForLongName[name]
			if ok {
				if err = me.handleOption(option, &state); err != nil {
					return err
				}
				continue // don't inc index
			} else {
				parts := strings.SplitN(name, "=", 2)
				if len(parts) == 2 { // --option=value
					option, ok := state.optionForLongName[parts[0]]
					if ok {
						state.args[0] = parts[1] // just keep the value
						if err = me.handleOption(option, &state); err != nil {
							return err
						}
						continue // don't inc index
					}
				}
				return me.handleError(fmt.Sprintf(
					"unrecognized option %s", arg))
			}
		} else if strings.HasPrefix(arg, "-") {
			name := strings.TrimPrefix(arg, "-")
			option, ok := state.optionForShortName[name]
			if ok {
				if err = me.handleOption(option, &state); err != nil {
					return err
				}
				continue // don't inc index
			} else {
				parts := strings.SplitN(name, "=", 2)
				if len(parts) == 2 { // -a=value or -abc=value
					state.args[0] = parts[1] // just keep the value
					flags := []rune(parts[0])
					if len(flags) == 1 { // -a=value
						name := string(flags[0])
						option, ok := state.optionForShortName[name]
						if ok {
							if err = me.handleOption(option, &state); err != nil {
								return err
							}
						}
					} else { // -abc=value
						for _, flag := range flags {
							name := string(flag)
							option, ok := state.optionForShortName[name]
							if ok {
								if err = me.handleOption(option, &state); err != nil {
									return err
								}
							}
						}
					}
					continue // don't inc index
				}
				return me.handleError(fmt.Sprintf(
					"unrecognized option %s", arg))
			}
		} else if state.hasSubCommands && !state.hadSubCommand {
			// is it a subcommand? - only allow one subcommand (excl. main)
			state.hadSubCommand = true
			cmd, ok := state.subCommandForName[arg]
			if ok {
				state.subcommand = cmd
				state.optionForLongName, state.optionForShortName =
					state.subcommand.optionsForNames()
			} else { // must be positionals from now on
				if err = me.checkPositionals(&state); err != nil {
					return err
				}
				break
			}
		} else { // handle positionals
			if err = me.checkPositionals(&state); err != nil {
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

func (me *Parser) getSubCommands() []string {
	keys := make([]string, 0, len(me.SubCommands))
	for key := range me.SubCommands {
		keys = append(keys, key)
	}
	return keys
}

func (me *Parser) checkPositionals(state *parserState) error {
	size := len(state.args)
	if size == 0 {
		if me.PositionalCount == One {
			return me.handleError(
				"expected one positional argument, got none")
		} else if me.PositionalCount == OneOrMore {
			return me.handleError(
				"expected at least one positional argument, got none")
		}
	} else if size == 1 && me.PositionalCount == Zero {
		return me.handleError("no positional arguments expected, got one")
	} else if size > 1 {
		if me.PositionalCount == Zero {
			return me.handleError(fmt.Sprintf(
				"no positional arguments expected, got %d", size))
		} else if me.PositionalCount == ZeroOrOne {
			return me.handleError(fmt.Sprintf(
				"expected at most one positional argument, got %d", size))
		} else if me.PositionalCount == One {
			return me.handleError(fmt.Sprintf(
				"expected one positional argument, got %d", size))
		}
	}
	me.Positionals = state.args
	return nil
}

func (me *Parser) handleOption(option *Option, state *parserState) error {
	// TODO set the option's value & if necessary keep reading args (& inc
	// index) until next - or --
	// If the option accepts anything other than Zero & the args[index] item
	// doesn't start with - then that's a value, ..., and so on
	// NOTE should leave the index ready at the next item
	return nil
}

func (me *Parser) handleError(msg string) error {
	msg = fmt.Sprintf("error: %s", msg)
	if me.QuitOnError {
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(2)
	}
	return errors.New(msg)
}
