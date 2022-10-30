// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

type SubCommand struct {
	longName  string
	shortName rune
	help      string
	options   []Optioner
}

// Can't change long name or help after creation
func newMainSubCommand() *SubCommand {
	return &SubCommand{longName: "", shortName: noShortName, help: "",
		options: make([]Optioner, 0)}
}

func newSubCommand(name, help string) *SubCommand {
	return &SubCommand{longName: name, shortName: noShortName, help: help,
		options: make([]Optioner, 0)}
}

func (me *SubCommand) LongName() string {
	return me.longName
}

func (me *SubCommand) ShortName() rune {
	return me.shortName
}

func (me *SubCommand) SetShortName(c rune) {
	me.shortName = c
}

func (me *SubCommand) Flag(name, help string) *FlagOption {
	option := newFlagOption(name, help)
	me.options = append(me.options, option)
	return option
}

func (me *SubCommand) Int(name, help string, theDefaultValue int) *IntOption {
	option := newIntOption(name, help, theDefaultValue)
	me.options = append(me.options, option)
	return option
}

func (me *SubCommand) IntInRange(name, help string,
	minimum, maximum, theDefaultValue int) *IntOption {
	option := newIntOption(name, help, theDefaultValue)
	option.validator = makeIntRangeValidator(minimum, maximum)
	me.options = append(me.options, option)
	return option
}

func (me *SubCommand) Real(name, help string,
	theDefaultValue float64) *RealOption {
	option := newRealOption(name, help, theDefaultValue)
	me.options = append(me.options, option)
	return option
}

func (me *SubCommand) RealInRange(name, help string,
	minimum, maximum, theDefaultValue float64) *RealOption {
	option := newRealOption(name, help, theDefaultValue)
	option.validator = makeRealRangeValidator(minimum, maximum)
	me.options = append(me.options, option)
	return option
}

func (me *SubCommand) Str(name, help, theDefaultValue string) *StrOption {
	option := newStrOption(name, help, theDefaultValue)
	me.options = append(me.options, option)
	return option
}

func (me *SubCommand) Choice(name, help string, choices []string,
	theDefaultValue string) *StrOption {
	option := newStrOption(name, help, theDefaultValue)
	option.validator = makeChoiceValidator(choices)
	me.options = append(me.options, option)
	return option
}

func (me *SubCommand) Strs(name, help string) *StrsOption {
	option := newStrsOption(name, help)
	me.options = append(me.options, option)
	return option
}

func (me *SubCommand) optionsForNames() (map[string]Optioner,
	map[string]Optioner) {
	optionForLongName := make(map[string]Optioner, len(me.options))
	optionForShortName := make(map[string]Optioner, len(me.options))
	for _, option := range me.options {
		if option.LongName() != "" {
			optionForLongName[option.LongName()] = option
		}
		if option.ShortName() != 0 {
			optionForShortName[string(option.ShortName())] = option
		}
	}
	return optionForLongName, optionForShortName
}
