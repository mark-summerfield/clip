// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

type SubCommand struct {
	longName  string
	shortName rune
	help      string
	options   []*Option
}

// Can't change long name or help after creation
func newMainSubCommand() *SubCommand {
	return &SubCommand{longName: "", shortName: noShortName, help: "",
		options: make([]*Option, 0)}
}

func newSubCommand(name, help string) *SubCommand {
	return &SubCommand{longName: name, shortName: noShortName, help: help,
		options: make([]*Option, 0)}
}

func (me *SubCommand) SetShortName(c rune) {
	me.shortName = c
}

func (me *SubCommand) Flag(name, help string) *Option {
	option := me.newOption(name, help, Flag)
	option.value = false
	option.valueCount = Zero
	return option
}

func (me *SubCommand) Int(name, help string, defaultValue int) *Option {
	option := me.newOption(name, help, Int)
	option.defaultValue = defaultValue
	return option
}

func (me *SubCommand) IntInRange(name, help string,
	minimum, maximum, defaultValue int) *Option {
	option := me.newOption(name, help, Int)
	option.validator = makeIntRangeValidator(minimum, maximum)
	option.defaultValue = defaultValue
	return option
}

func (me *SubCommand) Real(name, help string, defaultValue float64) *Option {
	option := me.newOption(name, help, Real)
	option.defaultValue = defaultValue
	return option
}

func (me *SubCommand) RealInRange(name, help string,
	minimum, maximum, defaultValue float64) *Option {
	option := me.newOption(name, help, Real)
	option.validator = makeRealRangeValidator(minimum, maximum)
	option.defaultValue = defaultValue
	return option
}

func (me *SubCommand) Str(name, help, defaultValue string) *Option {
	option := me.newOption(name, help, Str)
	option.defaultValue = defaultValue
	return option
}

func (me *SubCommand) Choice(name, help string, choices []string,
	defaultValue string) *Option {
	option := me.newOption(name, help, Str)
	option.validator = makeChoiceValidator(choices)
	option.defaultValue = defaultValue
	return option
}

func (me *SubCommand) Strs(name, help string) *Option {
	option := me.newOption(name, help, Strs)
	option.valueCount = OneOrMore
	return option
}

func (me *SubCommand) newOption(name, help string,
	valueType ValueType) *Option {
	option := newOption(name, help, valueType)
	me.options = append(me.options, option)
	return option
}

func (me *SubCommand) optionsForNames() (map[string]*Option,
	map[string]*Option) {
	optionForLongName := make(map[string]*Option, len(me.options))
	optionForShortName := make(map[string]*Option, len(me.options))
	for _, option := range me.options {
		if option.longName != "" {
			optionForLongName[option.longName] = option
		}
		if option.shortName != 0 {
			optionForShortName[string(option.shortName)] = option
		}
	}
	return optionForLongName, optionForShortName
}
