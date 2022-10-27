// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

import "fmt"

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

func newOption(name, help string, valueType ValueType) *Option {
	shortName, longName := namesForName(name)
	return &Option{LongName: longName, ShortName: shortName, Help: help,
		ValueCount: One, ValueType: valueType}
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

func (me *Option) AsInt() int {
	if me.ValueType != Int {
		panic(fmt.Sprintf("AsInt() called on type %s", me.ValueType))
	}
	if me.Value == nil {
		if me.DefaultValue == nil {
			return 0
		}
		return me.DefaultValue.(int)
	}
	return me.Value.(int)
}

func (me *Option) AsReal() float64 {
	if me.ValueType != Real {
		panic(fmt.Sprintf("AsReal() called on type %s", me.ValueType))
	}
	if me.Value == nil {
		if me.DefaultValue == nil {
			return 0.0
		}
		return me.DefaultValue.(float64)
	}
	return me.Value.(float64)
}

func (me *Option) AsStr() string {
	if me.ValueType != Str {
		panic(fmt.Sprintf("AsStr() called on type %s", me.ValueType))
	}
	if me.Value == nil {
		if me.DefaultValue == nil {
			return ""
		}
		return me.DefaultValue.(string)
	}
	return me.Value.(string)
}

func (me *Option) AsStrs() []string {
	if me.ValueType != Strs {
		panic(fmt.Sprintf("AsStrs() called on type %s", me.ValueType))
	}
	if me.Value == nil {
		if me.DefaultValue == nil {
			return []string{}
		}
		return me.DefaultValue.([]string)
	}
	return me.Value.([]string)
}
