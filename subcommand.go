// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

type SubCommand struct {
	longName  string
	shortName rune
	help      string
	options   []optioner
}

func newMainSubCommand() *SubCommand {
	return &SubCommand{longName: "", shortName: noShortName, help: "",
		options: make([]optioner, 0)}
}

func newSubCommand(name, help string) *SubCommand {
	return &SubCommand{longName: name, shortName: noShortName, help: help,
		options: make([]optioner, 0)}
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
	me.registerNewOption(option)
	return option
}

func (me *SubCommand) Int(name, help string, theDefault int) *IntOption {
	option := newIntOption(name, help, theDefault)
	me.registerNewOption(option)
	return option
}

func (me *SubCommand) IntInRange(name, help string,
	minimum, maximum, theDefault int) *IntOption {
	option := newIntOption(name, help, theDefault)
	option.validator = makeIntRangeValidator(minimum, maximum)
	me.registerNewOption(option)
	return option
}

func (me *SubCommand) Real(name, help string,
	theDefault float64) *RealOption {
	option := newRealOption(name, help, theDefault)
	me.registerNewOption(option)
	return option
}

func (me *SubCommand) RealInRange(name, help string,
	minimum, maximum, theDefault float64) *RealOption {
	option := newRealOption(name, help, theDefault)
	option.validator = makeRealRangeValidator(minimum, maximum)
	me.registerNewOption(option)
	return option
}

func (me *SubCommand) Str(name, help, theDefault string) *StrOption {
	option := newStrOption(name, help, theDefault)
	me.registerNewOption(option)
	return option
}

func (me *SubCommand) Choice(name, help string, choices []string,
	theDefault string) *StrOption {
	option := newStrOption(name, help, theDefault)
	option.validator = makeChoiceValidator(choices)
	me.registerNewOption(option)
	return option
}

func (me *SubCommand) Strs(name, help string) *StrsOption {
	option := newStrsOption(name, help)
	me.registerNewOption(option)
	return option
}

func (me *SubCommand) registerNewOption(option optioner) {
	me.options = append(me.options, option)
}

func (me *SubCommand) optionsForNames() (map[string]optioner,
	map[string]optioner) {
	optionForLongName := make(map[string]optioner, len(me.options))
	optionForShortName := make(map[string]optioner, len(me.options))
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
