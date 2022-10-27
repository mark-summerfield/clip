// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

import (
	"fmt"
	"strings"
)

type SubCommand struct {
	LongName  string
	ShortName rune
	Help      string
	Options   []*Option
}

func newMainSubCommand() *SubCommand {
	return &SubCommand{LongName: "", ShortName: noShortName, Help: "",
		Options: make([]*Option, 0)}
}

func newSubCommand(name, help string) *SubCommand {
	return &SubCommand{LongName: name, ShortName: noShortName, Help: help,
		Options: make([]*Option, 0)}
}

func (me *SubCommand) Flag(name, help string) *Option {
	option := me.newOption(name, help, Flag)
	option.ValueCount = Zero
	return option
}

func (me *SubCommand) Int(name, help string) *Option {
	return me.newOption(name, help, Int)
}

func (me *SubCommand) IntInRange(name, help string,
	minimum, maximum int) *Option {
	option := me.newOption(name, help, Int)
	option.Validator = makeIntRangeValidator(minimum, maximum)
	return option
}

func (me *SubCommand) Real(name, help string) *Option {
	return me.newOption(name, help, Real)
}

func (me *SubCommand) RealInRange(name, help string,
	minimum, maximum float64) *Option {
	option := me.newOption(name, help, Real)
	option.Validator = makeRealRangeValidator(minimum, maximum)
	return option
}

func (me *SubCommand) Str(name, help string) *Option {
	return me.newOption(name, help, Str)
}

func (me *SubCommand) Choice(name, help string, choices []string) *Option {
	option := me.newOption(name, help, Str)
	option.Validator = makeChoiceValidator(choices)
	return option
}

func (me *SubCommand) Strs(name, help string) *Option {
	option := me.newOption(name, help, Strs)
	option.ValueCount = OneOrMore
	return option
}

func (me *SubCommand) newOption(name, help string,
	valueType ValueType) *Option {
	option := newOption(name, help, valueType)
	me.Options = append(me.Options, option)
	return option
}

func (me *SubCommand) optionsForNames() (map[string]*Option,
	map[string]*Option) {
	optionForLongName := make(map[string]*Option, len(me.Options))
	optionForShortName := make(map[string]*Option, len(me.Options))
	for _, option := range me.Options {
		if option.LongName != "" {
			optionForLongName[option.LongName] = option
		}
		if option.ShortName != 0 {
			optionForShortName[string(option.ShortName)] = option
		}
	}
	return optionForLongName, optionForShortName
}

func (me *SubCommand) Debug(indent int) {
	tab := strings.Repeat(" ", indent)
	fmt.Printf("%sLongName=%s\n", tab, me.LongName)
	fmt.Printf("%sShortName=%s\n", tab, string(me.ShortName))
	fmt.Printf("%sHelp=%s\n", tab, me.Help)
	for i, option := range me.Options {
		fmt.Printf("%sOption #%d:\n", tab, i)
		option.Debug(indent + 4)
	}
}
