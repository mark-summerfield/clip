// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

import (
	"fmt"
	"strconv"
)

type Option struct {
	LongName     string
	ShortName    rune
	Help         string
	Required     bool
	ValueCount   ValueCount
	VarName      string // e.g., -o|--outfile FILENAME
	DefaultValue any    // not valid if ValueType == Strs
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

func (me *Option) Size() int {
	if me.Value == nil || me.ValueType == Flag {
		return 0
	}
	if me.ValueType == Strs {
		return len(me.Value.([]string))
	}
	return 1
}

func (me *Option) AddValue(value string) error {
	if me.Validator != nil {
		if !me.Validator(value) {
			return fmt.Errorf("invalid value for %s: %s", me.LongName,
				value)
		}
	}
	switch me.ValueType {
	case Flag:
		return fmt.Errorf("flag %s got unexpected value %s", me.LongName,
			value)
	case Int:
		i, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("option %s expected an int value, got %s",
				me.LongName, value)
		}
		me.Value = i
	case Real:
		r, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("option %s expected a real value, got %s",
				me.LongName, value)
		}
		me.Value = r
	case Str:
		me.Value = value
	case Strs:
		if me.Value == nil {
			me.Value = make([]string, 0, 1)
		}
		me.Value = append(me.Value.([]string), value)
	default:
		panic("invalid ValueType #2")
	}
	return nil
}
