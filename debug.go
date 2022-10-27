// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

import (
	"fmt"
	"sort"
	"strings"
)

func (me *Parser) debug() {
	fmt.Println("Parser")
	fmt.Printf("    %v v%v\n", me.AppName, me.AppVersion)
	fmt.Printf("    QuitOnError=%t\n", me.QuitOnError)
	keys := me.getSubCommands()
	sort.Strings(keys)
	for _, name := range keys {
		subcommand := me.SubCommands[name]
		if name == mainSubCommand {
			name = "<MAIN>"
		}
		fmt.Printf("    SubCommand=%v\n", name)
		subcommand.debug(8)
	}
	fmt.Printf("    PositionalCount=%s\n", me.PositionalCount)
	fmt.Printf("    PositionalVarName=%s\n", me.PositionalVarName)
	fmt.Printf("    Positionals=%v\n", me.Positionals)
}

func (me *Option) debug(indent int) {
	tab := strings.Repeat(" ", indent)
	fmt.Printf("%sLongName=%s\n", tab, me.LongName)
	fmt.Printf("%sShortName=%s\n", tab, string(me.ShortName))
	fmt.Printf("%sHelp=%s\n", tab, me.Help)
	fmt.Printf("%sRequired=%t\n", tab, me.Required)
	fmt.Printf("%sValueCount=%s\n", tab, me.ValueCount)
	fmt.Printf("%sVarName=%s\n", tab, me.VarName)
	fmt.Printf("%sDefaultValue=%v\n", tab, me.DefaultValue)
	fmt.Printf("%sValue=%v\n", tab, me.Value)
	fmt.Printf("%sValueType=%s\n", tab, me.ValueType)
	if me.Validator == nil {
		fmt.Printf("%sValidator=nil\n", tab)
	} else {
		fmt.Printf("%sValidator=func(any) bool\n", tab)
	}
}

func (me *SubCommand) debug(indent int) {
	tab := strings.Repeat(" ", indent)
	fmt.Printf("%sLongName=%s\n", tab, me.LongName)
	fmt.Printf("%sShortName=%s\n", tab, string(me.ShortName))
	fmt.Printf("%sHelp=%s\n", tab, me.Help)
	for i, option := range me.Options {
		fmt.Printf("%sOption #%d:\n", tab, i)
		option.debug(indent + 4)
	}
}
