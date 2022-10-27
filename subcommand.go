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
	return &SubCommand{LongName: "", ShortName: NoShortName, Help: "",
		Options: make([]*Option, 0)}
}

func newSubCommand(name, help string) *SubCommand {
	return &SubCommand{LongName: name, ShortName: NoShortName, Help: help,
		Options: make([]*Option, 0)}
}

func (me *SubCommand) Bool(name, help string) *Option {
	option := newOption(name, help)
	option.ValueCount = ZeroOrOne
	option.ValueType = Bool
	me.Options = append(me.Options, option)
	return option
}

// TODO Int() Real() Str() Strs() Choice()

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
