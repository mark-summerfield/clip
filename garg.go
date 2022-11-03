// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package garg

import (
	_ "embed"
	"fmt"
	tsize "github.com/kopoli/go-terminal-size"
	"github.com/mark-summerfield/gong"
	"os"
	"path"
	"strconv"
	"strings"
)

//go:embed Version.dat
var Version string

type Parser struct {
	Positionals           []string
	HelpName              string
	VersionName           string
	Description           string
	EndNotes              string
	shortVersionName      rune
	appName               string
	appVersion            string
	subCommands           map[string]*SubCommand
	subCommandNames       []string // so that help is in creation order
	mainSubCommand        *SubCommand
	PositionalCount       PositionalCount
	PositionalDescription string
	positionalVarName     string
	useLowerhForHelp      bool
	width                 int
}

// NewParser creates a new command line parser.
// It uses the executable's basename for the AppName and has no version
// option.
// See also NewParserVersion and NewParserUser.
func NewParser() Parser {
	return NewParserUser(appName(), "")
}

func appName() string {
	if len(os.Args) > 0 {
		return path.Base(os.Args[0])
	}
	return "<app>"
}

// NewParserVersion creates a new command line parser.
// It uses the executable's basename for the AppName and a version
// option with the given version.
// See also NewParser and NewParserUser.
func NewParserVersion(version string) Parser {
	return NewParserUser(appName(), version)
}

// NewParserUser creates a new command line parser.
// If appname == "" the executable's basename will be used.
// If version == "" no version option will be available.
// See also NewParser and NewParserVersion.
func NewParserUser(appname, version string) Parser {
	if appname == "" {
		appname = appName()
	}
	width := 80
	size, err := tsize.GetSize()
	if err == nil && size.Width >= 38 {
		width = size.Width
	}
	mainSubCommand := newMainSubCommand()
	subcommands := make(map[string]*SubCommand)
	subcommands[mainSubCommandName] = mainSubCommand
	return Parser{appName: appname, appVersion: strings.TrimSpace(version),
		subCommands: subcommands, subCommandNames: make([]string, 0),
		PositionalCount:   ZeroOrMorePositionals,
		positionalVarName: "FILENAME", HelpName: "help",
		VersionName: "version", useLowerhForHelp: true,
		mainSubCommand: mainSubCommand, width: width}
}

func (me *Parser) AppName() string {
	return me.appName
}

func (me *Parser) Version() string {
	return me.appVersion
}

func (me *Parser) SetPositionalVarName(name string) error {
	if err := checkName(name, "positional var"); err != nil {
		return err
	}
	me.positionalVarName = name
	return nil
}

func (me *Parser) SubCommand(name, help string) *SubCommand {
	subcommand, err := newSubCommand(name, help)
	if err != nil && subcommand.firstDelayedError == "" {
		subcommand.firstDelayedError = err.Error()
	}
	me.subCommands[name] = subcommand
	me.subCommandNames = append(me.subCommandNames, name)
	return subcommand
}

func (me *Parser) Flag(name, help string) *FlagOption {
	return me.mainSubCommand.Flag(name, help)
}

func (me *Parser) Int(name, help string, theDefault int) *IntOption {
	return me.mainSubCommand.Int(name, help, theDefault)
}

func (me *Parser) IntInRange(name, help string, minimum, maximum,
	theDefault int) *IntOption {
	return me.mainSubCommand.IntInRange(name, help, minimum, maximum,
		theDefault)
}

func (me *Parser) Real(name, help string, theDefault float64) *RealOption {
	return me.mainSubCommand.Real(name, help, theDefault)
}

func (me *Parser) RealInRange(name, help string, minimum, maximum,
	theDefault float64) *RealOption {
	return me.mainSubCommand.RealInRange(name, help, minimum, maximum,
		theDefault)
}

func (me *Parser) Str(name, help, theDefault string) *StrOption {
	return me.mainSubCommand.Str(name, help, theDefault)
}

func (me *Parser) Choice(name, help string, choices []string,
	theDefault string) *StrOption {
	return me.mainSubCommand.Choice(name, help, choices, theDefault)
}

func (me *Parser) Strs(name, help string) *StrsOption {
	return me.mainSubCommand.Strs(name, help)
}

func (me *Parser) Ints(name, help string) *IntsOption {
	return me.mainSubCommand.Ints(name, help)
}

func (me *Parser) Reals(name, help string) *RealsOption {
	return me.mainSubCommand.Reals(name, help)
}

func (me *Parser) Parse() error {
	return me.ParseArgs(os.Args[1:])
}

func (me *Parser) ParseLine(line string) error {
	return me.ParseArgs(strings.Fields(line))
}

func (me *Parser) ParseArgs(args []string) error {
	if err := me.checkForDelayedError(); err != nil {
		return err
	}
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
			if me.isSubcommandHelp(subcommand, currentOption) { // may not return
				return nil
			}
			if me.isVersion(subcommand, currentOption) { // may not return
				return nil
			}
			if option, ok := currentOption.(*FlagOption); ok {
				option.value = true
			}
		} else { // Value
			if currentOption != nil && currentOption.wantsValue() {
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

func (me *Parser) checkForDelayedError() error {
	for _, subcommand := range me.subCommands {
		if subcommand.firstDelayedError != "" {
			exitFunc(2, fmt.Sprintf("error %s",
				subcommand.firstDelayedError))
		}
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

// This allows for user to write: `myapp asubcommand help` as well as
// `myapp asubcommand -h|--help` (handled elsewhere)
func (me *Parser) isSubcommandHelp(subcommand *SubCommand, option optioner) bool {
	if subcommand.longName != mainSubCommandName &&
		option.LongName() == me.HelpName {
		me.onHelp(subcommand) // doesn't return
		return true
	}
	return false
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
			return tokens, me.handleError(eUnrecognizedOption,
				fmt.Sprintf("unrecognized option --%s", name))
		}
	} else { // --option
		option, ok := state.optionForLongName[name]
		if ok {
			tokens = append(tokens, newNameToken(name, option))
		} else {
			return tokens, me.handleError(eUnrecognizedOption,
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
	isFlag := true
	for i, c := range text {
		name := string(c)
		option, ok := state.optionForShortName[name]
		if ok {
			tokens = append(tokens, newNameToken(name, option))
			_, isFlag = option.(*FlagOption)
			if !isFlag && i+1 < len(text) {
				value := text[i+1:] // -aValue -abcValue
				tokens = append(tokens, newValueToken(value))
			}
		} else if pendingValue == "" && !isFlag {
			last := len(tokens) - 1
			rest := text[i:]
			if last >= 0 && rest != tokens[last].text {
				return tokens, me.handleError(eUnexpectedValue, fmt.Sprintf(
					"unexpected value %s", rest))
			}
			break
		} else {
			return tokens, me.handleError(eUnrecognizedOption,
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
	text, err := me.helpText(subcommand.LongName())
	if err != nil {
		exitFunc(1, err.Error())
	}
	exitFunc(0, text)
}

// error should always be nil.
// The name should be "" for help on the main options (and list of
// subcommands if any), or the name of a subcommand.
func (me *Parser) helpText(name string) (string, error) {
	var text string
	if subcommand, ok := me.subCommands[name]; ok {
		if len(me.subCommands) == 1 { // No subcommands
			text = me.mainHelpText(subcommand)
		} else { // Has subcommands
			if subcommand.longName == mainSubCommandName {
				text = me.mainHelpTextWithSubCommands(subcommand)
			} else {
				text = me.subcommandHelpText(subcommand)
			}
		}
		return text, nil
	}
	return "", fmt.Errorf("#%d:BUG: no main subcommand or subcommand", eBug)
}

func (me *Parser) mainHelpText(subcommand *SubCommand) string {
	hasOptions := len(subcommand.options) > 0
	text := me.usageLine(hasOptions, len(me.subCommands) > 1, "")
	text = me.maybeWithDescriptionAndPositionals(text)
	//if hasOptions {
	//	text = me.optionsHelp(text, subcommand)
	//}
	return text
}

func (me *Parser) mainHelpTextWithSubCommands(subcommand *SubCommand) string {
	hasOptions := len(subcommand.options) > 0
	text := me.usageLine(hasOptions, len(me.subCommands) > 1, "")
	text = me.maybeWithDescriptionAndPositionals(text)
	if hasOptions {
		text = me.optionsHelp(text, subcommand)
	}
	// TODO list subcommands
	return text
}

func (me *Parser) subcommandHelpText(subcommand *SubCommand) string {
	hasOptions := len(subcommand.options) > 0
	text := me.usageLine(hasOptions, len(me.subCommands) > 1, "")
	if hasOptions {
		text = me.optionsHelp(text, subcommand)
	}
	// TODO
	return text
}

func (me *Parser) usageLine(hasOptions, hasSubCommands bool,
	subcommandName string) string {
	text := fmt.Sprintf("usage: %s", me.appName)
	if hasSubCommands {
		text = fmt.Sprintf("%s [SUBCOMMAND]", text)
	}
	if hasOptions {
		text = fmt.Sprintf("%s [OPTIONS]", text)
	}
	if subcommandName != "" {
		text = fmt.Sprintf("%s %s", text, subcommandName)
	}
	switch me.PositionalCount {
	case ZeroPositionals: // do nothing
	case ZeroOrOnePositionals:
		text = fmt.Sprintf("%s [%s]", text, me.positionalVarName)
	case ZeroOrMorePositionals: // any count is valid
		text = fmt.Sprintf("%s [%s [%s ...]]", text, me.positionalVarName,
			me.positionalVarName)
	case OnePositional:
		text = fmt.Sprintf("%s <%s>", text, me.positionalVarName)
	case OneOrMorePositionals:
		text = fmt.Sprintf("%s <%s> [%s [%s ...]]", text,
			me.positionalVarName, me.positionalVarName,
			me.positionalVarName)
	case TwoPositionals:
		text = fmt.Sprintf("%s <%s> <%s>", text, me.positionalVarName,
			me.positionalVarName)
	case ThreePositionals:
		text = fmt.Sprintf("%s <%s> <%s> <%s>", text, me.positionalVarName,
			me.positionalVarName, me.positionalVarName)
	case FourPositionals:
		text = fmt.Sprintf("%s <%s> <%s> <%s> <%s>", text,
			me.positionalVarName, me.positionalVarName,
			me.positionalVarName, me.positionalVarName)
	}
	return text + "\n"
}

func (me *Parser) maybeWithDescriptionAndPositionals(text string) string {
	if me.Description != "" {
		desc := gong.TextWrap(me.Description, me.width)
		text = fmt.Sprintf("%s\n%s\n", text, strings.Join(desc, "\n"))
	}
	if me.PositionalCount != ZeroPositionals {
		text = fmt.Sprintf("%s\narguments:\n  ", text)
	}
	switch me.PositionalCount {
	case ZeroPositionals: // do nothing
	case ZeroOrOnePositionals:
		text = fmt.Sprintf("%s [%s]", text, me.positionalVarName)
	case ZeroOrMorePositionals: // any count is valid
		text = fmt.Sprintf("%s [%s [%s ...]]", text, me.positionalVarName,
			me.positionalVarName)
	case OnePositional:
		text = fmt.Sprintf("%s <%s>", text, me.positionalVarName)
	case OneOrMorePositionals:
		text = fmt.Sprintf("%s <%s> [%s [%s ...]]", text,
			me.positionalVarName, me.positionalVarName,
			me.positionalVarName)
	case TwoPositionals:
		text = fmt.Sprintf("%s <%s> <%s>", text, me.positionalVarName,
			me.positionalVarName)
	case ThreePositionals:
		text = fmt.Sprintf("%s <%s> <%s> <%s>", text, me.positionalVarName,
			me.positionalVarName, me.positionalVarName)
	case FourPositionals:
		text = fmt.Sprintf("%s <%s> <%s> <%s> <%s>", text,
			me.positionalVarName, me.positionalVarName,
			me.positionalVarName, me.positionalVarName)
	}
	if me.PositionalCount != ZeroPositionals {
		text = fmt.Sprintf("%s  %s\n", text, me.PositionalDescription)
	}
	return text
}

func (me *Parser) optionsHelp(text string, subcommand *SubCommand) string {
	/*
		maxFirst := 0
		maxSecond := 0
		pairs := make([]pair, 0, len(subcommand.options))
		for _, option := range subcommand.options {
			// TODO first is short (if present) long (args depending on
			// valuecount)
			// second is desc
		}
	*/
	return text
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
	switch me.PositionalCount {
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
				me.PositionalCount, count))
	}
	return nil
}

func (me *Parser) checkValues(options []optioner) error {
	for _, option := range options {
		if msg := option.check(); msg != "" {
			return me.handleError(eInvalidValue, msg)
		}
	}
	return nil
}

func (me *Parser) handleError(code int, msg string) error {
	exitFunc(2, fmt.Sprintf("error #%d: %s", code, msg))
	return nil // never returns
}

func (me *Parser) OnError(err error) {
	exitFunc(2, err.Error())
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

func defaultExitFunc(exitCode int, msg string) {
	if exitCode == 0 {
		fmt.Println(msg)
	} else {
		fmt.Fprintln(os.Stderr, msg)
	}
	os.Exit(exitCode)
}

var exitFunc = defaultExitFunc
