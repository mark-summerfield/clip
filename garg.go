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
	option.value = false
	option.valueCount = Zero
	return option
}

func (me *Parser) Int(name, help string) *Option {
	return me.newOption(name, help, Int)
}

func (me *Parser) IntInRange(name, help string,
	minimum, maximum int) *Option {
	option := me.newOption(name, help, Int)
	option.validator = makeIntRangeValidator(minimum, maximum)
	return option
}

func (me *Parser) Real(name, help string) *Option {
	return me.newOption(name, help, Real)
}

func (me *Parser) RealInRange(name, help string,
	minimum, maximum float64) *Option {
	option := me.newOption(name, help, Real)
	option.validator = makeRealRangeValidator(minimum, maximum)
	return option
}

func (me *Parser) Str(name, help string) *Option {
	return me.newOption(name, help, Str)
}

func (me *Parser) Choice(name, help string, choices []string) *Option {
	option := me.newOption(name, help, Str)
	option.validator = makeChoiceValidator(choices)
	return option
}

func (me *Parser) Strs(name, help string) *Option {
	option := me.newOption(name, help, Strs)
	option.valueCount = OneOrMore
	return option
}

func (me *Parser) newOption(name, help string, valueType ValueType) *Option {
	option := newOption(name, help, valueType)
	me.SubCommands[mainSubCommand].options = append(
		me.SubCommands[mainSubCommand].options, option)
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
	subcommand, tokens, err := me.tokenize(args)
	if err != nil {
		return err
	}
	var currentOption *Option
	inPositionals := false
	expect := Zero // ValueCount - how many values we expect to follow opt
	for _, token := range tokens {
		if token.positionalsFollow {
			inPositionals = true
		}
		if inPositionals {
			me.addPositional(token.text)
		} else if !token.isValue() { // Option
			currentOption = token.option
			expect = currentOption.valueCount
			if currentOption.valueType == Flag {
				currentOption.value = true
			} else if currentOption.valueType != Strs {
				currentOption.value = currentOption.defaultValue
			}
		} else { // Value
			switch expect {
			case Zero:
				inPositionals = me.addPositional(token.text)
			case ZeroOrOne:
				if currentOption.Size() == 1 {
					inPositionals = me.addPositional(token.text)
				} else {
					currentOption.AddValue(token.text)
				}
			case ZeroOrMore:
				currentOption.AddValue(token.text)
			case One:
				if currentOption.Size() == 0 {
					currentOption.AddValue(token.text)
				} else {
					inPositionals = me.addPositional(token.text)
				}
			case OneOrMore:
				currentOption.AddValue(token.text)
			default:
				panic("invalid ValueCount #2")
			}
		}
	}
	if err := me.checkPositionals(); err != nil {
		return err
	}
	return me.checkValues(subcommand.options)
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
	main := me.SubCommands[mainSubCommand]
	for _, option := range main.options {
		if option.longName == me.HelpName {
			panic("only auto-generated help is supported")
		} else if option.longName == me.VersionName {
			panic("only auto-generated version is supported")
		}
		if me.use_h_for_help && option.shortName == 'h' {
			me.use_h_for_help = false
		}
		if option.shortName == 'v' {
			me.use_v_for_version = false
		} else if option.shortName == 'V' {
			seen_V = true
		}
	}
	if !me.use_v_for_version && !seen_V {
		me.use_V_for_version = true
	}
}

func (me *Parser) tokenize(args []string) (*SubCommand, []token, error) {
	var err error
	state := me.initializeTokenState()
	tokens := make([]token, 0, len(args))
	for i, arg := range args {
		if arg == me.HelpName || (me.use_h_for_help && arg == "-h") {
			me.onHelp() // doesn't return
			return nil, nil, nil
		}
		if arg == me.VersionName || (me.use_v_for_version && arg == "-v") ||
			(me.use_V_for_version && arg == "-V") {
			me.onVersion() // doesn't return
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
			tokens, err = me.handleShortOption(arg, tokens, &state)
			if err != nil {
				return state.subcommand, tokens, err
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
		subcommand:        me.SubCommands[mainSubCommand],
		subCommandForName: me.getSubCommandsForNames(),
		hasSubCommands:    len(me.SubCommands) > 1,
		hadSubCommand:     false,
	}
	state.optionForLongName, state.optionForShortName =
		state.subcommand.optionsForNames()
	return state
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
			if option.valueType != Flag && i+1 < len(text) {
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
	cmdForName := make(map[string]*SubCommand, len(me.SubCommands)*2)
	for long, command := range me.SubCommands {
		if long != mainSubCommand {
			cmdForName[long] = command
			if command.shortName != 0 {
				cmdForName[string(command.shortName)] = command
			}
		}
	}
	return cmdForName
}

func (me *Parser) onHelp() {
	fmt.Printf("usage: %s TODO", me.AppName) // TODO
	os.Exit(0)
}

func (me *Parser) onVersion() {
	fmt.Printf("%s v%s", me.AppName, me.AppVersion)
	os.Exit(0)
}

func (me *Parser) checkPositionals() error {
	size := len(me.Positionals)
	switch me.PositionalCount {
	case Zero:
		if size > 0 {
			return me.handleError(20,
				fmt.Sprintf("expected no positional arguments, got %d",
					size))
		}
	case ZeroOrOne:
		if size > 1 {
			return me.handleError(22,
				fmt.Sprintf(
					"expected zero or one positional arguments, got %d",
					size))
		}
	case ZeroOrMore: // any size is valid
	case One:
		if size != 1 {
			return me.handleError(24,
				fmt.Sprintf("expected one positional argument, got %d",
					size))
		}
	case OneOrMore:
		if size == 0 {
			return me.handleError(26,
				fmt.Sprintf(
					"expected at least one positional argument, got %d",
					size))
		}
	default:
		panic("invalid ValueCount #3")
	}
	return nil
}

func (me *Parser) checkValues(options []*Option) error {
	for _, option := range options {
		option.setDefaultIfAppropriate()
		if option.required && option.value == nil {
			return me.handleError(30,
				fmt.Sprintf("expected a value for %s", option.longName))
		}
		size := option.Size()
		switch option.valueCount {
		case Zero:
			if option.valueType != Flag {
				panic(fmt.Sprintf("nonflag option %s with zero ValueCount",
					option.longName))
			}
		case ZeroOrOne:
			if size > 1 {
				return me.handleError(32,
					fmt.Sprintf(
						"expected zero or one values for %s, got %d",
						option.longName, size))
			}
		case ZeroOrMore: // any size is valid
		case One:
			if size != 1 {
				return me.handleError(34,
					fmt.Sprintf("expected exactly one value for %s, got %d",
						option.longName, size))
			}
		case OneOrMore:
			if size == 0 {
				return me.handleError(36,
					fmt.Sprintf(
						"expected at least one value for %s, got %d",
						option.longName, size))
			}
		default:
			panic("invalid ValueCount #4")
		}
	}
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
