// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

import (
	"fmt"
	"strings"
)

type Option struct {
	LongName     string
	ShortName    rune
	Help         string
	Required     bool
	ValueCount   ValueCount
	VarName      string // e.g., -o|--outfile FILENAME
	DefaultValue any
	Value        any
	ValueType    ValueType
	Validator    Validator
}

func newOption(name, help string) *Option {
	shortName, longName := namesForName(name)
	return &Option{LongName: longName, ShortName: shortName, Help: help}
}

func (me *Option) AsBool() bool {
	if me.ValueType != Flag {
		panic(fmt.Sprintf("AsBool() called on type %s", me.ValueType))
	}
	if me.Value == nil {
		if me.DefaultValue == nil {
			return false
		}
		return me.DefaultValue.(bool)
	}
	return me.Value.(bool)
}

// TODO AsT...

func (me *Option) Debug(indent int) {
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
