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
		SubCommands: subcommands, PositionalCount: ZeroOrMore,
		PositionalVarName: "FILENAME"}
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
	tokens, err := me.tokenize(args)
	if err == nil {
		fmt.Printf("TOKENS for %s: %s\n", args, tokens)
	} else {
		return err
	}
	index := 0
	for {
		if index >= len(tokens) {
			break
		}
		token := tokens[index]
		index++
		if !token.isValue { // Name
		} else { // Value
		}
		// TODO
	}
	return nil
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

func (me *Parser) tokenize(args []string) ([]Token, error) {
	state := tokenState{
		subcommand:        me.SubCommands[mainSubCommand],
		subCommandForName: me.getSubCommandsForNames(),
		hasSubCommands:    len(me.SubCommands) > 1,
		hadSubCommand:     false,
	}
	state.optionForLongName, state.optionForShortName =
		state.subcommand.optionsForNames()
	tokens := make([]Token, 0, len(args))
	for i, arg := range args {
		if arg == "--" { // --
			for _, v := range args[i+1:] {
				tokens = append(tokens, NewValueToken(v))
			}
			break
		}
		if strings.HasPrefix(arg, "--") { // --option --option=value
			name := strings.TrimPrefix(arg, "--")
			parts := strings.SplitN(name, "=", 2)
			if len(parts) == 2 { // --option=value
				tokens = append(tokens, NewNameToken(parts[0]))
				tokens = append(tokens, NewValueToken(parts[1]))
			} else { // --option
				tokens = append(tokens, NewNameToken(name))
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
					tokens = append(tokens, NewNameToken(name))
					if option.ValueType != Flag && i+1 < len(text) {
						value := text[i+1:] // -aValue -abcValue
						tokens = append(tokens, NewValueToken(value))
					}
				}
			}
			if pendingValue != "" {
				tokens = append(tokens, NewValueToken(pendingValue))
			}
		} else if state.hasSubCommands && !state.hadSubCommand { // subcmd?
			// is it a subcommand? - only allow one subcommand (excl. main)
			state.hadSubCommand = true
			cmd, ok := state.subCommandForName[arg]
			if ok {
				state.subcommand = cmd
				state.optionForLongName, state.optionForShortName =
					state.subcommand.optionsForNames()
			} else { // positionals
				tokens = append(tokens, NewValueToken(arg))
				for _, v := range args[i+1:] {
					tokens = append(tokens, NewValueToken(v))
				}
				break
			}
		} else {
			tokens = append(tokens, NewValueToken(arg))
		}
	}
	return tokens, nil
}
