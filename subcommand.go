// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

type SubCommand struct {
	longName          string
	shortName         rune
	help              string
	options           []optioner
	firstDelayedError string
}

func newMainSubCommand() *SubCommand {
	return &SubCommand{longName: "", shortName: noShortName, help: "",
		options: make([]optioner, 0)}
}

// Always returns a *SubCommand; _and_ either nil or error
func newSubCommand(name, help string) (*SubCommand, error) {
	name, err := validatedName(name, "subcommand")
	return &SubCommand{longName: name, shortName: noShortName, help: help,
		options: make([]optioner, 0)}, err
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
	option, err := newFlagOption(name, help)
	me.registerNewOption(option, err)
	return option
}

func (me *SubCommand) Int(name, help string, theDefault int) *IntOption {
	option, err := newIntOption(name, help, theDefault)
	me.registerNewOption(option, err)
	return option
}

func (me *SubCommand) IntInRange(name, help string, minimum, maximum,
	theDefault int) *IntOption {
	option, err := newIntOption(name, help, theDefault)
	option.Validator = makeIntRangeValidator(minimum, maximum)
	me.registerNewOption(option, err)
	return option
}

func (me *SubCommand) Real(name, help string,
	theDefault float64) *RealOption {
	option, err := newRealOption(name, help, theDefault)
	me.registerNewOption(option, err)
	return option
}

func (me *SubCommand) RealInRange(name, help string, minimum, maximum,
	theDefault float64) *RealOption {
	option, err := newRealOption(name, help, theDefault)
	option.Validator = makeRealRangeValidator(minimum, maximum)
	me.registerNewOption(option, err)
	return option
}

func (me *SubCommand) Str(name, help, theDefault string) *StrOption {
	option, err := newStrOption(name, help, theDefault)
	me.registerNewOption(option, err)
	return option
}

func (me *SubCommand) Choice(name, help string, choices []string,
	theDefault string) *StrOption {
	option, err := newStrOption(name, help, theDefault)
	option.Validator = makeChoiceValidator(choices)
	me.registerNewOption(option, err)
	return option
}

func (me *SubCommand) Strs(name, help string) *StrsOption {
	option, err := newStrsOption(name, help)
	me.registerNewOption(option, err)
	return option
}

func (me *SubCommand) Ints(name, help string) *IntsOption {
	option, err := newIntsOption(name, help)
	me.registerNewOption(option, err)
	return option
}

func (me *SubCommand) Reals(name, help string) *RealsOption {
	option, err := newRealsOption(name, help)
	me.registerNewOption(option, err)
	return option
}

func (me *SubCommand) registerNewOption(option optioner, err error) {
	me.options = append(me.options, option)
	if err != nil && me.firstDelayedError == "" {
		me.firstDelayedError = err.Error()
	}
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
