// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package clip

import (
	"fmt"
	"github.com/mark-summerfield/gong"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Parser struct {
	HelpName          string
	VersionName       string
	ShortDesc         string
	LongDesc          string
	EndNotes          string
	shortVersionName  rune
	appName           string
	appVersion        string
	options           []optioner
	firstDelayedError string
	Positionals       []string
	PositionalCount   PositionalCount
	PositionalHelp    string
	positionalVarName string
	useLowerhForHelp  bool
	width             int
}

// NewParser creates a new command line parser.
// It uses the executable's basename for the AppName and has no version
// option.
// See also NewParserVersion and NewParserUser.
func NewParser() Parser {
	return NewParserUser(appName(), "")
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
	return Parser{appName: appname, appVersion: strings.TrimSpace(version),
		options: make([]optioner, 0), PositionalCount: ZeroOrMorePositionals,
		positionalVarName: "FILE", HelpName: "help", VersionName: "version",
		useLowerhForHelp: true, width: GetWidth()}
}

func (me *Parser) AppName() string {
	return me.appName
}

func (me *Parser) SetAppName(appName string) {
	if appName != "" {
		me.appName = appName
	}
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

func (me *Parser) Flag(name, help string) *FlagOption {
	option, err := newFlagOption(name, help)
	me.registerNewOption(option, err)
	return option
}

func (me *Parser) Int(name, help string, theDefault int) *IntOption {
	option, err := newIntOption(name, help, theDefault)
	me.registerNewOption(option, err)
	return option
}

func (me *Parser) IntInRange(name, help string, minimum, maximum,
	theDefault int) *IntOption {
	option, err := newIntOption(name, help, theDefault)
	option.Validator = makeIntRangeValidator(minimum, maximum)
	me.registerNewOption(option, err)
	return option
}

func (me *Parser) Real(name, help string,
	theDefault float64) *RealOption {
	option, err := newRealOption(name, help, theDefault)
	me.registerNewOption(option, err)
	return option
}

func (me *Parser) RealInRange(name, help string, minimum, maximum,
	theDefault float64) *RealOption {
	option, err := newRealOption(name, help, theDefault)
	option.Validator = makeRealRangeValidator(minimum, maximum)
	me.registerNewOption(option, err)
	return option
}

func (me *Parser) Str(name, help, theDefault string) *StrOption {
	option, err := newStrOption(name, help, theDefault)
	me.registerNewOption(option, err)
	return option
}

func (me *Parser) Choice(name, help string, choices []string,
	theDefault string) *StrOption {
	option, err := newStrOption(name, help, theDefault)
	option.Validator = makeChoiceValidator(choices)
	me.registerNewOption(option, err)
	return option
}

func (me *Parser) Strs(name, help string) *StrsOption {
	option, err := newStrsOption(name, help)
	me.registerNewOption(option, err)
	return option
}

func (me *Parser) Ints(name, help string) *IntsOption {
	option, err := newIntsOption(name, help)
	me.registerNewOption(option, err)
	return option
}

func (me *Parser) Reals(name, help string) *RealsOption {
	option, err := newRealsOption(name, help)
	me.registerNewOption(option, err)
	return option
}

func (me *Parser) registerNewOption(option optioner, err error) {
	me.options = append(me.options, option)
	if err != nil && me.firstDelayedError == "" {
		me.firstDelayedError = err.Error()
	}
}

func (me *Parser) optionsForNames() (map[string]optioner,
	map[string]optioner) {
	optionForLongName := make(map[string]optioner, len(me.options))
	optionForShortName := make(map[string]optioner, len(me.options))
	for _, option := range me.options {
		if option.LongName() != "" {
			optionForLongName[option.LongName()] = option
		}
		if option.ShortName() != NoShortName {
			optionForShortName[string(option.ShortName())] = option
		}
	}
	return optionForLongName, optionForShortName
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
	tokens, err := me.tokenize(args)
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
			me.onHelp() // doesn't return
		} else if token.kind == nameTokenKind { // Option
			currentOption = token.option
			if me.isVersion(currentOption) { // may not return
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
	return me.checkValues()
}

func (me *Parser) prepareHelpAndVersionOptions() error {
	usevForVersion := true
	useVForVersion := false
	seenV := false
	for _, option := range me.options {
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
		versionOpt := me.Flag(me.VersionName, "Show version and quit")
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
	if me.firstDelayedError != "" {
		exitFunc(2, fmt.Sprintf("error %s", me.firstDelayedError))
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

func (me *Parser) isVersion(option optioner) bool {
	if option.LongName() == me.VersionName || (me.shortVersionName !=
		NoShortName && me.shortVersionName == option.ShortName()) {
		me.onVersion() // doesn't return
		return true
	}
	return false
}

func (me *Parser) tokenize(args []string) ([]token, error) {
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
				return tokens, err
			}
		} else if strings.HasPrefix(arg, "-") {
			if _, err := strconv.ParseFloat(arg, 64); err == nil {
				tokens = append(tokens, newValueToken(arg)) // -int | -real
			} else {
				tokens, err = me.handleShortOption(arg, tokens, &state)
				if err != nil {
					return tokens, err
				}
			}
		} else {
			tokens = append(tokens, newValueToken(arg))
		}
	}
	return tokens, nil
}

func (me *Parser) initializeTokenState() tokenState {
	state := tokenState{}
	state.optionForLongName, state.optionForShortName = me.optionsForNames()
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

func (me *Parser) onHelp() {
	text := me.usageLine()
	text = me.maybeWithDescriptionAndPositionals(text)
	text += me.optionsHelp()
	if me.EndNotes != "" {
		text += strings.Join(gong.TextWrap(me.EndNotes, me.width), "\n")
	}
	text = strings.TrimSuffix(text, "\n")
	exitFunc(0, text)
}

func (me *Parser) usageLine() string {
	text := fmt.Sprintf("usage: %s [OPTIONS]", me.appName)
	if me.PositionalCount != ZeroPositionals {
		text = fmt.Sprintf("%s %s", text,
			positionalCountText(me.PositionalCount, me.positionalVarName))
	}
	return text + "\n"
}

func (me *Parser) maybeWithDescriptionAndPositionals(text string) string {
	if me.LongDesc != "" {
		desc := gong.TextWrap(me.LongDesc, me.width)
		text = fmt.Sprintf("%s\n%s\n", text, strings.Join(desc, "\n"))
	}
	if me.PositionalCount != ZeroPositionals {
		posCountText := positionalCountText(me.PositionalCount,
			me.positionalVarName)
		text = fmt.Sprintf("%s\npositional arguments:\n%s%s", text,
			columnGap, posCountText)
		if me.PositionalHelp != "" {
			text += columnGap + ArgHelp(
				utf8.RuneCountInString(posCountText), me.width,
				me.PositionalHelp)
		} else {
			text += "\n"
		}
	}
	return text
}

func (me *Parser) optionsHelp() string {
	shorts := 0
	maxLeft := 0
	data := make([]datum, 0, len(me.options))
	for _, option := range me.options {
		n, arg := initialArgText(option)
		shorts += n
		arg += optArgText(option)
		lenArg := utf8.RuneCountInString(arg)
		if lenArg > maxLeft {
			maxLeft = lenArg
		}
		data = append(data, datum{arg: arg, lenArg: lenArg,
			help: option.Help()})

	}
	help := columnGap + "-h, --" + me.HelpName
	lenArg := utf8.RuneCountInString(help)
	if lenArg > maxLeft {
		maxLeft = lenArg
	}
	data = append(data, datum{arg: help, lenArg: lenArg,
		help: "Show help and quit"})
	gapWidth := utf8.RuneCountInString(columnGap)
	text := "\noptional arguments:\n"
	allFit := prepareOptionsData(maxLeft, gapWidth, me.width, shorts, data)
	text += optionsDataText(allFit, maxLeft, gapWidth, me.width, data)
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

func (me *Parser) checkValues() error {
	for _, option := range me.options {
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
	if option.ShortName() != NoShortName {
		return me.handleError(eMissing,
			fmt.Sprintf("option -%c (or --%s) is required",
				option.ShortName(), option.LongName()))
	}
	return me.handleError(eMissing, fmt.Sprintf("option --%s is required",
		option.LongName()))
}
