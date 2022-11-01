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

func (me *SubCommand) Flag(name, help string) (*FlagOption, error) {
	option, err := newFlagOption(name, help)
	if err != nil {
		return nil, err
	}
	me.registerNewOption(option)
	return option, nil
}

func (me *SubCommand) Int(name, help string, theDefault int) (*IntOption,
	error) {
	option, err := newIntOption(name, help, theDefault)
	if err != nil {
		return nil, err
	}
	me.registerNewOption(option)
	return option, nil
}

func (me *SubCommand) IntInRange(name, help string,
	minimum, maximum, theDefault int) (*IntOption, error) {
	option, err := newIntOption(name, help, theDefault)
	if err != nil {
		return nil, err
	}
	option.Validator = makeIntRangeValidator(minimum, maximum)
	me.registerNewOption(option)
	return option, nil
}

func (me *SubCommand) Real(name, help string,
	theDefault float64) (*RealOption, error) {
	option, err := newRealOption(name, help, theDefault)
	if err != nil {
		return nil, err
	}
	me.registerNewOption(option)
	return option, nil
}

func (me *SubCommand) RealInRange(name, help string,
	minimum, maximum, theDefault float64) (*RealOption, error) {
	option, err := newRealOption(name, help, theDefault)
	if err != nil {
		return nil, err
	}
	option.Validator = makeRealRangeValidator(minimum, maximum)
	me.registerNewOption(option)
	return option, nil
}

func (me *SubCommand) Str(name, help, theDefault string) (*StrOption,
	error) {
	option, err := newStrOption(name, help, theDefault)
	if err != nil {
		return nil, err
	}
	me.registerNewOption(option)
	return option, nil
}

func (me *SubCommand) Choice(name, help string, choices []string,
	theDefault string) (*StrOption, error) {
	option, err := newStrOption(name, help, theDefault)
	if err != nil {
		return nil, err
	}
	option.Validator = makeChoiceValidator(choices)
	me.registerNewOption(option)
	return option, nil
}

func (me *SubCommand) Strs(name, help string) (*StrsOption, error) {
	option, err := newStrsOption(name, help)
	if err != nil {
		return nil, err
	}
	me.registerNewOption(option)
	return option, nil
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
